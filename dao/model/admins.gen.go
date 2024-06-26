// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAdmin = "admins"

// Admin mapped from table <admins>
type Admin struct {
	ID        string     `gorm:"column:id;primaryKey" json:"id"`
	AvatarURL *string    `gorm:"column:avatarUrl" json:"avatarUrl"`
	Username  string     `gorm:"column:username;not null" json:"username"`
	Password  string     `gorm:"column:password;not null" json:"password"`
	CreatedAt *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updatedAt;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName Admin's table name
func (*Admin) TableName() string {
	return TableNameAdmin
}
