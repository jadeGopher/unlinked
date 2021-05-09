package handlers

import (
	"context"
	"go.uber.org/zap"
	"unlinked/cmd/errs"
	"unlinked/entities"
	"unlinked/proto"
)

func (h *handlers) GetUserProfileInfo(
	ctx context.Context,
	req *proto.GetUserProfileInfoRequest,
) (_ *proto.GetUserProfileInfoResponse, err error) {
	if err = req.Validate(); err != nil {
		h.logger.Error("Error Validate")
		return nil, err
	}

	var user = &entities.User{}
	if user, err = h.userService.GetByID(ctx, req.ProfileId); err != nil {
		h.logger.Error("Error userService.GetByID", zap.Error(err))
		return nil, errs.InternalServerError
	}

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

	return &proto.GetUserProfileInfoResponse{UserProfile: &proto.UserProfile{
		Id:             user.ID,
		Name:           user.Name,
		Avatar:         user.Avatar,
		RegisteredAt:   user.CreatedAt.Unix(),
		FollowersCount: friendsCount,
		PhotosCount:    photosCount,
	}}, nil
}
