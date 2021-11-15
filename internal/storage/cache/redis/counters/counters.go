package counters

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/social_network/internal/model"
)

func (c *Cache) GetUnreadCount(ctx context.Context, user string) (int64, error) {
	var unreadCount int64

	// transaction func
	txf := func(tx *redis.Tx) error {
		cursorCnt, err := tx.Get(ctx,
			getCounterName(model.CounterCmdCursorInc, user),
		).Int64()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("get user cursor_count %w", err)
		}

		totalCnt, err := tx.Get(ctx,
			getCounterName(model.CounterCmdTotalInc, user),
		).Int64()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("get user total_count %w", err)
		}

		unreadCount = totalCnt - cursorCnt
		return nil
	}

	for i := 0; i < txMaxRetries; i++ {
		err := c.rc.Watch(
			ctx,
			txf,
			getCounterName(model.CounterCmdCursorInc, user),
		)
		if err == redis.TxFailedErr {
			continue
		} else if err != nil {
			return 0, fmt.Errorf("Watch: %w", err)
		}
		return unreadCount, nil
	}
	return 0, errors.New("tx max retries reached")
}

func (c *Cache) IncCounters(ctx context.Context, command string, subs []string) error {

	for _, s := range subs {
		err := c.rc.Incr(ctx,
			getCounterName(command, s),
		).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) UpdateCursorCounter(ctx context.Context, subs []string) error {
	for _, user := range subs {
		err := c.updateUserCursorCounter(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) updateUserCursorCounter(ctx context.Context, user string) error {
	// transaction func
	txf := func(tx *redis.Tx) error {
		// get user total_count
		cnt, err := tx.Get(ctx,
			getCounterName(model.CounterCmdTotalInc, user),
		).Int()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("get user total_count %w", err)
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// update user cursor_count to the value of total_count
			err := pipe.Set(ctx,
				getCounterName(model.CounterCmdCursorInc, user),
				cnt,
				0,
			).Err()
			if err != nil {
				return fmt.Errorf("update user cursor_count %w", err)
			}
			return nil
		})
		return err
	}

	for i := 0; i < txMaxRetries; i++ {
		err := c.rc.Watch(
			ctx,
			txf,
			getCounterName(model.CounterCmdCursorInc, user),
		)
		if err == redis.TxFailedErr {
			continue
		} else if err != nil {
			return fmt.Errorf("Watch: %w", err)
		}
		return nil
	}
	return errors.New("tx max retries reached")
}
