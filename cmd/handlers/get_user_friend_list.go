package handlers

import (
	"context"
	"errors"
	"unlinked/proto"
)

func (h *handlers) GetUserFriendsList(
	ctx context.Context,
	req *proto.GetUserFriendsListRequest,
) (*proto.GetUserFriendsListResponse, error) {
	return nil, errors.New("not implemented")
}
