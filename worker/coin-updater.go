package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/october93/engine/coinmanager"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	rpccontext "github.com/october93/engine/rpc/context"
	"github.com/october93/engine/rpc/notifications"
	"github.com/october93/engine/store"
)

type pusher interface {
	NewNotification(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error
	UpdateCoinBalance(ctx context.Context, userID globalid.ID, newBalances *model.CoinBalances) error
}

type CoinUpdater struct {
	store                *store.Store
	notificationExporter *notifications.Notifications
	notifier             *Notifier
	pusher               pusher
	log                  log.Logger
	shutdown             chan struct{}

	coinManager *coinmanager.CoinManager

	config *CoinUpdaterConfig
}

func NewCoinUpdater(store *store.Store, log log.Logger, nE *notifications.Notifications, n *Notifier, p pusher, config CoinUpdaterConfig, cm *coinmanager.CoinManager) *CoinUpdater {
	return &CoinUpdater{
		store:                store,
		notificationExporter: nE,
		notifier:             n,
		pusher:               p,
		log:                  log,
		shutdown:             make(chan struct{}),
		coinManager:          cm,
		config:               &config,
	}
}

func (c *CoinUpdater) Run() {
	c.updateLeaderboard()
}

func (c *CoinUpdater) RunPopularPostUpdates() {
	c.givePopularCardRewards()
}

func (c *CoinUpdater) Stop() {
	close(c.shutdown)
}

// Find posts which:
// Is the first post/comment by a user that received >= 1 like or comment
// Has received >= 5 likes or replies in the last 24h
func (c *CoinUpdater) givePopularCardRewards() {
	// update leaderboard
	ctx := context.WithValue(context.Background(), rpccontext.RequestID, globalid.Next()) //nolint

	// Get first post/comment by a user that received >= 1 like or comment
	firstPostUsers, err := c.store.GetUsersNeedingFirstPostNotification()
	if err != nil {
		c.log.Error(err)
		return
	}
	fmt.Println(firstPostUsers)
	// Get cards with >= 5 likes or replies in the last 24h
	popularPostUsers, err := c.store.UsersWithPopularPostsSinceTime(time.Now().UTC().Add(-24 * time.Hour))
	if err != nil {
		c.log.Error(err)
		return
	}
	fmt.Println(popularPostUsers)
	notifs := make(map[globalid.ID]*model.Notification)

	for _, user := range popularPostUsers {
		// Give coins
		err = c.coinManager.ProcessTransaction(globalid.Nil, user.ID, globalid.Nil, model.CoinTransactionType_PopularPost)
		if err != nil {
			c.log.Error(err)
			continue
		}

		cB, berr := c.store.GetCurrentBalance(user.ID)
		if berr != nil {
			c.log.Error(berr)
		}

		if cB != nil {
			err = c.pusher.UpdateCoinBalance(ctx, user.ID, cB)
			if err != nil {
				c.log.Error(err)
			}
		}

		n := &model.Notification{
			UserID: user.ID,
			Type:   model.PopularPostType,
		}

		err = c.store.SaveNotification(n)
		if err != nil {
			c.log.Error(err)
			continue
		}

		notifs[user.ID] = n
	}

	for _, user := range firstPostUsers {
		_, ok := notifs[user.ID]

		if !ok {
			// Give coins
			err = c.coinManager.ProcessTransaction(globalid.Nil, user.ID, globalid.Nil, model.CoinTransactionType_FirstPostActivity)
			if err != nil {
				c.log.Error(err)
				continue
			}

			n := &model.Notification{
				UserID: user.ID,
				Type:   model.FirstPostActivityType,
			}

			err = c.store.SaveNotification(n)
			if err != nil {
				c.log.Error(err)
				continue
			}

			notifs[user.ID] = n
		}
	}

	for _, notif := range notifs {
		fmt.Println(notif)
		// export the notification
		exNotif, err := c.notificationExporter.ExportNotification(notif)
		if err != nil {
			c.log.Error(err)
			continue
		}

		// push new notif to open clients
		err = c.pusher.NewNotification(ctx, nil, exNotif)
		if err != nil {
			c.log.Error(err)
			continue
		}

		// notify via push
		err = c.notifier.NotifyPush(exNotif)
		if err != nil {
			c.log.Error(err)
		}
	}
}

func (c *CoinUpdater) updateLeaderboard() {
	// update leaderboard
	ctx := context.WithValue(context.Background(), rpccontext.RequestID, globalid.Next()) //nolint
	updateFrequency := int64(24)
	if c.config.LeaderboardUpdateFreqency > 0 {
		updateFrequency = c.config.LeaderboardUpdateFreqency
	}
	toTime := time.Now().UTC()
	fromTime := toTime.Add(-time.Duration(updateFrequency) * time.Hour)

	leaderboardRanks, err := c.store.GetLeaderboardRanksFromTransactions(fromTime, toTime)
	if err != nil {
		c.log.Error(err)
		return
	}

	rankToAssign := 1

	// assign ranks including ties up to
	placeLimit := 25
	if c.config.LeaderboardPlacesLimit > 0 {
		placeLimit = c.config.LeaderboardPlacesLimit
	}

	rankedLeaderboardRanks := make([]*model.LeaderboardRank, 0)

	for i := 0; i < len(leaderboardRanks) && rankToAssign <= placeLimit; i++ {
		currentRank := leaderboardRanks[i]

		currentRank.Rank = int64(rankToAssign)
		rankedLeaderboardRanks = append(rankedLeaderboardRanks, currentRank)

		// if this isn't the last rank
		if i+1 < len(leaderboardRanks) {
			nextRank := leaderboardRanks[i+1]

			if currentRank.CoinsEarned != nextRank.CoinsEarned {
				rankToAssign = i + 2
			}
		}
	}

	// clear existing rankings
	err = c.store.ClearLeaderboardRankings()
	if err != nil {
		c.log.Error(err)
	}

	// save ranks and send notifs
	for _, rank := range rankedLeaderboardRanks {
		// save the rank
		err = c.store.SaveLeaderboardRank(rank)
		if err != nil {
			c.log.Error(err)
			continue
		}

		// Save the notifications
		notif := &model.Notification{
			ID:     globalid.Next(),
			UserID: rank.UserID,
			Type:   model.LeaderboardRankType,
		}

		err = c.store.SaveNotification(notif)
		if err != nil {
			c.log.Error(err)
			continue
		}

		err = c.store.SaveLeaderboardNotificationData(notif.ID, int(rank.Rank))
		if err != nil {
			c.log.Error(err)
			continue
		}

		if notif != nil {
			// export the notification
			exNotif, verr := c.notificationExporter.ExportNotification(notif)
			if verr != nil {
				c.log.Error(verr)
			}

			// push new notif to open clients
			err = c.pusher.NewNotification(ctx, nil, exNotif)
			if err != nil {
				c.log.Error(err)
			}

			// notify via push
			err = c.notifier.NotifyPush(exNotif)
			if err != nil {
				c.log.Error(err)
			}
		}

		// Give coin award
		var rewardTyp model.CoinTransactionType
		switch rank.Rank {
		case 1:
			rewardTyp = model.CoinTransactionType_LeaderboardFirst
		case 2:
			rewardTyp = model.CoinTransactionType_LeaderboardSecond
		case 3:
			rewardTyp = model.CoinTransactionType_LeaderboardThird
		default:
			if rank.Rank > 4 && rank.Rank >= 10 {
				rewardTyp = model.CoinTransactionType_LeaderboardTopTen
			} else {
				rewardTyp = model.CoinTransactionType_LeaderboardRanked
			}
		}

		err = c.coinManager.ProcessTransaction(globalid.Nil, rank.UserID, globalid.Nil, rewardTyp)
		if err != nil {
			c.log.Error(err)
		}

		cB, err := c.store.GetCurrentBalance(rank.UserID)
		if err != nil {
			c.log.Error(err)
		}

		if cB != nil {
			err := c.pusher.UpdateCoinBalance(ctx, rank.UserID, cB)
			if err != nil {
				c.log.Error(err)
			}
		}
	}
}
