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

type redisConn interface {
	redis.Cmdable
	Watch(context.Context, func(tx *redis.Tx) error, ...string) error
}

type Cache struct {
	rc     redisConn
	logger logger.Logger
}

func New(nodes []string, isCluster bool, password string, l logger.Logger) storage.Cache {
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
