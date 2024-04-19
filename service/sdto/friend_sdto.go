package sdto

import "api.backend.xjco2913/dao/model"

type FollowInput struct {
	FollowerId  string
	FollowingId string
}

type GetAllFollowerOutput struct {
	Followers []*model.User
}
