package helper

import "time"

type logTime struct {
	TimeStart time.Time
}

func NewLogTime(time time.Time) *logTime {
	return &logTime{TimeStart: time}
}

func (t logTime) GetTimeSince() float64 {
	return time.Since(t.TimeStart).Seconds()
}
