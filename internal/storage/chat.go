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
