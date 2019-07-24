package rpc_test

import (
	"context"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/worker/emailsender"
)

type testStore struct {
	DeleteSessionf                      func(globalid.ID) error
	SaveCardf                           func(*model.Card) error
	GetCardf                            func(globalid.ID) (*model.Card, error)
	GetThreadf                          func(id, forUser globalid.ID) ([]*model.Card, error)
	GetUserByEmailf                     func(email string) (*model.User, error)
	GetUserByUsernamef                  func(username string) (*model.User, error)
	GetUsersByUsernamesf                func(username []string) ([]*model.User, error)
	GetUserf                            func(nodeID globalid.ID) (u *model.User, err error)
	GetUsersf                           func() ([]*model.User, error)
	GetUsersByIDf                       func([]globalid.ID) (res []*model.User, err error)
	SaveInvitef                         func(invite *model.Invite) error
	GetInviteByTokenf                   func(token string) (*model.Invite, error)
	DeleteInvitef                       func(id globalid.ID) error
	CountGivenReactionsForUserf         func(posterID globalid.ID, timeFrom, timeTo time.Time, onlyLikes bool) (c int, err error)
	SaveSessionf                        func(session *model.Session) error
	SaveUserf                           func(u *model.User) (err error)
	SaveResetTokenf                     func(rt *model.ResetToken) error
	GetResetTokenf                      func(userID globalid.ID) (*model.ResetToken, error)
	GetNotifcationsf                    func(userID globalid.ID, pageSize, pageNumber int) ([]*model.Notification, error)
	UpdateNotificationsSeenf            func(ids []globalid.ID) error
	UpdateNotificationsOpenedf          func(ids []globalid.ID) error
	GetAnonymousAliasf                  func(id globalid.ID) (*model.AnonymousAlias, error)
	GetAnonymousAliasesf                func() ([]*model.AnonymousAlias, error)
	GetUnusedAliasf                     func(cardID globalid.ID) (*model.AnonymousAlias, error)
	GetPostedCardsForNodef              func(nodeID globalid.ID, skip, count int) ([]*model.Card, error)
	GetPostedCardsForNodeIncludingAnonf func(nodeID globalid.ID, count, skip int) ([]*model.Card, error)
	DeleteCardf                         func(id globalid.ID) error
	SaveOAuthAccountf                   func(oaa *model.OAuthAccount) error
	GetOAuthAccountBySubjectf           func(subject string) (*model.OAuthAccount, error)
	GetOnSwitchesf                      func() ([]*model.FeatureSwitch, error)
	LatestForTypef                      func(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error)
	SaveNotificationf                   func(m *model.Notification) error
	GetAnonymousAliasByUsernamef        func(username string) (*model.AnonymousAlias, error)
	DeleteNotificationsForCardf         func(cardID globalid.ID) error
	UnseenNotificationsCountf           func(userID globalid.ID) (int, error)
	GetUserCardIDsf                     func(id globalid.ID, num int, createdAt time.Time) ([]globalid.ID, error)
	LastPostInThreadWasAnonymousf       func(userID, threadRootID globalid.ID) (bool, error)
	GetThreadCountf                     func(id globalid.ID) (int, error)
	SaveWaitlistEntryf                  func(m *model.WaitlistEntry) error
	GetInvitesForUserf                  func(id globalid.ID) ([]*model.Invite, error)
	GetInvitef                          func(id globalid.ID) (*model.Invite, error)
	GetEngagementf                      func(cardID globalid.ID) (*model.Engagement, error)
	DeleteWaitlistEntryf                func(email string) error
	GetSafeUsersByPagef                 func(forUser globalid.ID, pageSize, pageNumber int, searchString string) ([]*model.User, error)
	SaveMentionf                        func(m *model.Mention) error
	GetMentionExportDataf               func(n *model.Notification) (*datastore.MentionNotificationExportData, error)
	SaveNotificationMentionf            func(m *model.NotificationMention) error
	GetCommentExportDataf               func(n *model.Notification) (*datastore.CommentNotificationExportData, error)
	GetLikeNotificationExportDataf      func(n *model.Notification) (*datastore.LikeNotificationExportData, error)
	SaveNotificationCommentf            func(m *model.NotificationComment) error
	SaveReactionForNotificationf        func(notifID, userID, cardID globalid.ID) error
	SubscribeToCardf                    func(userID, cardID globalid.ID, typ string) error
	UnsubscribeFromCardf                func(userID, cardID globalid.ID, typ string) error
	SubscribersForCardf                 func(cardID globalid.ID, typ string) ([]globalid.ID, error)
	SubscribedToTypesf                  func(userID, cardID globalid.ID) ([]string, error)
	DeleteMentionsForCardf              func(cardID globalid.ID) error
	ReassignInviterForGroupf            func(rootInviteID, newUserID globalid.ID) error
	GroupInvitesByTokenf                func(tokens []string, groupID globalid.ID) error
	ReassignInvitesByTokenf             func(tokens []string, userID globalid.ID) error
	SaveFollowerf                       func(followerID, followeeID globalid.ID) error
	DeleteFollowerf                     func(followerID, followeeID globalid.ID) error
	GetFollowingf                       func(userID globalid.ID) ([]*model.User, error)
	ClearEmptyNotificationsf            func() error
	ClearEmptyNotificationsForUserf     func(userID globalid.ID) error
	SaveNotificationFollowf             func(notifID, followerID, followeeID globalid.ID) error
	GetFollowExportDataf                func(n *model.Notification) (*datastore.FollowNotificationExportData, error)
	IsFollowingf                        func(followerID, followeeID globalid.ID) (bool, error)
	GetNotificationf                    func(notifID globalid.ID) (*model.Notification, error)
	BlockUserf                          func(blockingUser, blockedUser globalid.ID) error
	RestorePreviousRankForUserf         func(userID, cardID globalid.ID) error
	BlockAnonUserInThreadf              func(blockingUser, blockedAlias, threadID globalid.ID) error
	GetCardsInFeedf                     func(userID globalid.ID) ([]model.FeedEntry, error)
	FeaturedCommentsForUserf            func(userID, cardID globalid.ID) ([]*model.Card, error)
	GetRankedImmediateRepliesf          func(cardID, forUser globalid.ID) ([]*model.Card, error)
	GetFlatRepliesf                     func(id, forUser globalid.ID, latestFirst bool, limitTo int) ([]*model.Card, error)
	GetChannelf                         func(id globalid.ID) (*model.Channel, error)
	GetChannelsf                        func() ([]*model.Channel, error)
	GetSubscribedChannelsf              func(userID globalid.ID) ([]*model.Channel, error)
	AddUserToDefaultChannelsf           func(userID globalid.ID) error
	GetCardsForChannelf                 func(channelID globalid.ID, count, skip int, forUser globalid.ID) ([]*model.Card, error)
	JoinChannelf                        func(userID, channelID globalid.ID) error
	LeaveChannelf                       func(userID, channelID globalid.ID) error
	MuteChannelf                        func(userID, channelID globalid.ID) error
	UnmuteChannelf                      func(userID, channelID globalid.ID) error
	MuteUserf                           func(userID, mutedUserID globalid.ID) error
	UnmuteUserf                         func(userID, mutedUserID globalid.ID) error
	MuteThreadf                         func(userID, threadRootID globalid.ID) error
	UnmuteThreadf                       func(userID, threadRootID globalid.ID) error
	FollowDefaultUsersf                 func(userID globalid.ID) error
	AddCardsToTopOfFeedf                func(userID globalid.ID, cardRankedIDs []globalid.ID) error
	ResetUserFeedTopf                   func(userID globalid.ID) error
	GetFeedCardsFromCurrentTopf         func(userID globalid.ID, perPage, page int) ([]*model.Card, error)
	SetCardVisitedf                     func(userID, cardID globalid.ID) error
	NewContentAvailableForUserf         func(userID, cardID globalid.ID) (bool, error)
	SetFeedLastUpdatedForUserf          func(userID globalid.ID, t time.Time) error
	GetIntroCardIDsf                    func() ([]globalid.ID, error)
	SaveChannelf                        func(m *model.Channel) error
	SaveScoreModificationf              func(m *model.ScoreModification) error

	GetFeedCardsFromCurrentTopWithQueryf func(userID globalid.ID, perPage, page int, searchString string) ([]*model.Card, error)
	GetChannelByHandleg                  func(handle string) (*model.Channel, error)

	AddCardToPopularRanksf      func(card *model.Card, startingScore float64) error
	UpdatePopularRankScoref     func(cardID globalid.ID, scoreUpdate float64) error
	GetPopularRankCardsForUserf func(userID globalid.ID, page, perPage int) ([]*model.Card, error)
	UpdatePopularRanksForUserf  func(userID globalid.ID) error

	GetChannelsForUserf func(userID globalid.ID) ([]*model.Channel, error)

	AwardCoinsf          func(userID globalid.ID, amount int64) error
	PayCoinAmountf       func(userID globalid.ID, amount int64) error
	AwardTemporaryCoinsf func(userID globalid.ID, amount int64) error

	GetCurrentBalancef func(userID globalid.ID) (*model.CoinBalances, error)

	GetAnnouncementForNotificationf func(n *model.Notification) (*model.Announcement, error)
	LeaveAllChannelsf               func(userID globalid.ID) error
	GetIsSubscribedf                func(userID, channelID globalid.ID) (bool, error)
	GetChannelByHandlef             func(handle string) (*model.Channel, error)

	CountPostsByAliasInThreadf func(aliasID, threadRootID globalid.ID) (int, error)

	GetUserReactionf               func(userID, cardID globalid.ID) (*model.UserReaction, error)
	DeleteUserReactionForTypef     func(userID, cardID globalid.ID, typ model.UserReactionType) (int64, error)
	SaveUserReactionf              func(m *model.UserReaction) error
	DeleteReactionForNotificationf func(notifID, userID, cardID globalid.ID) error
	GetPopularRankForCardf         func(cardID globalid.ID) (*model.PopularRankEntry, error)
	SavePopularRankf               func(m *model.PopularRankEntry) error
	UpdatePopularRankForCardf      func(cardID globalid.ID, viewCountChange, upCountChange, downCountChange, commentCountChange int64, scoreModChange float64) error
	UpdateViewsForCardsf           func(cardIDs []globalid.ID) error

	GetCoinsReceivedNotificationDataf func(notif *model.Notification) (*model.CoinReward, error)
	GetLeaderboardRankingsf           func(count, skip int) ([]*model.LeaderboardRank, error)
	CountCardsInFeedf                 func(userID globalid.ID) (int64, error)
	BuildInitialFeedf                 func(userID globalid.ID) error
	GetRankableCardsForUserf          func(userID globalid.ID) ([]*model.PopularRankEntry, error)
	UpdateCardRanksForUserf           func(userID globalid.ID, cardIDs []globalid.ID) error

	GetChannelInfosf    func(userID globalid.ID) ([]*model.ChannelUserInfo, error)
	GetSubscriberCountf func(channelID globalid.ID) (int, error)

	GetLeaderboardNotificationExportDataf func(notifID globalid.ID) (*datastore.LeaderboardNotificationExportData, error)
	UpdateUniqueCommentersForCardf        func(cardID globalid.ID) error

	UpdateAllNotificationsSeenf           func(userID globalid.ID) error
	GetPopularPostNotificationExportDataf func(n *model.Notification) (*datastore.PopularPostNotificationExportData, error)

	GetCardsForPopularRankSincef       func(t time.Time) ([]*model.PopularRankEntry, error)
	UpdatePopularRanksWithListf        func(userID globalid.ID, ids []globalid.ID) error
	SaveUserTipf                       func(m *model.UserTip) error
	AssignAliasForUserTipsInThreadf    func(userID, threadRootID, aliasID globalid.ID) error
	FeaturedCommentsForUserByCardIDsf  func(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.Card, error)
	GetAnonymousAliasesByIDf           func(ids []globalid.ID) ([]*model.AnonymousAlias, error)
	GetChannelsByIDf                   func(ids []globalid.ID) ([]*model.Channel, error)
	GetEngagementsf                    func(cardIDs []globalid.ID) (map[globalid.ID]*model.Engagement, error)
	GetThreadCountsf                   func(ids []globalid.ID) (map[globalid.ID]int, error)
	GetUserReactionsf                  func(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.UserReaction, error)
	IsFollowingsf                      func(followerID globalid.ID, followeeID []globalid.ID) (map[globalid.ID]bool, error)
	IsSubscribedToChannelsf            func(userID globalid.ID, channelIDs []globalid.ID) (map[globalid.ID]bool, error)
	NewContentAvailableForUserByCardsf func(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error)
	SubscribedToCardsf                 func(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error)
}

func newStore() *testStore {
	return &testStore{}
}
func (s *testStore) Close() error {
	return nil
}
func (s *testStore) SaveUser(u *model.User) (err error) {
	if s.SaveUserf == nil {
		return nil
	}
	return s.SaveUserf(u)
}

func (s *testStore) GetPopularPostNotificationExportData(n *model.Notification) (*datastore.PopularPostNotificationExportData, error) {
	if s.GetPopularPostNotificationExportDataf == nil {
		return nil, nil
	}
	return s.GetPopularPostNotificationExportDataf(n)
}

func (s *testStore) AssignAliasForUserTipsInThread(userID, threadRootID, aliasID globalid.ID) error {
	if s.AssignAliasForUserTipsInThreadf == nil {
		return nil
	}
	return s.AssignAliasForUserTipsInThreadf(userID, threadRootID, aliasID)
}

func (s *testStore) SaveUserTip(m *model.UserTip) error {
	if s.SaveUserTipf == nil {
		return nil
	}
	return s.SaveUserTipf(m)
}

func (s *testStore) GetCardsForPopularRankSince(t time.Time) ([]*model.PopularRankEntry, error) {
	if s.GetCardsForPopularRankSincef == nil {
		return nil, nil
	}
	return s.GetCardsForPopularRankSince(t)
}

func (s *testStore) UpdatePopularRanksWithList(userID globalid.ID, ids []globalid.ID) error {
	if s.UpdatePopularRanksWithListf == nil {
		return nil
	}
	return s.UpdatePopularRanksWithListf(userID, ids)
}

func (s *testStore) GetLeaderboardNotificationExportData(notifID globalid.ID) (*datastore.LeaderboardNotificationExportData, error) {
	if s.GetLeaderboardNotificationExportDataf == nil {
		return nil, nil
	}
	return s.GetLeaderboardNotificationExportDataf(notifID)
}

func (s *testStore) UpdateAllNotificationsSeen(userID globalid.ID) error {
	if s.UpdateAllNotificationsSeenf == nil {
		return nil
	}
	return s.UpdateAllNotificationsSeenf(userID)
}

func (s *testStore) DeleteSession(id globalid.ID) error {
	if s.DeleteSessionf == nil {
		return nil
	}
	return s.DeleteSessionf(id)
}

func (s *testStore) GetCard(id globalid.ID) (*model.Card, error) {
	if s.GetCardf == nil {
		return nil, nil
	}
	return s.GetCardf(id)
}

func (s *testStore) GetThread(id, forUser globalid.ID) ([]*model.Card, error) {
	return s.GetThreadf(id, forUser)
}

func (s *testStore) GetUser(nodeID globalid.ID) (u *model.User, err error) {
	if s.GetUserf == nil {
		return nil, nil
	}
	return s.GetUserf(nodeID)
}
func (s *testStore) GetUserByEmail(email string) (*model.User, error) { return s.GetUserByEmailf(email) }
func (s *testStore) GetUserByUsername(username string) (*model.User, error) {
	return s.GetUserByUsernamef(username)
}
func (s *testStore) GetUsersByUsernames(usernames []string) ([]*model.User, error) {
	return s.GetUsersByUsernamesf(usernames)
}
func (s *testStore) GetUsers() (res []*model.User, err error) { return s.GetUsersf() }

func (s *testStore) GetChannelInfos(userID globalid.ID) ([]*model.ChannelUserInfo, error) {
	if s.GetChannelInfosf == nil {
		return nil, nil
	}
	return s.GetChannelInfosf(userID)
}

func (s *testStore) UpdateUniqueCommentersForCard(cardID globalid.ID) error {
	if s.UpdateUniqueCommentersForCardf == nil {
		return nil
	}
	return s.UpdateUniqueCommentersForCardf(cardID)
}

func (s *testStore) GetSubscriberCount(channelID globalid.ID) (int, error) {
	if s.GetSubscriberCountf == nil {
		return 0, nil
	}
	return s.GetSubscriberCountf(channelID)
}

func (s *testStore) SaveInvite(invite *model.Invite) error {
	return s.SaveInvitef(invite)
}
func (s *testStore) GetInviteByToken(token string) (*model.Invite, error) {
	return s.GetInviteByTokenf(token)
}
func (s *testStore) DeleteInvite(id globalid.ID) error {
	return s.DeleteInvitef(id)
}

func (s *testStore) SaveResetToken(rt *model.ResetToken) error {
	return s.SaveResetTokenf(rt)
}
func (s *testStore) GetResetToken(userID globalid.ID) (*model.ResetToken, error) {
	return s.GetResetTokenf(userID)
}
func (s *testStore) GetNotifications(userID globalid.ID, pageSize, pageNumber int) ([]*model.Notification, error) {
	return s.GetNotifcationsf(userID, pageSize, pageNumber)
}
func (s *testStore) UpdateNotificationsSeen(ids []globalid.ID) error {
	return s.UpdateNotificationsSeenf(ids)
}

func (s *testStore) UpdateNotificationsOpened(ids []globalid.ID) error {
	if s.UpdateNotificationsOpenedf == nil {
		return nil
	}
	return s.UpdateNotificationsOpenedf(ids)
}

func (s *testStore) GetAnonymousAlias(id globalid.ID) (*model.AnonymousAlias, error) {
	return s.GetAnonymousAliasf(id)
}

func (s *testStore) GetAnonymousAliases() ([]*model.AnonymousAlias, error) {
	return s.GetAnonymousAliasesf()
}

func (s *testStore) GetUnusedAlias(cardID globalid.ID) (*model.AnonymousAlias, error) {
	return s.GetUnusedAliasf(cardID)
}

func (s *testStore) GetPostedCardsForNode(nodeID globalid.ID, skip, count int) ([]*model.Card, error) {
	return s.GetPostedCardsForNodef(nodeID, skip, count)
}

func (s *testStore) SaveSession(session *model.Session) error { return s.SaveSessionf(session) }

func (s *testStore) SaveCard(c *model.Card) error {
	return s.SaveCardf(c)
}
func (s *testStore) DeleteCard(id globalid.ID) error {
	return s.DeleteCardf(id)
}

func (s *testStore) SaveOAuthAccount(oaa *model.OAuthAccount) error {
	return s.SaveOAuthAccountf(oaa)
}

func (s *testStore) GetOAuthAccountBySubject(subject string) (*model.OAuthAccount, error) {
	return s.GetOAuthAccountBySubjectf(subject)
}

func (s *testStore) GetOnSwitches() ([]*model.FeatureSwitch, error) {
	return s.GetOnSwitchesf()
}

func (s testStore) LatestForType(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error) {
	if s.LatestForTypef == nil {
		return &model.Notification{ID: globalid.Nil}, nil
	}
	return s.LatestForTypef(userID, targetID, typ, unopenedOnly)
}
func (s testStore) SaveNotification(m *model.Notification) error {
	return s.SaveNotificationf(m)
}

func (s testStore) SaveScoreModification(m *model.ScoreModification) error {
	if s.SaveScoreModificationf == nil {
		return nil
	}
	return s.SaveScoreModificationf(m)
}

func (s testStore) GetAnnouncementForNotification(n *model.Notification) (*model.Announcement, error) {
	if s.GetAnnouncementForNotificationf == nil {
		return nil, nil
	}
	return s.GetAnnouncementForNotificationf(n)
}

func (s testStore) GetAnonymousAliasByUsername(username string) (*model.AnonymousAlias, error) {
	return s.GetAnonymousAliasByUsernamef(username)
}

func (s testStore) DeleteNotificationsForCard(cardID globalid.ID) error {
	return s.DeleteNotificationsForCardf(cardID)
}

func (s testStore) UnseenNotificationsCount(userID globalid.ID) (int, error) {
	return s.UnseenNotificationsCountf(userID)
}

func (s testStore) GetPostedCardsForNodeIncludingAnon(nodeID globalid.ID, count, skip int) ([]*model.Card, error) {
	return s.GetPostedCardsForNodeIncludingAnonf(nodeID, count, skip)
}

func (s testStore) GetAnonymousAliasLastUsed(userID, threadRootID globalid.ID) (bool, error) {
	return s.LastPostInThreadWasAnonymousf(userID, threadRootID)
}

func (s testStore) GetThreadCount(id globalid.ID) (int, error) {
	return s.GetThreadCountf(id)
}

func (s testStore) SaveWaitlistEntry(m *model.WaitlistEntry) error {
	return s.SaveWaitlistEntryf(m)
}

func (s testStore) GetInvitesForUser(id globalid.ID) ([]*model.Invite, error) {
	if s.GetInvitesForUserf == nil {
		return nil, nil
	}
	return s.GetInvitesForUserf(id)
}

func (s testStore) GetInvite(id globalid.ID) (*model.Invite, error) {
	return s.GetInvitef(id)
}

func (s testStore) GetEngagement(cardID globalid.ID) (*model.Engagement, error) {
	return s.GetEngagementf(cardID)
}

func (s *testStore) DeleteWaitlistEntry(email string) error {
	return s.DeleteWaitlistEntryf(email)
}

func (s *testStore) GetSafeUsersByPage(forUser globalid.ID, pageSize, pageNumber int, searchString string) ([]*model.User, error) {
	if s.GetSafeUsersByPagef == nil {
		return nil, nil
	}
	return s.GetSafeUsersByPagef(forUser, pageSize, pageNumber, searchString)
}

func (s *testStore) SaveMention(m *model.Mention) error {
	return s.SaveMentionf(m)
}
func (s *testStore) GetMentionExportData(n *model.Notification) (*datastore.MentionNotificationExportData, error) {
	return s.GetMentionExportDataf(n)
}
func (s *testStore) SaveNotificationMention(m *model.NotificationMention) error {
	return s.SaveNotificationMentionf(m)
}

func (s *testStore) GetCommentExportData(n *model.Notification) (*datastore.CommentNotificationExportData, error) {
	return s.GetCommentExportDataf(n)
}
func (s *testStore) GetLikeNotificationExportData(n *model.Notification) (*datastore.LikeNotificationExportData, error) {
	return s.GetLikeNotificationExportDataf(n)
}
func (s *testStore) SaveNotificationComment(m *model.NotificationComment) error {
	return s.SaveNotificationCommentf(m)
}
func (s *testStore) SaveReactionForNotification(notifID, userID, cardID globalid.ID) error {
	if s.SaveReactionForNotificationf == nil {
		return nil
	}
	return s.SaveReactionForNotificationf(notifID, userID, cardID)
}

func (s *testStore) SubscribeToCard(userID, cardID globalid.ID, typ string) error {
	return s.SubscribeToCardf(userID, cardID, typ)
}
func (s *testStore) UnsubscribeFromCard(userID, cardID globalid.ID, typ string) error {
	if s.UnsubscribeFromCardf == nil {
		return nil
	}
	return s.UnsubscribeFromCardf(userID, cardID, typ)
}
func (s *testStore) SubscribersForCard(cardID globalid.ID, typ string) ([]globalid.ID, error) {
	return s.SubscribersForCardf(cardID, typ)
}

func (s *testStore) SubscribedToTypes(userID, cardID globalid.ID) ([]string, error) {
	return s.SubscribedToTypesf(userID, cardID)
}
func (s *testStore) DeleteMentionsForCard(cardID globalid.ID) error {
	return s.DeleteMentionsForCardf(cardID)
}

func (s *testStore) ReassignInviterForGroup(rootInviteID, newUserID globalid.ID) error {
	return s.ReassignInviterForGroupf(rootInviteID, newUserID)
}

func (s *testStore) GroupInvitesByToken(tokens []string, groupID globalid.ID) error {
	return s.GroupInvitesByTokenf(tokens, groupID)
}

func (s *testStore) ReassignInvitesByToken(tokens []string, userID globalid.ID) error {
	return s.ReassignInvitesByTokenf(tokens, userID)
}

func (s *testStore) SaveFollower(followerID, followeeID globalid.ID) error {
	if s.SaveFollowerf == nil {
		return nil
	}
	return s.SaveFollowerf(followerID, followeeID)
}

func (s *testStore) DeleteFollower(followerID, followeeID globalid.ID) error {
	return s.DeleteFollowerf(followerID, followeeID)
}

func (s *testStore) GetFollowing(userID globalid.ID) ([]*model.User, error) {
	return s.GetFollowingf(userID)
}

func (s *testStore) ClearEmptyNotifications() error {
	return s.ClearEmptyNotificationsf()
}

func (s *testStore) ClearEmptyNotificationsForUser(userID globalid.ID) error {
	return s.ClearEmptyNotificationsForUserf(userID)
}

func (s *testStore) SaveNotificationFollow(notifID, followerID, followeeID globalid.ID) error {
	return s.SaveNotificationFollowf(notifID, followerID, followeeID)
}
func (s *testStore) GetFollowExportData(n *model.Notification) (*datastore.FollowNotificationExportData, error) {
	return s.GetFollowExportDataf(n)
}
func (s *testStore) IsFollowing(followerID, followeeID globalid.ID) (bool, error) {
	if s.IsFollowingf == nil {
		return false, nil
	}
	return s.IsFollowingf(followerID, followeeID)
}

func (s *testStore) GetNotification(notifID globalid.ID) (*model.Notification, error) {
	return s.GetNotificationf(notifID)
}

func (s *testStore) BlockUser(blockingUser, blockedUser globalid.ID) error {
	if s.BlockUserf == nil {
		return nil
	}
	return s.BlockUserf(blockingUser, blockedUser)
}

func (s *testStore) GetCoinsReceivedNotificationData(notif *model.Notification) (*model.CoinReward, error) {
	if s.GetCoinsReceivedNotificationDataf == nil {
		return nil, nil
	}
	return s.GetCoinsReceivedNotificationDataf(notif)
}

func (s *testStore) BlockAnonUserInThread(blockingUser, blockedAlias, threadID globalid.ID) error {
	if s.BlockAnonUserInThreadf == nil {
		return nil
	}
	return s.BlockAnonUserInThreadf(blockingUser, blockedAlias, threadID)
}

func (s *testStore) FeaturedCommentsForUser(userID, cardID globalid.ID) ([]*model.Card, error) {
	if s.FeaturedCommentsForUserf == nil {
		return nil, nil
	}
	return s.FeaturedCommentsForUserf(userID, cardID)
}

func (s *testStore) GetCurrentBalance(userID globalid.ID) (*model.CoinBalances, error) {
	if s.GetCurrentBalancef == nil {
		return nil, nil
	}
	return s.GetCurrentBalancef(userID)
}

func (s *testStore) GetRankedImmediateReplies(cardID, forUser globalid.ID) ([]*model.Card, error) {
	if s.GetRankedImmediateRepliesf == nil {
		return nil, nil
	}
	return s.GetRankedImmediateRepliesf(cardID, forUser)
}

func (s *testStore) GetFlatReplies(id, forUser globalid.ID, latestFirst bool, limitTo int) ([]*model.Card, error) {
	if s.GetFlatRepliesf == nil {
		return nil, nil
	}
	return s.GetFlatRepliesf(id, forUser, latestFirst, limitTo)
}

func (s *testStore) GetChannel(id globalid.ID) (*model.Channel, error) {
	if s.GetChannelf == nil {
		return nil, nil
	}
	return s.GetChannelf(id)
}
func (s *testStore) GetSubscribedChannels(userID globalid.ID) ([]*model.Channel, error) {
	if s.GetSubscribedChannelsf == nil {
		return nil, nil
	}
	return s.GetSubscribedChannelsf(userID)
}
func (s *testStore) AddUserToDefaultChannels(userID globalid.ID) error {
	if s.AddUserToDefaultChannelsf == nil {
		return nil
	}
	return s.AddUserToDefaultChannelsf(userID)
}

func (s *testStore) GetCardsForChannel(channelID globalid.ID, count, skip int, forUser globalid.ID) ([]*model.Card, error) {
	if s.GetCardsForChannelf == nil {
		return nil, nil
	}
	return s.GetCardsForChannelf(channelID, count, skip, forUser)
}

func (s *testStore) GetChannelsForUser(userID globalid.ID) ([]*model.Channel, error) {
	if s.GetChannelsForUserf == nil {
		return nil, nil
	}
	return s.GetChannelsForUserf(userID)
}

func (s *testStore) JoinChannel(userID, channelID globalid.ID) error {
	if s.JoinChannelf == nil {
		return nil
	}
	return s.JoinChannelf(userID, channelID)
}

func (s *testStore) CountCardsInFeed(userID globalid.ID) (int64, error) {
	if s.CountCardsInFeedf == nil {
		return 0, nil
	}
	return s.CountCardsInFeedf(userID)
}

func (s *testStore) BuildInitialFeed(userID globalid.ID) error {
	if s.BuildInitialFeedf == nil {
		return nil
	}
	return s.BuildInitialFeedf(userID)
}

func (s *testStore) LeaveChannel(userID, channelID globalid.ID) error {
	if s.LeaveChannelf == nil {
		return nil
	}
	return s.LeaveChannelf(userID, channelID)
}

func (s *testStore) MuteChannel(userID, channelID globalid.ID) error {
	if s.MuteChannelf == nil {
		return nil
	}
	return s.MuteChannelf(userID, channelID)
}

func (s *testStore) UnmuteChannel(userID, channelID globalid.ID) error {
	if s.UnmuteChannelf == nil {
		return nil
	}
	return s.UnmuteChannelf(userID, channelID)
}

func (s *testStore) MuteUser(userID, mutedUserID globalid.ID) error {
	if s.MuteUserf == nil {
		return nil
	}
	return s.MuteUserf(userID, mutedUserID)
}

func (s *testStore) UnmuteUser(userID, mutedUserID globalid.ID) error {
	if s.UnmuteUserf == nil {
		return nil
	}
	return s.UnmuteUserf(userID, mutedUserID)
}

func (s *testStore) MuteThread(userID, threadRootID globalid.ID) error {
	if s.MuteThreadf == nil {
		return nil
	}
	return s.MuteThreadf(userID, threadRootID)
}

func (s *testStore) UnmuteThread(userID, threadRootID globalid.ID) error {
	if s.UnmuteThreadf == nil {
		return nil
	}
	return s.UnmuteThreadf(userID, threadRootID)
}

func (s *testStore) AddCardsToTopOfFeed(userID globalid.ID, cardRankedIDs []globalid.ID) error {
	if s.AddCardsToTopOfFeedf == nil {
		return nil
	}
	return s.AddCardsToTopOfFeedf(userID, cardRankedIDs)
}

func (s *testStore) ResetUserFeedTop(userID globalid.ID) error {
	if s.ResetUserFeedTopf == nil {
		return nil
	}
	return s.ResetUserFeedTopf(userID)
}

func (s *testStore) GetFeedCardsFromCurrentTop(userID globalid.ID, perPage, page int) ([]*model.Card, error) {
	if s.GetFeedCardsFromCurrentTopf == nil {
		return nil, nil
	}
	return s.GetFeedCardsFromCurrentTopf(userID, perPage, page)
}

func (s *testStore) SetCardVisited(userID, cardID globalid.ID) error {
	if s.SetCardVisitedf == nil {
		return nil
	}
	return s.SetCardVisitedf(userID, cardID)
}

func (s *testStore) NewContentAvailableForUser(userID, cardID globalid.ID) (bool, error) {
	if s.NewContentAvailableForUserf == nil {
		return false, nil
	}
	return s.NewContentAvailableForUserf(userID, cardID)
}

func (s *testStore) LeaveAllChannels(userID globalid.ID) error {
	if s.LeaveAllChannelsf == nil {
		return nil
	}
	return s.LeaveAllChannelsf(userID)
}

func (s *testStore) FollowDefaultUsers(userID globalid.ID) error {
	if s.FollowDefaultUsersf == nil {
		return nil
	}
	return s.FollowDefaultUsersf(userID)
}
func (s *testStore) SetFeedLastUpdatedForUser(userID globalid.ID, t time.Time) error {
	if s.SetFeedLastUpdatedForUserf == nil {
		return nil
	}
	return s.SetFeedLastUpdatedForUserf(userID, t)
}

func (s *testStore) GetIntroCardIDs() ([]globalid.ID, error) {
	if s.GetIntroCardIDsf == nil {
		return nil, nil
	}
	return s.GetIntroCardIDsf()
}

func (s *testStore) SaveChannel(m *model.Channel) error {
	if s.SaveChannelf == nil {
		return nil
	}
	return s.SaveChannelf(m)
}

func (s *testStore) AddCardToPopularRanks(card *model.Card, startingScore float64) error {
	if s.AddCardToPopularRanksf == nil {
		return nil
	}
	return s.AddCardToPopularRanksf(card, startingScore)
}

func (s *testStore) UpdatePopularRankScore(cardID globalid.ID, scoreUpdate float64) error {
	if s.UpdatePopularRankScoref == nil {
		return nil
	}
	return s.UpdatePopularRankScoref(cardID, scoreUpdate)
}

func (s *testStore) GetPopularRankCardsForUser(userID globalid.ID, page, perPage int) ([]*model.Card, error) {
	if s.GetPopularRankCardsForUserf == nil {
		return nil, nil
	}
	return s.GetPopularRankCardsForUserf(userID, page, perPage)
}

func (s *testStore) UpdatePopularRanksForUser(userID globalid.ID) error {
	if s.UpdatePopularRanksForUserf == nil {
		return nil
	}
	return s.UpdatePopularRanksForUserf(userID)
}

func (s *testStore) GetFeedCardsFromCurrentTopWithQuery(userID globalid.ID, perPage, page int, searchString string) ([]*model.Card, error) {
	if s.GetFeedCardsFromCurrentTopWithQueryf == nil {
		return nil, nil
	}
	return s.GetFeedCardsFromCurrentTopWithQueryf(userID, perPage, page, searchString)
}

func (s *testStore) AwardCoins(userID globalid.ID, amount int64) error {
	if s.AwardCoinsf == nil {
		return nil
	}
	return s.AwardCoinsf(userID, amount)
}

func (s *testStore) PayCoinAmount(userID globalid.ID, amount int64) error {
	if s.PayCoinAmountf == nil {
		return nil
	}
	return s.PayCoinAmountf(userID, amount)
}

func (s *testStore) AwardTemporaryCoins(userID globalid.ID, amount int64) error {
	if s.AwardTemporaryCoinsf == nil {
		return nil
	}
	return s.AwardTemporaryCoinsf(userID, amount)
}

func (s *testStore) GetChannelByHandle(handle string) (*model.Channel, error) {
	if s.GetChannelByHandlef == nil {
		return nil, nil
	}
	return s.GetChannelByHandlef(handle)
}

func (s *testStore) GetIsSubscribed(userID, channelID globalid.ID) (bool, error) {
	if s.GetIsSubscribedf == nil {
		return false, nil
	}
	return s.GetIsSubscribedf(userID, channelID)
}

func (s *testStore) CountPostsByAliasInThread(aliasID, threadRootID globalid.ID) (int, error) {
	if s.CountPostsByAliasInThreadf == nil {
		return 0, nil
	}
	return s.CountPostsByAliasInThreadf(aliasID, threadRootID)
}

func (s *testStore) GetUserReaction(userID, cardID globalid.ID) (*model.UserReaction, error) {
	if s.GetUserReactionf == nil {
		return nil, nil
	}
	return s.GetUserReactionf(userID, cardID)
}
func (s *testStore) DeleteUserReactionForType(userID, cardID globalid.ID, typ model.UserReactionType) (int64, error) {
	if s.DeleteUserReactionForTypef == nil {
		return 0, nil
	}
	return s.DeleteUserReactionForTypef(userID, cardID, typ)
}

func (s *testStore) SaveUserReaction(m *model.UserReaction) error {
	if s.SaveUserReactionf == nil {
		return nil
	}
	return s.SaveUserReactionf(m)
}

func (s *testStore) DeleteReactionForNotification(notifID, userID, cardID globalid.ID) error {
	if s.DeleteReactionForNotificationf == nil {
		return nil
	}
	return s.DeleteReactionForNotificationf(notifID, userID, cardID)
}

func (s *testStore) GetPopularRankForCard(cardID globalid.ID) (*model.PopularRankEntry, error) {
	if s.GetPopularRankForCardf == nil {
		return nil, nil
	}
	return s.GetPopularRankForCardf(cardID)
}

func (s *testStore) SavePopularRank(m *model.PopularRankEntry) error {
	if s.SavePopularRankf == nil {
		return nil
	}
	return s.SavePopularRankf(m)
}
func (s *testStore) UpdatePopularRankForCard(cardID globalid.ID, viewCountChange, upCountChange, downCountChange, commentCountChange int64, scoreModChange float64) error {
	if s.UpdatePopularRankForCardf == nil {
		return nil
	}
	return s.UpdatePopularRankForCardf(cardID, viewCountChange, upCountChange, downCountChange, commentCountChange, scoreModChange)
}

func (s *testStore) GetLeaderboardRankings(count, skip int) ([]*model.LeaderboardRank, error) {
	if s.GetLeaderboardRankingsf == nil {
		return nil, nil
	}
	return s.GetLeaderboardRankingsf(count, skip)
}

func (s *testStore) UpdateViewsForCards(cardIDs []globalid.ID) error {
	if s.UpdateViewsForCardsf == nil {
		return nil
	}
	return s.UpdateViewsForCardsf(cardIDs)
}

func (s *testStore) GetRankableCardsForUser(userID globalid.ID) ([]*model.PopularRankEntry, error) {
	if s.GetRankableCardsForUserf == nil {
		return nil, nil
	}
	return s.GetRankableCardsForUserf(userID)
}

func (s *testStore) UpdateCardRanksForUser(userID globalid.ID, cardIDs []globalid.ID) error {
	if s.UpdateCardRanksForUserf == nil {
		return nil
	}
	return s.UpdateCardRanksForUserf(userID, cardIDs)
}

func (s *testStore) FeaturedCommentsForUserByCardIDs(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.Card, error) {
	if s.FeaturedCommentsForUserByCardIDsf == nil {
		return nil, nil
	}
	return s.FeaturedCommentsForUserByCardIDsf(userID, cardIDs)
}
func (s *testStore) GetAnonymousAliasesByID(ids []globalid.ID) ([]*model.AnonymousAlias, error) {
	if s.GetAnonymousAliasesByIDf == nil {
		return nil, nil
	}
	return s.GetAnonymousAliasesByIDf(ids)
}

func (s *testStore) GetChannelsByID(ids []globalid.ID) ([]*model.Channel, error) {
	if s.GetChannelsByIDf == nil {
		return nil, nil
	}
	return s.GetChannelsByIDf(ids)
}

func (s *testStore) GetEngagements(cardIDs []globalid.ID) (map[globalid.ID]*model.Engagement, error) {
	if s.GetEngagementsf == nil {
		return nil, nil
	}
	return s.GetEngagementsf(cardIDs)
}

func (s *testStore) GetThreadCounts(ids []globalid.ID) (map[globalid.ID]int, error) {
	if s.GetThreadCountsf == nil {
		return nil, nil
	}
	return s.GetThreadCountsf(ids)
}

func (s *testStore) GetUserReactions(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.UserReaction, error) {
	if s.GetUserReactionsf == nil {
		return nil, nil
	}
	return s.GetUserReactionsf(userID, cardIDs)
}

func (s *testStore) GetUsersByID(ids []globalid.ID) (res []*model.User, err error) {
	if s.GetUsersByIDf == nil {
		return nil, nil
	}
	return s.GetUsersByID(ids)
}

func (s *testStore) IsFollowings(followerID globalid.ID, followeeIDs []globalid.ID) (map[globalid.ID]bool, error) {
	if s.IsFollowingsf == nil {
		return nil, nil
	}
	return s.IsFollowingsf(followerID, followeeIDs)
}

func (s *testStore) IsSubscribedToChannels(userID globalid.ID, channelIDs []globalid.ID) (map[globalid.ID]bool, error) {
	if s.IsSubscribedToChannelsf == nil {
		return nil, nil
	}
	return s.IsSubscribedToChannelsf(userID, channelIDs)
}

func (s *testStore) NewContentAvailableForUserByCards(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error) {
	if s.NewContentAvailableForUserByCardsf == nil {
		return nil, nil
	}
	return s.NewContentAvailableForUserByCardsf(userID, cardIDs)
}

func (s *testStore) SubscribedToCards(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error) {
	if s.SubscribedToCardsf == nil {
		return nil, nil
	}
	return s.SubscribedToCardsf(userID, cardIDs)
}

type testWorker struct {
	// private queues for testing
	EnqueueMailJobf func(*emailsender.Job) error
}

func newWorker() *testWorker {
	return &testWorker{}
}

func (tw *testWorker) EnqueueMailJob(job *emailsender.Job) error {
	return tw.EnqueueMailJobf(job)
}

type mockImageProcessor struct {
	SaveBase64ProfileImagef      func(data string) (string, string, error)
	SaveBase64CoverImagef        func(data string) (string, string, error)
	SaveBase64CardImagef         func(data string) (string, string, error)
	SaveBase64CardContentImagef  func(data string) (string, string, error)
	GenerateDefaultProfileImagef func() (string, string, error)
	DownloadCardImagef           func(url string) (string, string, error)
	DownloadProfileImagef        func(url string) (string, string, error)
	BlendImagef                  func(src, gradient string) (string, error)
	GradientImagef               func(gradient string) (string, string, error)
}

func (mip *mockImageProcessor) SaveBase64ProfileImage(data string) (string, string, error) {
	return mip.SaveBase64ProfileImagef(data)
}

func (mip *mockImageProcessor) SaveBase64CoverImage(data string) (string, string, error) {
	return mip.SaveBase64CoverImagef(data)
}

func (mip *mockImageProcessor) SaveBase64CardImage(data string) (string, string, error) {
	return mip.SaveBase64CardImagef(data)
}

func (mip *mockImageProcessor) SaveBase64CardContentImage(data string) (string, string, error) {
	return mip.SaveBase64CardContentImagef(data)
}

func (mip *mockImageProcessor) GenerateDefaultProfileImage() (string, string, error) {
	return mip.GenerateDefaultProfileImagef()
}

func (mip *mockImageProcessor) DownloadCardImage(url string) (string, string, error) {
	return mip.DownloadCardImagef(url)
}

func (mip *mockImageProcessor) DownloadProfileImage(url string) (string, string, error) {
	return mip.DownloadProfileImagef(url)
}

func (mip *mockImageProcessor) BlendImage(src, gradient string) (string, error) {
	return mip.BlendImagef(src, gradient)
}

func (mip *mockImageProcessor) GradientImage(gradient string) (string, string, error) {
	return mip.GradientImagef(gradient)
}

type pusherMock struct {
	NewCardf            func(ctx context.Context, session *model.Session, card *model.CardResponse) error
	DeleteCardf         func(ctx context.Context, session *model.Session, id globalid.ID) error
	UpdateCardf         func(ctx context.Context, session *model.Session, card *model.CardResponse) error
	UpdateUserf         func(ctx context.Context, session *model.Session, user *model.ExportedUser) error
	NewNotificationf    func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error
	UpdateNotificationf func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error
	UpdateCoinBalancef  func(ctx context.Context, userID globalid.ID, newBalances *model.CoinBalances) error

	UpdateEngagementf func(ctx context.Context, session *model.Session, cardID globalid.ID) error
}

func (p *pusherMock) NewCard(ctx context.Context, session *model.Session, card *model.CardResponse) error {
	return p.NewCardf(ctx, session, card)
}

func (p *pusherMock) DeleteCard(ctx context.Context, session *model.Session, id globalid.ID) error {
	return p.DeleteCardf(ctx, session, id)
}

func (p *pusherMock) UpdateCoinBalance(ctx context.Context, userID globalid.ID, newBalances *model.CoinBalances) error {
	if p.UpdateCoinBalancef == nil {
		return nil
	}
	return p.UpdateCoinBalancef(ctx, userID, newBalances)
}

func (p *pusherMock) UpdateCard(ctx context.Context, session *model.Session, card *model.CardResponse) error {
	return p.UpdateCardf(ctx, session, card)
}

func (p *pusherMock) UpdateUser(ctx context.Context, session *model.Session, user *model.ExportedUser) error {
	return p.UpdateUserf(ctx, session, user)
}

func (p *pusherMock) NewNotification(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
	if p.NewNotificationf == nil {
		return nil
	}
	return p.NewNotificationf(ctx, session, notif)
}

func (p *pusherMock) UpdateNotification(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
	return p.UpdateNotificationf(ctx, session, notif)
}

func (p *pusherMock) UpdateEngagement(ctx context.Context, session *model.Session, cardID globalid.ID) error {
	return p.UpdateEngagementf(ctx, session, cardID)
}

type mockIndexer struct {
	IndexUserf          func(m *model.User) error
	RemoveIndexForUserf func(m *model.User) error
	IndexChannelf       func(m *model.Channel) error
}

func (mi *mockIndexer) IndexUser(m *model.User) error {
	if mi.IndexUserf == nil {
		return nil
	}
	return mi.IndexUserf(m)
}

func (mi *mockIndexer) RemoveIndexForUser(m *model.User) error {
	if mi.RemoveIndexForUserf == nil {
		return nil
	}
	return mi.RemoveIndexForUserf(m)
}

func (mi *mockIndexer) IndexChannel(m *model.Channel) error {
	if mi.IndexChannelf == nil {
		return nil
	}
	return mi.IndexChannelf(m)
}

type mockOAuth2 struct {
	ExtendTokenf func(ctx context.Context, token string) (rpc.AccessToken, error)
}

func (moa2 *mockOAuth2) ExtendToken(ctx context.Context, token string) (rpc.AccessToken, error) {
	return moa2.ExtendTokenf(ctx, token)
}

type mockAccessToken struct {
	FacebookUserf func() (*rpc.FacebookUser, error)
	Tokenf        func() string
	ExpiresAtf    func() int64
}

func (mat *mockAccessToken) FacebookUser() (*rpc.FacebookUser, error) {
	return mat.FacebookUserf()
}

func (mat *mockAccessToken) Token() string {
	return mat.Tokenf()
}

func (mat *mockAccessToken) ExpiresAt() int64 {
	return mat.ExpiresAtf()
}

type mockNotifications struct {
	ExportNotificationf func(n *model.Notification) (*model.ExportedNotification, error)
}

func (ns *mockNotifications) ExportNotification(n *model.Notification) (*model.ExportedNotification, error) {
	if ns.ExportNotificationf == nil {
		return nil, nil
	}
	return ns.ExportNotificationf(n)
}

type mockResponses struct {
	FeedCardResponsesf func(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error)
}

func (mr *mockResponses) FeedCardResponses(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error) {
	if mr.FeedCardResponsesf == nil {
		return nil, nil
	}
	return mr.FeedCardResponsesf(cards, viewerID)
}
