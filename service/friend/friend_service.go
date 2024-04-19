package friend

import (
	"context"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
)

type FriendService struct{}

var (
	friendService FriendService
)

func Service() *FriendService {
	return &friendService
}

func (f *FriendService) Follow(ctx context.Context, in *sdto.FollowInput) (*errorx.ServiceErr) {
	// check if the user to follow exist or not
	if !dao.IsUserExisted(ctx, in.FollowingId) {
		return errorx.NewServicerErr(
			errorx.ErrExternal,
			"following user is not found",
			nil,
		)
	}

	// follow
	err := dao.FollowById(ctx, in.FollowerId, in.FollowingId)
	if err != nil {
		zlog.Error("error while store new follow relation", zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (f *FriendService) GetAllFollower(ctx context.Context, userId string) (*sdto.GetAllFollowerOutput, *errorx.ServiceErr) {
	followers, err := dao.GetFollowersByUserID(ctx, userId)
	if err != nil {
		zlog.Error("error while get all followers by user ID", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// get avatar
	for _, follower := range followers {
		if follower.AvatarURL == nil || *follower.AvatarURL == "" {
			continue
		}
		avatarUrl, err := minio.GetUserAvatarUrl(ctx, *follower.AvatarURL)
		if err != nil {
			zlog.Error("error while get user avatar", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		follower.AvatarURL = &avatarUrl
	} 

	return &sdto.GetAllFollowerOutput{
		Followers: followers,
	}, nil
}

func (f *FriendService) GetAllFollowing(ctx context.Context, userId string) (*sdto.GetAllFollowingOutput, *errorx.ServiceErr) {
	followings, err := dao.GetFollowingsByUserID(ctx, userId)
	if err != nil {
		zlog.Error("error while get followings by user id", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// get avatar
	for _, following := range followings {
		if following.AvatarURL == nil || *following.AvatarURL == "" {
			continue
		}
		avatarUrl, err := minio.GetUserAvatarUrl(ctx, *following.AvatarURL)
		if err != nil {
			zlog.Error("error while get user avatar", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		following.AvatarURL = &avatarUrl
	}

	return &sdto.GetAllFollowingOutput{
		Followings: followings,
	}, nil
}