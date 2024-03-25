package sdto

type CreateUserInput struct {
	Username string
	Password string
	Gender   int32
	Birthday string
	Region   string
}

type CreateUserOutput struct {
	UserID string
	Token  string
}

type AuthenticateInput struct {
	Username string
	Password string
}

type AuthenticateOutput struct {
	UserID    string
	Token     string
	Gender    int32
	Birthday  string
	Region    string
	AvatarURL string
}

type GetAllOutput struct {
	UserID         string `json:"userId"`
	Username       string `json:"username"`
	Gender         int32  `json:"gender"`
	Birthday       string `json:"birthday"`
	Region         string `json:"region"`
	MembershipTime int64  `json:"membershipTime"`
	AvatarURL      string `json:"avatarUrl"`
}

type GetByIDOutput struct {
	UserID         string `json:"userId"`
	Username       string `json:"username"`
	Gender         int32  `json:"gender"`
	Birthday       string `json:"birthday"`
	Region         string `json:"region"`
	MembershipTime int64  `json:"membershipTime"`
	AvatarURL      string `json:"avatarUrl"`
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
