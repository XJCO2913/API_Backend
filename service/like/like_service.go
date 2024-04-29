package like

import (
	"context"
	"errors"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LikeService struct{}

var (
	likeService LikeService
)

func Service() *LikeService {
	return &likeService
}

func (s *LikeService) Create(ctx context.Context, input *sdto.CreateLikeInput) *errorx.ServiceErr {
	_, err := dao.GetMomentByID(ctx, input.MomentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("Moment not found by moment ID", zap.String("momentID", input.MomentID))
			return errorx.NewServicerErr(errorx.ErrExternal, "Moment not found by moment ID", nil)
		} else {
			zlog.Error("Failed to retrieve moment by moment ID", zap.String("momentID", input.MomentID), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	_, err = dao.GetLikeByIDs(ctx, input.UserID, input.MomentID)
	if err != gorm.ErrRecordNotFound {
		zlog.Error("User already liked this moment", zap.String("userID", input.UserID), zap.String("momentID", input.MomentID))
		return errorx.NewServicerErr(
			errorx.ErrExternal,
			"Like already exists",
			nil,
		)
	}

	err = dao.CreateNewLike(ctx, &model.Like{
		UserID:   input.UserID,
		MomentID: input.MomentID,
	})
	if err != nil {
		zlog.Error("Error while create new like", zap.String("userID", input.UserID), zap.String("momentID", input.MomentID), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (s *LikeService) DeleteByIDs(ctx context.Context, userID, momentID string) *errorx.ServiceErr {
	_, err := dao.GetLikeByIDs(ctx, userID, momentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("Like not found by user ID and moment ID", zap.String("userID", userID), zap.String("momentID", momentID))
			return errorx.NewServicerErr(errorx.ErrExternal, "Like not found by user ID and moment ID", nil)
		} else {
			zlog.Error("Failed to retrieve like by user ID and moment ID", zap.String("userID", userID), zap.String("momentID", momentID), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	err = dao.DeleteLikeByIDs(ctx, userID, momentID)
	if err != nil {
		zlog.Error("Error while delete the like", zap.String("userID", userID), zap.String("momentID", momentID), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}
