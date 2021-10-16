package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/social_network/internal/model"
)

func (c *Cache) AddPostToSubscribers(ctx context.Context, post model.Post, subs []string) error {
	c.logger.Debugw("AddPostToSubscribers invoked",
		"user", post.UserName,
		"post", post.Header,
		"cut_subs", subs)

	for _, u := range subs {
		err := c.appendUserFeed(ctx, post, u)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) appendUserFeed(ctx context.Context, post model.Post, userName string) error {
	// transaction func

	c.logger.Debugw("updateUserFeed invoked",
		"post_user", post.UserName,
		"post_name", post.Header,
		"subscriber", userName)

	b, err := json.Marshal(post)
	if err != nil {
		return err
	}

	txf := func(tx *redis.Tx) error {
		// Operation is committed only if the watched keys remain unchanged.
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// add to stream
			err := pipe.XAdd(ctx,
				&redis.XAddArgs{
					Stream: feedName(userName),
					MaxLen: maxFeedLen,
					ID:     strconv.FormatInt(post.Created.UnixNano(), 10),
					Values: []interface{}{
						"post_data", b,
					},
				},
			).Err()
			if err != nil {
				c.logger.Errorw("XAdd", "error", err)
			}
			return err
		})
		return err
	}

	for i := 0; i < txMaxRetries; i++ {
		err := c.rc.Watch(
			ctx,
			txf,
			feedName(userName),
		)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}

	return errors.New("tx max retries reached")
}
