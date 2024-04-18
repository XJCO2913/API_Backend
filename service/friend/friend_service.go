package friend

import (
	"context"

	"api.backend.xjco2913/dao"
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