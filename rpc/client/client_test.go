package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/test"
)

func TestTimeout(t *testing.T) {
	server, address := setupTestServer(t)
	defer func() {
		err := server.Close()
		if err != nil {
			t.Errorf("Close(): %v", err)
		}
	}()
	logger, err := log.NewLogger(false, "debug")
	if err != nil {
		t.Fatalf("NewLogger(): %v", err)
	}
	config := NewConfig()
	config.Address = fmt.Sprintf("ws://%s/", address)
	config.Timeout = 0
	client, err := NewClient(config, logger)
	if err != nil {
		t.Fatalf("NewClient(): %v", err)
	}
	done := make(chan struct{})
	go func() {
		req := rpc.AuthRequest{Params: rpc.AuthParams{Username: "chad", Password: "secret"}}
		_, err := client.Auth(context.Background(), req)
		if err != ErrTimeout {
			t.Errorf("Auth(): unexpected error %v", err)
		}
		done <- struct{}{}
	}()
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("client call did not return in time")
	}
}

func TestConcurrency(t *testing.T) {
	server, address := setupTestServer(t)
	defer func() {
		err := server.Close()
		if err != nil {
			t.Errorf("Close(): %v", err)
		}
	}()
	logger, err := log.NewLogger(false, "debug")
	if err != nil {
		t.Fatalf("NewLogger(): %v", err)
	}
	config := NewConfig()
	config.Address = fmt.Sprintf("ws://%s/", address)
	client, err := NewClient(config, logger)
	if err != nil {
		t.Fatalf("NewClient(): %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		req := rpc.AuthRequest{Params: rpc.AuthParams{Username: "chad", Password: "secret"}}
		_, err := client.Auth(context.Background(), req)
		if err != nil && err.Error() != "mock websocket server" {
			t.Errorf("Auth(): %v", err)
		}
		wg.Done()
	}()
	go func() {
		req := rpc.AuthRequest{Params: rpc.AuthParams{Username: "chad", Password: "secret"}}
		_, err := client.Auth(context.Background(), req)
		if err != nil && err.Error() != "mock websocket server" {
			t.Errorf("Auth(): %v", err)
		}
		wg.Done()
	}()
	wg.Wait()
}

func setupTestServer(t *testing.T) (*http.Server, string) {
	t.Helper()
	address := fmt.Sprintf("localhost:%d", test.FreePort(t))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := &websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Upgrade(): %v", err)
		}
		for {
			var m protocol.Message
			_, b, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("ReadMessage(): %v", err)
			}
			err = json.Unmarshal(b, &m)
			if err != nil {
				t.Fatalf("Unmarshal(): %v", err)
			}
			var resp protocol.Message
			resp.Ack = m.RequestID
			resp.Err = "mock websocket server"
			resp.Data = []byte(`{}`)
			time.Sleep(5 * time.Millisecond)
			err = c.WriteJSON(resp)
			if err != nil {
				t.Fatalf("WriteJSON(): %v", err)
			}
		}
	})
	server := &http.Server{Addr: address, Handler: handler}
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		t.Errorf("Listen(): %v", err)
	}
	go func() {
		t.Log(server.Serve(ln))
	}()
	return server, address
}
