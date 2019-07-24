package debug

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/october93/engine/kit/log"
)

// Profiler records runtime profiling data using pprof and makes it available.
type Profiler struct {
	config Config
	log    log.Logger
}

// NewProfiler returns a new instance of profiler.
func NewProfiler(c Config, l log.Logger) *Profiler {
	return &Profiler{config: c, log: l}
}

// Start listens on the configured port and serves the runtime profiling data
// over HTTP. This is a blocking method and should be run concurrently.
func (p *Profiler) Start() error {
	address := fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)
	return http.ListenAndServe(address, nil)
}
