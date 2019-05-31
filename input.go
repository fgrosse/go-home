package main

import (
	"github.com/faiface/pixel/pixelgl"
)

func (app *App) handleInput(win *pixelgl.Window, dt float64) {
	const moveSpeed float64 = 100

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
}
