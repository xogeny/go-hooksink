package hooksink

type PushHandler interface {
	Push(HubMessage)
}
