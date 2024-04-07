package sdto

import (
	"time"
)

type CreateActivityInput struct {
	Name        string
	Description *string
	RouteID     int32
	CoverData   []byte
	StartDate   time.Time
	EndDate     time.Time
	Tags        string
	Level       string
}

type GetAllActivityOutput struct {
	ActivityID  string
	Name        string
	Description string
	RouteID     int32
	CoverURL    string
	StartDate   string
	EndDate     string
	Tags        string
	NumberLimit int32
	Fee         int32
	CreatedAt   string
}

type GetActivityByIDOutput struct {
	ActivityID  string
	Name        string
	Description string
	RouteID     int32
	CoverURL    string
	StartDate   string
	EndDate     string
	Tags        string
	NumberLimit int32
	Fee         int32
	CreatedAt   string
}
