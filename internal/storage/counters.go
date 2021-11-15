package storage

import (
	"context"

	"github.com/tyghr/social_network/internal/model"
)

type CountersStorage struct {
	queue CountersQueue
	cache CountersCache
}

func NewCounters(queue CountersQueue, cache CountersCache) *CountersStorage {
	return &CountersStorage{
		queue: queue,
		cache: cache,
	}
}

func (s *CountersStorage) Q() CountersQueue {
	return s.queue
}

func (s *CountersStorage) C() CountersCache {
	return s.cache
}

type CountersQueue interface {
	ReadCountersBuckets(ctx context.Context) (<-chan model.CounterCmd, error)
}

type CountersCache interface {
	GetUnreadCount(ctx context.Context, user string) (int64, error)
	IncCounters(ctx context.Context, command string, subs []string) error
	UpdateCursorCounter(ctx context.Context, subs []string) error
}
