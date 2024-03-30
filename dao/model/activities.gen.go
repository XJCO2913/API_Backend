// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameActivity = "activities"

// Activity mapped from table <activities>
type Activity struct {
	ID          int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ActivityID  string     `gorm:"column:activityId;not null" json:"activityId"`
	Name        string     `gorm:"column:name;not null" json:"name"`
	Description *string    `gorm:"column:description" json:"description"`
	RouteID     int32      `gorm:"column:routeId;not null" json:"routeId"`
	CoverURL    string     `gorm:"column:coverUrl;not null" json:"coverUrl"`
	Type        int32      `gorm:"column:type;not null;comment:member only or not" json:"type"` // member only or not
	StartDate   time.Time  `gorm:"column:startDate;not null" json:"startDate"`
	EndDate     time.Time  `gorm:"column:endDate;not null" json:"endDate"`
	Tags        string     `gorm:"column:tags;not null" json:"tags"`
	NumberLimit int32      `gorm:"column:numberLimit;not null" json:"numberLimit"`
	Fee         int32      `gorm:"column:fee;not null;comment:can be free" json:"fee"` // can be free
	CreatedAt   *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"column:updatedAt;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName Activity's table name
func (*Activity) TableName() string {
	return TableNameActivity
}
