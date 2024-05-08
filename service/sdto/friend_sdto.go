package sdto

import "api.backend.xjco2913/dao/model"

type FollowInput struct {
	FollowerId  string
	FollowingId string
}

type Follower struct {
	UserID     string `json:"userId"`
	AvatarUrl  string `json:"avatarUrl"`
	Username   string `json:"username"`
	Region     string `json:"region"`
	IsFollowed bool   `json:"isFollowed"`
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
