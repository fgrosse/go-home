package main

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func (app *App) handleInput(win *pixelgl.Window, dt float64) {
	const speed = 2 * float64(time.Hour)

	switch {
	case win.Pressed(pixelgl.KeyRight):
		app.render.timeShift += time.Duration(speed * dt)
	case win.Pressed(pixelgl.KeyLeft):
		app.render.timeShift -= time.Duration(speed * dt)
	}
}
