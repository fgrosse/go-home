package main

import "time"

type Config struct {
	WorkDuration  time.Duration
	LunchDuration time.Duration
	DayEnd        ClockTime

	WindowWidth  int
	WindowHeight int
	VSync        bool
}

type ClockTime struct {
	Hour, Minute int
}

func (t ClockTime) Time(ref time.Time) time.Time {
	year, month, day := ref.Date()
	return time.Date(year, month, day, t.Hour, t.Minute, 0, 0, ref.Location())
}
