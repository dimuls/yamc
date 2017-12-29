package store

import (
	"errors"
	"time"
)

type ticker struct {
	period  time.Duration
	f       func()
	stopper chan struct{}
}

func newTicker(period time.Duration, f func()) (*ticker, error) {
	if f == nil {
		return nil, errors.New("nil ticker function")
	}
	return &ticker{
		period:  period,
		f:       f,
		stopper: nil,
	}, nil
}

func (t *ticker) running() bool {
	return t.stopper != nil
}

func (t *ticker) start() error {
	if t.running() {
		return errors.New("already started")
	}
	t.stopper = make(chan struct{})
	ticker := time.NewTicker(t.period)
	go func() {
		for {
			select {
			case <-ticker.C:
				t.f()
			case <-t.stopper:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

func (t *ticker) stop() error {
	if !t.running() {
		return errors.New("already stopped")
	}
	close(t.stopper)
	t.stopper = nil
	return nil
}
