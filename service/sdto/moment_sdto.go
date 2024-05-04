package sdto

import (
	"time"

	"api.backend.xjco2913/dao/model"
)

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

type CreateMomentGPXInput struct {
	UserID  string
	Content string
	GPXData []byte
}

type FeedMomentInput struct {
	UserID     string
	LatestTime int64
}

type FeedMomentOutput struct {
	Moments       []*model.Moment
	AuthorInfoMap map[string]*model.User
	NextTime      int64
	GPXRouteText  map[int][][]string
}

type MomentUser struct {
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
}

type MomentComment struct {
	Id        string     `json:"id"`
	Author    MomentUser `json:"author"`
	CreatedAt time.Time  `json:"createdAt"`
	Message   string     `json:"message"`
}

type GetLikesOutput struct {
	PersonLikes []MomentUser
}

type GetCommentListOutput struct {
	CommentList []MomentComment
}
