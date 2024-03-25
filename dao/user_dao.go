package dao

import (
	"context"
	"strings"

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

func DeleteUsersByID(ctx context.Context, userIDs string) ([]string, []string, error) {
	ids := strings.Split(userIDs, "|")
	var deletedIDs []string
	var notFoundIDs []string

	for _, id := range ids {
		_, err := GetUserByID(ctx, id)
		if err != nil {
			notFoundIDs = append(notFoundIDs, id)
			continue
		}

		u := query.Use(DB).User
		result, err := u.WithContext(ctx).Where(u.UserID.Eq(id)).Delete()
		if err != nil {
			return nil, nil, err
		}

		if result.RowsAffected > 0 {
			deletedIDs = append(deletedIDs, id)
		} else {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	return deletedIDs, notFoundIDs, nil
}

func UpdateUserByID(ctx context.Context, userID string, updates map[string]interface{}) error {
	_, err := GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	u := query.Use(DB).User

	_, err = u.WithContext(ctx).Where(u.UserID.Eq(userID)).Updates(updates)
	if err != nil {
		return err
	}

	return nil
}
