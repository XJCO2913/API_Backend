package sdto

import "time"

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
    UserID   string
	Token    string
    Gender   int32
    Birthday string
    Region   string
}

type AuthError struct {
    Msg               string
    RemainingAttempts int64
    LockExpires       time.Time
}

func (e *AuthError) Error() string {
    return e.Msg
}