package mysql

import (
	"context"
	"fmt"
)

func (db *DB) Subscribe(ctx context.Context, user, subscribeTo string) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	_, err := db.NamedExecContext(ctx,
		`INSERT INTO subscribers (user, subscriber)
		VALUES ((SELECT id FROM users WHERE users.username=:user), (SELECT id FROM users WHERE users.username=:subscriber))
		ON DUPLICATE KEY UPDATE subscriber=subscriber;`,
		map[string]interface{}{
			"user":       subscribeTo,
			"subscriber": user,
		},
	)
	if err != nil {
		return fmt.Errorf("update subscribers: %v", err)
	}

	return db.SetFeedRebuildFlag(ctx, []string{user})
}

func (db *DB) Unsubscribe(ctx context.Context, user, subscribeFrom string) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	_, err := db.NamedExecContext(ctx,
		`DELETE FROM subscribers (user, subscriber)
		VALUES ((SELECT id FROM users WHERE users.username=:user), (SELECT id FROM users WHERE users.username=:subscriber))`,
		map[string]interface{}{
			"user":       subscribeFrom,
			"subscriber": user,
		},
	)
	if err != nil {
		return fmt.Errorf("update subscribers: %v", err)
	}

	return db.SetFeedRebuildFlag(ctx, []string{user})
}

func (db *DB) GetSubscribers(ctx context.Context, user string) ([]string, error) {
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	subscribers := []string{}
	sqlQuery := `SELECT DISTINCT(u.username) FROM subscribers s
		LEFT JOIN users u ON s.subscriber=u.id
		WHERE s.user=(SELECT id FROM users WHERE username=?)
		`
	rows, err := db.QueryxContext(
		ctx,
		sqlQuery,
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, s)
	}

	return subscribers, nil
}

func (db *DB) GetSubscriptions(ctx context.Context, user string) ([]string, error) {
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	subscriptions := []string{}
	sqlQuery := `SELECT DISTINCT(u.username) FROM subscribers s
		LEFT JOIN users u ON s.user=u.id
		WHERE s.subscriber=(SELECT id FROM users WHERE username=?)`
	rows, err := db.QueryxContext(
		ctx,
		sqlQuery,
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}
