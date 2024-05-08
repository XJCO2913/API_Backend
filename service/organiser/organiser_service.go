package organiser

import (
	"context"
	"errors"
	"fmt"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/dao/redis"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrganiserService struct{}

var (
	organiserService OrganiserService
)

func Service() *OrganiserService {
	return &organiserService
}

func (o *OrganiserService) GetAll(ctx context.Context) (*sdto.GetAllOrganisersOutput, *errorx.ServiceErr) {
	organiserModels, err := dao.GetAllOrganisers(ctx)
	if err != nil {
		zlog.Error("error while get all organisers", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// org status map
	statusMap := map[int32]string{
		-1: "refused",
		1:  "untreated",
		2:  "agreed",
	}

	res := make([]sdto.Organiser, len(organiserModels))
	for i, organiserModel := range organiserModels {
		// get user model
		userModel, err := dao.GetUserByID(ctx, organiserModel.UserID)
		if err != nil {
			zlog.Error("error while get user by id", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		// get avatar url
		var avatarUrl string
		if userModel.AvatarURL != nil && *userModel.AvatarURL != "" {
			avatarUrl, err = minio.GetUserAvatarUrl(ctx, *userModel.AvatarURL)
			if err != nil {
				zlog.Error("error while get user avatar url", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		res[i] = sdto.Organiser{
			UserID:         userModel.UserID,
			Username:       userModel.Username,
			AvatarUrl:      avatarUrl,
			MembershipTime: userModel.MembershipTime,
			Status:         statusMap[organiserModel.Status],
		}
	}

	return &sdto.GetAllOrganisersOutput{
		Organisers: res,
	}, nil
}

func (o *OrganiserService) Agree(ctx context.Context, userId string) *errorx.ServiceErr {
	userModel, err := dao.GetUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.NewServicerErr(
				400,
				"User not found",
				nil,
			)
		}

		zlog.Error("error while get user by id", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// update org status
	err = dao.UpdateOrgStatus(ctx, userId, 2)
	if err != nil {
		zlog.Error("error while update org status", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// async delete token cache
	go func() {
		cacheKey := fmt.Sprintf("jwt:%s", userModel.Username)
		redis.RDB().Del(ctx, cacheKey).Err()
	}()

	return nil
}

func (o *OrganiserService) Refuse(ctx context.Context, userId string) *errorx.ServiceErr {
	userModel, err := dao.GetUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.NewServicerErr(
				400,
				"User not found",
				nil,
			)
		}

		zlog.Error("error while get user by id", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// update org status
	err = dao.UpdateOrgStatus(ctx, userId, -1)
	if err != nil {
		zlog.Error("error while update org status", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// async delete token cache
	go func() {
		cacheKey := fmt.Sprintf("jwt:%s", userModel.Username)
		redis.RDB().Del(ctx, cacheKey).Err()
	}()

	return nil
}

func (o *OrganiserService) Apply(ctx context.Context, userId string) *errorx.ServiceErr {
	_, err := dao.GetUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.NewServicerErr(
				400,
				"User not found",
				nil,
			)
		}

		zlog.Error("error while get user by id", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// check if already exist
	_, err = dao.GetOrganiserByID(ctx, userId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Error("error while get org by userId", zap.Error(err))
			return errorx.NewInternalErr()
		}
	} else {
		return errorx.NewServicerErr(
			400,
			"You have already applied, please wait for admin review",
			nil,
		)
	}

	err = dao.CreateNewOrg(ctx, userId)
	if err != nil {
		zlog.Error("error while create new org record", zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}