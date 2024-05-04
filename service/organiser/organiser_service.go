package organiser

import (
	"context"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
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
