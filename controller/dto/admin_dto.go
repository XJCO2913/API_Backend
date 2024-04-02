package dto

type AdminLoginReq struct {
	Name string `binding:"required"`
	Password string `binding:"required"`
}