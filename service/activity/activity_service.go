package activity

import (
	"context"

	"api.backend.xjco2913/dao"
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
	// Check if activities already exist or not
	activity, err := dao.FindActivityByName(ctx, in.Name)
	if err != gorm.ErrRecordNotFound || activity != nil {
		return errorx.NewServicerErr(
			errorx.ErrExternal,
			"Activity already exist",
			nil,
		)
	}

	// Generate uuid for activityId
	uuid, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generating uuid for activity: " + err.Error())
		return errorx.NewInternalErr()
	}

	newActivityID := uuid.String()

	err = dao.CreateNewActivity(ctx, &model.Activity{
		ActivityID:  newActivityID,
		Name:        in.Name,
		Description: in.Description,
		RouteID:     in.RouteID,
		CoverURL:    in.CoverURL,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		Tags:        in.Tags,
		NumberLimit: in.NumberLimit,
		Fee:         in.Fee,
	})
	if err != nil {
		zlog.Error("Error while create new activity"+err.Error(), zap.String("username", in.Name))
		return errorx.NewInternalErr()
	}

	return nil
}
