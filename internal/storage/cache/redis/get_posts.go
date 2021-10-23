package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/social_network/internal/storage"
)

// []redis.XMessage
func (c *Cache) GetSubscriptionPosts(ctx context.Context, user, id string) ([]storage.CachePost, error) {
	_, err := c.rc.Ping(ctx).Result()
	if err != nil {
		c.logger.Errorw("ping redis", "error", err.Error())
	}

	msgs := []storage.CachePost{}

	c.logger.Debugw("reading stream",
		"user", user)

	res, err := c.rc.XRead(ctx, &redis.XReadArgs{
		Streams: []string{
			feedName(user), id,
			resetFeedName(user), strconv.FormatInt(time.Now().UnixNano(), 10),
		},
		Block: -1,
	}).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		c.logger.Debugw("xread nil")
		return msgs, nil
	}

	for _, stream := range res {
		for _, sm := range stream.Messages {
			msgs = append(msgs, &ChatMessage{sm})
		}
	}

	// sort.Sort(storage.ByMsgID(msgs))

	return msgs, nil
}
