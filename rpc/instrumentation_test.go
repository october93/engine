package rpc

import (
	"context"
	"testing"

	"github.com/october93/engine/kit/log"
)

func TestExtractUsername(t *testing.T) {
	extractUsernameTests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "valid",
			ctx:      context.WithValue(context.Background(), "username", "chad"),
			expected: "chad",
		},
		{
			name:     "no username",
			ctx:      context.Background(),
			expected: "No User",
		},
		{
			name:     "invalid type",
			ctx:      context.WithValue(context.Background(), "username", 93),
			expected: "",
		},
	}
	for _, tt := range extractUsernameTests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			im := instrumentingMiddleware{l: log.NopLogger()}
			username := im.extractUsername(tt.ctx)
			if username != tt.expected {
				t.Errorf("extractUsername(): expected %s, actual: %s", tt.expected, username)
			}
		})
	}
}
