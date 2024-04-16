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

func CountParticipantsByActivityID(ctx context.Context, activityID string) (int64, error) {
	a := query.Use(DB).ActivityUser

	count, err := a.WithContext(ctx).Where(a.ActivityID.Eq(activityID)).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func FindUsersByActivityID(ctx context.Context, activityID string) ([]*model.ActivityUser, error) {
	a := query.Use(DB).ActivityUser

	activityUsers, err := a.WithContext(ctx).Where(a.ActivityID.Eq(activityID)).Find()
	if err != nil {
		return nil, err
	}

	return activityUsers, nil
}
