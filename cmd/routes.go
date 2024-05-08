package main

import (
	"context"
	"time"

	"api.backend.xjco2913/controller/activity"
	"api.backend.xjco2913/controller/admin"
	"api.backend.xjco2913/controller/comment"
	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/controller/friend"
	"api.backend.xjco2913/controller/like"
	"api.backend.xjco2913/controller/moment"
	"api.backend.xjco2913/controller/notify"
	"api.backend.xjco2913/controller/organiser"
	"api.backend.xjco2913/controller/user"
	"api.backend.xjco2913/controller/ws"
	"api.backend.xjco2913/middleware"
	userService "api.backend.xjco2913/service/user"
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
	friendController := friend.NewFriendController()
	websocketController := ws.NewWebsocketController()
	likeController := like.NewLikeController()
	commentController := comment.NewCommentController()
	organiserController := organiser.NewOrganiserController()
	notifyController := notify.NewNotifyController()

	// Global middleware
	// Prometheus
	r.Use(middleware.PrometheusRequests())
	r.Use(middleware.PrometheusDuration())
	r.Use(middleware.PrometheusResErr())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// WebSocket
	r.GET("/ws", gin.WrapF(websocketController.HandleConnections))

	// Demo admin login api, only for test
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
			"userID":         "123123123",
			"isAdmin":        true,
			"isOrganiser":    true,
			"membershipType": 2,
			"exp":            time.Now().Add(24 * time.Hour).Unix(),
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
		// Group middleware
		api.Use(middleware.VerifyToken())

		api.POST("/user/refresh", userController.RefreshToken)
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
		api.GET("/test", func(c *gin.Context) {
			userID := c.GetString("userID")
			isAdmin := c.GetBool("isAdmin")

			c.JSON(200, gin.H{
				"userID":  userID,
				"isAdmin": isAdmin,
			})
		})

		api.POST("/user/avatar", userController.UploadAvatar)

		// Admin
		admin := api.Group("/admin")
		{
			admin.POST("/login", adminController.Login)
		}

		// Moments
		moment := api.Group("/moment")
		{
			moment.POST("/create", momentController.Create)
			moment.GET("/feed", momentController.Feed)
			moment.POST("/like", likeController.Create)
			moment.DELETE("/unlike", likeController.DeleteByIDs)
			moment.POST("/comment", commentController.Create)
		}

		// Activity
		activity := api.Group("/activity")
		{
			activity.POST("/create", activityController.Create)
			activity.GET("", activityController.GetByID)
			activity.GET("/all", activityController.GetAll)
			activity.GET("/feed", activityController.Feed)
			activity.DELETE("", activityController.DeleteByID)
			activity.POST("/signup", activityController.SignUpByActivityID)
			activity.GET("/user", activityController.GetByUserID)
			activity.GET("/creator", activityController.GetByCreatorID)
			activity.GET("/profit", activityController.GetProfitWithOption)
			activity.GET("/tags", activityController.TagsInfo)
			activity.GET("/counts", activityController.Counts)
			activity.POST("/route", activityController.UploadRoute)
			activity.GET("/route", activityController.GetRouteByIDs)
		}

		// Friend
		friend := api.Group("/friend")
		{
			friend.POST("/follow", friendController.Follow)
			friend.GET("/follower", friendController.GetAllFollower)
			friend.GET("/following", friendController.GetAllFollowing)
			friend.GET("/", friendController.GetAll)
		}

		organiser := api.Group("/org")
		{
			organiser.GET("", organiserController.GetAll)
			organiser.POST("/agree", organiserController.Agree)
			organiser.POST("/refuse", organiserController.Refuse)
			organiser.POST("/apply", organiserController.Apply)
		}

		notify := api.Group("/notify")
		{
			notify.GET("/pull", notifyController.Pull)
			notify.POST("/route", notifyController.ShareRoute)
			notify.POST("/org", notifyController.OrgResult)
		}

		// mock data api
		mock := api.Group("/mock")
		{
			mock.GET("/shareList", func(c *gin.Context) {
				resp, sErr := userService.Service().MockUserList(context.Background())
				if sErr != nil {
					c.JSON(sErr.Code(), dto.CommonRes{
						StatusCode: -1,
						StatusMsg:  sErr.Error(),
					})
					return
				}

				c.JSON(200, dto.CommonRes{
					StatusCode: 0,
					StatusMsg:  "Get user list successfully",
					Data:       resp.MockUserList,
				})
			})
		}
	}

	return r
}
