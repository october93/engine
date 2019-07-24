package test

import (
	"net"
	"testing"
)

func FreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Listen(): %v", err)
	}
	defer func() {
		err = l.Close()
		if err != nil {
			t.Fatalf("Close(): %v", err)
		}
	}()
	return l.Addr().(*net.TCPAddr).Port

}
