package main

import (
	"time"

	"api.backend.xjco2913/controller/user"
	"api.backend.xjco2913/middleware"
	"api.backend.xjco2913/util/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	userController := user.NewUserController()

	// global middleware
	r.Use(cors.Default())

	// demo admin login api, only for test
	r.GET("/api/getAdmin", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		if username != "deck" {
			c.JSON(400, gin.H{
				"username": username,
				"password": password,
				"msg":      "not a admin",
			})
			return
		} else if password != "20030416" {
			c.JSON(400, gin.H{
				"username": username,
				"password": password,
				"msg":      "wrong password",
			})
			return
		}

		claims := jwt.MapClaims{
			"userID":  "123123123",
			"isAdmin": true,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := config.Get("jwt.secret")
		tokenStr, _ := token.SignedString([]byte(secret))

		c.JSON(200, gin.H{
			"username": username,
			"password": password,
			"token":    tokenStr,
		})
	})

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
				"userID":  userID,
				"isAdmin": isAdmin,
			})
		})
	}

	return r
}
