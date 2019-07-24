package gql

import (
	"context"
	"fmt"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (g *GraphQL) Query_cards(ctx context.Context, from *time.Time, to *time.Time, ids []globalid.ID, introCards *bool) (cards []model.Card, err error) {
	var cardPtrs []*model.Card
	if from != nil || to != nil {
		start := time.Time{}
		end := time.Time{}
		if from != nil {
			start = *from
		}
		if to != nil {
			end = *to
		}
		cardPtrs, err = g.Store.GetCardsByInterval(start, end)
	} else if ids != nil {
		cardPtrs, err = g.Store.GetCardsWithDeletedByID(ids)
	} else if introCards != nil && *introCards {
		cardPtrs, err = g.Store.GetIntroCards()
	} else {
		cardPtrs, err = g.Store.GetCards()
	}
	if err != nil {
		return nil, err
	}
	cards = make([]model.Card, len(cardPtrs))
	for i, card := range cardPtrs {
		cards[i] = *card
	}
	return
}

func (g *GraphQL) Card_apparentAuthor(ctx context.Context, card *model.Card) (author *model.Author, err error) {
	if card.AliasID != globalid.Nil {
		anonymousAlias, err := ctxLoaders(ctx).AliasByID.Load(card.AliasID)
		if err != nil {
			return nil, err
		}
		if anonymousAlias == nil {
			return nil, fmt.Errorf("no error, but anon author is nil %v", card.AliasID)
		}
		author = anonymousAlias.Author()
	} else if card.OwnerID != globalid.Nil {
		user, err := ctxLoaders(ctx).UserByID.Load(card.OwnerID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, fmt.Errorf("no error, but author is nil %v", card.OwnerID)
		}
		author = user.Author()
	} else {
		return nil, fmt.Errorf("card %v has no apparent author (anonId: %v, ownerId: %v)", card.ID, card.AliasID, card.OwnerID)
	}
	return
}
