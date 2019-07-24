//go:generate gorunpkg github.com/vektah/dataloaden -keys "globalid.ID" github.com/october93/engine/model.User
//go:generate gorunpkg github.com/vektah/dataloaden -keys "globalid.ID" github.com/october93/engine/model.AnonymousAlias
//go:generate gorunpkg github.com/vektah/dataloaden -keys "globalid.ID" github.com/october93/engine/model.Card
//go:generate gorunpkg github.com/vektah/dataloaden -keys "globalid.ID" github.com/october93/engine/model.Invite
//go:generate gorunpkg github.com/vektah/dataloaden -keys "IDTimeRangeKey" github.com/october93/engine/model.UserEngagement
//go:generate gorunpkg github.com/vektah/dataloaden -keys "globalid.ID" github.com/october93/engine/model.FeedEntriesByID
//go:generate gorunpkg github.com/vektah/dataloaden -keys "globalid.ID" github.com/october93/engine/model.DBTime

package dataloader

import (
	"fmt"
	"reflect"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store"
)

var _ = fmt.Sprint() //because fmt

type Loaders struct {
	UserByID             *UserLoader
	AliasByID            *AnonymousAliasLoader
	CardByID             *CardLoader
	InvitesByID          *InviteLoader
	UserEngagementLoader *UserEngagementLoader
	LastActiveAtByID     *DBTimeLoader
}

func NewLoaders(store *store.Store, wait time.Duration) Loaders {
	ldrs := Loaders{}

	ldrs.UserByID = &UserLoader{
		wait:     wait,
		maxBatch: 100,
		fetch: func(keys []globalid.ID) ([]*model.User, []error) {
			users, err := store.GetUsersByID(keys)
			return users, []error{err}
		},
	}
	ldrs.InvitesByID = &InviteLoader{
		wait:     wait,
		maxBatch: 100,
		fetch: func(keys []globalid.ID) ([]*model.Invite, []error) {
			invites, err := store.GetInvitesByID(keys)
			return invites, []error{err}
		},
	}
	ldrs.AliasByID = &AnonymousAliasLoader{
		wait:     wait,
		maxBatch: 100,
		fetch: func(keys []globalid.ID) ([]*model.AnonymousAlias, []error) {
			aliases, err := store.GetAnonymousAliasesByID(keys)
			return aliases, []error{err}
		},
	}
	ldrs.CardByID = &CardLoader{
		wait:     wait,
		maxBatch: 100,
		fetch: func(keys []globalid.ID) ([]*model.Card, []error) {
			stats, err := store.GetCardsByIDs(keys)
			return stats, []error{err}
		},
	}
	ldrs.LastActiveAtByID = &DBTimeLoader{
		wait:     wait,
		maxBatch: 100,
		fetch: func(keys []globalid.ID) ([]*model.DBTime, []error) {
			times, err := store.BatchGetLastActiveAt(keys)
			return times, []error{err}
		},
	}
	ldrs.UserEngagementLoader = &UserEngagementLoader{
		wait:     wait,
		maxBatch: 100,
		fetch: func(keys []IDTimeRangeKey) ([]*model.UserEngagement, []error) {
			result := make([]*model.UserEngagement, len(keys))
			engagements := make(map[globalid.ID]*model.UserEngagement, len(keys))

			from := keys[0].from
			to := keys[0].to

			for _, key := range keys {
				engagements[key.id] = &model.UserEngagement{UserID: key.id}
			}

			daysActive, err := store.GetDaysActive(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range daysActive {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].DaysActive = c.Count
				if c.Count >= 2 {
					engagements[c.UserID].Score += 1
				}
			}
			postCount, err := store.GetPostCount(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range postCount {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].PostCount = c.Count
				if c.Count >= 1 {
					engagements[c.UserID].Score += 1
				}
			}
			commentCount, err := store.GetCommentCount(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range commentCount {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].CommentCount = c.Count
				if c.Count >= 1 {
					engagements[c.UserID].Score += 1
				}
			}
			reactedCount, err := store.GetReactedCount(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range reactedCount {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].ReactedCount = c.Count
				if c.Count >= 1 {
					engagements[c.UserID].Score += 1
				}
			}
			receivedReactionsCount, err := store.GetReceivedReactionsCount(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range receivedReactionsCount {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].ReceivedReactionsCount = c.Count
				if c.Count >= 1 {
					engagements[c.UserID].Score += 1
				}
			}
			followedUsersCount, err := store.GetFollowedUsersCount(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range followedUsersCount {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].FollowedUsersCount = c.Count
				if c.Count >= 1 {
					engagements[c.UserID].Score += 1
				}
			}
			followedCount, err := store.GetFollowedCount(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range followedCount {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].FollowedCount = c.Count
				if c.Count >= 1 {
					engagements[c.UserID].Score += 1
				}
			}
			invitedCount, err := store.GetInvitedCount(from, to)
			if err != nil {
				return nil, []error{err}
			}
			for _, c := range invitedCount {
				if _, ok := engagements[c.UserID]; !ok {
					continue
				}
				engagements[c.UserID].InvitedCount = c.Count
				if c.Count >= 1 {
					engagements[c.UserID].Score += 1
				}
			}

			t := reflect.TypeOf(model.UserEngagement{})
			numFields := float64(t.NumField() - 2)

			for _, e := range engagements {
				e.Score = e.Score / numFields
			}
			for i, key := range keys {
				result[i] = engagements[key.id]
			}
			return result, nil
		},
	}
	return ldrs
}
