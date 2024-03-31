package sdto

type CreateUserInput struct {
	Username string
	Password string
	Gender   int32
	Birthday string
	Region   string
}

type AuthenticateInput struct {
	Username string
	Password string
}

type AuthenticateOutput struct {
	Token string `json:"token"`
}

type GetAllOutput struct {
	UserID         string `json:"userId"`
	Username       string `json:"username"`
	Gender         int32  `json:"gender"`
	Birthday       string `json:"birthday"`
	Region         string `json:"region"`
	MembershipTime int64  `json:"membershipTime"`
	AvatarURL      string `json:"avatarUrl"`
	OrganiserID    string `json:"organiserId"`
	MembershipType int32  `json:"membershipType"`
	IsSubscribed   int32  `json:"isSubscribed"`
}

type GetByIDOutput struct {
	UserID         string `json:"userId"`
	Username       string `json:"username"`
	Gender         int32  `json:"gender"`
	Birthday       string `json:"birthday"`
	Region         string `json:"region"`
	MembershipTime int64  `json:"membershipTime"`
	AvatarURL      string `json:"avatarUrl"`
	OrganiserID    string `json:"organiserId"`
	MembershipType int32  `json:"membershipType"`
	IsSubscribed   int32  `json:"isSubscribed"`
}

type GetAllStatusOutput struct {
	UserID   string `json:"userId"`
	IsBanned bool   `json:"isBanned"`
}

type UpdateUserInput struct {
	Username *string
	Password *string
	Gender   *int32
	Birthday *string
	Region   *string
}

type UploadAvatarInput struct {
	UserId     string
	AvatarData []byte
}

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
