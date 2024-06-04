package storage

import "github.com/LilLebowski/loyalty-system/internal/db"

type Storage struct {
}

func Init(databasePath string) *Storage {
	db.Init(databasePath)
	return &Storage{}
}
