package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CheckIn  time.Time `yaml:"check_in"`
	CheckOut time.Time `yaml:"-"`
	EndOfDay time.Time `yaml:"-"`

	WorkDuration  time.Duration `yaml:"work_duration"`
	LunchDuration time.Duration `yaml:"lunch_duration"`
	DayEnd        ClockTime     `yaml:"day_end"`

	UI    UIConfig `yaml:"ui"`
	Debug bool     `yaml:"-"`

	path string `yaml:"-"`
}

type UIConfig struct {
	FPS               int       `yaml:"fps"`
	WindowWidth       int       `yaml:"width"`
	WindowHeight      int       `yaml:"height"`
	WindowPos         pixel.Vec `yaml:"pos"`
	ShowRemainingTime bool      `yaml:"show_remaining_time"`
}

func (app *App) loadConfig(debug *bool, path *string) func() {
	return func() {
		app.logger = newLogger(*debug)
		if app.initErr != nil {
			return
		}

		app.logger.Debug("Running in debug mode")

		var r io.Reader
		if f, err := os.Open(*path); err == nil {
			app.logger.Info("Loading configuration", zap.String("path", *path))
			r = f
			defer f.Close()
		} else if os.IsNotExist(err) {
			app.logger.Info("Configuration file not found. Creating new file", zap.String("path", *path))
		} else {
			app.initErr = errors.Wrap(err, "failed to open config file")
			return
		}

		app.conf, app.initErr = LoadConfig(r, app.logger, *path, *debug)
		if app.initErr != nil {
			return
		}

		app.initErr = app.conf.Save()
	}
}

func LoadConfig(r io.Reader, logger *zap.Logger, path string, debug bool) (Config, error) {
	conf := Config{path: path}
	if r != nil {
		dec := yaml.NewDecoder(r)
		dec.KnownFields(true)
		err := dec.Decode(&conf)
		if err != nil {
			return conf, errors.Wrap(err, "failed to decode config")
		}
	}

	if conf.UI.WindowWidth == 0 {
		conf.UI.WindowWidth = 512
	}
	if conf.UI.WindowHeight == 0 {
		conf.UI.WindowHeight = 32
	}

	if conf.UI.WindowPos.X == 0 && conf.UI.WindowPos.Y == 0 {
		var displayWidth float64 = 1920 // TODO: make dynamic
		conf.UI.WindowPos = pixel.Vec{
			X: displayWidth/2 - float64(conf.UI.WindowWidth)/2,
			Y: 35,
		}
	}

	if conf.WorkDuration == 0 {
		conf.WorkDuration = 8 * time.Hour
	}
	if conf.LunchDuration == 0 {
		conf.LunchDuration = 1 * time.Hour
	}
	if conf.DayEnd.Hour == 0 {
		conf.DayEnd = ClockTime{Hour: 20, Minute: 00}
	}
	if conf.UI.FPS == 0 {
		conf.UI.FPS = 10
	}
	if conf.CheckIn.IsZero() || isDifferentDay(conf.CheckIn, time.Now()) {
		conf.CheckIn = time.Now()
		logger.Info("Detected start of new day", zap.String("date", conf.CheckIn.Format("2006-01-02")))
	}

	conf.CheckIn = conf.CheckIn.Round(time.Second)
	conf.CheckOut = conf.CheckIn.Add(conf.WorkDuration).Add(conf.LunchDuration)
	conf.EndOfDay = conf.DayEnd.Time(conf.CheckIn)
	conf.Debug = debug

	return conf, nil
}

func isDifferentDay(a, b time.Time) bool {
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()
	return yearA != yearB || monthA != monthB || dayA != dayB
}

func (conf Config) Save() error {
	conf.UI.WindowPos.X = math.Round(conf.UI.WindowPos.X)
	conf.UI.WindowPos.Y = math.Round(conf.UI.WindowPos.Y)

	data, err := yaml.Marshal(conf)
	if err != nil {
		return errors.Wrap(err, "failed to encode config as YAML")
	}

	err = ioutil.WriteFile(conf.path, data, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to save config")
	}

	return nil
}

func (conf Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if conf.CheckIn.IsZero() {
		enc.AddString("check_in", "-")
	} else {
		enc.AddString("check_in", conf.CheckIn.Format("2006-01-02 15:04"))
	}

	enc.AddDuration("work_duration", conf.WorkDuration)
	enc.AddDuration("lunch_duration", conf.LunchDuration)
	enc.AddString("day_end", conf.DayEnd.String())

	return nil
}

type ClockTime struct {
	Hour, Minute int
}

func (t *ClockTime) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), ":")
	if len(parts) != 2 {
		return errors.New(`ClockTime string is not formatted as "hh:mm"`)
	}

	var err error
	t.Hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return errors.Errorf("hour part is not an integer")
	}

	t.Minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return errors.Errorf("minute part is not an integer")
	}

	if t.Hour < 0 || t.Hour > 24 {
		return errors.Errorf("invalid hour")
	}

	if t.Minute < 0 || t.Minute > 59 {
		return errors.Errorf("invalid minute")
	}

	return nil
}

func (t ClockTime) Time(ref time.Time) time.Time {
	year, month, day := ref.Date()
	return time.Date(year, month, day, t.Hour, t.Minute, 0, 0, ref.Location())
}

func (t ClockTime) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

func (t ClockTime) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}
