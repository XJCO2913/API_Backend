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
				StatusMsg:  "Authorization header is missing",
			})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CommonRes{
				StatusCode: -1,
				StatusMsg:  "Authorization header format must br Bearer {token}",
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
				StatusMsg:  "Error while get jwt claims",
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

		// register token payload into context
		var (
			userID         string
			isAdmin        bool
			isOrganiser    bool
			membershipType float64
		)

		if val, ok := claims["userID"]; ok {
			userID = val.(string)
		} else {
			userID = ""
		}
		
		if val, ok := claims["isAdmin"]; ok {
			isAdmin = val.(bool)
		} else {
			isAdmin = false
		}

		if val, ok := claims["isOrganiser"]; ok {
			isOrganiser = val.(bool)
		} else {
			isOrganiser = false
		}

		if val, ok := claims["membershipType"]; ok {
			membershipType = val.(float64)
		} else {
			membershipType = 0.0
		}

		ctx.Set("userID", userID)
		ctx.Set("isAdmin", isAdmin)
		ctx.Set("isOrganiser", isOrganiser)
		ctx.Set("membershipType", membershipType)

		ctx.Next()
	}
}
