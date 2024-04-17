package user

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/microservice/kitex_gen/rpc/xjco2913/auth"
	"api.backend.xjco2913/microservice/kitex_gen/rpc/xjco2913/auth/authservice"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/user"
	"github.com/cloudwego/kitex/client"
	"api.backend.xjco2913/util"
	"github.com/gin-gonic/gin"
)

type UserController struct{
	authCli authservice.Client
}

func NewUserController() *UserController {
	return &UserController{
		authCli: authservice.MustNewClient("rpc.xjco2913.auth", client.WithHostPorts("43.136.232.116:8888")),
	}
}

func (u *UserController) SignUp(c *gin.Context) {
	var req dto.UserSignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong params: " + err.Error(),
		})
		return
	}

	// check the gender, gender must be 0, 1, 2
	if *req.Gender < 0 || *req.Gender > 2 {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong params: gender field must be 0 or 1 or 2",
		})
		return
	}

	err := user.Service().Create(c.Request.Context(), &sdto.CreateUserInput{
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
			"userInfo": gin.H{
				"username": req.Username,
				"gender":   req.Gender,
				"birthday": req.Birthday,
				"region":   req.Region,
			},
		},
	})
}

func (u *UserController) Login(c *gin.Context) {
	var req dto.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong params: " + err.Error(),
		})
		return
	}

	out, err := u.authCli.Login(context.Background(), &auth.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		// internal error
		c.JSON(500, dto.CommonRes{
			StatusCode: -1,
			StatusMsg: err.Error(),
		})
		return
	} else if out.BaseResp.Code != 0 {
		data := gin.H{
			"remaining_attempts": out.BaseResp.Data["remaining_attempts"],
		}
		if lockTimestamp, ok := out.BaseResp.Data["lock_expires"]; ok {
			data["lock_expires"] = lockTimestamp
		}

		c.JSON(int(out.BaseResp.Code), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  out.BaseResp.Msg,
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
				"username": out.Username,
				"gender":   out.Gender,
				"birthday": out.Birthday,
				"region":   out.Region,
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
		isOrganiser := false
		if user.OrganiserID != "" {
			isOrganiser = true
		}
		userInfos[i] = gin.H{
			"userId":         user.UserID,
			"username":       user.Username,
			"avatarUrl":      user.AvatarURL,
			"isOrganiser":    isOrganiser,
			"gender":         user.Gender,
			"birthday":       user.Birthday,
			"region":         user.Region,
			"membershipTime": user.MembershipTime,
			"membershipType": user.MembershipType,
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

	isOrganiser := userDetail.OrganiserID != ""
	responseData := gin.H{
		"userId":         userDetail.UserID,
		"username":       userDetail.Username,
		"avatarUrl":      userDetail.AvatarURL,
		"isOrganiser":    isOrganiser,
		"gender":         userDetail.Gender,
		"birthday":       userDetail.Birthday,
		"region":         userDetail.Region,
		"membershipTime": userDetail.MembershipTime,
		"membershipType": userDetail.MembershipType,
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
			StatusMsg:  "Wrong params: " + err.Error(),
		})
		return
	}

	if req.Gender != nil && (*req.Gender < 0 || *req.Gender > 2) {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong params: gender field must be 0 or 1 or 2",
		})
		return
	}

	input := sdto.UpdateUserInput{
		Username: req.Username,
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

func (u *UserController) Subscribe(c *gin.Context) {
	queryUserID := c.Query("userID")
	membershipTypeStr := c.Query("membershipType")

	currentUserID, currentUserExists := c.Get("userID")
	if !currentUserExists || queryUserID != currentUserID.(string) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: UserID mismatch",
		})
		return
	}

	membershipType, err := strconv.Atoi(membershipTypeStr)
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong membership type",
		})
		return
	}

	serviceErr := user.Service().Subscribe(c.Request.Context(), queryUserID, membershipType)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Subscribe successfully",
	})
}

func (u *UserController) CancelByID(c *gin.Context) {
	queryUserID := c.Query("userID")

	currentUserID, currentUserExists := c.Get("userID")
	if !currentUserExists || queryUserID != currentUserID.(string) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: UserID mismatch",
		})
		return
	}

	serviceErr := user.Service().CancelByID(c.Request.Context(), queryUserID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Cancel subscription successfully",
	})
}

func (u *UserController) UploadAvatar(c *gin.Context) {
	userId := c.PostForm("userId")
	avatarFileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  fmt.Sprintf("Bad avatar file header: %s", err.Error()),
		})
		return
	}

	currentUserID, currentUserExists := c.Get("userID")
	isAdmin, isAdminExists := c.Get("isAdmin")

	// Check if the current user is an administrator,
	// otherwise check if the requested userID is the same as the current userID.
	if !isAdminExists || !currentUserExists || (!isAdmin.(bool) && userId != currentUserID.(string)) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: You cannot access this resource",
		})
		return
	}

	avatarFile, err := avatarFileHeader.Open()
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  fmt.Sprintf("Fail to get avatar file: %s", err.Error()),
		})
		return
	}
	defer avatarFile.Close()

	avatarBuf := bytes.NewBuffer(nil)
	if _, err := io.Copy(avatarBuf, avatarFile); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  fmt.Sprintf("Fail copy avatar data: %s", err.Error()),
		})
		return
	}

	errx := user.Service().UploadAvatar(c.Request.Context(), sdto.UploadAvatarInput{
		UserId:     userId,
		AvatarData: avatarBuf.Bytes(),
	})
	if errx != nil {
		c.JSON(errx.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  errx.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Upload avatar successfully",
	})
}

func (u *UserController) RefreshToken(c *gin.Context) {
	userID := c.GetString("userID")
	if util.IsEmpty(userID) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg: "missing userID in token",
		})
		return
	}

	res, sErr := user.Service().RefreshToken(c.Request.Context(), userID)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg: sErr.Error(),
		})
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg: "refresh token successfully",
		Data: gin.H{
			"newToken": res.NewToken,
		},
	})
}
