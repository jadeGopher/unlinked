package handlers

import (
	"context"
	"go.uber.org/zap"
	"unlinked/cmd/errs"
	"unlinked/entities"
	"unlinked/proto"
	"unlinked/tools/pagination"
)

func (h *handlers) GetPhotoReactions(
	ctx context.Context,
	req *proto.GetPhotoReactionsRequest,
) (_ *proto.GetPhotoReactionsResponse, err error) {
	if err = req.Validate(); err != nil {
		h.logger.Error("Error Validate")
		return nil, err
	}

	var reactionsCount int64
	if reactionsCount, err = h.reactionService.GetReactionsCountByPhotoID(
		ctx,
		req.PhotoId,
		req.ReactionId,
	); err != nil {
		h.logger.Error("Error reactionService.GetReactionsCountByPhotoID", zap.Error(err))
		return nil, errs.InternalServerError
	}

	pagination.GetPageCount(reactionsCount, req.PageSize)

	var users []*entities.User
	if users, err = h.userService.GetUsersByReactionIDUnderPhoto(
		ctx,
		req.PhotoId,
		req.ReactionId,
		req.PageSize,
		pagination.GetPageCount(reactionsCount, req.PageSize)*req.PageNumber,
	); err != nil {
		h.logger.Error("Error userService.GetUsersByReactionIDUnderPhoto", zap.Error(err))
		return nil, errs.InternalServerError
	}

	var protoReactions = make([]*proto.UserProfile, 0, req.PageSize)
	for _, user := range users {
		var friendsCount int64
		if friendsCount, err = h.userService.GetFriendsCountByID(ctx, user.ID); err != nil {
			h.logger.Error("Error userService.GetFriendsCountByID", zap.Error(err))
			return nil, errs.InternalServerError
		}

		var photosCount int64
		if photosCount, err = h.photoService.GetPhotosCountByUserID(ctx, user.ID); err != nil {
			h.logger.Error("Error photoService.GetPhotosCountByUserID", zap.Error(err))
			return nil, errs.InternalServerError
		}

		tmp := &proto.UserProfile{
			Id:             user.ID,
			Name:           user.Name,
			Avatar:         user.Avatar,
			RegisteredAt:   user.CreatedAt.Unix(),
			FollowersCount: photosCount,
			PhotosCount:    friendsCount,
		}
		protoReactions = append(protoReactions, tmp)
	}

	var resp = &proto.GetPhotoReactionsResponse{
		UserProfile:    protoReactions,
		ReactionsCount: reactionsCount,
	}

	return resp, nil
}
