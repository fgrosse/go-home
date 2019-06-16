package main

import (
	"math"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func (app *App) handleInput(win *pixelgl.Window, dt float64) {
	if win.Pressed(pixelgl.KeyLeftShift) {
		app.limitFPS = false

		const speed = 100.0 // pixel per second
		delta := speed * dt
		if win.Pressed(pixelgl.KeyRight) {
			app.conf.UI.WindowPos.X += delta
		} else if win.Pressed(pixelgl.KeyLeft) {
			app.conf.UI.WindowPos.X -= delta
		}

		if win.Pressed(pixelgl.KeyUp) {
			app.conf.UI.WindowPos.Y -= delta
		} else if win.Pressed(pixelgl.KeyDown) {
			app.conf.UI.WindowPos.Y += delta
		}

		currPos := app.win.GetPos()
		if math.Round(currPos.X) != math.Round(app.conf.UI.WindowPos.X) ||
			math.Round(currPos.Y) != math.Round(app.conf.UI.WindowPos.Y) {
			app.win.SetPos(app.conf.UI.WindowPos)
		}

	} else {
		app.limitFPS = true

		if app.conf.Debug {
			const speed = 2 * float64(time.Hour)
			switch {
			case win.Pressed(pixelgl.KeyRight):
				app.render.timeShift += time.Duration(speed * dt)
			case win.Pressed(pixelgl.KeyLeft):
				app.render.timeShift -= time.Duration(speed * dt)
			case win.Pressed(pixelgl.KeyEscape):
				app.shutdown = true
			}
		}
	}
}
