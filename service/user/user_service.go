package user

import (
)

type UserService struct {}

var (
	userService UserService
)

func Service() *UserService {
	return &userService
}