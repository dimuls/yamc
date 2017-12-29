package store

import (
	"errors"
	"time"
)

type Params struct {
	CleaningPeriod time.Duration
}

func (p Params) Validate() error {
	if p.CleaningPeriod < 100*time.Millisecond {
		return errors.New("too small cleaning period")
	}
	return nil
}
