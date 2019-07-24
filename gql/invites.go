package gql

import (
	"context"

	"github.com/october93/engine/model"
)

func (g *GraphQL) Invite_issuer(ctx context.Context, obj *model.Invite) (model.User, error) {
	user, err := ctxLoaders(ctx).UserByID.Load(obj.NodeID)
	return *user, err
}
