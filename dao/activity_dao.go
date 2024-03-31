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
