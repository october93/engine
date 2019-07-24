//go:generate go run ../cmd/rpcgen/rpcgen.go
//go:generate go run ../cmd/rpcgen/subcmd/doc.go ../cmd/rpcgen/subcmd/examples.go

package rpc

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/october93/engine/coinmanager"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	notifExport "github.com/october93/engine/rpc/notifications"
	"github.com/october93/engine/store"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/worker"
	"github.com/october93/engine/worker/emailsender"
)

const (
	Auth                      = "auth"
	Logout                    = "logout"
	Signup                    = "signup"
	ResetPassword             = "resetPassword"
	ValidateInviteCode        = "validateInviteCode"
	AddToWaitlist             = "addToWaitlist"
	GetThread                 = "getThread"
	GetCards                  = "getCards"
	GetCard                   = "get"
	ReactToCard               = "reactToCard"
	VoteOnCard                = "voteOnCard"
	PostCard                  = "postCard"
	NewInvite                 = "newInvite"
	RegisterDevice            = "registerDevice"
	UnregisterDevice          = "unregisterDevice"
	UpdateSettings            = "updateSettings"
	GetUser                   = "getUser"
	ValidateUsername          = "validateUsername"
	GetNotifications          = "getNotifications"
	UpdateNotifications       = "updateNotifications"
	GetAnonymousHandle        = "getAnonymousHandle"
	DeleteCard                = "deleteCard"
	FollowUser                = "followUser"
	UnfollowUser              = "unfollowUser"
	GetFollowingUsers         = "getFollowing"
	GetPostsForUser           = "getPostsForUser"
	GetTags                   = "getTags"
	GetFeaturesForUser        = "getFeatures"
	PreviewContent            = "previewContent"
	UploadImage               = "uploadImage"
	GetTaggableUsers          = "taggableUsers"
	ConnectUsers              = "connectUsers"
	NewUser                   = "newUser"
	GetUsers                  = "getUsers"
	ModifyCardScore           = "modifyCardScore"
	GetInvites                = "getInvites"
	GetOnboardingData         = "getOnboardingData"
	GetMyNetwork              = "getMyNetwork"
	UnsubscribeFromCard       = "unsubscribeFromCard"
	SubscribeToCard           = "subscribeToCard"
	GroupInvites              = "groupInvites"
	ReportCard                = "reportCard"
	BlockUser                 = "blockUser"
	GetChannels               = "getChannels"
	GetCardsForChannel        = "getCardsForChannel"
	UpdateChannelSubscription = "updateChannelSubscription"
	JoinChannel               = "joinChannel"
	LeaveChannel              = "leaveChannel"
	MuteChannel               = "muteChannel"
	UnmuteChannel             = "unmuteChannel"
	MuteUser                  = "muteUser"
	UnmuteUser                = "unmuteUser"
	MuteThread                = "muteThread"
	UnmuteThread              = "unmuteThread"
	CreateChannel             = "createChannel"
	GetPopularCards           = "getPopularCards"
	GetActionCosts            = "getActionCosts"
	UseInviteCode             = "useInviteCode"
	RequestValidation         = "requestValidation"
	ConfirmValidation         = "confirmValidation"
	ValidateChannelName       = "validateChannelName"
	GetChannel                = "getChannel"
	CanAffordAnonymousPost    = "canAffordAnonymousPost"
	GetLeaderboard            = "getLeaderboard"
	SubmitFeedback            = "submitFeedback"
	TipCard                   = "tipCard"
)

const filtered = "[filtered]"

// RPC is the set of API calls for communicating with the backend, specifically
// with the graph and creation, modification and deletion of domain model
// entities.
type RPC interface {
	// Clients
	Auth(ctx context.Context, req AuthRequest) (*AuthResponse, error)
	ResetPassword(ctx context.Context, req ResetPasswordRequest) (*ResetPasswordResponse, error)
	ValidateInviteCode(ctx context.Context, req ValidateInviteCodeRequest) (*ValidateInviteCodeResponse, error)
	AddToWaitlist(ctx context.Context, req AddToWaitlistRequest) (*AddToWaitlistResponse, error)
	// Authenticated Client RPCs
	Logout(ctx context.Context, req LogoutRequest) (*LogoutResponse, error)
	GetThread(ctx context.Context, req GetThreadRequest) (*GetThreadResponse, error)
	GetCards(ctx context.Context, req GetCardsRequest) (*GetCardsResponse, error)
	GetCard(ctx context.Context, req GetCardRequest) (*GetCardResponse, error)
	ReactToCard(ctx context.Context, req ReactToCardRequest) (*ReactToCardResponse, error)
	VoteOnCard(ctx context.Context, req VoteOnCardRequest) (*VoteOnCardResponse, error)
	PostCard(ctx context.Context, req PostCardRequest) (*PostCardResponse, error)
	NewInvite(ctx context.Context, req NewInviteRequest) (*NewInviteResponse, error)
	RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*RegisterDeviceResponse, error)
	UnregisterDevice(ctx context.Context, req UnregisterDeviceRequest) (*UnregisterDeviceResponse, error)
	UpdateSettings(ctx context.Context, req UpdateSettingsRequest) (*UpdateSettingsResponse, error)
	GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error)
	ValidateUsername(ctx context.Context, req ValidateUsernameRequest) (*ValidateUsernameResponse, error)
	GetNotifications(ctx context.Context, req GetNotificationsRequest) (*GetNotificationsResponse, error)
	UpdateNotifications(ctx context.Context, req UpdateNotificationsRequest) (*UpdateNotificationsResponse, error)
	GetAnonymousHandle(ctx context.Context, Params GetAnonymousHandleRequest) (*GetAnonymousHandleResponse, error)
	DeleteCard(ctx context.Context, req DeleteCardRequest) (*DeleteCardResponse, error)
	FollowUser(ctx context.Context, req FollowUserRequest) (*FollowUserResponse, error)
	UnfollowUser(ctx context.Context, req UnfollowUserRequest) (*UnfollowUserResponse, error)
	GetFollowingUsers(ctx context.Context, req GetFollowingUsersRequest) (*GetFollowingUsersResponse, error)
	GetPostsForUser(ctx context.Context, req GetPostsForUserRequest) (*GetPostsForUserResponse, error)
	GetTags(ctx context.Context, req GetTagsRequest) (*GetTagsResponse, error)
	GetFeaturesForUser(ctx context.Context, req GetFeaturesForUserRequest) (*GetFeaturesForUserResponse, error)
	PreviewContent(ctx context.Context, req PreviewContentRequest) (*PreviewContentResponse, error)
	UploadImage(ctx context.Context, req UploadImageRequest) (*UploadImageResponse, error)
	GetTaggableUsers(ctx context.Context, req GetTaggableUsersRequest) (*GetTaggableUsersResponse, error)
	ModifyCardScore(ctx context.Context, req ModifyCardScoreRequest) (*ModifyCardScoreResponse, error)
	GetInvites(ctx context.Context, req GetInvitesRequest) (*GetInvitesResponse, error)
	GetOnboardingData(ctx context.Context, req GetOnboardingDataRequest) (*GetOnboardingDataResponse, error)
	GetMyNetwork(ctx context.Context, req GetMyNetworkRequest) (*GetMyNetworkResponse, error)
	UnsubscribeFromCard(ctx context.Context, req UnsubscribeFromCardRequest) (*UnsubscribeFromCardResponse, error)
	SubscribeToCard(ctx context.Context, req SubscribeToCardRequest) (*SubscribeToCardResponse, error)
	GroupInvites(ctx context.Context, req GroupInvitesRequest) (*GroupInvitesResponse, error)
	ReportCard(ctx context.Context, req ReportCardRequest) (*ReportCardResponse, error)
	BlockUser(ctx context.Context, req BlockUserRequest) (*BlockUserResponse, error)
	GetCardsForChannel(ctx context.Context, req GetCardsForChannelRequest) (*GetCardsForChannelResponse, error)
	UpdateChannelSubscription(ctx context.Context, req UpdateChannelSubscriptionRequest) (*UpdateChannelSubscriptionResponse, error)
	GetChannels(ctx context.Context, req GetChannelsRequest) (*GetChannelsResponse, error)
	JoinChannel(ctx context.Context, req JoinChannelRequest) (*JoinChannelResponse, error)
	LeaveChannel(ctx context.Context, req LeaveChannelRequest) (*LeaveChannelResponse, error)
	MuteChannel(ctx context.Context, req MuteChannelRequest) (*MuteChannelResponse, error)
	UnmuteChannel(ctx context.Context, req UnmuteChannelRequest) (*UnmuteChannelResponse, error)
	MuteUser(ctx context.Context, req MuteUserRequest) (*MuteUserResponse, error)
	UnmuteUser(ctx context.Context, req UnmuteUserRequest) (*UnmuteUserResponse, error)
	MuteThread(ctx context.Context, req MuteThreadRequest) (*MuteThreadResponse, error)
	UnmuteThread(ctx context.Context, req UnmuteThreadRequest) (*UnmuteThreadResponse, error)
	CreateChannel(ctx context.Context, req CreateChannelRequest) (*CreateChannelResponse, error)
	GetPopularCards(ctx context.Context, req GetPopularCardsRequest) (*GetPopularCardsResponse, error)
	GetActionCosts(ctx context.Context, req GetActionCostsRequest) (*GetActionCostsResponse, error)
	UseInviteCode(ctx context.Context, req UseInviteCodeRequest) (*UseInviteCodeResponse, error)
	RequestValidation(ctx context.Context, req RequestValidationRequest) (*RequestValidationResponse, error)
	ConfirmValidation(ctx context.Context, req ConfirmValidationRequest) (*ConfirmValidationResponse, error)
	ValidateChannelName(ctx context.Context, req ValidateChannelNameRequest) (*ValidateChannelNameResponse, error)
	GetChannel(ctx context.Context, req GetChannelRequest) (*GetChannelResponse, error)
	CanAffordAnonymousPost(ctx context.Context, req CanAffordAnonymousPostRequest) (*CanAffordAnonymousPostResponse, error)
	GetLeaderboard(ctx context.Context, req GetLeaderboardRequest) (*GetLeaderboardResponse, error)
	SubmitFeedback(ctx context.Context, req SubmitFeedbackRequest) (*SubmitFeedbackResponse, error)
	TipCard(ctx context.Context, req TipCardRequest) (*TipCardResponse, error)

	// Benchmarking
	ConnectUsers(ctx context.Context, req ConnectUsersRequest) (*ConnectUsersResponse, error)
	NewUser(ctx context.Context, req NewUserRequest) (*NewUserResponse, error)
	GetUsers(ctx context.Context, req GetUsersRequest) (*GetUsersResponse, error)
}

type workerQueue interface {
	EnqueueMailJob(job *emailsender.Job) error
}

type imageProcessor interface {
	SaveBase64CardImage(data string) (string, string, error)
	SaveBase64CardContentImage(data string) (string, string, error)
	SaveBase64ProfileImage(data string) (string, string, error)
	GenerateDefaultProfileImage() (string, string, error)
	SaveBase64CoverImage(data string) (string, string, error)
	DownloadProfileImage(url string) (string, string, error)
	DownloadCardImage(url string) (string, string, error)
	BlendImage(src, gradient string) (string, error)
	GradientImage(gradient string) (string, string, error)
}

type rpc struct {
	store          dataStore
	worker         workerQueue
	imageProcessor imageProcessor
	oauth2         OAuth2
	notifier       *worker.Notifier
	notifications  notifications
	mailTemplates  *emailsender.MailTemplates
	pusher         pusher
	settings       *model.Settings
	log            log.Logger
	config         *Config
	indexer        indexer
	coinmanager    *coinmanager.CoinManager
	responses      responses
}

// NewRPC returns a new instance of the server-side implementation of RPC.
func NewRPC(s dataStore, wq workerQueue, oa2 OAuth2, ip imageProcessor, c *Config, l log.Logger, mt *emailsender.MailTemplates, n *worker.Notifier, p pusher, settings *model.Settings, notifications notifications, i indexer, cm *coinmanager.CoinManager, r responses) RPC {
	return &rpc{
		store:          s,
		worker:         wq,
		oauth2:         oa2,
		imageProcessor: ip,
		config:         c,
		log:            l,
		mailTemplates:  mt,
		notifier:       n,
		notifications:  notifications,
		pusher:         p,
		settings:       settings,
		indexer:        i,
		coinmanager:    cm,
		responses:      r,
	}
}

type AuthRequest struct {
	Params  AuthParams
	Session *model.Session
}

type AuthParams struct {
	// Username is used for logging the user in.
	Username string `json:"username,omitempty"`
	// Password is used for logging the user in.
	Password string `json:"password,omitempty"`
	// Email is used for signin up a new user.
	Email string `json:"email,omitempty"`

	// AccessToken is used for signing in or signing up with Facebook OAuth.
	AccessToken string `json:"accessToken,omitempty"`
	// ResetToken is used for signing in after requesting to reset the password.
	ResetToken string `json:"resetToken,omitempty"`
	// InviteToken is mandatory for signing up.
	InviteToken string `json:"inviteToken,omitempty"`

	// FirstName is used for signing up.
	FirstName string `json:"firstName,omitempty"`
	// LastName is used for signing up.
	LastName string `json:"lastName,omitempty"`
	// ProfileImageData is the profile image for a new user encoded as Base64 string.
	ProfileImageData *string `json:"profilePicture,omitempty"`
	// CoverImageData is the cover image for a new user encoded as Base64 string.
	CoverImageData *string `json:"coverPicture,omitempty"`
	// Signup is used to identify that this is a signup
	IsSignup bool `json:"isSignup"`
}

func (p AuthParams) Sanitize() interface{} {
	if p.Password != "" {
		p.Password = filtered
	}
	if p.AccessToken != "" {
		p.AccessToken = filtered
	}
	return p
}

func (p AuthParams) Validate() error {
	return nil
}

type AuthResponse struct {
	// User is the user now authenticated with this connection.
	User *model.ExportedUser `json:"user"`
	// Session identifies the session of this connecetion. The session ID can
	// be used to authenticate immediately on connection establishment.
	Session *model.Session `json:"session"`
}

// Auth authenticates a user by either logging them in or signing them up with a new account.
func (r *rpc) Auth(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	if req.Params.ResetToken != "" {
		return r.resetTokenLogin(ctx, req)
	}
	if req.Params.AccessToken != "" {
		return r.accessTokenAuthentication(ctx, req)
	}

	var user *model.User
	var err error
	if strings.Contains(req.Params.Username, "@") {
		user, err = r.store.GetUserByEmail(req.Params.Username)
	} else {
		user, err = r.store.GetUserByUsername(req.Params.Username)
	}
	// user does not exist, attempt to sign up
	if errors.Cause(err) == sql.ErrNoRows && (req.Params.InviteToken != "" || req.Params.IsSignup) {
		return r.signup(ctx, req, nil)
	} else if errors.Cause(err) == sql.ErrNoRows && req.Params.InviteToken == "" && !req.Params.IsSignup {
		return nil, ErrWrongPassword
	} else if err != nil {
		return nil, err
	}

	matches, err := user.PasswordMatches(req.Params.Password)
	if err != nil {
		return nil, err
	}
	if !matches {
		return nil, ErrWrongPassword
	}

	if user.BlockedAt.Valid {
		return nil, ErrUserBlocked
	}

	req.Session.Identify(user)
	err = r.store.SaveSession(req.Session)
	if err != nil {
		return nil, err
	}
	result := &AuthResponse{
		User:    user.Export(user.ID),
		Session: req.Session,
	}
	return result, err
}

func (r *rpc) resetTokenLogin(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	user, err := r.store.GetUserByEmail(req.Params.Username)
	if err != nil || user == nil {
		return nil, ErrCouldNotRetrieveUser(err)
	}
	resetToken, err := r.store.GetResetToken(user.ID)
	if err != nil {
		return nil, ErrInvalidToken()
	}
	resetToken.Token = globalid.ID(req.Params.ResetToken)
	err = resetToken.Valid()
	if err != nil {
		return nil, err
	}
	req.Session.Identify(user)
	err = r.store.SaveSession(req.Session)
	if err != nil {
		return nil, err
	}
	result := &AuthResponse{
		User:    user.Export(user.ID),
		Session: req.Session,
	}
	return result, nil
}

func (r *rpc) signup(ctx context.Context, req AuthRequest, fbUser *FacebookUser) (res *AuthResponse, err error) {
	if r.settings.SignupsFrozen {
		return nil, errors.New("signups are frozen")
	}

	/* Check the user's invite code */
	var invite *model.Invite

	if req.Params.InviteToken != "" {
		// Find/validate the invite token used
		invite, err = r.store.GetInviteByToken(strings.ToUpper(req.Params.InviteToken))
		if err != nil {
			r.log.Info("Error validating invite code")
			r.log.Info(err.Error())
		}
	}

	if invite == nil && !req.Params.IsSignup {
		return nil, model.ErrInvalidInviteCode
	}

	/* Create the user or link to a duplicate */

	var user, duplicate *model.User
	if fbUser != nil {
		// check if there is already an account associated with this email address
		duplicate, err = r.store.GetUserByEmail(fbUser.Email)
		if err != nil && errors.Cause(err) != sql.ErrNoRows {
			return nil, err
		}
		// account with this email address is already exists, associate OAuth
		// account with this one.
		if duplicate != nil {
			return r.handleDuplicate(ctx, req, duplicate, fbUser)
		}
		user, err = r.facebookUser(ctx, req, fbUser)
		if err != nil {
			return nil, err
		}
	} else {
		user, err = r.user(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	user.CoinRewardLastUpdatedAt = model.NewDBTime(time.Now().UTC())

	// Save the user now so that it exists in the DB
	err = r.store.SaveUser(user)
	if err != nil {
		return nil, err
	}

	// delete waitlist entry if exists
	err = r.store.DeleteWaitlistEntry(user.Email)
	if err != nil {
		r.log.Error(err)
	}

	/* Create the welcome notif now */
	welcomeNotif := &model.Notification{
		ID:     globalid.Next(),
		UserID: user.ID,
		Type:   model.IntroductionType,
	}

	/* If you used an invite code */
	if invite != nil {
		// Set which invite was used
		user.JoinedFromInvite = invite.ID

		//if the inviter is shadowbanned, shadowban this user too
		var invitingUser *model.User
		invitingUser, err = r.store.GetUser(invite.NodeID)
		if err != nil {
			return nil, err
		}
		if invitingUser != nil && invitingUser.ShadowbannedAt.Valid {
			user.ShadowbannedAt = model.NewDBTime(time.Now().UTC())
		}

		// Reassign grouped invites
		if invite.GroupID != globalid.Nil {
			err = r.store.ReassignInviterForGroup(invite.ID, user.ID)
			if err != nil {
				return nil, err
			}
		}

		// subscribe to your inviter
		err = r.store.SaveFollower(user.ID, invite.NodeID)
		if err != nil {
			r.log.Error(err)
		}

		// Join the invite's channel if this is a channel invite
		if invite.ChannelID != globalid.Nil {
			err = r.store.JoinChannel(user.ID, invite.ChannelID)
			if err != nil {
				r.log.Error(err)
			}
		} else {
			err = r.store.AddUserToDefaultChannels(user.ID)
			if err != nil {
				r.log.Error(err)
			}
		}

		// "system invites" have no real inviter, don't give code rewards, don't send invites
		if !invite.SystemInvite {
			// assign the inviter for the welcome notification
			welcomeNotif.TargetID = invitingUser.ID
			// notify inviter of accepted inivte
			notif := &model.Notification{
				ID:       globalid.Next(),
				UserID:   invite.NodeID,
				TargetID: user.ID,
				Type:     model.InviteAcceptedType,
			}
			err = r.store.SaveNotification(notif)
			if err != nil {
				return nil, err
			}
			exNotif, eerr := r.notifications.ExportNotification(notif)
			if eerr != nil {
				return nil, eerr
			}

			err = r.pusher.NewNotification(ctx, req.Session, exNotif)
			if err != nil {
				r.log.Error(err)
			}
			err = r.notifier.NotifyPush(exNotif)
			if err != nil {
				r.log.Error(err)
			}

			/* reward inviter with tokens */
			nB, nberr := r.tokensForInvite(invite.NodeID)
			if nberr != nil {
				r.log.Error(nberr)
			}
			err = r.pusher.UpdateCoinBalance(ctx, invite.NodeID, nB)
			if err != nil {
				r.log.Error(err)
			}
		}
	} else { // if you didn't use an invite code
		// subscribe to default channels
		err = r.store.AddUserToDefaultChannels(user.ID)
		if err != nil {
			r.log.Error(err)
		}
	}

	// subscribe to default users
	err = r.store.FollowDefaultUsers(user.ID)
	if err != nil {
		r.log.Error(err)
	}

	// Save the welcome notification
	err = r.store.SaveNotification(welcomeNotif)
	if err != nil {
		return nil, err
	}

	/* Create the new user's invite code */
	newInvite, ierr := model.NewInvite(user.ID)
	if ierr != nil {
		return nil, ierr
	}

	newInvite.RemainingUses = 1

	serr := r.store.SaveInvite(newInvite)
	if serr != nil {
		return nil, serr
	}

	// save again because reasons
	err = r.store.SaveUser(user)
	if err != nil {
		return nil, err
	}

	err = r.notifier.NotifySlack("engagement", fmt.Sprintf("<https://october.app/user/%s|%s> has signed up.", user.Username, user.Username))
	if err != nil {
		r.log.Error(err)
	}

	/* Log the new user in and return */

	req.Session.Identify(user)
	err = r.store.SaveSession(req.Session)
	if err != nil {
		return nil, err
	}

	// Update token amounts
	nB, err := r.tokensForNewUser(user.ID)
	if err != nil {
		r.log.Error(err)
	}

	if nB != nil {
		user.CoinBalance = nB.CoinBalance
	}

	if invite != nil {
		nB, err = r.updateTokensForRedeemCode(user.ID)
		if err != nil {
			r.log.Error(err)
		}

		if nB != nil {
			user.CoinBalance = nB.CoinBalance
		}
	}

	result := &AuthResponse{
		User:    user.Export(user.ID),
		Session: req.Session,
	}

	/* Index in Algoila and notify slack of successful signup */
	go func() {
		ierr := r.indexer.IndexUser(user)
		if ierr != nil {
			r.log.Error(ierr)
		}
	}()

	return result, nil
}

func (r *rpc) handleDuplicate(ctx context.Context, req AuthRequest, user *model.User, fbUser *FacebookUser) (*AuthResponse, error) {
	oauthAccount := model.NewOAuthAccount(model.FacebookProvider, fbUser.ID, user.ID)
	err := r.store.SaveUser(user)
	if err != nil {
		return nil, err
	}

	err = r.store.SaveOAuthAccount(oauthAccount)
	if err != nil {
		return nil, err
	}
	req.Session.Identify(user)
	err = r.store.SaveSession(req.Session)
	if err != nil {
		return nil, err
	}
	result := &AuthResponse{
		User:    user.Export(user.ID),
		Session: req.Session,
	}
	return result, nil
}

func (r *rpc) accessTokenAuthentication(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	token, err := r.oauth2.ExtendToken(ctx, req.Params.AccessToken)
	if err != nil {
		return nil, err
	}
	fbUser, err := token.FacebookUser()
	if err != nil {
		return nil, err
	}

	var user *model.User
	oauthAccount, err := r.store.GetOAuthAccountBySubject(fbUser.ID)
	if errors.Cause(err) == sql.ErrNoRows {
		return r.signup(ctx, req, fbUser)
	} else if err != nil {
		return nil, err
	} else {
		user, err = r.store.GetUser(oauthAccount.UserID)
		if err != nil {
			return nil, err
		}
	}
	req.Session.Identify(user)
	err = r.store.SaveSession(req.Session)
	if err != nil {
		return nil, err
	}
	result := &AuthResponse{
		User:    user.Export(user.ID),
		Session: req.Session,
	}
	return result, nil
}

func (r *rpc) user(ctx context.Context, req AuthRequest) (*model.User, error) {
	user := model.NewUser(globalid.Next(), req.Params.Username, req.Params.Email, fmt.Sprintf("%s %s", req.Params.FirstName, req.Params.LastName))
	user.FirstName = req.Params.FirstName
	user.LastName = req.Params.LastName

	if req.Params.ProfileImageData != nil {
		url, _, err := r.imageProcessor.SaveBase64ProfileImage(*req.Params.ProfileImageData)
		if err != nil {
			return user, err
		}
		user.ProfileImagePath = url
	} else {
		var url string
		url, _, err := r.imageProcessor.GenerateDefaultProfileImage()
		if err != nil {
			return user, err
		}
		user.ProfileImagePath = url
	}

	if req.Params.CoverImageData != nil {
		var url string
		url, _, err := r.imageProcessor.SaveBase64CoverImage(*req.Params.CoverImageData)
		if err != nil {
			return user, err
		}
		user.CoverImagePath = url
	}
	err := user.SetPassword(req.Params.Password)
	return user, err
}

func (r *rpc) facebookUser(ctx context.Context, req AuthRequest, fbUser *FacebookUser) (*model.User, error) {
	username := req.Params.Username
	if username == "" {
		username = fmt.Sprintf("%s%s", fbUser.FirstName, fbUser.LastName)
		username = strings.ToLower(username)
	}

	email := fbUser.Email

	if email == "" {
		r.log.Error(errors.New("User with no email signed up"))
		email = fmt.Sprintf("%s@example.com", username)
	}

	user := model.NewUser(globalid.Next(), username, email, fmt.Sprintf("%s %s", fbUser.FirstName, fbUser.LastName))
	user.FirstName = fbUser.FirstName
	user.LastName = fbUser.LastName
	if fbUser.ProfileImagePath != "" {
		var url string
		url, _, err := r.imageProcessor.DownloadProfileImage(fbUser.ProfileImagePath)
		if err != nil {
			return user, err
		}
		user.ProfileImagePath = url
	}

	oauthAccount := model.NewOAuthAccount(model.FacebookProvider, fbUser.ID, user.ID)
	err := r.store.SaveUser(user)
	if err != nil {
		return user, err
	}
	err = r.store.SaveOAuthAccount(oauthAccount)
	return user, err
}

type ResetPasswordRequest struct {
	Params  ResetPasswordParams
	Session *model.Session
}

type ResetPasswordParams struct {
	// Email is used to identify the user's account.
	Email string `json:"email"`
}

func (p ResetPasswordParams) Sanitize() interface{} {
	return p
}

func (p ResetPasswordParams) Validate() error {
	if p.Email == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}

type ResetPasswordResponse struct{}

// ResetPassword requests to reset the password of a user. A reset token is
// generated and send to the user's email address. A request to this endpoint
// is always successful in order to prevent account details leaking.
func (r *rpc) ResetPassword(ctx context.Context, req ResetPasswordRequest) (*ResetPasswordResponse, error) {
	user, err := r.store.GetUserByEmail(req.Params.Email)
	if errors.Cause(err) == sql.ErrNoRows {
		// prevent email leakage
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token, err := model.NewResetToken(user.ID)
	if err != nil {
		return nil, err
	}
	err = r.store.SaveResetToken(token)
	if err != nil {
		return nil, err
	}
	data := struct {
		Token *model.ResetToken
		Email string
		Host  string
	}{
		Token: token,
		Email: req.Params.Email,
		Host:  r.config.WebEndpoint,
	}
	job := &emailsender.Job{
		TextTemplate: r.mailTemplates.PasswordReset.TextTemplate,
		HTMLTemplate: r.mailTemplates.PasswordReset.HTMLTemplate,
		From:         "no-reply@october.news",
		Sender:       "October News",
		To:           req.Params.Email,
		Recipient:    "",
		Subject:      "Reset Password",
		Data:         data,
	}
	return nil, r.worker.EnqueueMailJob(job)
}

type LogoutRequest struct {
	Params  LogoutParams
	Session *model.Session
}

type LogoutParams struct{}

func (p LogoutParams) Sanitize() interface{} {
	return p
}

func (p LogoutParams) Validate() error {
	return nil
}

type LogoutResponse struct{}

// Logout deauthenticates the current connection and deletes the session.
func (r *rpc) Logout(ctx context.Context, req LogoutRequest) (*LogoutResponse, error) {
	err := r.store.DeleteSession(req.Session.ID)
	if err != nil {
		return nil, err
	}
	req.Session.SetUser(nil)
	return nil, err
}

type DeleteCardRequest struct {
	Params  DeleteCardParams
	Session *model.Session
}

type DeleteCardParams struct {
	// CardID is the ID of the card to be deleted.
	CardID globalid.ID `json:"cardID"`
}

func (p DeleteCardParams) Sanitize() interface{} {
	return p
}

func (p DeleteCardParams) Validate() error {
	return nil
}

type DeleteCardResponse struct{}

// DeleteCard deletes the given card and makes it inacessible.
func (r *rpc) DeleteCard(ctx context.Context, req DeleteCardRequest) (*DeleteCardResponse, error) {
	card, err := r.store.GetCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}
	isOwnCard := card.OwnerID == req.Session.UserID

	if !isOwnCard && !req.Session.User.Admin {
		return nil, ErrForbidden()
	}
	err = r.store.DeleteNotificationsForCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}

	err = r.store.DeleteMentionsForCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}

	err = r.store.DeleteCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}

	if card.Reply() {
		// unsubcscribe from root card if you're self-deleting and it's a reply
		if isOwnCard {
			err = r.unsubscribeFromNotificationsForType(card.OwnerID, card.ThreadRootID, model.CommentType)
			if err != nil {
				return nil, err
			}
		}

		go r.pushUpdateEngagement(ctx, req.Session, card.ThreadRootID)

		// Popular Rank update
		err = r.store.UpdatePopularRankForCard(card.ThreadRootID, 0, 0, 0, -1, 0)
		if err != nil {
			r.log.Error(err)
		}

		err = r.store.UpdateUniqueCommentersForCard(card.ThreadRootID)
		if err != nil {
			r.log.Error(err)
		}

	}

	return nil, nil
}

type FollowUserRequest struct {
	Params  FollowUserParams
	Session *model.Session
}

type FollowUserParams struct {
	// UserID is the ID of the user to be followed.
	UserID globalid.ID `json:"userID"`
}

func (p FollowUserParams) Sanitize() interface{} {
	return p
}

func (p FollowUserParams) Validate() error {
	return nil
}

type FollowUserResponse struct{}

// FollowUser adds the given user to the list of followers.
func (r *rpc) FollowUser(ctx context.Context, req FollowUserRequest) (*FollowUserResponse, error) {
	if req.Params.UserID == req.Session.UserID {
		return nil, errors.New("can not follow self")
	}

	isFollowing, err := r.store.IsFollowing(req.Session.UserID, req.Params.UserID)
	if err != nil {
		return nil, err
	}

	if !isFollowing {
		err := r.store.SaveFollower(req.Session.UserID, req.Params.UserID)
		if err != nil {
			return nil, err
		}

		err = r.notifyForFollow(ctx, req.Session, req.Params.UserID)
		if err != nil {
			r.log.Error(err)
		}
	}

	return nil, nil
}

func (r *rpc) notifyForFollow(ctx context.Context, session *model.Session, followeeID globalid.ID) error {
	// get latest notification
	notif, err := r.store.LatestForType(followeeID, globalid.Nil, model.FollowType, true)

	// if there isn't one, make a new one
	newNotif := false
	if errors.Cause(err) == sql.ErrNoRows {
		newNotif = true
		notif = &model.Notification{
			ID:     globalid.Next(),
			UserID: followeeID,
			Type:   model.FollowType,
		}

		err = r.store.SaveNotification(notif)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = r.store.SaveNotificationFollow(notif.ID, session.UserID, followeeID)
	if err != nil {
		return err
	}

	if !newNotif {
		notif.UpdatedAt = time.Now().UTC()
		err = r.store.SaveNotification(notif)

		if err != nil {
			return err
		}
	}

	// export
	exNotif, err := r.notifications.ExportNotification(notif)
	if err != nil {
		return err
	}

	// notify via app notification
	err = r.pusher.NewNotification(ctx, session, exNotif)
	if err != nil {
		r.log.Error(err)
	}

	// notify via push
	if newNotif {
		err = r.notifier.NotifyPush(exNotif)
		if err != nil {
			r.log.Error(err)
		}
	}

	return nil
}

type UnfollowUserRequest struct {
	Params  UnfollowUserParams
	Session *model.Session
}

type UnfollowUserParams struct {
	// UserID is the ID of the user to be unfollowed.
	UserID globalid.ID `json:"userID"`
}

func (p UnfollowUserParams) Sanitize() interface{} {
	return p
}

func (p UnfollowUserParams) Validate() error {
	return nil
}

type UnfollowUserResponse struct{}

// UnfollowUser removes the given user of the list of followers.
func (r *rpc) UnfollowUser(ctx context.Context, req UnfollowUserRequest) (*UnfollowUserResponse, error) {
	err := r.store.DeleteFollower(req.Session.UserID, req.Params.UserID)
	if err != nil {
		return nil, err
	}

	return nil, r.store.ClearEmptyNotifications()
}

type GetFollowingUsersRequest struct {
	Params  GetFollowingUsersParams
	Session *model.Session
}

type GetFollowingUsersParams struct{}

func (p GetFollowingUsersParams) Sanitize() interface{} {
	return p
}

func (p GetFollowingUsersParams) Validate() error {
	return nil
}

type GetFollowingUsersResponse []*model.ExportedUser

// GetFollowingUsers returns a list of users which we the user is following.
func (r *rpc) GetFollowingUsers(ctx context.Context, req GetFollowingUsersRequest) (*GetFollowingUsersResponse, error) {
	following, err := r.store.GetFollowing(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	exportedFollowers := make([]*model.ExportedUser, len(following))

	for idx, follower := range following {
		exportedFollowers[idx] = follower.Export(req.Session.UserID)
	}

	response := GetFollowingUsersResponse(exportedFollowers)
	return &response, nil
}

type GetPostsForUserRequest struct {
	Params  GetPostsForUserParams
	Session *model.Session
}

type GetPostsForUserParams struct {
	// Username of the user's posts fo fetch.
	Username string `json:"username"`
	// UserID of the user's posts to fetch.
	UserID globalid.ID `json:"userID"`
	// PageSize specifies the number of cards to be returned per request.
	PageSize int `json:"pageSize"`
	// PageNumber specifies further pages to retrieve.
	PageNumber int `json:"pageNumber"`
}

func (p GetPostsForUserParams) Sanitize() interface{} {
	return p
}

func (p GetPostsForUserParams) Validate() error {
	return nil
}

type GetPostsForUserResponse struct {
	// Cards are the posts made by the requested user.
	Cards []*model.CardResponse `json:"cards"`
	// NextPage returns true if there is another page to fetch.
	NextPage bool `json:"hasNextPage"`
}

// GetPostsForUser retrieves the cards for a certain user. If this request is
// made on behalf of the currently authenticated user, anonymous cards be
// included as well.
func (r *rpc) GetPostsForUser(ctx context.Context, req GetPostsForUserRequest) (*GetPostsForUserResponse, error) {
	id := req.Params.UserID
	var cards []*model.Card
	var nextPage []*model.Card
	var err error
	if req.Params.Username != "" {
		user, uerr := r.store.GetUserByUsername(req.Params.Username)
		if uerr != nil {
			return nil, uerr
		}
		id = user.ID
	}

	if req.Session.UserID == id {
		cards, err = r.store.GetPostedCardsForNodeIncludingAnon(id, req.Params.PageSize, req.Params.PageNumber)

		if err != nil {
			return nil, err
		}
		nextPage, err = r.store.GetPostedCardsForNodeIncludingAnon(id, req.Params.PageSize, req.Params.PageNumber+1)
		if err != nil {
			return nil, err
		}
	} else {
		cards, err = r.store.GetPostedCardsForNode(id, req.Params.PageSize, req.Params.PageNumber)
		if err != nil {
			return nil, err
		}
		nextPage, err = r.store.GetPostedCardsForNode(id, req.Params.PageSize, req.Params.PageNumber+1)
		if err != nil {
			return nil, err
		}
	}
	cardResponses, err := r.cardResponses(cards, req.Session.UserID)
	if err != nil {
		return nil, err
	}
	return &GetPostsForUserResponse{
		Cards:    cardResponses,
		NextPage: len(nextPage) > 0,
	}, nil
}

type GetCardsRequest struct {
	Params  GetCardsParams
	Session *model.Session
}

type GetCardsParams struct {
	// PerPage specifies the number of cards returned per request.
	PerPage int `json:"pageSize"`
	// Page for pagination; the page to retrieve.
	Page int `json:"pageNumber"`
	// Optional search
	SearchString string `json:"searchString"`
}

func (p GetCardsParams) Sanitize() interface{} {
	return p
}

func (p GetCardsParams) Validate() error {
	if p.PerPage == 0 {
		return errors.New("pageSize must be > 0")
	}
	return nil
}

type GetCardsResponse struct {
	NewCardCount int                   `json:"newCardCount,omitempty"`
	Cards        []*model.CardResponse `json:"cards"`
	NextPage     bool                  `json:"hasNextPage"`
}

// GetCards retrieves the feed for the given user.
func (r *rpc) GetCards(ctx context.Context, req GetCardsRequest) (*GetCardsResponse, error) {
	newCardCount := 0
	if req.Params.Page == 0 {
		user, err := r.store.GetUser(req.Session.UserID)
		if err != nil {
			return nil, err
		}

		feedCt, err := r.store.CountCardsInFeed(req.Session.UserID)
		if err != nil {
			return nil, err
		}
		if feedCt <= 0 && !user.FeedLastUpdatedAt.Valid {
			// build initial poprank based feed
			err = r.store.BuildInitialFeed(req.Session.UserID)
			if err != nil {
				return nil, err
			}

			if !user.SeenIntroCards {
				introCardIDs, serr := r.store.GetIntroCardIDs()
				if serr != nil {
					return nil, serr
				}
				err = r.store.AddCardsToTopOfFeed(req.Session.UserID, introCardIDs)
				if err != nil {
					return nil, err
				}

				user.SeenIntroCards = true

				err = r.store.SaveUser(user)
				if err != nil {
					return nil, err
				}
			}

			err = r.store.SetFeedLastUpdatedForUser(req.Session.UserID, time.Now().UTC())
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil {
				return nil, err
			}
			if !user.DisableFeed {
				cards, err := GetMoreCards(r.store, req.Session.UserID)
				if err != nil {
					return nil, err
				}

				if len(cards) > 0 {
					newCardCount = len(cards)
					err = r.store.UpdateViewsForCards(cards)
					if err != nil {
						r.log.Error(err)
					}
				}

				err = r.store.AddCardsToTopOfFeed(req.Session.UserID, cards)
				if err != nil {
					return nil, err
				}

				err = r.store.SetFeedLastUpdatedForUser(req.Session.UserID, time.Now().UTC())
				if err != nil {
					return nil, err
				}
				err = r.store.ResetUserFeedTop(req.Session.UserID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	var cards []*model.Card
	var lookAhead []*model.Card
	var err error

	if req.Params.SearchString != "" {
		cards, err = r.store.GetFeedCardsFromCurrentTopWithQuery(req.Session.UserID, req.Params.PerPage, req.Params.Page, req.Params.SearchString)
		if err != nil {
			return nil, err
		}
		lookAhead, err = r.store.GetFeedCardsFromCurrentTopWithQuery(req.Session.UserID, req.Params.PerPage, req.Params.Page+1, req.Params.SearchString)
		if err != nil {
			return nil, err
		}
	} else {
		cards, err = r.store.GetFeedCardsFromCurrentTop(req.Session.UserID, req.Params.PerPage, req.Params.Page)
		if err != nil {
			return nil, err
		}
		lookAhead, err = r.store.GetFeedCardsFromCurrentTop(req.Session.UserID, req.Params.PerPage, req.Params.Page+1)
		if err != nil {
			return nil, err
		}
	}

	cardResponses, err := r.responses.FeedCardResponses(cards, req.Session.UserID)
	if err != nil {
		return nil, err
	}

	result := &GetCardsResponse{
		Cards:        cardResponses,
		NextPage:     len(lookAhead) != 0,
		NewCardCount: newCardCount,
	}

	return result, nil
}

func GetMoreCards(store dataStore, userID globalid.ID) ([]globalid.ID, error) {
	cardRanks, err := store.GetRankableCardsForUser(userID)
	if err != nil {
		return nil, err
	}

	cardsToAdd, leftoverCards, err := ChooseCardsFromPool(cardRanks)
	if err != nil {
		return nil, err
	}

	return cardsToAdd, store.UpdateCardRanksForUser(userID, leftoverCards)
}

// returns cards to put on feed, cards to return back to the pool, and err
func ChooseCardsFromPool(cardRanks []*model.PopularRankEntry) ([]globalid.ID, []globalid.ID, error) {
	const (
		HandSizeLimit = 10 //maximum number of cards to return
	)

	sort.Slice(cardRanks, func(i, j int) bool {
		return cardRanks[i].Rank() > cardRanks[j].Rank()
	})

	count := 0
	leftoverCards := []globalid.ID{}
	cardsToAdd := []globalid.ID{}
	for _, cardRank := range cardRanks {
		doAdd := count < HandSizeLimit && rand.Float64() < ConfidenceFromRankEntry(cardRank).ProbabilitySurfaced
		if doAdd {
			count++
			cardsToAdd = append(cardsToAdd, cardRank.CardID)
		} else {
			leftoverCards = append(leftoverCards, cardRank.CardID)
		}
	}

	return cardsToAdd, leftoverCards, nil
}

func ConfidenceFromRankEntry(entry *model.PopularRankEntry) model.ConfidenceData {
	const (
		LikePoints          = 1.0          //number of points a Like is worth
		DislikePoints       = 5.0          //number of points a Dislike is worth
		ViewPoints          = 0.05         //number of points a View (without a vote) is worth
		CommentPoints       = 1.5          //number of points a Comment is worth
		GoodThreshold       = 5.0          //minimum number of points to be confident a card is Good
		BadThreshold        = 15.0         //minimum number of points to be confident a card is Bad
		InitialP            = 0.5          //probability a brand new card will be surfaced
		ControversyExponent = 2.0          //higher values slow down the growth of InitialP with confidence
		SteepnessExponent   = 2.0          //higher values change probability more drastically with confidence
		ConfidenceOffset    = 3.8414588206 //this is the z^2 value for a 0.95 confidence interval
	)

	confidence := func(n float64) float64 {
		return n / (n + ConfidenceOffset)
	}

	controversialP := func(n float64) float64 {
		return InitialP + math.Pow(confidence(n), ControversyExponent)*(1-InitialP)
	}

	badLimit := func(n float64) float64 {
		return (1 - confidence(BadThreshold)) * (1 - math.Sqrt(BadThreshold/n))
	}

	badP := func(p, n float64) float64 {
		if n == 0 {
			return controversialP(n)
		}

		badLimit := badLimit(n)
		base := (p - badLimit) / (0.5 - badLimit)
		return controversialP(n) * math.Pow(base, SteepnessExponent)
	}

	goodLimit := func(n float64) float64 {
		return 1 - (1-confidence(GoodThreshold))*(1-math.Sqrt(GoodThreshold/n))
	}

	goodP := func(p, n float64) float64 {
		if n == 0 {
			return controversialP(n)
		}

		goodLimit := goodLimit(n)
		base := (goodLimit - p) / (goodLimit - 0.5)
		return 1 - (1-controversialP(n))*math.Pow(base, SteepnessExponent)
	}

	//turns p (ratio of good engagement / all engagement) into P (probability card is surfaced)
	getP := func(p, n float64) float64 {
		if p > 0.5 {
			if p > goodLimit(n) {
				return 1.0
			}
			return goodP(p, n)
		}

		if p < badLimit(n) {
			return 0.0
		}
		return badP(p, n)
	}

	res := model.ConfidenceData{
		ID:            entry.CardID,
		UpvoteCount:   entry.UpvoteCount,
		DownvoteCount: entry.DownvoteCount,
		CommentCount:  entry.CommentCount,
		ScoreMod:      entry.ScoreMod,
	}

	views := entry.Views - entry.UpvoteCount - entry.DownvoteCount
	if views > 0 { //account for fact that "seen" isn't too accurate right now
		res.ViewCount = views
	}

	var upMod, downMod float64
	if entry.ScoreMod > 0 {
		upMod = entry.ScoreMod
	} else {
		downMod = entry.ScoreMod
	}

	pos := (float64(entry.UpvoteCount)+upMod)*LikePoints + float64(entry.CommentCount)*CommentPoints
	res.EngagementScore = pos + (float64(entry.DownvoteCount)+downMod)*DislikePoints + float64(views)*ViewPoints
	if res.EngagementScore > 0 {
		res.Goodness = pos / res.EngagementScore
	} else { //when n==0, want p==0.5
		res.Goodness = 0.5
	}

	res.Confidence = confidence(res.EngagementScore)
	res.ProbabilitySurfaced = getP(res.Goodness, res.EngagementScore)
	res.Rank = entry.Rank()
	return res
}

type GetCardRequest struct {
	Params  GetCardParams
	Session *model.Session
}

type GetCardParams struct {
	// CardID identifies the card to be retrieved.
	CardID globalid.ID `json:"cardID"`
}

func (p GetCardParams) Sanitize() interface{} {
	return p
}

func (p GetCardParams) Validate() error {
	return nil
}

type GetCardResponse model.CardResponse

// GetCard retrieves a single card.
func (r *rpc) GetCard(ctx context.Context, req GetCardRequest) (*GetCardResponse, error) {
	card, err := r.store.GetCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}

	result, err := CardResponse(r.store, card, req.Session.UserID)
	if err != nil {
		return nil, err
	}
	return (*GetCardResponse)(result), nil
}

type GetThreadRequest struct {
	Params  GetThreadParams
	Session *model.Session
}

type GetThreadParams struct {
	// CardID identifies the root card of the thread.
	CardID globalid.ID `json:"cardID"`
	Nested bool        `json:"nested"`
}

func (p GetThreadParams) Sanitize() interface{} {
	return p
}

func (p GetThreadParams) Validate() error {
	return nil
}

type GetThreadResponse []*model.CardResponse

// GetThread retrieves all cards belonging to a thread ordered by their time of creation.
func (r *rpc) GetThread(ctx context.Context, req GetThreadRequest) (*GetThreadResponse, error) {
	var result []*model.CardResponse
	var err error

	if req.Params.Nested {
		// nested response
		card, verr := r.store.GetCard(req.Params.CardID)
		if verr != nil {
			return nil, verr
		}

		if card.ThreadReplyID == globalid.Nil {
			// top level, get immediate replies with latest comments
			cards, cerr := r.store.GetRankedImmediateReplies(req.Params.CardID, req.Session.UserID)
			fmt.Println(cards)
			fmt.Println(req.Params.CardID, req.Session.UserID)
			if cerr != nil {
				return nil, cerr
			}

			result, err = r.commentResponses(cards, req.Session.UserID)
			if err != nil {
				return nil, err
			}
		} else {
			// comment level, get all replies flat in asc chron order
			cards, cerr := r.store.GetFlatReplies(req.Params.CardID, req.Session.UserID, false, 0)
			if cerr != nil {
				return nil, cerr
			}

			result, err = r.cardResponses(cards, req.Session.UserID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// Not nested, return the old way
		cards, cerr := r.store.GetThread(req.Params.CardID, req.Session.UserID)
		if cerr != nil {
			return nil, cerr
		}
		result, err = r.cardResponses(cards, req.Session.UserID)
		if err != nil {
			return nil, err
		}
	}

	err = r.store.SetCardVisited(req.Session.UserID, req.Params.CardID)
	if err != nil {
		return nil, err
	}

	return (*GetThreadResponse)(&result), nil
}

// comment responses include latest comment
func (r *rpc) commentResponses(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error) {
	var err error
	result := make([]*model.CardResponse, len(cards))
	for i, card := range cards {
		result[i], err = commentResponse(r.store, card, viewerID)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

//  comment responses include latest comment
func commentResponse(store dataStore, card *model.Card, viewerID globalid.ID) (*model.CardResponse, error) {
	cRsp, err := CardResponse(store, card, viewerID)
	if err != nil {
		return nil, err
	}

	latestComments, err := store.GetFlatReplies(card.ID, viewerID, true, 1)
	if err != nil {
		return nil, err
	}

	if len(latestComments) > 0 {
		cRsp.LatestComment, err = CardResponse(store, latestComments[0], viewerID)
		if err != nil {
			return nil, err
		}
	}

	return cRsp, nil
}

func (r *rpc) cardResponses(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error) {
	var err error
	result := make([]*model.CardResponse, len(cards))
	for i, card := range cards {
		result[i], err = CardResponse(r.store, card, viewerID)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func CardResponse(store dataStore, card *model.Card, viewerID globalid.ID) (*model.CardResponse, error) {
	author, err := authorForCard(store, card)
	if err != nil {
		return nil, err
	}
	var replies int
	replies, err = store.GetThreadCount(card.ID)
	if err != nil {
		return nil, err
	}

	var reaction *model.Reaction
	var voteResponse *model.VoteResponse
	var userReaction *model.UserReaction

	if viewerID != globalid.Nil {
		userReaction, err = store.GetUserReaction(viewerID, card.ID)
		if err == nil && userReaction != nil {
			// if it errors out all 3 reactions fields should be nil
			reaction = userReaction.ToCardReaction()
			voteResponse = userReaction.ToVoteResponse()
		} else if errors.Cause(err) == sql.ErrNoRows {
			userReaction = nil
		}
	}

	isMine := viewerID == card.OwnerID

	threadRoot := card
	if card.ThreadRootID != globalid.Nil {
		threadRoot, err = store.GetCard(card.ThreadRootID)
		if err != nil {
			return nil, err
		}
	}

	var viewer *model.Viewer
	if threadRoot.AuthorToAlias[viewerID] != globalid.Nil {
		var anonymousAlias *model.AnonymousAlias
		viewer = &model.Viewer{}
		anonymousAlias, err = store.GetAnonymousAlias(threadRoot.AuthorToAlias[viewerID])
		if err != nil {
			return nil, err
		}
		var lastUsed bool
		lastUsed, err = store.GetAnonymousAliasLastUsed(viewerID, threadRoot.ID)
		if err != nil {
			return nil, err
		}
		viewer.AnonymousAlias = anonymousAlias
		viewer.AnonymousAliasLastUsed = lastUsed
	}

	engagement, err := store.GetEngagement(card.ID)
	if err != nil {
		return nil, err
	}

	subscribedTypes, err := store.SubscribedToTypes(viewerID, card.ID)
	if err != nil {
		return nil, err
	}

	var channel *model.Channel

	if card.ChannelID != globalid.Nil {
		channel, err = store.GetChannel(card.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	subscribedToUser, err := store.IsFollowing(viewerID, card.OwnerID)
	if err != nil {
		return nil, err
	}
	subscribedToChannel, err := store.GetIsSubscribed(viewerID, card.ChannelID)
	if err != nil {
		return nil, err
	}

	var rankingReason string

	if subscribedToUser && !subscribedToChannel && card.AliasID == globalid.Nil {
		rankingReason = fmt.Sprintf("Because you follow **@%v**", author.Username)
	}

	subscribed := len(subscribedTypes) > 0

	return &model.CardResponse{
		Card:           card.Export(),
		Author:         author,
		Viewer:         viewer,
		Channel:        channel,
		Replies:        replies,
		Reactions:      reaction,
		Engagement:     engagement,
		ViewerReaction: userReaction,
		Score:          0,
		Subscribed:     subscribed,
		IsMine:         isMine,
		Vote:           voteResponse,
		RankingReason:  rankingReason,
	}, nil
}

func featuredComment(store dataStore, card *model.Card, newContentAvailable bool) (*model.FeaturedComment, error) {
	author, err := authorForCard(store, card)
	if err != nil {
		return nil, err
	}
	return &model.FeaturedComment{
		Card:   card.Export(),
		Author: author,
		New:    newContentAvailable,
	}, nil
}

func authorForCard(store dataStore, card *model.Card) (*model.Author, error) {
	if card == nil {
		return nil, nil
	}

	if card.AliasID != globalid.Nil {
		anonymousAlias, err := store.GetAnonymousAlias(card.AliasID)
		if err != nil {
			return nil, err
		}
		return anonymousAlias.Author(), nil
	}

	user, err := store.GetUser(card.OwnerID)
	if err != nil {
		return nil, err
	}
	return user.Author(), nil
}

type ReactToCardRequest struct {
	Params  ReactToCardParams
	Session *model.Session
}

type ReactToCardParams struct {
	// CardID specifies the card to react to.
	CardID globalid.ID `json:"cardID"`
	// Anonymous determines if the reaction should be anonymous.
	Type model.UserReactionType `json:"type"`
	// Anonymous determines if the reaction should be anonymous.
	Anonymous bool `json:"anonymous"`
	// Undo if set to true it will delete the reaction again.
	Undo bool `json:"undo"`

	//OLD PARAMS
	Strength float64 `json:"strength"`
	// Reaction must be any of: boost, bury, seen
	Reaction model.ReactionType `json:"reaction"`
	// something horrible to deal with the bullshit fE does
	FromLegacy bool `json:"fromLegacy"`
}

func (p ReactToCardParams) Sanitize() interface{} {
	return p
}

func (p ReactToCardParams) Validate() error {
	return nil
}

type ReactToCardResponse struct {
	// AnonymousAlias is set when this was an anonymous reaction.
	AnonymousAlias *model.AnonymousAlias `json:"anonymousAlias"`
	NewBalances    *model.CoinBalances   `json:"newBalances"`
}

// ReactToCard creates a reaction by the given user for the specifed card.
// Reactions can be explicit, as exposed by the user interface: boost, bury, or
// implicit like seen.
//
// When anonymous is set to true, an anonymous alias will be generated or
// retrieved from the thread and returned.
//
func (r *rpc) ReactToCard(ctx context.Context, req ReactToCardRequest) (*ReactToCardResponse, error) {
	var response ReactToCardResponse
	var reaction *model.UserReaction
	var err error

	isLegacy := req.Params.Type == "" || req.Params.FromLegacy
	popRankNetUp := int64(0)
	popRankNetDown := int64(0)
	isLike := false
	isDislike := false
	sendPushNotification := false

	// set legacy reactions
	if req.Params.Reaction == model.Boost || req.Params.Reaction == model.Like {
		req.Params.Type = model.ReactionLike
	}

	// Set is like or is dislike
	if req.Params.Type == model.ReactionLike {
		isLike = true
	} else if req.Params.Type == model.ReactionDislike {
		isDislike = true
	}

	// get the card reacted to
	card, cardErr := r.store.GetCard(req.Params.CardID)
	if cardErr != nil {
		return nil, cardErr
	}

	// If this is an undo
	if req.Params.Undo {
		if isLike {
			fmt.Println("likeundo")
		} else {
			fmt.Println("dislikeundo")
		}

		// delete the reaction of the type specified if it exists
		rows, derr := r.store.DeleteUserReactionForType(req.Session.UserID, req.Params.CardID, req.Params.Type)
		if derr != nil {
			return nil, derr
		}

		// set popular rank update strength
		if isLike && rows > 0 {
			popRankNetUp--
		} else if isDislike && rows > 0 {
			popRankNetDown--
		}

		if isLike {
			// unsubscribe from notifs
			err = r.unsubscribeFromNotificationsForType(req.Session.UserID, card.ID, model.CommentType)
			if err != nil {
				return nil, err
			}
		}
	} else if isLike || isDislike {
		// get existing reaction
		reaction, err = r.store.GetUserReaction(req.Session.UserID, req.Params.CardID)

		if err != nil && errors.Cause(err) != sql.ErrNoRows {
			return nil, err
		} else if errors.Cause(err) == sql.ErrNoRows || reaction == nil {
			// create the reaction if it doesn't already exist
			reaction = &model.UserReaction{
				UserID: req.Session.UserID,
				CardID: req.Params.CardID,
				Type:   req.Params.Type,
			}
		}

		// undo the existing reaction's popular rank impact if necessary

		if reaction.Type != req.Params.Type && !isLegacy {
			if reaction.Type == model.ReactionLike {
				popRankNetUp--
			} else {
				popRankNetDown--
			}
		}
		reaction.Type = req.Params.Type

		// get aliases for likes
		if isLike && reaction.AliasID == globalid.Nil && req.Params.Anonymous {
			var alias *model.AnonymousAlias
			alias, err = r.generateAnonymousAlias(card, req.Session.UserID)
			if err != nil {
				return nil, err
			}
			response.AnonymousAlias = alias
			reaction.AliasID = alias.ID
		} else if reaction.AliasID != globalid.Nil && req.Params.Anonymous {
			reaction.AliasID = globalid.Nil
		}

		// save the reaction
		if err = r.store.SaveUserReaction(reaction); err != nil {
			return nil, err
		}

		if isLike {
			if card.OwnerID != req.Session.UserID {
				// update/create a notification for the like for the card owner
				sendPushNotification, err = r.saveNotificationForReaction(ctx, req.Session, card, reaction)
				if err != nil {
					r.log.Error(err)
				}
			}

			// set popular rank update strength
			popRankNetUp++

			// subscribe to notifications
			err = r.subscribeToNotificationsForType(req.Session.UserID, card.ID, model.CommentType)
			if err != nil {
				return nil, err
			}
		} else {
			// set popular rank update strength
			popRankNetDown++
			// unsubscribeFromNotifications
			err = r.unsubscribeFromNotificationsForType(req.Session.UserID, card.ID, model.CommentType)
			if err != nil {
				return nil, err
			}
		}

	} else {
		// unsupported reaction, ignore
		return nil, nil
	}

	// update notifications
	err = r.updateReactionNotification(ctx, req.Session, card, sendPushNotification)
	if err != nil && err != notifExport.ErrNotificationEmpty {
		r.log.Error(err)
	}

	// update popular rank
	if (popRankNetUp != 0.0 || popRankNetDown != 0.0) && card.ThreadRootID == globalid.Nil {
		err = r.store.UpdatePopularRankForCard(card.ID, 0, popRankNetUp, popRankNetDown, 0, 0)
		if err != nil {
			r.log.Error(err)
		}
	}

	if isLike && !req.Params.Undo && req.Session.UserID != card.OwnerID {
		var cB *model.CoinBalances
		cB, err = r.updateTokensForReceiveLike(card.OwnerID, card.ID)
		if err != nil {
			r.log.Error(err)
		}

		err = r.pusher.UpdateCoinBalance(ctx, card.OwnerID, cB)
		if err != nil {
			r.log.Error(err)
		}
	}

	cardID := card.ThreadRootID
	if !card.IsComment() {
		cardID = card.ID
	}

	go r.pushUpdateEngagement(ctx, req.Session, cardID)

	return &response, nil
}

func (r *rpc) saveNotificationForReaction(ctx context.Context, session *model.Session, card *model.Card, reaction *model.UserReaction) (bool, error) {
	if session.UserID == card.OwnerID {
		return false, nil
	}
	notif, err := r.store.LatestForType(card.OwnerID, card.ID, model.BoostType, false)
	newNotif := false

	if reaction != nil && reaction.Type == model.ReactionLike {
		if errors.Cause(err) == sql.ErrNoRows {
			newNotif = true
			notif = &model.Notification{
				ID:            globalid.Next(),
				UserID:        card.OwnerID,
				TargetID:      reaction.CardID,
				TargetAliasID: reaction.AliasID,
				Type:          model.BoostType,
			}
		} else if err != nil {
			return false, err
		} else {
			notif.UpdatedAt = time.Now().UTC()
		}

		err = r.store.SaveNotification(notif)
		if err != nil {
			return false, err
		}

		err = r.store.SaveReactionForNotification(notif.ID, reaction.UserID, reaction.CardID)
		if err != nil {
			return false, err
		}
	} else if reaction != nil && reaction.Type == model.ReactionDislike {
		err := r.store.DeleteReactionForNotification(notif.ID, reaction.UserID, reaction.CardID)
		if err != nil {
			return false, err
		}
		err = r.store.ClearEmptyNotificationsForUser(card.OwnerID)
		if err != nil {
			return false, err
		}
	}

	return newNotif, nil
}

func (r *rpc) updateReactionNotification(ctx context.Context, session *model.Session, card *model.Card, sendPush bool) error {
	notif, err := r.store.LatestForType(card.OwnerID, card.ID, model.BoostType, false)

	if err != nil {
		return err
	} else if notif != nil {
		exNotif, err := r.notifications.ExportNotification(notif)
		if err != nil {
			return err
		}

		if sendPush {
			err = r.notifier.NotifyPush(exNotif)
			if err != nil {
				r.log.Error(err)
			}
		}

		err = r.pusher.NewNotification(ctx, session, exNotif)
		if err != nil {
			r.log.Error(err)
		}
	}

	return nil
}

type VoteOnCardRequest struct {
	Params  VoteOnCardParams
	Session *model.Session
}

type VoteOnCardParams struct {
	// CardID specifies the card to react to.
	CardID globalid.ID `json:"cardID"`
	// Up/down
	Type model.VoteType `json:"type"`
	// Undo if set to true it will delete the reaction again.
	Undo bool `json:"undo"`
}

func (p VoteOnCardParams) Sanitize() interface{} {
	return p
}

func (p VoteOnCardParams) Validate() error {
	return nil
}

type VoteOnCardResponse struct {
	NewBalances *model.CoinBalances `json:"newBalances"`
}

func (r *rpc) VoteOnCard(ctx context.Context, req VoteOnCardRequest) (*VoteOnCardResponse, error) {
	if req.Params.Type == model.Down {
		rtcParams := ReactToCardParams{
			CardID:     req.Params.CardID,
			Type:       model.ReactionDislike,
			Anonymous:  false,
			Undo:       req.Params.Undo,
			FromLegacy: true,
		}

		rctReq := ReactToCardRequest{
			Params:  rtcParams,
			Session: req.Session,
		}

		resp, err := r.ReactToCard(ctx, rctReq)

		if err != nil {
			return nil, err
		}

		return &VoteOnCardResponse{
			NewBalances: resp.NewBalances,
		}, nil
	}

	return nil, nil
}

func (r *rpc) processMentions(ctx context.Context, req PostCardRequest, card *model.Card) (map[globalid.ID]bool, error) {
	// Get mentions via regex
	mentionRE := regexp.MustCompile(`[^\\]?[@!][a-z0-9]+\w`)
	mentions := mentionRE.FindAllString(card.Content, -1)

	// track users we've already processed mentions for to avoid redundant notification
	alreadyProcessed := make(map[globalid.ID]bool)

	for _, v := range mentions {
		username := v[strings.IndexAny(v, "@!"):]
		var notifyUserID globalid.ID
		var notifyAliasID globalid.ID

		// try getting the user by real name
		notifyUser, err := r.store.GetUserByUsername(username[1:])
		if err == nil {
			notifyUserID = notifyUser.ID
		}

		// if the userID is still null, try getting by anon alias
		if notifyUser == nil {
			reply, err := r.store.GetCard(card.ThreadRootID)
			if err != nil {
				return nil, err
			}
			notifyAlias, err := r.store.GetAnonymousAliasByUsername(username[1:])
			if err == nil {
				notifyAliasID = notifyAlias.ID
				for userID, aliasID := range reply.AuthorToAlias {

					if notifyAlias.ID == aliasID {
						notifyUserID = userID
					}
				}
			}
		}

		// save the mention and notify the user if needed
		if notifyUserID != globalid.Nil && notifyUserID != card.OwnerID && !alreadyProcessed[notifyUserID] {
			mention := &model.Mention{
				ID:             globalid.Next(),
				InCard:         card.ID,
				MentionedUser:  notifyUserID,
				MentionedAlias: notifyAliasID,
			}

			err := r.store.SaveMention(mention)
			if err != nil {
				return nil, err
			}

			alreadyProcessed[notifyUserID] = true

			err = r.notifyForMention(ctx, req.Session, mention)
			if err != nil {
				return nil, err
			}
		}
	}
	return alreadyProcessed, nil
}

func (r *rpc) notifyForMention(ctx context.Context, session *model.Session, mention *model.Mention) error {
	// create notification
	newNotif := &model.Notification{
		ID:     globalid.Next(),
		UserID: mention.MentionedUser,
		Type:   model.MentionType,
	}

	err := r.store.SaveNotification(newNotif)
	if err != nil {
		return err
	}

	// attach mention to notification
	notifMention := &model.NotificationMention{
		NotificationID: newNotif.ID,
		MentionID:      mention.ID,
	}

	err = r.store.SaveNotificationMention(notifMention)

	if err != nil {
		return err
	}

	// export
	exNotif, err := r.notifications.ExportNotification(newNotif)

	if err != nil {
		return err
	}

	// Notify via app
	err = r.pusher.NewNotification(ctx, session, exNotif)

	if err != nil {
		return err
	}

	// Notify via push
	return r.notifier.NotifyPush(exNotif)
}

type PostCardRequest struct {
	Params  PostCardParams
	Session *model.Session
}

type PostCardParams struct {
	// AuthorID can be used by admins to overwrite the user ID of the card.
	AuthorID globalid.ID `json:"authorID"`
	// Anonymous is set to true in order to publish a card anonymously.
	Anonymous bool `json:"anonymous"`
	// ReplyCardID is set if this card is a reply to another card.
	ReplyCardID globalid.ID `json:"replyCardID,omitempty"`
	// ReplyCardID is set if this card is a reply to another card.
	ChannelID globalid.ID `json:"channelID,omitempty"`
	// URL associates the card with a link.
	URL string `json:"url"`
	// Content is the card content written in Slate Markdown syntax.
	Content string `json:"content"`
	// BackgroundImage is a Base64 encoded background image to be displayed behind the card.
	BackgroundImage string `json:"bgImage"`
	// BackgroundImageURL points to an image to be displayed behind the card.
	BackgroundImageURL string `json:"bgImageURL"`
	// BackgroundColor specifies a color to be used as a background image.
	BackgroundColor string `json:"bgColor"`
}

func (p PostCardParams) Sanitize() interface{} {
	return p
}

func (p PostCardParams) Validate() error {
	trimmed := strings.TrimSpace(p.Content)
	if len(trimmed) == 0 {
		return errors.New("post has no content")
	}
	return nil
}

type PostCardResponse struct {
	Card        *model.CardView     `json:"card"`
	Author      *model.Author       `json:"author"`
	Channel     *model.Channel      `json:"channel,omitempty"`
	NewBalances *model.CoinBalances `json:"newBalances,omitempty"`
}

// PostCard creates a new card and posts it to the users following the author.
func (r *rpc) PostCard(ctx context.Context, req PostCardRequest) (*PostCardResponse, error) {
	var err error
	var replyCard *model.Card
	// check costs
	hasPostingAlias := false
	if req.Params.ReplyCardID != globalid.Nil {
		replyCard, err = r.store.GetCard(req.Params.ReplyCardID)
		if err != nil {
			return nil, err
		}
		threadRootID := replyCard.ID
		if replyCard.ThreadRootID != globalid.Nil {
			threadRootID = replyCard.ThreadRootID
		}
		hasPostingAlias, err = r.userHasPostingAliasInThread(req.Session.UserID, threadRootID)
		if err != nil {
			return nil, err
		}
	}
	isReply := req.Params.ReplyCardID != globalid.Nil
	isAnonymous := req.Params.Anonymous

	shouldChargeForAnonID := !hasPostingAlias && isAnonymous

	if shouldChargeForAnonID {
		if isReply && !r.userCanAffordThreadAlias(req.Session.UserID) {
			return nil, model.ErrInsufficientBalance
		} else if !r.userCanAffordPostAlias(req.Session.UserID) {
			return nil, model.ErrInsufficientBalance
		}
	}

	usr, err := r.store.GetUser(req.Session.GetUser().ID)
	if err != nil {
		return nil, err
	}

	authorID := req.Session.UserID
	if usr.Admin && req.Params.AuthorID != globalid.Nil {
		authorID = req.Params.AuthorID
	}

	card := &model.Card{
		ID:              globalid.Next(),
		OwnerID:         authorID,
		BackgroundColor: req.Params.BackgroundColor,
		URL:             req.Params.URL,
		AuthorToAlias:   model.IdentityMap{},
		CreatedAt:       time.Now().UTC(),
	}

	reply, err := r.assignReply(ctx, card, authorID, req)
	if err != nil {
		return nil, err
	}

	if req.Params.Anonymous {
		var alias *model.AnonymousAlias
		alias, err = r.generateAnonymousAlias(card, authorID)
		if err != nil {
			return nil, err
		}
		card.AliasID = alias.ID
		card.Author = alias.Author()
	} else {
		var user *model.User
		user, err = r.store.GetUser(authorID)
		if err != nil {
			return nil, err
		}
		card.Author = user.Author()
	}

	parsedContent, err := r.parseAndSaveContentImages(req.Params.Content)
	if err != nil {
		return nil, err
	}

	card.Content = parsedContent

	// shadowbanned users post shadowbanned cards
	if usr.ShadowbannedAt.Valid {
		card.ShadowbannedAt = model.NewDBTime(time.Now().UTC())
	}

	var postToChannel *model.Channel

	if req.Params.ChannelID != globalid.Nil {
		postToChannel, err = r.store.GetChannel(req.Params.ChannelID)
		if err != nil {
			return nil, err
		}

		card.ChannelID = req.Params.ChannelID
	}

	err = r.handleBackgroundImage(card, req.Params.BackgroundColor, req.Params.BackgroundImage, req.Params.BackgroundImageURL)
	if err != nil {
		return nil, err
	}

	err = r.store.SaveCard(card)
	if err != nil {
		return nil, err
	}

	mentioned, merr := r.processMentions(ctx, req, card)
	if merr != nil {
		r.log.Error(merr)
	}

	if card.Reply() && !card.ShadowbannedAt.Valid {
		var subs []globalid.ID
		subs, err = r.store.SubscribersForCard(card.ThreadRootID, model.CommentType)
		if errors.Cause(err) != sql.ErrNoRows && err != nil {
			return nil, err
		}

		// background for now to improve performance
		go func() {
			for _, v := range subs {
				if v != authorID && !mentioned[v] {
					err = r.saveCommentAction(ctx, req.Session, card, v)
					if err != nil {
						r.log.Error(err)
					}
				}
			}
		}()
	}

	err = r.subscribeToNotificationsForType(card.OwnerID, card.ID, model.BoostType)
	if err != nil {
		return nil, err
	}

	subToCard := card.ID

	if card.Reply() {
		subToCard = card.ThreadRootID
	}

	err = r.subscribeToNotificationsForType(card.OwnerID, subToCard, model.CommentType)
	if err != nil {
		return nil, err
	}

	// create a feedEntry for the author
	if card.ThreadRootID == globalid.Nil {
		if err = r.store.AddCardsToTopOfFeed(req.Session.UserID, []globalid.ID{card.ID}); err != nil {
			r.log.Error(err)
		}
	}

	err = r.notifier.NotifySlackAboutCard(card, reply, card.Author)
	if err != nil {
		r.log.Error(err)
	}
	result := &PostCardResponse{
		Card:   card.Export(),
		Author: card.Author,
	}

	if req.Params.ChannelID != globalid.Nil {
		c, chanerr := r.store.GetChannel(req.Params.ChannelID)
		if chanerr != nil {
			return nil, chanerr
		}
		result.Channel = c
	}

	if card.Reply() {
		go r.pushUpdateEngagement(ctx, req.Session, card.ThreadRootID)

		cardResp, berr := CardResponse(r.store, card, req.Session.UserID)
		if berr != nil {
			r.log.Error(err)
		}
		if !card.ShadowbannedAt.Valid {
			go r.pushNewCard(ctx, req.Session, cardResp)
		}
	}

	// if there's no channel (comments) or the posted-to channel is not private
	if req.Params.ChannelID == globalid.Nil || !postToChannel.Private {
		if card.Reply() {
			// update the root's commentrank
			err = r.store.UpdatePopularRankForCard(card.ThreadRootID, 0, 0, 0, 1, 0)
			if err != nil {
				r.log.Error(err)
			}
		} else {
			// create a popular rank for this card
			popRank := &model.PopularRankEntry{
				CardID: card.ID,
			}

			err = r.store.SavePopularRank(popRank)
			if err != nil {
				r.log.Error(err)
			}
		}
	}

	if shouldChargeForAnonID {
		if card.Reply() {
			result.NewBalances, err = r.updateTokensForBuyThreadAlias(card.OwnerID)
			if err != nil {
				r.log.Error(err)
			}
		} else {
			result.NewBalances, err = r.updateTokensForBuyPostAlias(card.OwnerID)
			if err != nil {
				r.log.Error(err)
			}
		}
	}

	// award the reply owner with new cards
	if card.Reply() && replyCard != nil && req.Session.UserID != replyCard.OwnerID {
		var cB *model.CoinBalances
		cB, err = r.updateTokensForReceiveComment(replyCard.OwnerID, replyCard.ID)
		if err != nil {
			r.log.Error(err)
		}

		err = r.pusher.UpdateCoinBalance(ctx, replyCard.OwnerID, cB)
		if err != nil {
			r.log.Error(err)
		}
	}
	return result, nil
}

func (r *rpc) pushUpdateEngagement(ctx context.Context, session *model.Session, cardID globalid.ID) {
	err := r.pusher.UpdateEngagement(ctx, session, cardID)
	if err != nil {
		r.log.Error(err)
	}
}

func (r *rpc) pushNewCard(ctx context.Context, session *model.Session, cardResp *model.CardResponse) {
	err := r.pusher.NewCard(ctx, session, cardResp)
	if err != nil {
		r.log.Error(err)
	}
}

func (r *rpc) saveCommentAction(ctx context.Context, session *model.Session, commentCard *model.Card, notifiedUser globalid.ID) error {
	// get latest notification
	notif, err := r.store.LatestForType(notifiedUser, commentCard.ThreadRootID, model.CommentType, true)

	// if there isn't one, make a new one
	newNotif := false
	if errors.Cause(err) == sql.ErrNoRows {
		newNotif = true
		notif = &model.Notification{
			ID:       globalid.Next(),
			UserID:   notifiedUser,
			TargetID: commentCard.ThreadRootID,
			Type:     model.CommentType,
		}

		err = r.store.SaveNotification(notif)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// create a new notifications_comments
	notificationComment := &model.NotificationComment{
		NotificationID: notif.ID,
		CardID:         commentCard.ID,
	}

	err = r.store.SaveNotificationComment(notificationComment)
	if err != nil {
		return err
	}

	if !newNotif {
		notif.UpdatedAt = time.Now().UTC()
		err = r.store.SaveNotification(notif)

		if err != nil {
			return err
		}
	}

	// export
	exNotif, err := r.notifications.ExportNotification(notif)
	if err != nil {
		return err
	}

	// notify via app notification
	err = r.pusher.NewNotification(ctx, session, exNotif)
	if err != nil {
		r.log.Error(err)
	}

	// notify via push
	if newNotif {
		err = r.notifier.NotifyPush(exNotif)
		if err != nil {
			r.log.Error(err)
		}
	}

	return nil
}

func (r *rpc) assignReply(ctx context.Context, card *model.Card, authorID globalid.ID, req PostCardRequest) (*model.Card, error) {
	if req.Params.ReplyCardID != globalid.Nil {
		reply, err := r.store.GetCard(req.Params.ReplyCardID)
		if err != nil {
			return nil, err
		}

		card.ReplyTo(reply)
		author, err := r.author(reply)
		if err != nil {
			return nil, err
		}
		reply.Author = author
		return reply, nil
	}
	return nil, nil
}

func (r *rpc) author(card *model.Card) (*model.Author, error) {
	if card.AliasID != globalid.Nil {
		var anonymousAlias *model.AnonymousAlias
		anonymousAlias, err := r.store.GetAnonymousAlias(card.AliasID)
		if err != nil {
			return nil, err
		}
		return anonymousAlias.Author(), nil
	} else {
		var user *model.User
		user, err := r.store.GetUser(card.OwnerID)
		if err != nil {
			return nil, err
		}
		return user.Author(), nil
	}
}

func (r *rpc) userHasPostingAliasInThread(userID, threadRootID globalid.ID) (bool, error) {
	card, err := r.store.GetCard(threadRootID)
	if err != nil {
		return false, err
	}

	aliasID := card.AuthorToAlias[userID]

	if aliasID != globalid.Nil {
		postCount, pcerr := r.store.CountPostsByAliasInThread(aliasID, threadRootID)
		if pcerr != nil {
			return false, err
		}

		if postCount > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (r *rpc) generateAnonymousAlias(card *model.Card, authorID globalid.ID) (*model.AnonymousAlias, error) { // bool return is "should charge for alias"
	threadRoot := card

	if card.IsComment() {
		var err error
		threadRoot, err = r.store.GetCard(card.ThreadRootID)
		if err != nil {
			return nil, err
		}
	}

	var alias *model.AnonymousAlias
	aliasID := threadRoot.AuthorToAlias[authorID]
	if aliasID != globalid.Nil {
		var err error
		alias, err = r.store.GetAnonymousAlias(aliasID)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		alias, err = r.store.GetUnusedAlias(threadRoot.ID)
		if err != nil {
			return nil, err
		}
		threadRoot.AuthorToAlias[authorID] = alias.ID
		err = r.store.SaveCard(threadRoot)
		if err != nil {
			return nil, err
		}

		err = r.store.AssignAliasForUserTipsInThread(authorID, threadRoot.ID, alias.ID)
		if err != nil {
			return nil, err
		}
	}
	return alias, nil
}

func (r *rpc) handleBackgroundImage(card *model.Card, backgroundColor, backgroundImage, backgroundImageURL string) error {
	if backgroundColor != "" && backgroundImage == "" && backgroundImageURL == "" {
		url, _, err := r.imageProcessor.GradientImage(backgroundColor)
		if err != nil {
			return err
		}
		card.BackgroundImagePath = url
		return nil
	}
	if backgroundImage != "" {
		url, filename, err := r.imageProcessor.SaveBase64CardImage(backgroundImage)
		if err != nil {
			return err
		}
		if backgroundColor != "" {
			url, err = r.imageProcessor.BlendImage(filename, backgroundColor)
			if err != nil {
				return err
			}
		}
		card.BackgroundImagePath = url
	}
	if backgroundImageURL != "" {
		url, filename, err := r.imageProcessor.DownloadCardImage(backgroundImageURL)
		if err != nil {
			return err
		}
		if backgroundColor != "" {
			url, err = r.imageProcessor.BlendImage(filename, backgroundColor)
			if err != nil {
				return err
			}
		}
		card.BackgroundImagePath = url
	}
	return nil
}

func (r *rpc) parseAndSaveContentImages(postMarkdown string) (string, error) {
	md := postMarkdown

	// get all the image tag URLs
	// matches image tags in the form of `!(text)[url]`
	findRE := regexp.MustCompile(`\!\[[^\(\)]*\]\([^\[\]]*\)`)

	// matches text between brackets
	extractRE := regexp.MustCompile(`\([^\(\)]*\)`)

	matchedTags := findRE.FindAllString(md, -1)

	urls := []string{}

	for _, tag := range matchedTags {
		urls = append(urls, strings.Trim(extractRE.FindString(tag), "()"))
	}

	linkbarRE := regexp.MustCompile(`%%%[^%]*%%%`)

	matchedBars := linkbarRE.FindAllString(md, -1)

	for _, bar := range matchedBars {
		tokens := strings.Split(bar, "\n")

		if len(tokens) > 2 && strings.HasPrefix(tokens[2], "http") {
			urls = append(urls, tokens[2])
		}
	}

	for _, url := range urls {
		// save the image
		newURL, _, err := r.imageProcessor.DownloadCardImage(url)
		if err != nil {
			return "", err
		}

		// replace the new URL in the post's content
		md = strings.Replace(md, url, newURL, 1)
	}

	return md, nil
}

type NewInviteRequest struct {
	Params  NewInviteParams
	Session *model.Session
}

type NewInviteParams struct {
	// Invites specifies the number of invites the generated token will be
	// valid for.
	Invites int `json:"invites"`
}

func (p NewInviteParams) Sanitize() interface{} {
	return p
}

func (p NewInviteParams) Validate() error {
	return nil
}

type NewInviteResponse model.Invite

// NewInvite creates a new invite to be given out for users to sign up. An
// invite is bound to the user creating it.
func (r *rpc) NewInvite(ctx context.Context, req NewInviteRequest) (*NewInviteResponse, error) {
	invite, err := model.NewInvite(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	if req.Params.Invites > 0 {
		invite.RemainingUses = req.Params.Invites
	}
	err = r.store.SaveInvite(invite)
	if err != nil {
		return nil, err
	}
	return (*NewInviteResponse)(invite), nil
}

type RegisterDeviceRequest struct {
	Params  RegisterDeviceParams
	Session *model.Session
}

type RegisterDeviceParams struct {
	// Token is generated by Firebase and identifies the device.
	Token string `json:"token"`
	// Platform to identify the mobile operating system, either iOS or Android.
	Platform string `json:"platform"`
}

func (p RegisterDeviceParams) Sanitize() interface{} {
	return p
}

func (p RegisterDeviceParams) Validate() error {
	return nil
}

type RegisterDeviceResponse struct{}

// RegisterDevices registers a mobile device to the current user so that it
// will receive push notifications.
func (r *rpc) RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*RegisterDeviceResponse, error) {
	user, err := r.store.GetUser(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	if user.Devices == nil {
		user.Devices = make(map[string]model.Device)
	}

	device := model.NewDevice(req.Params.Token, req.Params.Platform)
	if _, ok := user.Devices[req.Params.Token]; !ok {
		user.Devices[req.Params.Token] = *device
		return nil, r.store.SaveUser(user)
	}
	return nil, nil
}

type UnregisterDeviceRequest struct {
	Params  UnregisterDeviceParams
	Session *model.Session
}

type UnregisterDeviceParams struct {
	// Token is generated by Firebase and identifies the device.
	Token string `json:"token"`
}

func (p UnregisterDeviceParams) Sanitize() interface{} {
	return p
}

func (p UnregisterDeviceParams) Validate() error {
	return nil
}

type UnregisterDeviceResponse struct{}

// UnregisterDevices removes a mobile device from the user's list of devices
// again so that it won't receive any further push notifications.
func (r *rpc) UnregisterDevice(ctx context.Context, req UnregisterDeviceRequest) (*UnregisterDeviceResponse, error) {
	user, err := r.store.GetUser(req.Session.UserID)
	if err != nil {
		return nil, err
	}
	if _, ok := user.Devices[req.Params.Token]; ok {
		delete(user.Devices, req.Params.Token)
		return nil, r.store.SaveUser(user)
	}
	return nil, nil
}

type UpdateSettingsRequest struct {
	Params  UpdateSettingsParams
	Session *model.Session
}

type UpdateSettingsParams struct {
	// Username for changing the user's username.
	Username *string `json:"username"`
	// Password for changing the user's password.
	Password *string `json:"password"`
	// Email for changing the user's email.
	Email *string `json:"email"`
	// AllowEmailNotifications specifies whether the user wants to receive
	// email notifications.
	AllowEmailNotifications *bool `json:"allowEmail"`
	// Bio is a short description of the user.
	Bio *string `json:"bio"`
	// Display is how the user is displayed everywhere.
	DisplayName *string `json:"displayName"`
	// FirstName for changing the user's first name.
	FirstName *string `json:"firstName"`
	// LastName for changing the user's last name.
	LastName *string `json:"lastName"`
	// RemoveProfileImage is set to true if the profile image should be deleted
	// altogether.
	RemoveProfileImage *bool `json:"removeProfileImage"`
	// RemoveCoverImage is set to true if the cover image should be deleted
	// altogether.
	RemoveCoverImage *bool `json:"removeCoverImage"`
	// ImageData updates the user's profile picture by providing a Base64 encoded image.
	ImageData *string `json:"imageData"`
	// CoverImageData updates the user's cover picture by providing a Base64
	// encoded image.
	CoverImageData *string `json:"coverImageData"`
}

func (p UpdateSettingsParams) Sanitize() interface{} {
	filtered := filtered
	if p.Password != nil {
		p.Password = &filtered
	}
	if p.ImageData != nil {
		p.ImageData = &filtered
	}
	if p.CoverImageData != nil {
		p.CoverImageData = &filtered
	}
	return p
}

func (p UpdateSettingsParams) Validate() error {
	return nil
}

type UpdateSettingsResponse model.ExportedUser

// UpdateSettings updates various details about a user. Each field is optional.
func (r *rpc) UpdateSettings(ctx context.Context, req UpdateSettingsRequest) (*UpdateSettingsResponse, error) {
	user, err := r.store.GetUser(req.Session.UserID)
	if err != nil {
		return nil, err
	}
	if req.Params.Username != nil {
		usernameLowerCase := strings.ToLower(*req.Params.Username)
		var blacklist map[string]bool
		blacklist, err = r.blacklistedUsernames(ctx)
		if err != nil {
			return nil, err
		}
		err = model.ValidateUsername(usernameLowerCase, blacklist)
		if err != nil {
			return nil, err
		}
		user.Username = usernameLowerCase
	}
	if req.Params.Password != nil {
		err = user.SetPassword(*req.Params.Password)
		if err != nil {
			return nil, err
		}
	}
	if req.Params.Email != nil {
		user.Email = *req.Params.Email
	}
	if req.Params.AllowEmailNotifications != nil {
		user.AllowEmail = *req.Params.AllowEmailNotifications
	}
	if req.Params.DisplayName != nil {
		user.DisplayName = *req.Params.DisplayName
	}
	if req.Params.Bio != nil {
		user.Bio = *req.Params.Bio
	}
	if req.Params.FirstName != nil {
		user.FirstName = *req.Params.FirstName
	}
	if req.Params.LastName != nil {
		user.LastName = *req.Params.LastName
	}

	user.DisplayName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	if req.Params.ImageData != nil {
		var url string
		url, _, err = r.imageProcessor.SaveBase64ProfileImage(*req.Params.ImageData)
		if err != nil {
			return nil, err
		}
		user.ProfileImagePath = url
	} else if req.Params.RemoveProfileImage != nil && *req.Params.RemoveProfileImage {
		var url string
		url, _, err = r.imageProcessor.GenerateDefaultProfileImage()
		if err != nil {
			return nil, err
		}
		user.ProfileImagePath = url
	}

	if req.Params.CoverImageData != nil {
		var url string
		url, _, err = r.imageProcessor.SaveBase64CoverImage(*req.Params.CoverImageData)
		if err != nil {
			return nil, err
		}
		user.CoverImagePath = url
	} else if req.Params.RemoveCoverImage != nil && *req.Params.RemoveCoverImage {
		user.ProfileImagePath = ""
	}

	err = r.store.SaveUser(user)
	if err != nil {
		return nil, err
	}
	result := user.Export(user.ID)
	extResult := user.Export(globalid.Nil)

	err = r.pusher.UpdateUser(ctx, req.Session, extResult)
	if err != nil {
		r.log.Error(err)
	}

	err = r.indexer.IndexUser(user)
	if err != nil {
		r.log.Error(err)
	}

	return (*UpdateSettingsResponse)(result), nil
}

type NewUserRequest struct {
	Params  NewUserParams
	Session *model.Session
}

type NewUserParams struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	DisplayName        string `json:"displayName"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Token              string `json:"token"`
	ProfilePicturePath string `json:"profilePicturePath"`
	CoverPicturePath   string `json:"coverPicturePath"`
	Firstname          string `json:"firstname"`
	Lastname           string `json:"lastname"`
}

func (p NewUserParams) Sanitize() interface{} {
	p.Password = filtered
	return p
}

func (p NewUserParams) Validate() error {
	return nil
}

type NewUserResponse struct {
	ID globalid.ID `json:"id"`
}

func (r *rpc) NewUser(ctx context.Context, req NewUserRequest) (*NewUserResponse, error) {
	user := model.NewUser(globalid.Next(), req.Params.Username, req.Params.Email, req.Params.DisplayName)
	user.FirstName = req.Params.FirstName
	user.LastName = req.Params.LastName
	user.CoinBalance = 20000
	err := user.SetPassword(req.Params.Password)
	if err != nil {
		return nil, err
	}

	if req.Params.ProfilePicturePath != "" {
		var url string
		url, _, err = r.imageProcessor.DownloadProfileImage(req.Params.ProfilePicturePath)
		if err != nil {
			return nil, err
		}
		user.ProfileImagePath = url
	}

	if req.Params.CoverPicturePath != "" {
		var url string
		url, _, err = r.imageProcessor.DownloadProfileImage(req.Params.CoverPicturePath)
		if err != nil {
			return nil, err
		}
		user.CoverImagePath = url
	}

	err = r.store.SaveUser(user)
	if err != nil {
		return nil, mapErrors(err)
	}

	return &NewUserResponse{ID: user.ID}, nil
}

type ValidateInviteCodeRequest struct {
	Params  ValidateInviteCodeParams
	Session *model.Session
}

type ValidateInviteCodeParams struct {
	// Token is the invite code.
	Token string `json:"token"`
}

func (p ValidateInviteCodeParams) Sanitize() interface{} {
	return p
}

func (p ValidateInviteCodeParams) Validate() error {
	return nil
}

type ValidateInviteCodeResponse struct{}

// ValidateInviteCode checks whether a given invite code can be used for
// signing up a new user. Invite codes are limited in the number of times they
// can be used to sign up a new user. This endpoint checks whether the code
// exists and if there are any uses left.
func (r *rpc) ValidateInviteCode(ctx context.Context, req ValidateInviteCodeRequest) (*ValidateInviteCodeResponse, error) {
	_, err := r.store.GetInviteByToken(strings.ToUpper(req.Params.Token))
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		return nil, model.ErrInvalidInviteCode
	} else if err != nil {
		return nil, err
	}
	return nil, nil
}

type AddToWaitlistRequest struct {
	Params  AddToWaitlistParams
	Session *model.Session
}

type AddToWaitlistParams struct {
	// Email identifies the user entering the waitlist.
	Email string `json:"email,omitempty"`
	// Name adds a human description to the waitlist entry, ideally first and
	// last name.
	Name string `json:"name,omitempty"`
	// AccessToken uses Facebook account details to sign up to the waitlist.
	AccessToken string `json:"accessToken,omitempty"`
}

func (p AddToWaitlistParams) Sanitize() interface{} {
	return p
}

func (p AddToWaitlistParams) Validate() error {
	if p.Email != "" && p.AccessToken != "" {
		return errors.New("provide either email or accessToken but not both")
	}
	if p.Email == "" && p.AccessToken == "" {
		return errors.New("provide either email or accessToken")
	}
	return nil
}

type AddToWaitlistResponse struct {
	// AccessToken is the extended token of the access token used in the
	// request.
	AccessToken string `json:"accessToken"`
	// AccessTokenExpiresAt is the number of seconds until this access token
	// expires.
	AccessTokenExpiresAt int64 `json:"accessTokenExpiresAt"`
}

// AddToWaitlist lets users without an invite code for sign up to enter a
// waitlist instead. Sign up to the waitlist can either happen via email and
// name or via Facebook access token.
//
// If waitlist sign up happens via access token, an extended access token plus
// the time when it expires is returned.
func (r *rpc) AddToWaitlist(ctx context.Context, req AddToWaitlistRequest) (*AddToWaitlistResponse, error) {
	var response AddToWaitlistResponse
	entry := &model.WaitlistEntry{Email: req.Params.Email, Name: req.Params.Name}
	if req.Params.AccessToken != "" {
		accessToken, err := r.oauth2.ExtendToken(ctx, req.Params.AccessToken)
		if err != nil {
			return nil, err
		}
		user, err := accessToken.FacebookUser()
		if err != nil {
			return nil, err
		}
		entry = &model.WaitlistEntry{
			Email: user.Email,
			Name:  user.Name(),
		}
		response.AccessToken = accessToken.Token()
		response.AccessTokenExpiresAt = accessToken.ExpiresAt()
	}
	err := r.store.SaveWaitlistEntry(entry)
	if err != nil {
		return nil, err
	}
	name := req.Params.Name
	if name == "" {
		name = req.Params.Email
	}

	err = r.notifier.NotifySlack("waitlist", fmt.Sprintf("%s signed up for the waitlist.", name))
	if err != nil {
		r.log.Error(err)
	}
	return &response, nil
}

type GetUsersRequest struct {
	Params  GetUsersParams
	Session *model.Session
}

type GetUsersParams struct{}

func (p GetUsersParams) Sanitize() interface{} {
	return p
}

func (p GetUsersParams) Validate() error {
	return nil
}

type GetUsersResponse []*model.User

func (r *rpc) GetUsers(ctx context.Context, req GetUsersRequest) (*GetUsersResponse, error) {
	users, err := r.store.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		user.ProfileImagePath = fmt.Sprintf("%v/%v", r.config.ProfileImagesPath, user.ProfileImagePath)
	}
	return (*GetUsersResponse)(&users), nil
}

type GetUserRequest struct {
	Params  GetUserParams
	Session *model.Session
}

type GetUserParams struct {
	Username string      `json:"username"`
	UserID   globalid.ID `json:"userID,omitempty"`
}

func (p GetUserParams) Sanitize() interface{} {
	return p
}

func (p GetUserParams) Validate() error {
	return nil
}

type GetUserResponse model.ExportedUser

func (r *rpc) GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	var user *model.User
	var err error
	if req.Params.Username != "" {
		user, err = r.store.GetUserByUsername(req.Params.Username)
	} else {
		user, err = r.store.GetUser(req.Params.UserID)
	}
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) || (user.Username == store.RootUser.Username) {
		return nil, datastore.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return (*GetUserResponse)(user.Export(req.Session.UserID)), nil
}

type ValidateUsernameRequest struct {
	Params  ValidateUsernameParams
	Session *model.Session
}

type ValidateUsernameParams struct {
	// Username is the username to be checked.
	Username string `json:"username"`
}

func (p ValidateUsernameParams) Sanitize() interface{} {
	return p
}

func (p ValidateUsernameParams) Validate() error {
	return nil
}

type ValidateUsernameResponse struct{}

// ValidateUsername checks whether a given username can be used for a new user
// or to change the username of an existing user.
func (r *rpc) ValidateUsername(ctx context.Context, req ValidateUsernameRequest) (*ValidateUsernameResponse, error) {
	usernameLowerCase := strings.ToLower(req.Params.Username)
	blacklist, err := r.blacklistedUsernames(ctx)
	if err != nil {
		return nil, err
	}
	err = model.ValidateUsername(req.Params.Username, blacklist)
	if err != nil {
		return nil, err
	}
	// check to see if someone already has this username
	u, err := r.store.GetUserByUsername(usernameLowerCase)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, err
	}
	if u != nil {
		return nil, ErrUsernameTaken
	}
	return nil, nil
}

func (r *rpc) blacklistedUsernames(ctx context.Context) (map[string]bool, error) {
	aliases, err := r.store.GetAnonymousAliases()
	if err != nil {
		return nil, err
	}
	blacklist := make(map[string]bool)
	for _, alias := range aliases {
		blacklist[alias.Username] = true
	}
	return blacklist, nil
}

type GetNotificationsRequest struct {
	Params  GetNotificationsParams
	Session *model.Session
}

type GetNotificationsParams struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}

func (p GetNotificationsParams) Sanitize() interface{} {
	return p
}

func (p GetNotificationsParams) Validate() error {
	return nil
}

type GetNotificationsResponse struct {
	Notifications []*model.ExportedNotification `json:"notifications"`
	NextPage      bool                          `json:"nextPage"`
	UnseenCount   int                           `json:"unseenCount"`
}

func (r *rpc) GetNotifications(ctx context.Context, req GetNotificationsRequest) (*GetNotificationsResponse, error) {
	err := r.store.ClearEmptyNotificationsForUser(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	notifications, err := r.store.GetNotifications(req.Session.UserID, req.Params.PageSize, req.Params.PageNumber)
	if err != nil {
		return nil, err
	}

	exportedNotifs := []*model.ExportedNotification{}

	for _, n := range notifications {
		notif, nerr := r.notifications.ExportNotification(n)
		if nerr != nil {
			r.log.Error(nerr)
		} else {
			exportedNotifs = append(exportedNotifs, notif)
		}
	}

	nextNotifs, err := r.store.GetNotifications(req.Session.UserID, req.Params.PageSize, req.Params.PageNumber+1)
	if err != nil {
		return nil, err
	}
	unseenCount, err := r.store.UnseenNotificationsCount(req.Session.UserID)
	if err != nil {
		return nil, nil
	}
	return &GetNotificationsResponse{
		Notifications: exportedNotifs,
		NextPage:      len(nextNotifs) > 0,
		UnseenCount:   unseenCount,
	}, nil
}

type UpdateNotificationsRequest struct {
	Params  UpdateNotificationsParams
	Session *model.Session
}

type UpdateNotificationsParams struct {
	IDs    []globalid.ID `json:"ids"`
	Seen   bool          `json:"seen"`
	Opened bool          `json:"opened"`
}

func (p UpdateNotificationsParams) Sanitize() interface{} {
	return p
}

func (p UpdateNotificationsParams) Validate() error {
	return nil
}

type UpdateNotificationsResponse struct{}

func (r *rpc) UpdateNotifications(ctx context.Context, req UpdateNotificationsRequest) (*UpdateNotificationsResponse, error) {
	if req.Params.Seen {
		err := r.store.UpdateAllNotificationsSeen(req.Session.UserID)
		if err != nil {
			return nil, err
		}
	}
	if req.Params.Opened {
		err := r.store.UpdateNotificationsOpened(req.Params.IDs)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

type GetAnonymousHandleRequest struct {
	Params  GetAnonymousHandleParams
	Session *model.Session
}

type GetAnonymousHandleParams struct {
	ThreadRootID globalid.ID `json:"forThread"`
}

func (p GetAnonymousHandleParams) Sanitize() interface{} {
	return p
}

func (p GetAnonymousHandleParams) Validate() error {
	return nil
}

type GetAnonymousHandleResponse struct {
	Alias    *model.AnonymousAlias `json:"alias"`
	LastUsed bool                  `json:"wasLastUsed"`
}

func (r *rpc) GetAnonymousHandle(ctx context.Context, req GetAnonymousHandleRequest) (*GetAnonymousHandleResponse, error) {
	if req.Params.ThreadRootID != globalid.Nil {
		threadRoot, err := r.store.GetCard(req.Params.ThreadRootID)
		if err != nil {
			return nil, err
		}
		existingID, ok := threadRoot.AuthorToAlias[req.Session.UserID]
		if ok {
			alias, aliasErr := r.store.GetAnonymousAlias(existingID)
			if aliasErr != nil {
				return nil, aliasErr
			}
			wasAnon, wasAnonErr := r.store.GetAnonymousAliasLastUsed(req.Session.UserID, req.Params.ThreadRootID)
			return &GetAnonymousHandleResponse{
				Alias:    alias,
				LastUsed: wasAnon,
			}, wasAnonErr
		}
	}
	alias, err := r.store.GetUnusedAlias(req.Params.ThreadRootID)
	return &GetAnonymousHandleResponse{
		Alias:    alias,
		LastUsed: false,
	}, err
}

type GetTagsRequest struct {
	Params  GetTagsParams
	Session *model.Session
}

type GetTagsParams struct{}

func (p GetTagsParams) Sanitize() interface{} {
	return p
}

func (p GetTagsParams) Validate() error {
	return nil
}

type GetTagsResponse struct{}

func (r *rpc) GetTags(ctx context.Context, req GetTagsRequest) (*GetTagsResponse, error) {
	// deprecated but maybe still called
	return nil, nil
}

type GetFeaturesForUserRequest struct {
	Params  GetFeaturesForUserParams
	Session *model.Session
}

type GetFeaturesForUserParams struct {
}

func (p GetFeaturesForUserParams) Sanitize() interface{} {
	return p
}

func (p GetFeaturesForUserParams) Validate() error {
	return nil
}

type GetFeaturesForUserResponse []string

func (r *rpc) GetFeaturesForUser(ctx context.Context, req GetFeaturesForUserRequest) (*GetFeaturesForUserResponse, error) {
	onSwitches, err := r.store.GetOnSwitches()

	if err != nil {
		return nil, err
	}

	features := make([]string, 0)

	for _, sw := range onSwitches {
		// these are strings SO CASE MATTERS FUCK
		upper, _ := globalid.Parse(strings.ToUpper(string(req.Session.UserID)))
		if sw.State == "on" || sw.TestingUsers[upper] {
			features = append(features, sw.Name)
		}
	}

	return (*GetFeaturesForUserResponse)(&features), nil
}

type PreviewContentRequest struct {
	Params  PreviewContentParams
	Session *model.Session
}

type PreviewContentParams struct {
	URL string `json:"url"`
}

func (p PreviewContentParams) Sanitize() interface{} {
	return p
}

func (p PreviewContentParams) Validate() error {
	return nil
}

type PreviewContentResponse struct {
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Type      string   `json:"type"`
	ImageURLs []string `json:"imageURLs"`
}

func (r *rpc) PreviewContent(ctx context.Context, req PreviewContentRequest) (*PreviewContentResponse, error) {
	result, err := r.extractWithEmbedly(ctx, req)
	if err != nil {
		r.log.Info("Embedly failed, falling back to Diffbot", "err", err)
		return r.extractWithDiffbot(ctx, req)
	}
	return result, nil
}

func (r *rpc) extractWithEmbedly(ctx context.Context, req PreviewContentRequest) (*PreviewContentResponse, error) {
	var result PreviewContentResponse
	endpoint := fmt.Sprintf("https://api.embedly.com/1/extract?key=%s&url=%s&title=og&description=og&meta_images=1", r.config.EmbedlyToken, req.Params.URL)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type embedlyResponse struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Images      []struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"images"`
		ErrorMessage string `json:"error_message,omitempty"`
		ErrorCode    int    `json:"error_code"`
	}

	var response embedlyResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		return nil, err
	}
	if response.ErrorCode != 0 {
		return nil, errors.New(response.ErrorMessage)
	}

	result.Title = response.Title
	result.Summary = firstWords(response.Description, 200)
	result.Type = response.Type
	for _, image := range response.Images {
		if image.Height >= 280 && image.Width >= 360 {
			result.ImageURLs = append(result.ImageURLs, image.URL)
		}
	}

	return &result, nil
}

func (r *rpc) extractWithDiffbot(ctx context.Context, req PreviewContentRequest) (*PreviewContentResponse, error) {
	if r.config.DiffbotToken == "" {
		return nil, errors.New("diffbot token not specified")
	}
	var result PreviewContentResponse

	endpoint := fmt.Sprintf("https://api.diffbot.com/v3/article?token=%s&url=%s", r.config.DiffbotToken, req.Params.URL)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type diffbotResponse struct {
		Objects []struct {
			Title  string `json:"title"`
			Text   string `json:"text"`
			Images []struct {
				URL string `json:"url"`
			}
		} `json:"objects"`
		Error string `json:"error"`
	}
	var response diffbotResponse

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	if len(response.Objects) == 0 {
		return nil, errors.New("no preview available")
	}
	metadata := response.Objects[0]
	result.Title = metadata.Title
	result.Summary = firstWords(metadata.Text, 200)
	result.Summary = strings.Replace(result.Summary, "\n", " ", -1)
	for _, image := range metadata.Images {
		result.ImageURLs = append(result.ImageURLs, image.URL)
	}
	return &result, nil
}

// FirstWords retrieves count words from value detecting words through
// separation of spaces.
func firstWords(value string, count int) string {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count--
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i]
			}
		}
	}
	// Return the entire string.
	return value
}

type UploadImageRequest struct {
	Params  UploadImageParams
	Session *model.Session
}

type UploadImageParams struct {
	CardImageData        string `json:"cardImageData"`
	CardImageColor       string `json:"cardImageColor"`
	CardContentImageData string `json:"cardContentImageData"`
	ProfileImageData     string `json:"profileImageData"`
	CoverImageData       string `json:"coverImageData"`
}

func (p UploadImageParams) Sanitize() interface{} {
	if p.CardImageData != "" {
		p.CardImageData = filtered
	}
	if p.CardContentImageData != "" {
		p.CardContentImageData = filtered
	}
	if p.ProfileImageData != "" {
		p.ProfileImageData = filtered
	}
	if p.CoverImageData != "" {
		p.CoverImageData = filtered
	}
	return p
}

func (p UploadImageParams) Validate() error {
	return nil
}

type UploadImageResponse string

func (r *rpc) UploadImage(ctx context.Context, req UploadImageRequest) (*UploadImageResponse, error) {
	if req.Params.CardImageData != "" {
		url, filepath, err := r.imageProcessor.SaveBase64CardImage(req.Params.CardImageData)
		if err != nil {
			return nil, err
		}
		if req.Params.CardImageColor != "" {
			url, err = r.imageProcessor.BlendImage(filepath, req.Params.CardImageColor)
			if err != nil {
				return nil, err
			}
		}
		return (*UploadImageResponse)(&url), nil
	}
	if req.Params.CardContentImageData != "" {
		url, _, err := r.imageProcessor.SaveBase64CardContentImage(req.Params.CardContentImageData)
		if err != nil {
			return nil, err
		}
		return (*UploadImageResponse)(&url), nil
	}
	if req.Params.ProfileImageData != "" {
		url, _, err := r.imageProcessor.SaveBase64ProfileImage(req.Params.ProfileImageData)
		if err != nil {
			return nil, err
		}
		return (*UploadImageResponse)(&url), nil
	}
	if req.Params.CoverImageData != "" {
		url, _, err := r.imageProcessor.SaveBase64CoverImage(req.Params.CoverImageData)
		if err != nil {
			return nil, err
		}
		return (*UploadImageResponse)(&url), nil
	}
	return nil, ErrNoImageData
}

type ConnectUsersRequest struct {
	Params  ConnectUsersParams
	Session *model.Session
}

type ConnectUsersParams struct {
	Users []string `json:"users"`
}

func (p ConnectUsersParams) Sanitize() interface{} {
	return p
}

func (p ConnectUsersParams) Validate() error {
	return nil
}

type ConnectUsersResponse struct{}

func (r *rpc) ConnectUsers(ctx context.Context, req ConnectUsersRequest) (*ConnectUsersResponse, error) {
	users, err := r.store.GetUsersByUsernames(req.Params.Users)
	if err != nil {
		return nil, err
	}
	for _, u1 := range users {
		for _, u2 := range users {
			if u1.ID != u2.ID {
				err = r.store.SaveFollower(u1.ID, u2.ID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return nil, nil
}

type ModifyCardScoreRequest struct {
	Params  ModifyCardScoreParams
	Session *model.Session
}

type ModifyCardScoreParams struct {
	CardID   globalid.ID `json:"cardID"`
	Strength float64     `json:"strength"`
}

func (p ModifyCardScoreParams) Sanitize() interface{} {
	return p
}

func (p ModifyCardScoreParams) Validate() error {
	return nil
}

type ModifyCardScoreResponse struct{}

func (r *rpc) ModifyCardScore(ctx context.Context, req ModifyCardScoreRequest) (*ModifyCardScoreResponse, error) {
	err := r.store.SaveScoreModification(&model.ScoreModification{
		CardID:   req.Params.CardID,
		UserID:   req.Session.UserID,
		Strength: req.Params.Strength,
	})
	if err != nil {
		return nil, err
	}

	return nil, r.store.UpdatePopularRankForCard(req.Params.CardID, 0, 0, 0, 0, req.Params.Strength)
}

type GetInvitesRequest struct {
	Params  GetInvitesParams
	Session *model.Session
}

type GetInvitesParams struct{}

func (p GetInvitesParams) Sanitize() interface{} {
	return p
}

func (p GetInvitesParams) Validate() error {
	return nil
}

type GetInvitesResponse struct {
	Invites []*model.Invite `json:"invites"`
}

func (r *rpc) GetInvites(ctx context.Context, req GetInvitesRequest) (*GetInvitesResponse, error) {
	invites, err := r.store.GetInvitesForUser(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	return &GetInvitesResponse{
		Invites: invites,
	}, nil
}

type GetOnboardingDataRequest struct {
	Params  GetOnboardingDataParams
	Session *model.Session
}

type GetOnboardingDataParams struct{}

func (p GetOnboardingDataParams) Sanitize() interface{} {
	return p
}

func (p GetOnboardingDataParams) Validate() error {
	return nil
}

type GetOnboardingDataResponse struct {
	NetworkProfilePictures []string            `json:"networkProfilePictures"`
	InvitingUser           *model.ExportedUser `json:"invitingUser,omitempty"`
}

func (r *rpc) GetOnboardingData(ctx context.Context, req GetOnboardingDataRequest) (*GetOnboardingDataResponse, error) {
	users, err := r.store.GetUsers()
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(users))
	for i, v := range users {
		ret[i] = v.ProfileImagePath
	}

	invite, err := r.store.GetInvite(req.Session.User.JoinedFromInvite)
	var exportedInvitingUser *model.ExportedUser
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		// this erroring out is fine, just means nobody invited you
		r.log.Info("user with no inviter is onboarding")
	} else {
		invitingUser, err := r.store.GetUser(invite.NodeID)
		if err != nil {
			return nil, err
		}
		exportedInvitingUser = invitingUser.Export(req.Session.UserID)
	}

	return &GetOnboardingDataResponse{
		NetworkProfilePictures: ret,
		InvitingUser:           exportedInvitingUser,
	}, nil
}

type GetMyNetworkRequest struct {
	Params  GetMyNetworkParams
	Session *model.Session
}

type GetMyNetworkParams struct {
	PageSize   int    `json:"pageSize"`
	PageNumber int    `json:"pageNumber"`
	Search     string `json:"search"`
}

func (p GetMyNetworkParams) Sanitize() interface{} {
	return p
}

func (p GetMyNetworkParams) Validate() error {
	return nil
}

type MyNetworkUser struct {
	User      *model.ExportedUser `json:"user"`
	Rank      int                 `json:"rank"`
	MovedUp   bool                `json:"movedUp"`
	Following bool                `json:"following"`
}

type GetMyNetworkResponse struct {
	Users       []MyNetworkUser `json:"users"`
	LastUpdated int64           `json:"lastUpdated"`
	NextPage    bool            `json:"nextPage"`
}

func (r *rpc) GetMyNetwork(ctx context.Context, req GetMyNetworkRequest) (*GetMyNetworkResponse, error) {
	currentRankings, err := r.store.GetSafeUsersByPage(req.Session.UserID, req.Params.PageSize, req.Params.PageNumber, req.Params.Search)
	if err != nil {
		return nil, err
	}

	nextPage, err := r.store.GetSafeUsersByPage(req.Session.UserID, req.Params.PageSize, req.Params.PageNumber+1, req.Params.Search)
	if err != nil {
		return nil, err
	}

	followingUsers, err := r.store.GetFollowing(req.Session.UserID)
	if err != nil {
		r.log.Error(err)
	}

	following := make(map[globalid.ID]bool, len(followingUsers))
	for _, followee := range followingUsers {
		following[followee.ID] = true
	}

	exportedUsers := make([]MyNetworkUser, len(currentRankings))
	for i, user := range currentRankings {
		exportedUsers[i] = MyNetworkUser{
			User:      user.Export(req.Session.UserID),
			Rank:      0,
			MovedUp:   false,
			Following: following[user.ID],
		}
	}

	year, month, day := time.Now().Date()
	thisMorning := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Unix()

	return &GetMyNetworkResponse{
		Users:       exportedUsers,
		LastUpdated: thisMorning,
		NextPage:    len(nextPage) > 0,
	}, nil
}

type GetTaggableUsersRequest struct {
	Session *model.Session
	Params  GetTaggableUsersParams
}

type GetTaggableUsersParams struct {
	CardID globalid.ID `json:"cardID"`
}

func (p GetTaggableUsersParams) Sanitize() interface{} {
	return p
}

func (p GetTaggableUsersParams) Validate() error {
	return nil
}

type GetTaggableUsersResponse []model.TaggableUser

func (r *rpc) GetTaggableUsers(ctx context.Context, req GetTaggableUsersRequest) (*GetTaggableUsersResponse, error) {
	users, err := r.store.GetUsers()
	if err != nil {
		return nil, err
	}

	taggables := make([]model.TaggableUser, 0)
	for _, user := range users {
		if user.Username == store.RootUser.Username {
			continue
		}
		taggables = append(taggables, user.TaggableUser())
	}

	if req.Params.CardID != globalid.Nil {
		uniques := make(map[globalid.ID]bool)

		rootCard, err := r.store.GetCard(req.Params.CardID)
		if err != nil {
			return nil, err
		}

		if rootCard.AliasID != globalid.Nil {
			alias, err := r.store.GetAnonymousAlias(rootCard.AliasID)
			if err != nil {
				return nil, err
			}
			uniques[rootCard.AliasID] = true
			taggables = append(taggables, alias.TaggableUser())
		}

		for _, aliasID := range rootCard.AuthorToAlias {
			if !uniques[aliasID] {
				alias, err := r.store.GetAnonymousAlias(aliasID)
				if err != nil {
					return nil, err
				}
				uniques[rootCard.AliasID] = true
				taggables = append(taggables, model.TaggableUser{
					Username:       alias.Username,
					DisplayName:    alias.DisplayName,
					ProfilePicture: alias.ProfileImagePath,
					Anonymous:      true,
				})
			}
		}
	}

	return (*GetTaggableUsersResponse)(&taggables), nil
}

type GroupInvitesRequest struct {
	Params  GroupInvitesParams
	Session *model.Session
}

type GroupInvitesParams struct {
	InviteTokens []string `json:"tokens"`
}

func (p GroupInvitesParams) Sanitize() interface{} {
	return p
}

func (p GroupInvitesParams) Validate() error {
	return nil
}

type GroupInvitesResponse struct{}

func (r *rpc) GroupInvites(ctx context.Context, req GroupInvitesRequest) (*GroupInvitesResponse, error) {
	groupID := globalid.Next()

	err := r.store.ReassignInvitesByToken(req.Params.InviteTokens, req.Session.UserID)
	if err != nil {
		return nil, err
	}
	if len(req.Params.InviteTokens) > 1 {
		return nil, r.store.GroupInvitesByToken(req.Params.InviteTokens, groupID)
	}

	return nil, nil
}

type UnsubscribeFromCardRequest struct {
	Params  UnsubscribeFromCardParams
	Session *model.Session
}

type UnsubscribeFromCardParams struct {
	CardID globalid.ID `json:"cardID"`
}

func (p UnsubscribeFromCardParams) Sanitize() interface{} {
	return p
}

func (p UnsubscribeFromCardParams) Validate() error {
	return nil
}

type UnsubscribeFromCardResponse model.CardResponse

func (r *rpc) UnsubscribeFromCard(ctx context.Context, req UnsubscribeFromCardRequest) (*UnsubscribeFromCardResponse, error) {
	err := r.unsubscribeFromNotificationsForType(req.Session.UserID, req.Params.CardID, model.BoostType)

	if err != nil {
		return nil, err
	}

	err = r.unsubscribeFromNotificationsForType(req.Session.UserID, req.Params.CardID, model.CommentType)

	if err != nil {
		return nil, err
	}

	card, err := r.store.GetCard(req.Params.CardID)

	if err != nil {
		return nil, err
	}

	cR, err := CardResponse(r.store, card, req.Session.UserID)
	if err != nil {
		return nil, err
	}

	return (*UnsubscribeFromCardResponse)(cR), nil
}

type SubscribeToCardRequest struct {
	Params  SubscribeToCardParams
	Session *model.Session
}

type SubscribeToCardParams struct {
	CardID globalid.ID `json:"cardID"`
}

func (p SubscribeToCardParams) Sanitize() interface{} {
	return p
}

func (p SubscribeToCardParams) Validate() error {
	return nil
}

type SubscribeToCardResponse model.CardResponse

func (r *rpc) SubscribeToCard(ctx context.Context, req SubscribeToCardRequest) (*SubscribeToCardResponse, error) {
	card, err := r.store.GetCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}

	err = r.subscribeToNotificationsForType(req.Session.UserID, req.Params.CardID, model.CommentType)

	if err != nil {
		return nil, err
	}

	if req.Session.UserID == card.OwnerID {
		err = r.subscribeToNotificationsForType(req.Session.UserID, req.Params.CardID, model.BoostType)
		if err != nil {
			return nil, err
		}
	}

	cR, err := CardResponse(r.store, card, req.Session.UserID)
	if err != nil {
		return nil, err
	}

	return (*SubscribeToCardResponse)(cR), nil
}

type ReportCardRequest struct {
	Params  ReportCardParams
	Session *model.Session
}

type ReportCardParams struct {
	CardID globalid.ID `json:"cardID"`
}

func (p ReportCardParams) Sanitize() interface{} {
	return p
}

func (p ReportCardParams) Validate() error {
	return nil
}

type ReportCardResponse struct{}

func (r *rpc) ReportCard(ctx context.Context, req ReportCardRequest) (*ReportCardResponse, error) {
	card, err := r.store.GetCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}

	author, err := r.store.GetUser(card.OwnerID)
	if err != nil {
		return nil, err
	}

	return nil, r.notifier.NotifySlack("reports", fmt.Sprintf("%v reported <https://october.app/post/%v|%v's card>:\n```\n%v\n```", req.Session.User.DisplayName, req.Params.CardID, author.DisplayName, card.Content))
}

type SubmitFeedbackRequest struct {
	Params  SubmitFeedbackParams
	Session *model.Session
}

type SubmitFeedbackParams struct {
	Body string `json:"body"`
}

func (p SubmitFeedbackParams) Sanitize() interface{} {
	return p
}

func (p SubmitFeedbackParams) Validate() error {
	return nil
}

type SubmitFeedbackResponse struct{}

func (r *rpc) SubmitFeedback(ctx context.Context, req SubmitFeedbackRequest) (*SubmitFeedbackResponse, error) {
	return nil, r.notifier.NotifySlack("feedback", fmt.Sprintf("%v (<https://october.app/user/%v|@%v>) submitted feedback:\n```\n%v\n```", req.Session.User.DisplayName, req.Session.User.Username, req.Session.User.Username, req.Params.Body))
}

type TipCardRequest struct {
	Params  TipCardParams
	Session *model.Session
}

type TipCardParams struct {
	CardID    globalid.ID `json:"cardID"`
	Amount    int         `json:"amount"`
	Anonymous bool        `json:"anonymous"`
}

func (p TipCardParams) Sanitize() interface{} {
	return p
}

func (p TipCardParams) Validate() error {
	return nil
}

type TipCardResponse struct {
	NewBalance *model.CoinBalances `json:"newBalances"`
}

func (r *rpc) TipCard(ctx context.Context, req TipCardRequest) (*TipCardResponse, error) {
	if !r.userCanAffordTip(req.Session.UserID, req.Params.Amount) {
		return nil, model.ErrInsufficientBalance
	}

	card, err := r.store.GetCard(req.Params.CardID)
	if err != nil {
		return nil, err
	}

	if req.Session.UserID == card.OwnerID {
		// just ignore self-tips
		bal, berr := r.store.GetCurrentBalance(req.Session.UserID)
		if berr != nil {
			return nil, berr
		}

		return &TipCardResponse{
			NewBalance: bal,
		}, nil
	}

	tip := &model.UserTip{
		UserID:    req.Session.UserID,
		CardID:    req.Params.CardID,
		Amount:    req.Params.Amount,
		Anonymous: req.Params.Anonymous,
	}

	// set the alias if there's an existing one
	if req.Params.Anonymous {
		authorToAlias := card.AuthorToAlias
		if card.ThreadRootID != globalid.Nil {
			root, threadErr := r.store.GetCard(card.ThreadRootID)
			if threadErr != nil {
				return nil, threadErr
			}

			authorToAlias = root.AuthorToAlias
		}

		aliasID, ok := authorToAlias[req.Session.UserID]

		if ok {
			tip.AliasID = aliasID
		}
	}

	err = r.store.SaveUserTip(tip)
	if err != nil {
		return nil, err
	}

	tipperNewBalance, tippedNewBalance, err := r.updateTokensForTip(req.Session.UserID, card.ID, card.OwnerID, req.Params.Amount)
	if err != nil {
		return nil, err
	}

	err = r.pusher.UpdateCoinBalance(ctx, card.OwnerID, tippedNewBalance)
	if err != nil {
		r.log.Error(err)
	}

	go r.pushUpdateEngagement(ctx, req.Session, req.Params.CardID)

	return &TipCardResponse{
		NewBalance: tipperNewBalance,
	}, nil
}

type BlockUserRequest struct {
	Params  BlockUserParams
	Session *model.Session
}

type BlockUserParams struct {
	UserID    globalid.ID `json:"userID"`
	AliasID   globalid.ID `json:"aliasID"`
	ForThread globalid.ID `json:"forThread"`
}

func (p BlockUserParams) Sanitize() interface{} {
	return p
}

func (p BlockUserParams) Validate() error {
	return nil
}

type BlockUserResponse struct{}

func (r *rpc) BlockUser(ctx context.Context, req BlockUserRequest) (*BlockUserResponse, error) {
	if req.Params.UserID != globalid.Nil && req.Params.UserID != req.Session.UserID {
		return nil, r.store.BlockUser(req.Session.UserID, req.Params.UserID)
	}
	if req.Params.AliasID != globalid.Nil && req.Params.ForThread != globalid.Nil {
		return nil, r.store.BlockAnonUserInThread(req.Session.UserID, req.Params.AliasID, req.Params.ForThread)
	}
	return nil, nil
}

type GetCardsForChannelRequest struct {
	Params  GetCardsForChannelParams
	Session *model.Session
}

type GetCardsForChannelParams struct {
	ChannelID   globalid.ID `json:"channelID"`
	ChannelName string      `json:"channelName"`
	PageSize    int         `json:"pageSize"`
	PageNumber  int         `json:"pageNumber"`
}

func (p GetCardsForChannelParams) Sanitize() interface{} {
	return p
}

func (p GetCardsForChannelParams) Validate() error {
	return nil
}

type GetCardsForChannelResponse struct {
	Cards               []*model.CardResponse `json:"cards"`
	NextPage            bool                  `json:"hasNextPage"`
	Channel             *model.Channel        `json:"channel"`
	SubscribedToChannel bool                  `json:"subscribed"`
}

func (r *rpc) GetCardsForChannel(ctx context.Context, req GetCardsForChannelRequest) (*GetCardsForChannelResponse, error) {
	channelID := req.Params.ChannelID
	var chann *model.Channel
	var err error

	// try by name
	if req.Params.ChannelName != "" {
		chann, err = r.store.GetChannelByHandle(strings.ToLower(req.Params.ChannelName))

		// if it worked, set the ID, otherwise try by ID
		if chann != nil {
			channelID = chann.ID
		} else if errors.Cause(err) == sql.ErrNoRows {
			_, err = r.store.GetChannel(req.Params.ChannelID)
			if err != nil {
				return nil, err
			}
		}
	} else if channelID != globalid.Nil {
		chann, err = r.store.GetChannel(channelID)
		if err != nil {
			return nil, err
		}
	}

	subbed, err := r.store.GetIsSubscribed(req.Session.UserID, channelID)
	if err != nil {
		return nil, err
	}

	cards, err := r.store.GetCardsForChannel(channelID, req.Params.PageSize, req.Params.PageNumber, req.Session.UserID)
	if err != nil {
		return nil, err
	}

	lookAhead, err := r.store.GetCardsForChannel(channelID, req.Params.PageSize, req.Params.PageNumber+1, req.Session.UserID)
	if err != nil {
		return nil, err
	}

	cardResponses, err := r.cardResponses(cards, req.Session.UserID)
	if err != nil {
		return nil, err
	}
	result := &GetCardsForChannelResponse{
		Cards:               cardResponses,
		NextPage:            len(lookAhead) != 0,
		Channel:             chann,
		SubscribedToChannel: subbed,
	}

	return result, nil
}

type UpdateChannelSubscriptionRequest struct {
	Params  UpdateChannelSubscriptionParams
	Session *model.Session
}

type UpdateChannelSubscriptionParams struct {
	ChannelID  globalid.ID `json:"channelID"`
	Subscribed bool        `json:"subscribed"`
}

func (p UpdateChannelSubscriptionParams) Sanitize() interface{} {
	return p
}

func (p UpdateChannelSubscriptionParams) Validate() error {
	return nil
}

type UpdateChannelSubscriptionResponse struct{}

func (r *rpc) UpdateChannelSubscription(ctx context.Context, req UpdateChannelSubscriptionRequest) (*UpdateChannelSubscriptionResponse, error) {
	if req.Params.Subscribed {
		return nil, r.store.JoinChannel(req.Session.UserID, req.Params.ChannelID)
	}

	return nil, r.store.LeaveChannel(req.Session.UserID, req.Params.ChannelID)
}

type GetChannelsRequest struct {
	Params  GetChannelsParams
	Session *model.Session
}

type GetChannelsParams struct {
	OnlySubscribed bool `json:"onlySubscribed"`
	OnlyPostable   bool `json:"onlyPostable"`
	HideEmpty      bool `json:"hideEmpty"`
}

func (p GetChannelsParams) Sanitize() interface{} {
	return p
}

func (p GetChannelsParams) Validate() error {
	return nil
}

type UserChannel struct {
	Channel     *model.Channel `json:"channel"`
	MemberCount int            `json:"memberCount"`
	Subscribed  bool           `json:"subscribed"`
}

type GetChannelsResponse []*UserChannel

func (r *rpc) GetChannels(ctx context.Context, req GetChannelsRequest) (*GetChannelsResponse, error) {
	chans, err := r.store.GetChannelsForUser(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	infos, err := r.store.GetChannelInfos(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	infosMap := make(map[globalid.ID]*model.ChannelUserInfo)

	for _, info := range infos {
		infosMap[info.ChannelID] = info
	}

	chansWithSub := make([]*UserChannel, 0)
	for _, chann := range chans {
		// if you're subscribed OR onlySubscribed is false
		onlyShowSubscribed := req.Params.OnlySubscribed || req.Params.OnlyPostable
		chanInfo, ok := infosMap[chann.ID]

		if ok && (chanInfo.Subscribed || !onlyShowSubscribed) && (chanInfo.MemberCount > 0 || !req.Params.HideEmpty) {
			chanWithSub := &UserChannel{
				Channel:     chann,
				Subscribed:  chanInfo.Subscribed,
				MemberCount: chanInfo.MemberCount,
			}
			chansWithSub = append(chansWithSub, chanWithSub)
		}
	}
	resp := GetChannelsResponse(chansWithSub)
	return &resp, nil
}

type JoinChannelRequest struct {
	Params  JoinChannelParams
	Session *model.Session
}

type JoinChannelParams struct {
	ChannelID globalid.ID `json:"channelID"`
}

func (p JoinChannelParams) Sanitize() interface{} {
	return p
}

func (p JoinChannelParams) Validate() error {
	return nil
}

type JoinChannelResponse struct{}

func (r *rpc) JoinChannel(ctx context.Context, req JoinChannelRequest) (*JoinChannelResponse, error) {
	return nil, r.store.JoinChannel(req.Session.UserID, req.Params.ChannelID)
}

type LeaveChannelRequest struct {
	Params  LeaveChannelParams
	Session *model.Session
}

type LeaveChannelParams struct {
	ChannelID globalid.ID `json:"channelID"`
}

func (p LeaveChannelParams) Sanitize() interface{} {
	return p
}

func (p LeaveChannelParams) Validate() error {
	return nil
}

type LeaveChannelResponse struct{}

func (r *rpc) LeaveChannel(ctx context.Context, req LeaveChannelRequest) (*LeaveChannelResponse, error) {
	return nil, r.store.LeaveChannel(req.Session.UserID, req.Params.ChannelID)
}

type MuteChannelRequest struct {
	Params  MuteChannelParams
	Session *model.Session
}

type MuteChannelParams struct {
	ChannelID globalid.ID `json:"channelID"`
}

func (p MuteChannelParams) Sanitize() interface{} {
	return p
}

func (p MuteChannelParams) Validate() error {
	return nil
}

type MuteChannelResponse struct{}

func (r *rpc) MuteChannel(ctx context.Context, req MuteChannelRequest) (*MuteChannelResponse, error) {
	return nil, r.store.MuteChannel(req.Session.UserID, req.Params.ChannelID)
}

type UnmuteChannelRequest struct {
	Params  UnmuteChannelParams
	Session *model.Session
}

type UnmuteChannelParams struct {
	ChannelID globalid.ID `json:"channelID"`
}

func (p UnmuteChannelParams) Sanitize() interface{} {
	return p
}

func (p UnmuteChannelParams) Validate() error {
	return nil
}

type UnmuteChannelResponse struct{}

func (r *rpc) UnmuteChannel(ctx context.Context, req UnmuteChannelRequest) (*UnmuteChannelResponse, error) {
	return nil, r.store.UnmuteChannel(req.Session.UserID, req.Params.ChannelID)
}

type MuteUserRequest struct {
	Params  MuteUserParams
	Session *model.Session
}

type MuteUserParams struct {
	// UserID is the ID of the user to be followed.
	UserID globalid.ID `json:"userID"`
}

func (p MuteUserParams) Sanitize() interface{} {
	return p
}

func (p MuteUserParams) Validate() error {
	return nil
}

type MuteUserResponse struct{}

// FollowUser adds the given user to the list of followers.
func (r *rpc) MuteUser(ctx context.Context, req MuteUserRequest) (*MuteUserResponse, error) {
	return nil, r.store.MuteUser(req.Session.UserID, req.Params.UserID)
}

type UnmuteUserRequest struct {
	Params  UnmuteUserParams
	Session *model.Session
}

type UnmuteUserParams struct {
	// UserID is the ID of the user to be followed.
	UserID globalid.ID `json:"userID"`
}

func (p UnmuteUserParams) Sanitize() interface{} {
	return p
}

func (p UnmuteUserParams) Validate() error {
	return nil
}

type UnmuteUserResponse struct{}

// FollowUser adds the given user to the list of followers.
func (r *rpc) UnmuteUser(ctx context.Context, req UnmuteUserRequest) (*UnmuteUserResponse, error) {
	return nil, r.store.UnmuteUser(req.Session.UserID, req.Params.UserID)
}

type MuteThreadRequest struct {
	Params  MuteThreadParams
	Session *model.Session
}

type MuteThreadParams struct {
	// UserID is the ID of the user to be followed.
	CardID globalid.ID `json:"cardID"`
}

func (p MuteThreadParams) Sanitize() interface{} {
	return p
}

func (p MuteThreadParams) Validate() error {
	return nil
}

type MuteThreadResponse struct{}

// FollowUser adds the given user to the list of followers.
func (r *rpc) MuteThread(ctx context.Context, req MuteThreadRequest) (*MuteThreadResponse, error) {
	return nil, r.store.MuteThread(req.Session.UserID, req.Params.CardID)
}

type UnmuteThreadRequest struct {
	Params  UnmuteThreadParams
	Session *model.Session
}

type UnmuteThreadParams struct {
	// UserID is the ID of the user to be followed.
	CardID globalid.ID `json:"cardID"`
}

func (p UnmuteThreadParams) Sanitize() interface{} {
	return p
}

func (p UnmuteThreadParams) Validate() error {
	return nil
}

type UnmuteThreadResponse struct{}

// FollowUser adds the given user to the list of followers.
func (r *rpc) UnmuteThread(ctx context.Context, req UnmuteThreadRequest) (*UnmuteThreadResponse, error) {
	return nil, r.store.UnmuteUser(req.Session.UserID, req.Params.CardID)
}

func (r *rpc) subscribeToNotificationsForType(userID, cardID globalid.ID, typ string) error {
	return r.store.SubscribeToCard(userID, cardID, typ)
}

func (r *rpc) unsubscribeFromNotificationsForType(userID, cardID globalid.ID, typ string) error {
	err := r.store.UnsubscribeFromCard(userID, cardID, typ)
	if err != nil {
		return err
	}

	notif, err := r.store.LatestForType(userID, cardID, typ, true)

	if errors.Cause(err) != sql.ErrNoRows && err != nil {
		return err
	}

	// if there's no open notif, we're already unsubscribed
	if errors.Cause(err) == sql.ErrNoRows {
		return nil
	}

	return r.store.UpdateNotificationsOpened([]globalid.ID{notif.ID})
}

type CreateChannelRequest struct {
	Params  CreateChannelParams
	Session *model.Session
}

type CreateChannelParams struct {
	// UserID is the ID of the user to be followed.
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p CreateChannelParams) Sanitize() interface{} {
	return p
}

func (p CreateChannelParams) Validate() error {
	return nil
}

type CreateChannelResponse struct {
	Channel    *model.Channel      `json:"channel"`
	NewBalance *model.CoinBalances `json:"newBalances"`
}

// FollowUser adds the given user to the list of followers.
func (r *rpc) CreateChannel(ctx context.Context, req CreateChannelRequest) (*CreateChannelResponse, error) {
	if !r.userCanAffordChannel(req.Session.UserID) {
		return nil, model.ErrInsufficientBalance
	}

	err := model.ValidateChannelName(req.Params.Name)
	if err != nil {
		return nil, err
	}

	newChan := &model.Channel{
		Name:        req.Params.Name,
		Description: req.Params.Description,
		OwnerID:     req.Session.UserID,
		Handle:      strings.ToLower(req.Params.Name),
	}

	existingChan, err := r.store.GetChannelByHandle(newChan.Handle)

	if existingChan != nil || err == nil {
		return nil, ErrChannelExists
	} else if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, err
	}

	err = r.store.SaveChannel(newChan)
	if err != nil {
		return nil, err
	}

	// Auto-join a channel you created
	err = r.store.JoinChannel(req.Session.UserID, newChan.ID)
	if err != nil {
		return nil, err
	}

	err = r.indexer.IndexChannel(newChan)
	if err != nil {
		r.log.Error(err)
	}

	nB, err := r.updateTokensForBuyChannel(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	defer r.notifySlackChannelCreation(req.Session.User, newChan)

	return &CreateChannelResponse{
		Channel:    newChan,
		NewBalance: nB,
	}, nil
}

func (r *rpc) notifySlackChannelCreation(user *model.User, channel *model.Channel) {
	msg := fmt.Sprintf("%s created a new channel <https://october.app/channel/%v|%s>.", user.Username, channel.ID, channel.Handle)
	err := r.notifier.NotifySlack("engagement", msg)
	if err != nil {
		r.log.Error(err)
	}
}

type GetPopularCardsRequest struct {
	Params  GetPopularCardsParams
	Session *model.Session
}

type GetPopularCardsParams struct {
	// PerPage specifies the number of cards returned per request.
	PerPage int `json:"pageSize"`
	// Page for pagination; the page to retrieve.
	Page int `json:"pageNumber"`
}

func (p GetPopularCardsParams) Sanitize() interface{} {
	return p
}

func (p GetPopularCardsParams) Validate() error {
	if p.PerPage == 0 {
		return errors.New("pageSize must be > 0")
	}
	return nil
}

type GetPopularCardsResponse model.CardsResponse

// GetCards retrieves the feed for the given user.
func (r *rpc) GetPopularCards(ctx context.Context, req GetPopularCardsRequest) (*GetPopularCardsResponse, error) {
	if req.Params.Page == 0 {
		/*err := r.store.UpdatePopularRanksForUser(req.Session.UserID)
		if err != nil {
			return nil, err
		}*/

		// pop ranks for cards for the last 2 weeks
		popRanks, err := r.store.GetCardsForPopularRankSince(time.Now().Add(time.Hour * -24 * 14))
		if err != nil {
			return nil, err
		}

		sort.Slice(popRanks, func(i, j int) bool {
			return popRanks[i].Rank() > popRanks[j].Rank()
		})

		ids := make([]globalid.ID, len(popRanks))

		for i, popRanks := range popRanks {
			ids[i] = popRanks.CardID
		}

		err = r.store.UpdatePopularRanksWithList(req.Session.UserID, ids)
		if err != nil {
			return nil, err
		}
	}

	cards, err := r.store.GetPopularRankCardsForUser(req.Session.UserID, req.Params.PerPage, req.Params.Page)
	if err != nil {
		return nil, err
	}
	lookAhead, err := r.store.GetPopularRankCardsForUser(req.Session.UserID, req.Params.PerPage, req.Params.Page+1)
	if err != nil {
		return nil, err
	}

	cardResponses, err := r.responses.FeedCardResponses(cards, req.Session.UserID)
	if err != nil {
		return nil, err
	}
	result := GetPopularCardsResponse(model.CardsResponse{
		Cards:    cardResponses,
		NextPage: len(lookAhead) != 0,
	})
	return &result, nil
}

type GetActionCostsRequest struct {
	Params  GetActionCostsParams
	Session *model.Session
}

type GetActionCostsParams struct{}

func (p GetActionCostsParams) Sanitize() interface{} {
	return p
}

func (p GetActionCostsParams) Validate() error {
	return nil
}

type GetActionCostsResponse struct {
	LikeCost          int64 `json:"likeCost"`
	AnonLikeCost      int64 `json:"anonLikeCost"`
	DislikeCost       int64 `json:"dislikeCost"`
	PostCost          int64 `json:"postCost"`
	AnonPostCost      int64 `json:"anonPostCost"`
	CommentCost       int64 `json:"commentCost"`
	AnonCommentCost   int64 `json:"anonCommentCost"`
	CreateChannelCost int64 `json:"createChannelCost"`
	ThreadAliasCost   int64 `json:"threadAliasCost"`
	PostAliasCost     int64 `json:"postAliasCost"`
	UnitsPerCoin      int64 `json:"unitsPerCoin"`
}

func (r *rpc) GetActionCosts(ctx context.Context, req GetActionCostsRequest) (*GetActionCostsResponse, error) {
	return &GetActionCostsResponse{
		LikeCost:          0,
		AnonLikeCost:      0,
		DislikeCost:       0,
		PostCost:          0,
		AnonPostCost:      0,
		CommentCost:       0,
		AnonCommentCost:   0,
		ThreadAliasCost:   int64(r.coinmanager.Config.BoughtThreadAlias),
		PostAliasCost:     int64(r.coinmanager.Config.BoughtPostAlias),
		CreateChannelCost: int64(r.coinmanager.Config.BoughtChannel),
		UnitsPerCoin:      r.config.UnitsPerCoin,
	}, nil
}

type UseInviteCodeRequest struct {
	Params  UseInviteCodeParams
	Session *model.Session
}

type UseInviteCodeParams struct {
	Token string `json:"token"`
}

func (p UseInviteCodeParams) Sanitize() interface{} {
	return p
}

func (p UseInviteCodeParams) Validate() error {
	if p.Token == "" {
		return model.ErrInvalidInviteCode
	}
	return nil
}

type UseInviteCodeResponse struct {
	NewBalance *model.CoinBalances `json:"newBalances"`
}

func (r *rpc) UseInviteCode(ctx context.Context, req UseInviteCodeRequest) (*UseInviteCodeResponse, error) {
	user, err := r.store.GetUser(req.Session.UserID)
	if err != nil {
		return nil, err
	}

	if user.JoinedFromInvite != globalid.Nil {
		return nil, ErrUserAlreadyRedeemedCode
	}

	// Find/validate the invite token used
	invite, err := r.store.GetInviteByToken(strings.ToUpper(req.Params.Token))
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		return nil, model.ErrInvalidInviteCode
	} else if err != nil {
		return nil, err
	}

	// get the inviting user
	invitingUser, err := r.store.GetUser(invite.NodeID)
	if err != nil {
		return nil, err
	}

	// update the welcome notification
	notif, err := r.store.LatestForType(req.Session.UserID, globalid.Nil, model.IntroductionType, false)
	if err == nil && notif != nil {
		notif.TargetID = invitingUser.ID

		err = r.store.SaveNotification(notif)
		if err != nil {
			return nil, err
		}
	}

	// Set which invite was used
	user.JoinedFromInvite = invite.ID

	if invitingUser != nil && invitingUser.ShadowbannedAt.Valid {
		user.ShadowbannedAt = model.NewDBTime(time.Now().UTC())
	}

	// Reassign grouped invites
	if invite.GroupID != globalid.Nil {
		err = r.store.ReassignInviterForGroup(invite.ID, user.ID)
		if err != nil {
			return nil, err
		}
	}

	// subscribe to your inviter
	err = r.store.SaveFollower(user.ID, invite.NodeID)
	if err != nil {
		return nil, err
	}

	// Join the invite's channel if this is a channel invite
	if invite.ChannelID != globalid.Nil {
		err = r.store.LeaveAllChannels(user.ID)
		if err != nil {
			return nil, err
		}
		err = r.store.JoinChannel(user.ID, invite.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	err = r.store.SaveUser(user)
	if err != nil {
		return nil, err
	}

	// notify inviter of accepted inivte
	notif = &model.Notification{
		ID:       globalid.Next(),
		UserID:   invite.NodeID,
		TargetID: user.ID,
		Type:     model.InviteAcceptedType,
	}
	err = r.store.SaveNotification(notif)
	if err != nil {
		return nil, err
	}
	exNotif, eerr := r.notifications.ExportNotification(notif)
	if eerr != nil {
		return nil, eerr
	}
	err = r.pusher.NewNotification(ctx, req.Session, exNotif)
	if err != nil {
		r.log.Error(err)
	}
	err = r.notifier.NotifyPush(exNotif)
	if err != nil {
		r.log.Error(err)
	}

	/* reward inviter with tokens */
	newBalance, err := r.tokensForInvite(invite.NodeID)
	if err != nil {
		r.log.Error(err)
	}
	err = r.pusher.UpdateCoinBalance(ctx, invite.NodeID, newBalance)
	if err != nil {
		r.log.Error(err)
	}

	// Add the code reward to their coins
	newBalance, err = r.updateTokensForRedeemCode(req.Session.UserID)
	if err != nil {
		r.log.Error(err)
	}

	return &UseInviteCodeResponse{
		NewBalance: newBalance,
	}, nil
}

type RequestValidationRequest struct {
	Params  RequestValidationParams
	Session *model.Session
}

type RequestValidationParams struct {
	CountryCode string `json:"countryCode"`
	PhoneNumber string `json:"phoneNumber"`
}

func (p RequestValidationParams) Sanitize() interface{} {
	return p
}

func (p RequestValidationParams) Validate() error {
	return nil
}

type RequestValidationResponse struct{}

func (r *rpc) RequestValidation(ctx context.Context, req RequestValidationRequest) (*RequestValidationResponse, error) {
	if r.config.AutoVerify {
		return nil, nil
	}

	requestURL := "https://api.authy.com/protected/json/phones/verification/start"
	values := url.Values{
		"api_key":      {r.config.AuthyAPIKey},
		"via":          {"sms"},
		"phone_number": {req.Params.PhoneNumber},
		"country_code": {req.Params.CountryCode},
	}

	body := bytes.NewBufferString(values.Encode())

	httpreq, err := http.NewRequest("POST", requestURL, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(httpreq)
	if err != nil {
		return nil, err
	}

	var dat map[string]interface{}
	repbody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			r.log.Error(cerr)
		}
	}()

	if err := json.Unmarshal(repbody, &dat); err != nil {
		return nil, err
	}
	successful, ok := dat["success"]

	if !ok || !successful.(bool) {
		return nil, errors.New("Failed to send code")
	}

	return nil, nil
}

type ConfirmValidationRequest struct {
	Params  ConfirmValidationParams
	Session *model.Session
}

type ConfirmValidationParams struct {
	CountryCode string `json:"countryCode"`
	PhoneNumber string `json:"phoneNumber"`
	Code        string `json:"code"`
}

func (p ConfirmValidationParams) Sanitize() interface{} {
	return p
}

func (p ConfirmValidationParams) Validate() error {
	if p.CountryCode == "" || p.PhoneNumber == "" || p.Code == "" {
		return errors.New("Invalid Code")
	}
	return nil
}

type ConfirmValidationResponse struct{}

func (r *rpc) ConfirmValidation(ctx context.Context, req ConfirmValidationRequest) (*ConfirmValidationResponse, error) {
	if r.config.AutoVerify {
		user, uerr := r.store.GetUser(req.Session.UserID)
		if uerr != nil {
			return nil, uerr
		}
		user.IsVerified = true
		return nil, r.store.SaveUser(user)
	}

	endpoint := fmt.Sprintf("https://api.authy.com/protected/json/phones/verification/check?phone_number=%v&country_code=%v&verification_code=%v", req.Params.PhoneNumber, req.Params.CountryCode, req.Params.Code)
	httpreq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	httpreq.Header.Set("X-Authy-API-Key", r.config.AuthyAPIKey)
	client := &http.Client{}
	resp, err := client.Do(httpreq)
	if err != nil {
		return nil, err
	}

	repbody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			r.log.Error(cerr)
		}
	}()

	var dat map[string]interface{}
	err = json.Unmarshal(repbody, &dat)
	if err != nil {
		return nil, err
	}

	successful, ok := dat["success"]

	if ok && successful.(bool) {
		user, uerr := r.store.GetUser(req.Session.UserID)
		if uerr != nil {
			return nil, uerr
		}

		user.IsVerified = true

		uerr = r.store.SaveUser(user)
		if uerr != nil {
			return nil, uerr
		}
	}
	return nil, nil
}

type ValidateChannelNameRequest struct {
	Params  ValidateChannelNameParams
	Session *model.Session
}

type ValidateChannelNameParams struct {
	// Username is the username to be checked.
	Name string `json:"channelName"`
}

func (p ValidateChannelNameParams) Sanitize() interface{} {
	return p
}

func (p ValidateChannelNameParams) Validate() error {
	return nil
}

type ValidateChannelNameResponse struct{}

// ValidateUsername checks whether a given username can be used for a new user
// or to change the username of an existing user.
func (r *rpc) ValidateChannelName(ctx context.Context, req ValidateChannelNameRequest) (*ValidateChannelNameResponse, error) {
	err := model.ValidateChannelName(req.Params.Name)
	if err != nil {
		return nil, err
	}
	// check to see if someone already has this username
	c, err := r.store.GetChannelByHandle(strings.ToLower(req.Params.Name))
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, err
	}
	if c != nil {
		return nil, ErrChannelExists
	}
	return nil, nil
}

type GetChannelRequest struct {
	Params  GetChannelParams
	Session *model.Session
}

type GetChannelParams struct {
	Name string      `json:"channelName,omitempty"`
	ID   globalid.ID `json:"channelID,omitempty"`
}

func (p GetChannelParams) Sanitize() interface{} {
	return p
}

func (p GetChannelParams) Validate() error {
	return nil
}

type GetChannelResponse UserChannel

func (r *rpc) GetChannel(ctx context.Context, req GetChannelRequest) (*GetChannelResponse, error) {
	var chann *model.Channel
	var err error
	if req.Params.Name != "" {
		chann, err = r.store.GetChannelByHandle(strings.ToLower(req.Params.Name))

		if errors.Cause(err) == sql.ErrNoRows {
			chann, err = r.store.GetChannel(req.Params.ID)
		}
	} else {
		chann, err = r.store.GetChannel(req.Params.ID)
	}

	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		return nil, datastore.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	subbed, err := r.store.GetIsSubscribed(req.Session.UserID, chann.ID)
	if err != nil {
		return nil, err
	}

	subCount, err := r.store.GetSubscriberCount(chann.ID)
	if err != nil {
		return nil, err
	}

	return (*GetChannelResponse)(&UserChannel{
		Channel:     chann,
		Subscribed:  subbed,
		MemberCount: subCount,
	}), nil
}

type CanAffordAnonymousPostRequest struct {
	Params  CanAffordAnonymousPostParams
	Session *model.Session
}

type CanAffordAnonymousPostParams struct {
	ThreadRootID globalid.ID `json:"threadRootID"`
}

func (p CanAffordAnonymousPostParams) Sanitize() interface{} {
	return p
}

func (p CanAffordAnonymousPostParams) Validate() error {
	return nil
}

type CanAffordAnonymousPostResponse struct{}

func (r *rpc) CanAffordAnonymousPost(ctx context.Context, req CanAffordAnonymousPostRequest) (*CanAffordAnonymousPostResponse, error) {
	var err error
	// check costs
	hasPostingAlias := false
	if req.Params.ThreadRootID != globalid.Nil {
		hasPostingAlias, err = r.userHasPostingAliasInThread(req.Session.UserID, req.Params.ThreadRootID)
		if err != nil {
			return nil, err
		}
	}
	isReply := req.Params.ThreadRootID != globalid.Nil

	if !hasPostingAlias {
		if isReply && !r.userCanAffordThreadAlias(req.Session.UserID) {
			return nil, model.ErrInsufficientBalance
		} else if !r.userCanAffordPostAlias(req.Session.UserID) {
			return nil, model.ErrInsufficientBalance
		}
	}
	return nil, nil
}

type GetLeaderboardRequest struct {
	Params  GetLeaderboardParams
	Session *model.Session
}

type GetLeaderboardParams struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}

func (p GetLeaderboardParams) Sanitize() interface{} {
	return p
}

func (p GetLeaderboardParams) Validate() error {
	return nil
}

type LeaderboardUser struct {
	User        *model.ExportedUser `json:"user"`
	Rank        int                 `json:"rank"`
	CoinsEarned int                 `json:"coinsEarned"`
}

type GetLeaderboardResponse struct {
	Users    []*LeaderboardUser `json:"users"`
	NextPage bool               `json:"hasNextPage"`
}

func (r *rpc) GetLeaderboard(ctx context.Context, req GetLeaderboardRequest) (*GetLeaderboardResponse, error) {
	currentRankings, err := r.store.GetLeaderboardRankings(req.Params.PageSize, req.Params.PageNumber)
	if err != nil {
		return nil, err
	}

	nextPage, err := r.store.GetLeaderboardRankings(req.Params.PageSize, req.Params.PageNumber+1)
	if err != nil {
		return nil, err
	}

	exportedUsers := make([]*LeaderboardUser, len(currentRankings))
	for i, ranking := range currentRankings {
		user, err := r.store.GetUser(ranking.UserID)
		if err != nil {
			return nil, err
		}

		exportedUsers[i] = &LeaderboardUser{
			User:        user.Export(req.Session.UserID),
			Rank:        int(ranking.Rank),
			CoinsEarned: int(ranking.CoinsEarned),
		}
	}

	return &GetLeaderboardResponse{
		Users:    exportedUsers,
		NextPage: len(nextPage) > 0,
	}, nil
}

func extractUsername(ctx context.Context) string {
	value := ctx.Value("username")
	if value == nil {
		return "no user"
	}
	username, ok := value.(string)
	if !ok {
		return "invalid username field"
	}
	return username
}
