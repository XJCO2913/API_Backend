package like

import (
	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/like"
	"api.backend.xjco2913/service/sdto"
	"github.com/gin-gonic/gin"
)

type LikeController struct{}

func NewLikeController() *LikeController {
	return &LikeController{}
}

func (lc *LikeController) Create(c *gin.Context) {
	userID, userIDExists := c.Get("userID")
	momentID := c.Query("momentID")

	if !userIDExists {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "User ID is required",
		})
		return
	}

	input := &sdto.CreateLikeInput{
		UserID:   userID.(string),
		MomentID: momentID,
	}

	err := like.Service().Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Like created successfully",
	})
}
