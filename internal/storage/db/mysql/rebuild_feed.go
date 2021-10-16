package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/tyghr/social_network/internal/model"
)

func (db *DB) SetFeedRebuildFlag(ctx context.Context, users []string) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	args := make([]interface{}, len(users))
	for i, u := range users {
		args[i] = u
	}
	_, err := db.ExecContext(ctx,
		`UPDATE users SET rebuild_feed_flag='1'
		WHERE username IN (?`+strings.Repeat(",?", len(args)-1)+`)`,
		args...,
	)
	if err != nil {
		return fmt.Errorf("set rebuild_feed_flag: %v", err)
	}
	return nil
}

func (db *DB) GetFeedRebuild(ctx context.Context, user string) (bool, []model.Post, error) {
	if err := db.PingContext(ctx); err != nil {
		return false, nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	var nr bool
	sqlQuery := `SELECT rebuild_feed_flag FROM users
		WHERE username=?`
	err := db.QueryRowxContext(
		ctx,
		sqlQuery,
		user,
	).Scan(&nr)
	if err != nil {
		return false, nil, fmt.Errorf("get rebuild_feed_flag: %v", err)
	}

	if !nr {
		return nr, nil, nil
	}

	_, err = db.ExecContext(ctx,
		`UPDATE users SET rebuild_feed_flag='0'
		WHERE username=?`,
		user,
	)
	if err != nil {
		return false, nil, fmt.Errorf("reset rebuild_feed_flag: %v", err)
	}

	posts, err := db.GetPosts(ctx, model.PostFilter{UserName: user})

	return nr, posts, err
}
