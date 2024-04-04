package sdto

type CreateMomentInput struct {
	UserID  string
	Content string
}

type CreateMomentImageInput struct {
	UserID    string
	Content   string
	ImageData []byte
}
