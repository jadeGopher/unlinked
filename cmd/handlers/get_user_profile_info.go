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
	var user = &entities.User{}
	if user, err = h.userService.GetByID(ctx, req.ProfileId); err != nil {
		h.logger.Error("userService.GetByID", zap.Error(err))
		return nil, errs.InternalServerError
	}

	return &proto.GetUserProfileInfoResponse{UserProfile: &proto.UserProfile{
		Id:             user.ID,
		Name:           user.Name,
		Photo:          user.Avatar,
		RegisteredAt:   user.CreatedAt.Unix(),
		FollowersCount: 0,
		PhotosCount:    0,
	}}, nil
}
