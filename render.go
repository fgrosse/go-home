package main

import (
	"image/color"
	"io/ioutil"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

type Render struct {
	Width, Height float64
	CheckIn       time.Time
	CheckOut      time.Time
	EOD           time.Time
	Atlas         *text.Atlas
	MarkerColor   color.Color
}

func NewRender(conf UIConfig, checkIn, checkOut, endOfDay time.Time) (*Render, error) {
	fnt, err := loadTTF("assets/GlacialIndifference-Regular.ttf", 10)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load font")
	}

	return &Render{
		Width:       float64(conf.WindowWidth),
		Height:      float64(conf.WindowHeight),
		CheckIn:     checkIn,
		CheckOut:    checkOut,
		EOD:         endOfDay,
		Atlas:       text.NewAtlas(fnt, text.ASCII),
		MarkerColor: colornames.Mediumblue,
	}, nil
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

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
	now := time.Now()

	r.drawGradient(t)
	r.drawCells(t)
	r.drawText(t)
	r.drawMarker(t, now)
	r.drawRectangle(t)
}

func (r *Render) drawGradient(t pixel.Target) {
	rect := imdraw.New(nil)
	rect.Color = pixel.RGB(0, 1, 0)
	rect.Push(pixel.V(0, 0))
	rect.Push(pixel.V(0, r.Height))

	rect.Color = pixel.RGB(1, 0, 0)
	rect.Push(pixel.V(r.Width, r.Height))
	rect.Push(pixel.V(r.Width, 0))
	rect.Polygon(0)
	rect.Draw(t)
}

func (r *Render) drawCells(t pixel.Target) {
	numCells := 9
	for i := 0; i < numCells; i++ {
		x := r.Width * float64(i) / float64(numCells)
		line := imdraw.New(nil)
		line.Color = color.Black
		line.Push(pixel.V(x, 0))
		line.Push(pixel.V(x, r.Height))
		line.Line(1)
		line.Draw(t)
	}
}

func (r *Render) drawMarker(t pixel.Target, now time.Time) {
	posX := r.position(now)

	imd := imdraw.New(nil)
	imd.Color = r.MarkerColor
	imd.Push(pixel.V(posX, 2))
	imd.Push(pixel.V(posX, r.Height-2))
	imd.Line(2)
	imd.Draw(t)

	imd = imdraw.New(nil)
	imd.Color = r.MarkerColor
	imd.Push(pixel.V(posX-5, 2))
	imd.Push(pixel.V(posX+1+5, 2))
	imd.Push(pixel.V(posX+1, 5+2))
	imd.Push(pixel.V(posX, 5+2))
	imd.Polygon(0)
	imd.Draw(t)

	imd = imdraw.New(nil)
	imd.Color = r.MarkerColor
	imd.Push(pixel.V(posX-5, r.Height-2))
	imd.Push(pixel.V(posX+1+5, r.Height-2))
	imd.Push(pixel.V(posX+1, r.Height-5-2))
	imd.Push(pixel.V(posX, r.Height-5-2))
	imd.Polygon(0)
	imd.Draw(t)

	txt := text.New(pixel.V(posX+5, 4), r.Atlas)
	txt.Color = r.MarkerColor
	txt.WriteString(now.Format("15:04"))
	txt.Draw(t, pixel.IM)
}

func (r *Render) drawRectangle(t pixel.Target) {
	var borderWidth float64 = 2

	rect := imdraw.New(nil)
	rect.Color = color.Black
	rect.EndShape = imdraw.RoundEndShape
	rect.Push(pixel.V(1, 1))
	rect.Push(pixel.V(1, r.Height-1))
	rect.Push(pixel.V(r.Width-1, r.Height-1))
	rect.Push(pixel.V(r.Width-1, 1))
	rect.Push(pixel.V(1, 1))
	rect.Line(borderWidth)
	rect.Draw(t)
}

func (r *Render) drawText(t pixel.Target) {
	txt := text.New(pixel.V(4, 4), r.Atlas)
	txt.Color = color.Black
	txt.WriteString(r.CheckIn.Format("15:04"))
	txt.Draw(t, pixel.IM)

	txt = text.New(pixel.V(4+r.position(r.CheckOut), 4), r.Atlas)
	txt.Color = color.Black
	txt.WriteString(r.CheckOut.Format("15:04"))
	txt.Draw(t, pixel.IM)
}

func (r *Render) position(t time.Time) float64 {
	in := float64(r.CheckIn.Unix())
	out := float64(r.EOD.Unix())
	nowUnix := float64(t.Unix())

	totalSec := out - in
	diffSec := nowUnix - in
	d := diffSec / totalSec
	return r.Width * d
}
