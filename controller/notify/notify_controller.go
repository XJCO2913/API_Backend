package notify

import (
	"context"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/notify"
	"api.backend.xjco2913/service/sdto"
	"github.com/gin-gonic/gin"
)

type NotifyController struct{}

func NewNotifyController() *NotifyController {
	return &NotifyController{}
}

func (n *NotifyController) Pull(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userID in token",
		})
		return
	}

	resp, sErr := notify.Service().Pull(context.Background(), userId)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Pull notification successfully",
		Data:       resp.NotificationList,
	})
}

func (n *NotifyController) ShareRoute(c *gin.Context) {
	var req dto.ShareRouteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong Params: " + err.Error(),
		})
		return
	}

	sErr := notify.Service().ShareRoute(context.Background(), &sdto.ShareRouteInput{
		ReceiverID: req.ReceiverID,
		RouteData:  req.RouteData,
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
		StatusMsg:  "Share route successfully",
	})
}

func (n *NotifyController) OrgResult(c *gin.Context) {
	var req dto.OrgResultReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong Params: " + err.Error(),
		})
		return
	}

	sErr := notify.Service().OrgResult(context.Background(), &sdto.OrgResultInput{
		ReceiverID: req.ReceiverID,
		IsAgreed:   req.IsAgreed,
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
		StatusMsg:  "Push organiser apply result successfully",
	})
}
