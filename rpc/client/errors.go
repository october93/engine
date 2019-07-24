package client

import "errors"

// ErrTimeout happens when there is no response for a given request in the
// specified time frame.
var ErrTimeout = errors.New("request timed out")
