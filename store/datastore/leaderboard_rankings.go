package datastore

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/october93/engine/model"
)

func (store *Store) ClearLeaderboardRankings() error {
	_, err := store.Exec(`DELETE FROM leaderboard_rankings`)
	return errors.Wrap(err, "ClearLeaderboardRankings failed")
}

func (store *Store) SaveLeaderboardRank(m *model.LeaderboardRank) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}

	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}

	_, err := store.NamedExec(
		`INSERT INTO leaderboard_rankings
		(
			user_id,
			rank,
			coins_earned,
			created_at
		)
		VALUES
		(
			:user_id,
			:rank,
			:coins_earned,
			:created_at
		)
		ON CONFLICT(user_id) DO UPDATE
		SET
			user_id = :user_id,
			rank = :rank,
			coins_earned = :coins_earned,
			created_at  = :created_at
		WHERE leaderboard_rankings.user_id = :user_id`, m)
	return errors.Wrap(err, "SaveInvite failed")
}

func (store *Store) GetLeaderboardRankings(count, skip int) ([]*model.LeaderboardRank, error) {
	ranks := []*model.LeaderboardRank{}
	err := store.Select(&ranks, "SELECT * FROM leaderboard_rankings ORDER BY rank, user_id LIMIT $1 OFFSET $2", count, count*skip)
	if err != nil {
		return nil, errors.Wrap(err, "GetRankings failed")
	}
	return ranks, nil
}

func (store *Store) GetLeaderboardRanksFromTransactions(fromTime, toTime time.Time) ([]*model.LeaderboardRank, error) {
	data := []*model.LeaderboardRank{}
	excludeTransactionTypes := []model.CoinTransactionType{
		model.CoinTransactionType_LeaderboardFirst,
		model.CoinTransactionType_LeaderboardSecond,
		model.CoinTransactionType_LeaderboardThird,
		model.CoinTransactionType_LeaderboardTopTen,
		model.CoinTransactionType_LeaderboardRanked,
		model.CoinTransactionType_InitialBalance,
	}

	query, args, err := sqlx.In(`
		SELECT
	    recipient_user_id as user_id,
			0 AS rank,
	    SUM(amount) as coins_earned,
	    now() AS created_at
		FROM
			coin_transactions
			LEFT JOIN users ON coin_transactions.recipient_user_id = users.id
			LEFT JOIN cards ON coin_transactions.card_id = cards.id
		WHERE
			coin_transactions.created_at >= ? AND coin_transactions.created_at <= ?
			AND
			type NOT IN (?)
			AND
			users.is_internal = false
			AND
			(coin_transactions.source_user_id != coin_transactions.recipient_user_id OR coin_transactions.source_user_id IS NULL OR coin_transactions.source_user_id IS NULL)
			AND
			cards.alias_id IS NULL
		GROUP BY recipient_user_id
		ORDER BY coins_earned DESC
		`, fromTime, toTime, excludeTransactionTypes)
	if err != nil {
		return nil, err
	}
	query = store.Rebind(query)
	err = store.Select(&data, query, args...)

	return data, errors.Wrap(err, "GetLeaderboardRanksFromTransactions failed")
}
