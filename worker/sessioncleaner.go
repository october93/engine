package worker

import (
	"time"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/store"
)

// SessionCleaner is responsible for invalidating sessions which have been
// expired.
type SessionCleaner struct {
	store       *store.Store
	interval    time.Duration
	stopChannel chan struct{}
	log         log.Logger
}

// NewSessionCleaner returns a new instance of SessionCleaner.
func NewSessionCleaner(s *store.Store, l log.Logger) *SessionCleaner {
	return &SessionCleaner{store: s, interval: 10 * time.Minute, log: l}
}

// Start starts the worker. This method should be called as a
// goroutine, otherwise it will block any further execution.
func (sc *SessionCleaner) Start() {
	sc.stopChannel = make(chan struct{})
	go func() {
		for {
			select {
			case <-sc.stopChannel:
				return
			case <-time.After(sc.interval):
				sc.cleanup()
			}
		}
	}()
}

func (sc *SessionCleaner) cleanup() {
	deleted, err := sc.store.DeleteExpiredSessions()
	if err != nil {
		sc.log.Bug(err)
	}
	if deleted > 0 {
		sc.log.Info("sessions expired", "deleted", deleted)
	}
}

// Stop shuts down the go routine launched in Start.
func (sc *SessionCleaner) Stop() {
	close(sc.stopChannel)
}
