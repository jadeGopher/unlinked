package reaction

import (
	"database/sql"
	"unlinked/entities"
)

type reactionService struct {
	db *sql.DB
}

func NewService(db *sql.DB) entities.ReactionService {
	return &reactionService{db: db}
}
