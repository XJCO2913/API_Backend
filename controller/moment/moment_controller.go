package moment

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

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
	// 只能传一个文件, 其他数量的文件都是逻辑错误的
	switch {
	case fileCnt == 0:
		// 文件数量为0
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "not found any file",
		})
		return

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
				StatusMsg: sErr.Error(),
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
	}
}
