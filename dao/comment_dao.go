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

func GetCommentsByMomentId(ctx context.Context, momentId string) ([]*model.Comment, error) {
	c := query.Use(DB).Comment

	return c.WithContext(ctx).Where(c.MomentID.Eq(momentId)).Find()
}
