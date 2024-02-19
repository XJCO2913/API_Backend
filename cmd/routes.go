package main

import (
	"api.backend.xjco2913/controller/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	userController := user.NewUserController()

	// global middleware
	r.Use(cors.Default())

	api := r.Group("/api")
	{
		api.POST("/user/register", userController.SignUp)
		api.POST("/user/login", userController.Login)
	}

	return r
}