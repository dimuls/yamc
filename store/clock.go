package store

import "time"

// Clock is a now time provider
type Clock interface {
	Now() time.Time
}

// SystemClock is clock providing system time
type SystemClock struct{}

// Now returns current time
func (_ SystemClock) Now() time.Time {
	return time.Now()
}
