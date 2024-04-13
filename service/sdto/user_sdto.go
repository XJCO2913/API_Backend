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
	Token    string
	Gender   int32
	Birthday string
	Region   string
}

type GetAllOutput struct {
	UserID         string
	Username       string
	Gender         int32
	Birthday       string
	Region         string
	MembershipTime int64
	AvatarURL      string
	OrganiserID    string
	MembershipType int32
	IsSubscribed   int32
}

type GetByIDOutput struct {
	UserID         string
	Username       string
	Gender         int32
	Birthday       string
	Region         string
	MembershipTime int64
	AvatarURL      string
	OrganiserID    string
	MembershipType int32
	IsSubscribed   int32
}

type GetAllStatusOutput struct {
	UserID   string
	IsBanned bool
}

type UpdateUserInput struct {
	Username *string
	Gender   *int32
	Birthday *string
	Region   *string
}

type UploadAvatarInput struct {
	UserId     string
	AvatarData []byte
}

type RefreshTokenOutput struct {
	NewToken string
}
