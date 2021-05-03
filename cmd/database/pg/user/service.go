package user

import (
	"database/sql"
	"unlinked/entities"
)

type userService struct {
	db *sql.DB
}

func NewService(db *sql.DB) entities.UserService {
	return &userService{db: db}
}
