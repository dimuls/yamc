package store

import (
	"errors"
	"time"
)

// ticker is timer ticker, runs function f every period time
type ticker struct {
	period  time.Duration
	f       func()
	stopper chan struct{}
}

// newTicker constructs new ticker with given period and working function f. Returns error if f is nil
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

// isRunning determines if ticker is running
func (t *ticker) isRunning() bool {
	return t.stopper != nil
}

// start starts ticker. Can be called multiple times. Returns error if ticker is already started
func (t *ticker) start() error {
	if t.isRunning() {
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

// stop stops ticker. Can be called multiple times. Returns errors if ticker already stopped
func (t *ticker) stop() error {
	if !t.isRunning() {
		return errors.New("already stopped")
	}
	close(t.stopper)
	t.stopper = nil
	return nil
}
