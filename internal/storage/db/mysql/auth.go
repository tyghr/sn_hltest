package mysql

import (
	"context"
	"fmt"

	"github.com/tyghr/social_network/internal/model"
)

func (db *DB) CheckAuth(ctx context.Context, username string, phash []byte) (bool, error) {
	if err := db.PingContext(ctx); err != nil {
		return false, fmt.Errorf("unable to connect to database: %v", err)
	}
	var un string
	sqlQuery := `SELECT username FROM users WHERE username=? AND phash=?`
	err := db.QueryRowContext(ctx, sqlQuery,
		username,
		phash,
	).Scan(&un)

	db.logger.Debugf("CheckAuth: %s %d", string(un), len(un))

	return un != "", err
}

func (db *DB) Register(ctx context.Context, user model.User) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := user.Validate(); err != nil {
		return err
	}
	_, err := db.NamedExecContext(ctx,
		`INSERT INTO users (username, phash, name, surname, birthdate, gender, city)
		VALUES (:username, :phash, :name, :surname, :birthdate, :gender, :city)`,
		// ON DUPLICATE KEY UPDATE
		// 	username = VALUES(username),
		// 	phash = VALUES(phash),
		// 	name = VALUES(name),
		// 	surname = VALUES(surname),
		// 	birthdate = VALUES(birthdate),
		// 	gender = VALUES(gender),
		// 	city = VALUES(city)`,
		user,
	)
	if err != nil {
		return err
	}

	for _, i := range user.Interests {
		_, err = db.NamedExecContext(ctx,
			`WITH
			di AS ( INSERT INTO d_interests(name) VALUES (:interest) ON CONFLICT(name) DO UPDATE SET name=EXCLUDED.name RETURNING id),
			INSERT INTO interests (user, interest)
			VALUES ((SELECT id FROM users WHERE users.username=:username), (SELECT id FROM di));`,
			// ON DUPLICATE KEY UPDATE
			// 	user = VALUES(user),
			// 	interest = VALUES(interest);`,
			map[string]string{
				"username": user.UserName,
				"interest": i,
			},
		)
		if err != nil {
			return err
		}
	}

	// TODO index for interests
	// TODO transactions

	return nil
}
