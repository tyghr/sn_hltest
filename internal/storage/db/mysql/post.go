package mysql

import (
	"context"
	"fmt"

	"github.com/tyghr/social_network/internal/model"
)

func (db *DB) GetPosts(ctx context.Context, filter model.PostFilter) ([]model.Post, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	posts := []model.Post{}

	sqlQuery := `SELECT header,created,updated,text
		FROM posts
		WHERE user=(SELECT id FROM users WHERE username=?)
		AND deleted='0'
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

func (db *DB) UpsertPost(ctx context.Context, post model.Post) (bool, error) {
	if err := db.PingContext(ctx); err != nil {
		return false, fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := post.Validate(); err != nil {
		return false, err
	}

	res, err := db.NamedExecContext(ctx,
		`INSERT INTO posts (user, header, text, created, updated)
			VALUES ((SELECT id FROM users WHERE username=:username), :header, :text, :created, :updated)
			ON DUPLICATE KEY UPDATE
				text = VALUES(text),
				updated = VALUES(updated);`,
		post,
	)
	if err != nil {
		return false, fmt.Errorf("upsert post: %v", err)
	}

	// 1 if the row is inserted as a new row
	// 2 if an existing row is updated
	// 0 if an existing row is set to its current values
	// If you specify the CLIENT_FOUND_ROWS flag to the mysql_real_connect() C API function when connecting to mysqld, the affected-rows value is 1 (not 0) if an existing row is set to its current values.
	r, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("get rows affected: %v", err)
	}

	return r == 2, err
}

func (db *DB) DeletePost(ctx context.Context, post model.Post) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := post.Validate(); err != nil {
		return err
	}

	_, err := db.NamedExecContext(ctx,
		`UPDATE posts
		SET deleted='1'
			WHERE user=(SELECT id FROM users WHERE username=:username) AND header=:header;`,
		post,
	)

	return err
}
