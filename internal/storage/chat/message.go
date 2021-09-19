package chat

import "github.com/go-redis/redis/v8"

type ChatMessage struct {
	redis.XMessage
}

func (m *ChatMessage) ID() string {
	return m.XMessage.ID
}

func (m *ChatMessage) Value(key string) interface{} {
	return m.XMessage.Values[key]
}

type ChannelMessage struct {
	cm *redis.Message
}

func (m *ChannelMessage) Message() string {
	return m.cm.Payload
}
