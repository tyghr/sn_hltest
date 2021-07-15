package storage

import (
	"context"
	"errors"

	"github.com/tyghr/social_network/internal/model"
)

var (
	ErrDBConnIsDead = errors.New("db conn is dead")
)

type DataBase interface {
	Close()

	CheckAuth(ctx context.Context, username string, phash []byte) (bool, error)
	Register(ctx context.Context, user model.User) error

	GetPosts(ctx context.Context, filter model.PostFilter) ([]model.Post, error)
	GetProfile(ctx context.Context, username string) (model.User, error)

	AddFriend(ctx context.Context, user1, user2 string) error

	SearchUser(ctx context.Context, filter model.UserFilter) ([]model.User, error)

	EditPost(ctx context.Context, post model.Post) error
	DeletePost(ctx context.Context, post model.Post) error
}
