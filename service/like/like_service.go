package like

import (
	"context"

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
	like, err := dao.GetLikeByIDs(ctx, input.UserID, input.MomentID)
	if err != gorm.ErrRecordNotFound || like != nil {
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
