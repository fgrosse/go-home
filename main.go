package main

import (
	"image/color"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	var (
		width float64 = 800
		height float64 = 32
		moveSpeed float64 = 100
		checkIn = time.Now()
		checkOut = checkIn.Add(30*time.Minute)
	)

	win, err := NewWindow(width, height)
	if err != nil {
		log.Fatal(err)
	}

	// winCenter := win.Bounds().Center()
	fps := time.Tick(time.Second / 30)

	rect := Rect(width, height)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		winPos := win.GetPos()
		switch {
		case win.Pressed(pixelgl.KeyRight):
			winPos.X += moveSpeed*dt
			win.SetPos(winPos)
		case win.Pressed(pixelgl.KeyLeft):
			winPos.X -= moveSpeed*dt
			win.SetPos(winPos)
		case win.Pressed(pixelgl.KeyDown):
			winPos.Y += moveSpeed*dt
			win.SetPos(winPos)
		case win.Pressed(pixelgl.KeyUp):
			winPos.Y -= moveSpeed*dt
			win.SetPos(winPos)
		}

		win.Clear(color.White)
		rect.Draw(win)
		DrawMarker(win, checkIn, checkOut)
		win.Update()

		// TODO
		<-fps
	}
}

func NewWindow(width, height float64) (*pixelgl.Window, error) {
	cfg := pixelgl.WindowConfig{
		Title:  "Go Home",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
		Undecorated: true,
		Resizable: true,
	}

	// TODO: we need GLFW 3.3 where we get the GLFW_TRANSPARENT_FRAMEBUFFER option
	// See: https://www.glfw.org/docs/3.3/window_guide.html#window_hints_wnd

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create window")
	}

	win.SetSmooth(true)

	var displayWidth, displayHeight float64 = 1920, 1080
	pos := pixel.Vec{
		X: displayWidth/2-width/2,
		Y: displayHeight-height-5,
	}

	win.SetPos(pos)
	win.Update()

	return win, nil
}

func Rect(width, height float64) *imdraw.IMDraw {
	rect := imdraw.New(nil)
	rect.Color = pixel.RGB(0, 1, 0)
	rect.Push(pixel.V(0, 0))
	rect.Push(pixel.V(0, height))

	rect.Color = pixel.RGB(1, 0, 0)
	rect.Push(pixel.V(width, height))
	rect.Push(pixel.V(width, 0))
	rect.Polygon(0)

	return rect
}

func DrawMarker(win *pixelgl.Window, checkIn, checkOut time.Time) {
	var (
		in  = float64(checkIn.Unix())
		out = float64(checkOut.Unix())
		now = float64(time.Now().Unix())
	)

	totalSec := out - in
	diffSec := now -in
	d := diffSec/totalSec

	bounds := win.Bounds()
	posX := bounds.Max.X*d

	imd := imdraw.New(nil)
	imd.Color = color.Black
	imd.Push(pixel.V(posX, 0))
	imd.Push(pixel.V(posX, bounds.Max.Y))
	imd.Line(1)
	imd.Draw(win)

	imd = imdraw.New(nil)
	imd.Color = color.Black
	imd.Push(pixel.V(posX-5, 0))
	imd.Push(pixel.V(posX+5, 0))
	imd.Push(pixel.V(posX, 5))
	imd.Polygon(0)
	imd.Draw(win)

	imd = imdraw.New(nil)
	imd.Color = color.Black
	imd.Push(pixel.V(posX-5, bounds.Max.Y))
	imd.Push(pixel.V(posX+5, bounds.Max.Y))
	imd.Push(pixel.V(posX, bounds.Max.Y-5))
	imd.Polygon(0)
	imd.Draw(win)
}
