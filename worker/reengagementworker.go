package worker

import (
	"time"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/store"
)

// SessionCleaner is responsible for invalidating sessions which have been
// expired.
type ReengagementWorker struct {
	store       *store.Store
	notifier    *Notifier
	interval    time.Duration
	stopChannel chan struct{}
	log         log.Logger
}

// NewSessionCleaner returns a new instance of SessionCleaner.
func NewReengagementWorker(s *store.Store, n *Notifier, l log.Logger) *ReengagementWorker {
	return &ReengagementWorker{store: s, interval: 4 * time.Hour, log: l}
}

// Start starts the worker. This method should be called as a
// goroutine, otherwise it will block any further execution.
func (rw *ReengagementWorker) Start() {
	rw.stopChannel = make(chan struct{})
	go func() {
		for {
			select {
			case <-rw.stopChannel:
				return
			case <-time.After(rw.interval):
				rw.cleanup()
			}
		}
	}()
}

func (rw *ReengagementWorker) cleanup() {
	users, err := rw.store.GetUsers()
	if err != nil {
		rw.log.Error(err)
	}

	for _, user := range users {
		lastActive, err := rw.store.GetLastActiveAt(user.ID)
		if err != nil {
			rw.log.Error(err)
		}

		if lastActive.Before(time.Now().Add(-96 * time.Hour)) {
			nerr := rw.notifier.NotifySlack("engagement", "%v (%v) hasn't logged in in more than 4 days, re-engage them!")
			if nerr != nil {
				rw.log.Error(nerr)
			}
		}
	}
}

// Stop shuts down the go routine launched in Start.
func (rw *ReengagementWorker) Stop() {
	close(rw.stopChannel)
}
