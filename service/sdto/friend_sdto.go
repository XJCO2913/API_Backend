package sdto

import "api.backend.xjco2913/dao/model"

type FollowInput struct {
	FollowerId  string
	FollowingId string
}

type Follower struct {
	UserID     string
	AvatarUrl  string
	Username   string
	Region     string
	IsFollowed bool
}

type GetAllFollowerOutput struct {
	Followers []*Follower
}

type GetAllFollowingOutput struct {
	Followings []*model.User
}

type GetAllFriendsOutput struct {
	Friends []*model.User
}

type FollowerCountOutput struct {
	Count int
}

type FollowingCountOutput struct {
	Count int
}
