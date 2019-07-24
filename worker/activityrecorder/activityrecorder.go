package activityrecorder

import (
	"encoding/json"
	"sync"

	nsq "github.com/bitly/go-nsq"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store"
	"github.com/october93/engine/worker"
)

type ActivityRecorder struct {
	wg       *sync.WaitGroup
	wc       *worker.Config
	sc       *store.Config
	store    *store.Store
	producer *nsq.Producer
	log      log.Logger
}

func NewActivityRecorder(wc *worker.Config, sc *store.Config, log log.Logger) (*ActivityRecorder, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(wc.NSQDAddress, config)
	if err != nil {
		return nil, err
	}
	return &ActivityRecorder{
		wg:       &sync.WaitGroup{},
		wc:       wc,
		sc:       sc,
		producer: producer,
		log:      log,
	}, nil
}

func NewActivityConsumer(wc *worker.Config, sc *store.Config, log log.Logger) *ActivityRecorder {
	return &ActivityRecorder{
		wg:  &sync.WaitGroup{},
		wc:  wc,
		sc:  sc,
		log: log,
	}
}

func (ar *ActivityRecorder) EnqueueJob(activity *model.Activity) error {
	body, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	return ar.producer.Publish("activity", body)
}

func (ar *ActivityRecorder) ConsumeJobs() error {
	ar.wg.Add(1)

	var err error
	ar.store, err = store.NewStore(ar.sc, ar.log)
	if err != nil {
		return err
	}
	config := nsq.NewConfig()
	q, err := nsq.NewConsumer("activity", "default", config)
	if err != nil {
		return err
	}
	q.AddHandler(nsq.HandlerFunc(ar.handleJob))
	err = q.ConnectToNSQLookupd(ar.wc.NSQLookupdAddress)
	if err != nil {
		return err
	}
	ar.wg.Wait()
	return nil
}

func (ar *ActivityRecorder) handleJob(message *nsq.Message) error {
	var activity *model.Activity
	err := json.Unmarshal(message.Body, &activity)
	if err != nil {
		return err
	}
	return ar.store.SaveActivity(activity)
}
