package push

import (
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/rpc/protocol"
)

type connections interface {
	WritersByUser(globalid.ID) []*protocol.PushWriter
	Writers() []*protocol.PushWriter
}

type store interface {
	GetEngagement(cardID globalid.ID) (*model.Engagement, error)
}
