package model

import (
	"encoding/json"
	"time"

	"github.com/october93/engine/kit/globalid"
)

type Activity struct {
	ID        globalid.ID     `db:"id"          json:"-"`
	RPC       string          `db:"rpc"         json:"rpc"`
	Data      json.RawMessage `db:"data"        json:"data,omitempty"`
	UserID    globalid.ID     `db:"user_id"     json:"userID,omitempty"`
	Error     string          `db:"error"       json:"error,omitempty"`
	CreatedAt time.Time       `db:"created_at"  json:"-"`
}
