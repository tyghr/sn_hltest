package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/social_network/internal/model"
)

func (c *Cache) RebuildFeed(ctx context.Context, user string, posts []model.Post) error {
	c.logger.Debugw("RebuildFeed invoked",
		"user", user,
	)

	err := c.reFill(ctx, user, posts)
	if err != nil {
		return fmt.Errorf("reFill: %v", err)
	}

	err = c.sendResetMsg(ctx, user)
	if err != nil {
		return fmt.Errorf("reset FE msg: %v", err)
	}

	return nil
}

func (c *Cache) reFill(ctx context.Context, user string, posts []model.Post) error {
	// transaction func

	txf := func(tx *redis.Tx) error {
		// Operation is committed only if the watched keys remain unchanged.
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// delete (XTRIM mystream MAXLEN 0)
			err := pipe.XTrimMaxLen(ctx, feedName(user), 0).Err()
			if err != nil {
				c.logger.Errorw("XTrimMaxLen", "error", err)
				return err
			}

			for _, post := range posts {
				b, err := json.Marshal(post)
				if err != nil {
					return err
				}

				// add to stream (fill_from_db)
				err = pipe.XAdd(ctx,
					&redis.XAddArgs{
						Stream: feedName(user),
						MaxLen: maxFeedLen,
						ID:     strconv.FormatInt(post.Created.UnixNano(), 10),
						Values: []interface{}{
							"post_data", b,
						},
					},
				).Err()
				if err != nil {
					c.logger.Errorw("XAdd", "error", err)
					return err
				}
			}
			return nil
		})
		return err
	}

	for i := 0; i < txMaxRetries; i++ {
		err := c.rc.Watch(
			ctx,
			txf,
			feedName(user),
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

func (c *Cache) sendResetMsg(ctx context.Context, user string) error {
	// reset FE msg
	err := c.rc.XAdd(ctx,
		&redis.XAddArgs{
			Stream: resetFeedName(user),
			MaxLen: 1,
			ID:     strconv.FormatInt(time.Now().UnixNano(), 10),
			Values: []interface{}{
				"reset", 1,
			},
		},
	).Err()
	return err
}
