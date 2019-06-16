package main

import (
	"image/color"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type App struct {
	*cobra.Command
	logger *zap.Logger
	conf   Config
	win    *pixelgl.Window
	render *Render

	initErr  error
	shutdown bool
	limitFPS bool
}

func NewApp() *App {
	app := &App{Command: &cobra.Command{
		Use: "go-home",
	}}

	app.SilenceUsage = true  // do not output usage in case of an error
	app.SilenceErrors = true // we log them manually in the main function
	app.Command.RunE = app.Run

	var (
		debug  bool
		config string
	)

	flags := app.PersistentFlags()
	flags.StringVar(&config, "config", os.ExpandEnv("$HOME/.go-home.yml"), "config file")
	flags.BoolVar(&debug, "debug", false, "enable debug mode")

	cobra.OnInitialize(app.loadConfig(&debug, &config))

	return app
}

func (app *App) Run(_ *cobra.Command, _ []string) error {
	if app.initErr != nil {
		return app.initErr
	}

	var err error
	app.logger.Info("Starting application", zap.Object("config", app.conf))
	err = app.createWindow(app.conf.UI)
	if err != nil {
		return errors.Wrap(err, "failed to create window")
	}

	app.render, err = NewRender(app.conf)
	if err != nil {
		return errors.Wrap(err, "failed to create renderer")
	}

	app.runLoop()

	err = app.conf.Save()
	if err != nil {
		app.logger.Error("Failed to save configuration on shutdown", zap.Error(err))
	}

	return nil
}

func (app *App) createWindow(conf UIConfig) error {
	width := float64(conf.WindowWidth)
	height := float64(conf.WindowHeight)

	cfg := pixelgl.WindowConfig{
		Title:       "Go Home",
		Bounds:      pixel.R(0, 0, width, height),
		VSync:       true,
		Undecorated: true,
		Resizable:   false,
		AlwaysOnTop: true,
	}

	// TODO: we need GLFW 3.3 where we get the GLFW_TRANSPARENT_FRAMEBUFFER option
	// See: https://www.glfw.org/docs/3.3/window_guide.html#window_hints_wnd

	var err error
	app.win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		return errors.WithStack(err)
	}

	app.win.SetSmooth(true)
	app.win.SetPos(app.conf.UI.WindowPos)
	app.win.Update()

	return nil
}

func (app *App) runLoop() {
	app.limitFPS = true // might be disabled when moving the window via Shift+Arrow
	fps := time.Tick(time.Second / time.Duration(app.conf.UI.FPS))
	last := time.Now()
	for !app.win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// TODO: deal with a new dawn

		app.win.Clear(color.White)
		app.handleInput(app.win, dt)
		app.render.Draw(app.win)
		app.win.Update()

		if app.shutdown {
			return
		}

		if app.limitFPS {
			<-fps
		}
	}
}
