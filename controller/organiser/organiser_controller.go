package organiser

import (
	"context"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/organiser"
	"github.com/gin-gonic/gin"
)

type OrganiserController struct{}

func NewOrganiserController() *OrganiserController {
	return &OrganiserController{}
}

func (o *OrganiserController) GetAll(c *gin.Context) {
	resp, sErr := organiser.Service().GetAll(context.Background())
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get all organisers successfully",
		Data:       resp.Organisers,
	})
}

func (o *OrganiserController) Agree(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userId",
		})
		return
	}

	sErr := organiser.Service().Agree(context.Background(), userId)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Agree org apply successfully",
	})
}

func (o *OrganiserController) Refuse(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userId",
		})
		return
	}

	sErr := organiser.Service().Refuse(context.Background(), userId)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Refuse org apply successfully",
	})
}

func (o *OrganiserController) Apply(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg: "Missing userID in token",
		})
		return
	}

	membershipType := c.GetFloat64("membershipType")
	if membershipType != 1 && membershipType != 2 {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg: "Only membership can apply for organiser",
		})
		return
	}

	sErr := organiser.Service().Apply(context.Background(), userId)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg: sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg: "Apply created successfully, waiting for administrator review",
	})
}