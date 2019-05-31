package main

import (
	"fmt"
	"os"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	var err error
	pixelgl.Run(func() {
		cmd := NewApp()
		err = cmd.Execute()
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}
