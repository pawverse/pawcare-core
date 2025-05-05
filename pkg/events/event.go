package events

type IEvent interface {
	Key() string
	EventName() string
}
