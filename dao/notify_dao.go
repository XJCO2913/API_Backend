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

func ReadNotificationsById(ctx context.Context, notificationId string) error {
	n := query.Use(DB).Notification

	notification, err := n.WithContext(ctx).Where(n.NotificationID.Eq(notificationId)).First()
	if err != nil {
		return err
	}
	// mark notification as read
	if notification.Status == 1 {
		return nil
	}
	notification.Status = 1
	_, err = n.WithContext(ctx).Where(n.NotificationID.Eq(notificationId)).Updates(notification)
	if err != nil {
		return err
	}

	return nil
}

func PushNotification(ctx context.Context, newNotification *model.Notification) error {
	n := query.Use(DB).Notification

	return n.WithContext(ctx).Create(newNotification)
}

func GetUnreadNotificationByUserId(ctx context.Context, userId string) ([]*model.Notification, error) {
	n := query.Use(DB).Notification

	return n.WithContext(ctx).Where(n.ReceiverID.Eq(userId)).Where(n.Status.Eq(-1)).Find()
}