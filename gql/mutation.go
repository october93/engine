package gql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
	"github.com/pkg/errors"
)

func (g *GraphQL) Mutation_setUserDefaultStatus(ctx context.Context, id globalid.ID, status bool) (string, error) {
	user, err := g.Store.GetUser(id)
	if err != nil {
		return "", err
	}
	user.IsDefault = status
	return "", g.Store.SaveUser(user)
}

func (g *GraphQL) Mutation_indexAllUsers(ctx context.Context, clearFirst *bool) (string, error) {
	cF := clearFirst != nil && *clearFirst
	return "", g.Indexer.IndexAllUsers(cF)
}

func (g *GraphQL) Mutation_recalculateLeaderboard(ctx context.Context, inviteReward *int) (string, error) {
	toTime := time.Now().UTC()
	fromTime := toTime.Add(-time.Duration(24) * time.Hour)
	rankings, err := g.Store.GetLeaderboardRanksFromTransactions(fromTime, toTime)
	if err != nil {
		return "", err
	}

	rankToAssign := 1

	// assign ranks including ties
	for i := 0; i < len(rankings); i++ {
		currentRank := rankings[i]

		currentRank.Rank = int64(rankToAssign)

		// if this isn't the last rank
		if i+1 < len(rankings) {
			nextRank := rankings[i+1]

			if currentRank.CoinsEarned != nextRank.CoinsEarned {
				rankToAssign = i + 2
			}
		}
	}

	err = g.Store.ClearLeaderboardRankings()
	if err != nil {
		return "", err
	}

	for _, rank := range rankings {
		err = g.Store.SaveLeaderboardRank(rank)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (g *GraphQL) Mutation_indexAllChannels(ctx context.Context, clearFirst *bool) (string, error) {
	cF := clearFirst != nil && *clearFirst
	return "", g.Indexer.IndexAllChannels(cF)
}

func (g *GraphQL) Mutation_setCardIntroStatus(ctx context.Context, id globalid.ID, status bool) (string, error) {
	if status {
		err := g.Store.DeleteCardFromFeeds(id)
		if err != nil {
			return "", err
		}
	}
	return "", g.Store.SetIntroCardStatus(id, status)
}

func (g *GraphQL) Mutation_createChannel(ctx context.Context, channel ChannelInput) (*model.Channel, error) {
	err := model.ValidateChannelName(channel.Name)
	if err != nil {
		return nil, err
	}
	handle := strings.ToLower(channel.Name)

	existingChan, err := g.Store.GetChannelByHandle(handle)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, err
	} else if existingChan != nil {
		return nil, errors.New("channel with this name already exists")
	}

	newChan := &model.Channel{
		ID:        globalid.Next(),
		Name:      channel.Name,
		Handle:    handle,
		IsDefault: channel.IsDefault,
		Private:   channel.IsPrivate != nil && *channel.IsPrivate,
	}

	err = g.Store.SaveChannel(newChan)
	if err != nil {
		return nil, err
	}

	err = g.Indexer.IndexChannel(newChan)
	if err != nil {
		return nil, err
	}

	return newChan, nil
}

func (g *GraphQL) Mutation_updateCoinBalances(ctx context.Context, userID globalid.ID, coinBalance, temporaryCoinBalance *int) (string, error) {
	user, err := g.Store.GetUser(userID)

	if err != nil {
		return "", err
	}

	if coinBalance != nil {
		user.CoinBalance += int64(*coinBalance) * 10000
	}

	if temporaryCoinBalance != nil {
		user.TemporaryCoinBalance += int64(*temporaryCoinBalance) * 10000
	}
	fmt.Println(coinBalance, temporaryCoinBalance, user.CoinBalance, user.TemporaryCoinBalance)

	return "", g.Store.SaveUser(user)
}

func (g *GraphQL) Mutation_updateChannel(ctx context.Context, id globalid.ID, channelUpdate ChannelInput) (*model.Channel, error) {
	channel, err := g.Store.GetChannel(id)
	if err != nil {
		return nil, err
	}

	if channel.Name != channelUpdate.Name {
		err = model.ValidateChannelName(channelUpdate.Name)
		if err != nil {
			return nil, err
		}

		existingChan, cerr := g.Store.GetChannelByHandle(strings.ToLower(channelUpdate.Name))
		if cerr != nil && errors.Cause(cerr) != sql.ErrNoRows {
			return nil, cerr
		} else if existingChan != nil {
			return nil, errors.New("channel with this name already exists")
		}

		channel.Name = channelUpdate.Name
		channel.Handle = strings.ToLower(channelUpdate.Name)
		if channelUpdate.IsPrivate != nil {
			channel.Private = *channelUpdate.IsPrivate
		}
	}

	channel.IsDefault = channelUpdate.IsDefault

	err = g.Store.SaveChannel(channel)
	if err != nil {
		return nil, err
	}

	err = g.Indexer.IndexChannel(channel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (g *GraphQL) Mutation_createChannelInvite(ctx context.Context, channelID, inviterID globalid.ID) (*model.Invite, error) {
	existingInvite, err := g.Store.GetInviteForChannelAndNode(channelID, inviterID)
	if errors.Cause(err) == sql.ErrNoRows {
		newInvite, berr := model.NewInvite(inviterID)
		if berr != nil {
			return nil, berr
		}
		newInvite.ChannelID = channelID
		err = g.Store.SaveInvite(newInvite)
		if err != nil {
			return nil, err
		}
		return newInvite, nil
	} else if err != nil {
		return nil, err
	}

	return existingInvite, nil
}

func (g *GraphQL) Mutation_shadowbanCards(ctx context.Context, ids []globalid.ID) (string, error) {
	for _, id := range ids {
		err := g.Store.ShadowbanCard(id)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (g *GraphQL) Mutation_unshadowbanCards(ctx context.Context, ids []globalid.ID) (string, error) {
	for _, id := range ids {
		err := g.Store.UnshadowbanCard(id)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (g *GraphQL) Mutation_deactivateInvite(ctx context.Context, id globalid.ID) (string, error) {
	inv, err := g.Store.GetInvite(id)
	if err != nil {
		return "", err
	}

	inv.RemainingUses = 0

	return "", g.Store.SaveInvite(inv)
}

func (g *GraphQL) Mutation_shadowbanUser(ctx context.Context, id globalid.ID) (string, error) {
	user, err := g.Store.GetUser(id)

	if err != nil {
		return "", err
	}

	user.ShadowbannedAt = model.NewDBTime(time.Now().UTC())

	err = g.Store.SaveUser(user)
	if err != nil {
		return "", err
	}

	err = g.Store.ShadowbanAllCardsForUser(id)
	if err != nil {
		return "", err
	}

	err = g.Indexer.RemoveIndexForUser(user)
	if err != nil {
		return "", err
	}

	return RETURN_SUCCESS, nil
}

func (g *GraphQL) Mutation_unshadowbanUser(ctx context.Context, id globalid.ID) (string, error) {
	user, err := g.Store.GetUser(id)

	if err != nil {
		return "", err
	}

	user.ShadowbannedAt = model.NilDBTime()
	err = g.Store.SaveUser(user)
	if err != nil {
		return "", err
	}

	err = g.Indexer.IndexUser(user)
	if err != nil {
		return "", err
	}

	return RETURN_SUCCESS, nil
}

func (g *GraphQL) Mutation_createAnnouncement(ctx context.Context, announcementIn AnnouncementInput, sendPush *bool) (*model.Announcement, error) {
	cardID := announcementIn.ForCard
	fromUser := announcementIn.FromUser
	messasge := announcementIn.Message
	toUsers := announcementIn.ToUsers

	*sendPush = sendPush != nil && *sendPush

	announcement := &model.Announcement{
		ID:      globalid.Next(),
		Message: messasge,
		UserID:  fromUser,
	}

	if cardID != nil {
		announcement.CardID = *cardID
	}

	err := g.Store.SaveAnnouncement(announcement)
	if err != nil {
		return nil, err
	}

	// figure out who its to
	if announcementIn.ToEveryone {
		toUsers, err = g.Store.GetUserIDs()
		if err != nil {
			return nil, err
		}
	}

	// TODO: this should be pushed to a worker
	for _, v := range toUsers {
		notif := &model.Notification{
			UserID:   v,
			TargetID: announcement.ID,
			Type:     model.AnnouncementType,
		}

		err = g.Store.SaveNotification(notif)
		if err != nil {
			return nil, err
		}

		eN, err := g.Notifications.ExportNotification(notif)
		if err != nil {
			return nil, err
		}

		err = g.Pusher.NewNotification(context.Background(), &model.Session{}, eN)
		if err != nil {
			g.log.Error(err)
		}

		if *sendPush {
			err = g.Notifier.NotifyPush(eN)
			if err != nil {
				return nil, err
			}
		}
	}

	return announcement, nil

}
func (g *GraphQL) Mutation_deleteAnnouncement(ctx context.Context, id globalid.ID) (string, error) {
	err := g.Store.DeleteAnnouncement(id)
	if err != nil {
		return "", err
	}

	return RETURN_SUCCESS, err
}
func (g *GraphQL) Mutation_deleteFeatureSwitch(ctx context.Context, featureID globalid.ID) (string, error) {
	return RETURN_SUCCESS, g.Store.DeleteSwitch(featureID)
}

func (g *GraphQL) Mutation_createFeatureSwitch(ctx context.Context, state *string, name string) (string, error) {
	useState := "off"

	if state != nil {
		useState = *state
	}

	fs := model.FeatureSwitch{
		Name:         name,
		State:        useState,
		TestingUsers: make(map[globalid.ID]bool),
	}

	err := g.Store.SaveFeatureSwitch(&fs)

	return "", err
}
func (g *GraphQL) Mutation_setFeatureSwitchState(ctx context.Context, featureID globalid.ID, state string) (string, error) {
	switch state {
	case "off", "on", "testing":
		err := g.Store.ChangeFeatureSwitchState(featureID, state)
		return "", err
	}
	return "", errors.New("Invalid Feature State")
}
func (g *GraphQL) Mutation_createInvite(ctx context.Context, userID globalid.ID, usesAllowed *int) (*model.Invite, error) {

	invite, err := model.NewInvite(userID)
	if err != nil {
		return nil, err
	}

	if usesAllowed != nil {
		invite.RemainingUses = *usesAllowed
	}

	err = g.Store.SaveInvite(invite)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

func (g *GraphQL) Mutation_createInvitesFromTokens(ctx context.Context, userID globalid.ID, tokens []string) ([]model.Invite, error) {
	invites := make([]model.Invite, len(tokens))
	for i, token := range tokens {
		invite, err := model.NewInvite(userID)

		if err != nil {
			return nil, err
		}

		invite.Token = token

		//NOTE: This is a hardcoded parameter to account for the fact that invited users
		// pay 2x default attention to their inviter. The same is set in rpc.NewInvite()
		invite.RemainingUses = 1
		invite.HideFromUser = true

		err = g.Store.SaveInvite(invite)
		if err != nil {
			return nil, err
		}
		invites[i] = *invite
	}

	return invites, nil
}

func (g *GraphQL) Mutation_blockUser(ctx context.Context, id globalid.ID, deleteCards *bool) (string, error) {
	user, err := g.Store.GetUser(id)
	if err != nil {
		return "", err
	}

	user.BlockedAt = model.NewDBTime(time.Now().UTC())

	err = g.Store.SaveUser(user)

	if err != nil {
		return "", err
	}

	if deleteCards != nil && *deleteCards {
		err = g.Store.DeleteAllCardsForUser(id)
		if err != nil {
			return "", err
		}
	}
	m := protocol.NewMessage(rpc.Logout)

	for _, writer := range g.router.Connections().WritersByUser(id) {
		err = m.Encode(ctx, writer)

		if err != nil {
			g.log.Warn("failed to log blocked user out")
		}
		g.router.Connections().Deauthenticate(writer)
	}

	err = g.Store.DeleteSessionsForUser(id)
	if err != nil {
		return "", err
	}

	return "", g.Indexer.RemoveIndexForUser(user)
}

func (g *GraphQL) Mutation_unblockUser(ctx context.Context, id globalid.ID) (string, error) {
	user, err := g.Store.GetUser(id)
	if err != nil {
		return "", err
	}

	user.BlockedAt = model.NilDBTime()

	err = g.Store.SaveUser(user)
	if err != nil {
		return "", err
	}

	return "", g.Indexer.IndexUser(user)
}

func (g *GraphQL) Mutation_createUser(ctx context.Context, userIn UserInput) (user *model.User, err error) {
	nodeID := globalid.Next()
	user = model.NewUser(nodeID, userIn.Username, userIn.Email, userIn.Displayname)

	err = user.SetPassword(userIn.Password)
	if err != nil {
		return nil, err
	}

	if userIn.ProfilePictureURL != nil {
		var url string
		url, _, err = g.ImageProcessor.DownloadProfileImage(*userIn.ProfilePictureURL)
		if err != nil {
			return nil, err
		}
		user.ProfileImagePath = url
	}

	if userIn.CoverPictureURL != nil {
		var url string
		url, _, err = g.ImageProcessor.DownloadProfileImage(*userIn.CoverPictureURL)
		if err != nil {
			return nil, err
		}
		user.CoverImagePath = url
	}

	err = g.Store.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (g *GraphQL) Mutation_resetPasswords(ctx context.Context, usernames []string) (string, error) {
	users, err := g.Store.GetUsersByUsernames(usernames)
	if err != nil {
		return "", err
	}
	for _, user := range users {
		req := rpc.ResetPasswordRequest{
			Params: rpc.ResetPasswordParams{Email: user.Email},
		}
		_, err = g.RPC.ResetPassword(ctx, req)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (g *GraphQL) Mutation_toggleFeatureForUser(ctx context.Context, username string, featurename string) (string, error) {
	feature, err := g.Store.GetSwitchByName(featurename)
	if err != nil {
		return "", err
	}

	user, err := g.Store.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	err = g.Store.ToggleUserForFeature(feature.ID, user.ID)
	if err != nil {
		return "", err
	}

	return "", err
}

func (g *GraphQL) Mutation_updateSettings(ctx context.Context, settings SettingsInput) (string, error) {
	if settings.SignupsFrozen != nil {
		g.settings.SignupsFrozen = *settings.SignupsFrozen
		err := g.Store.SaveSettings(g.settings)
		if err != nil {
			return "", err
		}
	}
	if settings.MaintenanceMode != nil {
		g.settings.MaintenanceMode = *settings.MaintenanceMode
		err := g.Store.SaveSettings(g.settings)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (g *GraphQL) Mutation_updateWaitlist(ctx context.Context, comment string, email string) (string, error) {
	return "", g.Store.UpdateWaitlistEntry(email, comment)
}

func (g *GraphQL) Mutation_sendTestPush(ctx context.Context, userID globalid.ID, forCardID *globalid.ID, message, typ string, action, actionData *string) (string, error) {

	eN := &model.ExportedNotification{
		ID:        globalid.Next(),
		UserID:    userID,
		ImagePath: g.config.TestPushNotificationIconPath,
		Message:   message,
		Timestamp: time.Now().UTC().Unix(),
		Seen:      false,
		Opened:    false,
		Type:      typ,
	}

	if forCardID != nil {
		eN.ShowOnCardID = *forCardID
	}

	if action != nil {
		eN.Action = *action
		if actionData != nil {
			var actionDataMap map[string]string
			err := json.Unmarshal([]byte(*actionData), &actionDataMap)

			if err != nil {
				return "", err
			}

			eN.ActionData = actionDataMap
		}
	}

	return "", g.Notifier.NotifyPush(eN)
}

func (g *GraphQL) Mutation_generateExampleFeed(ctx context.Context, limitToChannels []globalid.ID, limitToLastNHours *int) ([]globalid.ID, error) {
	chans := limitToChannels
	var err error
	if len(chans) == 0 {
		chans, err = g.Store.GetDefaultChannelIDs()
		if err != nil {
			return nil, err
		}
	}

	postedAfter := time.Now().UTC()
	if limitToLastNHours != nil {
		postedAfter = postedAfter.Add(-(time.Duration(*limitToLastNHours) * time.Hour))
	} else {
		postedAfter = postedAfter.Add(-24 * time.Hour)
	}

	cardRanks, err := g.Store.GetRankableCardsForExampleFeed(chans, postedAfter)
	if err != nil {
		return nil, err
	}

	chosenCards, _, err := rpc.ChooseCardsFromPool(cardRanks)
	return chosenCards, err

}

func (g *GraphQL) Mutation_previewUserFeed(ctx context.Context, userID globalid.ID) ([]globalid.ID, error) {
	cardRanks, err := g.Store.GetRankableCardsForUser(userID)
	if err != nil {
		return nil, err
	}

	chosenCards, _, err := rpc.ChooseCardsFromPool(cardRanks)

	return chosenCards, err
}
func (g *GraphQL) Mutation_getCardConfidenceData(ctx context.Context, userID globalid.ID) ([]model.ConfidenceData, error) {
	cardRanks, err := g.Store.GetRankableCardsForUser(userID)
	if err != nil {
		return nil, err
	}

	sort.Slice(cardRanks, func(i, j int) bool {
		return cardRanks[i].Rank() > cardRanks[j].Rank()
	})

	ret := make([]model.ConfidenceData, len(cardRanks))
	for i, n := range cardRanks {
		ret[i] = rpc.ConfidenceFromRankEntry(n)
	}
	return ret, err

}
