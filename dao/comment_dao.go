package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateNewComment(ctx context.Context, newComment *model.Comment) error {
	err := query.Use(DB).Comment.WithContext(ctx).Create(newComment)
	if err != nil {
		return err
	}

	return nil
}
