package entities

import (
	"context"
	"time"
)

type PhotoService interface {
	CreateTableIfNotExists(ctx context.Context) error
	Add(ctx context.Context, userID int64, url string) (int64, error)
}

type Photo struct {
	ID        int64
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
