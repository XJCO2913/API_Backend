package sdto

import (
	"time"
)

type CreateActivityInput struct {
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description" binding:"required"`
	RouteID     int32     `json:"routeId" binding:"required"`
	CoverData   []byte    `json:"coverData" binding:"required"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	EndDate     time.Time `json:"endDate" binding:"required"`
	Tags        string    `json:"tags"`
	Level       string    `json:"level" binding:"required"`
}

type GetAllActivityOutput struct {
	ActivityID  string `json:"activityId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	RouteID     int32  `json:"routeId"`
	CoverURL    string `json:"coverUrl"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Tags        string `json:"tags,omitempty"`
	NumberLimit int32  `json:"numberLimit"`
	Fee         int32  `json:"fee"`
}
