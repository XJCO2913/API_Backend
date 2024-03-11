package user

import (
	"context"
	"errors"
	"fmt"
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
	// Check whether the user exist or not
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

	// Keys for tracking login attempts and locks
	attemptKey := fmt.Sprintf("WrongPwd:%s", in.Username)
	lockKey := fmt.Sprintf("lock:%s", in.Username)

	// Check if it is locked
	lockedUntilStr, err := redis.RDB().Get(ctx, lockKey).Result()
	if err == nil {
		lockedUntil, err := time.Parse(time.RFC3339, lockedUntilStr)
		if err == nil && time.Now().Before(lockedUntil) {
			// User is locked
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
		// Password correct, reset attempts and unlock
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
			// Return remaining attempts without lock expiration
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

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		zlog.Error("Error while signing jwt: " + err.Error())
		return nil, errorx.NewInternalErr()
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

func (s *UserService) GetAll(ctx context.Context) ([]*sdto.UserDetail, error) {
    users, err := dao.GetAllUsers(ctx)
    if err != nil {
        return nil, err
    }

    userDtos := make([]*sdto.UserDetail, len(users))
    for i, user := range users {
		var birthday string
        if user.Birthday != nil {
            birthday = user.Birthday.Format("2006-01-02")
        }

        userDtos[i] = &sdto.UserDetail{
            UserID:   user.UserID,
            Username: user.Username,
            Gender:   user.Gender,
            Birthday: birthday,
            Region:   user.Region,
        }
    }

    return userDtos, nil
}
