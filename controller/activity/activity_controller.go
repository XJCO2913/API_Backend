package activity

import (
	"io"
	"time"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/activity"
	"api.backend.xjco2913/service/sdto"
	"github.com/gin-gonic/gin"
)

type ActivityController struct{}

func NewActivityController() *ActivityController {
	return &ActivityController{}
}

func (a *ActivityController) Create(c *gin.Context) {
	var req dto.CreateActivityReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong params: " + err.Error(),
		})
		return
	}

	file, serviceErr := c.FormFile("coverFile")
	if serviceErr != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Cover file is required",
		})
		return
	}
	fileContent, serviceErr := file.Open()
	if serviceErr != nil {
		c.JSON(500, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Failed to open cover file",
		})
		return
	}
	defer fileContent.Close()

	coverData, serviceErr := io.ReadAll(fileContent)
	if serviceErr != nil {
		c.JSON(500, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Failed to read cover file",
		})
		return
	}

	startDate, serviceErr := time.Parse(time.RFC3339, req.StartDate)
	if serviceErr != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Invalid start date format",
		})
		return
	}

	endDate, serviceErr := time.Parse(time.RFC3339, req.EndDate)
	if serviceErr != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Invalid end date format",
		})
		return
	}

	input := &sdto.CreateActivityInput{
		Name:        req.Name,
		Description: req.Description,
		RouteID:     req.RouteID,
		CoverData:   coverData,
		StartDate:   startDate,
		EndDate:     endDate,
		Tags:        req.Tags,
		NumberLimit: req.NumberLimit,
	}

	err := activity.Service().Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Create activity successfully",
	})
}
