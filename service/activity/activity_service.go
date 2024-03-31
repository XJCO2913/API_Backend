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
)

type ActivityService struct{}

var (
	activityService ActivityService
)

func Service() *ActivityService {
	return &activityService
}

func (s *ActivityService) CreateActivity(ctx context.Context, in *sdto.CreateActivityInput) *errorx.ServiceErr {
	// Generate UUID for activityId
	activityId, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generating UUID for activity: ", zap.Error(err))
		return errorx.NewInternalErr()
	}

	newActivity := model.Activity{
		ActivityID:  activityId.String(),
		Name:        in.Name,
		Description: in.Description,
		RouteID:     in.RouteID,
		CoverURL:    in.CoverURL,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		Tags:        in.Tags,
		NumberLimit: in.NumberLimit,
		Fee:         in.Fee,
	}

	// Call DAO layer to create new activity
	if err := dao.CreateNewActivity(ctx, &newActivity); err != nil {
		zlog.Error("Failed to create new activity", zap.Any("activity", newActivity), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}
