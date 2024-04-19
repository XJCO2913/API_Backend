package sdto

import "api.backend.xjco2913/dao/model"

type FollowInput struct {
	FollowerId  string
	FollowingId string
}

type GetAllFollowerOutput struct {
	Followers []*model.User
}

type GetAllFollowingOutput struct {
	Followings []*model.User
}

type GetAllFriendsOutput struct {
	Friends []*model.User
}
