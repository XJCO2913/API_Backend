package dao

import (
	"context"
	"strings"

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

func GetActivityUserByActivityIDs(ctx context.Context, activityIDs string) ([]*model.ActivityUser, error) {
	ids := strings.Split(activityIDs, "|")
	var activityUsers []*model.ActivityUser

	a := query.Use(DB).ActivityUser

	activityUsers, err := a.WithContext(ctx).Where(a.ActivityID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	return activityUsers, nil
}

func ActivityUserCount(ctx context.Context) (int64, error) {
	au := query.Use(DB).ActivityUser

	return au.WithContext(ctx).Count()
}

func GetFinalFeesByActivityId(ctx context.Context, activityId string) ([]*model.ActivityUser, error) {
	au := query.Use(DB).ActivityUser

	return au.WithContext(ctx).Where(au.ActivityID.Eq(activityId)).Find()
}

func UpdateActivityUserRoute(ctx context.Context, activityID, userID string, routeID int32) error {
	a := query.Use(DB).ActivityUser

	activityUser, err := FindActivityUserByIDs(ctx, activityID, userID)
	if err != nil {
		return err
	}

	activityUser.RouteID = &routeID
	err = a.WithContext(ctx).Save(activityUser)
	if err != nil {
		return err
	}

	return nil
}
