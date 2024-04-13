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
