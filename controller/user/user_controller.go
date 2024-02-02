package user

import (
	"net/http"

	"api.backend.xjco2913/controller/user/dto"
	"api.backend.xjco2913/service/user"
	"api.backend.xjco2913/service/user/sdto"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	
}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.CreateUserRes{
			StatusCode: -1,
			StatusMsg: "missing request params",
		})
		return
	}

	newUserID, err := user.Service().CreateUser(c.Request.Context(), &sdto.CreateUserInput{
		Name: req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.CreateUserRes{
			StatusCode: -1,
			StatusMsg: err.Error(),
		})
	}

	c.JSON(http.StatusOK, &dto.CreateUserRes{
		StatusCode: 0,
		StatusMsg: "successfully",
		Data: gin.H{
			"new_user_id": newUserID,
		},
	})
}