package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tyghr/social_network/internal/model"
)

func (db *DB) SearchUser(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	users := []model.User{}
	if err := db.PingContext(ctx); err != nil {
		return users, fmt.Errorf("unable to connect to database: %v", err)
	}

	qWhere := []string{}
	if filter.UserName != "" {
		qWhere = append(qWhere, fmt.Sprintf("username LIKE '%s'", filter.UserName))
	}
	if filter.Name != "" {
		qWhere = append(qWhere, fmt.Sprintf("name LIKE '%s'", filter.Name))
	}
	if filter.SurName != "" {
		qWhere = append(qWhere, fmt.Sprintf("surname LIKE '%s'", filter.SurName))
	}
	if filter.Gender == "M" {
		qWhere = append(qWhere, "gender LIKE '1'")
	} else if filter.Gender == "F" {
		qWhere = append(qWhere, "gender LIKE '0'")
	}
	if filter.City != "" {
		qWhere = append(qWhere, fmt.Sprintf("city LIKE '%s'", filter.City))
	}
	if filter.Interests != "" {
		// TODO
	}
	if filter.Friends != "" {
		// TODO
	}
	nilTime := time.Time{}
	if filter.AgeFrom != "" {
		// TODO
	} else if filter.BirthDateFrom != nilTime {
		// TODO
	}
	if filter.AgeTo != "" {
		// TODO
	} else if filter.BirthDateTo != nilTime {
		// TODO
	}

	sqlQuery := `SELECT username,name,surname,birthdate,gender,city FROM users`
	if len(qWhere) > 0 {
		sqlQuery += " WHERE " + strings.Join(qWhere, " AND ")
	}
	sqlQuery += " ORDER BY id ASC LIMIT 100" // TODO

	db.logger.Debugf("sqlQuery: %s", sqlQuery)

	rows, err := db.QueryxContext(ctx, sqlQuery)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		u := model.User{}
		err = rows.StructScan(&u)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
