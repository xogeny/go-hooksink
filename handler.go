package hooksink

// This is the handler for a 'push' event.  Not sure if I will ever
// support other event types (or whether that even makes sense), but I
// thought it was a good idea to give this handler a very specific
// name to protect for future enhancements, just in case.
type PushHandler interface {
	// The Push method gets the message from GitHub as well as
	// the values of any query parameters that were part of the
	// URL that was POSTED.
	Push(msg PushMessage, params map[string][]string)
}
