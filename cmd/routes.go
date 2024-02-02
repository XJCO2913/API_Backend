package main

import (
	"api.backend.xjco2913/controller/user"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	userController := user.NewUserController()

	api := r.Group("/api")
	api.POST("/user", userController.CreateUser)

	return r
}