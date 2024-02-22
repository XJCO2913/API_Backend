package user

import (
	"time"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/user"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) SignUp(c *gin.Context) {
	var req dto.UserSignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "wrong params: " + err.Error(),
		})
		return
	}

	out, err := user.Service().Create(c.Request.Context(), &sdto.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Gender:   req.Gender,
		Region:   req.Region,
		Birthday: req.Birthday,
	})
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Register successfully",
		Data: gin.H{
			"token": out.Token,
			"userInfo": gin.H{
				"userId":         out.UserID,
				"username":       req.Username,
				"avatarUrl":      "",
				"isOrganiser":    0,
				"membershipTime": time.Now().Unix(),
				"gender":         req.Gender,
				"birthday":       req.Birthday,
				"region":         req.Region,
			},
		},
	})
}

func (u *UserController) Login(c *gin.Context) {
    var req dto.UserLoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, dto.CommonRes{
            StatusCode: -1,
            StatusMsg:  "wrong params: " + err.Error(),
        })
        return
    }

    out, err := user.Service().Authenticate(c.Request.Context(), &sdto.AuthenticateInput{
        Username: req.Username,
        Password: req.Password,
    })
    if err != nil {
		// Check whether the password is entered incorrectly or is disabled
        if authErr, ok := err.(*sdto.AuthError); ok {
			data := gin.H{
				"remaining_attempts": authErr.RemainingAttempts,
			}
			if !authErr.LockExpires.IsZero() {
				data["lock_expires"] = authErr.LockExpires.Unix()
			}
			c.JSON(401, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  authErr.Msg,
				Data:       data,
			})
		} else {
            c.JSON(401, dto.CommonRes{
                StatusCode: -1,
                StatusMsg:  err.Error(),
            })
        }
		return
    }

    c.JSON(200, dto.CommonRes{
        StatusCode: 0,
        StatusMsg:  "Login successfully",
        Data: gin.H{
            "token": out.Token,
            "userInfo": gin.H{
                "userId":         out.UserID,
                "username":       req.Username,
				"avatarUrl":      "",
				"isOrganiser":    0,
				"membershipTime": time.Now().Unix(),
				"gender":   	  out.Gender,
                "birthday": 	  out.Birthday,
                "region":   	  out.Region,
            },
        },
    })
}
