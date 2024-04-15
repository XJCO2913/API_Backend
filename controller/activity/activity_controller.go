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
	isOrganiser, isOrganiserExists := c.Get("isOrganiser")
	userID, userIDExists := c.Get("userID")

	if !isOrganiserExists || !userIDExists || !isOrganiser.(bool) {
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
		CreatorID:   userID.(string),
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

func (a *ActivityController) GetAll(c *gin.Context) {
	userID, userIDExists := c.Get("userID")
	if !userIDExists {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "User ID is required",
		})
		return
	}

	userActivities, err := activity.Service().GetByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	registeredActivities := make(map[string]bool)
	for _, activity := range userActivities.Activities {
		registeredActivities[activity.ActivityID] = true
	}

	activities, err := activity.Service().GetAll(c.Request.Context())
	if err != nil {
		c.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	isAdmin, exists := c.Get("isAdmin")
	var discount int32 = 10
	if exists && isAdmin.(bool) {
		discount = 10
	} else {
		membershipTypeValue, exists := c.Get("membershipType")
		if !exists || membershipTypeValue == nil {
			c.JSON(500, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Membership type does not exist",
			})
			return
		}

		membershipType, convertSuccess := membershipTypeValue.(float64)
		if !convertSuccess {
			c.JSON(500, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Fail to convert MembershipType",
			})
			return
		}

		if membershipType == 2 {
			discount = 8
		}
	}

	activityInfos := make([]gin.H, len(activities))
	for i, activity := range activities {
		finalFee := activity.OriginalFee
		finalFee = finalFee * discount / 10

		// Check if the current user has registered for the activities
		isRegistered := registeredActivities[activity.ActivityID]

		activityInfos[i] = gin.H{
			"activityId":        activity.ActivityID,
			"name":              activity.Name,
			"description":       activity.Description,
			"coverUrl":          activity.CoverURL,
			"startDate":         activity.StartDate,
			"endDate":           activity.EndDate,
			"tags":              activity.Tags,
			"numberLimit":       activity.NumberLimit,
			"originalFee":       activity.OriginalFee,
			"finalFee":          finalFee,
			"createdAt":         activity.CreatedAt,
			"creatorID":         activity.CreatorID,
			"participantsCount": activity.ParticipantsCount,
			"isRegistered":      isRegistered,
		}
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get activities successfully",
		Data:       activityInfos,
	})
}

func (a *ActivityController) GetByID(c *gin.Context) {
	activityID := c.Query("activityID")

	userID, userIDExists := c.Get("userID")
	if !userIDExists {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "User ID is required",
		})
		return
	}

	userActivities, serviceErr := activity.Service().GetByUserID(c.Request.Context(), userID.(string))
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	// Check if the current user has registered for this activity
	isRegistered := false
	for _, userActivity := range userActivities.Activities {
		if userActivity.ActivityID == activityID {
			isRegistered = true
			break
		}
	}

	activity, serviceErr := activity.Service().GetByID(c.Request.Context(), activityID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	isAdmin, exists := c.Get("isAdmin")
	var discount int32 = 10
	if exists && isAdmin.(bool) {
		discount = 10
	} else {
		membershipTypeValue, exists := c.Get("membershipType")
		if !exists || membershipTypeValue == nil {
			c.JSON(500, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Membership type does not exist",
			})
			return
		}

		membershipType, convertSuccess := membershipTypeValue.(float64)
		if !convertSuccess {
			c.JSON(500, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Failed to convert MembershipType",
			})
			return
		}

		if membershipType == 2 {
			discount = 8
		}
	}

	finalFee := activity.OriginalFee
	finalFee = finalFee * discount / 10

	responseData := gin.H{
		"activityId":        activity.ActivityID,
		"name":              activity.Name,
		"description":       activity.Description,
		"coverUrl":          activity.CoverURL,
		"startDate":         activity.StartDate,
		"endDate":           activity.EndDate,
		"tags":              activity.Tags,
		"numberLimit":       activity.NumberLimit,
		"originalFee":       activity.OriginalFee,
		"finalFee":          finalFee,
		"createdAt":         activity.CreatedAt,
		"creatorID":         activity.CreatorID,
		"participantsCount": activity.ParticipantsCount,
		"isRegistered":      isRegistered,
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get activity successfully",
		Data:       responseData,
	})
}

func (a *ActivityController) Feed(c *gin.Context) {
	activities, err := activity.Service().Feed(c.Request.Context())
	if err != nil {
		c.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Feed activity successfully",
		Data:       activities,
	})
}

func (a *ActivityController) DeleteByID(c *gin.Context) {
	activityID := c.Query("activityID")

	isAdmin, isAdminExists := c.Get("isAdmin")
	if !isAdminExists || !isAdmin.(bool) {
		// Non-admins must be the creator to delete the activity
		activityDetail, serviceErr := activity.Service().GetByID(c.Request.Context(), activityID)
		if serviceErr != nil {
			c.JSON(serviceErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  serviceErr.Error(),
			})
			return
		}

		userID, userIDExists := c.Get("userID")
		if !userIDExists || activityDetail.CreatorID != userID.(string) {
			c.JSON(403, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Forbidden: You are not the creator of this activity",
			})
			return
		}
	}

	serviceErr := activity.Service().DeleteByID(c.Request.Context(), activityID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Delete activity(ies) successfully",
	})
}

func (a *ActivityController) SignUpByActivityID(c *gin.Context) {
	activityID := c.Query("activityID")

	userID, userIDExists := c.Get("userID")
	if !userIDExists {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "User ID does not exist",
		})
		return
	}

	membershipTypeValue, membershipTypeExists := c.Get("membershipType")
	if !membershipTypeExists {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Membership type does not exist",
		})
		return
	}

	membershipFloat, convertSuccess := membershipTypeValue.(float64)
	if !convertSuccess {
		c.JSON(500, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Failed to convert MembershipType",
		})
		return
	}
	membershipType := int64(membershipFloat)

	input := &sdto.SignUpActivityInput{
		UserID:         userID.(string),
		ActivityID:     activityID,
		MembershipType: membershipType,
	}

	serviceErr := activity.Service().SignUpByActivityID(c.Request.Context(), input)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Sign up for the activity successfully",
	})
}

func (a *ActivityController) GetByUserID(c *gin.Context) {
	userID := c.Query("userID")

	currentUserID, currentUserExists := c.Get("userID")

	// Check if the requested userID is the same as the current userID.
	if !currentUserExists || userID != currentUserID.(string) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: You are not allowed to access other users' activities",
		})
		return
	}

	activities, serviceErr := activity.Service().GetByUserID(c.Request.Context(), userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	var activitiesList []gin.H
	for _, activity := range activities.Activities {
		activitiesList = append(activitiesList, gin.H{
			"activityId":  activity.ActivityID,
			"name":        activity.Name,
			"description": activity.Description,
			"coverUrl":    activity.CoverURL,
			"startDate":   activity.StartDate,
			"endDate":     activity.EndDate,
			"tags":        activity.Tags,
			"numberLimit": activity.NumberLimit,
			"originalFee": activity.OriginalFee,
			"finalFee":    activity.FinalFee,
			"createdAt":   activity.CreatedAt,
			"creatorID":   activity.CreatorID,
		})
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get activity(ies) successfully",
		Data:       activitiesList,
	})
}
