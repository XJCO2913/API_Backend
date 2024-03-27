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

	_, err := user.Service().Create(c.Request.Context(), &sdto.CreateUserInput{
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
				"avatarUrl":      out.AvatarURL,
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
			StatusMsg:  "Forbidden: Only admins can access this resource",
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
			"avatarUrl":      user.AvatarURL,
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

	currentUserID, currentUserExists := c.Get("userID")
	isAdmin, isAdminExists := c.Get("isAdmin")

	// Check if the current user is an administrator,
	// otherwise check if the requested userID is the same as the current userID.
	if !isAdminExists || !currentUserExists || !isAdmin.(bool) && userID != currentUserID.(string) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
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
		"avatarUrl":      userDetail.AvatarURL,
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

func (u *UserController) DeleteByID(c *gin.Context) {
	userID := c.Query("userID")

	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	serviceErr := user.Service().DeleteByID(c.Request.Context(), userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Delete user(s) successfully",
	})
}

func (u *UserController) BanByID(c *gin.Context) {
	userID := c.Query("userID")

	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	serviceErr := user.Service().BanByID(c.Request.Context(), userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Ban User(s) successfully",
	})
}

func (u *UserController) UnbanByID(c *gin.Context) {
	userID := c.Query("userID")

	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	serviceErr := user.Service().UnbanByID(c.Request.Context(), userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Unban user(s) successfully",
	})
}

func (u *UserController) IsBanned(c *gin.Context) {
	userID := c.Query("userID")

	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	isBanned := user.Service().IsBanned(c.Request.Context(), userID)

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Check status successfully",
		Data: gin.H{
			"userId":   userID,
			"isBanned": isBanned,
		},
	})
}

func (u *UserController) GetAllStatus(c *gin.Context) {
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	userStatusList, serviceErr := user.Service().GetAllStatus(c.Request.Context())
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Check statuses successfully",
		Data:       userStatusList,
	})
}

func (u *UserController) UpdateByID(c *gin.Context) {
	userID := c.Query("userID")

	currentUserID, currentUserExists := c.Get("userID")
	isAdmin, isAdminExists := c.Get("isAdmin")

	// Check if the current user is an administrator,
	// otherwise check if the requested userID is the same as the current userID.
	if !isAdminExists || !currentUserExists || (!isAdmin.(bool) && userID != currentUserID.(string)) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "wrong params: " + err.Error(),
		})
		return
	}

	if req.Gender != nil && (*req.Gender < 0 || *req.Gender > 2) {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "wrong params: gender field must be 0 or 1 or 2",
		})
		return
	}

	input := sdto.UpdateUserInput{
		Username: req.Username,
		Password: req.Password,
		Gender:   req.Gender,
		Birthday: req.Birthday,
		Region:   req.Region,
	}

	serviceErr := user.Service().UpdateByID(c.Request.Context(), userID, input)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Update user successfully",
	})
}
