package entities

import (
	"context"
	"time"
)

type UserService interface {
	CreateTableIfNotExists(ctx context.Context) error
	Add(ctx context.Context, name, avatarURL string) (int64, error)
	GetByID(ctx context.Context, userID int64) (*User, error)
	GetFriendsCountByID(ctx context.Context, userID int64) (int64, error)
	GetFriendsByID(ctx context.Context, userID, limit, offset int64) ([]*User, error)
	GetUsersByReactionIDUnderPhoto(ctx context.Context, photoID, reactionID, limit, offset int64) ([]*User, error)
}

type User struct {
	ID        int64
	Name      string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
