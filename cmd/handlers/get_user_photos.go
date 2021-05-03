package handlers

import (
	"context"
	"errors"
	"unlinked/proto"
)

func (h *handlers) GetUserPhotos(
	ctx context.Context,
	req *proto.GetUserPhotosRequest,
) (*proto.GetUserPhotosResponse, error) {
	return nil, errors.New("not implemented")
}
