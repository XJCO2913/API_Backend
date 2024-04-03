package activity

import (
	"context"
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
	isOrganiser, exists := c.Get("isOrganiser")
	membershipType, exists := c.Get("isOrganiser")
	if !exists || !isOrganiser.(bool) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only organisers can access this resource",
		})
		return
	}

	var req dto.CreateActivityReq
	if err := c.Bind(&req); err != nil {
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

	startDate, serviceErr := time.Parse(time.DateOnly, req.StartDate)
	if serviceErr != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Invalid start date format",
		})
		return
	}

	endDate, serviceErr := time.Parse(time.DateOnly, req.EndDate)
	if serviceErr != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Invalid end date format",
		})
		return
	}

	// Check if the activity spans more than one year
	duration := endDate.Sub(startDate)
	if duration > 365*24*time.Hour {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "The duration of the activity cannot exceed one year",
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
		Level:       req.Level,
	}

	err := activity.Service().Create(context.WithValue(c.Request.Context(), "membershipType", membershipType), input)
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
