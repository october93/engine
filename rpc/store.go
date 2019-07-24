package rpc

import (
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store/datastore"
)

type dataStore interface {
	DeleteSession(id globalid.ID) error
	GetThread(id, forUser globalid.ID) ([]*model.Card, error)
	GetThreadCount(id globalid.ID) (int, error)
	GetThreadCounts(ids []globalid.ID) (map[globalid.ID]int, error)
	GetCard(id globalid.ID) (*model.Card, error)
	GetUser(nodeID globalid.ID) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetAnonymousAlias(id globalid.ID) (*model.AnonymousAlias, error)
	GetAnonymousAliases() ([]*model.AnonymousAlias, error)
	GetAnonymousAliasesByID(ids []globalid.ID) ([]*model.AnonymousAlias, error)
	GetUnusedAlias(cardID globalid.ID) (*model.AnonymousAlias, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUsers() (res []*model.User, err error)
	GetUsersByUsernames(usernames []string) ([]*model.User, error)
	GetUsersByID(ids []globalid.ID) ([]*model.User, error)
	DeleteInvite(id globalid.ID) error
	SaveSession(session *model.Session) error
	SaveUser(u *model.User) (err error)
	SaveInvite(invite *model.Invite) error
	GetInviteByToken(token string) (*model.Invite, error)
	GetPostedCardsForNode(nodeID globalid.ID, skip, count int) ([]*model.Card, error)
	GetPostedCardsForNodeIncludingAnon(nodeID globalid.ID, skip, count int) ([]*model.Card, error)
	SaveResetToken(rt *model.ResetToken) error
	GetResetToken(userID globalid.ID) (*model.ResetToken, error)
	SaveCard(c *model.Card) error
	GetNotifications(userID globalid.ID, page, skip int) ([]*model.Notification, error)
	UpdateNotificationsSeen(ids []globalid.ID) error
	UpdateNotificationsOpened(ids []globalid.ID) error
	DeleteCard(id globalid.ID) error
	SaveOAuthAccount(oaa *model.OAuthAccount) error
	GetOAuthAccountBySubject(subject string) (*model.OAuthAccount, error)
	GetOnSwitches() ([]*model.FeatureSwitch, error)
	LatestForType(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error)
	SaveNotification(m *model.Notification) error
	GetAnonymousAliasByUsername(username string) (*model.AnonymousAlias, error)
	DeleteNotificationsForCard(cardID globalid.ID) error
	UnseenNotificationsCount(userID globalid.ID) (int, error)
	GetInvitesForUser(id globalid.ID) ([]*model.Invite, error)
	GetAnonymousAliasLastUsed(userID, threadRootID globalid.ID) (bool, error)
	SaveWaitlistEntry(waiter *model.WaitlistEntry) error
	GetInvite(id globalid.ID) (*model.Invite, error)
	GetEngagement(cardID globalid.ID) (*model.Engagement, error)
	GetEngagements(cardIDs []globalid.ID) (map[globalid.ID]*model.Engagement, error)
	DeleteWaitlistEntry(email string) error
	SaveMention(m *model.Mention) error
	GetMentionExportData(n *model.Notification) (*datastore.MentionNotificationExportData, error)
	SaveNotificationMention(m *model.NotificationMention) error
	GetCommentExportData(n *model.Notification) (*datastore.CommentNotificationExportData, error)
	GetLikeNotificationExportData(n *model.Notification) (*datastore.LikeNotificationExportData, error)
	SaveNotificationComment(m *model.NotificationComment) error
	SubscribeToCard(userID, cardID globalid.ID, typ string) error
	UnsubscribeFromCard(userID, cardID globalid.ID, typ string) error
	SubscribersForCard(cardID globalid.ID, typ string) ([]globalid.ID, error)
	SubscribedToTypes(userID, cardID globalid.ID) ([]string, error)
	SubscribedToCards(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error)
	DeleteMentionsForCard(cardID globalid.ID) error
	ReassignInviterForGroup(rootInviteID, newUserID globalid.ID) error
	GroupInvitesByToken(tokens []string, groupID globalid.ID) error
	ReassignInvitesByToken(tokens []string, userID globalid.ID) error
	SaveFollower(followerID, followeeID globalid.ID) error
	DeleteFollower(followerID, followeeID globalid.ID) error
	GetFollowing(userID globalid.ID) ([]*model.User, error)
	ClearEmptyNotifications() error
	SaveNotificationFollow(notifID, followerID, followeeID globalid.ID) error
	GetFollowExportData(n *model.Notification) (*datastore.FollowNotificationExportData, error)
	IsFollowing(followerID, followeeID globalid.ID) (bool, error)
	IsFollowings(followerID globalid.ID, followeeID []globalid.ID) (map[globalid.ID]bool, error)
	GetNotification(notifID globalid.ID) (*model.Notification, error)
	BlockUser(blockingUser, blockedUser globalid.ID) error
	BlockAnonUserInThread(blockingUser, blockedAlias, threadID globalid.ID) error
	FeaturedCommentsForUser(userID, cardID globalid.ID) ([]*model.Card, error)
	FeaturedCommentsForUserByCardIDs(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.Card, error)
	GetRankedImmediateReplies(cardID, forUser globalid.ID) ([]*model.Card, error)
	GetFlatReplies(id, forUser globalid.ID, latestFirst bool, limitTo int) ([]*model.Card, error)
	GetChannel(id globalid.ID) (*model.Channel, error)
	GetChannelsForUser(userID globalid.ID) ([]*model.Channel, error)
	GetChannelsByID(ids []globalid.ID) ([]*model.Channel, error)
	GetSubscribedChannels(userID globalid.ID) ([]*model.Channel, error)
	AddUserToDefaultChannels(userID globalid.ID) error
	GetCardsForChannel(channelID globalid.ID, count, skip int, forUser globalid.ID) ([]*model.Card, error)
	JoinChannel(userID, channelID globalid.ID) error
	LeaveChannel(userID, channelID globalid.ID) error
	MuteChannel(userID, channelID globalid.ID) error
	UnmuteChannel(userID, channelID globalid.ID) error
	MuteUser(userID, mutedUserID globalid.ID) error
	UnmuteUser(userID, mutedUserID globalid.ID) error
	MuteThread(userID, threadRootID globalid.ID) error
	UnmuteThread(userID, threadRootID globalid.ID) error
	FollowDefaultUsers(userID globalid.ID) error
	AddCardsToTopOfFeed(userID globalid.ID, cardRankedIDs []globalid.ID) error
	ResetUserFeedTop(userID globalid.ID) error
	GetFeedCardsFromCurrentTop(userID globalid.ID, perPage, page int) ([]*model.Card, error)
	GetFeedCardsFromCurrentTopWithQuery(userID globalid.ID, perPage, page int, searchString string) ([]*model.Card, error)
	SetCardVisited(userID, cardID globalid.ID) error
	NewContentAvailableForUser(userID, cardID globalid.ID) (bool, error)
	NewContentAvailableForUserByCards(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]bool, error)
	SetFeedLastUpdatedForUser(userID globalid.ID, t time.Time) error
	GetIntroCardIDs() ([]globalid.ID, error)
	SaveChannel(m *model.Channel) error
	ClearEmptyNotificationsForUser(userID globalid.ID) error
	SaveScoreModification(m *model.ScoreModification) error

	GetPopularRankCardsForUser(userID globalid.ID, page, perPage int) ([]*model.Card, error)
	UpdatePopularRanksForUser(userID globalid.ID) error

	AwardCoins(userID globalid.ID, amount int64) error
	PayCoinAmount(userID globalid.ID, amount int64) error
	AwardTemporaryCoins(userID globalid.ID, amount int64) error
	GetCurrentBalance(userID globalid.ID) (*model.CoinBalances, error)

	GetAnnouncementForNotification(n *model.Notification) (*model.Announcement, error)
	LeaveAllChannels(userID globalid.ID) error

	GetSafeUsersByPage(forUser globalid.ID, pageSize, pageNumber int, searchString string) ([]*model.User, error)

	GetChannelByHandle(handle string) (*model.Channel, error)
	GetIsSubscribed(userID, channelID globalid.ID) (bool, error)
	IsSubscribedToChannels(userID globalid.ID, channelIDs []globalid.ID) (map[globalid.ID]bool, error)

	CountPostsByAliasInThread(aliasID, threadRootID globalid.ID) (int, error)

	SaveReactionForNotification(notifID, userID, cardID globalid.ID) error
	GetUserReaction(userID, cardID globalid.ID) (*model.UserReaction, error)
	GetUserReactions(userID globalid.ID, cardIDs []globalid.ID) (map[globalid.ID]*model.UserReaction, error)
	DeleteUserReactionForType(userID, cardID globalid.ID, typ model.UserReactionType) (int64, error)
	SaveUserReaction(m *model.UserReaction) error
	DeleteReactionForNotification(notifID, userID, cardID globalid.ID) error

	SavePopularRank(m *model.PopularRankEntry) error
	UpdatePopularRankForCard(cardID globalid.ID, viewCountChange, upCountChange, downCountChange, commentCountChange int64, scoreModChange float64) error
	UpdateViewsForCards(cardIDs []globalid.ID) error

	GetCoinsReceivedNotificationData(notif *model.Notification) (*model.CoinReward, error)
	GetLeaderboardRankings(count, skip int) ([]*model.LeaderboardRank, error)
	CountCardsInFeed(userID globalid.ID) (int64, error)
	BuildInitialFeed(userID globalid.ID) error
	GetRankableCardsForUser(userID globalid.ID) ([]*model.PopularRankEntry, error)
	UpdateCardRanksForUser(userID globalid.ID, cardIDs []globalid.ID) error

	GetChannelInfos(userID globalid.ID) ([]*model.ChannelUserInfo, error)
	GetSubscriberCount(channelID globalid.ID) (int, error)

	GetLeaderboardNotificationExportData(notifID globalid.ID) (*datastore.LeaderboardNotificationExportData, error)

	UpdateUniqueCommentersForCard(cardID globalid.ID) error
	UpdateAllNotificationsSeen(userID globalid.ID) error

	GetPopularPostNotificationExportData(n *model.Notification) (*datastore.PopularPostNotificationExportData, error)

	GetCardsForPopularRankSince(t time.Time) ([]*model.PopularRankEntry, error)
	UpdatePopularRanksWithList(userID globalid.ID, ids []globalid.ID) error

	SaveUserTip(m *model.UserTip) error
	AssignAliasForUserTipsInThread(userID, threadRootID, aliasID globalid.ID) error
}
