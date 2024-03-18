package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/redis"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	maxLoginAttempts = 5
	lockDuration     = 3 * time.Minute
)

type UserService struct{}

var (
	userService UserService
)

func Service() *UserService {
	return &userService
}

func (u *UserService) Create(ctx context.Context, in *sdto.CreateUserInput) (*sdto.CreateUserOutput, *errorx.ServiceErr) {
	// check if user already exist or not
	user, err := dao.FindUserByUsername(ctx, in.Username)
	if err != gorm.ErrRecordNotFound || user != nil {
		return nil, errorx.NewServicerErr(
			errorx.ErrExternal,
			"user already exist",
			nil,
		)
	}

	// generate uuid for userID
	uuid, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generate uuid: " + err.Error())
		return nil, errorx.NewInternalErr()
	}
	newUserID := uuid.String()

	// parse birthday
	var birthdayEntity *time.Time = nil
	if !util.IsEmpty(in.Birthday) {
		birthday, err := time.Parse("2006-01-02", in.Birthday)
		if err != nil {
			zlog.Error("Error while parse birthday " + in.Birthday)
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				"invalid birthday",
				nil,
			)
		}

		birthdayEntity = &birthday
	}

	// encrypt password
	hashPwd, err := util.EncryptPassword(in.Password)
	if err != nil {
		zlog.Error("Error while encrypt password " + in.Password)
		return nil, errorx.NewInternalErr()
	}

	// DB logic
	err = dao.CreateNewUser(ctx, &model.User{
		UserID:         newUserID,
		AvatarURL:      nil, // avatar not implement yet
		MembershipTime: time.Now().Unix(),
		Gender:         in.Gender,
		Region:         in.Region,
		Tags:           nil,
		Birthday:       birthdayEntity,
		Username:       in.Username,
		Password:       hashPwd,
	})
	if err != nil {
		zlog.Error("Error while create new user: "+err.Error(), zap.String("username", in.Username))
		return nil, errorx.NewInternalErr()
	}

	// sign token
	claims := jwt.MapClaims{
		"userID":  newUserID,
		"isAdmin": false,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.Get("jwt.secret")
	if util.IsEmpty(secret) {
		zlog.Error("jwt.secret is empty in config")
		return nil, errorx.NewInternalErr()
	}
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		zlog.Error("error while sign jwt: " + err.Error())
		return nil, errorx.NewInternalErr()
	}

	return &sdto.CreateUserOutput{
		UserID: newUserID,
		Token:  tokenStr,
	}, nil
}

func (u *UserService) Authenticate(ctx context.Context, in *sdto.AuthenticateInput) (*sdto.AuthenticateOutput, *errorx.ServiceErr) {
	user, err := dao.FindUserByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				"user not found",
				nil,
			)
		} else {
			zlog.Error("Error while finding user by username", zap.String("username", in.Username), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	// Check if the user is banned
	if u.IsBanned(ctx, user.UserID) {
		zlog.Warn("Attempted login by banned user", zap.String("userID", user.UserID))
		return nil, errorx.NewServicerErr(errorx.ErrExternal, "account is banned", nil)
	}

	attemptKey := fmt.Sprintf("WrongPwd:%s", in.Username)
	lockKey := fmt.Sprintf("lock:%s", in.Username)

	// Check if the user is locked
	lockedUntilStr, err := redis.RDB().Get(ctx, lockKey).Result()
	if err == nil {
		lockedUntil, err := time.Parse(time.RFC3339, lockedUntilStr)
		if err == nil && time.Now().Before(lockedUntil) {
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				fmt.Sprintf("account is locked until %v", lockedUntil),
				map[string]any{
					"remaining_attempts": 0,
					"lock_expires":       lockedUntil,
				},
			)
		}
	}

	if util.VerifyPassword(user.Password, in.Password) {
		redis.RDB().Del(ctx, attemptKey)
		redis.RDB().Del(ctx, lockKey)
	} else {
		// Increment login attempt and check for lock condition
		attempts, err := redis.RDB().Incr(ctx, attemptKey).Result()
		if err != nil {
			zlog.Error("Error incrementing login attempts", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
		redis.RDB().Expire(ctx, attemptKey, lockDuration)

		if attempts >= maxLoginAttempts {
			// Account should be locked, set the lock expiration
			lockExpiration := time.Now().Add(lockDuration)
			err := redis.RDB().Set(ctx, lockKey, lockExpiration.Format(time.RFC3339), lockDuration).Err()
			if err != nil {
				zlog.Error("Error while set lock key", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			zlog.Warn("Account locked due to too many failed login attempts", zap.String("username", in.Username))
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				fmt.Sprintf("account is locked until %v", lockExpiration),
				map[string]any{
					"remaining_attempts": 0,
					"lock_expires":       lockExpiration,
				},
			)
		} else {
			zlog.Info("Invalid login attempt", zap.String("username", in.Username))
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				fmt.Sprintf("invalid password, %d attempts remaining", maxLoginAttempts-attempts),
				map[string]any{
					"remaining_attempts": maxLoginAttempts - attempts,
				},
			)
		}
	}

	// first try to get jwt cache from redis
	// key format => jwt:username
	cacheTokenKey := fmt.Sprintf("jwt:%v", in.Username)
	cachedToken, err := redis.RDB().Get(ctx, cacheTokenKey).Result()
	if err != nil && err != redis.KEY_NOT_FOUND {
		// error occur
		zlog.Error("fail to get cached token", zap.Error(err), zap.String("cachedTokenKey", cacheTokenKey))
		return nil, errorx.NewInternalErr()
	}

	// if exist cached token, return it immediately
	var tokenStr string
	if err != redis.KEY_NOT_FOUND {
		tokenStr = cachedToken
	} else {
		// if not exist cached token, generate a new token
		// Sign token
		claims := jwt.MapClaims{
			"userID":  user.UserID,
			"isAdmin": false,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := config.Get("jwt.secret")

		if util.IsEmpty(secret) {
			zlog.Error("jwt.secret is empty in config")
			return nil, errorx.NewInternalErr()
		}

		tokenStr, err = token.SignedString([]byte(secret))
		if err != nil {
			zlog.Error("Error while signing jwt: " + err.Error())
			return nil, errorx.NewInternalErr()
		}

		// store the token into cache
		err = redis.RDB().Set(ctx, cacheTokenKey, tokenStr, 24*time.Hour).Err()
		if err != nil {
			zlog.Error("fail to store token into cache", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	var birthdayStr string
	if user.Birthday != nil {
		birthdayStr = user.Birthday.Format("2006-01-02")
	}

	return &sdto.AuthenticateOutput{
		UserID:   user.UserID,
		Token:    tokenStr,
		Gender:   user.Gender,
		Birthday: birthdayStr,
		Region:   user.Region,
	}, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*sdto.GetAllOutput, *errorx.ServiceErr) {
	users, err := dao.GetAllUsers(ctx)
	if err != nil {
		zlog.Error("Failed to retrieve all users", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	userDtos := make([]*sdto.GetAllOutput, len(users))
	for i, user := range users {
		var birthday string
		if user.Birthday != nil {
			birthday = user.Birthday.Format("2006-01-02")
		}

		userDtos[i] = &sdto.GetAllOutput{
			UserID:         user.UserID,
			Username:       user.Username,
			Gender:         user.Gender,
			Birthday:       birthday,
			Region:         user.Region,
			MembershipTime: user.MembershipTime,
		}
	}

	return userDtos, nil
}

func (s *UserService) GetByID(ctx context.Context, userID string) (*sdto.GetByIDOutput, *errorx.ServiceErr) {
	user, err := dao.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("User not found", zap.String("userID", userID))
			return nil, errorx.NewServicerErr(errorx.ErrExternal, "User not found", nil)
		} else {
			zlog.Error("Failed to retrieve user by ID", zap.String("userID", userID), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	var birthday string
	if user.Birthday != nil {
		birthday = user.Birthday.Format("2006-01-02")
	}

	userDto := &sdto.GetByIDOutput{
		UserID:         user.UserID,
		Username:       user.Username,
		Gender:         user.Gender,
		Birthday:       birthday,
		Region:         user.Region,
		MembershipTime: user.MembershipTime,
	}

	return userDto, nil
}

func (s *UserService) DeleteByID(ctx context.Context, userIDs string) *errorx.ServiceErr {
	ids := strings.Split(userIDs, "|")
	deletedIDs, notFoundIDs, err := dao.DeleteUsersByID(ctx, userIDs)

	if err != nil {
		zlog.Error("Failed to delete users", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// All specified users were not found
	if len(notFoundIDs) == len(ids) {
		zlog.Error("All specified users not found", zap.Strings("not_found_ids", notFoundIDs))
		return errorx.NewServicerErr(errorx.ErrExternal, "All specified users not found", map[string]any{"not_found_ids": notFoundIDs})
	}

	zlog.Info("Specified users deleted", zap.Strings("deleted_user_ids", deletedIDs))
	// Part of specified users were not found
	if len(notFoundIDs) > 0 {
		zlog.Warn("Some specified users not found", zap.Strings("not_found_ids", notFoundIDs))
	}

	return nil
}

func (s *UserService) BanByID(ctx context.Context, userID string) *errorx.ServiceErr {
	_, err := dao.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("User not found", zap.String("userID", userID))
			return errorx.NewServicerErr(errorx.ErrExternal, "User not found", nil)
		} else {
			zlog.Error("Error checking user existence", zap.String("userID", userID), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	if s.IsBanned(ctx, userID) {
		zlog.Warn("User already banned", zap.String("userID", userID))
		return errorx.NewServicerErr(errorx.ErrExternal, "User already banned", nil)
	}

	banKey := fmt.Sprintf("ban:%s", userID)
	// Disable never expire
	err = redis.RDB().Set(ctx, banKey, "banned", 0).Err()
	if err != nil {
		zlog.Error("Failed to ban user", zap.String("userID", userID), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (s *UserService) UnbanByID(ctx context.Context, userID string) *errorx.ServiceErr {
	_, err := dao.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("User not found", zap.String("userID", userID))
			return errorx.NewServicerErr(errorx.ErrExternal, "User not found", nil)
		} else {
			zlog.Error("Error checking user existence", zap.String("userID", userID), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	if !s.IsBanned(ctx, userID) {
		zlog.Warn("User not banned", zap.String("userID", userID))
		return errorx.NewServicerErr(errorx.ErrExternal, "User not banned", nil)
	}

	banKey := fmt.Sprintf("ban:%s", userID)
	err = redis.RDB().Del(ctx, banKey).Err()
	if err != nil {
		zlog.Error("Failed to unban user", zap.String("userID", userID), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (s *UserService) IsBanned(ctx context.Context, userID string) bool {
	banKey := fmt.Sprintf("ban:%s", userID)
	exists, err := redis.RDB().Exists(ctx, banKey).Result()
	if err != nil || exists == 0 {
		return false
	}
	return true
}
