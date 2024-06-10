package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
)

func Init(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databasePath)

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	return db, nil
}

func RunMigrations(cfg *config.Config) error {
	if cfg.DBPath == "" {
		return fmt.Errorf("got empty dbURI")
	}
	m, err := migrate.New(cfg.MigrateSourceURL, cfg.DBPath)
	if err != nil {
		log.Printf("Got err %s", err.Error())
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Got err %s", err.Error())
		return err
	}
	return nil
}
