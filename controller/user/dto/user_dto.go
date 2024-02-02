package dto

type CreateUserReq struct {
	Name string `json:"name"`
}

type CreateUserRes struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	Data       interface{} `json:"data"`
}
