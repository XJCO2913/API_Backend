package dto

type CreateActivityReq struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	RouteID     int32   `json:"routeId" binding:"required"`
	CoverURL    string  `json:"coverUrl" binding:"required"`
	StartDate   string  `json:"startDate" binding:"required"`
	EndDate     string  `json:"endDate" binding:"required"`
	Tags        string  `json:"tags" binding:"required"`
	NumberLimit int32   `json:"numberLimit" binding:"required"`
}
