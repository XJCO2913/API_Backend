package sdto

import "api.backend.xjco2913/dao/model"

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

type FeedMomentInput struct {
	UserID      string
	LatestTime int64
}

type FeedMomentOutput struct {
	Moments []*model.Moment
	NextTime int64
}
