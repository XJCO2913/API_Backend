package activity

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
	"gorm.io/gorm"
)

type ActivityService struct{}

var (
	activityService ActivityService
)

func Service() *ActivityService {
	return &activityService
}

func (a *ActivityService) Create(ctx context.Context, in *sdto.CreateActivityInput) *errorx.ServiceErr {
	activity, err := dao.FindActivityByName(ctx, in.Name)
	if err != gorm.ErrRecordNotFound || activity != nil {
		return errorx.NewServicerErr(
			errorx.ErrExternal,
			"Activity already exists",
			nil,
		)
	}

	coverName, uploadErr := a.UploadCover(ctx, in.CoverData)
	if uploadErr != nil {
		return uploadErr
	}

	// Generate a uuid for the new activity
	uuid, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generate uuid: " + err.Error())
		return errorx.NewInternalErr()
	}
	newActivityID := uuid.String()

	err = dao.CreateNewActivity(ctx, &model.Activity{
		ActivityID:  newActivityID,
		Name:        in.Name,
		Description: in.Description,
		RouteID:     in.RouteID,
		CoverURL:    coverName,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		Tags:        in.Tags,
		NumberLimit: in.NumberLimit,
		Fee:         in.Fee,
	})
	if err != nil {
		zlog.Error("Error while create new activity: "+err.Error(), zap.String("name", in.Name))
		return errorx.NewInternalErr()
	}

	return nil
}

func (a *ActivityService) UploadCover(ctx context.Context, coverData []byte) (string, *errorx.ServiceErr) {
	// Check by GetByID TBD...

	coverName, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generate uuid: " + err.Error())
		return "", errorx.NewInternalErr()
	}
	coverNameStr := coverName.String()

	// Upload the cover to minio
	err = minio.UploadActivityCover(ctx, coverNameStr, coverData)
	if err != nil {
		zlog.Error("Error while store activity cover into minio", zap.Error(err))
		return "", errorx.NewInternalErr()
	}

	return coverNameStr, nil
}
