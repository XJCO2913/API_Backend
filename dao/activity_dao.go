package dao

import (
	"context"
	"strings"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateNewActivity(ctx context.Context, newActivity *model.Activity) error {
	err := query.Use(DB).WithContext(ctx).Activity.Create(newActivity)
	if err != nil {
		return err
	}

	return nil
}

func FindActivityByName(ctx context.Context, name string) (*model.Activity, error) {
	a := query.Use(DB).Activity

	activity, err := a.WithContext(ctx).Where(a.Name.Eq(name)).First()
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func GetAllActivities(ctx context.Context) ([]*model.Activity, error) {
	a := query.Use(DB).Activity

	activities, err := a.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func GetActivityByID(ctx context.Context, activityID string) (*model.Activity, error) {
	a := query.Use(DB).Activity

	activity, err := a.WithContext(ctx).Where(a.ActivityID.Eq(activityID)).First()
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func DeleteActivitiesByID(ctx context.Context, activityIDs string) ([]string, []string, error) {
	ids := strings.Split(activityIDs, "|")
	var deletedIDs []string
	var notFoundIDs []string

	for _, id := range ids {
		_, err := GetActivityByID(ctx, id)
		if err != nil {
			notFoundIDs = append(notFoundIDs, id)
			continue
		}

		a := query.Use(DB).Activity
		result, err := a.WithContext(ctx).Where(a.ActivityID.Eq(id)).Delete()
		if err != nil {
			return nil, nil, err
		}

		if result.RowsAffected > 0 {
			deletedIDs = append(deletedIDs, id)
		} else {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	return deletedIDs, notFoundIDs, nil
}

func GetActivityLimit(ctx context.Context, limit int) ([]*model.Activity, error) {
	a := query.Use(DB).Activity

	res, err := a.WithContext(ctx).Limit(limit).Order(a.CreatedAt.Asc()).Find()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetActivitiesByUserID(ctx context.Context, userID string) ([]*model.Activity, error) {
	var activityUsers []*model.ActivityUser
	var activities []*model.Activity

	a := query.Use(DB).ActivityUser

	activityUsers, err := a.WithContext(ctx).Where(a.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}

	activityIDs := make([]string, len(activityUsers))
	for i, userActivity := range activityUsers {
		activityIDs[i] = userActivity.ActivityID
	}

	if len(activityIDs) > 0 {
		a := query.Use(DB).Activity

		activities, err = a.WithContext(ctx).Where(a.ActivityID.In(activityIDs...)).Find()
		if err != nil {
			return nil, err
		}
	}

	return activities, nil
}

func GetActivitiesByCreatorID(ctx context.Context, creatorID string) ([]*model.Activity, error) {
	var activities []*model.Activity

	a := query.Use(DB).Activity

	activities, err := a.WithContext(ctx).Where(a.CreatorID.Eq(creatorID)).Find()
	if err != nil {
		return nil, err
	}

	return activities, nil
}
