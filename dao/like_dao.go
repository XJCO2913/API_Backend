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

func GetLikeByID(ctx context.Context, userID, momentID string) (*model.Like, error) {
	a := query.Use(DB).Like

	like, err := a.WithContext(ctx).Where(
		a.UserID.Eq(userID),
		a.MomentID.Eq(momentID),
	).First()
	if err != nil {
		return nil, err
	}
	return like, nil
}
