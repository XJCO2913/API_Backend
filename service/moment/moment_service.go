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
	momentId := uuid.New()
	momentIdStr := momentId.String()
	_, err := dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		MomentID: momentIdStr,
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
	momentId := uuid.New()
	momentIdStr := momentId.String()
	_, err = dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		ImageURL: &imageNameStr,
		MomentID: momentIdStr,
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

func (m *MomentService) CreateWithVideo(ctx context.Context, in *sdto.CreateMomentVideoInput) *errorx.ServiceErr {
	videoName := uuid.New()
	videoNameStr := videoName.String()

	momentId := uuid.New()
	momentIdStr := momentId.String()
	newMomentID, err := dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		VideoURL: &videoNameStr,
		MomentID: momentIdStr,
	})
	if err != nil {
		zlog.Error("error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// async upload video
	go func() {
		err := minio.UploadMomentVideo(ctx, videoNameStr, in.VideoData)
		if err != nil {
			// error, remove moment record in mysql
			zlog.Error("error while async upload moment video", zap.Error(err))
			dao.DeleteMomentByID(ctx, newMomentID)
		}
	}()

	return nil
}