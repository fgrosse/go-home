package main

import (
	"log"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	conf := Config{
		WorkDuration:  8 * time.Hour,
		LunchDuration: 1 * time.Hour,
		DayEnd:        ClockTime{Hour: 20, Minute: 00},

		WindowWidth:  512,
		WindowHeight: 32,
		VSync:        true,
	}

	pixelgl.Run(func() {
		app, err := NewApp(conf)
		if err != nil {
			log.Fatal(err)
		}

		app.Run()
	})
}
