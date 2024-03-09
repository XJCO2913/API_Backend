package dao

import (
	"context"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func CreateNewUser(ctx context.Context, newUser *model.User) error {
	err := query.Use(DB).WithContext(ctx).User.Create(newUser)
	if err != nil {
		return err
	}

	return nil
}

func FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	u := query.Use(DB).User

	user, err := u.WithContext(ctx).Where(u.Username.Eq(username)).First()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers(ctx context.Context) ([]*model.User, error) {
    var users []*model.User
    err := DB.WithContext(ctx).Model(&model.User{}).Find(&users).Error
    if err != nil {
        return nil, err
    }
    return users, nil
}
