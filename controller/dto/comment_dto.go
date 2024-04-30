package dto

type CreateCommentReq struct {
	MomentID string `json:"momentId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}
