package main

import (
	"image/color"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
)

type App struct {
	win *pixelgl.Window
	render *Render
}

func NewApp(conf Config) (*App, error) {
	today := time.Now()
	year, month, day := today.Date()
	checkIn := time.Date(year, month, day, 9, 30, 0, 0, time.Local)
	checkOut := checkIn.Add(conf.WorkDuration).Add(conf.LunchDuration)
	endOfDay := conf.DayEnd.Time(today)

	win, err := newWindow(conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create window")
	}

	r, err := NewRender(conf, checkIn, checkOut, endOfDay)
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		win: win,
		render: r,
	}, nil
}

func (app *App) Run() {
		fps := time.Tick(time.Second / 30)
		last := time.Now()
		for !app.win.Closed() {
			dt := time.Since(last).Seconds()
			last = time.Now()

			app.win.Clear(color.White)
			HandleInput(app.win, dt)
			app.render.Draw(app.win)
			app.win.Update()

			<-fps
		}
}

func newWindow(conf Config) (*pixelgl.Window, error) {
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
		return nil, errors.WithStack(err)
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
