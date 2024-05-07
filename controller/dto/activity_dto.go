package dto

import "mime/multipart"

type CreateActivityReq struct {
	Name        string  `form:"name" binding:"required"`
	Description *string `form:"description" binding:"required"`
	RouteID     int32   `form:"routeId"`
	StartDate   string  `form:"startDate" binding:"required"`
	EndDate     string  `form:"endDate" binding:"required"`
	Tags        string  `form:"tags"`
	Level       string  `form:"level" binding:"required,oneof=small medium"`
}

type UploadRouteReq struct {
	ActivityID string                `form:"activityId" binding:"required"`
	GPXData    *multipart.FileHeader `form:"gpxData" binding:"required"`
}

type GetRouteReq struct {
	ActivityID string `json:"activityId" binding:"required"`
	UserID     string `json:"userId" binding:"required"`
}
