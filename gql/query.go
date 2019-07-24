package gql

import (
	"context"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc/protocol"
)

func (g *GraphQL) Query_sessions(ctx context.Context) ([]model.Session, error) {
	sessions, err := g.Store.GetSessions()
	if err != nil {
		return nil, err
	}
	ret := make([]model.Session, len(sessions))
	for i, session := range sessions {
		ret[i] = *session
	}
	return ret, nil
}

func (g *GraphQL) Query_settings(ctx context.Context) (*model.Settings, error) {
	return g.Store.GetSettings()
}

func (g *GraphQL) Query_channelEngagements(ctx context.Context) ([]model.ChannelEngagement, error) {
	a, err := g.Store.GetChannelEngagements()
	if err != nil {
		return nil, err
	}
	ret := make([]model.ChannelEngagement, len(a))
	for i, item := range a {
		ret[i] = *item
	}
	return ret, nil
}

func (g *GraphQL) Channel_isPrivate(ctx context.Context, obj *model.Channel) (bool, error) {
	return obj.Private, nil
}

func (g *GraphQL) Query_channels(ctx context.Context) ([]model.Channel, error) {
	a, err := g.Store.GetChannels()
	if err != nil {
		return nil, err
	}
	ret := make([]model.Channel, len(a))
	for i, item := range a {
		ret[i] = *item
	}
	return ret, nil
}

func (g *GraphQL) Query_invites(ctx context.Context) ([]model.Invite, error) {
	a, err := g.Store.GetInvites()
	if err != nil {
		return nil, err
	}
	ret := make([]model.Invite, len(a))
	for i, item := range a {
		ret[i] = *item
	}
	return ret, nil
}
func (g *GraphQL) Query_announcements(ctx context.Context) ([]model.Announcement, error) {
	a, err := g.Store.GetAnnouncements()
	if err != nil {
		return nil, err
	}
	ret := make([]model.Announcement, len(a))
	for i, item := range a {
		ret[i] = *item
	}
	return ret, nil
}
func (g *GraphQL) Query_featureSwitches(ctx context.Context) ([]model.FeatureSwitch, error) {
	a, err := g.Store.GetFeatureSwitches()
	if err != nil {
		return nil, err
	}
	ret := make([]model.FeatureSwitch, len(a))
	for i, item := range a {
		ret[i] = *item
	}
	return ret, nil
}
func (g *GraphQL) Query_waitlist(ctx context.Context) ([]model.WaitlistEntry, error) {
	a, err := g.Store.GetWaitlist()
	if err != nil {
		return nil, err
	}
	ret := make([]model.WaitlistEntry, len(a))
	for i, item := range a {
		ret[i] = *item
	}
	return ret, nil
}
func (g *GraphQL) Query_connections(ctx context.Context) ([]protocol.Connection, error) {
	writers := g.router.Connections().Writers()
	conns := make([]protocol.Connection, len(writers))
	for i, writer := range writers {
		conns[i] = *writer.Connection()
	}
	return conns, nil
}

func (g *GraphQL) Query_cardEngagement(ctx context.Context, from *string, to *string) ([]CardEngagement, error) {
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
	uucc, err := g.Store.GetUniqueUserCommentCount(startDate, endDate)
	if err != nil {
		return nil, err
	}

	trc, err := g.Store.GetTotalReplyCount(startDate, endDate)
	if err != nil {
		return nil, err
	}
	tlc, err := g.Store.GetTotalLikeCount(startDate, endDate)
	if err != nil {
		return nil, err
	}
	tdc, err := g.Store.GetTotalDislikeCount(startDate, endDate)
	if err != nil {
		return nil, err
	}
	cardEngagementByID := make(map[globalid.ID]*CardEngagement)
	for _, c := range uucc {
		ce, ok := cardEngagementByID[c.UserID]
		if !ok {
			ce = &CardEngagement{ID: c.UserID}
			cardEngagementByID[c.UserID] = ce
		}
		ce.UniqueUserCommentCount = c.Count
	}
	for _, c := range trc {
		ce, ok := cardEngagementByID[c.UserID]
		if !ok {
			ce = &CardEngagement{ID: c.UserID}
			cardEngagementByID[c.UserID] = ce
		}
		ce.TotalReplyCount = c.Count
	}
	for _, c := range tdc {
		ce, ok := cardEngagementByID[c.UserID]
		if !ok {
			ce = &CardEngagement{ID: c.UserID}
			cardEngagementByID[c.UserID] = ce
		}
		ce.TotalDislikeCount = c.Count
	}
	for _, c := range tlc {
		ce, ok := cardEngagementByID[c.UserID]
		if !ok {
			ce = &CardEngagement{ID: c.UserID}
			cardEngagementByID[c.UserID] = ce
		}
		ce.TotalLikeCount = c.Count
	}
	cardEngagement := make([]CardEngagement, 0)
	for _, ce := range cardEngagementByID {
		cardEngagement = append(cardEngagement, *ce)
	}
	return cardEngagement, nil
}
