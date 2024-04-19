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

func GetFollowersByUserID(ctx context.Context, userId string) ([]*model.User, error) {
	f := query.Use(DB).Follow

	follows, err := f.WithContext(ctx).Where(f.FollowingID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}

	followers := make([]*model.User, len(follows))
	for i, follow := range follows {
		follower, err := GetUserByID(ctx, follow.UserID)
		if err != nil {
			return nil, err
		}

		followers[i] = follower
	}

	return followers, nil
}