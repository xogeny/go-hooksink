package hooksink

/* This is the handler for a 'push' event.  Not sure if I will ever
   support other event types (or whether that even makes sense), but
   I thought I'd make this explicit. */
type PushHandler interface {
	Push(msg HubMessage, params map[string][]string)
}
