package store

import (
	"errors"
	"time"
)

// Params is a store parameters
type Params struct {
	// CleaningPeriod determines period of expired keys removing from store
	CleaningPeriod time.Duration

	// DumpingPeriod determines how often store will be dumped to the disk
	DumpingPeriod time.Duration
}

// Validate validates store parameters
func (p Params) Validate() error {
	if p.CleaningPeriod < 100*time.Millisecond {
		return errors.New("too small cleaning period, must be >= 100ms")
	}
	if p.DumpingPeriod < 60*time.Second {
		return errors.New("too small dumping period, must be >= 60s")
	}
	return nil
}
