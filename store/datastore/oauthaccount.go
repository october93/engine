package datastore

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

func (s *Store) SaveOAuthAccount(m *model.OAuthAccount) error {
	return saveOAuthAccount(s, m)
}

func (tx *Tx) SaveOAuthAccount(m *model.OAuthAccount) error {
	return saveOAuthAccount(tx, m)
}

func saveOAuthAccount(e sqlx.Ext, m *model.OAuthAccount) error {
	if m == nil {
		return errors.New("provier model cannot be be nil")
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
		`INSERT INTO oauth_accounts
		(
			id,
			provider,
			subject,
			user_id,
			created_at,
			updated_at
		)
		VALUES
		(
			:id,
			:provider,
			:subject,
			:user_id,
			:created_at,
			:updated_at
		)
		ON CONFLICT(id) DO UPDATE
		SET
			provider   = :provider,
			subject    = :subject,
			user_id    = :user_id,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE oauth_accounts.ID = :id`, m)
	return errors.Wrap(err, "SaveOAuthAccount failed")
}

// GetOAuthAccountBySubject returns the OAuth account from the database with
// the given subject. This is limited to Facebook OAuth accounts for now.
func (s *Store) GetOAuthAccountBySubject(subject string) (*model.OAuthAccount, error) {
	oauthAccount := model.OAuthAccount{}
	err := s.Get(&oauthAccount, "SELECT * FROM oauth_accounts WHERE provider = $1 AND subject = $2", model.FacebookProvider, subject)
	if err != nil {
		return nil, errors.Wrap(err, "GetOAuthAccountBySubject failed")
	}
	return &oauthAccount, nil
}
