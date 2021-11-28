package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migrateMysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(conf *config.DBConfig, lgr logger.Logger) error {
	lgr.Debugw("start migrations...")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true", conf.User, conf.Pass, conf.Host, conf.Port, conf.Name)
	instance, err := sql.Open("mysql", dbUrl)
	if err != nil {
		lgr.Errorw("failed migration")

		return fmt.Errorf("Connection error: %v", err)
	}

	driver, err := migrateMysql.WithInstance(instance, &migrateMysql.Config{})
	if err != nil {
		lgr.Errorw("failed migration")

		return fmt.Errorf("Error while configuring the driver: %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", conf.MigrationPath), "mysql", driver)
	if err != nil {
		lgr.Errorw("failed migration")

		return fmt.Errorf("Error while configuring the migrate instance: %v", err)
	}
	defer migration.Close()

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		lgr.Error("Migration UP error: ", err)

		if err := migration.Down(); err != nil {
			lgr.Errorw("failed migration")

			return fmt.Errorf("Migration error: %v", err)
		}

		lgr.Errorw("failed migration")

		return err
	}

	lgr.Debugw("successful migrations")

	return nil
}
