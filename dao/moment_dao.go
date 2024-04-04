package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateNewMoment(ctx context.Context, newMoment *model.Moment) (int32, error) {
	err := query.Use(DB).WithContext(ctx).Moment.Create(newMoment)
	if err != nil {
		return -1, err
	}

	return newMoment.ID, nil
}

func DeleteMomentByID(ctx context.Context, momentID int32) error {
	_, err := query.Use(DB).WithContext(ctx).Moment.Delete(&model.Moment{ID: momentID})
	if err != nil {
		return err
	}

	return nil
}

func GetAllMoment(ctx context.Context) ([]*model.Moment, error) {
	return query.Use(DB).WithContext(ctx).Moment.Find()
}