package friend

import (
	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/friend"
	"api.backend.xjco2913/service/sdto"
	"github.com/gin-gonic/gin"
)

type FriendController struct{}

func NewFriendController() *FriendController {
	return &FriendController{}
}

func (f *FriendController) Follow(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userID in token",
		})
		return
	}

	followId := c.Query("followingId")
	if followId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing follow user ID",
		})
		return
	}

	sErr := friend.Service().Follow(c.Request.Context(), &sdto.FollowInput{
		FollowerId:  userId,
		FollowingId: followId,
	})
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Follow user successfully",
	})
}

func (f *FriendController) GetAllFollower(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userID in token",
		})
		return
	}

	resp, sErr := friend.Service().GetAllFollower(c.Request.Context(), userId)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get followers successfully",
		Data:       resp.Followers,
	})
}

func (f *FriendController) GetAllFollowing(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userID in token",
		})
		return
	}

	resp, sErr := friend.Service().GetAllFollowing(c.Request.Context(), userId)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get all followings successfully",
		Data:       resp.Followings,
	})
}

func (f *FriendController) GetAll(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userID in token",
		})
		return
	}

	resp, sErr := friend.Service().GetAll(c.Request.Context(), userId)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: -1,
		StatusMsg:  "Get friends successfully",
		Data:       resp.Friends,
	})
}
