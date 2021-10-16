package storage

import "context"

type ChatMessage interface {
	ID() string
	Value(string) interface{}
}

type ChannelMessage interface {
	Message() string
}

type Chat interface {
	WriteMessage(context.Context, string, string, string) error
	GetMsgs(context.Context, string) (<-chan ChannelMessage, []ChatMessage, func() error)
}

// sort.Sort(ByMsgID([]ChatMessage{}))

type ByMsgID []ChatMessage

func (msg ByMsgID) Len() int {
	return len(msg)
}

func (msg ByMsgID) Swap(i, j int) {
	msg[i], msg[j] = msg[j], msg[i]
}

func (msg ByMsgID) Less(i, j int) bool {
	return msg[i].ID() < msg[j].ID()
}
