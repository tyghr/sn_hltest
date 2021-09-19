package chat

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/social_network/internal/storage"
)

type Chat struct {
	rc *redis.ClusterClient
}

func Init(nodes []string, password string) storage.Chat {
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    nodes,
		Password: password,
	})
	c.Ping(context.TODO())

	return &Chat{rc: c}
}

func getStreamName(chatID, userName string, userMsgCount int) string {
	return fmt.Sprintf("%s.{%s.%d}", chatID, userName, userMsgCount/1000)
}

func getChannelName(chatID string) string {
	return fmt.Sprintf("channel:%s", chatID)
}
