package sdto

type CreateLikeInput struct {
	UserID   string `json:"userID" binding:"required"`
	MomentID string `json:"momentID" binding:"required"`
}
