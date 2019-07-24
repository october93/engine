package subcmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/protocol"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Populates the database",
	Long:  "Populates the database and creates a graph",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := log.NewLogger(false, "info")
		if err != nil {
			exit(err)
		}
		initializer, err := NewInitializer(endpoint, users, skip, l)
		if err != nil {
			exit(err)
		}
		err = initializer.Run()
		if err != nil {
			exit(err)
		}
	},
}

type Initializer struct {
	client      rpc.RPC
	rootUser    *model.User
	rootSession *model.Session
	usernames   []string
	endpoint    string
	wg          sync.WaitGroup

	users int
	skip  int
}

func NewInitializer(endpoint string, users, skip int, l log.Logger) (*Initializer, error) {
	client, err := newClient(endpoint, l)
	if err != nil {
		return nil, err
	}
	return &Initializer{
		client:    client,
		users:     users,
		skip:      skip,
		usernames: make([]string, users),
		endpoint:  endpoint,
	}, err
}

func (init *Initializer) Run() error {
	err := init.login()
	if err != nil {
		return err
	}
	inviteCode, err := init.createInvite()
	if err != nil {
		return err
	}
	for i := 0; i < init.users; i++ {
		init.usernames[i] = username(i)
	}
	for i := init.skip; i < init.users; i++ {
		init.wg.Add(1)
		init.createUser(inviteCode, i)
	}
	init.wg.Wait()
	ctx := context.WithValue(context.Background(), protocol.SessionID, init.rootSession.ID)
	req := rpc.ConnectUsersRequest{Params: rpc.ConnectUsersParams{Users: init.usernames}}
	_, err = init.client.ConnectUsers(ctx, req)
	return err
}

func (init *Initializer) login() error {
	req := rpc.AuthRequest{Params: rpc.AuthParams{Username: "root", Password: "J5DM{wQZ}&Hbvjnc*$7sTe9DV&QQQWZL"}}
	resp, err := init.client.Auth(context.Background(), req)
	if err != nil {
		return err
	}
	init.rootUser = resp.User.Import()
	init.rootSession = resp.Session
	return nil
}

func (init *Initializer) createInvite() (string, error) {
	resp, err := init.client.NewInvite(context.Background(), rpc.NewInviteRequest{})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (init *Initializer) createUser(inviteCode string, i int) {
	email := fmt.Sprintf("user%d@october.news", i+1)
	firstName := "User"
	lastName := fmt.Sprintf("%d", i+1)

	ctx := context.Background()
	req := rpc.AuthRequest{Params: rpc.AuthParams{
		Username:    init.usernames[i],
		Password:    password,
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		InviteToken: inviteCode,
		IsSignup:    true,
	}}
	_, err := init.client.Auth(ctx, req)
	if err != nil && err.Error() != rpc.ErrUserAlreadyExists.Error() {
		exit(err)
	}
	init.wg.Done()
}
