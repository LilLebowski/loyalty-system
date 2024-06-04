package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Instance struct {
	db *sql.DB
}

func Init(databasePath string) (*Instance, error) {
	db, err := sql.Open("pgx", databasePath)

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	dbStore := &Instance{
		db: db,
	}

	return dbStore, nil
}
