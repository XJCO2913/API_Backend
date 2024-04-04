package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateNewMoment(ctx context.Context, newMoment *model.Moment) error {
	return query.Use(DB).WithContext(ctx).Moment.Create(newMoment)
}