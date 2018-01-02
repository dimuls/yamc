package store

import (
	"errors"
	"time"
)

// Params is a store parameters
type Params struct {
	// CleaningPeriod determines period of expired keys removing from store
	CleaningPeriod time.Duration
}

// Validate validates store parameters
func (p Params) Validate() error {
	if p.CleaningPeriod < 100*time.Millisecond {
		return errors.New("too small cleaning period")
	}
	return nil
}
