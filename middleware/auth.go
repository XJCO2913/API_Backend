package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"api.backend.xjco2913/controller/dto"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/api/user/login" || ctx.Request.URL.Path == "/api/user/register" {
			// login and register no need to auth token
			ctx.Next()
			return
		}

		if ctx.Request.URL.Path == "/api/admin/login" {
			ctx.Next()
			return
		}

		// verify token
		authHeader := ctx.GetHeader("Authorization")
		if util.IsEmpty(authHeader) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "authorization header is missing",
			})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "authorization header format must br Bearer {token}",
			})
			return
		}

		// parse  and verify token
		tokenString := bearerToken[1]
		jwtSecret := config.Get("jwt.secret")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// validate sign algorithm
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CommonRes{
					StatusCode: -1,
					StatusMsg:  "Invalid token signature",
				})
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  fmt.Sprintf("Invalid token: %v", err.Error()),
			})
			return
		}

		// if token is valid, check expried
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "error while get jwt claims",
			})
			return
		}

		exp := claims["exp"].(float64)
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Token expired",
			})
			return
		}

		// set userId into context
		userID := claims["userID"]
		isAdmin := claims["isAdmin"]
		isOrganiser := claims["isOrganiser"]
		membershipType := claims["membershipType"]
		ctx.Set("userID", userID)
		ctx.Set("isAdmin", isAdmin)
		ctx.Set("isOrganiser", isOrganiser)
		ctx.Set("membershipType", membershipType)

		ctx.Next()
	}
}
