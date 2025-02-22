package database

import (
	"database/sql"
	"fmt"
	"user/config"
	"user/consts"

	// Importing pgx for side effects, such as driver registration
	_ "github.com/jackc/pgx/v4" // PostgreSQL driver for pgx

	// Importing pq for side effects, such as driver registration
	_ "github.com/lib/pq" // PostgreSQL driver for lib/pq
)

// ConnectPsqlDB initializes a connection to a PostgreSQL database using the provided configuration.
func Connect(cfg config.Database) (*sql.DB, error) {
	datasource := ConnectionString(cfg)
	databaseType := consts.DatabaseType
	db, err := sql.Open(databaseType, datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %s err: %w", datasource, err)
	}
	db.SetMaxOpenConns(cfg.MaxActive)
	db.SetMaxIdleConns(cfg.MaxIdle)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db(ping): %s err: %w", datasource, err)
	}

	return db, nil
}

// ConnectionString constructs a PostgreSQL connection string using the provided configuration.
func ConnectionString(cfg config.Database) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=60 search_path=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DATABASE, cfg.Schema)
}
