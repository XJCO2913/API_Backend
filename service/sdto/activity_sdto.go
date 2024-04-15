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
	CreatorID   string
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
	OriginalFee int32
	FinalFee    int32
	CreatedAt   string
	CreatorID   string
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
	OriginalFee int32
	FinalFee    int32
	CreatedAt   string
	CreatorID   string
}

type ActivityFeed struct {
	ActivityID  string
	Name        string
	Description string
	CoverUrl    string
}

type ActivityFeedOutput struct {
	Activities []*ActivityFeed
}

type SignUpActivityInput struct {
	UserID         string
	ActivityID     string
	MembershipType int64
}

type GetActivitiesByUserIDOutput struct {
	Activities []*GetActivityByIDOutput
}
