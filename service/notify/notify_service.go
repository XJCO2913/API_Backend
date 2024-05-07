package notify

import (
	"context"
	"errors"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotifyService struct{}

var (
	notifyService NotifyService
)

func Service() *NotifyService {
	return &notifyService
}

func (n *NotifyService) Pull(ctx context.Context, userId string) (*sdto.PullNotificationOutput, *errorx.ServiceErr) {
	notifications, err := dao.PullNotificationByReceiverId(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &sdto.PullNotificationOutput{
				NotificationList: []*sdto.Notification{},
			}, nil
		}

		zlog.Error("error while pull notification by user id", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	res := make([]*sdto.Notification, len(notifications))
	for i, notification := range notifications {
		res[i] = &sdto.Notification{}

		// get sender
		sender := sdto.NotifyUser{}
		senderModel, err := dao.GetUserByID(ctx, notification.SenderID)
		if err != nil {
			zlog.Error("error while get user by id", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
		if senderModel.AvatarURL != nil && *senderModel.AvatarURL != "" {
			sender.AvatarUrl, err = minio.GetUserAvatarUrl(ctx, *senderModel.AvatarURL)
			if err != nil {
				zlog.Error("error while get user avatar url", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}
		sender.UserID = senderModel.UserID
		sender.Username = senderModel.Username

		// get route data
		var routeData [][]string
		if notification.Type == 2 && notification.RouteID != nil {
			path, err := dao.GetPathAsText(ctx, *notification.RouteID)
			if err != nil {
				zlog.Error("Error while get GPX route from mysql", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			pathText, err := util.GPXRoute(path)
			if err != nil {
				zlog.Error("Error while parse gpx route to text", zap.String("path", path))
				return nil, errorx.NewInternalErr()
			}

			routeData = util.GPXStrTo2DString(pathText)
		}

		res[i].NotificationID = notification.NotificationID
		res[i].Sender = &sender
		res[i].Route = routeData
		res[i].Type = notification.Type
		res[i].CreatedAt = notification.CreatedAt

		// mark notification as read
		err = dao.ReadNotificationsById(ctx, notification.NotificationID)
		if err != nil {
			zlog.Error("error while read mark notification as read", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	return &sdto.PullNotificationOutput{
		NotificationList: res,
	}, nil
}
