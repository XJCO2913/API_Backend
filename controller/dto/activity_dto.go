package dto

type CreateActivityReq struct {
	Name        string  `form:"name" binding:"required"`
	Description *string `form:"description" binding:"required"`
	RouteID     int32   `form:"routeId"`
	StartDate   string  `form:"startDate" binding:"required"`
	EndDate     string  `form:"endDate" binding:"required"`
	Tags        string  `form:"tags"`
	Level       string  `form:"level" binding:"required,oneof=small medium"`
}
