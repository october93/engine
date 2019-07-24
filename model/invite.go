package model

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/october93/engine/kit/globalid"
)

const (
	inviteTokenLength = 5
)

var ErrInvalidInviteCode = errors.New("invalid or expired invite code")

type Invite struct {
	ID            globalid.ID `db:"id" json:"-"`
	NodeID        globalid.ID `db:"node_id" json:"node_id"`
	ChannelID     globalid.ID `db:"channel_id" json:"-"`
	Token         string      `db:"token" json:"token"`
	RemainingUses int         `db:"remaining_uses" json:"remaining_uses"`
	HideFromUser  bool        `db:"hide_from_user" json:"-"`
	SystemInvite  bool        `db:"system_invite" json:"-"`
	GroupID       globalid.ID `db:"group_id" json:"-"`
	CreatedAt     time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at"`
}

// NewInvite returns a new Invite instance.
func NewInvite(nodeID globalid.ID) (*Invite, error) {
	token, err := randomHumanToken(inviteTokenLength)
	if err != nil {
		return nil, err
	}
	return &Invite{
		ID:            globalid.Next(),
		NodeID:        nodeID,
		Token:         token,
		RemainingUses: 1,
	}, nil
}

func NewInviteWithParams(id, nodeID globalid.ID, token string, createdAt, updatedAt time.Time) *Invite {
	return &Invite{
		ID:            id,
		NodeID:        nodeID,
		Token:         token,
		RemainingUses: 1,
	}
}

func randomHumanToken(n int) (string, error) {
	// omit 0, 1 and 8 since it can be confused with O, I and B
	// omit L, I and O since it can be confused with 1 and 0
	// omit U for accidental obscenity
	var alphabet = []rune("2345679ABCDEFGHJKMNPQRSTVWXYZ")
	b := make([]rune, n)
	for i := range b {
		max := big.NewInt(int64(len(alphabet)))
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[r.Int64()]
	}
	return string(b), nil
}
