package server_test

import (
	"context"

	"github.com/october93/engine/rpc"
)

type rpcMock struct {
	// clients
	Authf               func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error)
	ResetPasswordf      func(ctx context.Context, req rpc.ResetPasswordRequest) (*rpc.ResetPasswordResponse, error)
	ValidateInviteCodef func(ctx context.Context, req rpc.ValidateInviteCodeRequest) (*rpc.ValidateInviteCodeResponse, error)
	AddToWaitlistf      func(ctx context.Context, req rpc.AddToWaitlistRequest) (*rpc.AddToWaitlistResponse, error)

	// authenticated client rpc
	Logoutf                    func(ctx context.Context, req rpc.LogoutRequest) (*rpc.LogoutResponse, error)
	GetCardsf                  func(ctx context.Context, req rpc.GetCardsRequest) (*rpc.GetCardsResponse, error)
	GetCardf                   func(ctx context.Context, req rpc.GetCardRequest) (*rpc.GetCardResponse, error)
	ReactToCardf               func(ctx context.Context, req rpc.ReactToCardRequest) (*rpc.ReactToCardResponse, error)
	VoteOnCardf                func(ctx context.Context, req rpc.VoteOnCardRequest) (*rpc.VoteOnCardResponse, error)
	PostCardf                  func(ctx context.Context, req rpc.PostCardRequest) (*rpc.PostCardResponse, error)
	NewInvitef                 func(ctx context.Context, req rpc.NewInviteRequest) (*rpc.NewInviteResponse, error)
	RegisterDevicef            func(ctx context.Context, req rpc.RegisterDeviceRequest) (*rpc.RegisterDeviceResponse, error)
	UnregisterDevicef          func(ctx context.Context, req rpc.UnregisterDeviceRequest) (*rpc.UnregisterDeviceResponse, error)
	UpdateSettingsf            func(ctx context.Context, req rpc.UpdateSettingsRequest) (*rpc.UpdateSettingsResponse, error)
	GetUserf                   func(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error)
	ValidateUsernamef          func(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error)
	GetThreadf                 func(ctx context.Context, req rpc.GetThreadRequest) (*rpc.GetThreadResponse, error)
	GetNotificationsf          func(ctx context.Context, req rpc.GetNotificationsRequest) (*rpc.GetNotificationsResponse, error)
	UpdateNotificationsf       func(ctx context.Context, req rpc.UpdateNotificationsRequest) (*rpc.UpdateNotificationsResponse, error)
	GetAnonymousHandlef        func(ctx context.Context, req rpc.GetAnonymousHandleRequest) (*rpc.GetAnonymousHandleResponse, error)
	DeleteCardf                func(ctx context.Context, req rpc.DeleteCardRequest) (*rpc.DeleteCardResponse, error)
	FollowUserf                func(ctx context.Context, req rpc.FollowUserRequest) (*rpc.FollowUserResponse, error)
	UnfollowUserf              func(ctx context.Context, req rpc.UnfollowUserRequest) (*rpc.UnfollowUserResponse, error)
	GetFollowingUsersf         func(ctx context.Context, req rpc.GetFollowingUsersRequest) (*rpc.GetFollowingUsersResponse, error)
	GetPostsForUserf           func(ctx context.Context, req rpc.GetPostsForUserRequest) (*rpc.GetPostsForUserResponse, error)
	GetTagsf                   func(ctx context.Context, req rpc.GetTagsRequest) (*rpc.GetTagsResponse, error)
	GetFeaturesForUserf        func(ctx context.Context, req rpc.GetFeaturesForUserRequest) (*rpc.GetFeaturesForUserResponse, error)
	PreviewContentf            func(ctx context.Context, req rpc.PreviewContentRequest) (*rpc.PreviewContentResponse, error)
	UploadImagef               func(ctx context.Context, req rpc.UploadImageRequest) (*rpc.UploadImageResponse, error)
	GetTaggableUsersf          func(ctx context.Context, req rpc.GetTaggableUsersRequest) (*rpc.GetTaggableUsersResponse, error)
	ModifyCardScoref           func(ctx context.Context, req rpc.ModifyCardScoreRequest) (*rpc.ModifyCardScoreResponse, error)
	GetInvitesf                func(ctx context.Context, req rpc.GetInvitesRequest) (*rpc.GetInvitesResponse, error)
	GetOnboardingDataf         func(ctx context.Context, req rpc.GetOnboardingDataRequest) (*rpc.GetOnboardingDataResponse, error)
	GetMyNetworkf              func(ctx context.Context, req rpc.GetMyNetworkRequest) (*rpc.GetMyNetworkResponse, error)
	UnsubscribeFromCardf       func(ctx context.Context, req rpc.UnsubscribeFromCardRequest) (*rpc.UnsubscribeFromCardResponse, error)
	SubscribeToCardf           func(ctx context.Context, req rpc.SubscribeToCardRequest) (*rpc.SubscribeToCardResponse, error)
	GroupInvitesf              func(ctx context.Context, req rpc.GroupInvitesRequest) (*rpc.GroupInvitesResponse, error)
	ReportCardf                func(ctx context.Context, req rpc.ReportCardRequest) (*rpc.ReportCardResponse, error)
	BlockUserf                 func(ctx context.Context, req rpc.BlockUserRequest) (*rpc.BlockUserResponse, error)
	GetCardsForChannelf        func(ctx context.Context, req rpc.GetCardsForChannelRequest) (*rpc.GetCardsForChannelResponse, error)
	UpdateChannelSubscriptionf func(ctx context.Context, req rpc.UpdateChannelSubscriptionRequest) (*rpc.UpdateChannelSubscriptionResponse, error)
	GetChannelsf               func(ctx context.Context, req rpc.GetChannelsRequest) (*rpc.GetChannelsResponse, error)
	JoinChannelf               func(ctx context.Context, req rpc.JoinChannelRequest) (*rpc.JoinChannelResponse, error)
	LeaveChannelf              func(ctx context.Context, req rpc.LeaveChannelRequest) (*rpc.LeaveChannelResponse, error)
	MuteChannelf               func(ctx context.Context, req rpc.MuteChannelRequest) (*rpc.MuteChannelResponse, error)
	UnmuteChannelf             func(ctx context.Context, req rpc.UnmuteChannelRequest) (*rpc.UnmuteChannelResponse, error)
	MuteUserf                  func(ctx context.Context, req rpc.MuteUserRequest) (*rpc.MuteUserResponse, error)
	UnmuteUserf                func(ctx context.Context, req rpc.UnmuteUserRequest) (*rpc.UnmuteUserResponse, error)
	MuteThreadf                func(ctx context.Context, req rpc.MuteThreadRequest) (*rpc.MuteThreadResponse, error)
	UnmuteThreadf              func(ctx context.Context, req rpc.UnmuteThreadRequest) (*rpc.UnmuteThreadResponse, error)
	CreateChannelf             func(ctx context.Context, req rpc.CreateChannelRequest) (*rpc.CreateChannelResponse, error)
	GetPopularCardsf           func(ctx context.Context, req rpc.GetPopularCardsRequest) (*rpc.GetPopularCardsResponse, error)
	GetActionCostsf            func(ctx context.Context, req rpc.GetActionCostsRequest) (*rpc.GetActionCostsResponse, error)
	UseInviteCodef             func(ctx context.Context, req rpc.UseInviteCodeRequest) (*rpc.UseInviteCodeResponse, error)
	RequestValidationf         func(ctx context.Context, req rpc.RequestValidationRequest) (*rpc.RequestValidationResponse, error)
	ConfirmValidationf         func(ctx context.Context, req rpc.ConfirmValidationRequest) (*rpc.ConfirmValidationResponse, error)
	ValidateChannelNamef       func(ctx context.Context, req rpc.ValidateChannelNameRequest) (*rpc.ValidateChannelNameResponse, error)
	GetChannelf                func(ctx context.Context, req rpc.GetChannelRequest) (*rpc.GetChannelResponse, error)
	CanAffordAnonymousPostf    func(ctx context.Context, req rpc.CanAffordAnonymousPostRequest) (*rpc.CanAffordAnonymousPostResponse, error)
	GetLeaderboardf            func(ctx context.Context, req rpc.GetLeaderboardRequest) (*rpc.GetLeaderboardResponse, error)
	SubmitFeedbackf            func(ctx context.Context, req rpc.SubmitFeedbackRequest) (*rpc.SubmitFeedbackResponse, error)
	TipCardf                   func(ctx context.Context, req rpc.TipCardRequest) (*rpc.TipCardResponse, error)

	// admin panel
	ConnectUsersf func(ctx context.Context, req rpc.ConnectUsersRequest) (*rpc.ConnectUsersResponse, error)
	NewUserf      func(ctx context.Context, req rpc.NewUserRequest) (*rpc.NewUserResponse, error)
	GetUsersf     func(ctx context.Context, req rpc.GetUsersRequest) (*rpc.GetUsersResponse, error)
}

func (r *rpcMock) ValidateInviteCode(ctx context.Context, req rpc.ValidateInviteCodeRequest) (*rpc.ValidateInviteCodeResponse, error) {
	return r.ValidateInviteCodef(ctx, req)
}

func (r *rpcMock) Auth(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error) {
	return r.Authf(ctx, req)
}

func (r *rpcMock) ResetPassword(ctx context.Context, req rpc.ResetPasswordRequest) (*rpc.ResetPasswordResponse, error) {
	return r.ResetPasswordf(ctx, req)
}

func (r *rpcMock) AddToWaitlist(ctx context.Context, req rpc.AddToWaitlistRequest) (*rpc.AddToWaitlistResponse, error) {
	return r.AddToWaitlistf(ctx, req)
}

func (r *rpcMock) Logout(ctx context.Context, req rpc.LogoutRequest) (*rpc.LogoutResponse, error) {
	return r.Logoutf(ctx, req)
}

func (r *rpcMock) GetCards(ctx context.Context, req rpc.GetCardsRequest) (*rpc.GetCardsResponse, error) {
	return r.GetCardsf(ctx, req)
}

func (r *rpcMock) GetCard(ctx context.Context, req rpc.GetCardRequest) (*rpc.GetCardResponse, error) {
	return r.GetCardf(ctx, req)
}

func (r *rpcMock) ReactToCard(ctx context.Context, req rpc.ReactToCardRequest) (*rpc.ReactToCardResponse, error) {
	return r.ReactToCardf(ctx, req)
}

func (r *rpcMock) VoteOnCard(ctx context.Context, req rpc.VoteOnCardRequest) (*rpc.VoteOnCardResponse, error) {
	return r.VoteOnCardf(ctx, req)
}

func (r *rpcMock) PostCard(ctx context.Context, req rpc.PostCardRequest) (*rpc.PostCardResponse, error) {
	return r.PostCardf(ctx, req)
}

func (r *rpcMock) NewInvite(ctx context.Context, req rpc.NewInviteRequest) (*rpc.NewInviteResponse, error) {
	return r.NewInvitef(ctx, req)
}

func (r *rpcMock) RegisterDevice(ctx context.Context, req rpc.RegisterDeviceRequest) (*rpc.RegisterDeviceResponse, error) {
	return r.RegisterDevicef(ctx, req)
}

func (r *rpcMock) UnregisterDevice(ctx context.Context, req rpc.UnregisterDeviceRequest) (*rpc.UnregisterDeviceResponse, error) {
	return r.UnregisterDevicef(ctx, req)
}

func (r *rpcMock) UpdateSettings(ctx context.Context, req rpc.UpdateSettingsRequest) (*rpc.UpdateSettingsResponse, error) {
	return r.UpdateSettingsf(ctx, req)
}

func (r *rpcMock) GetUser(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
	return r.GetUserf(ctx, req)
}

func (r *rpcMock) ConnectUsers(ctx context.Context, req rpc.ConnectUsersRequest) (*rpc.ConnectUsersResponse, error) {
	return r.ConnectUsersf(ctx, req)
}

func (r *rpcMock) NewUser(ctx context.Context, req rpc.NewUserRequest) (*rpc.NewUserResponse, error) {
	return r.NewUserf(ctx, req)
}

func (r *rpcMock) GetUsers(ctx context.Context, req rpc.GetUsersRequest) (*rpc.GetUsersResponse, error) {
	return r.GetUsersf(ctx, req)
}

func (r *rpcMock) ValidateUsername(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error) {
	return r.ValidateUsernamef(ctx, req)
}

func (r *rpcMock) GetThread(ctx context.Context, req rpc.GetThreadRequest) (*rpc.GetThreadResponse, error) {
	return r.GetThreadf(ctx, req)
}

func (r *rpcMock) GetNotifications(ctx context.Context, req rpc.GetNotificationsRequest) (*rpc.GetNotificationsResponse, error) {
	return r.GetNotificationsf(ctx, req)
}

func (r *rpcMock) UpdateNotifications(ctx context.Context, req rpc.UpdateNotificationsRequest) (*rpc.UpdateNotificationsResponse, error) {
	return r.UpdateNotificationsf(ctx, req)
}

func (r *rpcMock) GetAnonymousHandle(ctx context.Context, req rpc.GetAnonymousHandleRequest) (*rpc.GetAnonymousHandleResponse, error) {
	return r.GetAnonymousHandlef(ctx, req)
}

func (r *rpcMock) DeleteCard(ctx context.Context, req rpc.DeleteCardRequest) (*rpc.DeleteCardResponse, error) {
	return r.DeleteCardf(ctx, req)
}

func (r *rpcMock) FollowUser(ctx context.Context, req rpc.FollowUserRequest) (*rpc.FollowUserResponse, error) {
	return r.FollowUserf(ctx, req)
}

func (r *rpcMock) UnfollowUser(ctx context.Context, req rpc.UnfollowUserRequest) (*rpc.UnfollowUserResponse, error) {
	return r.UnfollowUserf(ctx, req)
}

func (r *rpcMock) GetFollowingUsers(ctx context.Context, req rpc.GetFollowingUsersRequest) (*rpc.GetFollowingUsersResponse, error) {
	return r.GetFollowingUsersf(ctx, req)
}

func (r *rpcMock) GetPostsForUser(ctx context.Context, req rpc.GetPostsForUserRequest) (*rpc.GetPostsForUserResponse, error) {
	return r.GetPostsForUserf(ctx, req)
}

func (r *rpcMock) GetTags(ctx context.Context, req rpc.GetTagsRequest) (*rpc.GetTagsResponse, error) {
	return r.GetTagsf(ctx, req)
}

func (r *rpcMock) GetFeaturesForUser(ctx context.Context, req rpc.GetFeaturesForUserRequest) (*rpc.GetFeaturesForUserResponse, error) {
	return r.GetFeaturesForUserf(ctx, req)
}

func (r *rpcMock) PreviewContent(ctx context.Context, req rpc.PreviewContentRequest) (*rpc.PreviewContentResponse, error) {
	return r.PreviewContentf(ctx, req)
}

func (r *rpcMock) UploadImage(ctx context.Context, req rpc.UploadImageRequest) (*rpc.UploadImageResponse, error) {
	return r.UploadImagef(ctx, req)
}

func (r *rpcMock) GetTaggableUsers(ctx context.Context, req rpc.GetTaggableUsersRequest) (*rpc.GetTaggableUsersResponse, error) {
	return r.GetTaggableUsersf(ctx, req)
}

func (r *rpcMock) ModifyCardScore(ctx context.Context, req rpc.ModifyCardScoreRequest) (*rpc.ModifyCardScoreResponse, error) {
	return r.ModifyCardScoref(ctx, req)
}

func (r *rpcMock) GetInvites(ctx context.Context, req rpc.GetInvitesRequest) (*rpc.GetInvitesResponse, error) {
	return r.GetInvitesf(ctx, req)
}

func (r *rpcMock) GetOnboardingData(ctx context.Context, req rpc.GetOnboardingDataRequest) (*rpc.GetOnboardingDataResponse, error) {
	return r.GetOnboardingDataf(ctx, req)
}

func (r *rpcMock) GetMyNetwork(ctx context.Context, req rpc.GetMyNetworkRequest) (*rpc.GetMyNetworkResponse, error) {
	return r.GetMyNetworkf(ctx, req)
}

func (r *rpcMock) UnsubscribeFromCard(ctx context.Context, req rpc.UnsubscribeFromCardRequest) (*rpc.UnsubscribeFromCardResponse, error) {
	return r.UnsubscribeFromCardf(ctx, req)
}

func (r *rpcMock) SubscribeToCard(ctx context.Context, req rpc.SubscribeToCardRequest) (*rpc.SubscribeToCardResponse, error) {
	return r.SubscribeToCardf(ctx, req)
}

func (r *rpcMock) GroupInvites(ctx context.Context, req rpc.GroupInvitesRequest) (*rpc.GroupInvitesResponse, error) {
	return r.GroupInvitesf(ctx, req)
}

func (r *rpcMock) ReportCard(ctx context.Context, req rpc.ReportCardRequest) (*rpc.ReportCardResponse, error) {
	return r.ReportCardf(ctx, req)
}

func (r *rpcMock) BlockUser(ctx context.Context, req rpc.BlockUserRequest) (*rpc.BlockUserResponse, error) {
	return r.BlockUserf(ctx, req)
}
func (r *rpcMock) GetCardsForChannel(ctx context.Context, req rpc.GetCardsForChannelRequest) (*rpc.GetCardsForChannelResponse, error) {
	return r.GetCardsForChannelf(ctx, req)
}
func (r *rpcMock) UpdateChannelSubscription(ctx context.Context, req rpc.UpdateChannelSubscriptionRequest) (*rpc.UpdateChannelSubscriptionResponse, error) {
	return r.UpdateChannelSubscriptionf(ctx, req)
}

func (r *rpcMock) GetChannels(ctx context.Context, req rpc.GetChannelsRequest) (*rpc.GetChannelsResponse, error) {
	return r.GetChannelsf(ctx, req)
}
func (r *rpcMock) JoinChannel(ctx context.Context, req rpc.JoinChannelRequest) (*rpc.JoinChannelResponse, error) {
	return r.JoinChannelf(ctx, req)
}
func (r *rpcMock) LeaveChannel(ctx context.Context, req rpc.LeaveChannelRequest) (*rpc.LeaveChannelResponse, error) {
	return r.LeaveChannelf(ctx, req)
}
func (r *rpcMock) MuteChannel(ctx context.Context, req rpc.MuteChannelRequest) (*rpc.MuteChannelResponse, error) {
	return r.MuteChannelf(ctx, req)
}
func (r *rpcMock) UnmuteChannel(ctx context.Context, req rpc.UnmuteChannelRequest) (*rpc.UnmuteChannelResponse, error) {
	return r.UnmuteChannelf(ctx, req)
}
func (r *rpcMock) MuteUser(ctx context.Context, req rpc.MuteUserRequest) (*rpc.MuteUserResponse, error) {
	return r.MuteUserf(ctx, req)
}
func (r *rpcMock) UnmuteUser(ctx context.Context, req rpc.UnmuteUserRequest) (*rpc.UnmuteUserResponse, error) {
	return r.UnmuteUserf(ctx, req)
}
func (r *rpcMock) MuteThread(ctx context.Context, req rpc.MuteThreadRequest) (*rpc.MuteThreadResponse, error) {
	return r.MuteThreadf(ctx, req)
}
func (r *rpcMock) UnmuteThread(ctx context.Context, req rpc.UnmuteThreadRequest) (*rpc.UnmuteThreadResponse, error) {
	return r.UnmuteThreadf(ctx, req)
}

func (r *rpcMock) CreateChannel(ctx context.Context, req rpc.CreateChannelRequest) (*rpc.CreateChannelResponse, error) {
	return r.CreateChannelf(ctx, req)
}
func (r *rpcMock) GetPopularCards(ctx context.Context, req rpc.GetPopularCardsRequest) (*rpc.GetPopularCardsResponse, error) {
	return r.GetPopularCardsf(ctx, req)
}
func (r *rpcMock) ValidateChannelName(ctx context.Context, req rpc.ValidateChannelNameRequest) (*rpc.ValidateChannelNameResponse, error) {
	return r.ValidateChannelNamef(ctx, req)
}

func (r *rpcMock) GetChannel(ctx context.Context, req rpc.GetChannelRequest) (*rpc.GetChannelResponse, error) {
	return r.GetChannelf(ctx, req)
}
func (r *rpcMock) GetActionCosts(ctx context.Context, req rpc.GetActionCostsRequest) (*rpc.GetActionCostsResponse, error) {
	return r.GetActionCostsf(ctx, req)
}
func (r *rpcMock) UseInviteCode(ctx context.Context, req rpc.UseInviteCodeRequest) (*rpc.UseInviteCodeResponse, error) {
	return r.UseInviteCodef(ctx, req)
}

func (r *rpcMock) RequestValidation(ctx context.Context, req rpc.RequestValidationRequest) (*rpc.RequestValidationResponse, error) {
	return r.RequestValidationf(ctx, req)
}

func (r *rpcMock) ConfirmValidation(ctx context.Context, req rpc.ConfirmValidationRequest) (*rpc.ConfirmValidationResponse, error) {
	return r.ConfirmValidationf(ctx, req)
}
func (r *rpcMock) CanAffordAnonymousPost(ctx context.Context, req rpc.CanAffordAnonymousPostRequest) (*rpc.CanAffordAnonymousPostResponse, error) {
	return r.CanAffordAnonymousPostf(ctx, req)
}

func (r *rpcMock) GetLeaderboard(ctx context.Context, req rpc.GetLeaderboardRequest) (*rpc.GetLeaderboardResponse, error) {
	return r.GetLeaderboardf(ctx, req)
}

func (r *rpcMock) SubmitFeedback(ctx context.Context, req rpc.SubmitFeedbackRequest) (*rpc.SubmitFeedbackResponse, error) {
	return r.SubmitFeedbackf(ctx, req)
}

func (r *rpcMock) TipCard(ctx context.Context, req rpc.TipCardRequest) (*rpc.TipCardResponse, error) {
	return r.TipCardf(ctx, req)
}
