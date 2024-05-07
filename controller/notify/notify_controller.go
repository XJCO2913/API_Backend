package notify

import (
	"context"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/notify"
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
