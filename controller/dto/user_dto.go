package dto

type CommonRes struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	Data       interface{}
}

type UserSignUpReq struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Gender   *int32 `binding:"required"`
	Birthday string
	Region   string `binding:"required"`
}

type UserLoginReq struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type UserUpdateReq struct {
	// allows for partial updates
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Gender   *int32  `json:"gender,omitempty"`
	Birthday *string `json:"birthday,omitempty"`
	Region   *string `json:"region,omitempty"`
}
