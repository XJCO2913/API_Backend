package moment

import (
	"context"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util/zlog"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MomentService struct{}

var (
	momentService MomentService
)

func Service() *MomentService {
	return &momentService
}

func (m *MomentService) Create(ctx context.Context, in *sdto.CreateMomentInput) *errorx.ServiceErr {
	err := dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content: &in.Content,
	})
	if err != nil {
		zlog.Error("error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (m *MomentService) CreateWithImage(ctx context.Context, in *sdto.CreateMomentImageInput) *errorx.ServiceErr {
	imageName := uuid.New()
	err := minio.UploadMomentImage(ctx, imageName.String(), in.ImageData)
	if err != nil {
		zlog.Error("error while store moment image into minio", zap.Error(err))
		return errorx.NewInternalErr()
	}

	imageNameStr := imageName.String()
	err = dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		ImageURL: &imageNameStr,
	})
	if err != nil {
		zlog.Error("error while create new moment", zap.Error(err))

		// async remove image in minio
		go func() {
			minio.RemoveObjectFromMoment(ctx, imageNameStr)
		}()

		return errorx.NewInternalErr()
	}

	return nil
}
