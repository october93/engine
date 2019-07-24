package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type CoinTransactionType string

const (
	// Initial balance
	CoinTransactionType_InitialBalance CoinTransactionType = "initialBalance"

	// Invite transactions
	// use a code when signing up
	CoinTransactionType_UsedInvite CoinTransactionType = "usedInvite"
	// someone uses your code to sign up
	CoinTransactionType_InviteAccepted CoinTransactionType = "inviteAccepted"

	// Post transactions
	// someone liked a post or comment you authored
	CoinTransactionType_LikeReceived CoinTransactionType = "likeReceived"
	// someone replied to a post or comment you authored
	CoinTransactionType_ReplyReceived CoinTransactionType = "replyReceived"
	// one of your posts performend well
	CoinTransactionType_FirstPostActivity CoinTransactionType = "firstPostActivity"
	// one of your posts performend well
	CoinTransactionType_PopularPost CoinTransactionType = "popularPost"
	// Someone tipped a post or comment you authored
	CoinTransactionType_TipGiven CoinTransactionType = "tipGiven"

	// Leaderboard transactions
	// placed on the leaderboard
	CoinTransactionType_LeaderboardFirst  CoinTransactionType = "leaderboardFirst"
	CoinTransactionType_LeaderboardSecond CoinTransactionType = "leaderboardSecond"
	CoinTransactionType_LeaderboardThird  CoinTransactionType = "leaderboardThird"
	CoinTransactionType_LeaderboardTopTen CoinTransactionType = "leaderboardTopTen"
	CoinTransactionType_LeaderboardRanked CoinTransactionType = "leaderboardRanked"

	// Alias transactions
	// Bought an alias to reply in someone else's thread with
	CoinTransactionType_BoughtThreadAlias CoinTransactionType = "boughtThreadAlias"
	// Bought an alias to post a new post with
	CoinTransactionType_BoughtPostAlias CoinTransactionType = "boughtPostAlias"

	// Channel transactins
	// Bought a channel
	CoinTransactionType_BoughtChannel CoinTransactionType = "boughtChannel"
)

type CoinTransaction struct {
	ID globalid.ID `db:"id"`

	// ID of user paying coins (null means system generated coins)
	SourceUserID globalid.ID `db:"source_user_id"`
	// ID of user getting coins (null means paid to system)
	RecipientUserID globalid.ID `db:"recipient_user_id"`

	CardID    globalid.ID         `db:"card_id"`
	Amount    int                 `db:"amount"`
	Type      CoinTransactionType `db:"type"`
	UpdatedAt time.Time           `db:"updated_at"`
	CreatedAt time.Time           `db:"created_at"`
}
