package admin

import (
	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/admin"
	"api.backend.xjco2913/service/sdto"
	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func NewAdminController() *AdminController {
	return &AdminController{}
}

func (a *AdminController) Login(c *gin.Context) {
	var req dto.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg: "wrong params: " + err.Error(),
		})
		return
	}

	resp, err := admin.Service().Authenticate(c.Request.Context(), &sdto.AdminAuthenticateInput{
		Name: req.Name,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg: err.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg: "Admin login successfully",
		Data: gin.H{
			"token": resp.Token,
			"name": resp.Name,
			"adminID": resp.AdminId,
		},
	})
}