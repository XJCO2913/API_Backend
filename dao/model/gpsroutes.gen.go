// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameGPSRoute = "GPSRoutes"

// GPSRoute mapped from table <GPSRoutes>
type GPSRoute struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Path string `gorm:"column:path;not null" json:"path"`
}

// TableName GPSRoute's table name
func (*GPSRoute) TableName() string {
	return TableNameGPSRoute
}
