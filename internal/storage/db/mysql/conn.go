package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
)

type DB struct {
	*sqlx.DB
	logger logger.Logger
}

func OpenConn(conf *config.Config, l logger.Logger) (*DB, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", conf.DBuser, conf.DBpass, conf.DBhost, conf.DBport, conf.DBname))
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &DB{
		DB:     db,
		logger: l,
	}, nil
}

func (db *DB) Close() {
	db.DB.Close()
}
