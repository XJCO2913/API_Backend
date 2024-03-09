// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID             int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID         string     `gorm:"column:userId;not null" json:"userId"`
	AvatarURL      *string    `gorm:"column:avatarUrl" json:"avatarUrl"`
	MembershipTime int64      `gorm:"column:membershipTime;not null;comment:membership expired time, a unix timestamp" json:"membershipTime"` // membership expired time, a unix timestamp
	Gender         int32      `gorm:"column:gender;not null;comment:0 is male, 1 is female, 2 is prefer-not-to-say" json:"gender"`            // 0 is male, 1 is female, 2 is prefer-not-to-say
	Region         string     `gorm:"column:region;not null" json:"region"`
	Tags           *string    `gorm:"column:tags;comment:Multiple tags are separated using '|'" json:"tags"` // Multiple tags are separated using '|'
	Birthday       *time.Time `gorm:"column:birthday" json:"birthday"`
	CreatedAt      *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt      *time.Time `gorm:"column:updatedAt;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	Username       string     `gorm:"column:username;not null" json:"username"`
	Password       string     `gorm:"column:password;not null" json:"password"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
