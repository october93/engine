package protocol

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/october93/engine/kit/globalid"
)

// DefaultEncoder encodes the result and writes it to the given push writer.
//
// DefaultEncoder provides a default way of handling the response writing. The
// Request ID is used as a value in the ack field of the message. The original
// RPC, if this message is a response to a request, is used as a value in the
// callback field.
func DefaultEncoder(ctx context.Context, data interface{}, err error, pw *PushWriter) error {
	enc := json.NewEncoder(pw)
	requestID, ok := ctx.Value(RequestID).(globalid.ID)
	if !ok {
		return fmt.Errorf("invalid protocol.RequestID in context: %v", ctx.Value(RequestID))
	}
	callback, ok := ctx.Value(Callback).(string)
	if !ok {
		return fmt.Errorf("invalid protocol.Callback in context: %v", ctx.Value(Callback))
	}

	if err != nil {
		pw.log.Info("sending error response", "rpc", callback, "ack", requestID, "error", err.Error())
		return enc.Encode(Message{RPC: callback, Ack: requestID, Err: err.Error()})
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return enc.Encode(Message{RPC: callback, Ack: requestID, Err: err.Error()})
	}
	pw.log.Info("sending response", "rpc", callback, "ack", requestID, "response", json.RawMessage(raw))
	return enc.Encode(Message{RPC: callback, Ack: requestID, Data: raw})
}
