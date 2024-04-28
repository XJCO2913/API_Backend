package dao

import (
	"context"
	"time"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateNewMoment(ctx context.Context, newMoment *model.Moment) (int32, error) {
	err := query.Use(DB).Moment.WithContext(ctx).Create(newMoment)
	if err != nil {
		return -1, err
	}

	return newMoment.ID, nil
}

func DeleteMomentByID(ctx context.Context, momentID int32) error {
	_, err := query.Use(DB).Moment.WithContext(ctx).Delete(&model.Moment{ID: momentID})
	if err != nil {
		return err
	}

	return nil
}

func GetAllMoment(ctx context.Context) ([]*model.Moment, error) {
	return query.Use(DB).WithContext(ctx).Moment.Find()
}

func GetMomentByID(ctx context.Context, momentID string) (*model.Moment, error) {
	m := query.Use(DB).Moment

	moment, err := m.WithContext(ctx).Where(m.MomentID.Eq(momentID)).First()
	if err != nil {
		return nil, err
	}
	return moment, nil
}

func GetMomentsByTime(ctx context.Context, limit int, latestTime time.Time) ([]*model.Moment, error) {
	m := query.Use(DB).Moment

	moments, err := m.WithContext(ctx).Limit(limit).Order(m.CreatedAt.Desc()).Where(m.CreatedAt.Lt(latestTime)).Find()
	if err != nil {
		return nil, err
	}

	return moments, nil
}
