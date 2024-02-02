package user

import (
	"context"
	"errors"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
	"api.backend.xjco2913/service/user/sdto"
)

type UserService struct {}

var (
	userService UserService
)

func Service() *UserService {
	return &userService
}

func (s *UserService) CreateUser(ctx context.Context, in *sdto.CreateUserInput) (int, error) {
	if in == nil || in.Name == "" {
		return 0, errors.New("username cannot be empty")
	}

	newUser := model.User{
		Name: in.Name,
	}
	err := query.Use(dao.DB).User.WithContext(ctx).Create(&newUser)
	if err != nil {
		return 0, errors.New("internal database error")
	}

	return int(newUser.ID), nil
}