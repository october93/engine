package push

import (
	"context"

	nats "github.com/nats-io/go-nats"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/rpc/protocol"
)

const (
	pushToAllClients          = "push:all"
	pushToSpecificUserClients = "push:user"
)

type FanOut struct {
	bus   *nats.EncodedConn
	conns connections
	log   log.Logger
}

type NatsMessage struct {
	Message *protocol.Message `json:"message"`
	UserID  globalid.ID       `json:"userID"`
}

func NewFanOut(bus *nats.EncodedConn, conns connections, log log.Logger) (*FanOut, error) {
	fa := &FanOut{
		bus:   bus,
		conns: conns,
		log:   log,
	}
	_, err := bus.Subscribe(pushToAllClients, fa.ListenToAll)
	if err != nil {
		return nil, err
	}
	_, err = bus.Subscribe(pushToSpecificUserClients, fa.ListenToUserMessage)
	if err != nil {
		return nil, err
	}
	return fa, nil
}

func (fa *FanOut) PublishToAll(m *protocol.Message) error {
	return fa.bus.Publish(pushToAllClients, m)
}

func (fa *FanOut) PublishToUser(userID globalid.ID, m *protocol.Message) error {
	return fa.bus.Publish(pushToSpecificUserClients, NatsMessage{Message: m, UserID: userID})
}

func (fa *FanOut) ListenToAll(m *protocol.Message) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, protocol.RequestID, m.RequestID)
	for _, writer := range fa.conns.Writers() {
		err := m.Encode(ctx, writer)
		if err != nil {
			fa.log.Error(err)
		}
	}
}

func (fa *FanOut) ListenToUserMessage(m *NatsMessage) {
	ctx := context.WithValue(context.Background(), protocol.RequestID, m.Message.RequestID)
	for _, writer := range fa.conns.WritersByUser(m.UserID) {
		err := m.Message.Encode(ctx, writer)
		if err != nil {
			fa.log.Error(err)
		}
	}
}
