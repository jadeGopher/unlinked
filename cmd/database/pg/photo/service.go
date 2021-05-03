package photo

import (
	"database/sql"
	"unlinked/entities"
)

type photoService struct {
	db *sql.DB
}

func NewService(db *sql.DB) entities.PhotoService {
	return &photoService{db: db}
}
