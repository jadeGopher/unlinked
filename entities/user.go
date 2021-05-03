package entities

import (
	"context"
	"time"
)

type UserService interface {
	CreateTableIfNotExists(ctx context.Context) error
	Add(ctx context.Context, name, avatarURL string) (int64, error)
	GetByID(ctx context.Context, userID int64) (*User, error)
	GetFriendsByID(ctx context.Context, userID, limit, offset int64) ([]*User, error)
}

type User struct {
	ID        int64
	Name      string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
