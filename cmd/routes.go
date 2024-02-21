package main

import (
	"api.backend.xjco2913/controller/user"
	"api.backend.xjco2913/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	userController := user.NewUserController()

	// global middleware
	r.Use(cors.Default())

	api := r.Group("/api")
	{	
		// group middleware
		api.Use(middleware.VerifyToken())

		api.POST("/user/register", userController.SignUp)
		api.POST("/user/login", userController.Login)
		api.GET("/test", func(c *gin.Context) {
			userID := c.GetString("userID")
			isAdmin := c.GetBool("isAdmin")

			c.JSON(200, gin.H{
				"userID": userID,
				"isAdmin": isAdmin,
			})
		})
	}

	return r
}