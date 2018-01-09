package store

import "time"

// Clock is a now time provider
type Clock interface {
	now() time.Time
}

// SystemClock is clock providing system time
type SystemClock struct{}

// now returns current time
func (_ SystemClock) now() time.Time {
	return time.Now()
}
