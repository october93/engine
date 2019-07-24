package dataloader

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type IDTimeRangeKey struct {
	id   globalid.ID
	from time.Time
	to   time.Time
}

func NewIDTimeRangeKey(id globalid.ID, from, to time.Time) IDTimeRangeKey {
	return IDTimeRangeKey{
		id:   id,
		from: from,
		to:   to,
	}
}
