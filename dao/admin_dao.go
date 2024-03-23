package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func FindAdminByName(ctx context.Context, name string) (*model.Admin, error) {
	a := query.Use(DB).Admin

	admin, err := a.WithContext(ctx).Where(a.Username.Eq(name)).First()
	if err != nil {
		return nil, err
	}

	return admin, nil
}