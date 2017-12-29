package store

import "time"

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (_ SystemClock) Now() time.Time {
	return time.Now()
}