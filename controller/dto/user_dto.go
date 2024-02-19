package dto

type CommonRes struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	Data       interface{}
}

type UserSignUpReq struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Gender   int32
	Birthday string
	Region   string `binding:"required"`
}

type UserLoginReq struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}
