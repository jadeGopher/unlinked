package handlers

import (
	"context"
	"go.uber.org/zap"
	"unlinked/cmd/errs"
	"unlinked/entities"
	"unlinked/proto"
	"unlinked/tools/pagination"
)

func (h *handlers) GetUserFriendsList(
	ctx context.Context,
	req *proto.GetUserFriendsListRequest,
) (_ *proto.GetUserFriendsListResponse, err error) {
	if err = req.Validate(); err != nil {
		h.logger.Error("Error Validate")
		return nil, err
	}

	var friendsCount int64
	if friendsCount, err = h.userService.GetFriendsCountByID(ctx, req.ProfileId); err != nil {
		h.logger.Error("Error userService.GetFriendsCountByID", zap.Error(err))
		return nil, errs.InternalServerError
	}

	pagination.GetPageCount(friendsCount, req.PageSize)

	var users []*entities.User
	if users, err = h.userService.GetFriendsByID(
		ctx,
		req.ProfileId,
		req.PageSize,
		pagination.GetPageCount(friendsCount, req.PageSize)*req.PageNumber,
	); err != nil {
		h.logger.Error("Error userService.GetFriendsByID", zap.Error(err))
		return nil, errs.InternalServerError
	}

	var resp = &proto.GetUserFriendsListResponse{TotalFriendsCount: friendsCount}
	resp.Friends = make([]*proto.UserProfile, 0, req.PageSize)
	for _, user := range users {
		if friendsCount, err = h.userService.GetFriendsCountByID(ctx, user.ID); err != nil {
			h.logger.Error("Error userService.GetFriendsCountByID", zap.Error(err))
			return nil, errs.InternalServerError
		}

		var photosCount int64
		if photosCount, err = h.photoService.GetPhotosCountByUserID(ctx, user.ID); err != nil {
			h.logger.Error("Error photoService.GetPhotosCountByUserID", zap.Error(err))
			return nil, errs.InternalServerError
		}

		resp.Friends = append(resp.Friends, &proto.UserProfile{
			Id:             user.ID,
			Name:           user.Name,
			Avatar:         user.Avatar,
			RegisteredAt:   user.CreatedAt.Unix(),
			FollowersCount: friendsCount,
			PhotosCount:    photosCount,
		})
	}

	return resp, nil
}
