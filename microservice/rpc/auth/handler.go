package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/redis"
	"api.backend.xjco2913/microservice/consts"
	auth "api.backend.xjco2913/microservice/kitex_gen/rpc/xjco2913/auth"
	"api.backend.xjco2913/microservice/kitex_gen/rpc/xjco2913/base"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/zlog"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	maxLoginAttempts = 5
	lockDuration     = 3 * time.Minute
)

// LoginServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// Login implements the LoginServiceImpl interface.
func (s *AuthServiceImpl) Login(ctx context.Context, req *auth.LoginReq) (*auth.LoginResp, error) {
	resp := auth.NewLoginResp()
	resp.BaseResp = base.NewBaseResp()
	resp.BaseResp.Data = make(map[string]string)

	// check existence of user
	user, err := dao.FindUserByUsername(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "user not found"
			resp.BaseResp.Code = -1
			resp.BaseResp.Msg = errMsg
			return resp, nil
		} else {
			zlog.Error("Error while finding user by username", zap.String("username", req.GetUsername()), zap.Error(err))
			return nil, consts.ErrInternal
		}
	}

	// check if the user is banned
	banKey := fmt.Sprintf("ban:%s", user.UserID)
	exists, err := redis.RDB().Exists(ctx, banKey).Result()
	if err == nil && exists != 0 {
		// user is banned
		errMsg := "account is banned"
		resp.BaseResp.Code = -1
		resp.BaseResp.Msg = errMsg
		return resp, nil
	}

	/*
		wrong password attempts logic
	*/
	attemptKey := fmt.Sprintf("WrongPwd:%s", req.Username)
	lockKey := fmt.Sprintf("lock:%s", req.Username)

	// check if user is locked
	lockedUntilStr, err := redis.RDB().Get(ctx, lockKey).Result()
	if err == nil {
		lockedUntil, err := time.Parse(time.RFC3339, lockedUntilStr)
		if err == nil && time.Now().Before(lockedUntil) {
			// account is still locked
			errMsg := fmt.Sprintf("account is locked until %v", lockedUntil)
			resp.BaseResp.Code = -1
			resp.BaseResp.Msg = errMsg
			resp.BaseResp.Data = map[string]string{
				"remaining_attempts": "0",
				"lock_expires":       lockedUntil.String(),
			}
			return resp, nil
		}
	}

	/*
		verify password
	*/
	if util.VerifyPassword(user.Password, req.GetPassword()) {
		redis.RDB().Del(ctx, attemptKey)
		redis.RDB().Del(ctx, lockKey)
	} else {
		// Increment login attempt and check for lock condition
		attempts, err := redis.RDB().Incr(ctx, attemptKey).Result()
		if err != nil {
			zlog.Error("error incrementing login attempts", zap.Error(err))
			return nil, consts.ErrInternal
		}
		redis.RDB().Expire(ctx, attemptKey, lockDuration)

		if attempts >= maxLoginAttempts {
			// Account should be locked, set the lock expiration
			lockExpiration := time.Now().Add(lockDuration)
			err := redis.RDB().Set(ctx, lockKey, lockExpiration.Format(time.RFC3339), lockDuration).Err()
			if err != nil {
				zlog.Error("Error while set lock key", zap.Error(err))
				return nil, consts.ErrInternal
			}

			// lock account
			errMsg := fmt.Sprintf("account is locked until %v", lockExpiration)
			resp.BaseResp.Code = consts.ErrExternalCode
			resp.BaseResp.Msg = errMsg
			resp.BaseResp.Data = map[string]string{
				"remaining_attempts": "0",
				"lock_expires":       lockExpiration.String(),
			}
			return resp, nil
		} else {
			errMsg := fmt.Sprintf("invalid password, %d attempts remaining", maxLoginAttempts-attempts)
			resp.BaseResp.Code = consts.ErrExternalCode
			resp.BaseResp.Msg = errMsg
			resp.BaseResp.Data = map[string]string{
				"remaining_attempts": strconv.Itoa(int(maxLoginAttempts - attempts)),
			}
			return resp, nil
		}
	}

	/*
		sign token
	*/
	// First try to get jwt cache from redis
	// Key format => jwt:username
	cacheTokenKey := fmt.Sprintf("jwt:%v", req.GetUsername())
	cachedToken, err := redis.RDB().Get(ctx, cacheTokenKey).Result()
	if err != nil && err != redis.KEY_NOT_FOUND {
		// error occur
		errMsg := "fail to get cached token"
		zlog.Error(errMsg, zap.Error(err), zap.String("cacheTokenKey", cacheTokenKey))
		return nil, consts.ErrInternal
	}

	// if exist cached token, return it immediately
	var tokenStr string
	if err != redis.KEY_NOT_FOUND {
		tokenStr = cachedToken
	} else {
		organiser, err := dao.GetOrganiserByID(ctx, user.UserID)
		isOrganiser := false
		if err != nil && organiser != nil {
			isOrganiser = true
		}

		// sign new token
		claims := jwt.MapClaims{
			"userID":         user.UserID,
			"isAdmin":        false,
			"exp":            time.Now().Add(24 * time.Hour).Unix(),
			"isOrganiser":    isOrganiser,
			"membershipType": user.MembershipType,
		}

		tokenStr, err = util.GenerateJWTToken(claims)
		if err != nil {
			zlog.Error("Error while generating jwt", zap.Error(err))
			return nil, consts.ErrInternal
		}

		// store new token into cache
		err = redis.RDB().Set(ctx, cacheTokenKey, tokenStr, 24*time.Hour).Err()
		if err != nil {
			zlog.Error("Fail to store token into cache", zap.Error(err))
			return nil, consts.ErrInternal
		}
	}

	var birthdayStr string
	if user.Birthday != nil {
		birthdayStr = user.Birthday.Format(time.RFC822Z)
	}

	resp.BaseResp.Code = 0
	resp.Token = tokenStr
	resp.Gender = user.Gender
	resp.Birthday = birthdayStr
	resp.Region = user.Region
	return resp, nil
}

// RefreshToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *auth.RefreshTokenReq) (resp *auth.RefreshTokenResp, err error) {
	// TODO: Your code here...
	return
}
