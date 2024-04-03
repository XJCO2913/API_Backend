package dto

import "mime/multipart"

type CreateActivityReq struct {
	Name        string                `json:"name" binding:"required"`
	Description *string               `json:"description"`
	RouteID     int32                 `json:"routeId" binding:"required"`
	CoverFile   *multipart.FileHeader `json:"coverFile" binding:"required"`
	StartDate   string                `json:"startDate" binding:"required"`
	EndDate     string                `json:"endDate" binding:"required"`
	Tags        string                `json:"tags"`
	Level       string                `json:"level" binding:"required"`
}
