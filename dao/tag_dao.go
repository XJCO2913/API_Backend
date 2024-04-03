package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func GetAllTags(ctx context.Context) ([]*model.Tag, error) {
	t := query.Use(DB).Tag

	tags, err := t.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	return tags, nil
}
