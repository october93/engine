package worker

import (
	"runtime"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/october93/engine/metrics"
	"github.com/october93/engine/metrics/dogstatsd"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/kit/log"
)

// DatabaseMonitor sends selected metrics about data accessible from the
// database, i.e. number of users in the database.
type DatabaseMonitor struct {
	store       *datastore.Store
	interval    time.Duration
	stopChannel chan struct{}

	users      metrics.Gauge
	goroutines metrics.Gauge
	log        log.Logger
}

// NewDatabaseMonitor returns a new instance of DatabaseMonitor.
func NewDatabaseMonitor(ds *datastore.Store, s *statsd.Client, l log.Logger) *DatabaseMonitor {
	return &DatabaseMonitor{
		store:      ds,
		interval:   5 * time.Second,
		users:      dogstatsd.NewGauge("users", 1, s, l),
		goroutines: dogstatsd.NewGauge("goroutines", 1, s, l),
		log:        l,
	}
}

// Start launches a go routine which will continiously update database
// metrics until Stop is called.
func (dm *DatabaseMonitor) Start() {
	go func() {
		dm.stopChannel = make(chan struct{})
		for {
			select {
			case <-dm.stopChannel:
				return
			case <-time.After(dm.interval):
				dm.updateMetrics()
			}
		}
	}()
}

// Stop shuts down the go routine launched in Start.
func (dm *DatabaseMonitor) Stop() {
	dm.stopChannel <- struct{}{}
}

func (dm *DatabaseMonitor) updateMetrics() {
	count, err := dm.store.GetUserCount()
	if err != nil {
		dm.log.Error(err)
	}
	dm.users.Set(float64(count))
	dm.goroutines.Set(float64(runtime.NumGoroutine()))
}
