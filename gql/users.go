package gql

import (
	"context"
	"math"
	"time"

	"github.com/october93/engine/dataloader"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (g *GraphQL) Query_users(ctx context.Context, usernames []string) ([]model.User, error) {
	var users []*model.User
	var err error
	if len(usernames) > 0 {
		users, err = g.Store.GetUsersByUsernames(usernames)
		if err != nil {
			return nil, err
		}
	} else {
		users, err = g.Store.GetUsers()
		if err != nil {
			return nil, err
		}
	}
	ret := make([]model.User, len(users))
	for i := 0; i < len(users); i++ {
		ret[i] = *users[i]
	}
	return ret, nil
}

func (g *GraphQL) User_lastActiveAt(ctx context.Context, obj *model.User) (*model.DBTime, error) {
	return ctxLoaders(ctx).LastActiveAtByID.Load(obj.ID)
}

//hack to account for the fact that it's hard to make a dataloader for time.Time
func (g *GraphQL) ModelTime_time(ctx context.Context, obj *model.DBTime) (*time.Time, error) {
	if obj == nil || !obj.Valid {
		return nil, nil
	}
	return &obj.WEIRDNAME, nil
}

func (g *GraphQL) User_blocked(ctx context.Context, obj *model.User) (bool, error) {
	return obj.BlockedAt.Valid, nil
}

func (g *GraphQL) User_shadowbanned(ctx context.Context, obj *model.User) (bool, error) {
	return obj.ShadowbannedAt.Valid, nil
}

func (g *GraphQL) User_coinBalance(ctx context.Context, obj *model.User) (int, error) {
	return int(math.Floor(float64(obj.CoinBalance) / 10000.0)), nil
}
func (g *GraphQL) User_temporaryCoinBalance(ctx context.Context, obj *model.User) (int, error) {
	return int(math.Floor(float64(obj.TemporaryCoinBalance) / 10000.0)), nil
}

func (g *GraphQL) User_postCount(ctx context.Context, obj *model.User, from *time.Time, to *time.Time) (int, error) {
	start, end := containingWeek(time.Now(), START_OF_WEEK_DAY)

	if from != nil {
		start = *from
	}

	if to != nil {
		end = *to
	}

	posts, err := g.Store.GetCardsByNodeInRange(
		obj.ID,
		start,
		end,
	)

	if err != nil {
		return 0, err
	}
	return len(posts), nil
}

func (g *GraphQL) User_joinedFromInvite(ctx context.Context, obj *model.User) (*model.Invite, error) {
	if obj.JoinedFromInvite != globalid.Nil {
		return ctxLoaders(ctx).InvitesByID.Load(obj.JoinedFromInvite)
	}
	return nil, nil
}

func (g *GraphQL) User_engagement(ctx context.Context, user *model.User, from *string, to *string) (*model.UserEngagement, error) {
	ldrs := ctxLoaders(ctx)

	layout := "2006-01-02T15:04:05-07:00"
	startDate, err := time.Parse(layout, *from)
	if err != nil {
		return nil, err
	}
	startDate = startDate.UTC()
	endDate, err := time.Parse(layout, *to)
	if err != nil {
		return nil, err
	}
	endDate = endDate.UTC()

	engagement, err := ldrs.UserEngagementLoader.Load(dataloader.NewIDTimeRangeKey(user.ID, startDate, endDate))
	if err != nil {
		return nil, err
	}
	return engagement, nil
}
