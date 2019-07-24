package push

import (
	"context"
	"encoding/json"
	"errors"

	nats "github.com/nats-io/go-nats"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	rpccontext "github.com/october93/engine/rpc/context"
	"github.com/october93/engine/rpc/protocol"
)

const (
	newCard            = "newCard"
	deleteCard         = "deleteCard"
	updateCard         = "updateCard"
	updateUser         = "updateUser"
	newNotification    = "newNotification"
	updateNotification = "updateNotification"
	updateEngagement   = "updateEngagement"
	updateCoinBalance  = "updateCoinBalance"
)

type Pusher struct {
	fanOut *FanOut
	conns  connections
	store  store
	config *Config
	log    log.Logger
}

func NewPusher(conns connections, s store, c *Config, l log.Logger) (*Pusher, error) {
	p := &Pusher{
		conns:  conns,
		store:  s,
		config: c,
		log:    l,
	}
	if c.FanOut {
		nc, err := nats.Connect(c.NatsEndpoint)
		if err != nil {
			return nil, err
		}
		encodedConn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		if err != nil {
			return nil, err
		}
		fanOut, err := NewFanOut(encodedConn, conns, l)
		if err != nil {
			return nil, err
		}
		p.fanOut = fanOut
	}
	return p, nil
}

func NewMessage(ctx context.Context, rpc string, data interface{}) (*protocol.Message, error) {
	m := protocol.NewMessage(rpc)
	err := m.EncodePayload(data)
	if err != nil {
		return nil, err
	}
	requestID, ok := ctx.Value(rpccontext.RequestID).(globalid.ID)
	if !ok {
		return nil, errors.New("invalid RequestID")
	}
	m.RequestID = requestID
	return m, nil
}

func (p *Pusher) NewCard(ctx context.Context, session *model.Session, card *model.CardResponse) error {
	m, err := NewMessage(ctx, newCard, card)
	if err != nil {
		return err
	}
	for _, writer := range p.conns.Writers() {
		err = m.Encode(ctx, writer)
		if err != nil {
			p.log.Error(err)
		}
	}
	if p.config.FanOut {
		return p.fanOut.PublishToAll(m)
	}
	return nil
}

func (p *Pusher) NewNotification(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
	m, err := NewMessage(ctx, newNotification, notif)
	if err != nil {
		return err
	}
	for _, writer := range p.conns.WritersByUser(notif.UserID) {
		err = m.Encode(ctx, writer)
		if err != nil {
			p.log.Error(err)
		}
	}
	if p.config.FanOut {
		return p.fanOut.PublishToUser(notif.UserID, m)
	}
	return nil
}

func (p *Pusher) UpdateNotification(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
	m, err := NewMessage(ctx, updateNotification, notif)
	if err != nil {
		return err
	}
	for _, writer := range p.conns.WritersByUser(notif.UserID) {
		err = m.Encode(ctx, writer)
		if err != nil {
			p.log.Error(err)
		}
	}
	if p.config.FanOut {
		return p.fanOut.PublishToUser(notif.UserID, m)
	}
	return nil
}

func (p *Pusher) DeleteCard(ctx context.Context, session *model.Session, id globalid.ID) error {
	m, err := NewMessage(ctx, deleteCard, id)
	if err != nil {
		return err
	}
	for _, writer := range p.conns.Writers() {
		err = m.Encode(ctx, writer)
		if err != nil {
			p.log.Error(err)
		}
	}
	if p.config.FanOut {
		return p.fanOut.PublishToAll(m)
	}
	return nil
}

func (p *Pusher) UpdateCard(ctx context.Context, session *model.Session, card *model.CardResponse) error {
	m, err := NewMessage(ctx, updateCard, card)
	if err != nil {
		return err
	}
	for _, writer := range p.conns.Writers() {
		err = m.Encode(ctx, writer)
		if err != nil {
			p.log.Error(err)
		}
	}
	if p.config.FanOut {
		return p.fanOut.PublishToAll(m)
	}
	return nil
}

func (p *Pusher) UpdateCoinBalance(ctx context.Context, userID globalid.ID, newBalances *model.CoinBalances) error {
	m, err := NewMessage(ctx, updateCoinBalance, newBalances)
	if err != nil {
		return err
	}
	for _, writer := range p.conns.WritersByUser(userID) {
		err = m.Encode(ctx, writer)
		if err != nil {
			p.log.Error(err)
		}
	}
	if p.config.FanOut {
		return p.fanOut.PublishToUser(userID, m)
	}
	return nil
}

func (p *Pusher) UpdateUser(ctx context.Context, session *model.Session, user *model.ExportedUser) error {
	m, err := NewMessage(ctx, updateUser, user)
	if err != nil {
		return err
	}
	for _, writer := range p.conns.Writers() {
		err = m.Encode(ctx, writer)
		if err != nil {
			p.log.Error(err)
		}
	}
	if p.config.FanOut {
		return p.fanOut.PublishToAll(m)
	}
	return nil
}

type EngagementUpdate struct {
	CardID     globalid.ID       `json:"cardID"`
	Engagement *model.Engagement `json:"engagement"`
}

func (p *Pusher) UpdateEngagement(ctx context.Context, session *model.Session, cardID globalid.ID) error {
	engagement, err := p.store.GetEngagement(cardID)
	if err != nil {
		return err
	}
	update := EngagementUpdate{
		CardID:     cardID,
		Engagement: engagement,
	}
	for _, writer := range p.conns.Writers() {
		if writer.Session() != nil {
			userID := writer.Session().UserID
			m, err := NewMessage(ctx, updateEngagement, update)
			if err != nil {
				return err
			}
			b, err := json.Marshal(m)
			if err != nil {
				return err
			}
			_, err = writer.Write(b)
			if err != nil {
				p.log.Error(err)
			}
			if p.config.FanOut {
				err = p.fanOut.PublishToUser(userID, m)
				if err != nil {
					p.log.Error(err)
				}
			}
		}
	}
	return nil
}
