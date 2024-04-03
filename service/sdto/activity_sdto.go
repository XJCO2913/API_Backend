package sdto

import (
	"time"
)

type CreateActivityInput struct {
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description"`
	RouteID     int32     `json:"routeId" binding:"required"`
	CoverData   []byte    `json:"coverData"`
	Type        int32     `json:"type"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	EndDate     time.Time `json:"endDate" binding:"required"`
	Tags        string    `json:"tags"`
	NumberLimit int32     `json:"numberLimit" binding:"required"`
}
