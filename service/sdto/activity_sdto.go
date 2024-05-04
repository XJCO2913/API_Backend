package sdto

import (
	"time"
)

type CreateActivityInput struct {
	Name        string
	Description *string
	RouteID     int32
	CoverData   []byte
	GPXData     []byte
	StartDate   time.Time
	EndDate     time.Time
	Tags        string
	Level       string
	CreatorID   string
}

type GetAllActivityOutput struct {
	ActivityID        string
	Name              string
	Description       string
	RouteID           int32
	CoverURL          string
	StartDate         string
	EndDate           string
	Tags              string
	NumberLimit       int32
	OriginalFee       int32
	CreatedAt         string
	CreatorID         string
	ParticipantsCount int32
}

type ParticipantInfo struct {
	UserID         string
	Username       string
	Gender         int32
	Birthday       string
	Region         string
	MembershipTime int64
	AvatarURL      string
	MembershipType int32
}

type GetActivityByIDOutput struct {
	ActivityID        string
	Name              string
	Description       string
	RouteID           int32
	CoverURL          string
	GPXRoute          [][]string
	StartDate         string
	EndDate           string
	Tags              string
	NumberLimit       int32
	OriginalFee       int32
	CreatedAt         string
	CreatorID         string
	CreatorName       string
	ParticipantsCount int32
	Participants      []ParticipantInfo
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

type GetActivitiesByUserID struct {
	ActivityID  string
	Name        string
	Description string
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

type GetActivitiesByUserIDOutput struct {
	Activities []*GetActivitiesByUserID
}

type GetActivitiesByCreator struct {
	ActivityID        string
	Name              string
	Description       string
	CoverURL          string
	StartDate         string
	EndDate           string
	Tags              string
	NumberLimit       int32
	OriginalFee       int32
	CreatedAt         string
	CreatorID         string
	ParticipantsCount int32
}

type GetActivitiesByCreatorOutput struct {
	Activities []*GetActivitiesByCreator
}

type GetAllTagsInfoOutput struct {
	TotalCount int
	EachCount  map[string]int
}

type GetAllCountsOutput struct {
	ActivityCount    int
	ParticipantCount int
	MembershipCount  int
}

type GetProfitOutput struct {
	Profits []int
	Dates   []string
}

type UploadRouteInput struct {
	UserID     string
	ActivityID string
	GPXData    []byte
}

type GetRouteInput struct {
	ActivityID string
	UserID     string
}

type GetRouteOutput struct {
	GPXRouteText map[int][][]string
}
