package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/october93/engine/coinmanager"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/october93/engine/rpc/client"
	"github.com/october93/engine/rpc/notifications"
	"github.com/october93/engine/rpc/protocol"
	"github.com/october93/engine/rpc/push"
	"github.com/october93/engine/rpc/server"
	"github.com/october93/engine/store"
	datastore "github.com/october93/engine/store/datastore"
	"github.com/october93/engine/worker"
	"github.com/october93/engine/worker/emailsender"
	"golang.org/x/sync/errgroup"
)

type mockIndexer struct {
}

func (i *mockIndexer) IndexUser(m *model.User) error {
	return nil
}

func (i *mockIndexer) RemoveIndexForUser(m *model.User) error {
	return nil
}

func (i *mockIndexer) IndexChannel(m *model.Channel) error {
	return nil
}

func TestEngine(t *testing.T) {
	if testing.Short() || true {
		t.Skip("short testing detected. skipping test")
		return
	}
	resetWd := changeWd(t, "../../..")
	defer resetWd()

	var err error
	engine := setupEngine(t)
	defer func() {
		err = engine.store.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = datastore.DropDatabase(engine.config.Store.Datastore)
		if err != nil {
			t.Fatal(err)
		}
	}()

	rootClient := engine.client
	c := engine.config
	l := engine.log

	// use root user to create graph
	req := rpc.AuthRequest{Params: rpc.AuthParams{Username: "root", Password: c.Store.RootUserPassword}}
	_, err = rootClient.Auth(context.Background(), req)
	if err != nil {
		t.Fatalf("Login(): %v", err)
	}

	var wg sync.WaitGroup
	// create a fully connected graph of 10 users
	n := 10
	u := NewUsers(n)
	clients := make([]rpc.RPC, n)
	var g errgroup.Group
	for i := 0; i < n; i++ {
		i := i // bind i to each closure
		g.Go(func() error {
			username := fmt.Sprintf("user-%d", i+1)
			password := "secret"
			req := rpc.NewUserRequest{Params: rpc.NewUserParams{
				Username:    username,
				Password:    password,
				Email:       fmt.Sprintf("user.%d@october.news", i+1),
				DisplayName: fmt.Sprintf("User %d", i+1),
			}}
			_, err := rootClient.NewUser(context.Background(), req) // nolint: vetshadow
			if err != nil {
				return err
			}

			config := client.NewConfig()
			config.Address = fmt.Sprintf("ws://localhost:%d/deck_endpoint/", c.Server.Port)
			clients[i], err = client.NewClient(config, l)
			if err != nil {
				return err
			}

			resp, err := clients[i].Auth(context.Background(), rpc.AuthRequest{Params: rpc.AuthParams{Username: username, Password: password}})
			if err != nil {
				return err
			}
			u.SetUser(resp.User.Import(), i)
			u.SetSession(resp.Session, i)
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	_, err = rootClient.ConnectUsers(ctx, rpc.ConnectUsersRequest{Params: rpc.ConnectUsersParams{Users: u.Usernames()}})
	if err != nil {
		t.Fatalf("ConnectUsers(): %v", err)
	}

	// force initial feed build
	n = 10
	wg.Add(n - 1)
	// check user 1-10 get the card
	for i := 1; i < n; i++ {
		go func(i int) {
			req := rpc.GetCardsRequest{Params: rpc.GetCardsParams{PerPage: 10}}
			_, err := clients[i].GetCards(context.Background(), req) // nolint: vetshadow
			if err != nil {
				t.Errorf("GetCards(): %v", err)
			}

			wg.Done()
		}(i)
	}
	wg.Wait()

	time.Sleep(4 * time.Second)

	// user 1 posts a card
	card, err := clients[0].PostCard(ctx, rpc.PostCardRequest{Params: rpc.PostCardParams{AuthorID: u.User(0).ID, Content: "Lorem ipsum"}})
	if err != nil {
		t.Fatalf("PostCard(): %v", err)
	}
	time.Sleep(1 * time.Second)

	n = 10
	wg.Add(n - 1)
	// check user 1-10 get the card
	for i := 1; i < n; i++ {
		go func(i int) {
			req := rpc.GetCardsRequest{Params: rpc.GetCardsParams{PerPage: 10}}
			result, err := clients[i].GetCards(context.Background(), req) // nolint: vetshadow
			if err != nil {
				t.Errorf("GetCards(): %v", err)
			}
			cards := result.Cards
			if len(cards) != 1 {
				t.Errorf("GetCards(): expected length %d, actual: %d", 1, len(cards))
				wg.Done()
				return
			}
			if cards[0].Card.ID != card.Card.ID {
				t.Logf("GetCards(): %v", cards[0])
				t.Errorf("GetCards(): expected card %v, actual: %v", card.Card.ID, cards[0].Card.ID)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// user 1 posts a card
	_, err = clients[0].PostCard(ctx, rpc.PostCardRequest{Params: rpc.PostCardParams{AuthorID: u.User(0).ID, Content: "Article"}})
	if err != nil {
		t.Fatalf("PostCard(): %v", err)
	}
	result, err := clients[1].GetCards(context.Background(), rpc.GetCardsRequest{Params: rpc.GetCardsParams{PerPage: 10}})
	if err != nil {
		t.Fatalf("GetCards(): %v", err)
	}
	cards := result.Cards
	if len(cards) != 2 {
		t.Fatalf("GetCards(): expected length %d, actual: %d", 2, len(cards))
	}
	// check user 1-10 do not get a thread
	n = 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			req := rpc.GetThreadRequest{Params: rpc.GetThreadParams{CardID: cards[0].Card.ID}}
			resp, err := clients[i].GetThread(context.Background(), req) // nolint: vetshadow
			if err != nil {
				t.Errorf("GetThread(): %v", err)
			}
			thread := ([]*model.CardResponse)(*resp)
			if len(thread) != 0 {
				t.Errorf("GetThread(): expected length %d, actual: %d", 0, len(cards))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	// user 2 replies to the card
	reply, err := clients[1].PostCard(ctx, rpc.PostCardRequest{Params: rpc.PostCardParams{
		AuthorID:    u.User(1).ID,
		ReplyCardID: cards[0].Card.ID,
		Content:     "Comment 1",
	}})
	if err != nil {
		t.Fatalf("PostCard(): %v", err)
	}
	// user 3 replies to the reply of user 2
	reply, err = clients[2].PostCard(ctx, rpc.PostCardRequest{Params: rpc.PostCardParams{
		AuthorID:    u.User(2).ID,
		ReplyCardID: reply.Card.ID,
		Content:     "Comment 2",
	}})
	if err != nil {
		t.Fatalf("PostCard(): %v", err)
	}
	// user 4 replies to the reply of user 3
	_, err = clients[3].PostCard(ctx, rpc.PostCardRequest{Params: rpc.PostCardParams{
		AuthorID:    u.User(3).ID,
		ReplyCardID: reply.Card.ID,
		Content:     "Comment 3",
	}})
	if err != nil {
		t.Fatalf("PostCard(): %v", err)
	}
	// check user 2-10 do not get replies as new cards
	n = 10
	wg.Add(n - 1)
	for i := 1; i < n; i++ {
		go func(i int) {
			req := rpc.GetCardsRequest{Params: rpc.GetCardsParams{PerPage: 10}}
			result, err := clients[i].GetCards(context.Background(), req) // nolint: vetshadow
			if err != nil {
				t.Errorf("GetCards(): %v", err)
				wg.Done()
				return
			}
			cards := result.Cards // nolint: vetshadow
			if len(cards) != 2 {
				t.Errorf("GetCards(): expected length %d, actual: %d", 2, len(cards))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	n = 10
	wg.Add(n)
	// check user 1-10 get thread
	for i := 0; i < n; i++ {
		go func(i int) {
			req := rpc.GetThreadRequest{Params: rpc.GetThreadParams{CardID: cards[0].Card.ID}}
			resp, err := clients[i].GetThread(context.Background(), req) // nolint: vetshadow
			if err != nil {
				t.Errorf("GetThread(): %v", err)
			}
			thread := ([]*model.CardResponse)(*resp)
			if len(thread) != 3 {
				t.Errorf("GetThread(): expected length %d, actual: %d", 3, len(cards))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	invite, err := clients[2].NewInvite(context.Background(), rpc.NewInviteRequest{Params: rpc.NewInviteParams{Invites: 3}})
	if err != nil {
		t.Fatal(err)
	}
	clientConfig := client.NewConfig()
	clientConfig.Address = fmt.Sprintf("ws://localhost:%d/deck_endpoint/", c.Server.Port)
	clientt, err := client.NewClient(clientConfig, log.NopLogger())
	if err != nil {
		t.Fatal(err)
	}
	// signup
	_, err = clientt.Auth(context.Background(), rpc.AuthRequest{
		Params: rpc.AuthParams{
			Username:    "richard",
			Email:       "richard@piedpiper.com",
			FirstName:   "Richard",
			LastName:    "Hendricks",
			InviteToken: invite.Token,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := clientt.GetCards(context.Background(), rpc.GetCardsRequest{Params: rpc.GetCardsParams{PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Cards) != 1 {
		t.Fatalf("expected invitee to get a feed with %d cards, actual %d", 1, len(resp.Cards))
	}

}

func TestEngineStaging(t *testing.T) {
	ci := os.Getenv("CI")
	if ci != "true" || true {
		return
	}
	resetWd := changeWd(t, "..")
	defer resetWd()

	var err error
	config := NewConfig()
	if err = config.Load("staging.config.toml"); err != nil {
		t.Fatal(err)
	}

	// generate new Facebook test user to sign up on staging
	type testUser struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
		LoginURL    string `json:"login_url"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}

	buf := bytes.NewBuffer([]byte("installed=true&permissions=email&name=Richard%20Hendricks"))
	endpoint := fmt.Sprintf("https://graph.facebook.com/v3.0/%s/accounts/test-users?access_token=%s|%s", config.Server.RPC.FacebookAppID, config.Server.RPC.FacebookAppID, config.Server.RPC.FacebookAppSecret)
	resp, err := http.Post(endpoint, "", buf)
	if err != nil {
		t.Fatal(err)
	}

	var fbTestUser testUser
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&fbTestUser)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		r := recover()
		endpoint := fmt.Sprintf("https://graph.facebook.com/v3.0/%s?access_token=%s|%s", fbTestUser.ID, config.Server.RPC.FacebookAppID, config.Server.RPC.FacebookAppSecret)
		var req *http.Request
		req, err = http.NewRequest("DELETE", endpoint, nil)
		if err != nil {
			t.Log(err)
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Log(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Logf("expected status %d, actual %d", http.StatusOK, resp.StatusCode)
		}
		if r != nil {
			t.Fatal(r)
		}
	}()

	paul := newStagingClient(t)

	// login as Paul
	paulAuth, err := paul.Auth(context.Background(), rpc.AuthRequest{Params: rpc.AuthParams{Username: "paul", Password: "uJ@dZyVb7VJG230s"}})
	if err != nil {
		t.Fatal(err)
	}

	// create new invite
	invite, err := paul.NewInvite(context.Background(), rpc.NewInviteRequest{Params: rpc.NewInviteParams{Invites: 1}})
	if err != nil {
		t.Fatal(err)
	}

	c := newStagingClient(t)
	// test validate invite code before signing up
	_, err = c.ValidateInviteCode(context.Background(), rpc.ValidateInviteCodeRequest{Params: rpc.ValidateInviteCodeParams{
		Token: invite.Token,
	}})
	if err != nil {
		t.Fatal(err)
	}

	username := fmt.Sprintf("richardhendricks%d", rand.Int())

	// sign up test user with staging
	testUserAuth, err := c.Auth(context.Background(), rpc.AuthRequest{Params: rpc.AuthParams{
		AccessToken: fbTestUser.AccessToken,
		InviteToken: invite.Token,
		Username:    username,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if testUserAuth.User.Email != fbTestUser.Email {
		t.Fatalf("expected email %s, actual %s", fbTestUser.Email, testUserAuth.User.Email)
	}
	if testUserAuth.User.Username != username {
		t.Fatalf("expected username %s, actual %s", username, testUserAuth.User.Username)
	}

	// get user should reflect the logged in user
	userResp, err := c.GetUser(context.Background(), rpc.GetUserRequest{Params: rpc.GetUserParams{Username: username}})
	if err != nil {
		t.Fatal(err)
	}
	if userResp.FirstName != "Richard" {
		t.Fatalf("expected first name %s, actual %s", "Richard", userResp.FirstName)
	}
	if userResp.LastName != "Hendricks" {
		t.Fatalf("expected last name %s, actual %s", "Hendricks", userResp.LastName)
	}
	if userResp.Email == "" {
		t.Fatal("expected email not to be empty")
	}

	/*
		// get features should not be empty
		featureResp, err := c.GetFeaturesForUser(context.Background(), rpc.GetFeaturesForUserRequest{})
		if err != nil {
			t.Fatal(err)
		}

		features := *(*[]string)(featureResp)
		if len(features) == 0 {
			t.Fatal("expected features not to be empty")
		}
	*/

	// get invites should not fail
	_, err = c.GetInvites(context.Background(), rpc.GetInvitesRequest{})
	if err != nil {
		t.Fatal(err)
	}

	// newly signed up user should have notifications
	notificationsResp, err := c.GetNotifications(context.Background(), rpc.GetNotificationsRequest{Params: rpc.GetNotificationsParams{
		PageSize: 20,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(notificationsResp.Notifications) == 0 {
		t.Fatal("expected newly signed up user to have notifications")
	}
	if notificationsResp.UnseenCount == 0 {
		t.Fatal("expected newly signed up user to have unseen notifications")
	}

	// newly signed up user should have onboarding data
	onboardingResp, err := c.GetOnboardingData(context.Background(), rpc.GetOnboardingDataRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(onboardingResp.NetworkProfilePictures) == 0 {
		t.Fatal("expected newly signed up user to have users in their network")
	}
	if onboardingResp.InvitingUser.Username != "paul" {
		t.Fatal("expected paul to be the inviting user")
	}

	// newly signed up user should have a working networking tab
	networkResp, err := c.GetMyNetwork(context.Background(), rpc.GetMyNetworkRequest{Params: rpc.GetMyNetworkParams{
		PageSize: 20,
	}})
	if err != nil {
		t.Fatal(err)
	}

	if len(networkResp.Users) != 20 {
		t.Fatalf("expected 20 network users, actual %d", len(networkResp.Users))
	}

	// newly signed up user should have a feed
	_, err = c.GetCards(context.Background(), rpc.GetCardsRequest{Params: rpc.GetCardsParams{PerPage: 20}})
	if err != nil {
		t.Fatal(err)
	}
	// TODO: fix test
	//	if len(cardsResp.Cards) != 20 {
	//		t.Fatal("expected feed to consist of 20 cards")
	//	}

	// Paul should be able to introduce newly signed up user
	postCardResp, err := paul.PostCard(context.Background(), rpc.PostCardRequest{Params: rpc.PostCardParams{
		AuthorID: paulAuth.User.ID,
		Content:  fmt.Sprintf("Welcome @%s (Richard Hendricks). Richard is the CEO and and creator of Pied Piper", username),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if postCardResp.Author.Username != "paul" {
		t.Fatalf("expected author's username of introduction post to be paul, actual %s", postCardResp.Author.Username)
	}

	// newly signed up user should be able to boost card
	_, err = c.ReactToCard(context.Background(), rpc.ReactToCardRequest{Params: rpc.ReactToCardParams{
		CardID:   postCardResp.Card.ID,
		Reaction: model.Like,
	}})
	if err != nil {
		t.Fatal(err)
	}

	// newly signed up user should be able to comment on introduction post
	postResp, err := c.PostCard(context.Background(), rpc.PostCardRequest{Params: rpc.PostCardParams{
		AuthorID:    testUserAuth.User.ID,
		ReplyCardID: postCardResp.Card.ID,
		Content:     "Thanks Paul. Let's talk about how to integrate October into our new distributed(!) Internet",
	}})
	if err != nil {
		t.Fatal(err)
	}

	// get thread should reflect that the comment has worked
	threadResp, err := c.GetThread(context.Background(), rpc.GetThreadRequest{Params: rpc.GetThreadParams{CardID: postCardResp.Card.ID}})
	if err != nil {
		t.Fatal(err)
	}
	cards := *(*[]*model.CardResponse)(threadResp)
	if len(cards) != 1 {
		t.Fatalf("expected get thread to return 1 card, actual %d", len(cards))
	}

	// clean up post again
	_, err = c.DeleteCard(context.Background(), rpc.DeleteCardRequest{Params: rpc.DeleteCardParams{CardID: postResp.Card.ID}})
	if err != nil {
		t.Fatal(err)
	}

	// changed email address through update settings
	//	email := "octoberci@outlook.com"
	//	updateResp, err := c.UpdateSettings(context.Background(), rpc.UpdateSettingsRequest{Params: rpc.UpdateSettingsParams{
	//		Email: &email,
	//	}})
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	if updateResp.Email != email {
	//		t.Fatalf("expected updated email to be %s, actual %s", email, updateResp.Email)
	//	}

	// test reset password functionality
	/*
		_, err = c.ResetPassword(context.Background(), rpc.ResetPasswordRequest{Params: rpc.ResetPasswordParams{
			Email: email,
		}})
		if err != nil {
			t.Fatal(err)
		}

		content := fetchEmailContent(t)

		split := strings.Split(content, `<a href="`)
		split = strings.Split(split[1], `"`)
		sendGridURL := split[0]

		httpClient := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}

		resp, err = httpClient.Head(sendGridURL)
		if err != nil {
			t.Fatal(err)
		}
		resetURL, err := resp.Location()
		if err != nil {
			t.Fatal(err)
		}
		split = strings.Split(resetURL.String(), "token=")
		split = strings.Split(split[1], "&email=")
		resetToken := split[0]
		email = split[1]
		email, err = url.QueryUnescape(split[1])
		if err != nil {
			t.Fatal(err)
		}

		c2 := newStagingClient(t)
		authResp, err := c2.Auth(context.Background(), rpc.AuthRequest{Params: rpc.AuthParams{
			Username:   email,
			ResetToken: resetToken,
		}})
		if err != nil {
			t.Fatal(err)
		}

		if authResp.User.Email != email {
			t.Fatalf("expected email %s, actual %s", email, authResp.User.Email)
		}
		if authResp.User.Username != "richardhendricks" {
			t.Fatalf("expected username %s, actual %s", "richardhendricks", authResp.User.Username)
		}
	*/
}

func changeWd(t *testing.T, path string) func() {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir(path)
	if err != nil {
		t.Fatal(err)
	}
	return func() {
		err = os.Chdir(wd)
		if err != nil {
			t.Fatal(err)
		}
	}
}

/*
func fetchEmailContent(t *testing.T) string {
	// check emails
	emailClient, err := imapclient.DialTLS("outlook.office365.com:993", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer emailClient.Logout()

	err = emailClient.Login("octoberci@outlook.com", "*j7eTP$]Cfel}d0ZClK45{aO1yKxq26B5K8L")
	if err != nil {
		t.Fatal(err)
	}

	attempts := 0
	maxAttempts := 5

	var mbox *imap.MailboxStatus
	for {
		mbox, err = emailClient.Select("INBOX", false)
		if err != nil {
			t.Fatal(err)
		}

		if mbox.Messages >= 1 {
			break
		}

		attempts++
		if attempts == maxAttempts {
			t.Fatal("no messages in the mailbox of octoberci@outlook.com")
		}
		time.Sleep(time.Duration(attempts*2) * time.Second)
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(mbox.Messages)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- emailClient.Fetch(seqset, items, messages)
	}()

	var content string
	for msg := range messages {
		r := msg.GetBody(section)
		if r == nil {
			t.Fatal("Server didn't returned message body")
		}

		// Create a new mail reader
		var mr *mail.Reader
		mr, err = mail.CreateReader(r)
		if err != nil {
			t.Fatal(err)
		}

		// Process each message's part
		for {
			var p *mail.Part
			p, err = mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				t.Fatal(err)
			}

			switch p.Header.(type) {
			case mail.TextHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				content = string(b)
			}
		}
	}

	// delete message
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if err = emailClient.Store(seqset, item, flags, nil); err != nil {
		t.Fatal(err)
	}

	if err = <-done; err != nil {
		t.Fatal(err)
	}
	return content
}
*/
func newStagingClient(t *testing.T) rpc.RPC {
	t.Helper()

	config := client.NewConfig()
	config.Address = "wss://engine.staging.october.news/deck_endpoint/"
	c, err := client.NewClient(config, log.NopLogger())
	if err != nil {
		t.Fatal(err)
	}
	return c
}

type engine struct {
	store  *store.Store
	config *Config
	client rpc.RPC
	log    log.Logger
}

func setupEngine(t *testing.T) engine {
	t.Helper()
	config := NewConfig()
	if err := config.Load("test.config.toml"); err != nil {
		t.Fatal(err)
	}
	config.Server.Port = freePort(t)
	logger, err := log.NewLogger(false, log.Info)
	if err != nil {
		t.Fatal(err)
	}
	err = datastore.DropDatabase(config.Store.Datastore)
	if err != nil {
		t.Fatal(err)
	}
	err = datastore.SetupDatabase(config.Store.Datastore)
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.NewStore(&config.Store, logger)
	if err != nil {
		t.Fatal(err)
	}

	err = st.EnsureRootUser()
	if err != nil {
		t.Fatal(err)
	}

	notifier, err := worker.NewNotifier(st, &config.Worker, logger)
	if err != nil {
		t.Fatal(err)
	}
	w := struct {
		*worker.Notifier
		*worker.SessionCleaner
		*emailsender.EmailSender
	}{
		Notifier:       nil,
		SessionCleaner: nil,
		EmailSender:    nil,
	}

	conns := protocol.NewConnections(nil, log.NopLogger())
	rt := protocol.NewRouter(st, conns, &model.Settings{}, &protocol.Config{}, log.NopLogger())
	p, err := push.NewPusher(rt.Connections(), st, &push.Config{}, log.NopLogger())
	if err != nil {
		t.Fatal(err)
	}
	ip, err := rpc.NewImageProcessor(&config.Server.RPC, log.NopLogger())
	if err != nil {
		t.Fatal(err)
	}

	ns := notifications.NewNotifications(st, "", 10000)

	cm := coinmanager.NewCoinManager(st, &config.CoinManager)

	resps := rpc.NewResponses(st)

	r := rpc.NewRPC(st, w, nil, ip, &config.Server.RPC, logger, nil, notifier, p, &model.Settings{}, ns, &mockIndexer{}, cm, resps)
	err = server.NewServer(st, nil, r, rt, nil, &config.Server, logger, context.Background()).Open()
	if err != nil {
		t.Fatal(err)
	}
	clientConfig := client.NewConfig()
	clientConfig.Address = fmt.Sprintf("ws://localhost:%d/deck_endpoint/", config.Server.Port)
	clientt, err := client.NewClient(clientConfig, logger)
	if err != nil {
		t.Fatal(err)
	}
	return engine{
		store:  st,
		client: clientt,
		config: config,
		log:    logger,
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen(): %v", err)
	}
	defer func() {
		err = l.Close()
		if err != nil {
			t.Fatalf("Close(): %v", err)
		}
	}()
	return l.Addr().(*net.TCPAddr).Port

}

type Users struct {
	mu       sync.RWMutex
	users    []*model.User
	sessions []*model.Session
}

func NewUsers(n int) *Users {
	return &Users{
		users:    make([]*model.User, n),
		sessions: make([]*model.Session, n),
	}
}

func (u *Users) SetUser(user *model.User, i int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users[i] = user
}

func (u *Users) User(i int) *model.User {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.users[i]
}

func (u *Users) SetSession(session *model.Session, i int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.sessions[i] = session
}

func (u *Users) Session(i int) *model.Session {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.sessions[i]
}

func (u *Users) Usernames() []string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	result := make([]string, len(u.users))
	for i, user := range u.users {
		result[i] = user.Username
	}
	return result
}
