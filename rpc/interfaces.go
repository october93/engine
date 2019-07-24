package rpc

import (
	"context"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

type notifications interface {
	ExportNotification(n *model.Notification) (*model.ExportedNotification, error)
}

type pusher interface {
	NewCard(ctx context.Context, session *model.Session, card *model.CardResponse) error
	DeleteCard(ctx context.Context, session *model.Session, id globalid.ID) error
	UpdateCard(ctx context.Context, session *model.Session, card *model.CardResponse) error
	UpdateUser(ctx context.Context, session *model.Session, user *model.ExportedUser) error
	UpdateCoinBalance(ctx context.Context, userID globalid.ID, newBalances *model.CoinBalances) error
	NewNotification(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error
	UpdateNotification(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error
	UpdateEngagement(ctx context.Context, session *model.Session, cardID globalid.ID) error
}

type indexer interface {
	IndexUser(m *model.User) error
	RemoveIndexForUser(m *model.User) error
	IndexChannel(m *model.Channel) error
}

type responses interface {
	FeedCardResponses(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error)
}
