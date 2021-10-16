package mysql

import (
	"context"
	"fmt"
)

func (db *DB) GetRelations(ctx context.Context, user, f_user string) (bool, bool, error) {
	if err := db.PingContext(ctx); err != nil {
		return false, false, fmt.Errorf("unable to connect to database: %v", err)
	}

	var isF, isS bool
	sqlQuery := `SELECT EXISTS(
		SELECT * FROM friends
		WHERE user=(SELECT id FROM users WHERE username=?)
		AND friend=(SELECT id FROM users WHERE username=?)),
		EXISTS(
			SELECT * FROM subscribers
			WHERE user=(SELECT id FROM users WHERE username=?)
			AND subscriber=(SELECT id FROM users WHERE username=?))`
	err := db.QueryRowxContext(
		ctx,
		sqlQuery,
		user,
		f_user,
		f_user,
		user,
	).Scan(&isF, &isS)
	return isF, isS, err
}
