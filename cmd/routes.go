package main

import (
	"time"

	"api.backend.xjco2913/controller/activity"
	"api.backend.xjco2913/controller/admin"
	"api.backend.xjco2913/controller/moment"
	"api.backend.xjco2913/controller/user"
	"api.backend.xjco2913/middleware"
	"api.backend.xjco2913/util/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	userController := user.NewUserController()
	activityController := activity.NewActivityController()
	adminController := admin.NewAdminController()
	momentController := moment.NewMomentController()

	// global middleware
	// prometheus
	r.Use(middleware.PrometheusRequests())
	r.Use(middleware.PrometheusDuration())
	r.Use(middleware.PrometheusResErr())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
			"userID":      "123123123",
			"isAdmin":     true,
			"isOrganiser": true,
			"exp":         time.Now().Add(24 * time.Hour).Unix(),
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
		api.GET("/user", userController.GetByID)
		api.GET("/users", userController.GetAll)
		api.DELETE("/user", userController.DeleteByID)
		api.POST("/user/ban", userController.BanByID)
		api.POST("/user/unban", userController.UnbanByID)
		api.GET("/user/status", userController.IsBanned)
		api.GET("/user/statuses", userController.GetAllStatus)
		api.PATCH("/user", userController.UpdateByID)
		api.POST("/user/subscribe", userController.Subscribe)
		api.POST("/user/cancel", userController.CancelByID)
		api.POST("/activity/create", activityController.Create)
		api.GET("/activities", activityController.GetAll)
		api.GET("/test", func(c *gin.Context) {
			userID := c.GetString("userID")
			isAdmin := c.GetBool("isAdmin")

			c.JSON(200, gin.H{
				"userID":  userID,
				"isAdmin": isAdmin,
			})
		})

		api.POST("/user/avatar", userController.UploadAvatar)

		// admin
		admin := api.Group("/admin")
		{
			admin.POST("/login", adminController.Login)
		}

		// moments
		moment := api.Group("/moment")
		{
			moment.POST("/create", momentController.Create)
		}
	}

	return r
}
