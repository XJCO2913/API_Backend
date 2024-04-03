package dao

import (
	"context"

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

func GetActivityByID(ctx context.Context, id int32) (*model.Activity, error) {
	a := query.Use(DB).Activity

	activity, err := a.WithContext(ctx).Where(a.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return activity, nil
}
