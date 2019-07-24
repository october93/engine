package protocol

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/october93/engine/kit/log"
)

func TestWrite(t *testing.T) {
	var writeTests = []struct {
		data []string
	}{
		{[]string{"October", "is", "coming"}},
	}

	for _, tt := range writeTests {
		var buf bytes.Buffer
		writer := &PushWriter{writer: &buf, log: log.NopLogger()}

		var wg sync.WaitGroup
		wg.Add(len(tt.data))
		for _, input := range tt.data {
			go func(input string) {
				defer wg.Done()
				_, err := writer.Write([]byte(input))
				if err != nil {
					t.Logf("Write(): unexpected error: %v\ndata: %s\n", err, input)
				}
			}(input)
		}
		wg.Wait()

		for i := range tt.data {
			if !strings.Contains(buf.String(), tt.data[i]) {
				t.Errorf("Write(%s): expected buffer to contain %s, actual %s", tt.data[i], tt.data[i], buf.String())
			}
		}
	}
}
