package entities

import (
	"context"
	"time"
)

type PhotoService interface {
	CreateTableIfNotExists(ctx context.Context) error
	Add(ctx context.Context, userID int64, url string) (int64, error)
	GetPhotosCountByUserID(ctx context.Context, userID int64) (int64, error)
	GetPhotosByUserID(ctx context.Context, userID, limit, offset int64) ([]*Photo, error)
}

type Photo struct {
	ID        int64
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
