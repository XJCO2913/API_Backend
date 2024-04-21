package activity

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/service/gpx"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/zlog"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	ACTIVITY_FEED_LIMIT = 3
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

	coverName, uploadErr := a.UploadCover(ctx, in.CoverData)
	if uploadErr != nil {
		zlog.Error("Error while upload cover: "+uploadErr.Error(), zap.Error(uploadErr))
		return errorx.NewInternalErr()
	}

	// parse gpx data
	gpxResp, sErr := gpx.Service().ParseGPXData(ctx, &sdto.ParseGPXDataInput{
		GPXData: in.GPXData,
	})
	if sErr != nil {
		return sErr
	}

	// Generate a uuid for the new activity
	uuid, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generate uuid: " + err.Error())
		return errorx.NewInternalErr()
	}

	activityID := uuid.String()
	err = dao.CreateNewActivity(ctx, &model.Activity{
		ActivityID:  activityID,
		Name:        in.Name,
		Description: in.Description,
		RouteID:     gpxResp.RouteID,
		CoverURL:    coverName,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		Tags:        &joinedTags,
		NumberLimit: numberLimit,
		Fee:         finalFee,
		CreatorID:   in.CreatorID,
	})
	if err != nil {
		zlog.Error("Error while create activity: "+err.Error(), zap.String("name", in.Name))

		go func() {
			// Asynchronously delete the uploaded cover
			cleanupErr := minio.DeleteActivityCover(ctx, coverName)
			if cleanupErr != nil {
				zlog.Error("Failed to delete cover in Minio", zap.String("coverName", coverName), zap.Error(cleanupErr))
			}
		}()

		// delete gpx route record
		dao.DeleteRouteById(ctx, gpxResp.RouteID)

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

		var createdAtStr string
		if activity.CreatedAt != nil {
			createdAtStr = activity.CreatedAt.Format(time.RFC822)
		}

		participantsCount, err := dao.CountParticipantsByActivityID(ctx, activity.ActivityID)
		if err != nil {
			zlog.Error("Failed to count participants for the activity", zap.String("activityID", activity.ActivityID), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		activityDtos[i] = &sdto.GetAllActivityOutput{
			ActivityID:        activity.ActivityID,
			Name:              activity.Name,
			Description:       description,
			CoverURL:          coverURL,
			StartDate:         activity.StartDate.Format(time.RFC822),
			EndDate:           activity.EndDate.Format(time.RFC822),
			Tags:              tags,
			NumberLimit:       activity.NumberLimit,
			OriginalFee:       activity.Fee,
			CreatedAt:         createdAtStr,
			CreatorID:         activity.CreatorID,
			ParticipantsCount: int32(participantsCount),
		}
	}

	return activityDtos, nil
}

func (s *ActivityService) GetByID(ctx context.Context, activityID string) (*sdto.GetActivityByIDOutput, *errorx.ServiceErr) {
	activity, err := dao.GetActivityByID(ctx, activityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("Activity not found by activity ID", zap.String("activityID", activityID))
			return nil, errorx.NewServicerErr(errorx.ErrExternal, "Activity not found by activity ID", nil)
		} else {
			zlog.Error("Failed to retrieve activity by activity ID", zap.String("activityID", activityID), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	var description, tags string
	if activity.Description != nil {
		description = *activity.Description
	}

	if activity.Tags != nil {
		tags = *activity.Tags
	}

	coverURL := ""
	if activity.CoverURL != "" {
		coverURL, err = minio.GetActivityCoverUrl(ctx, activity.CoverURL)
		if err != nil {
			zlog.Error("Error while getting activity cover URL", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	var createdAtStr string
	if activity.CreatedAt != nil {
		createdAtStr = activity.CreatedAt.Format(time.RFC822)
	}

	participantsCount, err := dao.CountParticipantsByActivityID(ctx, activity.ActivityID)
	if err != nil {
		zlog.Error("Failed to count participants for activity", zap.String("activityID", activity.ActivityID), zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	participants, err := dao.GetUsersByActivityID(ctx, activityID)
	if err != nil {
		zlog.Error("Failed to retrieve participants for activity", zap.String("activityID", activityID), zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	creator, err := dao.GetUserByID(ctx, activity.CreatorID)
	if err != nil {
		zlog.Error("Failed to retrieve creator by ID", zap.String("userID", activity.CreatorID), zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	var participantInfos []sdto.ParticipantInfo
	for _, user := range participants {
		// get avatar url from minio
		avatarURL := ""
		if user.AvatarURL != nil || !util.IsEmpty(user.AvatarURL) {
			avatarURL, err = minio.GetUserAvatarUrl(ctx, *user.AvatarURL)
			if err != nil {
				return nil, errorx.NewInternalErr()
			}
		}

		participantInfos = append(participantInfos, sdto.ParticipantInfo{
			UserID:         user.UserID,
			Username:       user.Username,
			Gender:         user.Gender,
			Birthday:       user.Birthday.Format(time.RFC822),
			Region:         user.Region,
			MembershipTime: user.MembershipTime,
			AvatarURL:      avatarURL,
			MembershipType: user.MembershipType,
		})
	}

	// get gpx data
	path, err := dao.GetPathAsText(ctx, activity.RouteID)
	if err != nil {
		zlog.Error("error while get GPX route from mysql", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	pathText, err := util.GPXRoute(path)
	if err != nil {
		zlog.Error("error while parse gpx route to text", zap.String("path", path))
		return nil, errorx.NewInternalErr()
	}

	output := &sdto.GetActivityByIDOutput{
		ActivityID:        activity.ActivityID,
		Name:              activity.Name,
		Description:       description,
		CoverURL:          coverURL,
		GPXRoute:          util.GPXStrTo2DString(pathText),
		StartDate:         activity.StartDate.Format(time.RFC822),
		EndDate:           activity.EndDate.Format(time.RFC822),
		Tags:              tags,
		NumberLimit:       activity.NumberLimit,
		OriginalFee:       activity.Fee,
		CreatedAt:         createdAtStr,
		CreatorID:         activity.CreatorID,
		CreatorName:       creator.Username,
		ParticipantsCount: int32(participantsCount),
		Participants:      participantInfos,
	}

	return output, nil
}

func (s *ActivityService) DeleteByID(ctx context.Context, activityIDs string) *errorx.ServiceErr {
	ids := strings.Split(activityIDs, "|")
	deletedIDs, notFoundIDs, err := dao.DeleteActivitiesByID(ctx, activityIDs)

	if err != nil {
		zlog.Error("Failed to delete activities", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// All specified activities were not found
	if len(notFoundIDs) == len(ids) {
		zlog.Warn("All specified activities not found", zap.Strings("not_found_ids", notFoundIDs))
		return errorx.NewServicerErr(errorx.ErrExternal, "All specified activities not found", map[string]any{"not_found_ids": notFoundIDs})
	}

	zlog.Info("Specified activities deleted", zap.Strings("deleted_activity_ids", deletedIDs))
	// Part of specified activities were not found
	if len(notFoundIDs) > 0 {
		zlog.Warn("Some specified activities not found", zap.Strings("not_found_ids", notFoundIDs))
	}

	return nil
}

func (s *ActivityService) Feed(ctx context.Context) (*sdto.ActivityFeedOutput, *errorx.ServiceErr) {
	activitiesModels, err := dao.GetActivityLimit(ctx, ACTIVITY_FEED_LIMIT)
	if err != nil {
		zlog.Error("Error while feed activities", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	activities := make([]*sdto.ActivityFeed, len(activitiesModels))
	for i, activity := range activitiesModels {
		activities[i] = &sdto.ActivityFeed{
			ActivityID: activity.ActivityID,
			Name:       activity.Name,
		}

		if activity.Description != nil {
			activities[i].Description = *activity.Description
		}

		coverUrl, err := minio.GetActivityCoverUrl(ctx, activity.CoverURL)
		if err != nil {
			zlog.Error("Error while get activity cover URL", zap.Error(err), zap.String("activityID", activity.ActivityID))
			return nil, errorx.NewInternalErr()
		}
		activities[i].CoverUrl = coverUrl
	}

	return &sdto.ActivityFeedOutput{
		Activities: activities,
	}, nil
}

func (s *ActivityService) SignUpByActivityID(ctx context.Context, input *sdto.SignUpActivityInput) *errorx.ServiceErr {
	activity, err := dao.GetActivityByID(ctx, input.ActivityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("Activity not found by activity ID", zap.String("activityID", input.ActivityID))
			return errorx.NewServicerErr(errorx.ErrExternal, "Activity not found by activity ID", nil)
		} else {
			zlog.Error("Failed to retrieve activity by activity ID", zap.String("activityID", input.ActivityID), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	_, err = dao.FindActivityUserByIDs(ctx, input.ActivityID, input.UserID)
	if err == nil {
		zlog.Error("User already signed up for this activity", zap.String("userID", input.UserID), zap.String("activityID", input.ActivityID))
		return errorx.NewServicerErr(errorx.ErrExternal, "User already signed up for this activity", nil)
	}

	if activity.Fee > 0 && input.MembershipType == 0 {
		zlog.Error("Ordinary user attempts to sign up for a paid activity", zap.String("userID", input.UserID), zap.String("activityID", input.ActivityID))
		return errorx.NewServicerErr(errorx.ErrExternal, "Ordinary user cannot sign up for paid activities", nil)
	}

	finalFee := CalculateFinalFee(activity.Fee, input.MembershipType)
	if finalFee == -1 {
		zlog.Error("Invalid membership type or failed to calculate fee", zap.String("userID", input.UserID), zap.String("activityID", input.ActivityID))
		return errorx.NewInternalErr()
	}

	newUserActivity := &model.ActivityUser{
		ActivityID: input.ActivityID,
		UserID:     input.UserID,
		FinalFee:   finalFee,
	}
	err = dao.CreateActivityUser(ctx, newUserActivity)
	if err != nil {
		zlog.Error("Failed to create activity-user association", zap.String("userID", input.UserID), zap.String("activityID", input.ActivityID), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func CalculateFinalFee(baseFee int32, membershipType int64) int32 {
	switch membershipType {
	case 1:
		return baseFee
	case 2:
		return int32(float32(baseFee) * 0.8)
	default:
		return -1
	}
}

func (s *ActivityService) GetByUserID(ctx context.Context, userID string) (*sdto.GetActivitiesByUserIDOutput, *errorx.ServiceErr) {
	activities, err := dao.GetActivitiesByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("Activity not found by user ID", zap.String("userID", userID))
			return nil, errorx.NewServicerErr(errorx.ErrExternal, "Activity not found by user ID", nil)
		} else {
			zlog.Error("Failed to retrieve activities by user ID", zap.String("userID", userID), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	if len(activities) == 0 {
		return &sdto.GetActivitiesByUserIDOutput{Activities: []*sdto.GetActivitiesByUserID{}}, nil
	}

	var activitiesOutput []*sdto.GetActivitiesByUserID
	for _, activity := range activities {
		var description, tags, coverURL, createdAtStr string
		if activity.Description != nil {
			description = *activity.Description
		}

		if activity.Tags != nil {
			tags = *activity.Tags
		}

		// get cover url from minio
		if activity.CoverURL != "" {
			coverURL, err = minio.GetActivityCoverUrl(ctx, activity.CoverURL)
			if err != nil {
				zlog.Error("Error while get activity cover URL", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		if activity.CreatedAt != nil {
			createdAtStr = activity.CreatedAt.Format(time.RFC822)
		}

		activityUser, err := dao.FindActivityUserByIDs(ctx, activity.ActivityID, userID)
		if err != nil {
			zlog.Error("Failed to find activity user association", zap.String("userID", userID), zap.String("activityID", activity.ActivityID), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		activitiesOutput = append(activitiesOutput, &sdto.GetActivitiesByUserID{
			ActivityID:  activity.ActivityID,
			Name:        activity.Name,
			Description: description,
			CoverURL:    coverURL,
			StartDate:   activity.StartDate.Format(time.RFC822),
			EndDate:     activity.EndDate.Format(time.RFC822),
			Tags:        tags,
			NumberLimit: activity.NumberLimit,
			OriginalFee: activity.Fee,
			FinalFee:    activityUser.FinalFee,
			CreatedAt:   createdAtStr,
			CreatorID:   activity.CreatorID,
		})
	}

	return &sdto.GetActivitiesByUserIDOutput{Activities: activitiesOutput}, nil
}

func (s *ActivityService) GetByCreatorID(ctx context.Context, creatorID string) (*sdto.GetActivitiesByCreatorOutput, *errorx.ServiceErr) {
	activities, err := dao.GetActivitiesByCreatorID(ctx, creatorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("Activity not found by creator ID", zap.String("creatorID", creatorID))
			return nil, errorx.NewServicerErr(errorx.ErrExternal, "Activity not found by creator ID", nil)
		} else {
			zlog.Error("Failed to retrieve activities by creator ID", zap.String("creatorID", creatorID), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	if len(activities) == 0 {
		return &sdto.GetActivitiesByCreatorOutput{Activities: []*sdto.GetActivitiesByCreator{}}, nil
	}

	var activitiesOutput []*sdto.GetActivitiesByCreator
	for _, activity := range activities {
		var description, tags, coverURL, createdAtStr string
		if activity.Description != nil {
			description = *activity.Description
		}

		if activity.Tags != nil {
			tags = *activity.Tags
		}

		if activity.CoverURL != "" {
			coverURL, err = minio.GetActivityCoverUrl(ctx, activity.CoverURL)
			if err != nil {
				zlog.Error("Error while get activity cover URL", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		if activity.CreatedAt != nil {
			createdAtStr = activity.CreatedAt.Format(time.RFC822)
		}

		participantsCount, err := dao.CountParticipantsByActivityID(ctx, activity.ActivityID)
		if err != nil {
			zlog.Error("Failed to count participants for the activity", zap.String("activityID", activity.ActivityID), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		activitiesOutput = append(activitiesOutput, &sdto.GetActivitiesByCreator{
			ActivityID:        activity.ActivityID,
			Name:              activity.Name,
			Description:       description,
			CoverURL:          coverURL,
			StartDate:         activity.StartDate.Format(time.RFC822),
			EndDate:           activity.EndDate.Format(time.RFC822),
			Tags:              tags,
			NumberLimit:       activity.NumberLimit,
			OriginalFee:       activity.Fee,
			CreatedAt:         createdAtStr,
			CreatorID:         activity.CreatorID,
			ParticipantsCount: int32(participantsCount),
		})
	}

	return &sdto.GetActivitiesByCreatorOutput{Activities: activitiesOutput}, nil
}

func (s *ActivityService) ProfitWithinDateRange(ctx context.Context, startTimestamp, endTimestamp int64) (int32, *errorx.ServiceErr) {
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	// The time period must be more than a week and less than a year
	if end.Sub(start) < 7*24*time.Hour {
		zlog.Error("Date range is at least one week", zap.Int64("startTimestamp", startTimestamp), zap.Int64("endTimestamp", endTimestamp))
		return 0, errorx.NewServicerErr(errorx.ErrExternal, "Date range is at least one week", nil)
	}

	if end.Sub(start) > 365*24*time.Hour {
		zlog.Error("Date range is no more than one year", zap.Int64("startTimestamp", startTimestamp), zap.Int64("endTimestamp", endTimestamp))
		return 0, errorx.NewServicerErr(errorx.ErrExternal, "Date range is no more than one year", nil)
	}

	activities, err := dao.GetActivitiesWithinDateRange(ctx, start, end)
	if err != nil {
		zlog.Error("Failed to retrieve activities within date range", zap.Error(err))
		return 0, errorx.NewInternalErr()
	}

	if len(activities) == 0 {
		return 0, nil
	}

	var activityIDs []string
	for _, activity := range activities {
		activityIDs = append(activityIDs, activity.ActivityID)
	}
	activityIDsString := strings.Join(activityIDs, "|")

	activityUsers, err := dao.GetActivityUserByActivityIDs(ctx, activityIDsString)
	if err != nil {
		zlog.Error("Failed to retrieve activity users", zap.String("activityIDs", activityIDsString), zap.Error(err))
		return 0, errorx.NewInternalErr()
	}

	// Calculate total activity revenue
	var totalProfit int32 = 0
	for _, activityUser := range activityUsers {
		totalProfit += activityUser.FinalFee
	}

	return totalProfit, nil
}

func (s *ActivityService) GetAllTagsInfo(ctx context.Context) (*sdto.GetAllTagsInfoOutput, *errorx.ServiceErr) {
	activities, err := dao.GetAllActivities(ctx)
	if err != nil {
		zlog.Error("error while get all activities", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	var (
		totalCount = 0
		eachCount  = make(map[string]int)
	)
	for _, activity := range activities {
		if activity.Tags == nil || *activity.Tags == "" {
			continue
		}

		tagArr := strings.Split(*activity.Tags, "|")

		totalCount += len(tagArr)
		for _, tag := range tagArr {
			eachCount[tag]++
		}
	}

	return &sdto.GetAllTagsInfoOutput{
		TotalCount: totalCount,
		EachCount:  eachCount,
	}, nil
}
