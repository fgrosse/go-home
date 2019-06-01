package main

import (
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/fgrosse/go-home/assets"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

type Render struct {
	Width, Height     float64
	CheckIn           time.Time
	CheckOut          time.Time
	EOD               time.Time
	Atlas             *text.Atlas
	MarkerColor       color.Color
	BorderColor       color.Color
	ShowRemainingTime bool          // show how much time is left instead of the current time
	timeShift         time.Duration // for debugging
}

func NewRender(conf Config) (*Render, error) {
	fnt, err := loadTTF(conf.Debug, "/assets/GlacialIndifference-Regular.ttf", 12)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load font")
	}

	return &Render{
		Width:             float64(conf.UI.WindowWidth),
		Height:            float64(conf.UI.WindowHeight),
		CheckIn:           conf.CheckIn,
		CheckOut:          conf.CheckOut,
		EOD:               conf.EndOfDay,
		Atlas:             text.NewAtlas(fnt, text.ASCII),
		MarkerColor:       colornames.Mediumblue,
		ShowRemainingTime: conf.UI.ShowRemainingTime,
	}, nil
}

func loadTTF(debug bool, path string, size float64) (font.Face, error) {
	bytes := assets.FSMustByte(debug, path)
	fnt, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(fnt, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func (r *Render) Draw(t pixel.Target) {
	now := time.Now().Add(r.timeShift)
	progress := r.progress(now)

	markerTxt := r.markerText(progress, now)
	checkoutTxt := r.checkoutText(now)

	r.drawGradient(t, progress)
	r.drawCells(t, markerTxt, checkoutTxt)
	r.drawTargetMarker(t, progress, markerTxt)
	r.drawCurrentMarker(t, progress, now, markerTxt)
	r.drawText(t, now, markerTxt, checkoutTxt)
	r.drawRectangle(t, now)
}

func (r *Render) drawGradient(t pixel.Target, progress float64) {
	leftColor := pixel.RGB(0, 1, 0)
	rightColor := pixel.RGB(progress, 1-progress, 0)

	rect := imdraw.New(nil)
	rect.Color = leftColor
	rect.Push(pixel.V(0, 0))
	rect.Push(pixel.V(0, r.Height))

	rect.Color = rightColor
	rect.Push(pixel.V(r.Width*progress, r.Height))
	rect.Push(pixel.V(r.Width*progress, 0))
	rect.Polygon(0)
	rect.Draw(t)

	// draw gray scale gradient
	rect = imdraw.New(nil)
	gray := pixel.RGB(0.333, 0.333, 0.333)
	rect.Color = gray.Scaled(progress)
	rect.Push(pixel.V(r.Width*progress+1, 0))
	rect.Push(pixel.V(r.Width*progress+1, r.Height))

	rect.Color = gray
	rect.Push(pixel.V(r.Width, r.Height))
	rect.Push(pixel.V(r.Width, 0))
	rect.Polygon(0)
	rect.Draw(t)
}

func (r *Render) drawCells(t pixel.Target, markerTxt, checkoutTxt *text.Text) {
	numCells := 9
	txtBounds := markerTxt.Bounds()
	checkoutBounds := checkoutTxt.Bounds()
	for i := 0; i < numCells; i++ {
		x := r.Width * float64(i) / float64(numCells)
		line := imdraw.New(nil)
		line.Color = color.Black

		v := pixel.V(x, r.Height/2)
		if txtBounds.Contains(v) || checkoutBounds.Contains(v) {
			col := pixel.ToRGBA(line.Color)
			col.A = 0.1
			line.Color = col
		}

		line.Push(pixel.V(x, 0))
		line.Push(pixel.V(x, r.Height))
		line.Line(1)
		line.Draw(t)
	}
}

func (r *Render) drawTargetMarker(t pixel.Target, progress float64, markerTxt *text.Text) {
	posX := r.position(r.CheckOut)

	imd := imdraw.New(nil)
	col := pixel.RGB(0, 0, 0)
	if markerTxt.Bounds().Contains(pixel.V(posX, r.Height/2)) {
		col.A = 0.1
	}

	imd.Color = col
	imd.Push(pixel.V(posX, 2))
	imd.Push(pixel.V(posX, r.Height-2))
	imd.Line(3)
	imd.Draw(t)
}

func (r *Render) drawCurrentMarker(t pixel.Target, progress float64, now time.Time, txt *text.Text) {
	posX := progress * r.Width
	col := r.MarkerColor
	if now.After(r.CheckOut) {
		col = r.BorderColor
	}

	imd := imdraw.New(nil)
	imd.Color = col
	imd.Push(pixel.V(posX, 2))
	imd.Push(pixel.V(posX, r.Height-2))
	imd.Line(2)
	imd.Draw(t)

	imd = imdraw.New(nil)
	imd.Color = col
	imd.Push(pixel.V(posX-5, 2))
	imd.Push(pixel.V(posX+1+5, 2))
	imd.Push(pixel.V(posX+1, 5+2))
	imd.Push(pixel.V(posX, 5+2))
	imd.Polygon(0)
	imd.Draw(t)

	imd = imdraw.New(nil)
	imd.Color = col
	imd.Push(pixel.V(posX-5, r.Height-2))
	imd.Push(pixel.V(posX+1+5, r.Height-2))
	imd.Push(pixel.V(posX+1, r.Height-5-2))
	imd.Push(pixel.V(posX, r.Height-5-2))
	imd.Polygon(0)
	imd.Draw(t)
}

func (r *Render) markerText(progress float64, now time.Time) *text.Text {
	posX := progress * r.Width
	txt := text.New(pixel.V(posX+5, 4), r.Atlas)
	if progress > 0.7 {
		txt.Color = color.White
	} else {
		txt.Color = color.Black
	}

	if r.ShowRemainingTime {
		remaining := r.CheckOut.Sub(now).Round(time.Minute)
		s := remaining.String()
		s = s[:len(s)-2] // strip away trailing seconds
		txt.WriteString(s)
	} else {
		txt.WriteString(now.Format("15:04"))
	}

	return txt
}

func (r *Render) checkoutText(now time.Time) *text.Text {
	txt := text.New(pixel.V(4+r.position(r.CheckOut), 4), r.Atlas)
	if now.Before(r.CheckOut) {
		txt.Color = color.White
	} else {
		txt.Color = color.Black
	}

	txt.WriteString(r.CheckOut.Format("15:04"))
	return txt
}

func (r *Render) drawRectangle(t pixel.Target, now time.Time) {
	var borderWidth float64 = 2

	rect := imdraw.New(nil)
	if now.Before(r.CheckOut) {
		r.BorderColor = color.Black
	} else {
		// Pulsating red border
		totalMillis := int(float64(now.UnixNano()) / float64(time.Millisecond))
		scale := (1 + math.Sin(float64(totalMillis)/500)) / 2
		r.BorderColor = pixel.ToRGBA(colornames.Red).Mul(pixel.RGB(scale, scale, scale))
		borderWidth = 4
	}

	rect.Color = r.BorderColor
	rect.EndShape = imdraw.RoundEndShape
	rect.Push(pixel.V(1, 1))
	rect.Push(pixel.V(1, r.Height-1))
	rect.Push(pixel.V(r.Width-1, r.Height-1))
	rect.Push(pixel.V(r.Width-1, 1))
	rect.Push(pixel.V(1, 1))
	rect.Line(borderWidth)
	rect.Draw(t)
}

func (r *Render) drawText(t pixel.Target, now time.Time, markerTxt, checkoutTxt *text.Text) {
	txt := text.New(pixel.V(4, 4), r.Atlas)
	txt.Color = color.Black
	txt.WriteString(r.CheckIn.Format("15:04"))
	bounds := txt.Bounds()

	shift := r.Height/2 - bounds.Size().Y/2 - 1
	m := pixel.IM.Moved(pixel.V(0, shift))

	if bounds.Intersect(markerTxt.Bounds()) == pixel.ZR {
		txt.Draw(t, m)
	}

	if checkoutTxt.Bounds().Intersect(markerTxt.Bounds()) == pixel.ZR {
		checkoutTxt.Draw(t, m)
	}

	markerTxt.Draw(t, m)
}

func (r *Render) position(t time.Time) float64 {
	return r.Width * r.progress(t)
}

func (r *Render) progress(now time.Time) float64 {
	left := float64(r.CheckIn.Unix())
	right := float64(r.EOD.Unix())
	nowUnix := float64(now.Unix())

	totalSec := right - left
	diffSec := nowUnix - left
	return diffSec / totalSec
}
