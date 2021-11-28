package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
	"github.com/tyghr/social_network/internal/storage"
)

type DB struct {
	*sqlx.DB
	logger logger.Logger
}

func connect(dbUrl string) (db *sqlx.DB, err error) {
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("mysql", dbUrl)
		if err == nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
	return
}

func OpenConn(conf *config.DBConfig, l logger.Logger) (storage.DataBase, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", conf.User, conf.Pass, conf.Host, conf.Port, conf.Name)
	db, err := connect(dbUrl)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = runMigrations(conf, l)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:     db,
		logger: l,
	}, nil
}

func (db *DB) Close() {
	db.DB.Close()
}
