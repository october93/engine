package context

const (
	// RequestID is a unique identifier for every request in order to
	// implement distributed request tracing.
	RequestID = "RequstID"
	// SessionID is used in order to retrieve the current session for the given
	// context.
	SessionID = "SessionID"
	// Callback is used to identify and match up the RPC of the original
	// request. It is helpful for logging and the client.
	Callback = "Callback"
	// IPAddress is the context key used to identify the IP address of a
	// connection.
	IPAddress = "IPAddress"
	// IPAddress is the context key used to identify the user agent of a
	// connection.
	UserAgent = "UserAgent"
)
