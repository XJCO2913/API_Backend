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
	u := query.Use(DB).User

	users, err := u.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	u := query.Use(DB).User

	user, err := u.WithContext(ctx).Where(u.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUserByID(ctx context.Context, userID string) error {
	_, err := GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	u := query.Use(DB).User
	_, err = u.WithContext(ctx).Where(u.UserID.Eq(userID)).Delete()
	if err != nil {
		return err
	}

	return nil
}
