package sdto

type CreateActivityInput struct {
	ActivityId  string `json:"activityId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RouteId     int    `json:"routeId"`
	CoverUrl    string `json:"coverUrl"`
	Type        int    `json:"type"` // 0 for public, 1 for members only
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Tags        string `json:"tags"`
	NumberLimit int    `json:"numberLimit"`
	Fee         int    `json:"fee"`
}
