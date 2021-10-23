package storage

import (
	"context"
	"errors"

	"github.com/tyghr/social_network/internal/model"
)

var (
	ErrQueueConnIsDead = errors.New("queue conn is dead")
)

type Queue interface {
	Close()

	AddPostBuckets(ctx context.Context, post model.Post, subs []string) error
	ReadPostAppendBuckets(ctx context.Context) (<-chan model.PostBacket, error)

	PostRebuildSubsFeedRequest(ctx context.Context, subs []string) error
	ReadFeedRebuild(ctx context.Context) (<-chan string, error)
}
