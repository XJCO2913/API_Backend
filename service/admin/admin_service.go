package admin

import (
	"context"
	"errors"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminService struct{}

var (
	localAdminService AdminService
)

func Service() *AdminService {
	return &localAdminService
}

func (a *AdminService) Authenticate(ctx context.Context, in *sdto.AdminAuthenticateInput) (*sdto.AdminAuthenticateOuput, *errorx.ServiceErr) {
	admin, err := dao.FindAdminByName(ctx, in.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				"admin not found",
				nil,
			)
		} else {
			zlog.Error("Error while finding user by username", zap.String("admin-name", in.Name), zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	// verify password
	if in.Password != admin.Password {
		return nil, errorx.NewServicerErr(errorx.ErrExternal, "wrong password", nil)
	}

	// sign token
	claims := jwt.MapClaims{
		"userID": admin.ID,
		"isAdmin": true,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.Get("jwt.secret")

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		zlog.Error("Error while signing jwt: " + err.Error())
		return nil, errorx.NewInternalErr()
	}

	return &sdto.AdminAuthenticateOuput{
		Token: tokenStr,
		AdminId: admin.ID,
		Name: admin.Username,
	}, nil
} 