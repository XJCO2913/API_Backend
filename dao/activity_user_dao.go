package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateActivityUser(ctx context.Context, newUserActivity *model.ActivityUser) error {
	err := query.Use(DB).WithContext(ctx).ActivityUser.Create(newUserActivity)
	if err != nil {
		return err
	}

	return nil
}

func FindActivityUserByIDs(ctx context.Context, activityID, userID string) (*model.ActivityUser, error) {
	a := query.Use(DB).ActivityUser
	activityUser, err := a.WithContext(ctx).Where(a.ActivityID.Eq(activityID), a.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	return activityUser, nil
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
