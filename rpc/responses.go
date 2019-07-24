package rpc

import (
	"fmt"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// Responses is responsible for building complex responses. For instance card
// responses require several queries to build the output needed for the
// clients.
type Responses struct {
	store dataStore
}

type CardData struct {
	engagements          map[globalid.ID]*model.Engagement
	authors              map[globalid.ID]*model.Author
	channels             map[globalid.ID]*model.Channel
	featuredComments     map[globalid.ID]*model.Card
	userReactions        map[globalid.ID]*model.UserReaction
	threadCounts         map[globalid.ID]int
	subscribedToUsers    map[globalid.ID]bool
	channelSubscriptions map[globalid.ID]bool
	subscribedToCards    map[globalid.ID]bool
	newContentAvailable  map[globalid.ID]bool
}

func NewResponses(store dataStore) *Responses {
	return &Responses{
		store: store,
	}
}

func (r *Responses) FeedCardResponses(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error) {
	cardIDs := make([]globalid.ID, len(cards))
	channelIDs := make([]globalid.ID, len(cards))
	ownerIDs := make([]globalid.ID, len(cards))
	for i, card := range cards {
		cardIDs[i] = card.ID
		channelIDs[i] = card.ChannelID
		ownerIDs[i] = card.OwnerID
	}
	engagements, err := r.store.GetEngagements(cardIDs)
	if err != nil {
		return nil, err
	}
	authors, err := r.authors(cards)
	if err != nil {
		return nil, err
	}
	channels, err := r.store.GetChannelsByID(channelIDs)
	if err != nil {
		return nil, err
	}
	channelsByID := make(map[globalid.ID]*model.Channel)
	for _, channel := range channels {
		channelsByID[channel.ID] = channel
	}
	featuredComments, err := r.store.FeaturedCommentsForUserByCardIDs(viewerID, cardIDs)
	if err != nil {
		return nil, err
	}
	userReactions, err := r.store.GetUserReactions(viewerID, cardIDs)
	if err != nil {
		return nil, err
	}
	threadCounts, err := r.store.GetThreadCounts(cardIDs)
	if err != nil {
		return nil, err
	}
	followings, err := r.store.IsFollowings(viewerID, ownerIDs)
	if err != nil {
		return nil, err
	}
	subscriptions, err := r.store.IsSubscribedToChannels(viewerID, channelIDs)
	if err != nil {
		return nil, err
	}
	subscribedToCards, err := r.store.SubscribedToCards(viewerID, cardIDs)
	if err != nil {
		return nil, err
	}
	newContentAvailable, err := r.store.NewContentAvailableForUserByCards(viewerID, cardIDs)
	if err != nil {
		return nil, err
	}

	data := &CardData{
		engagements:          engagements,
		authors:              authors,
		channels:             channelsByID,
		featuredComments:     featuredComments,
		userReactions:        userReactions,
		threadCounts:         threadCounts,
		subscribedToUsers:    followings,
		channelSubscriptions: subscriptions,
		subscribedToCards:    subscribedToCards,
		newContentAvailable:  newContentAvailable,
	}

	result := make([]*model.CardResponse, len(cards))
	for i, card := range cards {
		cRsp, err := r.CardResponse(card, viewerID, data)
		if err != nil {
			return nil, err
		}
		featuredCommentsAreNew := data.newContentAvailable[card.ID]
		featuredCommentCard := data.featuredComments[card.ID]
		if featuredCommentCard != nil {
			fC, gerr := featuredComment(r.store, featuredCommentCard, featuredCommentsAreNew)
			if gerr != nil {
				return nil, gerr
			}
			cRsp.FeaturedComments = &model.FeaturedComments{
				Comments: []*model.FeaturedComment{fC},
			}
			if err != nil {
				return nil, err
			}
		}
		result[i] = cRsp
	}

	return result, nil
}

func (r *Responses) authors(cards []*model.Card) (map[globalid.ID]*model.Author, error) {
	userIDs := make([]globalid.ID, len(cards))
	aliasIDs := make([]globalid.ID, len(cards))
	for i, card := range cards {
		userIDs[i] = card.OwnerID
		aliasIDs[i] = card.AliasID
	}

	users, err := r.store.GetUsersByID(userIDs)
	if err != nil {
		return nil, err
	}
	aliases, err := r.store.GetAnonymousAliasesByID(aliasIDs)
	if err != nil {
		return nil, err
	}

	usersByID := make(map[globalid.ID]*model.User, len(users))
	for _, user := range users {
		usersByID[user.ID] = user
	}
	aliasesByID := make(map[globalid.ID]*model.AnonymousAlias, len(aliases))
	for _, alias := range aliases {
		aliasesByID[alias.ID] = alias
	}

	result := make(map[globalid.ID]*model.Author, len(cards))
	for _, card := range cards {
		if card.AliasID != globalid.Nil {
			result[card.ID] = aliasesByID[card.AliasID].Author()
		} else {
			result[card.ID] = usersByID[card.OwnerID].Author()
		}

	}
	return result, nil
}

func (r *Responses) CardResponses() ([]*model.CardResponse, error) {
	return nil, nil
}

func (r *Responses) CardResponse(card *model.Card, viewerID globalid.ID, data *CardData) (*model.CardResponse, error) {
	author := data.authors[card.ID]
	replies := data.threadCounts[card.ID]

	var reaction *model.Reaction
	var voteResponse *model.VoteResponse
	var userReaction *model.UserReaction

	if viewerID != globalid.Nil {
		userReaction = data.userReactions[card.ID]
		if userReaction != nil {
			// if it errors out all 3 reactions fields should be nil
			reaction = userReaction.ToCardReaction()
			voteResponse = userReaction.ToVoteResponse()
		}
	}

	var err error
	isMine := viewerID == card.OwnerID
	threadRoot := card
	if card.ThreadRootID != globalid.Nil {
		threadRoot, err = r.store.GetCard(card.ThreadRootID)
		if err != nil {
			return nil, err
		}
	}

	var viewer *model.Viewer
	if threadRoot.AuthorToAlias[viewerID] != globalid.Nil {
		var anonymousAlias *model.AnonymousAlias
		viewer = &model.Viewer{}
		anonymousAlias, err = r.store.GetAnonymousAlias(threadRoot.AuthorToAlias[viewerID])
		if err != nil {
			return nil, err
		}
		var lastUsed bool
		lastUsed, err = r.store.GetAnonymousAliasLastUsed(viewerID, threadRoot.ID)
		if err != nil {
			return nil, err
		}
		viewer.AnonymousAlias = anonymousAlias
		viewer.AnonymousAliasLastUsed = lastUsed
	}

	engagement := data.engagements[card.ID]
	var channel *model.Channel
	if card.ChannelID != globalid.Nil {
		channel = data.channels[card.ChannelID]
	}

	subscribedToUser := data.subscribedToUsers[card.OwnerID]
	subscribedToChannel := data.channelSubscriptions[card.ChannelID]

	var rankingReason string
	if subscribedToUser && !subscribedToChannel && card.AliasID == globalid.Nil {
		rankingReason = fmt.Sprintf("Because you follow **@%v**", author.Username)
	}

	subscribed := data.subscribedToCards[card.ID]

	return &model.CardResponse{
		Card:           card.Export(),
		Author:         author,
		Viewer:         viewer,
		Channel:        channel,
		Replies:        replies,
		Reactions:      reaction,
		Engagement:     engagement,
		ViewerReaction: userReaction,
		Score:          0,
		Subscribed:     subscribed,
		IsMine:         isMine,
		Vote:           voteResponse,
		RankingReason:  rankingReason,
	}, nil
}
