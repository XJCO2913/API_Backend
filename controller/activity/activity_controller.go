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

type contextKey string

const (
	keyMembershipType contextKey = "membershipType"
)

func (a *ActivityController) Create(c *gin.Context) {
	isOrganiser, exists := c.Get("isOrganiser")
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

	membershipType, _ := c.Get("membershipType")
	err := activity.Service().Create(context.WithValue(c.Request.Context(), keyMembershipType, membershipType), input)
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

func (a *ActivityController) GetAll(ctx *gin.Context) {
	isAdmin, exists := ctx.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		ctx.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	activities, err := activity.Service().GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	activityInfos := make([]gin.H, len(activities))
	for i, activity := range activities {
		activityInfos[i] = gin.H{
			"activityId":  activity.ActivityID,
			"name":        activity.Name,
			"description": activity.Description,
			// "routeId":     activity.RouteID,
			"coverUrl":    activity.CoverURL,
			"startDate":   activity.StartDate,
			"endDate":     activity.EndDate,
			"tags":        activity.Tags,
			"numberLimit": activity.NumberLimit,
			"fee":         activity.Fee,
		}
	}

	ctx.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get activities successfully",
		Data:       activityInfos,
	})
}

func (a *ActivityController) GetByID(c *gin.Context) {
	activityID := c.Query("activityID")

	// If the user is not an administrator,
	// check whether the activity creator ID and user ID are consistent. (TBD)

	activityDetail, serviceErr := activity.Service().GetByID(c.Request.Context(), activityID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	responseData := gin.H{
		"activityId":  activityDetail.ActivityID,
		"name":        activityDetail.Name,
		"description": activityDetail.Description,
		// "routeId":     activityDetail.RouteID,
		"coverUrl":    activityDetail.CoverURL,
		"startDate":   activityDetail.StartDate,
		"endDate":     activityDetail.EndDate,
		"tags":        activityDetail.Tags,
		"numberLimit": activityDetail.NumberLimit,
		"fee":         activityDetail.Fee,
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get activity successfully",
		Data:       responseData,
	})
}
