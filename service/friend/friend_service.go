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

func (f *FriendService) Follow(ctx context.Context, in *sdto.FollowInput) *errorx.ServiceErr {
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
	res := make([]*sdto.Follower, len(followers))
	for i, follower := range followers {
		var avatarUrl string
		if follower.AvatarURL != nil && *follower.AvatarURL != "" {
			avatarUrl, err = minio.GetUserAvatarUrl(ctx, *follower.AvatarURL)
			if err != nil {
				zlog.Error("error while get user avatar", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		// check if isFollowed
		isFollowed, err := dao.CheckIsFollowed(ctx, userId, follower.UserID)
		if err != nil {
			zlog.Error("error while check if follow follower or not", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		res[i] = &sdto.Follower{}
		res[i].UserID = follower.UserID
		res[i].Username = follower.Username
		res[i].AvatarUrl = avatarUrl
		res[i].Region = follower.Region
		res[i].IsFollowed = isFollowed
	}

	return &sdto.GetAllFollowerOutput{
		Followers: res,
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

func (f *FriendService) GetAll(ctx context.Context, userId string) (*sdto.GetAllFriendsOutput, *errorx.ServiceErr) {
	friends, err := dao.GetFriendsByUserID(ctx, userId)
	if err != nil {
		zlog.Error("error while get user friends", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// get avatar
	for _, friend := range friends {
		if friend.AvatarURL == nil || *friend.AvatarURL == "" {
			continue
		}
		avatarUrl, err := minio.GetUserAvatarUrl(ctx, *friend.AvatarURL)
		if err != nil {
			zlog.Error("error while get user avatar", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		friend.AvatarURL = &avatarUrl
	}

	return &sdto.GetAllFriendsOutput{
		Friends: friends,
	}, nil
}
