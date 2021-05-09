package handlers

import (
	"database/sql"
	"go.uber.org/zap"
	"unlinked/cmd/database/pg/photo"
	"unlinked/cmd/database/pg/reaction"
	"unlinked/cmd/database/pg/user"
	"unlinked/entities"
	"unlinked/proto"
)

type handlers struct {
	userService     entities.UserService
	photoService    entities.PhotoService
	reactionService entities.ReactionService
	logger          *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) proto.UnlinkedServiceServer {
	return &handlers{
		userService:     user.NewService(db),
		photoService:    photo.NewService(db),
		reactionService: reaction.NewService(db),
		logger:          logger,
	}
}
