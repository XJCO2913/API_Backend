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
	birthday, err := time.Parse("2006-01-02", in.Birthday)
	if err != nil {
		zlog.Error("Error while parse birthday " + in.Birthday)
		return nil, err
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
		Birthday:       &birthday,
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
		Token:  tokenStr, // jwt token not implement yet
	}, nil
}

func (u *UserService) Authenticate(ctx context.Context, in *sdto.AuthenticateInput) (*sdto.AuthenticateOutput, error) {
    user, err := dao.FindUserByUsername(ctx, in.Username)
    if err != nil {
        zlog.Error("Error while finding user by username", zap.String("username", in.Username), zap.Error(err))
		return nil, errors.New("an error occurred while processing your request")
    }

	// Find the user
    if user == nil {
        return nil, errors.New("user not found")
    }

	// Verify the password
    if !util.VerifyPassword(user.Password, in.Password) {
		zlog.Info("Invalid login attempt", zap.String("username", in.Username))
        return nil, errors.New("invalid password")
    }

    return &sdto.AuthenticateOutput{
        UserID:   user.UserID,
        Username: user.Username,
        Gender:   *user.Gender,
        Birthday: user.Birthday.Format("2006-01-02"),
        Region:   user.Region,
        Token:    "", // Token generation logic to be implemented
    }, nil
}
