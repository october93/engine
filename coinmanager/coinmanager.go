package coinmanager

import (
	"errors"

	"github.com/lib/pq"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store"
)

var ErrNoAmount = errors.New("transaction requiring amount has no amount")
var ErrInsufficientFunds = errors.New("insufficient funds for transaction")

type CoinManager struct {
	Store *store.Store

	Config *Config
}

// NewGraphQL returns a new instance of GraphQL.
func NewCoinManager(s *store.Store, c *Config) *CoinManager {
	return &CoinManager{
		Store:  s,
		Config: c,
	}
}

func (c *CoinManager) GetAmountForTransaction(typ model.CoinTransactionType) (int, error) {
	switch typ {
	case model.CoinTransactionType_InitialBalance:
		return c.Config.InitialBalance, nil
	case model.CoinTransactionType_UsedInvite:
		return c.Config.UsedInvite, nil
	case model.CoinTransactionType_InviteAccepted:
		return c.Config.InviteAccepted, nil
	case model.CoinTransactionType_LikeReceived:
		return c.Config.LikeReceived, nil
	case model.CoinTransactionType_ReplyReceived:
		return c.Config.ReplyReceived, nil
	case model.CoinTransactionType_FirstPostActivity:
		return c.Config.FirstPostActivity, nil
	case model.CoinTransactionType_PopularPost:
		return c.Config.PopularPost, nil
	case model.CoinTransactionType_LeaderboardFirst:
		return c.Config.LeaderboardFirst, nil
	case model.CoinTransactionType_LeaderboardSecond:
		return c.Config.LeaderboardSecond, nil
	case model.CoinTransactionType_LeaderboardThird:
		return c.Config.LeaderboardThird, nil
	case model.CoinTransactionType_LeaderboardTopTen:
		return c.Config.LeaderboardTopTen, nil
	case model.CoinTransactionType_LeaderboardRanked:
		return c.Config.LeaderboardRanked, nil
	case model.CoinTransactionType_BoughtThreadAlias:
		return c.Config.BoughtThreadAlias, nil
	case model.CoinTransactionType_BoughtPostAlias:
		return c.Config.BoughtPostAlias, nil
	case model.CoinTransactionType_BoughtChannel:
		return c.Config.BoughtChannel, nil
	}
	return 0, ErrNoAmount
}

func (c *CoinManager) ValidateTransaction(sourceUserID, recipientUserID globalid.ID, typ model.CoinTransactionType) error {
	// if source user is nil, the system provides the coins and is always valid
	if sourceUserID == globalid.Nil {
		return nil
	}

	amount, err := c.GetAmountForTransaction(typ)
	if err != nil {
		return err
	}

	user, err := c.Store.GetUser(sourceUserID)
	if err != nil {
		return err
	}

	if int(user.CoinBalance) < amount {
		return ErrInsufficientFunds
	}
	return nil
}

func (c *CoinManager) ValidateTransactionWithAmount(sourceUserID, recipientUserID globalid.ID, typ model.CoinTransactionType, amount int) error {
	// if source user is nil, the system provides the coins and is always valid
	if sourceUserID == globalid.Nil {
		return nil
	}

	user, err := c.Store.GetUser(sourceUserID)
	if err != nil {
		return err
	}

	if int(user.CoinBalance) < amount {
		return ErrInsufficientFunds
	}
	return nil
}

func (c *CoinManager) ProcessTransaction(sourceUserID, recipientUserID, cardID globalid.ID, typ model.CoinTransactionType) error {
	amount, err := c.GetAmountForTransaction(typ)
	if err != nil {
		return err
	}

	return c.SaveTransaction(sourceUserID, recipientUserID, cardID, typ, amount)
}

func (c *CoinManager) ProcessTransactionWithAmount(sourceUserID, recipientUserID, cardID globalid.ID, typ model.CoinTransactionType, amount int) error {
	return c.SaveTransaction(sourceUserID, recipientUserID, cardID, typ, amount)
}

func (c *CoinManager) SaveTransaction(sourceUserID, recipientUserID, cardID globalid.ID, typ model.CoinTransactionType, amount int) error {
	if sourceUserID != globalid.Nil {
		// charge the source
		err := c.Store.SubtractCoinsFromBalance(sourceUserID, amount)

		// return the right error if it's insufficient funds
		if pgerr, ok := err.(*pq.Error); ok && pgerr.Code == "23000" {
			return ErrInsufficientFunds
		} else if err != nil {
			return err
		}
	}

	if recipientUserID != globalid.Nil {
		// reward the recipient
		err := c.Store.AddCoinsToBalance(recipientUserID, amount)
		if err != nil {
			return err
		}
	}

	if cardID != globalid.Nil {
		// reward the recipient
		err := c.Store.AddToCoinsEarned(cardID, amount)
		if err != nil {
			return err
		}
	}

	// save the transaction since it was successful
	t := &model.CoinTransaction{
		SourceUserID:    sourceUserID,
		RecipientUserID: recipientUserID,
		CardID:          cardID,
		Amount:          amount,
		Type:            typ,
	}

	return c.Store.SaveCoinTransaction(t)
}
