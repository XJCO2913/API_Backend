package sdto

import "time"

type NotifyUser struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type Notification struct {
	NotificationID string      `json:"notificationId"`
	Sender         *NotifyUser `json:"sender"`
	Route          [][]string  `json:"route"`
	Type           int32       `json:"type"`
	CreatedAt      *time.Time  `json:"createdAt"`
}

type PullNotificationOutput struct {
	NotificationList []*Notification
}
