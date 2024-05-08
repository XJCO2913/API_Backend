package sdto

type CreateLikeInput struct {
	UserID   string
	MomentID string
}

type DeleteLikeInput struct {
	UserID   string
	MomentID string
}
