package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	fcm "github.com/edganiukov/fcm"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store"
)

// Notifier manages the notifications to be send out to the user. It
// identifies each node with a list of possible devices.
type Notifier struct {
	FirebaseClient *fcm.Client
	Store          *store.Store
	SlackWebhook   string
	log            log.Logger
	pendingPushMap PendingPushMap
}

// NewNotifier returns a NotificationWorker to issue notifications.
func NewNotifier(s *store.Store, config *Config, l log.Logger) (*Notifier, error) {
	var c *fcm.Client

	// only configure fcm if there's a key
	if config.FCMServerKey != "" {
		var err error
		c, err = fcm.NewClient(config.FCMServerKey)
		if err != nil {
			return nil, err
		}
	}

	return &Notifier{FirebaseClient: c, Store: s, SlackWebhook: config.SlackWebhook, log: l, pendingPushMap: PendingPushMap{pendingIDs: make(map[globalid.ID]bool)}}, nil
}

func (n *Notifier) NotifySlackAboutCard(card, reply *model.Card, author *model.Author) error {
	if n.SlackWebhook == "" {
		return nil
	}

	var text string
	if reply == nil {
		text = fmt.Sprintf("%s created a new post (<https://october.news/post/%v|Link>)", author.DisplayName, card.ID)
	} else {
		replID := reply.ThreadRootID
		if replID == globalid.Nil {
			replID = reply.ID
		}
		text = fmt.Sprintf("%s commented on %s's post (<https://october.news/post/%v|Link>)", author.DisplayName, reply.Author.DisplayName, replID)
	}
	return n.NotifySlack("engagement", text)
}

func (n *Notifier) NotifySlack(channel, text string) error {
	if n.SlackWebhook == "" {
		return nil
	}
	payload := struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
	}{
		Channel: channel,
		Text:    text,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(n.SlackWebhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

type PendingPushMap struct {
	sync.RWMutex
	pendingIDs map[globalid.ID]bool
}

func (p *PendingPushMap) IsPending(id globalid.ID) bool {
	return p.pendingIDs[id]
}

func (p *PendingPushMap) SetPending(id globalid.ID, to bool) {
	p.Lock()
	defer p.Unlock()
	p.pendingIDs[id] = to
}

// Notify is used to notify a given user (node) about a new card.
func (n *Notifier) NotifyPush(eN *model.ExportedNotification) error {
	if n.Store == nil {
		return nil
	}
	user, err := n.Store.GetUser(eN.UserID)

	if user.BlockedAt.Valid {
		// the user is blocked, return with no error
		return nil
	}

	if err != nil {
		return err
	}

	pushNotification := &fcm.Notification{
		Badge: "1",
		Body:  eN.PlainMessage(),
		Icon:  "ic_stat",
	}

	pushData := pushNotificationData{
		"type":       eN.Type,
		"id":         eN.ID,
		"action":     eN.Action,
		"actionData": eN.ActionData,
	}

	for _, d := range user.Devices {
		go func(d model.Device) {
			err := n.sendNotification(d, pushNotification, pushData)
			if err != nil {
				n.log.Error(err,
					"username", user.Username,
					"token", d.Token,
					"platform", d.Platform,
				)
				if err == fcm.ErrNotRegistered {
					user.PossibleUninstall = true
					devices := user.Devices
					delete(user.Devices, d.Token)
					user.Devices = devices
					serr := n.Store.SaveUser(user)
					if serr != nil {
						n.log.Error(err)
					}

					if len(user.Devices) <= 0 {
						serr = n.NotifySlack("engagement", fmt.Sprintf("<@U0URBUX2L> User %v (%v) may have uninstalled october!", user.DisplayName, user.Username))
						if serr != nil {
							n.log.Error(err)
						}
					}
				}
			}
		}(d)
	}
	return nil
}

type pushNotificationData map[string]interface{}

func (n *Notifier) sendNotification(device model.Device, pushNotification *fcm.Notification, data pushNotificationData) error {
	if n.FirebaseClient == nil {
		return nil
	}
	msg := &fcm.Message{
		Token:            device.Token,
		Notification:     pushNotification,
		Priority:         "high",
		ContentAvailable: true,
		Data:             data,
	}
	response, err := n.FirebaseClient.Send(msg)
	if err != nil {
		return err
	}

	// TODO (konrad) check for new registration_id
	if response.Failure == 1 {
		return response.Results[0].Error
	}

	return nil
}
