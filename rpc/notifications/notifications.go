package notifications

import (
	"database/sql"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store"
	"github.com/pkg/errors"
)

const andText = " and "

type NotificationUser interface {
	DisplayName() string
}

type Notifications struct {
	store              *store.Store
	systemProfileImage string
	unitsPerCoin       int64
}

var ErrNotificationEmpty = errors.New("Notification Empty")

func formatName(name string, anon bool) string {
	r := name
	if anon {
		r = "!" + r
	}
	return "**" + r + "**"
}

func emphasize(t string) string {
	return fmt.Sprintf("**%v**", t)
}

func formatNames(users []NotificationUser) string {
	// "Bob", "Bob and Joe", "Bob, Joe, and 2 others" text
	text := ""
	numUsers := len(users)

	if numUsers >= 1 {
		text = emphasize(users[0].DisplayName())
	}

	if numUsers >= 2 {
		if numUsers == 2 {
			text += andText
		} else {
			text += ", "
		}
		text += emphasize(users[1].DisplayName())
	}

	if numUsers >= 3 {
		text += fmt.Sprintf(", and %v other", numUsers-2)
	}

	// turns "other" into "others" for 4+ boosters
	if numUsers >= 4 {
		text += "s"
	}

	return text
}

func NewNotifications(store *store.Store, systemProfileImage string, unitsPerCoin int64) *Notifications {
	return &Notifications{
		store:              store,
		systemProfileImage: systemProfileImage,
		unitsPerCoin:       unitsPerCoin,
	}
}

func (ns *Notifications) ExportNotification(n *model.Notification) (*model.ExportedNotification, error) {
	switch n.Type {
	case model.CommentType:
		return ns.exportCommentNotification(n)
	case model.MentionType:
		return ns.exportMentionNotification(n)
	case model.BoostType:
		return ns.exportBoostNotification(n)
	case model.InviteAcceptedType:
		return ns.exportInviteAcceptedNotification(n)
	case model.AnnouncementType:
		return ns.exportAnnouncementNotification(n)
	case model.IntroductionType:
		return ns.exportIntroductionNotification(n)
	case model.FollowType:
		return ns.exportFollowNotification(n)
	case model.NewInvitesType:
		return ns.exportNewInvitesNotification(n)
	case model.CoinsReceivedType:
		return ns.exportCoinsReceivedNotification(n)
	case model.PopularPostType:
		return ns.exportPopularPostNotification(n)
	case model.FirstPostActivityType:
		return ns.exportFirstPostActivityNotification(n)
	case model.LeaderboardRankType:
		return ns.exportLeaderboardRankNotification(n)
	}

	return nil, errors.New("Export Notification Unknown Type Error")
}

func (ns *Notifications) exportPopularPostNotification(n *model.Notification) (*model.ExportedNotification, error) {
	notification := &model.ExportedNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		ImagePath: ns.systemProfileImage,
		Message:   `💥 BOOM! Your recent posts got lot of attention, here's an extra 20 coins.`,
		Timestamp: n.CreatedAt.Unix(),
		Type:      n.Type,
		Action:    model.OpenWalletAction,
	}

	if n.SeenAt.Valid {
		notification.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		notification.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}
	return notification, nil
}

func (ns *Notifications) exportFirstPostActivityNotification(n *model.Notification) (*model.ExportedNotification, error) {
	notification := &model.ExportedNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		ImagePath: ns.systemProfileImage,
		Message:   `👀 Your posts are starting to get some attention! Here’s an extra 10 coins, post more great stuff!`,
		Timestamp: n.CreatedAt.Unix(),
		Type:      n.Type,
		Action:    model.OpenWalletAction,
	}

	if n.SeenAt.Valid {
		notification.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		notification.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}
	return notification, nil
}

func numberWithOrdinal(n int) string {
	withOrd := strconv.Itoa(n)
	if strings.HasSuffix(withOrd, "11") || strings.HasSuffix(withOrd, "12") || strings.HasSuffix(withOrd, "13") {
		withOrd += "th"
	} else if strings.HasSuffix(withOrd, "1") {
		withOrd += "st"
	} else if strings.HasSuffix(withOrd, "2") {
		withOrd += "nd"
	} else if strings.HasSuffix(withOrd, "3") {
		withOrd += "rd"
	} else {
		withOrd += "th"
	}

	return withOrd
}

func (ns *Notifications) exportLeaderboardRankNotification(n *model.Notification) (*model.ExportedNotification, error) {
	notifData, err := ns.store.GetLeaderboardNotificationExportData(n.ID)

	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		return nil, ErrNotificationEmpty
	} else if err != nil {
		return nil, err
	}

	message := "🏅 You made the leaderboard on October today! You earned 10 Coins 💰💰💰. Think you can make the top 10 today?"

	if notifData.Rank == 1 {
		message = "🥇 Wow, you were the #1 top contributor on October yesterday! You earned 100 Coins 💰💰💰. Can you defend your position?"
	} else if notifData.Rank == 2 {
		message = "🥈 Game on! You made it to the runner-up slot on the leaderboard. You earned 50 Coins 💰💰. Can you make it to the top spot today?"
	} else if notifData.Rank == 3 {
		message = "🥉 Woo! You made third place on the leaderboard yesterday. You earned 30 Coins 💰💰💰. Can you make it to the top spot today?"
	} else if notifData.Rank > 3 && notifData.Rank <= 10 {
		message = fmt.Sprintf("🏅 Wow, you were the %v top contributor on October yesterday! You earned 20 Coins 💰💰💰. Can you hit the Top 3 today?", numberWithOrdinal(int(notifData.Rank)))
	}

	notification := &model.ExportedNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		ImagePath: ns.systemProfileImage,
		Message:   message,
		Timestamp: n.CreatedAt.Unix(),
		Type:      n.Type,
		Action:    model.NavigateToAction,
		ActionData: map[string]string{
			"mobileRouteName": "leaderboard",
			"webRouteName":    "leaderboard",
		},
	}

	if n.SeenAt.Valid {
		notification.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		notification.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}
	return notification, nil
}

func (ns *Notifications) exportCoinsReceivedNotification(n *model.Notification) (*model.ExportedNotification, error) {
	var message string
	var pluralizedCoins string

	notifData, err := ns.store.GetCoinsReceivedNotificationData(n)

	if err != nil {
		return nil, err
	}

	coinAmount := int64(math.Floor(float64(notifData.CoinsReceived) / float64(ns.unitsPerCoin)))

	if coinAmount > 1 {
		pluralizedCoins = fmt.Sprintf("%v Coins", coinAmount)
	} else {
		pluralizedCoins = "a Coin"
	}

	if !notifData.LastRewardedOn.Valid {
		// send a kickoff notification
		message = fmt.Sprintf("You got %v from your past posts and reactions!", pluralizedCoins)
	} else {
		// send a regular daily notification
		message = fmt.Sprintf("You got %v from your activity yesterday!", pluralizedCoins)
	}

	notification := &model.ExportedNotification{
		ID:         n.ID,
		UserID:     n.UserID,
		ImagePath:  ns.systemProfileImage,
		Message:    message,
		Timestamp:  n.CreatedAt.Unix(),
		Type:       n.Type,
		Action:     model.OpenWalletAction,
		ActionData: map[string]string{},
	}

	if n.SeenAt.Valid {
		notification.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		notification.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}
	return notification, nil
}

func (ns *Notifications) exportNewInvitesNotification(n *model.Notification) (*model.ExportedNotification, error) {
	if n.Type != model.NewInvitesType {
		return nil, errors.New("wrong notification type")
	}

	notification := &model.ExportedNotification{
		ID:         n.ID,
		UserID:     n.UserID,
		ImagePath:  ns.systemProfileImage,
		Message:    `You’ve received new invites to October. Invite someone awesome!`,
		Timestamp:  n.CreatedAt.Unix(),
		Type:       n.Type,
		Action:     model.OpenInvitesAction,
		ActionData: map[string]string{},
	}
	if n.SeenAt.Valid {
		notification.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		notification.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}
	return notification, nil
}

func (ns *Notifications) exportCommentNotification(n *model.Notification) (*model.ExportedNotification, error) {
	notifData, err := ns.store.GetCommentExportData(n)

	if errors.Cause(err) != sql.ErrNoRows && err != nil {
		return nil, err
	} else if errors.Cause(err) == sql.ErrNoRows || len(notifData.Comments) <= 0 {
		return nil, ErrNotificationEmpty
	}

	latestComment := notifData.Comments[0]

	// "Bob", "Bob and Joe", "Bob, Joe, and 2 others" text
	nU := make([]NotificationUser, len(notifData.Comments))
	for i, v := range notifData.Comments {
		nU[i] = NotificationUser(v)
	}
	commentersText := formatNames(nU)

	// your post or someone else's
	also := " also"
	ownerName := notifData.PosterName + "'s"

	if notifData.RootOwnerID == n.UserID {
		also = ""
		ownerName = "your"
	}

	if notifData.RootOwnerID == latestComment.AuthorID && !latestComment.IsAnonymous && !notifData.RootIsAnonymous {
		also = ""
		ownerName = "their"
	}

	// ": 'Content'" text
	sanitized := sanitizeAndTruncateContent(notifData.CardContent, 80)
	nonWhitespaceRegex := regexp.MustCompile(`[^\s]+`)
	postTail := "."
	if nonWhitespaceRegex.FindStringIndex(sanitized) != nil {
		postTail = fmt.Sprintf(`: "%v"`, sanitized)
	}

	message := fmt.Sprintf("%v%v commented on %v post%v", commentersText, also, ownerName, postTail)

	actionData := map[string]string{
		"threadRootID":        notifData.ThreadRootID.String(),
		"commentCardID":       notifData.LatestCommentID.String(),
		"commentCardUsername": latestComment.AuthorUsername,
	}

	if notifData.LatestCommentParentID != globalid.Nil && notifData.LatestCommentParentID != notifData.ThreadRootID {
		actionData["parentCommentID"] = notifData.LatestCommentParentID.String()
	}

	eNotif := model.ExportedNotification{
		ID:         n.ID,
		UserID:     n.UserID,
		ImagePath:  latestComment.ImagePath,
		Message:    message,
		Timestamp:  latestComment.Timestamp.Unix(),
		Type:       n.Type,
		Action:     model.OpenThreadAction,
		ActionData: actionData,
	}

	if n.SeenAt.Valid {
		eNotif.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		eNotif.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}

	return &eNotif, nil
}

func (ns *Notifications) exportBoostNotification(n *model.Notification) (*model.ExportedNotification, error) {
	notifData, err := ns.store.GetLikeNotificationExportData(n)

	if errors.Cause(err) != sql.ErrNoRows && err != nil {
		return nil, err
	} else if errors.Cause(err) == sql.ErrNoRows || len(notifData.Boosts) <= 0 {
		return nil, ErrNotificationEmpty
	}

	latestBoost := notifData.Boosts[0]

	// "Bob", "Bob and Joe", "Bob, Joe, and 2 others" text
	nU := make([]NotificationUser, len(notifData.Boosts))
	for i, v := range notifData.Boosts {
		nU[i] = NotificationUser(v)
	}
	boostersText := formatNames(nU)

	postOrComment := "post"
	if notifData.IsComment {
		postOrComment = "comment"
	}

	sanitized := sanitizeAndTruncateContent(notifData.CardContent, 80)
	nonWhitespaceRegex := regexp.MustCompile(`[^\s]+`)
	postTail := "."
	if nonWhitespaceRegex.FindStringIndex(sanitized) != nil {
		postTail = fmt.Sprintf(`: "%v"`, sanitized)
	}

	message := fmt.Sprintf(`%v liked your %v%v`, boostersText, postOrComment, postTail)

	actionData := map[string]string{
		"threadRootID": notifData.ThreadRootID.String(),
	}

	if notifData.IsComment {
		actionData["commentCardID"] = n.TargetID.String()
		actionData["commentCardUsername"] = notifData.AuthorUsername
		if notifData.ThreadReplyID != globalid.Nil && notifData.ThreadReplyID != notifData.ThreadRootID {
			actionData["parentCommentID"] = notifData.ThreadReplyID.String()
		}
	}

	eNotif := model.ExportedNotification{
		ID:         n.ID,
		UserID:     n.UserID,
		ImagePath:  latestBoost.ImagePath,
		Message:    message,
		Timestamp:  latestBoost.Timestamp.Unix(),
		Type:       n.Type,
		Action:     model.OpenThreadAction,
		ActionData: actionData,
	}

	if n.SeenAt.Valid {
		eNotif.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		eNotif.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}

	return &eNotif, nil
}

func (ns *Notifications) exportMentionNotification(notif *model.Notification) (*model.ExportedNotification, error) {
	// get datas
	notifData, err := ns.store.GetMentionExportData(notif)

	if err != nil {
		return nil, err
	}

	taggerName := formatName(notifData.Name, notifData.IsAnonymous)

	notifPicture := notifData.ImagePath

	messageTail := "tagged you in a post."

	if notifData.InComment {
		messageTail = "mentioned you in a comment."
	}

	message := fmt.Sprintf("%v %v", taggerName, messageTail)

	actionData := map[string]string{
		"threadRootID": notifData.ThreadRoot.String(),
	}

	if notifData.InComment {
		actionData["commentCardID"] = notifData.InCard.String()
		actionData["commentCardUsername"] = notifData.InCardAuthorUsername
	}

	if notifData.ThreadReply != globalid.Nil && notifData.ThreadRoot != notifData.ThreadReply {
		actionData["parentCommentID"] = notifData.ThreadReply.String()
	}

	eNotif := model.ExportedNotification{
		ID:         notif.ID,
		UserID:     notif.UserID,
		ImagePath:  notifPicture,
		Message:    message,
		Timestamp:  notif.CreatedAt.Unix(),
		Type:       notif.Type,
		Action:     model.OpenThreadAction,
		ActionData: actionData,
	}

	if notif.SeenAt.Valid {
		eNotif.Seen = notif.SeenAt.Time.Before(time.Now().UTC())
	}

	if notif.OpenedAt.Valid {
		eNotif.Opened = notif.OpenedAt.Time.Before(time.Now().UTC())
	}

	return &eNotif, nil
}

func (ns *Notifications) exportAnnouncementNotification(n *model.Notification) (*model.ExportedNotification, error) {
	if n.Type != model.AnnouncementType {
		return nil, errors.New("Bad Announcement Notification Format ")
	}

	announcement, err := ns.store.GetAnnouncementForNotification(n)
	if err != nil {
		return nil, err
	}

	user, err := ns.store.GetUser(announcement.UserID)
	if err != nil {
		return nil, err
	}

	message := announcement.Message

	if len(message) == 0 {
		card, err := ns.store.GetCard(announcement.CardID)

		if err != nil {
			return nil, err
		}

		content := sanitizeAndTruncateContent(card.Content, 80)

		message = fmt.Sprintf(`**%v** posted a new card "%v""`, user.DisplayName, content)
	}

	eNotif := model.ExportedNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		ImagePath: user.ProfileImagePath,
		Message:   message,
		Timestamp: announcement.CreatedAt.Unix(),
		Type:      n.Type,
	}

	if announcement.CardID != globalid.Nil {
		eNotif.ShowOnCardID = announcement.CardID
		eNotif.Action = model.OpenThreadAction
		eNotif.ActionData = map[string]string{
			"threadRootID": announcement.CardID.String(),
		}
	}

	if n.SeenAt.Valid {
		eNotif.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		eNotif.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}

	return &eNotif, nil
}

func (ns *Notifications) exportInviteAcceptedNotification(n *model.Notification) (*model.ExportedNotification, error) {
	if n.Type != model.InviteAcceptedType {
		return nil, errors.New("Bad invite accepted Notification Format ")
	}

	user, err := ns.store.GetUser(n.TargetID)
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf(`**%v** accepted your invitation to October.`, user.DisplayName)

	eNotif := model.ExportedNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		ImagePath: user.ProfileImagePath,
		Message:   message,
		Timestamp: n.CreatedAt.Unix(),
		Type:      n.Type,
		Action:    model.OpenUserProfileAction,
		ActionData: map[string]string{
			"username": user.Username,
			"id":       user.ID.String(),
		},
	}

	if n.SeenAt.Valid {
		eNotif.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		eNotif.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}

	return &eNotif, nil
}

func (ns *Notifications) exportIntroductionNotification(n *model.Notification) (*model.ExportedNotification, error) {
	if n.Type != model.IntroductionType {
		return nil, errors.New("bad introduction notification format")
	}

	notification := &model.ExportedNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		Timestamp: n.CreatedAt.Unix(),
		Type:      n.Type,
		Action:    model.OpenUserProfileAction,
	}

	if n.TargetID != globalid.Nil {
		inviter, err := ns.store.GetUser(n.TargetID)
		if err != nil {
			return nil, err
		}

		notification.ImagePath = inviter.ProfileImagePath
		notification.Message = fmt.Sprintf(`You were invited to October by **%v**!`, inviter.DisplayName)
		notification.ActionData = map[string]string{
			"username": inviter.Username,
			"id":       inviter.ID.String(),
		}
	} else {
		user, err := ns.store.GetUser(n.UserID)
		if err != nil {
			return nil, err
		}

		notification.ImagePath = ns.systemProfileImage
		notification.Message = "You joined October!"
		notification.ActionData = map[string]string{
			"username": user.Username,
			"id":       user.ID.String(),
		}
	}

	if n.SeenAt.Valid {
		notification.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		notification.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}
	return notification, nil
}

func (ns *Notifications) exportFollowNotification(n *model.Notification) (*model.ExportedNotification, error) {
	notifData, err := ns.store.GetFollowExportData(n)

	if errors.Cause(err) != sql.ErrNoRows && err != nil {
		return nil, err
	} else if errors.Cause(err) == sql.ErrNoRows || len(notifData.Followers) <= 0 {
		return nil, ErrNotificationEmpty
	}

	latestFollow := notifData.Followers[0]

	// "Bob", "Bob and Joe", "Bob, Joe, and 2 others" text
	nU := make([]NotificationUser, len(notifData.Followers))
	for i, v := range notifData.Followers {
		nU[i] = NotificationUser(v)
	}
	followersText := formatNames(nU)

	message := fmt.Sprintf("%v followed you.", followersText)

	eNotif := model.ExportedNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		ImagePath: latestFollow.ImagePath,
		Message:   message,
		Timestamp: n.UpdatedAt.Unix(),
		Type:      n.Type,
		Action:    model.OpenUserProfileAction,
		ActionData: map[string]string{
			"username": latestFollow.Username,
			"id":       latestFollow.ID.String(),
		},
	}

	if n.SeenAt.Valid {
		eNotif.Seen = n.SeenAt.Time.Before(time.Now().UTC())
	}

	if n.OpenedAt.Valid {
		eNotif.Opened = n.OpenedAt.Time.Before(time.Now().UTC())
	}

	return &eNotif, nil
}

// TODO (konrad): Regular expressions that do not contain any meta characters (things
// like `\d`) are just regular strings. Using the `regexp` with such
// expressions is unnecessarily complex and slow. Functions from the
// `bytes` and `strings` packages should be used instead.
func sanitizeAndTruncateContent(content string, length int) string {
	removeURL := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	removeImageAndLink := regexp.MustCompile(`[!]?\[[^\]]*\]\([^)]*\)`)
	removeLinkBlockTags := regexp.MustCompile("%%%\n") // nolint: megacheck
	removeInlineHeaders := regexp.MustCompile(`#+\s+`) // nolint: megacheck
	removeNewlines := regexp.MustCompile("\n")         // nolint: megacheck
	removeSoftbreaks := regexp.MustCompile("%n")       // nolint: megacheck
	removeTrailingWhitespace := regexp.MustCompile(`\s+$`)

	content = removeURL.ReplaceAllString(content, "")
	content = removeImageAndLink.ReplaceAllString(content, "")
	content = removeLinkBlockTags.ReplaceAllString(content, "")
	content = removeInlineHeaders.ReplaceAllString(content, "")
	content = removeNewlines.ReplaceAllString(content, " ")
	content = removeSoftbreaks.ReplaceAllString(content, " ")
	content = removeTrailingWhitespace.ReplaceAllString(content, "")

	if len(content) < length || length == 0 {
		return content
	}

	return content[0:length] + "..."
}
