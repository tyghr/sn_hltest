package mysql

import (
	"context"
	"fmt"

	"github.com/tyghr/social_network/internal/model"
)

func (db *DB) GetProfile(ctx context.Context, username string) (model.User, error) {
	user := model.User{}
	if err := db.PingContext(ctx); err != nil {
		return user, fmt.Errorf("unable to connect to database: %v", err)
	}

	sqlQuery := `SELECT username,name,surname,birthdate,gender,city FROM users WHERE username=?`
	err := db.QueryRowxContext(ctx, sqlQuery,
		username,
	).StructScan(&user)
	if err != nil {
		return user, err
	}

	// Interests
	interests := []string{}
	sqlQueryI := `SELECT DISTINCT(d.name) FROM interests i
		LEFT JOIN d_interests d ON i.interest=d.id
		WHERE i.user=(SELECT id FROM users WHERE username=?)
		ORDER by d.name`
	rowsI, err := db.QueryxContext(
		ctx,
		sqlQueryI,
		user.UserName,
	)
	if err != nil {
		return user, err
	}
	defer rowsI.Close()

	for rowsI.Next() {
		var i string
		if err := rowsI.Scan(&i); err != nil {
			return user, err
		}
		interests = append(interests, i)
	}
	user.Interests = interests

	// Friends
	friends := []string{}
	sqlQueryF := `SELECT DISTINCT(fu.username) FROM friends f
		LEFT JOIN users fu ON f.friend=fu.id
		WHERE f.user=(SELECT id FROM users WHERE username=?)`
	rowsF, err := db.QueryxContext(
		ctx,
		sqlQueryF,
		user.UserName,
	)
	if err != nil {
		return user, err
	}
	defer rowsF.Close()

	for rowsF.Next() {
		var f string
		if err := rowsF.Scan(&f); err != nil {
			return user, err
		}
		friends = append(friends, f)
	}
	user.Friends = friends

	// Subscriptions
	subscriptions := []string{}
	sqlQueryS := `SELECT DISTINCT(u.username) FROM subscribers s
		LEFT JOIN users u ON s.user=u.id
		WHERE s.subscriber=(SELECT id FROM users WHERE username=?)`
	rowsS, err := db.QueryxContext(
		ctx,
		sqlQueryS,
		user.UserName,
	)
	if err != nil {
		return user, err
	}
	defer rowsS.Close()

	for rowsS.Next() {
		var s string
		if err := rowsS.Scan(&s); err != nil {
			return user, err
		}
		subscriptions = append(subscriptions, s)
	}
	user.Subscriptions = subscriptions

	// Subscribers
	subscribers := []string{}
	sqlQueryS = `SELECT DISTINCT(u.username) FROM subscribers s
		LEFT JOIN users u ON s.subscriber=u.id
		WHERE s.user=(SELECT id FROM users WHERE username=?)`
	rowsS2, err := db.QueryxContext(
		ctx,
		sqlQueryS,
		user.UserName,
	)
	if err != nil {
		return user, err
	}
	defer rowsS2.Close()

	for rowsS2.Next() {
		var s string
		if err := rowsS2.Scan(&s); err != nil {
			return user, err
		}
		subscribers = append(subscribers, s)
	}
	user.Subscribers = subscribers

	return user, nil
}

func (db *DB) AddFriend(ctx context.Context, user, friend string) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	_, err := db.NamedExecContext(ctx,
		`INSERT INTO friends (user, friend)
		VALUES ((SELECT id FROM users WHERE users.username=:user), (SELECT id FROM users WHERE users.username=:friend))
		ON DUPLICATE KEY UPDATE friend=friend;`,
		map[string]interface{}{
			"user":   user,
			"friend": friend,
		},
	)
	// TODO index for interests
	return err
}
