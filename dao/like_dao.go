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

func GetLikeByIDs(ctx context.Context, userID, momentID string) (*model.Like, error) {
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

func DeleteLikeByIDs(ctx context.Context, userID, momentID string) error {
	a := query.Use(DB).Like

	_, err := a.WithContext(ctx).Where(
		a.UserID.Eq(userID),
		a.MomentID.Eq(momentID),
	).Delete()
	if err != nil {
		return err
	}

	return nil
}

func GetLikeByMomentId(ctx context.Context, momentId string) ([]*model.Like, error) {
	l := query.Use(DB).Like

	return l.WithContext(ctx).Where(l.MomentID.Eq(momentId)).Find()
}
