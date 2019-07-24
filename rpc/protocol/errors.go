package protocol

import "fmt"

// Error is a wrapper error in order to hide error information to the client.
type Error struct {
	msg string
}

// Error implements the error interface for Error.
func (e *Error) Error() string {
	return e.msg
}

// ErrRPCNotFound is used when a request with an unknown or non-registered RPC
// was made.
func ErrRPCNotFound(rpcName string) error {
	return &Error{msg: fmt.Sprintf("RPC %v not found", rpcName)}
}

// MarshalJSON ensures errors are encoded in the form of strings.
func (e Error) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, e.msg)), nil
}

// UnmarshalJSON ensures errors are read as strings and unmarshaled in the
// error again.
func (e *Error) UnmarshalJSON(b []byte) error {
	e.msg = string(b)
	return nil
}
