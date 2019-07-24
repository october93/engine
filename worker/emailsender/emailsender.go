package emailsender

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"sync"

	nsq "github.com/bitly/go-nsq"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/worker"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSender struct {
	producer *nsq.Producer
	config   *worker.Config
	log      log.Logger
	wg       *sync.WaitGroup
}

type Job struct {
	TextTemplate string      `json:"textTemplate"`
	HTMLTemplate string      `json:"htmlTemplate"`
	Subject      string      `json:"subject"`
	Sender       string      `json:"sender"`
	From         string      `json:"from"`
	To           string      `json:"to"`
	Recipient    string      `json:"recipient"`
	Data         interface{} `json:"data"`
}

// NewEmailSender returns a new instance of EmailWorker.
func NewEmailSender(config *worker.Config, log log.Logger) (*EmailSender, error) {
	producer, err := nsq.NewProducer(config.NSQDAddress, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return &EmailSender{
		producer: producer,
		wg:       &sync.WaitGroup{},
		config:   config,
		log:      log,
	}, nil
}

func NewEmailConsumer(config *worker.Config, log log.Logger) *EmailSender {
	return &EmailSender{
		wg:     &sync.WaitGroup{},
		config: config,
		log:    log,
	}
}

func (ew *EmailSender) EnqueueMailJob(job *Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	ew.log.Info("enqueuing email job", "recipient", job.To)
	return ew.producer.Publish("email", body)
}

func (ew *EmailSender) ConsumeJobs() error {
	ew.wg.Add(1)

	config := nsq.NewConfig()
	q, err := nsq.NewConsumer("email", "default", config)
	if err != nil {
		return err
	}
	q.AddHandler(ew.logMessage(ew.handleEmailJob))
	err = q.ConnectToNSQLookupd(ew.config.NSQLookupdAddress)
	if err != nil {
		return err
	}
	ew.wg.Wait()
	return nil
}

func (ew *EmailSender) logMessage(handler func(message *nsq.Message) error) nsq.HandlerFunc {
	return func(message *nsq.Message) error {
		ew.log.Info("processing message", "id", message.ID)
		err := handler(message)
		if err != nil {
			ew.log.Error(err, "id", message.ID)
		}
		return err
	}
}

func (ew *EmailSender) handleEmailJob(message *nsq.Message) error {
	var job Job
	err := json.Unmarshal(message.Body, &job)
	if err != nil {
		return err
	}

	var textContentBuf bytes.Buffer
	err = renderTemplate(&textContentBuf, "text/plain", job.TextTemplate, job.Data)
	if err != nil {
		return err
	}
	var htmlContentBuf bytes.Buffer
	err = renderTemplate(&htmlContentBuf, "text/html", job.HTMLTemplate, job.Data)
	if err != nil {
		return err
	}

	from := mail.NewEmail(job.Sender, job.From)
	to := mail.NewEmail(job.Recipient, job.To)
	textContent := mail.NewContent("text/plain", textContentBuf.String())
	htmlContent := mail.NewContent("text/html", htmlContentBuf.String())

	m := mail.NewV3MailInit(from, job.Subject, to, textContent, htmlContent)
	client := sendgrid.NewSendClient(ew.config.SendGridAPIKey)
	response, err := client.Send(m)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func renderTemplate(w io.Writer, mime, tmpl string, p interface{}) error {
	t, err := template.New(mime).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, p)
}

func (ew *EmailSender) Shutdown() error {
	ew.wg.Done()
	return nil
}
