package gql

import (
	"context"

	"github.com/october93/engine/model"
)

func (g *GraphQL) Announcement_user(ctx context.Context, obj *model.Announcement) (*model.User, error) {
	return ctxLoaders(ctx).UserByID.Load(obj.UserID)
}
func (g *GraphQL) Announcement_card(ctx context.Context, obj *model.Announcement) (*model.Card, error) {
	return ctxLoaders(ctx).CardByID.Load(obj.CardID)
}
