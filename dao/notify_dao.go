package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func PullNotificationByReceiverId(ctx context.Context, receiverId string) ([]*model.Notification, error) {
	n := query.Use(DB).Notification

	return n.WithContext(ctx).Where(n.ReceiverID.Eq(receiverId)).Find()
}