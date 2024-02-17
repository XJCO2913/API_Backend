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
