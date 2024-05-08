package moment

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/moment"
	"api.backend.xjco2913/service/sdto"
	"github.com/gin-gonic/gin"
)

type MomentController struct{}

func NewMomentController() *MomentController {
	return &MomentController{}
}

func (m *MomentController) Create(c *gin.Context) {
	userId, ok := c.Get("userID")
	if !ok {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing userID in token",
		})
		return
	}

	content := c.PostForm("content")

	videoFileHeader, errVideo := c.FormFile("videoFile")
	imageFileHeader, errImage := c.FormFile("imageFile")
	GPXFileHeader, errGPX := c.FormFile("gpxFile")

	// Get the number of files and available file headers
	fileHeaderMap := map[string]*multipart.FileHeader{
		"gpxFile":   GPXFileHeader,
		"imageFile": imageFileHeader,
		"videoFile": videoFileHeader,
	}

	fileCnt := 0
	validFileHeader := ""
	if errGPX == nil {
		fileCnt++
		validFileHeader = "gpxFile"
	}
	if errImage == nil {
		fileCnt++
		validFileHeader = "imageFile"
	}
	if errVideo == nil {
		fileCnt++
		validFileHeader = "videoFile"
	}

	// Process files
	// 1. When the content is empty, only one file can be transferred. Any other number of files (including 0) is a logical error.
	// 2. If the content is not empty, then the file does not need to be transferred.
	switch {
	case fileCnt == 0:
		if content == "" {
			// The number of files is 0
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Moment cannot be empty",
			})
			return
		}

		validFileHeader = "null"

	case fileCnt > 1:
		// Too many files
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Only allow upload one type of file",
		})
		return
	}

	fileHeader := fileHeaderMap[validFileHeader]
	switch validFileHeader {
	case "gpxFile":
		// Parser gpx file logic
		gpxFile, err := fileHeader.Open()
		if err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Failed to get gpx file: %s", err.Error()),
			})
			return
		}
		defer gpxFile.Close()

		gpxBuf := bytes.NewBuffer(nil)
		if _, err := io.Copy(gpxBuf, gpxFile); err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Failed copy image data: %s", err.Error()),
			})
			return
		}

		sErr := moment.Service().CreateWithGPX(context.Background(), &sdto.CreateMomentGPXInput{
			UserID:  userId.(string),
			Content: content,
			GPXData: gpxBuf.Bytes(),
		})
		if sErr != nil {
			c.JSON(sErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  sErr.Error(),
			})
			return
		}

	case "imageFile":
		imageFile, err := fileHeader.Open()
		if err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Failed to get image file: %s", err.Error()),
			})
			return
		}
		defer imageFile.Close()

		imageBuf := bytes.NewBuffer(nil)
		if _, err := io.Copy(imageBuf, imageFile); err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Failed to copy image data: %s", err.Error()),
			})
			return
		}

		sErr := moment.Service().CreateWithImage(c.Request.Context(), &sdto.CreateMomentImageInput{
			UserID:    userId.(string),
			Content:   content,
			ImageData: imageBuf.Bytes(),
		})
		if sErr != nil {
			c.JSON(sErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  sErr.Error(),
			})
			return
		}

	case "videoFile":
		videoFile, err := fileHeader.Open()
		if err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Failed to get video file: %s", err.Error()),
			})
			return
		}
		defer videoFile.Close()

		videoBuf := bytes.NewBuffer(nil)
		if _, err := io.Copy(videoBuf, videoFile); err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Failed copy video data: %s", err.Error()),
			})
			return
		}

		sErr := moment.Service().CreateWithVideo(context.Background(), &sdto.CreateMomentVideoInput{
			UserID:    userId.(string),
			Content:   content,
			VideoData: videoBuf.Bytes(),
		})
		if sErr != nil {
			c.JSON(sErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  sErr.Error(),
			})
			return
		}

	case "null":
		// Only content, no files
		sErr := moment.Service().Create(c.Request.Context(), &sdto.CreateMomentInput{
			UserID:  userId.(string),
			Content: content,
		})
		if sErr != nil {
			c.JSON(sErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  sErr.Error(),
			})
			return
		}
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Create new moment successfully",
	})
}

func (m *MomentController) Feed(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing user Id in token",
		})
		return
	}

	var latestTime int64
	var err error
	if latest := c.Query("latestTime"); latest == "" {
		latestTime = 0
	} else {
		latestTime, err = strconv.ParseInt(latest, 0, 64)
		if err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Invalid latestTime: " + err.Error(),
			})
			return
		}
	}

	res, sErr := moment.Service().Feed(context.Background(), &sdto.FeedMomentInput{
		UserID:     userId,
		LatestTime: latestTime,
	})
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	moments := make([]gin.H, len(res.Moments))
	for i := range res.Moments {
		momentID := res.Moments[i].MomentID
		moments[i] = gin.H{
			"id":        momentID,
			"createdAt": res.Moments[i].CreatedAt,
			"message":   res.Moments[i].Content,
			"authorInfo": gin.H{
				"ID":        res.AuthorInfoMap[momentID].UserID,
				"avatarUrl": res.AuthorInfoMap[momentID].AvatarURL,
				"name":      res.AuthorInfoMap[momentID].Username,
			},
		}

		// Check and set media url
		if res.Moments[i].ImageURL != nil {
			moments[i]["media_image"] = res.Moments[i].ImageURL
		}
		if res.Moments[i].VideoURL != nil {
			moments[i]["media_video"] = res.Moments[i].VideoURL
		}
		if GPXPath, ok := res.GPXRouteText[i]; ok {
			moments[i]["media"] = GPXPath
		}

		// Get moment liked person
		likeResp, sErr := moment.Service().GetLikesByMomentId(context.Background(), momentID)
		if sErr != nil {
			c.JSON(sErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  sErr.Error(),
			})
			return
		}
		moments[i]["personLikes"] = likeResp.PersonLikes

		// Get comment list
		commentResp, sErr := moment.Service().GetCommentListByMomentId(context.Background(), momentID)
		if sErr != nil {
			c.JSON(sErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  sErr.Error(),
			})
			return
		}
		moments[i]["comments"] = commentResp.CommentList

		isLiked, sErr := moment.Service().IsLiked(context.Background(), momentID, userId)
		if sErr != nil {
			c.JSON(sErr.Code(), dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  sErr.Error(),
			})
			return
		}
		moments[i]["isLiked"] = isLiked
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get feed moments successfully",
		Data: gin.H{
			"moments":  moments,
			"nextTime": res.NextTime,
		},
	})
}

func (m *MomentController) GetByUserID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Missing user ID in token",
		})
		return
	}

	res, sErr := moment.Service().GetByUserID(context.Background(), userID)
	if sErr != nil {
		c.JSON(sErr.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  sErr.Error(),
		})
		return
	}

	moments := make([]gin.H, len(res.Moments))
	for i, moment := range res.Moments {
		moments[i] = gin.H{
			"id":        moment.MomentID,
			"createdAt": moment.CreatedAt,
			"content":   moment.Content,
		}

		// Add media URLs if present
		if moment.ImageURL != nil {
			moments[i]["imageUrl"] = moment.ImageURL
		}
		if moment.VideoURL != nil {
			moments[i]["videoUrl"] = moment.VideoURL
		}
		if gpxText, ok := res.GPXRouteText[i]; ok {
			moments[i]["gpxRoute"] = gpxText
		}
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Get user moments successfully ",
		Data: gin.H{
			"moments": moments,
			"user": gin.H{
				"id":        res.User.UserID,
				"username":  res.User.Username,
				"avatarUrl": res.User.AvatarURL,
			},
		},
	})
}
