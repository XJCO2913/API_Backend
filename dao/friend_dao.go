package dao

import (
	"context"
	"errors"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
	"gorm.io/gorm"
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

	// get follower
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

func GetFollowingsByUserID(ctx context.Context, userId string) ([]*model.User, error) {
	f := query.Use(DB).Follow

	// get following
	follows, err := f.WithContext(ctx).Where(f.UserID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}

	followings := make([]*model.User, len(follows))
	for i, follow := range follows {
		following, err := GetUserByID(ctx, follow.FollowingID)
		if err != nil {
			return nil, err
		}

		followings[i] = following
	}

	return followings, nil
}

func GetFriendsByUserID(ctx context.Context, userId string) ([]*model.User, error) {
	f := query.Use(DB).Follow

	followers, err := f.WithContext(ctx).Where(f.FollowingID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}

	followings, err := f.WithContext(ctx).Where(f.UserID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}

	followersMap := make(map[string]bool)
	for _, follower := range followers {
		followersMap[follower.UserID] = true
	}

	friends := []*model.User{}
	for _, following := range followings {
		if !followersMap[following.FollowingID] {
			continue
		}

		friend, err := GetUserByID(ctx, following.FollowingID)
		if err != nil {
			return nil, err
		}

		friends = append(friends, friend)
	}

	return friends, nil
}

func CheckIsFollowed(ctx context.Context, userId, otherUserId string) (bool, error) {
	f := query.Use(DB).Follow

	_, err := f.WithContext(ctx).Where(f.UserID.Eq(userId)).Where(f.FollowingID.Eq(otherUserId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}