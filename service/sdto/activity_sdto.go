package sdto

import (
	"time"
)

type CreateActivityInput struct {
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	RouteID     int32     `json:"routeId"`
	CoverURL    string    `json:"coverUrl"`
	Type        int32     `json:"type"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Tags        string    `json:"tags"`
	NumberLimit int32     `json:"numberLimit"`
	Fee         int32     `json:"fee"`
}
