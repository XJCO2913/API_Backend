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

type CreateMomentVideoInput struct {
	UserID    string
	Content   string
	VideoData []byte
}
