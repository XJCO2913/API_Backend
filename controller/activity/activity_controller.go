package activity

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
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
	if err := c.ShouldBindJSON(&req); err != nil {
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

	// Get gpx file
	gpxFileHeader, err := c.FormFile("gpxFile")
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "GPX file is required",
		})
		return
	}

	gpxFile, err := gpxFileHeader.Open()
	if err != nil {
		c.JSON(500, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Failed to open gpx file",
		})
		return
	}
	defer gpxFile.Close()

	gpxBuf := bytes.NewBuffer(nil)
	if _, err := io.Copy(gpxBuf, gpxFile); err != nil {
		c.JSON(500, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  fmt.Sprintf("Failed to copy image data: %s", err.Error()),
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
		GPXData:     gpxBuf.Bytes(),
		StartDate:   startDate,
		EndDate:     endDate,
		Tags:        req.Tags,
		Level:       req.Level,
		CreatorID:   userID.(string),
	}

	sErr := activity.Service().Create(c.Request.Context(), input)
	if sErr != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
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

	participantsInfo := make([]gin.H, 0, len(activity.Participants))
	for _, participant := range activity.Participants {
		participantsInfo = append(participantsInfo, gin.H{
			"userID":         participant.UserID,
			"username":       participant.Username,
			"gender":         participant.Gender,
			"birthday":       participant.Birthday,
			"region":         participant.Region,
			"membershipTime": participant.MembershipTime,
			"avatarURL":      participant.AvatarURL,
			"membershipType": participant.MembershipType,
		})
	}

	responseData := gin.H{
		"activityId":        activity.ActivityID,
		"name":              activity.Name,
		"description":       activity.Description,
		"coverUrl":          activity.CoverURL,
		"media_gpx":         activity.GPXRoute,
		"startDate":         activity.StartDate,
		"endDate":           activity.EndDate,
		"tags":              activity.Tags,
		"numberLimit":       activity.NumberLimit,
		"originalFee":       activity.OriginalFee,
		"finalFee":          finalFee,
		"createdAt":         activity.CreatedAt,
		"creatorID":         activity.CreatorID,
		"creatorName":       activity.CreatorName,
		"participantsCount": activity.ParticipantsCount,
		"participants":      participantsInfo,
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

func (a *ActivityController) GetByCreatorID(c *gin.Context) {
	creatorID := c.Query("userID")

	currentCreatorID, currentCreatorExists := c.Get("userID")

	// Check if the requested CreatorID is the same as the current CreatorID.
	if !currentCreatorExists || creatorID != currentCreatorID.(string) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: You are not allowed to access other creators' activities",
		})
		return
	}

	activities, serviceErr := activity.Service().GetByCreatorID(c.Request.Context(), creatorID)
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
			"activityId":        activity.ActivityID,
			"name":              activity.Name,
			"description":       activity.Description,
			"coverUrl":          activity.CoverURL,
			"startDate":         activity.StartDate,
			"endDate":           activity.EndDate,
			"tags":              activity.Tags,
			"numberLimit":       activity.NumberLimit,
			"originalFee":       activity.OriginalFee,
			"createdAt":         activity.CreatedAt,
			"creatorID":         activity.CreatorID,
			"participantsCount": activity.ParticipantsCount,
		})
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get activity(ies) successfully",
		Data:       activitiesList,
	})
}

func (a *ActivityController) ProfitWithinDateRange(c *gin.Context) {
	startTimestampStr := c.Query("startTimestamp")
	endTimestampStr := c.Query("endTimestamp")

	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Forbidden: Only admins can access this resource",
		})
		return
	}

	startTimestamp, err := strconv.ParseInt(startTimestampStr, 10, 64)
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Invalid startTimestamp format",
		})
		return
	}

	endTimestamp, err := strconv.ParseInt(endTimestampStr, 10, 64)
	if err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Invalid endTimestamp format",
		})
		return
	}

	profit, serviceErr := activity.Service().ProfitWithinDateRange(c.Request.Context(), startTimestamp, endTimestamp)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	responseData := gin.H{
		"profit": profit,
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get activity profit successfully",
		Data:       responseData,
	})
}

func (a *ActivityController) TagsInfo(c *gin.Context) {
	resp, sErr := activity.Service().GetAllTagsInfo(c.Request.Context())
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get tags info successfully",
		Data: gin.H{
			"totalCount": resp.TotalCount,
			"eachCount":  resp.EachCount,
		},
	})
}

func (a *ActivityController) Counts(c *gin.Context) {
	resp, sErr := activity.Service().GetAllCounts(c.Request.Context())
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get all counts successfully",
		Data: gin.H{
			"activityCount":    resp.ActivityCount,
			"participantCount": resp.ParticipantCount,
			"membershipCount":  resp.MembershipCount,
		},
	})
}

func (a *ActivityController) GetProfitWithOption(c *gin.Context) {
	op := c.Query("option")
	if op != "week" && op != "month" && op != "year" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Invalid option",
		})
		return
	}

	resp, sErr := activity.Service().GetProfitWithOption(c.Request.Context(), op)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get profit successfully",
		Data: gin.H{
			"profits": resp.Profits,
			"dates":   resp.Dates,
		},
	})
}

func (ac *ActivityController) UploadRoute(c *gin.Context) {
	userID, userIDExists := c.Get("userID")
	if !userIDExists {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "User ID is required",
		})
		return
	}

	var req dto.UploadRouteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong params: " + err.Error(),
		})
		return
	}

	file, err := req.GPXData.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": -1,
			"status_msg":  "Failed to open GPX file",
		})
		return
	}
	defer file.Close()

	// Read file contents into byte slices
	gpxData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": -1,
			"status_msg":  "Failed to read GPX file",
		})
		return
	}

	input := &sdto.UploadRouteInput{
		UserID:     userID.(string),
		ActivityID: req.ActivityID,
		GPXData:    gpxData,
	}

	serviceErr := activity.Service().UploadRoute(c.Request.Context(), input)
	if serviceErr != nil {
		c.JSON(serviceErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  serviceErr.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Upload route successfully",
	})
}
