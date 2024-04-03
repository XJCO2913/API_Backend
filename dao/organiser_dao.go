package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func GetAllOrganisers(ctx context.Context) ([]*model.Organiser, error) {
	o := query.Use(DB).Organiser

	organisers, err := o.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	return organisers, nil
}

func GetOrganiserByID(ctx context.Context, userID string) (*model.Organiser, error) {
	o := query.Use(DB).Organiser

	organiser, err := o.WithContext(ctx).Where(o.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	return organiser, nil
}
