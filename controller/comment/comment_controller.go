package comment

import (
	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/service/comment"
	"api.backend.xjco2913/service/sdto"
	"github.com/gin-gonic/gin"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

func (cc *CommentController) Create(c *gin.Context) {
	userID, userIDExists := c.Get("userID")
	if !userIDExists {
		c.JSON(403, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "User ID is required",
		})
		return
	}

	var req dto.CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  "Wrong params: " + err.Error(),
		})
		return
	}

	input := &sdto.CreateCommentInput{
		AuthorID: userID.(string),
		MomentID: req.MomentID,
		Content:  req.Content,
	}

	err := comment.Service().Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(err.Code(), dto.CommonRes{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(200, dto.CommonRes{
		StatusCode: 0,
		StatusMsg:  "Create comment successfully",
	})
}
