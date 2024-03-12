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

	// check the gender, gender must be 0, 1, 2
	if *req.Gender < 0 || *req.Gender > 2 {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "wrong params: gender field must be 0 or 1 or 2",
		})
		return
	}

	out, err := user.Service().Create(c.Request.Context(), &sdto.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Gender:   *req.Gender,
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
		data := gin.H{
			"remaining_attempts": err.Get("remaining_attempts"),
		}
		if t, ok := err.Get("lock_expires").(time.Time); ok {
			data["lock_expires"] = t.Unix()
		}

		c.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
			Data:       data,
		})
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
				"gender":         out.Gender,
				"birthday":       out.Birthday,
				"region":         out.Region,
			},
		},
	})
}

func (u *UserController) GetAll(ctx *gin.Context) {
	isAdmin, exists := ctx.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		ctx.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource.",
		})
		return
	}

	users, err := user.Service().GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	userInfos := make([]gin.H, len(users))
	for i, user := range users {
		userInfos[i] = gin.H{
			"userId":         user.UserID,
			"username":       user.Username,
			"avatarUrl":      "",
			"isOrganiser":    0,
			"gender":         user.Gender,
			"birthday":       user.Birthday,
			"region":         user.Region,
			"membershipTime": user.MembershipTime,
		}
	}

	ctx.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get users successfully",
		Data:       userInfos,
	})
}

func (u *UserController) GetByID(c *gin.Context) {
	userID := c.Query("userID")

	currentUserID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	// Check if the current user is an administrator,
	// otherwise check if the requested userID is the same as the current userID.
	if !isAdmin.(bool) && userID != currentUserID.(string) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: You do not have permission to access this resource.",
		})
		return
	}

	userDetail, serviceErr := user.Service().GetByID(c.Request.Context(), userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	responseData := gin.H{
		"userId":         userDetail.UserID,
		"username":       userDetail.Username,
		"avatarUrl":      "",
		"isOrganiser":    0,
		"gender":         userDetail.Gender,
		"birthday":       userDetail.Birthday,
		"region":         userDetail.Region,
		"membershipTime": userDetail.MembershipTime,
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get user successfully",
		Data:       responseData,
	})
}
