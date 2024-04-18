package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func FollowById(ctx context.Context, followerId, followingId string) error {
	f := query.Use(DB).Follow

	newFollow := model.Follow{
		UserID:      followerId,
		FollowingID: followingId,
	}

	return f.WithContext(ctx).Create(&newFollow)
}
