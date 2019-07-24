package gql

import (
	"context"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (g *GraphQL) FeatureSwitch_testingUsers(ctx context.Context, obj *model.FeatureSwitch) ([]globalid.ID, error) {
	ids := make([]globalid.ID, 0)
	for i, v := range obj.TestingUsers {
		if v {
			ids = append(ids, i)
		}
	}
	return ids, nil
}

// TODO: this is not neede if we use the enum types in the model definition
func (g *GraphQL) FeatureSwitch_state(ctx context.Context, obj *model.FeatureSwitch) (FeatureState, error) {
	return FeatureState(obj.State), nil
}
