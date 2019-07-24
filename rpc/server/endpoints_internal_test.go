package server

import (
	"bytes"
	"context"
	"testing"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
)

type mockRegisterDevice struct {
	rpc.RPC
	token    string
	platform string
}

func (c *mockRegisterDevice) RegisterDevice(ctx context.Context, req rpc.RegisterDeviceRequest) (*rpc.RegisterDeviceResponse, error) {
	c.token = req.Params.Token
	c.platform = req.Params.Platform
	return nil, nil
}

func (c *mockRegisterDevice) UnregisterDevice(ctx context.Context, req rpc.UnregisterDeviceRequest) (*rpc.UnregisterDeviceResponse, error) {
	c.token = req.Params.Token
	return nil, nil
}

func TestRegisterDeviceEndpoint(t *testing.T) {
	var registerDeviceEndpointTests = []struct {
		request  string
		token    string
		platform string
	}{
		{
			request:  `{"token": "123", "platform": "ios"}`,
			token:    "123",
			platform: "ios",
		},
		{
			request:  `{}`,
			token:    "",
			platform: "",
		},
	}

	for _, tt := range registerDeviceEndpointTests {
		rpc := &mockRegisterDevice{}
		endpoint := RegisterDeviceEndpoint(rpc)
		ctx := context.Background()
		ctx = context.WithValue(ctx, protocol.RequestID, globalid.Next())
		ctx = context.WithValue(ctx, protocol.Callback, "")

		user := &model.User{ID: globalid.NewFixture().Next()}
		session := model.NewSession(user)
		writer := protocol.NewPushWriter(&protocol.Connection{}, &bytes.Buffer{}, log.NopLogger())
		writer.SetSession(session)

		err := endpoint(ctx, session, writer, protocol.Message{Data: []byte(tt.request)})
		if err != nil {
			t.Errorf("endpoint(%v): unexpected error: %v", tt.request, err)
		}
		if rpc.token != tt.token {
			t.Errorf("RegisterDeviceEndpoint(%s): expected token to be: %s, actual: %s", tt.request, rpc.token, tt.token)
		}
		if rpc.platform != tt.platform {
			t.Errorf("RegisterDeviceEndpoint(%s): expected platform to be: %s, actual: %s", tt.request, rpc.platform, tt.platform)
		}
	}
}

func TestUnregisterDeviceEndpoint(t *testing.T) {
	var unregisterDeviceEndpointTests = []struct {
		request string
		token   string
	}{
		{
			request: `{"token": "123", "platform": "ios"}`,
			token:   "123",
		},
		{
			request: `{}`,
			token:   "",
		},
	}

	for _, tt := range unregisterDeviceEndpointTests {
		rpc := &mockRegisterDevice{}
		endpoint := UnregisterDeviceEndpoint(rpc)
		ctx := context.Background()
		ctx = context.WithValue(ctx, protocol.RequestID, globalid.Next())
		ctx = context.WithValue(ctx, protocol.Callback, "")

		user := &model.User{ID: globalid.NewFixture().Next()}
		session := model.NewSession(user)
		writer := protocol.NewPushWriter(&protocol.Connection{}, &bytes.Buffer{}, log.NopLogger())
		writer.SetSession(session)

		err := endpoint(ctx, session, writer, protocol.Message{Data: []byte(tt.request)})
		if err != nil {
			t.Errorf("endpoint(%v): unexpected error: %v", tt.request, err)
		}
		if rpc.token != tt.token {
			t.Errorf("UnregisterDeviceEndpoint(%s): expected token to be: %s, actual: %s", tt.request, rpc.token, tt.token)
		}
	}
}
