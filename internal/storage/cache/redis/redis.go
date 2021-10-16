package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/storage"
)

var (
	txMaxRetries = 1000
	// postTimeFormat       = "02/01/2006 15:04"
	maxFeedLen int64 = 1000
)

type Cache struct {
	rc     *redis.ClusterClient
	logger logger.Logger
}

func New(nodes []string, password string, l logger.Logger) storage.Cache {
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    nodes,
		Password: password,
	})

	l.Debugw("redis create conn", "nodes", nodes)

	_, err := c.Ping(context.TODO()).Result()
	if err != nil {
		l.Fatalw("Ping redis", "error", err.Error())
	}

	return &Cache{
		rc:     c,
		logger: l,
	}
}

func (c *Cache) Close() {

}

func feedName(name string) string {
	return fmt.Sprintf("userfeed:%s", name)
}

func resetFeedName(name string) string {
	return fmt.Sprintf("resetfeed:%s", name)
}
