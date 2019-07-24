package datastore

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// SaveInvite saves an invite to the store
func (store *Store) SaveInvite(m *model.Invite) error {
	return saveInvite(store, m)
}

// SaveInvite saves an invite to the store
func (tx *Tx) SaveInvite(m *model.Invite) error {
	return saveInvite(tx, m)
}

func saveInvite(e sqlx.Ext, m *model.Invite) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	if m.ID == globalid.Nil {
		m.ID = globalid.Next()
	}
	tn := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = tn
	}
	m.UpdatedAt = tn

	_, err := sqlx.NamedExec(e,
		`INSERT INTO invites
		(
			id,
			node_id,
			token,
			remaining_uses,
			hide_from_user,
			group_id,
			system_invite,
			channel_id,
			created_at,
			updated_at
		)
		VALUES
		(
			:id,
			:node_id,
			:token,
			:remaining_uses,
			:hide_from_user,
			:group_id,
			:system_invite,
			:channel_id,
			:created_at,
			:updated_at
		)
		ON CONFLICT(id) DO UPDATE
		SET
			id          = :id,
			node_id     = :node_id,
			token       = :token,
			remaining_uses = :remaining_uses,
			hide_from_user = :hide_from_user,
			group_id = :group_id,
			channel_id = :channel_id,
			system_invite = :system_invite,
			created_at  = :created_at,
			updated_at   = :updated_at
		WHERE invites.id = :id`, m)
	return errors.Wrap(err, "SaveInvite failed")
}

// GetInvites returns all invites from the store
func (store *Store) GetInvites() ([]*model.Invite, error) {
	invites := []*model.Invite{}
	err := store.Select(&invites, "SELECT * FROM invites ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}

	return invites, errors.Wrap(err, "GetInvites failed")
}

// GetInvites returns all invites from the store
func (store *Store) GetInviteForChannelAndNode(channelID, nodeID globalid.ID) (*model.Invite, error) {
	var invites model.Invite
	err := store.Get(&invites, "SELECT * FROM invites WHERE channel_id = $1 AND node_id = $2 ORDER BY updated_at DESC LIMIT 1", channelID, nodeID)
	if err != nil {
		return nil, err
	}

	return &invites, errors.Wrap(err, "GetInvites failed")
}

// Groups invites together
func (store *Store) GroupInvitesByToken(tokens []string, groupID globalid.ID) error {
	if len(tokens) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`UPDATE invites SET group_id=? WHERE token IN (?)`, groupID, tokens)
	if err != nil {
		return errors.Wrap(err, "GroupInvitesByToken failed")
	}
	query = store.Rebind(query)
	_, err = store.Exec(query, args...)
	return errors.Wrap(err, "GroupInvitesByToken failed")
}

func (store *Store) ReassignInvitesByToken(tokens []string, userID globalid.ID) error {
	if len(tokens) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`UPDATE invites SET node_id=? WHERE token IN (?)`, userID, tokens)
	if err != nil {
		return errors.Wrap(err, "ReassignInvitesByToken failed")
	}
	query = store.Rebind(query)
	_, err = store.Exec(query, args...)
	return errors.Wrap(err, "ReassignInvitesByToken failed")
}

func (store *Store) ReassignInviterForGroup(rootInviteID, newUserID globalid.ID) error {
	_, err := store.Exec(`UPDATE invites SET node_id = $2, gives_invites = false, group_id = NULL WHERE group_id = (SELECT group_id FROM invites WHERE id = $1) AND id != $1`, rootInviteID, newUserID)
	return errors.Wrap(err, "ReassignInviterForGroup failed")
}

// GetInvites returns all invites from the store
func (store *Store) GetInvitesForUser(id globalid.ID) ([]*model.Invite, error) {
	invites := []*model.Invite{}
	err := store.Select(&invites, "SELECT * FROM invites WHERE node_id = $1 AND remaining_uses > 0 AND hide_from_user = false ORDER BY created_at, token", id)
	if err != nil {
		return nil, errors.Wrap(err, "GetInvitesForUser failed")
	}

	return invites, errors.Wrap(err, "GetInvitesForUser failed")
}

// DeleteInvite deletes the invite with the provided ID
func (store *Store) DeleteInvite(id globalid.ID) error {
	_, err := store.Exec("DELETE FROM invites WHERE id = $1", id)
	return errors.Wrap(err, "DeleteInvite failed")
}

// GetInviteByToken returns the invite for the provided token
func (store *Store) GetInviteByToken(token string) (*model.Invite, error) {
	invite := model.Invite{}
	err := store.Get(&invite, "SELECT * FROM invites where token=$1;", token)
	return &invite, errors.Wrap(err, "GetInviteByToken failed")
}

// GetInviteByToken returns the invite for the provided token
func (store *Store) GetInvite(id globalid.ID) (*model.Invite, error) {
	invite := model.Invite{}
	err := store.Get(&invite, "SELECT * FROM invites where id = $1", id)
	return &invite, errors.Wrap(err, "GetInvite failed")
}

func (store *Store) GetInvitesByID(ids []globalid.ID) ([]*model.Invite, error) {
	query, args, err := sqlx.In(`SELECT invites.* FROM unnest(ARRAY[?]::uuid[]) WITH ORDINALITY AS r(id, rn) LEFT OUTER JOIN invites USING (id) ORDER BY r.rn;`, ids)
	if err != nil {
		return nil, errors.Wrap(err, "GetInvitesByID failed")
	}
	query = store.Rebind(query)
	users := []*model.Invite{}
	err = store.Select(&users, query, args...)
	if len(users) != len(ids) {
		return nil, fmt.Errorf("GetInvitesByID failed: unexpected number of invites returned, likely requesting an invalid id\n")
	}
	return users, errors.Wrap(err, "GetInvitesByID failed")
}
