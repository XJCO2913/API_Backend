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
