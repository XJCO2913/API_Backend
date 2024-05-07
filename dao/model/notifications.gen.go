// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameNotification = "notifications"

// Notification mapped from table <notifications>
type Notification struct {
	NotificationID string     `gorm:"column:notificationId;primaryKey" json:"notificationId"`
	ReceiverID     string     `gorm:"column:receiverId;not null" json:"receiverId"`
	SenderID       string     `gorm:"column:senderId;not null" json:"senderId"`
	RouteID        *int32     `gorm:"column:routeId" json:"routeId"`
	Type           int32      `gorm:"column:type;not null;default:1;comment:1 is admin notification, 2 is route notification" json:"type"` // 1 is admin notification, 2 is route notification
	Status         int32      `gorm:"column:status;not null;default:-1;comment:-1 is unread, 1 is read" json:"status"`                     // -1 is unread, 1 is read
	CreatedAt      *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt      *time.Time `gorm:"column:updatedAt;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName Notification's table name
func (*Notification) TableName() string {
	return TableNameNotification
}
