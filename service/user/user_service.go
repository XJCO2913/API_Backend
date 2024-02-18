package user

import (
	"context"
	"errors"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct{}

var (
	userService UserService
)

func Service() *UserService {
	return &userService
}

func (u *UserService) Create(ctx context.Context, in *sdto.CreateUserInput) (*sdto.CreateUserOutput, error) {
	// check if user already exist or not
	user, err := dao.FindUserByUsername(ctx, in.Username)
	if err != gorm.ErrRecordNotFound || user != nil {
		return nil, errors.New("user already exist")
	}

	// generate uuid for userID
	uuid, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generate uuid: " + err.Error())
		return nil, err
	}
	newUserID := uuid.String()

	// parse birthday
	var birthdayEntity *time.Time = nil
	if !util.IsEmpty(in.Birthday) {
		birthday, err := time.Parse("2006-01-02", in.Birthday)
		if err != nil {
			zlog.Error("Error while parse birthday " + in.Birthday)
			return nil, err
		}

		birthdayEntity = &birthday
	}

	// encrypt password
	hashPwd, err := util.EncryptPassword(in.Password)
	if err != nil {
		zlog.Error("Error while encrypt password " + in.Password)
		return nil, err
	}

	// DB logic
	err = dao.CreateNewUser(ctx, &model.User{
		UserID:         newUserID,
		AvatarURL:      nil, // avatar not implement yet
		MembershipTime: time.Now().Unix(),
		Gender:         &in.Gender,
		Region:         in.Region,
		Tags:           nil,
		Birthday:       birthdayEntity,
		Username:       in.Username,
		Password:       hashPwd,
	})
	if err != nil {
		zlog.Error("Error while create new user: "+err.Error(), zap.String("username", in.Username))
		return nil, err
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
		return nil, errors.New("internal error")
	}
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		zlog.Error("error while sign jwt: " + err.Error())
		return nil, errors.New("internal error")
	}

	return &sdto.CreateUserOutput{
		UserID: newUserID,
		Token:  tokenStr,
	}, nil
}

func (u *UserService) Authenticate(ctx context.Context, in *sdto.AuthenticateInput) (*sdto.AuthenticateOutput, error) {
	// check whether the user exist or not
    user, err := dao.FindUserByUsername(ctx, in.Username)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        } else {
            zlog.Error("Error while finding user by username", zap.String("username", in.Username), zap.Error(err))
            return nil, errors.New("an error occurred while processing your request")
        }
    }

	// Verify the password
    if !util.VerifyPassword(user.Password, in.Password) {
		zlog.Info("Invalid login attempt", zap.String("username", in.Username))
        return nil, errors.New("invalid password")
    }

	// sign token
	claims := jwt.MapClaims{
        "userID":  user.UserID,
        "isAdmin": false,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.Get("jwt.secret")
	if util.IsEmpty(secret) {
        zlog.Error("jwt.secret is empty in config")
        return nil, errors.New("internal error")
    }
	tokenStr, err := token.SignedString([]byte(secret))
    if err != nil {
        zlog.Error("Error while signing jwt: " + err.Error())
        return nil, errors.New("internal error")
    }

	var gender int32
	if user.Gender != nil {
		gender = *user.Gender
	}

	var birthdayStr string
	if user.Birthday != nil {
		birthdayStr = user.Birthday.Format("2006-01-02")
	}

    return &sdto.AuthenticateOutput{
        UserID:   user.UserID,
		Token:    tokenStr,
        Gender:   gender,
        Birthday: birthdayStr,
        Region:   user.Region,
    }, nil
}
