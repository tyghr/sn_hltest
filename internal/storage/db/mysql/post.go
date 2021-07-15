package mysql

import (
	"context"

	"github.com/tyghr/social_network/internal/model"
)

func (db *DB) GetPosts(ctx context.Context, filter model.PostFilter) ([]model.Post, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	posts := []model.Post{}

	sqlQuery := `SELECT header,updated,text
		FROM posts WHERE user=(SELECT id FROM users WHERE username=?)
		ORDER by updated DESC`

	rowsList, err := db.QueryxContext(
		ctx,
		sqlQuery,
		filter.UserName,
	)
	if err != nil {
		return nil, err
	}
	defer rowsList.Close()

	for rowsList.Next() {
		p := model.Post{}
		if err := rowsList.StructScan(&p); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (db *DB) EditPost(ctx context.Context, post model.Post) error {
	return nil
}
func (db *DB) DeletePost(ctx context.Context, post model.Post) error {
	return nil
}
