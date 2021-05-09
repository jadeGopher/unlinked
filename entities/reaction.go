package entities

import (
	"context"
	"time"
)

type ReactionService interface {
	CreateTableIfNotExists(ctx context.Context) error
	Add(ctx context.Context, reactionName string) (int64, error)
	AddReactionToPhoto(ctx context.Context, photoID, reactionID, userID int64) error
	GetAllReactionsCountByPhotoID(ctx context.Context, photoID int64) ([]*ReactionInfo, error)
	GetReactionsCountByPhotoID(ctx context.Context, photoID, reactionID int64) (int64, error)
}

type Reaction struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReactionInfo struct {
	ID    int64
	Name  string
	Count int64
}
