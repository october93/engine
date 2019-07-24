package server_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	jd "github.com/josephburnett/jd/lib"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/rpc/server"
	"github.com/october93/engine/store/datastore"
)

type testStore struct {
	DeleteSessionf              func(globalid.ID) error
	GetUserByEmailf             func(email string) (*model.User, error)
	GetUserByUsernamef          func(username string) (*model.User, error)
	GetUserf                    func(nodeID globalid.ID) (u *model.User, err error)
	GetUsersf                   func() (res []*model.User, err error)
	SaveInvitef                 func(invite *model.Invite) error
	GetInviteByTokenf           func(token string) (*model.Invite, error)
	DeleteInvitef               func(id globalid.ID) error
	UpdateInvitef               func(invite *model.Invite) error
	CountGivenReactionsForUserf func(posterID globalid.ID, timeFrom, timeTo time.Time, onlyLikes bool) (c int, err error)
	SaveSessionf                func(session *model.Session) error
	SaveStoref                  func(*model.User) error
	SaveUserf                   func(u *model.User) (err error)
	SaveResetTokenf             func(rt *model.ResetToken) error
	GetResetTokenf              func(userID globalid.ID) (*model.ResetToken, error)
	GetSessionsf                func() ([]*model.Session, error)
	GetSessionf                 func(id globalid.ID) (*model.Session, error)
	GetAnonymousAliasf          func(id globalid.ID) (*model.AnonymousAlias, error)
	GetUnusedAliasf             func(used []globalid.ID) (*model.AnonymousAlias, error)
	GetThreadCountf             func(id globalid.ID) (int, error)
	GetOnSwitchesf              func() ([]*model.FeatureSwitch, error)
	GetHistoricalFeedCardsf     func(id globalid.ID, num int) ([]globalid.ID, error)
	GetUserPostIDsf             func(id globalid.ID, num int) ([]globalid.ID, error)
}

func (s *testStore) Close() error {
	return nil
}

func (s *testStore) SaveUser(u *model.User) (err error) {
	if s.SaveStoref == nil {
		return nil
	}
	return s.SaveStoref(u)
}

func (s *testStore) DeleteSession(id globalid.ID) error {
	if s.DeleteSessionf == nil {
		return nil
	}
	return s.DeleteSessionf(id)
}

func (s *testStore) GetUser(nodeID globalid.ID) (u *model.User, err error) {
	if s.GetUserf == nil {
		return &model.User{}, nil
	}
	return s.GetUserf(nodeID)
}
func (s *testStore) GetUserByEmail(email string) (*model.User, error) { return s.GetUserByEmailf(email) }
func (s *testStore) GetUserByUsername(username string) (*model.User, error) {
	return s.GetUserByUsernamef(username)
}
func (s *testStore) GetUsers() (res []*model.User, err error) { return s.GetUsersf() }
func (s *testStore) SaveInvite(invite *model.Invite) error {
	return s.SaveInvitef(invite)
}
func (s *testStore) GetInviteByToken(token string) (*model.Invite, error) {
	return s.GetInviteByTokenf(token)
}
func (s *testStore) DeleteInvite(id globalid.ID) error {
	return s.DeleteInvitef(id)
}
func (s *testStore) UpdateInvite(invite *model.Invite) error {
	return s.UpdateInvitef(invite)
}

func (s *testStore) CountGivenReactionsForUser(posterID globalid.ID, timeFrom, timeTo time.Time, onlyLikes bool) (c int, err error) {
	return s.CountGivenReactionsForUserf(posterID, timeFrom, timeTo, onlyLikes)
}
func (s *testStore) SaveResetToken(rt *model.ResetToken) error {
	return s.SaveResetTokenf(rt)
}
func (s *testStore) GetResetToken(userID globalid.ID) (*model.ResetToken, error) {
	return s.GetResetTokenf(userID)
}

func (s *testStore) SaveSession(session *model.Session) error { return s.SaveSessionf(session) }

func (s *testStore) GetSessions() ([]*model.Session, error) {
	return s.GetSessionsf()
}

func (s *testStore) GetSession(id globalid.ID) (*model.Session, error) {
	return s.GetSessionf(id)
}

func (s *testStore) GetAnonymousAlias(id globalid.ID) (*model.AnonymousAlias, error) {
	return s.GetAnonymousAliasf(id)
}
func (s *testStore) GetUnusedAlias(used []globalid.ID) (*model.AnonymousAlias, error) {
	return s.GetUnusedAliasf(used)
}
func (s *testStore) GetThreadCount(id globalid.ID) (int, error) {
	return s.GetThreadCountf(id)
}

func (r *testStore) GetOnSwitches() ([]*model.FeatureSwitch, error) {
	return r.GetOnSwitchesf()
}

func (s *testStore) GetHistoricalFeedCards(id globalid.ID, num int) ([]globalid.ID, error) {
	return s.GetHistoricalFeedCardsf(id, num)
}

func (s *testStore) GetUserPostIDs(id globalid.ID, num int) ([]globalid.ID, error) {
	return s.GetUserPostIDsf(id, num)
}

func TestReactToCardEndpoint(t *testing.T) {
	var tests = []struct {
		name         string
		data         string
		callback     string
		user         *model.User
		cardID       globalid.ID
		response     string
		expectError  bool
		reactToCardf func(ctx context.Context, req rpc.ReactToCardRequest) (*rpc.ReactToCardResponse, error)
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:     "throw react error",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			cardID:   globalid.Next(),
			data:     `{"cardID":"!ID!", "reaction":"like"}`,
			callback: "foo",
			reactToCardf: func(ctx context.Context, req rpc.ReactToCardRequest) (*rpc.ReactToCardResponse, error) {
				return nil, errors.New("boom!")
			},
			response: `{"rpc":"foo","ack":"!ID!","error":"boom!"}`,
		},

		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			cardID:   globalid.Next(),
			data:     `{"cardID":"!ID!", "reaction":"like"}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!ID!","data":null}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.data = strings.Replace(tt.data, "!ID!", string(tt.cardID), -1)
			r := &rpcMock{}
			r.ReactToCardf = tt.reactToCardf
			if r.ReactToCardf == nil {
				r.ReactToCardf = func(ctx context.Context, req rpc.ReactToCardRequest) (*rpc.ReactToCardResponse, error) {
					return nil, nil
				}
			}
			endpoint := server.ReactToCardEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!ID!", string(tt.user.ID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestGetUsersEndpoint(t *testing.T) {
	var tests = []struct {
		name        string
		callback    string
		user        *model.User
		newUsers    []*model.User
		response    string
		expectError bool
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty body",
			expectError: true,
		},
		{
			name:        "user with empty data",
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:     "user and callback with empty",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			callback: "foo",
			response: `{"rpc":"foo","ack":"!userID!","data":[]}`,
		},
		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			newUsers: []*model.User{model.NewUser(globalid.ID("29657849-8bb5-4f94-8561-e8388dde76f5"), "andrew", "andrew@yahoo.com", "rob")},
			callback: "foo",
			response: `{"rpc":"foo","ack":"!userID!","data":[{"nodeId":"29657849-8bb5-4f94-8561-e8388dde76f5","displayname":"rob","profileimg_path":"","cover_image_path":"","userBio":"","email":"andrew@yahoo.com","username":"andrew","botchedSignup":false}]}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			if tt.newUsers != nil {
				r.GetUsersf = func(ctx context.Context, req rpc.GetUsersRequest) (*rpc.GetUsersResponse, error) {
					return (*rpc.GetUsersResponse)(&tt.newUsers), nil
				}
			}
			if r.GetUsersf == nil {
				r.GetUsersf = func(ctx context.Context, req rpc.GetUsersRequest) (*rpc.GetUsersResponse, error) {
					return (*rpc.GetUsersResponse)(&[]*model.User{}), nil
				}
			}
			endpoint := server.GetUsersEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint: unexpected error: %v", err)
			}
			tt.response = strings.Replace(tt.response, "!userID!", string(tt.user.ID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestGetUserEndpoint(t *testing.T) {
	var tests = []struct {
		name        string
		callback    string
		data        string
		user        *model.User
		response    string
		getUserf    func(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error)
		expectError bool
		regexError  string
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty body",
			expectError: true,
		},
		{
			name:        "user with empty data",
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:       "user with invalid data",
			callback:   "callback",
			data:       `{"foo":}`,
			user:       model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			regexError: `"error":".*invalid character.*"`,
		},
		{
			name:     "get user error",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "callback",
			response: `{"rpc":"callback","ack":"!nodeID!","error":"boom!"}`,
			getUserf: func(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:       "user not found",
			user:       model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:       `{}`,
			callback:   "callback",
			regexError: `"error":"user not found"`,
			getUserf: func(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
				return nil, datastore.ErrUserNotFound
			},
		},
		{
			name:     "user querying themselves",
			user:     model.NewUser("7a98d5b0-8b11-446b-b056-164f76d20ec8", "rob", "rob@yahoo.com", "rob"),
			data:     `{"username":"rob"}`,
			callback: "callback",
			getUserf: func(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
				user := &model.User{
					ID:               "7a98d5b0-8b11-446b-b056-164f76d20ec8",
					DisplayName:      "rob",
					FirstName:        "Rob",
					LastName:         "Pike",
					Email:            "rob@yahoo.com",
					ProfileImagePath: "cdn/image.gif",
					Bio:              "I write go.",
					SearchKey:        "123=",
				}
				return (*rpc.GetUserResponse)(user.Export("7a98d5b0-8b11-446b-b056-164f76d20ec8")), nil
			},
			response: `{"rpc":"callback","ack":"!userID!","data":{"nodeId":"7a98d5b0-8b11-446b-b056-164f76d20ec8","displayname":"rob","firstName":"Rob","lastName":"Pike","profileimg_path":"cdn/image.gif","cover_image_path":"","userBio":"I write go.","email":"rob@yahoo.com","username":"","searchKey":"123=","botchedSignup":false}}`,
		},
		{
			name:     "user querying another user",
			user:     model.NewUser("7a98d5b0-8b11-446b-b056-164f76d20ec8", "rob", "rob@yahoo.com", "rob"),
			data:     `{"username":"dave"}`,
			callback: "callback",
			getUserf: func(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
				user := &model.User{
					ID:               "7a98d5b0-8b11-446b-b056-164f76d20ec8",
					DisplayName:      "Dave",
					FirstName:        "Dave",
					LastName:         "Cheney",
					ProfileImagePath: "cdn/image.jpg",
					Bio:              "Not available for discussions about generic types",
					SearchKey:        "123=",
				}
				return (*rpc.GetUserResponse)(user.Export(globalid.Nil)), nil

			},
			response: `{"rpc":"callback","ack":"!userID!","data":{"nodeId":"!userID!","displayname":"Dave","firstName":"Dave","lastName":"Cheney","profileimg_path":"cdn/image.jpg","cover_image_path":"","userBio":"Not available for discussions about generic types","username":"","botchedSignup":false}}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.GetUserf = tt.getUserf
			if r.GetUserf == nil {
				r.GetUserf = func(ctx context.Context, req rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
					return nil, nil
				}
			}
			endpoint := server.GetUserEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint: unexpected error: %v", err)
			}
			tt.response = strings.Replace(tt.response, "!userID!", string(tt.user.ID), -1)
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			if tt.regexError != "" {
				jsonError := regexp.MustCompile(tt.regexError)
				if !jsonError.MatchString(b.String()) {
					t.Fatalf("expected response: \n\t%s\nto match:\n\t%s", b.String(), tt.regexError)
				}
				// done here
				return
			}
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestConnectUsersEndpoint(t *testing.T) {
	var tests = []struct {
		name          string
		users         []string
		data          string
		callback      string
		user          *model.User
		response      string
		connectUsersf func(ctx context.Context, req rpc.ConnectUsersRequest) (*rpc.ConnectUsersResponse, error)
		expectError   bool
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "user with invalid data",
			data:        `{"foo":}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:     "no users",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{"users":[]}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":null}`,
		},
		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     fmt.Sprintf(`{"users":["%s"]}`, globalid.Next().String()),
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":null}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.ConnectUsersf = tt.connectUsersf
			if r.ConnectUsersf == nil {
				r.ConnectUsersf = func(ctx context.Context, req rpc.ConnectUsersRequest) (*rpc.ConnectUsersResponse, error) {
					return nil, nil
				}
			}
			endpoint := server.ConnectUsersEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestAuthLoginEndpoint(t *testing.T) {
	var tests = []struct {
		name        string
		requestID   globalid.ID
		authf       func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error)
		data        string
		user        *model.User
		response    string
		callback    string
		expectError bool
	}{
		{
			name:        "no body",
			callback:    "callback",
			expectError: true,
		},
		{
			name:        "empty data",
			callback:    "callback",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "invalid data",
			callback:    "callback",
			data:        `{"foo":}`,
			expectError: true,
		},
		{
			name:      "login error",
			requestID: globalid.Next(),
			callback:  "callback",
			data:      "{}",
			response:  `{"rpc":"callback","ack":"!nodeID!","error":"boom!"}`,
			authf: func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:      "valid",
			requestID: globalid.Next(),
			callback:  "callback",
			user:      &model.User{ID: globalid.Next(), Email: "rob@rob.com", DisplayName: "rob", ProfileImagePath: "cdn/rob.gif", FirstName: "Rob", LastName: "Pike"},
			data:      `{"username":"rob","password":"go"}`,
			response:  `{"rpc":"callback","ack":"!nodeID!","data":{"user":{"nodeId":"!nodeID!","displayname":"!displayname!","firstName":"Rob","lastName":"Pike","profileimg_path":"!picturepath!","cover_image_path":"","userBio":"","email":"rob@rob.com","username":"","botchedSignup":false},"session":{"id":"!sessionID!"}}}`,
		},
		{
			name:      "admin",
			requestID: globalid.Next(),
			callback:  "callback",
			user:      &model.User{ID: globalid.Next(), DisplayName: "pike", ProfileImagePath: "cdn/pike.gif", FirstName: "Rob", LastName: "Pike", Admin: true},
			data:      `{"username":"rob","password":"go"}`,
			response:  `{"rpc":"callback","ack":"!nodeID!","data":{"user":{"nodeId":"!nodeID!","displayname":"!displayname!","firstName":"Rob","lastName":"Pike","profileimg_path":"!picturepath!","cover_image_path":"","userBio":"","username":"","admin":true,"botchedSignup":false},"session":{"id":"!sessionID!"}}}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.Authf = tt.authf
			if r.Authf == nil {
				r.Authf = func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error) {
					if tt.user != nil {
						return &rpc.AuthResponse{
							User:    tt.user.Export(tt.user.ID),
							Session: &model.Session{ID: tt.requestID},
						}, nil
					}
					return nil, errors.New("not found")
				}
			}
			endpoint := server.AuthEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}
			if tt.user == nil && tt.requestID != globalid.Nil {
				ctx = context.WithValue(ctx, protocol.RequestID, tt.requestID)
			}
			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			if tt.user != nil {
				tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
				tt.response = strings.Replace(tt.response, "!sessionID!", string(tt.requestID), -1)
				tt.response = strings.Replace(tt.response, "!displayname!", tt.user.DisplayName, -1)
				tt.response = strings.Replace(tt.response, "!picturepath!", tt.user.ProfileImagePath, -1)
			}
			if tt.user == nil && tt.requestID != globalid.Nil {
				tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.requestID), -1)
			}
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestLogoutEndpoint(t *testing.T) {
	var tests = []struct {
		name        string
		session     *model.Session
		logoutf     func(ctx context.Context, req rpc.LogoutRequest) (*rpc.LogoutResponse, error)
		data        string
		response    string
		callback    string
		expectError bool
	}{
		{
			name:        "no body",
			callback:    "callback",
			expectError: true,
		},
		{
			name:     "logout error",
			session:  &model.Session{ID: globalid.Next()},
			callback: "callback",
			data:     `{"idToken":"7a98d5b0-8b11-446b-b056-164f76d20ec8"}`,
			response: `{"rpc":"callback","ack":"!sessionID!","error":"boom!"}`,
			logoutf: func(ctx context.Context, req rpc.LogoutRequest) (*rpc.LogoutResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid",
			session:  &model.Session{ID: globalid.Next()},
			callback: "callback",
			response: `{"rpc":"callback","ack":"!sessionID!","data":null}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.Logoutf = tt.logoutf
			if r.Logoutf == nil {
				r.Logoutf = func(ctx context.Context, req rpc.LogoutRequest) (*rpc.LogoutResponse, error) {
					return nil, nil
				}
			}
			endpoint := server.LogoutEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			if tt.session != nil {
				pw.SetSession(tt.session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.session.ID)
			}

			err := endpoint(ctx, tt.session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!sessionID!", string(tt.session.ID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestGetCardsEndpoint(t *testing.T) {
	now := time.Now().UTC()
	var tests = []struct {
		name        string
		users       []string
		data        string
		callback    string
		user        *model.User
		response    string
		getCardsf   func(ctx context.Context, req rpc.GetCardsRequest) (*rpc.GetCardsResponse, error)
		getuserf    func(id globalid.ID) (*model.User, error)
		getaliasf   func(id globalid.ID) (*model.AnonymousAlias, error)
		cards       []*model.Card
		expectError bool
		regexError  string
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "user with invalid data",
			data:        `{"foo":}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:       "invalid pageSize",
			user:       model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:       `{"pageSize":"boom","pageNumber":2}`,
			callback:   "foo",
			regexError: `"error":"json.*"`,
		},
		{
			name:       "invalid pageNumber",
			user:       model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:       `{"pageSize":-1,"pageNumber":"foo"}`,
			callback:   "foo",
			regexError: `"error":"json.*"`,
		},
		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{"pageSize":1,"pageNumber":2}`,
			callback: "foo",
			cards:    []*model.Card{&model.Card{Title: "New Card", OwnerID: "dontmatter", Content: "some content...", URL: "http://somesite.com/url", BackgroundColor: "indigo-blue", CreatedAt: now, ThreadLevel: 3}},
			response: `{"rpc":"foo","ack":"!nodeID!","data":{"cards":[{"card":{"cardID":"","threadLevel":3,"coinsEarned":0,"title":"New Card","body":"some content...","url":"http://somesite.com/url","bgColor":"indigo-blue","background_image_path":"","anonymous":false,"post_timestamp":!timestamp_unix!},"author":{"nodeId":"!nodeID!","displayname":"rob","username":"rob","profileimg_path":"","isAnonymous":false},"replies":3,"score":0,"subscribed":false}],"hasNextPage":false}}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.GetCardsf = tt.getCardsf
			if r.GetCardsf == nil {
				r.GetCardsf = func(ctx context.Context, req rpc.GetCardsRequest) (*rpc.GetCardsResponse, error) {
					cards := make([]*model.CardResponse, len(tt.cards))
					for i, card := range tt.cards {
						cards[i] = &model.CardResponse{Card: card.Export(), Author: tt.user.Author(), Replies: 3}
					}
					return &rpc.GetCardsResponse{Cards: cards}, nil
				}
			}

			teststore := &testStore{}
			teststore.GetUserf = tt.getuserf
			if teststore.GetUserf == nil {
				teststore.GetUserf = func(id globalid.ID) (*model.User, error) {
					return tt.user, nil
				}
			}

			teststore.GetAnonymousAliasf = tt.getaliasf
			if teststore.GetAnonymousAliasf == nil {
				teststore.GetAnonymousAliasf = func(id globalid.ID) (*model.AnonymousAlias, error) {
					return nil, nil
				}
			}

			if teststore.GetThreadCountf == nil {
				teststore.GetThreadCountf = func(id globalid.ID) (int, error) {
					return 3, nil
				}
			}

			endpoint := server.GetCardsEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			tt.response = strings.Replace(tt.response, "!timestamp_unix!", fmt.Sprintf("%d", now.Unix()), -1)
			if tt.regexError != "" {
				jsonError := regexp.MustCompile(tt.regexError)
				if !jsonError.MatchString(b.String()) {
					t.Fatalf("expected response: \n\t%s\nto match:\n\t%s", b.String(), tt.regexError)
				}
				// done here
				return
			}
			got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response
			diffJSON(got, exp, t)
		})
	}
}

func diffJSON(got, exp string, t *testing.T) {
	t.Helper()
	if got != exp {
		a, err := jd.ReadJsonString(got)
		if err != nil {
			t.Fatal(err)
		}
		b, err := jd.ReadJsonString(exp)
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf("unexpected result, expected: %v\ngot: %v\ndiff: %v", exp, got, a.Diff(b).Render())
	}
}

func TestGetCardEndpoint(t *testing.T) {
	now := time.Now().UTC()
	var tests = []struct {
		name        string
		data        string
		callback    string
		user        *model.User
		response    string
		getcardf    func(ctx context.Context, req rpc.GetCardRequest) (*rpc.GetCardResponse, error)
		getuserf    func(id globalid.ID) (*model.User, error)
		getaliasf   func(id globalid.ID) (*model.AnonymousAlias, error)
		card        *model.Card
		expectError bool
		regexError  string
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "user with invalid data",
			data:        `{"foo":}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:       "invalid card id",
			callback:   "foo",
			data:       `{"cardID":""}`,
			user:       model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			regexError: `"error":"uuid.*"`,
		},
		{
			name:     "get cards error",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","error":"boom!"}`,
			getcardf: func(ctx context.Context, req rpc.GetCardRequest) (*rpc.GetCardResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":{"card":{"cardID":"","threadLevel":3,"coinsEarned":0,"title":"New Card","body":"some content...","url":"http://somesite.com/url","bgColor":"indigo-blue","background_image_path":"","anonymous":false,"post_timestamp":!timestamp_unix!},"author":{"nodeId":"!nodeID!","displayname":"rob","username":"rob","profileimg_path":"","isAnonymous":false},"replies":3,"score":0,"subscribed":false}}`,
			card:     &model.Card{Title: "New Card", OwnerID: "thisdoesntmatter", ThreadLevel: 3, Content: "some content...", URL: "http://somesite.com/url", BackgroundColor: "indigo-blue", CreatedAt: now},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.GetCardf = tt.getcardf
			if r.GetCardf == nil {
				r.GetCardf = func(ctx context.Context, req rpc.GetCardRequest) (*rpc.GetCardResponse, error) {
					if tt.card != nil {
						card := model.CardResponse{Card: tt.card.Export(), Author: tt.user.Author(), Replies: 3}
						return (*rpc.GetCardResponse)(&card), nil
					}
					return nil, nil

				}
			}

			teststore := &testStore{}
			teststore.GetUserf = tt.getuserf
			if teststore.GetUserf == nil {
				teststore.GetUserf = func(id globalid.ID) (*model.User, error) {
					return tt.user, nil
				}
			}

			teststore.GetAnonymousAliasf = tt.getaliasf
			if teststore.GetAnonymousAliasf == nil {
				teststore.GetAnonymousAliasf = func(id globalid.ID) (*model.AnonymousAlias, error) {
					return nil, nil
				}
			}
			if teststore.GetThreadCountf == nil {
				teststore.GetThreadCountf = func(id globalid.ID) (int, error) {
					return 3, nil
				}
			}

			endpoint := server.GetCardEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			tt.response = strings.Replace(tt.response, "!timestamp_unix!", fmt.Sprintf("%d", now.Unix()), -1)
			if tt.regexError != "" {
				jsonError := regexp.MustCompile(tt.regexError)
				if !jsonError.MatchString(b.String()) {
					t.Fatalf("expected response: \n\t%s\nto match:\n\t%s", b.String(), tt.regexError)
				}
				// done here
				return
			}

			got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response
			diffJSON(got, exp, t)
		})
	}
}

func TestPostCardEndpoint(t *testing.T) {
	var tests = []struct {
		nodeID, tagID, cardID, replyCardID, aliasID globalid.ID
		name                                        string
		data                                        string
		callback                                    string
		user                                        *model.User
		alias                                       *model.AnonymousAlias
		response                                    string
		expectError                                 bool
		postcardf                                   func(ctx context.Context, req rpc.PostCardRequest) (*rpc.PostCardResponse, error)
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "user with invalid data",
			data:        `{"foo":"bar"}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:     "post error",
			nodeID:   globalid.Next(),
			tagID:    globalid.Next(),
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{"nodeId":"!nodeID!","postBody":"","layoutdata":{},"tagID":"!tagID!"}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","error":"post has no content"}`,
			postcardf: func(ctx context.Context, req rpc.PostCardRequest) (*rpc.PostCardResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid",
			nodeID:   globalid.Next(),
			tagID:    globalid.Next(),
			cardID:   globalid.Next(),
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{"nodeId":"!nodeID!","content":"Lorem upsim","bgColor":"indigo-blue","tagID":"!tagID!"}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":{"card":{"cardID":"!cardID!","threadLevel":3,"coinsEarned":0,"title":"","body":"","url":"","bgColor":"indigo-blue","background_image_path":"","anonymous":false,"post_timestamp":1226358000},"author":{"nodeId":"!nodeID!","displayname":"rob","username":"rob","profileimg_path":"","isAnonymous":false}}}`,
		},
		{
			name:        "valid with reply",
			nodeID:      globalid.Next(),
			tagID:       globalid.Next(),
			cardID:      globalid.Next(),
			replyCardID: globalid.Next(),
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:        `{"nodeId":"!nodeID!","replyCardID":"!replyCardID!","content":"Lorem ipsum","tagID":"!tagID!"}`,
			callback:    "foo",
			response:    `{"rpc":"foo","ack":"!nodeID!","data":{"card":{"cardID":"!cardID!","threadLevel":3,"coinsEarned":0,"title":"","body":"","url":"","bgColor":"","background_image_path":"","anonymous":false,"post_timestamp":1226358000},"author":{"nodeId":"!nodeID!","displayname":"rob","username":"rob","profileimg_path":"","isAnonymous":false}}}`,
		},
		{
			name:        "anonymous",
			nodeID:      globalid.Next(),
			tagID:       globalid.Next(),
			cardID:      globalid.Next(),
			replyCardID: globalid.Next(),
			aliasID:     globalid.Next(),
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			alias:       &model.AnonymousAlias{Username: "coffee", DisplayName: "Anonymous"},
			data:        `{"nodeId":"!nodeID!","replyCardID":"!replyCardID!","aliasID":"!aliasID!","content":"Lorem ipsum","tagID":"!tagID!"}`,
			callback:    "foo",
			response:    `{"rpc":"foo","ack":"!nodeID!","data":{"card":{"cardID":"!cardID!","threadLevel":3,"coinsEarned":0,"title":"","body":"","url":"","bgColor":"","background_image_path":"","anonymous":false,"post_timestamp":1226358000},"author":{"nodeId":"!nodeID!","displayname":"Anonymous","username":"coffee","profileimg_path":"","isAnonymous":true}}}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.data = strings.Replace(tt.data, "!nodeID!", string(tt.nodeID), -1)
			tt.data = strings.Replace(tt.data, "!tagID!", string(tt.tagID), -1)
			tt.data = strings.Replace(tt.data, "!replyCardID!", string(tt.replyCardID), -1)
			tt.data = strings.Replace(tt.data, "!aliasID!", string(tt.aliasID), -1)
			r := &rpcMock{}
			r.PostCardf = tt.postcardf
			if r.PostCardf == nil {
				r.PostCardf = func(ctx context.Context, req rpc.PostCardRequest) (*rpc.PostCardResponse, error) {
					if tt.user == nil {
						return nil, nil
					}
					author := tt.user.Author()
					if tt.alias != nil {
						author = tt.alias.Author()
						author.ID = tt.user.ID
					}
					card := model.Card{
						ID:              tt.cardID,
						ThreadLevel:     3,
						BackgroundColor: req.Params.BackgroundColor,
						CreatedAt:       time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC),
						Author:          author,
					}
					return &rpc.PostCardResponse{
						Card:   card.Export(),
						Author: author,
					}, nil
				}
			}
			endpoint := server.PostCardEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			tt.response = strings.Replace(tt.response, "!cardID!", string(tt.cardID), -1)

			got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response
			diffJSON(got, exp, t)
		})
	}
}

func TestUpdateSettingsEndpoint(t *testing.T) {
	var tests = []struct {
		name            string
		callback        string
		data            string
		user            *model.User
		response        string
		updateSettingsf func(ctx context.Context, req rpc.UpdateSettingsRequest) (*rpc.UpdateSettingsResponse, error)
		expectError     bool
		regexError      string
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty body",
			expectError: true,
		},
		{
			name:        "user with empty data",
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "no user",
			expectError: true,
			data:        `{}`,
		},
		{
			name:       "user with invalid data",
			callback:   "callback",
			data:       `{"foo":}`,
			user:       model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			regexError: `"error":".*invalid character.*"`,
		},
		{
			name:     "update settings error",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "callback",
			response: `{"rpc":"callback","ack":"!nodeID!","error":"boom!"}`,
			updateSettingsf: func(ctx context.Context, req rpc.UpdateSettingsRequest) (*rpc.UpdateSettingsResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid",
			user:     model.NewUser("2edde20a-7e93-49ae-8b53-931a4ad40a1d", "rob", "rob@yahoo.com", "rob"),
			data:     `{"username":"pike","password":"supersecret","email":"rob@yahoo.com","cover_image_path":"","bio":"go developer","imageData":"cdn/image.gif"}`,
			callback: "callback",
			updateSettingsf: func(ctx context.Context, req rpc.UpdateSettingsRequest) (*rpc.UpdateSettingsResponse, error) {
				if req.Params.Username == nil || *req.Params.Username != "pike" {
					return nil, errors.New("unexpected username")
				}
				if req.Params.Password == nil || *req.Params.Password != "supersecret" {
					return nil, errors.New("unexpected new password")
				}
				if req.Params.Email == nil || *req.Params.Email != "rob@yahoo.com" {
					return nil, errors.New("unexpected email")
				}
				if req.Params.Bio == nil || *req.Params.Bio != "go developer" {
					return nil, errors.New("unexpected bio")
				}
				if req.Params.ImageData == nil || *req.Params.ImageData != "cdn/image.gif" {
					return nil, errors.New("unexpected image data")
				}
				user := &model.User{
					ID:               "2edde20a-7e93-49ae-8b53-931a4ad40a1d",
					Username:         "rob",
					Email:            "rob@yahoo.com",
					DisplayName:      "commander",
					Bio:              "go developer",
					ProfileImagePath: "cdn/image.gif",
					FirstName:        "Rob",
					LastName:         "Pike",
				}
				return (*rpc.UpdateSettingsResponse)(user.Export(user.ID)), nil
			},
			response: `{"rpc":"callback","ack":"!userID!","data":{"nodeId":"!userID!","displayname":"commander","firstName":"Rob","lastName":"Pike","profileimg_path":"cdn/image.gif","cover_image_path":"","userBio":"go developer","email":"rob@yahoo.com","username":"rob","botchedSignup":false}}`,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.UpdateSettingsf = tt.updateSettingsf
			if r.UpdateSettingsf == nil {
				r.UpdateSettingsf = func(ctx context.Context, req rpc.UpdateSettingsRequest) (*rpc.UpdateSettingsResponse, error) {
					return nil, nil
				}
			}
			endpoint := server.UpdateSettingsEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint: unexpected error: %v", err)
			}
			tt.response = strings.Replace(tt.response, "!userID!", string(tt.user.ID), -1)
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			if tt.regexError != "" {
				jsonError := regexp.MustCompile(tt.regexError)
				if !jsonError.MatchString(b.String()) {
					t.Fatalf("expected response: \n\t%s\nto match:\n\t%s", b.String(), tt.regexError)
				}
				// done here
				return
			}
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestInviteEndpoint(t *testing.T) {
	var tests = []struct {
		name        string
		data        string
		callback    string
		user        *model.User
		response    string
		newInvitef  func(ctx context.Context, req rpc.NewInviteRequest) (*rpc.NewInviteResponse, error)
		expectError bool
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "user with invalid data",
			data:        `{"foo":}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:     "invite error",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","error":"boom!"}`,
			newInvitef: func(ctx context.Context, req rpc.NewInviteRequest) (*rpc.NewInviteResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":{"node_id":"ad36dc78-9a68-445b-86c6-97406c87cbb4","token":"token","remaining_uses":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`,

			newInvitef: func(ctx context.Context, req rpc.NewInviteRequest) (*rpc.NewInviteResponse, error) {
				return (*rpc.NewInviteResponse)(model.NewInviteWithParams(
					globalid.ID("0ba38723-501b-42df-a9fc-b3472c2b24f8"),
					globalid.ID("ad36dc78-9a68-445b-86c6-97406c87cbb4"),
					"token",
					time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC),
					time.Date(2008, time.November, 10, 23, 0, 0, 0, time.UTC),
				)), nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.NewInvitef = tt.newInvitef
			if r.NewInvitef == nil {
				r.NewInvitef = func(ctx context.Context, req rpc.NewInviteRequest) (*rpc.NewInviteResponse, error) {
					return (*rpc.NewInviteResponse)(&model.Invite{}), nil
				}
			}
			endpoint := server.NewInviteEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestNewUserEndpoint(t *testing.T) {
	var tests = []struct {
		name        string
		data        string
		callback    string
		user        *model.User
		response    string
		newUserf    func(ctx context.Context, req rpc.NewUserRequest) (*rpc.NewUserResponse, error)
		expectError bool
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "user with invalid data",
			data:        `{"foo":}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:     "new user error",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","error":"boom!"}`,
			newUserf: func(ctx context.Context, req rpc.NewUserRequest) (*rpc.NewUserResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{"username":"commander","password":"secret","email":"andrew@yahoo.com","displayname":"pike","token":"token"}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":{"id":"0ba38723-501b-42df-a9fc-b3472c2b24f8"}}`,
			newUserf: func(ctx context.Context, req rpc.NewUserRequest) (*rpc.NewUserResponse, error) {
				if req.Params.Username != "commander" {
					return nil, errors.New("unexpected usename")
				}
				if req.Params.Password != "secret" {
					return nil, errors.New("unexpected password")
				}
				if req.Params.Email != "andrew@yahoo.com" {
					return nil, errors.New("unexpected email")
				}
				if req.Params.DisplayName != "pike" {
					return nil, errors.New("unexpected displayName")
				}
				id := globalid.ID("0ba38723-501b-42df-a9fc-b3472c2b24f8")
				return &rpc.NewUserResponse{ID: id}, nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.NewUserf = tt.newUserf
			if r.NewUserf == nil {
				r.NewUserf = func(ctx context.Context, req rpc.NewUserRequest) (*rpc.NewUserResponse, error) {
					return nil, nil
				}
			}
			endpoint := server.NewUserEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestAuthSignupEndpoint(t *testing.T) {
	var tests = []struct {
		name        string
		data        string
		callback    string
		user        *model.User
		response    string
		authf       func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error)
		expectError bool
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:        "user with empty data",
			data:        `{}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:        "user with invalid data",
			data:        `{"foo":}`,
			user:        model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			expectError: true,
		},
		{
			name:     "signup error",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","error":"boom!"}`,
			authf: func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid",
			user:     model.NewUser(globalid.Next(), "rob", "rob@yahoo.com", "rob"),
			data:     `{"username":"commander","password":"secret","token":"token","firstName":"Rob","lastName":"Pike","accessToken":"ABC123"}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":{"user":{"displayname":"","firstName":"Rob","lastName":"Pike","profileimg_path":"","cover_image_path":"","userBio":"","username":"commander","botchedSignup":false},"session":{"id":"0ba38723-501b-42df-a9fc-b3472c2b24f8"}}}`,
			authf: func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error) {
				if req.Params.Username != "commander" {
					return nil, errors.New("unexpected usename")
				}
				if req.Params.Password != "secret" {
					return nil, errors.New("unexpected password")
				}
				if req.Params.AccessToken != "ABC123" {
					return nil, errors.New("unexpected access token")
				}
				passwordHash, err := model.HashPassword(req.Params.Password)
				if err != nil {
					return nil, err
				}
				user := model.User{Username: req.Params.Username, PasswordHash: passwordHash, FirstName: req.Params.FirstName, LastName: req.Params.LastName}
				result := &rpc.AuthResponse{
					User:    user.Export(globalid.Nil),
					Session: &model.Session{ID: globalid.ID("0ba38723-501b-42df-a9fc-b3472c2b24f8")},
				}
				return result, nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.Authf = tt.authf
			if r.Authf == nil {
				r.Authf = func(ctx context.Context, req rpc.AuthRequest) (*rpc.AuthResponse, error) {
					return nil, nil
				}
			}
			endpoint := server.AuthEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			session := model.NewSession(tt.user)
			if tt.user != nil {
				pw.SetSession(session)
				ctx = context.WithValue(ctx, protocol.RequestID, tt.user.ID)
			}

			err := endpoint(ctx, session, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(tt.user.ID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}

func TestValidateUsernameEndpoint(t *testing.T) {
	var tests = []struct {
		name              string
		data              string
		username          string
		token             string
		callback          string
		response          string
		sessionID         globalid.ID
		validateusernamef func(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error)
		getinvitetokenf   func(token string) (*model.Invite, error)
		getsessionf       func(id globalid.ID) (*model.Session, error)
		expectError       bool
	}{
		{
			name:        "no body",
			expectError: true,
		},
		{
			name:        "empty data",
			data:        `{}`,
			expectError: true,
		},
		{
			name:      "signup error",
			data:      `{"token":"WCQ9W"}`,
			callback:  "foo",
			sessionID: globalid.Next(),
			response:  `{"rpc":"foo","ack":"!nodeID!","error":"boom!"}`,
			getinvitetokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{}, nil
			},
			validateusernamef: func(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error) {
				return nil, errors.New("boom!")
			},
		},
		{
			name:     "valid username",
			data:     `{"username":"gopher"}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","data":null}`,
			getinvitetokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{}, nil
			},
			validateusernamef: func(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error) {
				if req.Params.Username != "gopher" {
					return nil, fmt.Errorf("unexpected usename: %s", req.Params.Username)
				}
				return nil, nil
			},
		},
		{
			name:     "invalid username",
			data:     `{"username":"12345","token":"WCQ9W"}`,
			callback: "foo",
			response: `{"rpc":"foo","ack":"!nodeID!","error":"invalid username"}`,
			getinvitetokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{}, nil
			},
			validateusernamef: func(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error) {
				return nil, errors.New("invalid username")
			},
		},
		{
			name:      "no token but session",
			data:      `{"username":"gopher"}`,
			callback:  "foo",
			response:  `{"rpc":"foo","ack":"!nodeID!","error":"invalid username"}`,
			sessionID: globalid.Next(),
			getinvitetokenf: func(token string) (*model.Invite, error) {
				return &model.Invite{}, nil
			},
			getsessionf: func(id globalid.ID) (*model.Session, error) {
				return &model.Session{}, nil
			},
			validateusernamef: func(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error) {
				return nil, errors.New("invalid username")
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rpcMock{}
			r.ValidateUsernamef = tt.validateusernamef
			if r.ValidateUsernamef == nil {
				r.ValidateUsernamef = func(ctx context.Context, req rpc.ValidateUsernameRequest) (*rpc.ValidateUsernameResponse, error) {
					return nil, nil
				}
			}

			teststore := &testStore{}
			teststore.GetInviteByTokenf = tt.getinvitetokenf
			if teststore.GetInviteByTokenf == nil {
				teststore.GetInviteByTokenf = func(token string) (*model.Invite, error) {
					return nil, nil
				}
			}
			teststore.GetSessionf = tt.getsessionf
			if teststore.GetSessionf == nil {
				teststore.GetSessionf = func(id globalid.ID) (*model.Session, error) {
					return nil, nil
				}
			}

			endpoint := server.ValidateUsernameEndpoint(r)
			ctx := context.Background()
			if tt.callback != "" {
				ctx = context.WithValue(ctx, protocol.Callback, tt.callback)
			}
			if tt.sessionID != globalid.Nil {
				ctx = context.WithValue(ctx, protocol.SessionID, tt.sessionID)
			}
			requestID := globalid.Next()
			ctx = context.WithValue(ctx, protocol.RequestID, requestID)
			var b bytes.Buffer
			pw := protocol.NewPushWriter(&protocol.Connection{}, &b, log.NopLogger())
			err := endpoint(ctx, nil, pw, protocol.Message{Data: []byte(tt.data)})
			if tt.expectError && err == nil {
				t.Fatal("expected error")
			}
			if tt.expectError {
				// done with this test
				return
			}
			if err != nil {
				t.Fatalf("endpoint(%v): unexpected error: %v", tt.data, err)
			}
			tt.response = strings.Replace(tt.response, "!nodeID!", string(requestID), -1)
			if got, exp := strings.TrimSuffix(b.String(), "\n"), tt.response; got != exp {
				t.Errorf("unexpected result:\ngot: %s\nexp: %s\n\n", got, exp)
			}
		})
	}
}
