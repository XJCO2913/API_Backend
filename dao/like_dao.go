package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateNewLike(ctx context.Context, newLike *model.Like) error {
	err := query.Use(DB).Like.WithContext(ctx).Create(newLike)
	if err != nil {
		return err
	}

	return nil
}
