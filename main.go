package main

import (
	"image/color"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	conf := Config{
		WorkDuration:  8 * time.Hour,
		LunchDuration: 1 * time.Hour,
		DayEnd:        ClockTime{Hour: 20, Minute: 00},

		WindowWidth:  512,
		WindowHeight: 32,
		VSync:        true,
	}

	today := time.Now()
	year, month, day := today.Date()
	checkIn := time.Date(year, month, day, 9, 30, 0, 0, time.Local)
	checkOut := checkIn.Add(conf.WorkDuration).Add(conf.LunchDuration)
	endOfDay := conf.DayEnd.Time(today)

	win, err := NewWindow(conf)
	if err != nil {
		log.Fatal(err)
	}

	r, err := NewRender(conf, checkIn, checkOut, endOfDay)
	if err != nil {
		log.Fatal(err)
	}

	fps := time.Tick(time.Second / 30)
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(color.White)
		HandleInput(win, dt)
		r.Draw(win)
		win.Update()

		<-fps
	}
}

func NewWindow(conf Config) (*pixelgl.Window, error) {
	width := float64(conf.WindowWidth)
	height := float64(conf.WindowHeight)

	cfg := pixelgl.WindowConfig{
		Title:       "Go Home",
		Bounds:      pixel.R(0, 0, width, height),
		VSync:       conf.VSync,
		Undecorated: true,
		Resizable:   false,
		Floating:    true,
		AutoIconify: true,
	}

	// TODO: we need GLFW 3.3 where we get the GLFW_TRANSPARENT_FRAMEBUFFER option
	// See: https://www.glfw.org/docs/3.3/window_guide.html#window_hints_wnd

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create window")
	}

	win.SetSmooth(true)

	var displayWidth float64 = 1920 // TODO: make dynamic
	pos := pixel.Vec{
		X: displayWidth/2 - width/2,
		Y: 29,
	}

	win.SetPos(pos)
	win.Update()

	return win, nil
}
