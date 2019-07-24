package rpc

import (
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

//For Reactions
func (r *rpc) tokensForNewUser(userID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(globalid.Nil, userID, globalid.Nil, model.CoinTransactionType_InitialBalance)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) tokensForInvite(userID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(globalid.Nil, userID, globalid.Nil, model.CoinTransactionType_InviteAccepted)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) updateTokensForRedeemCode(userID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(globalid.Nil, userID, globalid.Nil, model.CoinTransactionType_UsedInvite)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) userCanAffordThreadAlias(userID globalid.ID) bool {
	if r.coinmanager != nil {
		return r.coinmanager.ValidateTransaction(userID, globalid.Nil, model.CoinTransactionType_BoughtThreadAlias) == nil
	}
	return true
}

func (r *rpc) userCanAffordPostAlias(userID globalid.ID) bool {
	if r.coinmanager != nil {
		return r.coinmanager.ValidateTransaction(userID, globalid.Nil, model.CoinTransactionType_BoughtPostAlias) == nil
	}
	return true
}

func (r *rpc) updateTokensForBuyThreadAlias(userID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(userID, globalid.Nil, globalid.Nil, model.CoinTransactionType_BoughtThreadAlias)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) updateTokensForBuyPostAlias(userID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(userID, globalid.Nil, globalid.Nil, model.CoinTransactionType_BoughtPostAlias)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) userCanAffordChannel(userID globalid.ID) bool {
	if r.coinmanager != nil {
		return r.coinmanager.ValidateTransaction(userID, globalid.Nil, model.CoinTransactionType_BoughtChannel) == nil
	}
	return true
}

func (r *rpc) updateTokensForBuyChannel(userID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(userID, globalid.Nil, globalid.Nil, model.CoinTransactionType_BoughtChannel)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) updateTokensForReceiveLike(userID, cardID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(globalid.Nil, userID, cardID, model.CoinTransactionType_LikeReceived)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) updateTokensForReceiveComment(userID, cardID globalid.ID) (*model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransaction(globalid.Nil, userID, cardID, model.CoinTransactionType_ReplyReceived)
		if err != nil {
			return nil, err
		}
	}
	return r.store.GetCurrentBalance(userID)
}

func (r *rpc) userCanAffordTip(userID globalid.ID, amount int) bool {
	if r.coinmanager != nil {
		return r.coinmanager.ValidateTransactionWithAmount(userID, globalid.Nil, model.CoinTransactionType_TipGiven, amount) == nil
	}
	return true
}

func (r *rpc) updateTokensForTip(userID, tippedCardID, tippedCardOwnerID globalid.ID, amount int) (*model.CoinBalances, *model.CoinBalances, error) {
	if r.coinmanager != nil {
		err := r.coinmanager.ProcessTransactionWithAmount(userID, tippedCardOwnerID, tippedCardID, model.CoinTransactionType_TipGiven, amount)
		if err != nil {
			return nil, nil, err
		}
	}
	tipperBalance, err := r.store.GetCurrentBalance(userID)
	if err != nil {
		return nil, nil, err
	}
	tippedBalance, err := r.store.GetCurrentBalance(tippedCardOwnerID)
	if err != nil {
		return nil, nil, err
	}
	return tipperBalance, tippedBalance, nil

}
