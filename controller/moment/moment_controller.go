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
			StatusMsg:  "missing userID in token",
		})
		return
	}

	content := c.PostForm("content")

	videoFileHeader, errVideo := c.FormFile("videoFile")
	imageFileHeader, errImage := c.FormFile("imageFile")
	GPXFileHeader, errGPX := c.FormFile("gpxFile")

	// 获取文件数量和可用文件头
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

	// 处理文件
	// 1. 在content为空的情况下只能传一个文件, 其他数量(包括0)的文件都是逻辑错误的
	// 2. 如果content不为空, 那么可以不传文件
	switch {
	case fileCnt == 0:
		if content == "" {
			// 文件数量为0
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "moment cannot be empty",
			})
			return
		}

		validFileHeader = "null"

	case fileCnt > 1:
		// 文件数量过多
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "only allow upload one type of file",
		})
		return
	}

	fileHeader := fileHeaderMap[validFileHeader]
	switch validFileHeader {
	case "gpxFile":
		// parser gpx file logic
		gpxFile, err := fileHeader.Open()
		if err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Fail to get gpx file: %s", err.Error()),
			})
			return
		}
		defer gpxFile.Close()

		gpxBuf := bytes.NewBuffer(nil)
		if _, err := io.Copy(gpxBuf, gpxFile); err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Fail copy image data: %s", err.Error()),
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
				StatusMsg:  fmt.Sprintf("Fail to get image file: %s", err.Error()),
			})
			return
		}
		defer imageFile.Close()

		imageBuf := bytes.NewBuffer(nil)
		if _, err := io.Copy(imageBuf, imageFile); err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Fail copy image data: %s", err.Error()),
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
				StatusMsg:  fmt.Sprintf("Fail to get video file: %s", err.Error()),
			})
			return
		}
		defer videoFile.Close()

		videoBuf := bytes.NewBuffer(nil)
		if _, err := io.Copy(videoBuf, videoFile); err != nil {
			c.JSON(400, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Fail copy video data: %s", err.Error()),
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
		// 只有content, 没有文件
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
			StatusMsg:  "missing user Id in token",
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
				StatusMsg:  "invalid latestTime: " + err.Error(),
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
		moments[i] = gin.H{
			"id":        res.Moments[i].MomentID,
			"createdAt": res.Moments[i].CreatedAt,
			"message":   res.Moments[i].Content,
		}

		// check and set media url
		if res.Moments[i].ImageURL != nil {
			moments[i]["media_image"] = res.Moments[i].ImageURL
		}
		if res.Moments[i].VideoURL != nil {
			moments[i]["media_video"] = res.Moments[i].VideoURL
		}
		if GPXPath, ok := res.GPXRouteText[i]; ok {
			moments[i]["media"] = GPXPath
		}
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
