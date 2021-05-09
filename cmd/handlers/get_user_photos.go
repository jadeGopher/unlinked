package handlers

import (
	"context"
	"go.uber.org/zap"
	"unlinked/cmd/errs"
	"unlinked/entities"
	"unlinked/proto"
	"unlinked/tools/pagination"
)

func (h *handlers) GetUserPhotos(
	ctx context.Context,
	req *proto.GetUserPhotosRequest,
) (_ *proto.GetUserPhotosResponse, err error) {
	if err = req.Validate(); err != nil {
		h.logger.Error("Error Validate")
		return nil, err
	}

	var photosCount int64
	if photosCount, err = h.photoService.GetPhotosCountByUserID(ctx, req.ProfileId); err != nil {
		h.logger.Error("Error photoService.GetPhotosCountByUserID", zap.Error(err))
		return nil, errs.InternalServerError
	}

	pagination.GetPageCount(photosCount, req.PageSize)

	var photos []*entities.Photo
	if photos, err = h.photoService.GetPhotosByUserID(
		ctx,
		req.ProfileId,
		req.PageSize,
		pagination.GetPageCount(photosCount, req.PageSize)*req.PageNumber,
	); err != nil {
		h.logger.Error("Error photoService.GetPhotosByUserID", zap.Error(err))
		return nil, errs.InternalServerError
	}

	var resp = &proto.GetUserPhotosResponse{TotalPhotosCount: photosCount}

	for _, photo := range photos {
		var reactionsInfo []*entities.ReactionInfo
		if reactionsInfo, err = h.reactionService.GetAllReactionsCountByPhotoID(ctx, photo.ID); err != nil {
			h.logger.Error("Error reactionService.GetAllReactionsCountByPhotoID", zap.Error(err))
			return nil, errs.InternalServerError
		}
		var protoReactionsInfo = make([]*proto.ReactionInfo, 0, len(reactionsInfo))
		for _, r := range reactionsInfo {
			protoReactionInfoElem := &proto.ReactionInfo{
				Id:    r.ID,
				Name:  r.Name,
				Count: r.Count,
			}
			protoReactionsInfo = append(protoReactionsInfo, protoReactionInfoElem)
		}

		tmp := &proto.Photo{
			Id:            photo.ID,
			Url:           photo.Url,
			ReactionsInfo: protoReactionsInfo,
			CreatedAt:     photo.CreatedAt.Unix(),
		}
		resp.Photos = append(resp.Photos, tmp)
	}

	return resp, nil
}
