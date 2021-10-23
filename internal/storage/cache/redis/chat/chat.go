package chat

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/storage"
)

type redisConn interface {
	redis.Cmdable
	Watch(context.Context, func(tx *redis.Tx) error, ...string) error
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

type Chat struct {
	rc     redisConn
	logger logger.Logger
}

func New(nodes []string, isCluster bool, password string, l logger.Logger) storage.Chat {
	var c redisConn
	if isCluster {
		c = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    nodes,
			Password: password,
		})
	} else {
		if len(nodes) == 0 {
			l.Fatalw("redis nodelist empty")
		}
		c = redis.NewClient(&redis.Options{
			Addr:     nodes[0],
			Password: password,
		})
	}

	l.Debugw("redis create conn", "nodes", nodes)

	_, err := c.Ping(context.TODO()).Result()
	if err != nil {
		l.Fatalw("Ping redis", "error", err.Error())
	}

	return &Chat{
		rc:     c,
		logger: l,
	}
}

func getStreamName(chatID, userName string, userMsgCount int) string {
	return fmt.Sprintf("%s.{%s.%d}", chatID, userName, userMsgCount/1000)
}

func getChannelName(chatID string) string {
	return fmt.Sprintf("channel:%s", chatID)
}
