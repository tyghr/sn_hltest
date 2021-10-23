package storage

import (
	"context"
	"errors"

	"github.com/tyghr/social_network/internal/model"
)

var (
	ErrCacheConnIsDead = errors.New("cache conn is dead")
)

type CachePost interface {
	ID() string
	Value(string) interface{}
}

type Cache interface {
	Close()

	AddPostToSubscribers(ctx context.Context, post model.Post, subs []string) error

	GetSubscriptionPosts(ctx context.Context, userToken, id string) ([]CachePost, error)

	RebuildFeed(ctx context.Context, sub string, posts []model.Post) error
}
