package subcmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc"
	"github.com/spf13/cobra"
)

var (
	load         int
	postComments bool
	getCards     bool
)

func init() {
	RootCmd.PersistentFlags().IntVar(&load, "load", 10, "number of go routines performing requests")
	loadCmd.Flags().BoolVarP(&postComments, "comments", "c", false, "Post comments instead of top-level posts")
	loadCmd.Flags().BoolVarP(&getCards, "getCards", "g", false, "Run get cards concurrently")
	RootCmd.AddCommand(loadCmd)
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Starts a load test",
	Long:  "Starts a load test with Engine by posting continiously cards",
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			err := http.ListenAndServe("localhost:7070", nil)
			if err != nil {
				exit(err)
			}
		}()
		l, err := log.NewLogger(false, "info")
		if err != nil {
			exit(err)
		}
		lt, err := NewLoadTester(endpoint, users, load, skip, l)
		if err != nil {
			exit(err)
		}
		if err := lt.Run(); err != nil {
			exit(err)
		}
	},
}

type LoadTester struct {
	clients  []rpc.RPC
	users    []*model.User
	sessions []*model.Session
	wg       sync.WaitGroup

	skip     int
	numUsers int
	load     int

	loginLatency           Latency
	postCardLatency        Latency
	postCommentLatency     Latency
	getCardsLatency        Latency
	reactToCardLatency     Latency
	voteOnCardLatency      Latency
	getPopularCardsLatency Latency
}

type Latency struct {
	mu          sync.Mutex
	requests    int64
	nanoseconds int64
	rps         int64
	last        time.Duration
	start       time.Time
}

func (l *Latency) Update(d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.nanoseconds += d.Nanoseconds()
	l.last = d
	l.requests += 1
	l.rps = int64(float64(l.requests) / time.Since(l.start).Seconds())
}

func (l *Latency) Average() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.requests == 0 {
		return time.Duration(0)
	}
	return time.Duration(l.nanoseconds / l.requests)
}

func (l *Latency) StartTimer() {
	l.start = time.Now()
}

func NewLoadTester(endpoint string, users, load, skip int, l log.Logger) (*LoadTester, error) {
	lt := &LoadTester{
		clients:  make([]rpc.RPC, users),
		users:    make([]*model.User, users),
		sessions: make([]*model.Session, users),
		skip:     skip,
		numUsers: users,
		load:     load,
	}
	for i := 0; i < users; i++ {
		client, err := newClient(endpoint, l)
		if err != nil {
			return nil, err
		}
		lt.clients[i] = client
	}
	return lt, nil
}

func (lt *LoadTester) Run() error {
	for i := lt.skip; i < lt.numUsers; i++ {
		lt.wg.Add(1)
		go lt.login(i)
	}
	lt.wg.Wait()
	go lt.printLatency()
	for i := 0; i < lt.load; i++ {
		lt.wg.Add(1)

		if !postComments {
			go lt.postCard((skip + i) % lt.numUsers)
		}

		num := (skip + i) % lt.numUsers
		if getCards {
			go lt.getCardsLoop(num)
			go lt.getPopularCards(num)
		}

		cards := lt.getCards(i)
		go lt.postComment((skip+i)%lt.numUsers, cards.Cards)
		go lt.reactToCards((skip+i)%lt.numUsers, cards.Cards)
		go lt.voteOnCards((skip+i)%lt.numUsers, cards.Cards)
	}
	lt.wg.Wait()
	return nil
}

func (lt *LoadTester) login(i int) {
	defer func(begin time.Time) {
		lt.loginLatency.Update(time.Since(begin))

	}(time.Now())

	req := rpc.AuthRequest{Params: rpc.AuthParams{
		Username: username(i),
		Password: password,
	}}
	resp, err := lt.clients[i].Auth(context.Background(), req)
	if err != nil {
		exit(err)
	}
	lt.users[i] = resp.User.Import()
	lt.sessions[i] = resp.Session
	lt.wg.Done()
}

func (lt *LoadTester) postCard(i int) {
	lt.postCardLatency.StartTimer()
	for {
		begin := time.Now()
		req := rpc.PostCardRequest{
			Params: rpc.PostCardParams{
				AuthorID: lt.users[i].ID,
				Content:  "Lorem ipsum",
			}}
		_, err := lt.clients[i].PostCard(context.Background(), req)
		if err != nil {
			exit(fmt.Errorf("PostCard(): %v", err))
		}
		lt.postCardLatency.Update(time.Since(begin))
	}
}

func (lt *LoadTester) postComment(i int, cards []*model.CardResponse) {
	lt.postCommentLatency.StartTimer()
	for {
		begin := time.Now()
		req := rpc.PostCardRequest{
			Params: rpc.PostCardParams{
				AuthorID:    lt.users[i].ID,
				Content:     "Lorem ipsum",
				ReplyCardID: cards[0].Card.ID,
			}}
		_, err := lt.clients[i].PostCard(context.Background(), req)
		if err != nil {
			exit(fmt.Errorf("PostCard(): %v", err))
		}
		lt.postCommentLatency.Update(time.Since(begin))
	}
}

func (lt *LoadTester) getCardsLoop(i int) {
	lt.getCardsLatency.StartTimer()
	for {
		begin := time.Now()
		req := rpc.GetCardsRequest{
			Params: rpc.GetCardsParams{
				PerPage: 20,
				Page:    2,
			},
		}
		_, err := lt.clients[i].GetCards(context.Background(), req)
		if err != nil {
			exit(fmt.Errorf("GetCards(): %v", err))
		}
		lt.getCardsLatency.Update(time.Since(begin))
	}
}

func (lt *LoadTester) getPopularCards(i int) {
	lt.getPopularCardsLatency.StartTimer()
	for {
		begin := time.Now()
		req := rpc.GetPopularCardsRequest{
			Params: rpc.GetPopularCardsParams{
				PerPage: 20,
				Page:    2,
			},
		}
		_, err := lt.clients[i].GetPopularCards(context.Background(), req)
		if err != nil {
			exit(fmt.Errorf("GetPopularCards(): %v", err))
		}
		lt.getPopularCardsLatency.Update(time.Since(begin))
	}
}

func (lt *LoadTester) getCards(i int) *rpc.GetCardsResponse {
	lt.getCardsLatency.StartTimer()
	begin := time.Now()
	req := rpc.GetCardsRequest{
		Params: rpc.GetCardsParams{
			PerPage: 10,
			Page:    2,
		},
	}
	result, err := lt.clients[i].GetCards(context.Background(), req)
	if err != nil {
		exit(fmt.Errorf("GetCards(): %v", err))
	}
	lt.getCardsLatency.Update(time.Since(begin))
	return result
}

func (lt *LoadTester) reactToCards(i int, cards []*model.CardResponse) {
	for _, card := range cards {
		begin := time.Now()
		req := rpc.ReactToCardRequest{Params: rpc.ReactToCardParams{CardID: card.Card.ID, Reaction: model.Like}}
		_, err := lt.clients[i].ReactToCard(context.Background(), req)
		if err != nil {
			exit(fmt.Errorf("ReactToCard(): %v", err))
		}
		lt.reactToCardLatency.Update(time.Since(begin))
	}
}

func (lt *LoadTester) voteOnCards(i int, cards []*model.CardResponse) {
	for _, card := range cards {
		begin := time.Now()
		req := rpc.VoteOnCardRequest{Params: rpc.VoteOnCardParams{CardID: card.Card.ID, Type: model.Up}}
		_, err := lt.clients[i].VoteOnCard(context.Background(), req)
		if err != nil {
			exit(fmt.Errorf("VoteOnCard(): %v", err))
		}
		lt.voteOnCardLatency.Update(time.Since(begin))
	}
}

func (lt *LoadTester) printLatency() {
	for {
		time.Sleep(3 * time.Second)
		fmt.Printf("Login: %v Requests: %d req/second: %d\n", lt.loginLatency.last, lt.loginLatency.requests, lt.loginLatency.rps)
		fmt.Printf("PostCard: %v Requests: %d req/second: %d\n", lt.postCardLatency.last, lt.postCardLatency.requests, lt.postCardLatency.rps)
		fmt.Printf("PostComment: %v Requests: %d req/second: %d\n", lt.postCommentLatency.last, lt.postCommentLatency.requests, lt.postCommentLatency.rps)
		fmt.Printf("GetCards: %v Requests: %d req/second: %d\n", lt.getCardsLatency.last, lt.getCardsLatency.requests, lt.getCardsLatency.rps)
		fmt.Printf("ReactToCard: %v Requests: %d req/second: %d\n", lt.reactToCardLatency.last, lt.reactToCardLatency.requests, lt.reactToCardLatency.rps)
		fmt.Printf("VoteOnCard: %v Requests: %d req/second: %d\n", lt.voteOnCardLatency.last, lt.voteOnCardLatency.requests, lt.voteOnCardLatency.rps)
		fmt.Printf("GetPopularCardsLatency: %v Requests: %d req/second: %d\n\n", lt.getPopularCardsLatency.last, lt.getPopularCardsLatency.requests, lt.getPopularCardsLatency.rps)
	}
}
