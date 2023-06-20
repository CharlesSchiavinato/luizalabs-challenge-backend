package migration

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(config *util.Config) error {
	migration, err := migrate.New(config.DBMigrationURL, config.DBURL)

	if err != nil {
		return err
	}

	err = migration.Up()

	if err != migrate.ErrNoChange {
		return err
	}

	return nil
}
