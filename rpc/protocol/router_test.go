package protocol

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kr/pretty"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"

	"testing"
)

func TestRoute(t *testing.T) {
	if testing.Short() {
		t.Skip("short testing detected. skipping test")
		return
	}
	logger, err := log.NewLogger(false, "debug")
	if err != nil {
		t.Fatal(err)
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	//nolint
	defer os.RemoveAll(tmpDir)

	cfg := store.NewConfig()
	cfg.Datastore = datastore.NewTestConfig()
	cfg.Datastore.Database = "engine_rpc_protocol_route_test"

	db := test.DBInit(t, cfg.Datastore)
	// drop the database after the test is finsihed
	defer func() {
		if e := db.Dialect.DropDB(); e != nil {
			t.Fatalf("drop database failed: %s", e)
		}
	}()

	pretty.Println(cfg) // nolint
	str, err := store.NewStore(&cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if e := str.Close(); e != nil {
			t.Fatalf("store.Close() failed.: %s", e)
		}
	}()

	routeTests := []struct {
		context  context.Context
		request  string
		response string
	}{
		{
			context:  context.Background(),
			request:  `{"rpc": "unknown"}`,
			response: `{"error":"RPC unknown not found"}`,
		},
	}

	for _, tt := range routeTests {
		config := NewConfig()
		buf := &bytes.Buffer{}
		r := NewRouter(str, NewConnections(nil, log.NopLogger()), &model.Settings{}, &config, log.NopLogger())
		r.Route(tt.context, model.NewSession(nil), &PushWriter{writer: buf, log: log.NopLogger()}, tt.request)
		if strings.TrimSuffix(buf.String(), "\n") != tt.response {
			t.Errorf("Route(%s): expected error %v, actual %v", tt.request, tt.response, buf.String())
		}
	}
}

func PingEndpoint() MessageEndpoint {
	return func(ctx context.Context, session *model.Session, pw *PushWriter, m Message) error {
		_, err := pw.Write([]byte("Pong!"))
		return err
	}
}

func TestUpgradeConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("short testing detected. skipping test")
		return
	}
	logger, err := log.NewLogger(false, "debug")
	if err != nil {
		t.Fatal(err)
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	//nolint
	defer os.RemoveAll(tmpDir)

	cfg := store.NewConfig()
	cfg.Datastore = datastore.NewTestConfig()
	cfg.Datastore.Database = "engine_rpc_protocol_route_upgrade_connection_test"

	db := test.DBInit(t, cfg.Datastore)
	// drop the database after the test is finsihed
	defer func() {
		if e := db.Dialect.DropDB(); e != nil {
			t.Fatalf("drop database failed: %s", e)
		}
	}()

	str, err := store.NewStore(&cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if e := str.Close(); e != nil {
			t.Fatalf("store.Close() failed.: %s", e)
		}
	}()

	rconfig := NewConfig()
	router := NewRouter(str, NewConnections(nil, log.NopLogger()), &model.Settings{}, &rconfig, log.NopLogger())
	router.RegisterRPC("ping?", PingEndpoint())

	ts := httptest.NewServer(router.UpgradeConnection())
	defer ts.Close()

	client1 := webSocketClient(t, extractURL(t, ts))
	client2 := webSocketClient(t, extractURL(t, ts))

	pingRPC := `{"rpc":"ping?"}`
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := client1.WriteMessage(websocket.TextMessage, []byte(pingRPC))
		if err != nil {
			t.Errorf("WriteMessage(): unexpected error %v", err)
			return
		}
		messageType, p, err := client1.ReadMessage()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if messageType != websocket.TextMessage {
			t.Errorf("expected message type %d, actual %d", websocket.TextMessage, messageType)
		}
		if string(p) != "Pong!" {
			t.Errorf("expected response %s, actual %s", "Pong!", string(p))
		}
		wg.Done()
	}()

	go func() {
		err := client2.WriteMessage(websocket.TextMessage, []byte(pingRPC))
		if err != nil {
			t.Errorf("WriteMessage(): unexpected error %v", err)
			return
		}
		messageType, p, err := client2.ReadMessage()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if messageType != websocket.TextMessage {
			t.Errorf("expected message type %d, actual %d", websocket.TextMessage, messageType)
		}
		if string(p) != "Pong!" {
			t.Errorf("expected response %s, actual %s", "Pong!", string(p))
		}
		wg.Done()
	}()
	wg.Wait()
}

func TestHandleConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("short testing detected. skipping test")
		return
	}
	logger, err := log.NewLogger(false, "debug")
	if err != nil {
		t.Fatal(err)
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	//nolint
	defer os.RemoveAll(tmpDir)

	cfg := store.NewConfig()
	cfg.Datastore = datastore.NewTestConfig()
	cfg.Datastore.Database = "engine_rpc_protocol_route_handle_connection_test"

	db := test.DBInit(t, cfg.Datastore)
	// drop the database after the test is finsihed
	defer test.DBCleanup(t, db)

	str, err := store.NewStore(&cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if e := str.Close(); e != nil {
			t.Fatalf("store.Close() failed.: %s", e)
		}
	}()

	rconfig := NewConfig()
	router := NewRouter(str, NewConnections(nil, log.NopLogger()), &model.Settings{}, &rconfig, log.NopLogger())
	router.RegisterRPC("ping?", PingEndpoint())

	ts := httptest.NewServer(router.UpgradeConnection())
	defer ts.Close()
	client := webSocketClient(t, extractURL(t, ts))

	go router.HandleConnection(context.Background(), client, httptest.NewRequest("", "/deck_endpoint/", nil))

	err = client.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"ping?"}`))
	if err != nil {
		t.Fatalf("WriteMessage(): unexpected error %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	err = client.Close()
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		t.Errorf("Close(): unexpected error: %v", err)
	}
	if router.conns.Count() != 0 {
		t.Errorf("HandleConnection(): expected writeres to be cleaned up, actual: %d writers(s)", router.conns.Count())
	}
}

func extractURL(t *testing.T, ts *httptest.Server) string {
	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	url.Scheme = "ws"
	return url.String()
}

func webSocketClient(t *testing.T, url string) *websocket.Conn {
	dialer := websocket.DefaultDialer
	conn, resp, err := dialer.Dial(url, nil)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("Dial(): expected response status code %d, actual %d", 101, resp.StatusCode)
	}
	return conn
}
