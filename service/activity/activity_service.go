package activity

import (
	"context"
	"strconv"
	"strings"

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

	var ExtraFee int32
	var joinedTags string
	if in.Tags != "" {
		// Split tag IDs and accumulate extra fee
		tagIDs := strings.Split(in.Tags, "|")
		var validTagIDs []string
		for _, tagID := range tagIDs {
			id, err := strconv.Atoi(tagID)
			if err != nil {
				zlog.Error("Failed to convert tagID to int", zap.String("tagID", tagID), zap.Error(err))
				continue
			}

			tag, err := dao.GetTagByID(ctx, int32(id))
			if err != nil {
				zlog.Error("Failed to retrieve tag by ID", zap.Int("tagID", id), zap.Error(err))
				continue
			}

			ExtraFee += tag.Price
			validTagIDs = append(validTagIDs, tagID)
		}
		if len(validTagIDs) > 0 {
			joinedTags = strings.Join(validTagIDs, "|")
		}
	}

	var baseFee int32
	var numberLimit int32
	// Set number limit and basic fee based on level
	switch in.Level {
	case "small":
		numberLimit = 10
		baseFee = 0
	case "medium":
		numberLimit = 30
		baseFee = 10
	// case "large": TBD
	default:
		zlog.Error("Unsupported level: " + in.Level)
		return errorx.NewServicerErr(
			errorx.ErrExternal,
			"Unsupported activity level",
			nil,
		)
	}

	finalFee := baseFee + ExtraFee

	membershipType, exists := ctx.Value("membershipType").(int)
	if !exists {
		zlog.Error("MembershipType not found")
	} else {
		// 20% off for second level members
		if membershipType == 2 {
			finalFee = finalFee * 8 / 10
		}
	}

	coverName, uploadErr := a.UploadCover(ctx, in.CoverData)
	if uploadErr != nil {
		zlog.Error("Error while upload cover: "+uploadErr.Error(), zap.Error(uploadErr))
		return errorx.NewInternalErr()
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
		RouteID:     1,
		CoverURL:    coverName,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		Tags:        &joinedTags,
		NumberLimit: numberLimit,
		Fee:         finalFee,
	})
	if err != nil {
		zlog.Error("Error while create new activity: "+err.Error(), zap.String("name", in.Name))
		return errorx.NewInternalErr()
	}

	return nil
}

func (a *ActivityService) UploadCover(ctx context.Context, coverData []byte) (string, *errorx.ServiceErr) {
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

func (s *ActivityService) GetAll(ctx context.Context) ([]*sdto.GetAllActivityOutput, *errorx.ServiceErr) {
	activities, err := dao.GetAllActivities(ctx)
	if err != nil {
		zlog.Error("Failed to retrieve all activities", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	activityDtos := make([]*sdto.GetAllActivityOutput, len(activities))
	for i, activity := range activities {
		var description, tags string
		if activity.Description != nil {
			description = *activity.Description
		}
		if activity.Tags != nil {
			tags = *activity.Tags
		}

		// get cover url from minio
		coverURL := ""
		if activity.CoverURL != "" {
			coverURL, err = minio.GetActivityCoverUrl(ctx, activity.CoverURL)
			if err != nil {
				zlog.Error("Error while get activity cover URL", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		activityDtos[i] = &sdto.GetAllActivityOutput{
			ActivityID:  activity.ActivityID,
			Name:        activity.Name,
			Description: description,
			RouteID:     activity.RouteID,
			CoverURL:    coverURL,
			StartDate:   activity.StartDate.Format("2006-01-02"),
			EndDate:     activity.EndDate.Format("2006-01-02"),
			Tags:        tags,
			NumberLimit: activity.NumberLimit,
			Fee:         activity.Fee,
		}
	}

	return activityDtos, nil
}
