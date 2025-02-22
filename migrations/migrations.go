package migrations

import (
	"database/sql"
	"errors"
	"os"
	"user/consts"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Required for migration file source
	_ "github.com/lib/pq"                                // PostgreSQL driver for database connections
)

func Up(db *sql.DB) error {
	driverDB, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	// migration instance creation
	migration, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		consts.DatabaseType, driverDB)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func Down(db *sql.DB) error {
	driverDB, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	// migration instance creation
	migration, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		consts.DatabaseType, driverDB)
	if err != nil {
		return err
	}

	err = migration.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
