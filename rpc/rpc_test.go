package rpc_test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/store"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/worker"
	"github.com/october93/engine/worker/emailsender"
)

var ErrFailure = errors.New("failure")

func TestAuth(t *testing.T) {
	now := model.NewDBTime(time.Now().UTC())

	user := &model.User{
		ID:                      "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
		Username:                "chad",
		Email:                   "chad@october.news",
		FirstName:               "Chad",
		LastName:                "",
		DisplayName:             "Chad ",
		ProfileImagePath:        "chad.jpg",
		CoinRewardLastUpdatedAt: now,
	}

	err := user.SetPassword("secret")
	if err != nil {
		t.Fatal(err)
	}

	session := model.NewSession(nil)
	resetToken, err := model.NewResetToken(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	ErrViolateUniqueConstraint := errors.New(`"pq: duplicate key value violates unique constraint "users_username_idx"`)

	tests := []struct {
		name                         string
		params                       rpc.AuthParams
		session                      *model.Session
		getUserByUsernamef           func(username string) (*model.User, error)
		getUserByEmailf              func(email string) (*model.User, error)
		getResetTokenf               func(userID globalid.ID) (*model.ResetToken, error)
		extendTokenf                 func(ctx context.Context, token string) (rpc.AccessToken, error)
		getOAuthAccountBySubjectf    func(subject string) (*model.OAuthAccount, error)
		downloadProfileImagef        func(url string) (string, string, error)
		saveOAuthAccountf            func(oaa *model.OAuthAccount) error
		getUserf                     func(userID globalid.ID) (*model.User, error)
		saveSessionf                 func(session *model.Session) error
		getInviteByTokenf            func(token string) (*model.Invite, error)
		saveBase64ProfileImagef      func(data string) (string, string, error)
		saveUserf                    func(user *model.User) error
		generateDefaultProfileImagef func() (string, string, error)
		saveInvitef                  func(invite *model.Invite) error
		deleteWaitlistEntryf         func(email string) error
		saveNotificationf            func(notification *model.Notification) error
		exportNotificationf          func(notification *model.Notification) (*model.ExportedNotification, error)
		newNotificationf             func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error
		saveCardf                    func(card *model.Card) error
		createIndexf                 func(card *model.Card) error
		newCardf                     func(ctx context.Context, session *model.Session, card *model.CardResponse) error
		err                          error
		expected                     *rpc.AuthResponse
	}{
		{
			name: "signup",
			params: rpc.AuthParams{
				InviteToken:      "6XDAB",
				Username:         "chad",
				Password:         "secret",
				Email:            "chad@october.news",
				FirstName:        "Chad",
				LastName:         "",
				ProfileImageData: &user.ProfileImagePath,
			},
			session: session,
			getUserByUsernamef: func(username string) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
			getUserByEmailf: func(email string) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{
					Token:         token,
					RemainingUses: 1,
					NodeID:        globalid.ID("1234"),
				}, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return &model.User{}, nil
			},
			saveBase64ProfileImagef: func(data string) (string, string, error) {
				return user.ProfileImagePath, "", nil
			},
			saveUserf: func(u *model.User) error {
				u.ID = user.ID
				return nil
			},

			generateDefaultProfileImagef: func() (string, string, error) {
				return "chad.png", "", nil
			},
			saveInvitef: func(invite *model.Invite) error {
				return nil
			},
			deleteWaitlistEntryf: func(email string) error {
				return nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			saveNotificationf: func(notification *model.Notification) error {
				return nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			newNotificationf: func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
				return nil
			},
			saveCardf: func(card *model.Card) error {
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: &rpc.AuthResponse{
				User:    user.Export(user.ID),
				Session: session,
			},
		},
		{
			name: "facebook signup",
			session: &model.Session{
				ID: "6e22088e-e193-4e86-b3c8-cf544dadef29",
			},
			params: rpc.AuthParams{
				AccessToken: "secret",
				InviteToken: "6XDAB",
			},
			extendTokenf: func(ctx context.Context, token string) (rpc.AccessToken, error) {
				accessToken := &mockAccessToken{}
				accessToken.FacebookUserf = func() (*rpc.FacebookUser, error) {
					return &rpc.FacebookUser{
						ID:               "5feb7922-2acf-4fb6-9a46-ac4e177d55ca",
						Email:            "chad@october.news",
						FirstName:        "Chad",
						LastName:         "Unicorn",
						ProfileImagePath: "chad.jpg",
					}, nil
				}
				return accessToken, nil
			},
			getOAuthAccountBySubjectf: func(subject string) (*model.OAuthAccount, error) {
				return nil, sql.ErrNoRows
			},
			getUserByEmailf: func(email string) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{
					Token:         token,
					RemainingUses: 1,
					NodeID:        globalid.ID("1234"),
				}, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return &model.User{}, nil
			},
			downloadProfileImagef: func(url string) (string, string, error) {
				return "chad.jpg", "", nil
			},
			saveOAuthAccountf: func(oaa *model.OAuthAccount) error {
				return nil
			},
			saveUserf: func(u *model.User) error {
				u.ID = user.ID
				u.CoinRewardLastUpdatedAt = user.CoinRewardLastUpdatedAt
				return nil
			},
			saveInvitef: func(invite *model.Invite) error {
				return nil
			},
			deleteWaitlistEntryf: func(email string) error {
				return nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			saveNotificationf: func(notification *model.Notification) error {
				return nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			newNotificationf: func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
				return nil
			},
			saveCardf: func(card *model.Card) error {
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: &rpc.AuthResponse{
				User: &model.ExportedUser{
					ID:               "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
					Username:         "chadunicorn",
					Email:            "chad@october.news",
					FirstName:        "Chad",
					LastName:         "Unicorn",
					DisplayName:      "Chad Unicorn",
					ProfileImagePath: "chad.jpg",
				},
				Session: &model.Session{
					ID:     "6e22088e-e193-4e86-b3c8-cf544dadef29",
					UserID: "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
					User: &model.User{
						ID:                      "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
						Username:                "chadunicorn",
						Email:                   "chad@october.news",
						FirstName:               "Chad",
						LastName:                "Unicorn",
						DisplayName:             "Chad Unicorn",
						ProfileImagePath:        "chad.jpg",
						CoinRewardLastUpdatedAt: now,
					},
				},
			},
		},
		{
			name: "facebook signup overwrite username",
			session: &model.Session{
				ID: "6e22088e-e193-4e86-b3c8-cf544dadef29",
			},
			params: rpc.AuthParams{
				AccessToken: "secret",
				InviteToken: "6XDAB",
				Username:    "chad",
			},
			extendTokenf: func(ctx context.Context, token string) (rpc.AccessToken, error) {
				accessToken := &mockAccessToken{}
				accessToken.FacebookUserf = func() (*rpc.FacebookUser, error) {
					return &rpc.FacebookUser{
						ID:               "5feb7922-2acf-4fb6-9a46-ac4e177d55ca",
						Email:            "chad@october.news",
						FirstName:        "Chad",
						LastName:         "Unicorn",
						ProfileImagePath: "chad.jpg",
					}, nil
				}
				return accessToken, nil
			},
			getOAuthAccountBySubjectf: func(subject string) (*model.OAuthAccount, error) {
				return nil, sql.ErrNoRows
			},
			getUserByEmailf: func(email string) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{
					Token:         token,
					RemainingUses: 1,
					NodeID:        globalid.ID("1234"),
				}, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return &model.User{}, nil
			},
			downloadProfileImagef: func(url string) (string, string, error) {
				return "chad.jpg", "", nil
			},
			saveOAuthAccountf: func(oaa *model.OAuthAccount) error {
				return nil
			},
			saveUserf: func(u *model.User) error {
				u.ID = user.ID
				u.CoinRewardLastUpdatedAt = user.CoinRewardLastUpdatedAt
				return nil
			},
			saveInvitef: func(invite *model.Invite) error {
				return nil
			},
			deleteWaitlistEntryf: func(email string) error {
				return nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			saveNotificationf: func(notification *model.Notification) error {
				return nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			newNotificationf: func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
				return nil
			},
			saveCardf: func(card *model.Card) error {
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: &rpc.AuthResponse{
				User: &model.ExportedUser{
					ID:               "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
					Username:         "chad",
					Email:            "chad@october.news",
					FirstName:        "Chad",
					LastName:         "Unicorn",
					DisplayName:      "Chad Unicorn",
					ProfileImagePath: "chad.jpg",
				},
				Session: &model.Session{
					ID:     "6e22088e-e193-4e86-b3c8-cf544dadef29",
					UserID: "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
					User: &model.User{
						ID:                      "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
						Username:                "chad",
						Email:                   "chad@october.news",
						FirstName:               "Chad",
						LastName:                "Unicorn",
						DisplayName:             "Chad Unicorn",
						ProfileImagePath:        "chad.jpg",
						CoinRewardLastUpdatedAt: now,
					},
				},
			},
		},
		/*
			{
				name: "botched facebook signup",
				session: &model.Session{
					ID: "6e22088e-e193-4e86-b3c8-cf544dadef29",
				},
				params: rpc.AuthParams{
					AccessToken: "secret",
					InviteToken: "6XDAB",
				},
				extendTokenf: func(ctx context.Context, token string) (rpc.AccessToken, error) {
					accessToken := &mockAccessToken{}
					accessToken.FacebookUserf = func() (*rpc.FacebookUser, error) {
						return &rpc.FacebookUser{
							ID:               "5feb7922-2acf-4fb6-9a46-ac4e177d55ca",
							Email:            "chad@october.news",
							FirstName:        "Chad",
							LastName:         "Unicorn",
							ProfileImagePath: "chad.jpg",
						}, nil
					}
					return accessToken, nil
				},
				getOAuthAccountBySubjectf: func(subject string) (*model.OAuthAccount, error) {
					return nil, sql.ErrNoRows
				},
				getUserByEmailf: func(email string) (*model.User, error) {
					return nil, sql.ErrNoRows
				},
				newNodef: func() graph.Node {
					node := &mockNode{}
					node.IDf = func() globalid.ID {
						return user.ID
					}
					node.GetConfigf = func() graph.NodeConfig {
						return graph.DefaultConfig().NodeConfig
					}
					node.SetConfigf = func(config graph.NodeConfig) {}
					node.WhoWeFollowf = func() map[globalid.ID]bool {
						return nil
					}
					node.WhoFollowsUsf = func() map[globalid.ID]bool {
						return nil
					}
					return node
				},
				getInviteByTokenf: func(token string) (*model.Invite, error) {
					return &model.Invite{
						Token:         token,
						RemainingUses: 1,
					}, nil
				},
				getNodef: func(nodeID globalid.ID) graph.Node {
					node := &mockNode{}
					node.IDf = func() globalid.ID {
						return user.ID
					}
					node.WhoFollowsUsf = func() map[globalid.ID]bool {
						return nil
					}
					return node
				},
				downloadProfileImagef: func(url string) (string, string, error) {
					return "chad.jpg", "", nil
				},
				saveOAuthAccountf: func(oaa *model.OAuthAccount) error {
					return nil
				},
				saveNodef: func(node graph.Node) error {
					return nil
				},
				saveUserf: func(u *model.User) error {
					u.ID = user.ID
					return nil
				},
				saveInvitef: func(invite *model.Invite) error {
					return nil
				},
				deleteWaitlistEntryf: func(email string) error {
					return nil
				},
				saveSessionf: func(session *model.Session) error {
					return nil
				},
				saveNotificationf: func(notification *model.Notification) error {
					return nil
				},
				exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
					return &model.ExportedNotification{}, nil
				},
				newNotificationf: func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
					return nil
				},
				saveCardf: func(card *model.Card) error {
					return nil
				},
				createIndexf: func(card *model.Card) error {
					return nil
				},
				newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
					return nil
				},
				expected: &rpc.AuthResponse{
					User: &model.ExportedUser{
						ID:               "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
						Username:         "chadunicorn",
						Email:            "chad@october.news",
						FirstName:        "Chad",
						LastName:         "Unicorn",
						DisplayName:      "Chad Unicorn",
						ProfileImagePath: "chad.jpg",
						BotchedSignup:    true,
					},
					Session: &model.Session{
						ID:     "6e22088e-e193-4e86-b3c8-cf544dadef29",
						UserID: "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
						User: &model.User{
							ID:               "5a6c1446-8a9a-46f3-b3cb-200449bdb7fa",
							Username:         "chadunicorn",
							Email:            "chad@october.news",
							FirstName:        "Chad",
							LastName:         "Unicorn",
							DisplayName:      "Chad Unicorn",
							ProfileImagePath: "chad.jpg",
							BotchedSignup:    true,
						},
					},
				},
			},
		*/
		{
			name: "botched facebook signup fails",
			session: &model.Session{
				ID: "6e22088e-e193-4e86-b3c8-cf544dadef29",
			},
			params: rpc.AuthParams{
				AccessToken: "secret",
				InviteToken: "6XDAB",
			},
			extendTokenf: func(ctx context.Context, token string) (rpc.AccessToken, error) {
				accessToken := &mockAccessToken{}
				accessToken.FacebookUserf = func() (*rpc.FacebookUser, error) {
					return &rpc.FacebookUser{
						ID:               "5feb7922-2acf-4fb6-9a46-ac4e177d55ca",
						Email:            "chad@october.news",
						FirstName:        "Chad",
						LastName:         "Unicorn",
						ProfileImagePath: "chad.jpg",
					}, nil
				}
				return accessToken, nil
			},
			getOAuthAccountBySubjectf: func(subject string) (*model.OAuthAccount, error) {
				return nil, sql.ErrNoRows
			},
			getUserByEmailf: func(email string) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{
					Token:         token,
					RemainingUses: 1,
					NodeID:        globalid.ID("1234"),
				}, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return &model.User{}, nil
			},
			downloadProfileImagef: func(url string) (string, string, error) {
				return "chad.jpg", "", nil
			},
			saveOAuthAccountf: func(oaa *model.OAuthAccount) error {
				return nil
			},
			saveUserf: func(u *model.User) error {
				return ErrViolateUniqueConstraint
			},
			saveInvitef: func(invite *model.Invite) error {
				return nil
			},
			deleteWaitlistEntryf: func(email string) error {
				return nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			saveNotificationf: func(notification *model.Notification) error {
				return nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			newNotificationf: func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
				return nil
			},
			saveCardf: func(card *model.Card) error {
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: nil,
			err:      ErrViolateUniqueConstraint,
		},
		{
			name:    "username login",
			session: session,
			params: rpc.AuthParams{
				Username: user.Username,
				Password: "secret",
			},
			getUserByUsernamef: func(username string) (*model.User, error) {
				return user, nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			expected: &rpc.AuthResponse{
				User:    user.Export(user.ID),
				Session: session,
			},
		},
		{
			name:    "reset password login",
			session: session,
			params: rpc.AuthParams{
				ResetToken: resetToken.Token.String(),
			},
			getUserByEmailf: func(email string) (*model.User, error) {
				return user, nil
			},
			getResetTokenf: func(userID globalid.ID) (*model.ResetToken, error) {
				return resetToken, nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			expected: &rpc.AuthResponse{
				User:    user.Export(user.ID),
				Session: session,
			},
		},
		{
			name:    "access token login",
			session: session,
			params: rpc.AuthParams{
				AccessToken: "secret",
			},
			extendTokenf: func(ctx context.Context, token string) (rpc.AccessToken, error) {
				accessToken := &mockAccessToken{}
				accessToken.FacebookUserf = func() (*rpc.FacebookUser, error) {
					return &rpc.FacebookUser{
						ID:               "5feb7922-2acf-4fb6-9a46-ac4e177d55ca",
						Email:            "chad@october.news",
						FirstName:        "Chad",
						LastName:         "Unicorn",
						ProfileImagePath: "chad.jpg",
					}, nil
				}
				return accessToken, nil
			},
			getOAuthAccountBySubjectf: func(subject string) (*model.OAuthAccount, error) {
				return &model.OAuthAccount{
					UserID: user.ID,
				}, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				if userID != user.ID {
					t.Errorf("expected userID %v, actual %v", user.ID, userID)
				}
				return user, nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			expected: &rpc.AuthResponse{
				User:    user.Export(user.ID),
				Session: session,
			},
		},
		{
			name: "wrong password",
			params: rpc.AuthParams{
				Username: user.Username,
				Password: "123",
			},
			getUserByUsernamef: func(username string) (*model.User, error) {
				return user, nil
			},
			saveSessionf: func(session *model.Session) error {
				return nil
			},
			err: rpc.ErrWrongPassword,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserByUsernamef = tt.getUserByUsernamef
			r.store.GetUserByEmailf = tt.getUserByEmailf
			r.store.GetResetTokenf = tt.getResetTokenf
			r.store.GetOAuthAccountBySubjectf = tt.getOAuthAccountBySubjectf
			r.store.GetUserf = tt.getUserf
			r.store.GetInviteByTokenf = tt.getInviteByTokenf
			r.store.SaveUserf = tt.saveUserf
			r.store.SaveInvitef = tt.saveInvitef
			r.store.DeleteWaitlistEntryf = tt.deleteWaitlistEntryf
			r.store.SaveSessionf = tt.saveSessionf
			r.store.SaveNotificationf = tt.saveNotificationf
			r.store.SaveOAuthAccountf = tt.saveOAuthAccountf
			r.store.SaveCardf = tt.saveCardf
			r.pusher.NewNotificationf = tt.newNotificationf
			r.pusher.NewCardf = tt.newCardf
			r.oauth2.ExtendTokenf = tt.extendTokenf
			r.imageProcessor.SaveBase64ProfileImagef = tt.saveBase64ProfileImagef
			r.imageProcessor.DownloadProfileImagef = tt.downloadProfileImagef
			r.imageProcessor.GenerateDefaultProfileImagef = tt.generateDefaultProfileImagef
			r.notifications.ExportNotificationf = tt.exportNotificationf

			req := rpc.AuthRequest{
				Session: tt.session,
				Params:  tt.params,
			}
			resp, err := r.Auth(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestResetPassword(t *testing.T) {
	tests := []struct {
		name            string
		email           string
		getUserByEmailf func(email string) (*model.User, error)
		saveResetTokenf func(resetToken *model.ResetToken) error
		enqueueMailJobf func(job *emailsender.Job) error
		err             error
	}{
		{
			name:  "valid",
			email: "chad@october.news",
			getUserByEmailf: func(email string) (*model.User, error) {
				return model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn"), nil
			},
			saveResetTokenf: func(resetToken *model.ResetToken) error {
				return nil
			},
			enqueueMailJobf: func(job *emailsender.Job) error {
				if job.To != "chad@october.news" {
					t.Errorf("expected recipient to be set to %s, actual %s", "chad@october.news", job.Recipient)
				}
				return nil
			},
		},
		{
			name: "unknown email",
			getUserByEmailf: func(email string) (*model.User, error) {
				return nil, sql.ErrNoRows
			},
		},
		{
			name: "query fails",
			getUserByEmailf: func(email string) (*model.User, error) {
				return nil, ErrFailure
			},
			err: ErrFailure,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			r := newRPC(t)
			r.store.GetUserByEmailf = tt.getUserByEmailf
			r.store.SaveResetTokenf = tt.saveResetTokenf
			r.worker.EnqueueMailJobf = tt.enqueueMailJobf

			req := rpc.ResetPasswordRequest{
				Params: rpc.ResetPasswordParams{
					Email: tt.email,
				},
			}
			_, err := r.ResetPassword(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
		})
	}
}

func TestValidateInviteCode(t *testing.T) {
	tests := []struct {
		name              string
		getInviteByTokenf func(token string) (*model.Invite, error)
		err               error
	}{
		{
			name: "valid",
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				return model.NewInvite(globalid.Next())
			},
		},
		{
			name: "expired token",
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				invite, err := model.NewInvite(globalid.Next())
				return invite, err
			},
			err: nil,
		},
		/*{
			name: "no remaining uses left",
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				invite, err := model.NewInvite(globalid.Next())
				invite.RemainingUses = 0
				return invite, err
			},
			err: model.ErrInvalidInviteCode,
		},
		{
			name: "negative remaining uses left",
			getInviteByTokenf: func(token string) (*model.Invite, error) {
				invite, err := model.NewInvite(globalid.Next())
				invite.RemainingUses = -1
				return invite, err
			},
			err: model.ErrInvalidInviteCode,
		},*/
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			r := newRPC(t)
			r.store.GetInviteByTokenf = tt.getInviteByTokenf

			req := rpc.ValidateInviteCodeRequest{}
			_, err := r.ValidateInviteCode(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
		})
	}
}

func TestAddToWaitlist(t *testing.T) {
	expiresAt := time.Now().Add(90 * 24 * time.Hour).Unix()

	tests := []struct {
		name               string
		saveWaitlistEntryf func(waitlistEntry *model.WaitlistEntry) error
		extendTokenf       func(ctx context.Context, token string) (rpc.AccessToken, error)
		params             rpc.AddToWaitlistParams
		err                error
		expected           *rpc.AddToWaitlistResponse
	}{
		{
			name: "valid",
			params: rpc.AddToWaitlistParams{
				Name:  "Chad Unicorn",
				Email: "chad@october.news",
			},
			saveWaitlistEntryf: func(waitlistEntry *model.WaitlistEntry) error {
				if waitlistEntry.Email != "chad@october.news" {
					t.Errorf("expected waitlist entry to have email set to %s, actual %s", "chad@october.news", waitlistEntry.Email)
				}
				if waitlistEntry.Name != "Chad Unicorn" {
					t.Errorf("expected waitlist entry to have name set to %s, actual %s", "Chad Unicorn", waitlistEntry.Name)
				}
				return nil
			},
			expected: &rpc.AddToWaitlistResponse{},
		},
		{
			name: "access token",
			params: rpc.AddToWaitlistParams{
				AccessToken: "123=",
			},
			extendTokenf: func(ctx context.Context, token string) (rpc.AccessToken, error) {
				accessToken := &mockAccessToken{}
				accessToken.FacebookUserf = func() (*rpc.FacebookUser, error) {
					return &rpc.FacebookUser{
						Email:     "richard@piedpiper.com",
						FirstName: "Richard",
						LastName:  "Hendricks",
					}, nil
				}
				accessToken.Tokenf = func() string {
					return "12345="
				}
				accessToken.ExpiresAtf = func() int64 {
					return expiresAt
				}
				return accessToken, nil
			},
			saveWaitlistEntryf: func(waitlistEntry *model.WaitlistEntry) error {
				if waitlistEntry.Email != "richard@piedpiper.com" {
					t.Errorf("expected waitlist entry to have email set to %s, actual %s", "richard@piedpiper.com", waitlistEntry.Email)
				}
				if waitlistEntry.Name != "Richard Hendricks" {
					t.Errorf("expected waitlist entry to have name set to %s, actual %s", "Richard Hendricks", waitlistEntry.Name)
				}
				return nil
			},
			expected: &rpc.AddToWaitlistResponse{
				AccessToken:          "12345=",
				AccessTokenExpiresAt: expiresAt,
			},
		},
		{
			name: "failure",
			saveWaitlistEntryf: func(waitlistEntry *model.WaitlistEntry) error {
				return ErrFailure
			},
			err: ErrFailure,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			r := newRPC(t)
			r.store.SaveWaitlistEntryf = tt.saveWaitlistEntryf
			r.oauth2.ExtendTokenf = tt.extendTokenf

			req := rpc.AddToWaitlistRequest{Params: tt.params}
			resp, err := r.AddToWaitlist(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(resp, tt.expected) {
				t.Fatalf("expected a different response: %v", pretty.Diff(resp, tt.expected))
			}
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name           string
		user           *model.User
		deleteSessionf func(id globalid.ID) error
		err            error
	}{
		{
			name: "valid",
			user: model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn"),
			deleteSessionf: func(id globalid.ID) error {
				return nil
			},
		},
		{
			name: "unknown session id",
			user: model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn"),
			deleteSessionf: func(id globalid.ID) error {
				return sql.ErrNoRows
			},
			err: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			req := rpc.LogoutRequest{
				Session: model.NewSession(tt.user),
			}
			r.store.DeleteSessionf = tt.deleteSessionf
			_, err := r.Logout(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if tt.err != nil {
				if req.Session.User == nil {
					t.Fatal("expected Logout not to remove user from session")
				}
			} else {
				if req.Session.User != nil {
					t.Fatalf("expected Logout to removed user from session, actual: %v", req.Session.User)
				}
			}
		})
	}
}

func TestGetThread(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	erlich := model.NewUser("eb0b382b-e58d-4653-9dbc-e7f47aca65d5", "erlich", "erlich@october.news", "Erlich Bachman")
	card := model.NewCard(user.ID, "Title", "Content")
	reaction := &model.UserReaction{
		CardID: card.ID,
		UserID: user.ID,
		Type:   model.ReactionLike,
	}

	tests := []struct {
		name             string
		user             *model.User
		getThreadf       func(id, forUser globalid.ID) ([]*model.Card, error)
		getUserf         func(userID globalid.ID) (*model.User, error)
		getThreadCountf  func(cardID globalid.ID) (int, error)
		getUserReactionf func(userID, cardID globalid.ID) (*model.UserReaction, error)
		getEngagementf   func(cardID globalid.ID) (*model.Engagement, error)
		err              error
		expected         []*model.CardResponse
	}{
		{
			name: "valid",
			user: user,
			getThreadf: func(id, forUser globalid.ID) ([]*model.Card, error) {
				return []*model.Card{&card}, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 1, nil
			},
			getUserReactionf: func(userID, cardID globalid.ID) (*model.UserReaction, error) {
				return reaction, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return &model.Engagement{
					EngagedUsers: []*model.EngagedUser{
						&model.EngagedUser{
							UserID: reaction.UserID,
						},
					},
					Count: 1,
				}, nil
			},
			expected: []*model.CardResponse{
				&model.CardResponse{
					Author: user.Author(),
					Card:   card.Export(),
					Vote:   &model.VoteResponse{Type: "up"},
					Reactions: &model.Reaction{
						ID:        "33c79287-be80-4fbd-a605-90a90ca3253c",
						NodeID:    "33c79287-be80-4fbd-a605-90a90ca3253c",
						AliasID:   "",
						CardID:    reaction.CardID,
						Reaction:  "boost",
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					},
					Engagement: &model.Engagement{
						EngagedUsers: []*model.EngagedUser{
							&model.EngagedUser{
								UserID: reaction.UserID,
							},
						},
						Count: 1,
					},
					ViewerReaction: reaction,
					Replies:        1,
					IsMine:         true,
				},
			},
		},
		{
			name: "cards not authored by user",
			user: erlich,
			getThreadf: func(id, forUser globalid.ID) ([]*model.Card, error) {
				return []*model.Card{&card}, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 1, nil
			},
			getUserReactionf: func(userID, cardID globalid.ID) (*model.UserReaction, error) {
				return nil, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return &model.Engagement{
					EngagedUsers: []*model.EngagedUser{
						&model.EngagedUser{
							UserID: reaction.UserID,
						},
					},
					Count: 1,
				}, nil
			},
			expected: []*model.CardResponse{
				&model.CardResponse{
					Author: user.Author(),
					Card:   card.Export(),
					Engagement: &model.Engagement{
						EngagedUsers: []*model.EngagedUser{
							&model.EngagedUser{
								UserID: reaction.UserID,
							},
						},
						Count: 1,
					},
					Replies: 1,
				},
			},
		},
		{
			name: "empty thread",
			user: user,
			getThreadf: func(id, forUser globalid.ID) ([]*model.Card, error) {
				return []*model.Card{}, nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getUserReactionf: func(userID, cardID globalid.ID) (*model.UserReaction, error) {
				return nil, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			expected: []*model.CardResponse{},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetThreadf = tt.getThreadf
			r.store.GetUserf = tt.getUserf
			r.store.GetThreadCountf = tt.getThreadCountf
			r.store.GetUserReactionf = tt.getUserReactionf
			r.store.GetEngagementf = tt.getEngagementf

			r.store.SubscribedToTypesf = func(userID, cardID globalid.ID) ([]string, error) {
				return nil, nil
			}

			req := rpc.GetThreadRequest{
				Session: model.NewSession(tt.user),
			}
			resp, err := r.GetThread(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if resp == nil && tt.expected != nil {
				t.Fatalf("expected a result, actual %v", tt.expected)
			}
			cards := ([]*model.CardResponse)(*resp)
			if len(cards) != len(tt.expected) {
				t.Fatalf("expected response to be of length %d, actual %d", len(tt.expected), len(cards))
			}
			if !reflect.DeepEqual(cards, tt.expected) {
				t.Fatalf("expected a different response: %v", pretty.Diff(cards, tt.expected))
			}
		})
	}
}

func TestGetCards(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	card := model.NewCard(user.ID, "Title", "Content")
	reaction := &model.UserReaction{
		CardID: card.ID,
		UserID: user.ID,
		Type:   model.ReactionLike,
	}

	tests := []struct {
		name                        string
		params                      rpc.GetCardsParams
		user                        *model.User
		getUserf                    func(userID globalid.ID) (*model.User, error)
		getThreadCountf             func(cardID globalid.ID) (int, error)
		getUserReactionf            func(userID, cardID globalid.ID) (*model.UserReaction, error)
		cacheFeedForUserf           func(userID globalid.ID, ids []globalid.ID) error
		getEngagementf              func(cardID globalid.ID) (*model.Engagement, error)
		getFeedCardsFromCurrentTopf func(userID globalid.ID, perPage, page int) ([]*model.Card, error)
		feedCardResponsesf          func(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error)
		err                         error
		expected                    *rpc.GetCardsResponse
	}{
		{
			name:   "first page",
			params: rpc.GetCardsParams{Page: 0},
			user:   user,
			getFeedCardsFromCurrentTopf: func(userID globalid.ID, perPage, page int) ([]*model.Card, error) {
				return []*model.Card{&card}, nil
			},
			getThreadCountf: func(globalid.ID) (int, error) {
				return 0, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			getUserReactionf: func(userID, CardID globalid.ID) (*model.UserReaction, error) {
				return reaction, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			cacheFeedForUserf: func(userID globalid.ID, ids []globalid.ID) error {
				return nil
			},
			feedCardResponsesf: func(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error) {
				return []*model.CardResponse{
					&model.CardResponse{
						Author: user.Author(),
						Card:   card.Export(),
						Vote:   &model.VoteResponse{Type: "up"},
						Reactions: &model.Reaction{
							ID:        "33c79287-be80-4fbd-a605-90a90ca3253c",
							NodeID:    "33c79287-be80-4fbd-a605-90a90ca3253c",
							AliasID:   "",
							CardID:    reaction.CardID,
							Reaction:  "boost",
							CreatedAt: time.Time{},
							UpdatedAt: time.Time{},
						},
						IsMine:         true,
						ViewerReaction: reaction,
						Replies:        0,
					},
				}, nil
			},
			expected: &rpc.GetCardsResponse{
				NextPage: true,
				Cards: []*model.CardResponse{
					&model.CardResponse{
						Author: user.Author(),
						Card:   card.Export(),
						Vote:   &model.VoteResponse{Type: "up"},
						Reactions: &model.Reaction{
							ID:        "33c79287-be80-4fbd-a605-90a90ca3253c",
							NodeID:    "33c79287-be80-4fbd-a605-90a90ca3253c",
							AliasID:   "",
							CardID:    reaction.CardID,
							Reaction:  "boost",
							CreatedAt: time.Time{},
							UpdatedAt: time.Time{},
						},
						IsMine:         true,
						ViewerReaction: reaction,
						Replies:        0,
					},
				},
			},
		},
		{
			name:   "second page",
			params: rpc.GetCardsParams{Page: 1},
			user:   user,
			getFeedCardsFromCurrentTopf: func(userID globalid.ID, perPage, page int) ([]*model.Card, error) {
				return []*model.Card{&card}, nil
			},
			getThreadCountf: func(globalid.ID) (int, error) {
				return 0, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			getUserReactionf: func(userID, CardID globalid.ID) (*model.UserReaction, error) {
				return reaction, nil
			},
			cacheFeedForUserf: func(userID globalid.ID, ids []globalid.ID) error {
				return nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			feedCardResponsesf: func(cards []*model.Card, viewerID globalid.ID) ([]*model.CardResponse, error) {
				return []*model.CardResponse{
					&model.CardResponse{
						Author: user.Author(),
						Card:   card.Export(),
						Vote:   &model.VoteResponse{Type: "up"},
						Reactions: &model.Reaction{
							ID:        "33c79287-be80-4fbd-a605-90a90ca3253c",
							NodeID:    "33c79287-be80-4fbd-a605-90a90ca3253c",
							AliasID:   "",
							CardID:    reaction.CardID,
							Reaction:  "boost",
							CreatedAt: time.Time{},
							UpdatedAt: time.Time{},
						},
						ViewerReaction: reaction,
						Replies:        0,
						IsMine:         true,
					},
				}, nil
			},
			expected: &rpc.GetCardsResponse{
				NextPage: true,
				Cards: []*model.CardResponse{
					&model.CardResponse{
						Author: user.Author(),
						Card:   card.Export(),
						Vote:   &model.VoteResponse{Type: "up"},
						Reactions: &model.Reaction{
							ID:        "33c79287-be80-4fbd-a605-90a90ca3253c",
							NodeID:    "33c79287-be80-4fbd-a605-90a90ca3253c",
							AliasID:   "",
							CardID:    reaction.CardID,
							Reaction:  "boost",
							CreatedAt: time.Time{},
							UpdatedAt: time.Time{},
						},
						ViewerReaction: reaction,
						Replies:        0,
						IsMine:         true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserf = tt.getUserf
			r.store.GetThreadCountf = tt.getThreadCountf
			r.store.GetUserReactionf = tt.getUserReactionf
			r.store.GetEngagementf = tt.getEngagementf
			r.store.GetFeedCardsFromCurrentTopf = tt.getFeedCardsFromCurrentTopf

			r.store.SubscribedToTypesf = func(userID, cardID globalid.ID) ([]string, error) {
				return nil, nil
			}
			r.responses.FeedCardResponsesf = tt.feedCardResponsesf

			req := rpc.GetCardsRequest{
				Session: model.NewSession(tt.user),
				Params:  tt.params,
			}
			resp, err := r.GetCards(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if resp == nil && tt.expected != nil {
				t.Fatalf("expected a result, actual %v", tt.expected)
			}
			cards := resp
			if len(cards.Cards) != len(tt.expected.Cards) {
				t.Fatalf("expected response to be of length %d, actual %d", len(tt.expected.Cards), len(cards.Cards))
			}
			if !reflect.DeepEqual(cards, tt.expected) {
				t.Fatalf("expected a different response: %v", pretty.Diff(cards, tt.expected))
			}
		})
	}
}

func TestGetCard(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	alias := &model.AnonymousAlias{
		ID:               "40c5a86a-ab1a-4174-9096-afdc40dd862b",
		Username:         "egg",
		DisplayName:      "Anonymous",
		ProfileImagePath: "egg.png",
	}
	card := &model.Card{
		ID:      "ce08c980-df23-48c6-8434-eac3c81fc845",
		OwnerID: user.ID,
		AuthorToAlias: model.IdentityMap{
			user.ID: alias.ID,
		},
	}
	reaction := &model.UserReaction{
		CardID: card.ID,
		UserID: user.ID,
		Type:   model.ReactionLike,
	}

	tests := []struct {
		name                       string
		user                       *model.User
		params                     rpc.GetCardParams
		getCardf                   func(cardID globalid.ID) (*model.Card, error)
		getUserf                   func(userID globalid.ID) (*model.User, error)
		getThreadCountf            func(cardID globalid.ID) (int, error)
		getUserReactionf           func(userID, cardID globalid.ID) (*model.UserReaction, error)
		getEngagementf             func(cardID globalid.ID) (*model.Engagement, error)
		getAnonymousAliasf         func(aliasID globalid.ID) (*model.AnonymousAlias, error)
		getAnonymousAliasLastUsedf func(userID, threadRootID globalid.ID) (bool, error)
		err                        error
		expected                   *model.CardResponse
	}{
		{
			name:   "valid",
			user:   user,
			params: rpc.GetCardParams{CardID: "ce08c980-df23-48c6-8434-eac3c81fc845"},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				return card, nil
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getUserReactionf: func(userID, CardID globalid.ID) (*model.UserReaction, error) {
				return reaction, nil
			},
			getAnonymousAliasLastUsedf: func(userID, threadRootID globalid.ID) (bool, error) {
				return true, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			getAnonymousAliasf: func(aliasID globalid.ID) (*model.AnonymousAlias, error) {
				if aliasID != alias.ID {
					t.Errorf("expected aliasID: %v, actual %v", alias.ID, aliasID)
				}
				return alias, nil
			},
			expected: &model.CardResponse{
				Card:   card.Export(),
				Author: user.Author(),
				Vote:   &model.VoteResponse{Type: "up"},
				Reactions: &model.Reaction{
					ID:        "33c79287-be80-4fbd-a605-90a90ca3253c",
					NodeID:    "33c79287-be80-4fbd-a605-90a90ca3253c",
					AliasID:   "",
					CardID:    "ce08c980-df23-48c6-8434-eac3c81fc845",
					Reaction:  "boost",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Viewer: &model.Viewer{
					AnonymousAlias:         alias,
					AnonymousAliasLastUsed: true,
				},
				ViewerReaction: reaction,
				Replies:        0,
				IsMine:         true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			r := newRPC(t)
			r.store.GetCardf = tt.getCardf
			r.store.GetUserf = tt.getUserf
			r.store.GetThreadCountf = tt.getThreadCountf
			r.store.GetUserReactionf = tt.getUserReactionf
			r.store.GetEngagementf = tt.getEngagementf
			r.store.GetAnonymousAliasf = tt.getAnonymousAliasf
			r.store.LastPostInThreadWasAnonymousf = tt.getAnonymousAliasLastUsedf

			r.store.SubscribedToTypesf = func(userID, cardID globalid.ID) ([]string, error) {
				return nil, nil
			}

			req := rpc.GetCardRequest{
				Session: model.NewSession(tt.user),
				Params:  tt.params,
			}
			resp, err := r.GetCard(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			card := (*model.CardResponse)(resp)
			if !reflect.DeepEqual(card, tt.expected) {
				t.Fatalf("expected a different response: %v", pretty.Diff(card, tt.expected))
			}
		})
	}
}

func TestReactToCard(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	card := &model.Card{
		ID:      "ce08c980-df23-48c6-8434-eac3c81fc845",
		OwnerID: "990fd10f-8e3b-48c4-8a58-f862773d0352",
	}

	tests := []struct {
		name                     string
		params                   rpc.ReactToCardParams
		getUserf                 func(userID globalid.ID) (*model.User, error)
		countGraphReactionf      func(userID, cardID globalid.ID) (int, error)
		getCardf                 func(cardID globalid.ID) (*model.Card, error)
		getUnusedAliasf          func(cardID globalid.ID) (*model.AnonymousAlias, error)
		saveCardf                func(c *model.Card) error
		getAnonymousAliasf       func(id globalid.ID) (*model.AnonymousAlias, error)
		latestForTypef           func(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error)
		exportNotificationf      func(notification *model.Notification) (*model.ExportedNotification, error)
		updateNotificationf      func(ctx context.Context, session *model.Session, notification *model.ExportedNotification) error
		newNotificationf         func(ctx context.Context, session *model.Session, notification *model.ExportedNotification) error
		updateEngagementf        func(ctx context.Context, session *model.Session, cardID globalid.ID) error
		clearEmptyNotificationsf func() error
		err                      error
		expected                 *rpc.ReactToCardResponse
	}{
		{
			name: "boost",
			params: rpc.ReactToCardParams{
				CardID:   card.ID,
				Reaction: model.Boost,
				Strength: 1.0,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			countGraphReactionf: func(userID, cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				return card, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			saveCardf: func(c *model.Card) error {
				t.Error("unexpected call")
				return nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			latestForTypef: func(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error) {
				return &model.Notification{}, nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			updateNotificationf: func(ctx context.Context, session *model.Session, notification *model.ExportedNotification) error {
				return nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
			expected: &rpc.ReactToCardResponse{},
		},
		{
			name: "boost own card",
			params: rpc.ReactToCardParams{
				CardID:   "6ebffcc0-bf5c-437c-8468-af7eb7ecf03e",
				Reaction: model.Boost,
				Strength: 1.0,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			countGraphReactionf: func(userID, cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				return &model.Card{
					ID:      "6ebffcc0-bf5c-437c-8468-af7eb7ecf03e",
					OwnerID: user.ID,
				}, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			saveCardf: func(c *model.Card) error {
				t.Error("unexpected call")
				return nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
			expected: &rpc.ReactToCardResponse{},
		},
		{
			name: "boost anonymously",
			params: rpc.ReactToCardParams{
				CardID:    "6ebffcc0-bf5c-437c-8468-af7eb7ecf03e",
				Reaction:  model.Boost,
				Strength:  1.0,
				Anonymous: true,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			countGraphReactionf: func(userID, cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				return &model.Card{
					ID:            "6ebffcc0-bf5c-437c-8468-af7eb7ecf03e",
					OwnerID:       user.ID,
					AuthorToAlias: make(model.IdentityMap),
				}, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				return &model.AnonymousAlias{
					ID:          "60807b6f-305e-4da2-a5fa-e4e5be1e9bd7",
					DisplayName: "Anonymous",
					Username:    "mouse",
				}, nil
			},
			saveCardf: func(c *model.Card) error {
				return nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
			expected: &rpc.ReactToCardResponse{
				AnonymousAlias: &model.AnonymousAlias{
					ID:          "60807b6f-305e-4da2-a5fa-e4e5be1e9bd7",
					DisplayName: "Anonymous",
					Username:    "mouse",
				},
			},
		},
		{
			name: "boost anonymously twice",
			params: rpc.ReactToCardParams{
				CardID:    "6ebffcc0-bf5c-437c-8468-af7eb7ecf03e",
				Reaction:  model.Boost,
				Strength:  1.0,
				Anonymous: true,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			countGraphReactionf: func(userID, cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				return &model.Card{
					ID:      "6ebffcc0-bf5c-437c-8468-af7eb7ecf03e",
					OwnerID: user.ID,
					AuthorToAlias: model.IdentityMap{
						user.ID: "60807b6f-305e-4da2-a5fa-e4e5be1e9bd7",
					},
				}, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			saveCardf: func(c *model.Card) error {
				t.Error("unexpected call")
				return nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				return &model.AnonymousAlias{
					ID:          "60807b6f-305e-4da2-a5fa-e4e5be1e9bd7",
					DisplayName: "Anonymous",
					Username:    "mouse",
				}, nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
			expected: &rpc.ReactToCardResponse{
				AnonymousAlias: &model.AnonymousAlias{
					ID:          "60807b6f-305e-4da2-a5fa-e4e5be1e9bd7",
					DisplayName: "Anonymous",
					Username:    "mouse",
				},
			},
		},
		{
			name: "undo reaction",
			params: rpc.ReactToCardParams{
				CardID:   card.ID,
				Reaction: model.Boost,
				Strength: 1.0,
				Undo:     true,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			countGraphReactionf: func(userID, cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				return card, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			saveCardf: func(c *model.Card) error {
				t.Error("unexpected call")
				return nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				t.Error("unexpected call")
				return nil, nil
			},
			latestForTypef: func(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error) {
				return &model.Notification{}, nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			updateNotificationf: func(ctx context.Context, session *model.Session, notification *model.ExportedNotification) error {
				return nil
			},
			newNotificationf: func(ctx context.Context, session *model.Session, notification *model.ExportedNotification) error {
				return nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
			clearEmptyNotificationsf: func() error {
				return nil
			},
			expected: &rpc.ReactToCardResponse{},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserf = tt.getUserf
			r.store.GetCardf = tt.getCardf
			r.store.GetUnusedAliasf = tt.getUnusedAliasf
			r.store.SaveCardf = tt.saveCardf
			r.store.GetAnonymousAliasf = tt.getAnonymousAliasf
			r.store.LatestForTypef = tt.latestForTypef
			r.store.ClearEmptyNotificationsf = tt.clearEmptyNotificationsf
			r.pusher.UpdateNotificationf = tt.updateNotificationf
			r.pusher.NewNotificationf = tt.updateNotificationf
			r.pusher.UpdateEngagementf = tt.updateEngagementf
			r.notifications.ExportNotificationf = tt.exportNotificationf

			r.store.SaveNotificationf = func(m *model.Notification) error {
				return nil
			}

			r.store.SubscribeToCardf = func(userID, cardID globalid.ID, typ string) error {
				return nil
			}

			r.store.UnsubscribeFromCardf = func(userID, cardID globalid.ID, typ string) error {
				return nil
			}

			r.store.UpdateNotificationsOpenedf = func(ids []globalid.ID) error {
				return nil
			}

			r.store.ClearEmptyNotificationsf = func() error {
				return nil
			}

			req := rpc.ReactToCardRequest{
				Params:  tt.params,
				Session: model.NewSession(user),
			}
			resp, err := r.ReactToCard(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(resp, tt.expected) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(resp, tt.expected))
			}
		})
	}
}

func TestPostCard(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	alias1 := &model.AnonymousAlias{
		ID:          "ca51f608-31ae-43af-8684-644a80c50c50",
		DisplayName: "Anonymous",
		Username:    "egg",
	}
	alias2 := &model.AnonymousAlias{
		ID:               "5858f233-b6e1-40f5-a0ac-a73be339f1d6",
		Username:         "mouse",
		DisplayName:      "Anonymous",
		ProfileImagePath: "mouse.png",
	}

	tests := []struct {
		name                       string
		params                     rpc.PostCardParams
		getUserf                   func(userID globalid.ID) (*model.User, error)
		getUnusedAliasf            func(cardID globalid.ID) (*model.AnonymousAlias, error)
		getAnonymousAliasf         func(id globalid.ID) (*model.AnonymousAlias, error)
		notifiableForCardf         func(cardID globalid.ID) ([]globalid.ID, error)
		getCardf                   func(cardID globalid.ID) (*model.Card, error)
		getUserByUsernamef         func(username string) (*model.User, error)
		saveNotificationf          func(notification *model.Notification) error
		exportNotificationf        func(notification *model.Notification) (*model.ExportedNotification, error)
		newNotificationf           func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error
		saveCardf                  func(card *model.Card) error
		createIndexf               func(card *model.Card) error
		getThreadCountf            func(cardID globalid.ID) (int, error)
		getEngagementf             func(cardID globalid.ID) (*model.Engagement, error)
		getAnonymousAliasLastUsedf func(userID, threadRootID globalid.ID) (bool, error)
		newCardf                   func(ctx context.Context, session *model.Session, card *model.CardResponse) error
		updateEngagementf          func(ctx context.Context, session *model.Session, cardID globalid.ID) error

		err      error
		expected *rpc.PostCardResponse
	}{
		{
			name: "new post",
			params: rpc.PostCardParams{
				Content: "Lorem ipsum",
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn"), nil
			},
			notifiableForCardf: func(cardID globalid.ID) ([]globalid.ID, error) {
				return nil, nil
			},
			saveCardf: func(card *model.Card) error {
				card.CreatedAt = time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC)
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			getAnonymousAliasLastUsedf: func(userID, threadRootID globalid.ID) (bool, error) {
				return false, nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: &rpc.PostCardResponse{
				Author: user.Author(),
				Card: &model.CardView{
					Content:   "Lorem ipsum",
					CreatedAt: time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		{
			name: "comment",
			params: rpc.PostCardParams{
				Content:     "Lorem ipsum",
				ReplyCardID: "87bfe749-00a4-4844-ba61-7665f4c2cbb3",
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn"), nil
			},
			notifiableForCardf: func(cardID globalid.ID) ([]globalid.ID, error) {
				return nil, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				expected := globalid.ID("87bfe749-00a4-4844-ba61-7665f4c2cbb3")
				if cardID != expected {
					t.Errorf("expected cardID: %v, actual %v", expected, cardID)
				}
				return &model.Card{
					ID: expected,
				}, nil
			},
			saveCardf: func(card *model.Card) error {
				card.CreatedAt = time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC)
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			getAnonymousAliasLastUsedf: func(userID, threadRootID globalid.ID) (bool, error) {
				return false, nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
			expected: &rpc.PostCardResponse{
				Author: user.Author(),
				Card: &model.CardView{
					Content:       "Lorem ipsum",
					ThreadRootID:  "87bfe749-00a4-4844-ba61-7665f4c2cbb3",
					ThreadReplyID: "87bfe749-00a4-4844-ba61-7665f4c2cbb3",
					CreatedAt:     time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		{
			name: "post anonymously",
			params: rpc.PostCardParams{
				Content:   "Lorem ipsum",
				Anonymous: true,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn"), nil
			},
			notifiableForCardf: func(cardID globalid.ID) ([]globalid.ID, error) {
				return nil, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				return alias1, nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				if id != alias1.ID {
					t.Errorf("expected id: %v, actual %v", alias1.ID, id)
				}
				return alias1, nil
			},
			saveCardf: func(card *model.Card) error {
				card.CreatedAt = time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC)
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			getAnonymousAliasLastUsedf: func(userID, threadRootID globalid.ID) (bool, error) {
				return false, nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: &rpc.PostCardResponse{
				Author: alias1.Author(),
				Card: &model.CardView{
					CreatedAt: time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC).Unix(),
					Content:   "Lorem ipsum",
				},
			},
		},
		{
			name: "reply anonymously second time",
			params: rpc.PostCardParams{
				Content:     "Lorem ipsum",
				Anonymous:   true,
				ReplyCardID: "0763f798-1afc-4f65-9d95-c9c75f320bb9",
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn"), nil
			},
			notifiableForCardf: func(cardID globalid.ID) ([]globalid.ID, error) {
				return nil, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				threadRootID := globalid.ID("0763f798-1afc-4f65-9d95-c9c75f320bb9")
				if cardID != threadRootID {
					t.Errorf("expected cardID %v, actual %v", threadRootID, cardID)
				}
				return &model.Card{
					ID: threadRootID,
					AuthorToAlias: model.IdentityMap{
						user.ID: alias1.ID,
					},
				}, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				t.Errorf("unexpected call to GetUnusedAlias")
				return alias2, nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				return alias1, nil
			},
			saveCardf: func(card *model.Card) error {
				card.CreatedAt = time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC)
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			getAnonymousAliasLastUsedf: func(userID, threadRootID globalid.ID) (bool, error) {
				return false, nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
			expected: &rpc.PostCardResponse{
				Author: alias1.Author(),
				Card: &model.CardView{
					CreatedAt:     time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC).Unix(),
					Content:       "Lorem ipsum",
					ThreadReplyID: "0763f798-1afc-4f65-9d95-c9c75f320bb9",
					ThreadRootID:  "0763f798-1afc-4f65-9d95-c9c75f320bb9",
				},
			},
		},
		{
			name: "mention",
			params: rpc.PostCardParams{
				Content: "@richard Lorem ipsum",
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn"), nil
			},
			getUserByUsernamef: func(username string) (*model.User, error) {
				return &model.User{
					ID:          "ebcd9f8b-e873-4e4f-8383-72233439248a",
					Username:    "richard",
					DisplayName: "Richard Hendricks",
				}, nil
			},
			saveNotificationf: func(notification *model.Notification) error {
				return nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			newNotificationf: func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
				return nil
			},
			notifiableForCardf: func(cardID globalid.ID) ([]globalid.ID, error) {
				return nil, nil
			},
			saveCardf: func(card *model.Card) error {
				card.CreatedAt = time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC)
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			getAnonymousAliasLastUsedf: func(userID, threadRootID globalid.ID) (bool, error) {
				return false, nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: &rpc.PostCardResponse{
				Author: user.Author(),
				Card: &model.CardView{
					Content:   "@richard Lorem ipsum",
					CreatedAt: time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		{
			name: "mention anonymously",
			params: rpc.PostCardParams{
				Content:   "@richard Lorem ipsum",
				Anonymous: true,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn"), nil
			},
			getUserByUsernamef: func(username string) (*model.User, error) {
				return &model.User{
					ID:          "ebcd9f8b-e873-4e4f-8383-72233439248a",
					Username:    "richard",
					DisplayName: "Richard Hendricks",
				}, nil
			},
			saveNotificationf: func(notification *model.Notification) error {
				return nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			newNotificationf: func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
				return nil
			},
			notifiableForCardf: func(cardID globalid.ID) ([]globalid.ID, error) {
				return nil, nil
			},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				return alias1, nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				return alias1, nil
			},
			saveCardf: func(card *model.Card) error {
				card.CreatedAt = time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC)
				return nil
			},
			createIndexf: func(card *model.Card) error {
				return nil
			},
			getThreadCountf: func(cardID globalid.ID) (int, error) {
				return 0, nil
			},
			getEngagementf: func(cardID globalid.ID) (*model.Engagement, error) {
				return nil, nil
			},
			getAnonymousAliasLastUsedf: func(userID, threadRootID globalid.ID) (bool, error) {
				return false, nil
			},
			newCardf: func(ctx context.Context, session *model.Session, card *model.CardResponse) error {
				return nil
			},
			expected: &rpc.PostCardResponse{
				Author: alias1.Author(),
				Card: &model.CardView{
					Content:   "@richard Lorem ipsum",
					CreatedAt: time.Date(2016, 9, 30, 5, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserf = tt.getUserf
			r.store.GetUserByUsernamef = tt.getUserByUsernamef
			r.store.SaveNotificationf = tt.saveNotificationf
			r.store.GetCardf = tt.getCardf
			r.store.GetUnusedAliasf = tt.getUnusedAliasf
			r.store.GetAnonymousAliasf = tt.getAnonymousAliasf
			r.store.SaveCardf = tt.saveCardf
			r.store.GetThreadCountf = tt.getThreadCountf
			r.store.GetEngagementf = tt.getEngagementf
			r.store.LastPostInThreadWasAnonymousf = tt.getAnonymousAliasLastUsedf
			r.pusher.NewCardf = tt.newCardf
			r.pusher.UpdateEngagementf = tt.updateEngagementf
			r.pusher.NewNotificationf = tt.newNotificationf
			r.notifications.ExportNotificationf = tt.exportNotificationf

			r.store.SaveMentionf = func(m *model.Mention) error {
				return nil
			}
			r.store.SaveNotificationMentionf = func(nm *model.NotificationMention) error {
				return nil
			}

			r.store.SubscribersForCardf = func(cardID globalid.ID, typ string) ([]globalid.ID, error) {
				return nil, nil
			}

			r.store.SubscribeToCardf = func(userID, cardID globalid.ID, typ string) error {
				return nil
			}
			r.store.SubscribedToTypesf = func(userID, cardID globalid.ID) ([]string, error) {
				return nil, nil
			}

			req := rpc.PostCardRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.PostCard(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			tt.expected.Card.ID = resp.Card.ID
			if !reflect.DeepEqual(resp, tt.expected) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(resp, tt.expected))
			}

		})
	}
}

func TestNewInvite(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	invite, err := model.NewInvite(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		params      rpc.NewInviteParams
		err         error
		saveInvitef func(invite *model.Invite) error
		expected    *model.Invite
	}{
		{
			name: "valid",
			params: rpc.NewInviteParams{
				Invites: 3,
			},
			saveInvitef: func(invite *model.Invite) error {
				invite.RemainingUses = 3
				return nil
			},
			expected: invite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.SaveInvitef = tt.saveInvitef

			req := rpc.NewInviteRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.NewInvite(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			got := (*model.Invite)(resp)
			invite.ID = got.ID
			invite.Token = got.Token
			if !reflect.DeepEqual(tt.expected, invite) {
				t.Fatalf("unexpected result, diff :%v", pretty.Diff(tt.expected, invite))
			}
		})
	}
}

func TestRegisterDevice(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	tests := []struct {
		name      string
		params    rpc.RegisterDeviceParams
		getUserf  func(userID globalid.ID) (*model.User, error)
		saveUserf func(user *model.User) error
		err       error
	}{
		{
			name: "valid",
			params: rpc.RegisterDeviceParams{
				Token:    "H36FX",
				Platform: "Android",
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			saveUserf: func(user *model.User) error {
				if len(user.Devices) != 1 {
					t.Errorf("unexpected device count, expected: %d, actual %d", 1, len(user.Devices))
				}
				device := model.Device{
					Token:    "H36FX",
					Platform: "Android",
				}
				if !reflect.DeepEqual(device, user.Devices[device.Token]) {
					t.Errorf("unexpected device, diff: %v", pretty.Diff(device, user.Devices[device.Token]))
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserf = tt.getUserf
			r.store.SaveUserf = tt.saveUserf

			req := rpc.RegisterDeviceRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			_, err := r.RegisterDevice(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}

		})
	}
}

func TestUnregisterDevice(t *testing.T) {
	user := model.NewUser("33c79287-be80-4fbd-a605-90a90ca3253c", "chad", "chad@october.news", "Chad Unicorn")
	tests := []struct {
		name      string
		params    rpc.UnregisterDeviceParams
		getUserf  func(userID globalid.ID) (*model.User, error)
		saveUserf func(user *model.User) error
		err       error
	}{
		{
			name: "valid",
			getUserf: func(userID globalid.ID) (*model.User, error) {
				device := model.Device{
					Token:    "H36FX",
					Platform: "Android",
				}
				user.Devices = make(map[string]model.Device)
				user.Devices[device.Token] = device
				return user, nil
			},
			saveUserf: func(user *model.User) error {
				if len(user.Devices) != 0 {
					t.Errorf("unexpected device count, expected: %d, actual %d", 0, len(user.Devices))
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserf = tt.getUserf
			r.store.SaveUserf = tt.saveUserf

			req := rpc.UnregisterDeviceRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			_, err := r.UnregisterDevice(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}

		})
	}
}

func TestUpdateSettings(t *testing.T) {
	user1 := &model.User{
		ID:          "33c79287-be80-4fbd-a605-90a90ca3253c",
		Username:    "chad",
		Email:       "chad@october.news",
		DisplayName: "Chad Unicorn",
		FirstName:   "Chad",
		LastName:    "Unicorn",
	}

	user2 := &model.User{
		ID:          "33c79287-be80-4fbd-a605-90a90ca3253c",
		Username:    "chadunicorn",
		Email:       "chad@october.news",
		DisplayName: "Chad Unicorn",
		FirstName:   "Chad",
		LastName:    "Unicorn",
	}

	uppercase := "ChadUnicorn"

	tests := []struct {
		name                 string
		params               rpc.UpdateSettingsParams
		getUserf             func(userID globalid.ID) (*model.User, error)
		getAnonymousAliasesf func() ([]*model.AnonymousAlias, error)
		saveUserf            func(user *model.User) error
		updateUserf          func(ctx context.Context, session *model.Session, user *model.ExportedUser) error
		expected             *rpc.UpdateSettingsResponse
		err                  error
	}{
		{
			name:   "valid",
			params: rpc.UpdateSettingsParams{},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user1, nil
			},
			saveUserf: func(user *model.User) error {
				return nil
			},
			updateUserf: func(ctx context.Context, session *model.Session, user *model.ExportedUser) error {
				return nil
			},
			expected: (*rpc.UpdateSettingsResponse)(user1.Export(user1.ID)),
		},
		{
			name: "lowercase username",
			params: rpc.UpdateSettingsParams{
				Username: &uppercase,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user2, nil
			},
			getAnonymousAliasesf: func() ([]*model.AnonymousAlias, error) {
				return []*model.AnonymousAlias{}, nil
			},
			saveUserf: func(user *model.User) error {
				return nil
			},
			updateUserf: func(ctx context.Context, session *model.Session, user *model.ExportedUser) error {
				return nil
			},
			expected: (*rpc.UpdateSettingsResponse)(user2.Export(user2.ID)),
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserf = tt.getUserf
			r.store.GetAnonymousAliasesf = tt.getAnonymousAliasesf
			r.store.SaveUserf = tt.saveUserf
			r.pusher.UpdateUserf = tt.updateUserf

			req := rpc.UpdateSettingsRequest{
				Session: model.NewSession(user1),
				Params:  tt.params,
			}
			resp, err := r.UpdateSettings(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	user := &model.User{
		ID:          "33c79287-be80-4fbd-a605-90a90ca3253c",
		Username:    "chad",
		Email:       "chad@october.news",
		DisplayName: "Chad Unicorn",
		FirstName:   "Chad",
		LastName:    "Unicorn",
	}

	tests := []struct {
		name               string
		params             rpc.GetUserParams
		getUserf           func(userID globalid.ID) (*model.User, error)
		getUserByUsernamef func(username string) (*model.User, error)
		err                error
		expected           *rpc.GetUserResponse
	}{
		{
			name: "by ID",
			params: rpc.GetUserParams{
				UserID: user.ID,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				return user, nil
			},
			getUserByUsernamef: func(username string) (*model.User, error) {
				t.Error("unexpected call to GetUserByUsername")
				return nil, nil
			},
			expected: (*rpc.GetUserResponse)(user.Export(user.ID)),
		},
		{
			name: "by username",
			params: rpc.GetUserParams{
				Username: user.Username,
			},
			getUserf: func(userID globalid.ID) (*model.User, error) {
				t.Error("unexpected call to GetUser")
				return nil, nil
			},
			getUserByUsernamef: func(username string) (*model.User, error) {
				return user, nil
			},
			expected: (*rpc.GetUserResponse)(user.Export(user.ID)),
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUserf = tt.getUserf
			r.store.GetUserByUsernamef = tt.getUserByUsernamef

			req := rpc.GetUserRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetUser(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name                 string
		params               rpc.ValidateUsernameParams
		getAnonymousAliasesf func() ([]*model.AnonymousAlias, error)
		getUserByUsernamef   func(username string) (*model.User, error)
		err                  error
	}{
		{
			name: "valid",
			params: rpc.ValidateUsernameParams{
				Username: "chad",
			},
			getAnonymousAliasesf: func() ([]*model.AnonymousAlias, error) {
				return []*model.AnonymousAlias{}, nil
			},
			getUserByUsernamef: func(username string) (*model.User, error) {
				return nil, nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetAnonymousAliasesf = tt.getAnonymousAliasesf
			r.store.GetUserByUsernamef = tt.getUserByUsernamef

			req := rpc.ValidateUsernameRequest{
				Params: tt.params,
			}
			_, err := r.ValidateUsername(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
		})
	}
}

func TestGetNotifications(t *testing.T) {
	user := &model.User{
		ID:          "33c79287-be80-4fbd-a605-90a90ca3253c",
		Username:    "chad",
		Email:       "chad@october.news",
		DisplayName: "Chad Unicorn",
		FirstName:   "Chad",
		LastName:    "Unicorn",
	}
	notification := &model.Notification{
		ID:     "7902dfec-285b-4333-87e8-03a37a840355",
		UserID: user.ID,
	}

	tests := []struct {
		name                      string
		params                    rpc.GetNotificationsParams
		getNotificationsf         func(userID globalid.ID, pageSize, pageNumber int) ([]*model.Notification, error)
		exportNotificationf       func(notification *model.Notification) (*model.ExportedNotification, error)
		unseenNotificationsCountf func(userID globalid.ID) (int, error)
		err                       error
		expected                  *rpc.GetNotificationsResponse
	}{
		{
			name: "valid",
			getNotificationsf: func(userID globalid.ID, pageSize, pageNumber int) ([]*model.Notification, error) {
				return []*model.Notification{notification}, nil
			},
			exportNotificationf: func(notification *model.Notification) (*model.ExportedNotification, error) {
				return &model.ExportedNotification{}, nil
			},
			unseenNotificationsCountf: func(userID globalid.ID) (int, error) {
				return 0, nil
			},
			expected: &rpc.GetNotificationsResponse{
				Notifications: []*model.ExportedNotification{&model.ExportedNotification{}},
				NextPage:      true,
				UnseenCount:   0,
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetNotifcationsf = tt.getNotificationsf
			r.store.UnseenNotificationsCountf = tt.unseenNotificationsCountf
			r.notifications.ExportNotificationf = tt.exportNotificationf

			r.store.ClearEmptyNotificationsForUserf = func(userID globalid.ID) error {
				return nil
			}

			req := rpc.GetNotificationsRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetNotifications(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestUpdateNotifications(t *testing.T) {
	user := &model.User{
		ID:          "33c79287-be80-4fbd-a605-90a90ca3253c",
		Username:    "chad",
		Email:       "chad@october.news",
		DisplayName: "Chad Unicorn",
		FirstName:   "Chad",
		LastName:    "Unicorn",
	}

	tests := []struct {
		name                       string
		params                     rpc.UpdateNotificationsParams
		updateNotificationsSeenf   func(notificationIDs []globalid.ID) error
		updateNotificationsOpenedf func(notificationIDs []globalid.ID) error
		err                        error
	}{
		{
			name: "valid",
			params: rpc.UpdateNotificationsParams{
				Opened: true,
				Seen:   true,
			},
			updateNotificationsSeenf: func(notificationIDs []globalid.ID) error {
				return nil
			},
			updateNotificationsOpenedf: func(notificationIDs []globalid.ID) error {
				return nil
			},
		},
		{
			name: "none",
			params: rpc.UpdateNotificationsParams{
				Opened: false,
				Seen:   false,
			},
			updateNotificationsSeenf: func(notificationIDs []globalid.ID) error {
				t.Error("unexpected call to UpdateNotificationsSeen")
				return nil
			},
			updateNotificationsOpenedf: func(notificationIDs []globalid.ID) error {
				t.Error("unexpected call to UpdateNotificationsOpened")
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.UpdateNotificationsSeenf = tt.updateNotificationsSeenf
			r.store.UpdateNotificationsOpenedf = tt.updateNotificationsOpenedf

			req := rpc.UpdateNotificationsRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			_, err := r.UpdateNotifications(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
		})
	}
}

func TestGetAnonymousHandle(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")
	tests := []struct {
		name                       string
		params                     rpc.GetAnonymousHandleParams
		getUnusedAliasf            func(cardID globalid.ID) (*model.AnonymousAlias, error)
		getCardf                   func(id globalid.ID) (*model.Card, error)
		getAnonymousAliasf         func(id globalid.ID) (*model.AnonymousAlias, error)
		getAnonymousAliasLastUsedf func(userID, threadRootID globalid.ID) (bool, error)
		expected                   *rpc.GetAnonymousHandleResponse
		err                        error
	}{
		{
			name:   "new post",
			params: rpc.GetAnonymousHandleParams{},
			getUnusedAliasf: func(cardID globalid.ID) (*model.AnonymousAlias, error) {
				return &model.AnonymousAlias{}, nil
			},
			getCardf: func(id globalid.ID) (*model.Card, error) {
				userID := globalid.ID("b7788839-14e9-43d8-95c3-e173bdf4764d")
				aliasID := globalid.Next()
				return &model.Card{AuthorToAlias: model.IdentityMap{userID: aliasID}}, nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				return &model.AnonymousAlias{ID: "3c45aeba-fa83-4a4b-a776-5f36d966778f"}, nil
			},
			expected: &rpc.GetAnonymousHandleResponse{
				Alias:    &model.AnonymousAlias{},
				LastUsed: false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUnusedAliasf = tt.getUnusedAliasf
			r.store.GetCardf = tt.getCardf
			r.store.GetAnonymousAliasf = tt.getAnonymousAliasf
			r.store.LastPostInThreadWasAnonymousf = tt.getAnonymousAliasLastUsedf

			req := rpc.GetAnonymousHandleRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetAnonymousHandle(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestDeleteCard(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")

	tests := []struct {
		name                        string
		params                      rpc.DeleteCardParams
		getCardf                    func(globalid.ID) (*model.Card, error)
		deleteNotificationsForCardf func(cardID globalid.ID) error
		deleteCardf                 func(cardID globalid.ID) error
		pusherDeleteCardf           func(ctx context.Context, session *model.Session, id globalid.ID) error
		updateEngagementf           func(ctx context.Context, session *model.Session, cardID globalid.ID) error
		err                         error
	}{
		{
			name: "top-level card",
			getCardf: func(globalid.ID) (*model.Card, error) {
				card := &model.Card{ID: "6f83a34f-4a4e-476b-a93e-2e951fa8ca1c", OwnerID: user.ID}
				return card, nil
			},
			deleteNotificationsForCardf: func(cardID globalid.ID) error {
				return nil
			},
			deleteCardf: func(cardID globalid.ID) error {
				return nil
			},
			pusherDeleteCardf: func(ctx context.Context, session *model.Session, id globalid.ID) error {
				return nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				t.Error("unexpected innovocation")
				return nil
			},
		},
		{
			name: "comment",
			getCardf: func(globalid.ID) (*model.Card, error) {
				card := &model.Card{
					ID:            "6f83a34f-4a4e-476b-a93e-2e951fa8ca1c",
					ThreadRootID:  "d28fd6d6-8247-407a-b846-8259246f08e5",
					ThreadReplyID: "70c79bd3-905c-42e3-a87c-eeeeeec79ea4",
					OwnerID:       user.ID,
				}
				return card, nil
			},
			deleteNotificationsForCardf: func(cardID globalid.ID) error {
				return nil
			},
			deleteCardf: func(cardID globalid.ID) error {
				return nil
			},
			pusherDeleteCardf: func(ctx context.Context, session *model.Session, id globalid.ID) error {
				return nil
			},
			updateEngagementf: func(ctx context.Context, session *model.Session, cardID globalid.ID) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetCardf = tt.getCardf
			r.store.DeleteNotificationsForCardf = tt.deleteNotificationsForCardf
			r.store.DeleteCardf = tt.deleteCardf
			r.pusher.DeleteCardf = tt.pusherDeleteCardf
			r.pusher.UpdateEngagementf = tt.updateEngagementf

			r.store.DeleteMentionsForCardf = func(cardID globalid.ID) error {
				return nil
			}
			r.store.GetUserf = func(userID globalid.ID) (*model.User, error) {
				return &model.User{DisplayName: "Test User"}, nil
			}

			req := rpc.DeleteCardRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			_, err := r.DeleteCard(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
		})
	}
}

func TestFollowUser(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")

	tests := []struct {
		name          string
		params        rpc.FollowUserParams
		saveFollowerf func(followerID, followeeID globalid.ID) error
		isFollowingf  func(followerID, followeeID globalid.ID) (bool, error)
		err           error
	}{
		{
			name: "valid",
			params: rpc.FollowUserParams{
				UserID: "1e295ea4-a43a-49b0-aa95-bbb670914593",
			},
			saveFollowerf: func(followerID, followeeID globalid.ID) error {
				return nil
			},
			isFollowingf: func(followerID, followeeID globalid.ID) (bool, error) {
				return false, nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.SaveFollowerf = tt.saveFollowerf
			r.store.IsFollowingf = tt.isFollowingf

			r.store.SaveNotificationFollowf = func(notifID, followerID, followeeID globalid.ID) error {
				return nil
			}

			r.store.SaveNotificationf = func(n *model.Notification) error {
				return nil
			}

			r.store.LatestForTypef = func(userID, targetID globalid.ID, typ string, unopenedOnly bool) (*model.Notification, error) {
				return &model.Notification{ID: globalid.Next()}, nil
			}

			r.store.ClearEmptyNotificationsf = func() error {
				return nil
			}

			r.store.GetFollowExportDataf = func(n *model.Notification) (*datastore.FollowNotificationExportData, error) {
				return nil, nil
			}

			r.notifications.ExportNotificationf = func(n *model.Notification) (*model.ExportedNotification, error) {
				return nil, nil
			}
			r.pusher.NewNotificationf = func(ctx context.Context, session *model.Session, notif *model.ExportedNotification) error {
				return nil
			}

			req := rpc.FollowUserRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			_, err := r.FollowUser(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
		})
	}
}

func TestModifyCardScore(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")
	tests := []struct {
		name                   string
		params                 rpc.ModifyCardScoreParams
		saveScoreModificationf func(m *model.ScoreModification) error
		err                    error
		expected               *rpc.ModifyCardScoreResponse
	}{
		{
			name:   "valid",
			params: rpc.ModifyCardScoreParams{},
			saveScoreModificationf: func(m *model.ScoreModification) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.SaveScoreModificationf = tt.saveScoreModificationf
			req := rpc.ModifyCardScoreRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.ModifyCardScore(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestUnfollowUser(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")

	tests := []struct {
		name            string
		params          rpc.UnfollowUserParams
		deleteFollowerf func(followerID, followeeID globalid.ID) error
		err             error
	}{
		{
			name: "valid",
			params: rpc.UnfollowUserParams{
				UserID: "1e295ea4-a43a-49b0-aa95-bbb670914593",
			},
			deleteFollowerf: func(followerID, followeeID globalid.ID) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.DeleteFollowerf = tt.deleteFollowerf
			r.store.ClearEmptyNotificationsf = func() error {
				return nil
			}

			req := rpc.UnfollowUserRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			_, err := r.UnfollowUser(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
		})
	}
}

func TestGetFollowingUsers(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")

	tests := []struct {
		name          string
		params        rpc.GetFollowingUsersParams
		getFollowingf func(userID globalid.ID) ([]*model.User, error)
		err           error
		expected      *rpc.GetFollowingUsersResponse
	}{
		{
			name:   "valid",
			params: rpc.GetFollowingUsersParams{},
			getFollowingf: func(userID globalid.ID) ([]*model.User, error) {
				return []*model.User{user}, nil
			},
			expected: (*rpc.GetFollowingUsersResponse)(&[]*model.ExportedUser{user.Export(user.ID)}),
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetFollowingf = tt.getFollowingf

			req := rpc.GetFollowingUsersRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetFollowingUsers(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestGetPostsForUser(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")

	tests := []struct {
		name                   string
		params                 rpc.GetPostsForUserParams
		getPostedCardsForNodef func(nodeID globalid.ID, skip, count int) ([]*model.Card, error)
		err                    error
		expected               *rpc.GetPostsForUserResponse
	}{
		{
			name:   "valid",
			params: rpc.GetPostsForUserParams{},
			getPostedCardsForNodef: func(nodeID globalid.ID, skip, count int) ([]*model.Card, error) {
				return nil, nil
			},
			expected: &rpc.GetPostsForUserResponse{
				NextPage: false,
				Cards:    []*model.CardResponse{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetPostedCardsForNodef = tt.getPostedCardsForNodef

			req := rpc.GetPostsForUserRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetPostsForUser(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestGetFeaturesForUser(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")

	tests := []struct {
		name           string
		params         rpc.GetFeaturesForUserParams
		getOnSwitchesf func() ([]*model.FeatureSwitch, error)
		err            error
		expected       *rpc.GetFeaturesForUserResponse
	}{
		{
			name:   "valid",
			params: rpc.GetFeaturesForUserParams{},
			getOnSwitchesf: func() ([]*model.FeatureSwitch, error) {
				return nil, nil
			},
			expected: &rpc.GetFeaturesForUserResponse{},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetOnSwitchesf = tt.getOnSwitchesf

			req := rpc.GetFeaturesForUserRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetFeaturesForUser(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestUploadImage(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")
	expected := rpc.UploadImageResponse("http://assets.october.news/images/a1f3199d-eaf4-481d-83cf-22398df1b54a.png")

	tests := []struct {
		name                        string
		params                      rpc.UploadImageParams
		saveBase64CardContentImagef func(data string) (string, string, error)
		saveBase64ProfileImagef     func(data string) (string, string, error)
		err                         error
		expected                    *rpc.UploadImageResponse
	}{
		{
			name: "card content image",
			params: rpc.UploadImageParams{
				CardContentImageData: "123=",
			},
			saveBase64CardContentImagef: func(data string) (string, string, error) {
				return string(expected), "", nil
			},
			expected: &expected,
		},
		{
			name: "profile image",
			params: rpc.UploadImageParams{
				ProfileImageData: "123=",
			},
			saveBase64ProfileImagef: func(data string) (string, string, error) {
				return string(expected), "", nil
			},
			expected: &expected,
		},
		{
			name:   "no parameter",
			params: rpc.UploadImageParams{},
			err:    rpc.ErrNoImageData,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.imageProcessor.SaveBase64CardContentImagef = tt.saveBase64CardContentImagef
			r.imageProcessor.SaveBase64ProfileImagef = tt.saveBase64ProfileImagef

			req := rpc.UploadImageRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.UploadImage(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestGetTaggableUsers(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")
	richard := &model.User{
		ID:               "0b0ce4ed-abed-4357-accb-8909036d01b2",
		Email:            "richard@piedpiper.com",
		Username:         "richard",
		DisplayName:      "Richard Hendricks",
		ProfileImagePath: "richard.png",
	}
	erlich := &model.User{
		ID:               "52df8f20-8b15-4d9c-a3c7-2d5417df2346",
		Email:            "erlich@piedpiper.com",
		Username:         "erlich",
		DisplayName:      "Erlich Bachman",
		ProfileImagePath: "erlich.png",
	}
	dinesh := &model.User{
		ID:               "179d1823-f025-4222-8ccc-bc18ba98cdf5",
		Email:            "dinesh@piedpiper.com",
		Username:         "dinesh",
		DisplayName:      "Dinesh Chugtai",
		ProfileImagePath: "dinesh.png",
	}
	alias1 := &model.AnonymousAlias{
		ID:               "40c5a86a-ab1a-4174-9096-afdc40dd862b",
		Username:         "egg",
		DisplayName:      "Anonymous",
		ProfileImagePath: "egg.png",
	}
	alias2 := &model.AnonymousAlias{
		ID:               "5858f233-b6e1-40f5-a0ac-a73be339f1d6",
		Username:         "mouse",
		DisplayName:      "Anonymous",
		ProfileImagePath: "mouse.png",
	}
	card := &model.Card{
		ID:      "da6e9385-08e9-405e-97ed-05ea6c0cb473",
		AliasID: alias1.ID,
		AuthorToAlias: model.IdentityMap{
			richard.ID: alias1.ID,
			erlich.ID:  alias2.ID,
		},
	}

	tests := []struct {
		name               string
		params             rpc.GetTaggableUsersParams
		getUsersf          func() ([]*model.User, error)
		getCardf           func(cardID globalid.ID) (*model.Card, error)
		getAnonymousAliasf func(id globalid.ID) (*model.AnonymousAlias, error)
		err                error
		expected           *rpc.GetTaggableUsersResponse
	}{
		{
			name:   "without card",
			params: rpc.GetTaggableUsersParams{},
			getUsersf: func() ([]*model.User, error) {
				return []*model.User{store.RootUser, richard, erlich, dinesh}, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				return card, nil
			},
			expected: &rpc.GetTaggableUsersResponse{
				richard.TaggableUser(),
				erlich.TaggableUser(),
				dinesh.TaggableUser(),
			},
		},
		{
			name: "with card",
			params: rpc.GetTaggableUsersParams{
				CardID: "da6e9385-08e9-405e-97ed-05ea6c0cb473",
			},
			getUsersf: func() ([]*model.User, error) {
				return []*model.User{store.RootUser, richard, erlich, dinesh}, nil
			},
			getCardf: func(cardID globalid.ID) (*model.Card, error) {
				if cardID != card.ID {
					t.Errorf("expected cardID %v, actual %v", card.ID, cardID)
				}
				return card, nil
			},
			getAnonymousAliasf: func(id globalid.ID) (*model.AnonymousAlias, error) {
				if id == alias1.ID {
					return alias1, nil
				}
				if id == alias2.ID {
					return alias2, nil
				}
				t.Errorf("expected id %v, actual %v", card.ID, id)
				return nil, nil
			},
			expected: &rpc.GetTaggableUsersResponse{
				richard.TaggableUser(),
				erlich.TaggableUser(),
				dinesh.TaggableUser(),
				alias1.TaggableUser(),
				alias2.TaggableUser(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUsersf = tt.getUsersf
			r.store.GetCardf = tt.getCardf
			r.store.GetAnonymousAliasf = tt.getAnonymousAliasf

			req := rpc.GetTaggableUsersRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetTaggableUsers(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestGetInvites(t *testing.T) {
	user := model.NewUser(globalid.Next(), "chad", "chad@october.news", "Chad Unicorn")

	tests := []struct {
		name               string
		params             rpc.GetInvitesParams
		getInvitesForUserf func(id globalid.ID) ([]*model.Invite, error)
		err                error
		expected           *rpc.GetInvitesResponse
	}{
		{
			name:   "valid",
			params: rpc.GetInvitesParams{},
			getInvitesForUserf: func(id globalid.ID) ([]*model.Invite, error) {
				return nil, nil
			},
			expected: &rpc.GetInvitesResponse{},
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetInvitesForUserf = tt.getInvitesForUserf

			req := rpc.GetInvitesRequest{
				Session: model.NewSession(user),
				Params:  tt.params,
			}
			resp, err := r.GetInvites(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

func TestGetOnboardingData(t *testing.T) {
	chad := &model.User{
		ID:               "33c79287-be80-4fbd-a605-90a90ca3253c",
		Username:         "chad",
		Email:            "chad@october.news",
		DisplayName:      "Chad Unicorn",
		FirstName:        "Chad",
		LastName:         "Unicorn",
		ProfileImagePath: "chad.jpg",
	}

	erlich := &model.User{
		ID:               "06d41407-107f-45ee-95f3-c43cd17950fe",
		Username:         "erlich",
		Email:            "erlich@october.news",
		DisplayName:      "Erlich Bachman",
		FirstName:        "Erlich",
		LastName:         "Bachman",
		ProfileImagePath: "erlich.jpg",
	}

	tests := []struct {
		name       string
		params     rpc.GetOnboardingDataParams
		getUsersf  func() ([]*model.User, error)
		getInvitef func(id globalid.ID) (*model.Invite, error)
		getUserf   func(nodeID globalid.ID) (u *model.User, err error)
		err        error
		expected   *rpc.GetOnboardingDataResponse
	}{
		{
			name:   "valid",
			params: rpc.GetOnboardingDataParams{},
			getUsersf: func() ([]*model.User, error) {
				return []*model.User{chad, erlich}, nil
			},
			getInvitef: func(id globalid.ID) (*model.Invite, error) {
				return &model.Invite{
					NodeID: erlich.ID,
				}, nil
			},
			getUserf: func(nodeID globalid.ID) (u *model.User, err error) {
				if nodeID == chad.ID {
					return chad, nil
				} else if nodeID == erlich.ID {
					return erlich, nil
				}
				t.Errorf("unexpected userID: %v", nodeID)
				return chad, nil
			},
			expected: &rpc.GetOnboardingDataResponse{
				NetworkProfilePictures: []string{"chad.jpg", "erlich.jpg"},
				InvitingUser:           erlich.Export(chad.ID),
			},
		},
		{
			name:   "user without invite",
			params: rpc.GetOnboardingDataParams{},
			getUsersf: func() ([]*model.User, error) {
				return []*model.User{chad}, nil
			},
			getInvitef: func(id globalid.ID) (*model.Invite, error) {
				return nil, sql.ErrNoRows
			},
			getUserf: func(nodeID globalid.ID) (u *model.User, err error) {
				if nodeID != chad.ID {
					t.Errorf("expected userID: %v, actual %v", chad.ID, nodeID)
				}
				return chad, nil
			},
			expected: &rpc.GetOnboardingDataResponse{
				NetworkProfilePictures: []string{"chad.jpg"},
				InvitingUser:           nil,
			},
		},
		{
			name:   "get users fails",
			params: rpc.GetOnboardingDataParams{},
			getUsersf: func() ([]*model.User, error) {
				return nil, ErrFailure
			},
			err: ErrFailure,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRPC(t)
			r.store.GetUsersf = tt.getUsersf
			r.store.GetInvitef = tt.getInvitef
			r.store.GetUserf = tt.getUserf

			req := rpc.GetOnboardingDataRequest{
				Session: model.NewSession(chad),
				Params:  tt.params,
			}
			resp, err := r.GetOnboardingData(context.Background(), req)
			if err != tt.err {
				t.Fatalf("expected error %v, actual %v", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expected, resp) {
				t.Fatalf("unexpected result, diff: %v", pretty.Diff(tt.expected, resp))
			}
		})
	}
}

type rpcMock struct {
	rpc.RPC

	oauth2         *mockOAuth2
	imageProcessor *mockImageProcessor
	pusher         *pusherMock
	store          *testStore
	worker         *testWorker
	notifications  *mockNotifications
	responses      *mockResponses
}

func newRPC(t *testing.T) *rpcMock {
	t.Helper()

	s := newStore()
	w := newWorker()
	c := rpc.NewConfig()
	ip := &mockImageProcessor{}
	p := &pusherMock{}
	oa2 := &mockOAuth2{}
	resps := &mockResponses{}

	m := &emailsender.MailTemplates{
		PasswordReset: &emailsender.MailTemplate{},
		UserInvite:    &emailsender.MailTemplate{},
	}
	ns := &mockNotifications{}

	r := rpc.NewRPC(s, w, oa2, ip, &c, log.NopLogger(), m, &worker.Notifier{}, p, &model.Settings{}, ns, &mockIndexer{}, nil, resps)
	return &rpcMock{
		RPC:            r,
		oauth2:         oa2,
		imageProcessor: ip,
		pusher:         p,
		store:          s,
		worker:         w,
		notifications:  ns,
		responses:      resps,
	}
}

func TestSanitize(t *testing.T) {
	params := rpc.AuthParams{
		Password:    "secret",
		AccessToken: "ABC",
	}
	sanitized := params.Sanitize().(rpc.AuthParams)
	if sanitized.Password != "[filtered]" {
		t.Errorf("expected password to be sanitized")
	}
	if params.Password != "secret" {
		t.Errorf("Sanitize() mutated the orginal parameters")
	}
}
