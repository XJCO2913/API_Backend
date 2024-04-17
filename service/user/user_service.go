package user

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/redis"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util"
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

func (u *UserService) Create(ctx context.Context, in *sdto.CreateUserInput) *errorx.ServiceErr {
	// Check if user already exist or not
	user, err := dao.FindUserByUsername(ctx, in.Username)
	if err != gorm.ErrRecordNotFound || user != nil {
		return errorx.NewServicerErr(
			errorx.ErrExternal,
			"User already exist",
			nil,
		)
	}

	// Generate uuid for userID
	uuid, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generate uuid: " + err.Error())
		return errorx.NewInternalErr()
	}
	newUserID := uuid.String()

	// Parse birthday
	var birthdayEntity *time.Time = nil
	if !util.IsEmpty(in.Birthday) {
		birthday, err := time.Parse(time.RFC822Z, in.Birthday)
		if err != nil {
			zlog.Error("Error while parse birthday " + in.Birthday)
			return errorx.NewServicerErr(
				errorx.ErrExternal,
				"Invalid birthday format",
				nil,
			)
		}

		birthdayEntity = &birthday
	}

	// encrypt password
	hashPwd, err := util.EncryptPassword(in.Password)
	if err != nil {
		zlog.Error("Error while encrypt password " + in.Password)
		return errorx.NewInternalErr()
	}

	// DB logic to create new user
	err = dao.CreateNewUser(ctx, &model.User{
		UserID:         newUserID,
		AvatarURL:      nil,
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
		return errorx.NewInternalErr()
	}

	return nil
}

func (u *UserService) Authenticate(ctx context.Context, in *sdto.AuthenticateInput) (*sdto.AuthenticateOutput, *errorx.ServiceErr) {
	user, err := dao.FindUserByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				"User not found",
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
				fmt.Sprintf("Account is locked until %v", lockedUntil),
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
				fmt.Sprintf("Account is locked until %v", lockExpiration),
				map[string]any{
					"remaining_attempts": 0,
					"lock_expires":       lockExpiration,
				},
			)
		} else {
			zlog.Info("Invalid login attempt", zap.String("username", in.Username))
			return nil, errorx.NewServicerErr(
				errorx.ErrExternal,
				fmt.Sprintf("Invalid password, %d attempts remaining", maxLoginAttempts-attempts),
				map[string]any{
					"remaining_attempts": maxLoginAttempts - attempts,
				},
			)
		}
	}

	// First try to get jwt cache from redis
	// Key format => jwt:username
	cacheTokenKey := fmt.Sprintf("jwt:%v", in.Username)
	cachedToken, err := redis.RDB().Get(ctx, cacheTokenKey).Result()
	if err != nil && err != redis.KEY_NOT_FOUND {
		// error occur
		zlog.Error("Fail to get cached token", zap.Error(err), zap.String("cachedTokenKey", cacheTokenKey))
		return nil, errorx.NewInternalErr()
	}

	// If exist cached token, return it immediately
	var tokenStr string
	if err != redis.KEY_NOT_FOUND {
		tokenStr = cachedToken
	} else {
		organiser, err := dao.GetOrganiserByID(ctx, user.UserID)
		isOrganiser := false
		if err == nil && organiser != nil {
			isOrganiser = true
		}

		// If not exist cached token, prepare claims for new token
		claims := jwt.MapClaims{
			"userID":         user.UserID,
			"isAdmin":        false,
			"exp":            time.Now().Add(24 * time.Hour).Unix(),
			"isOrganiser":    isOrganiser,
			"membershipType": user.MembershipType,
		}

		// Generate new token
		tokenStr, err = util.GenerateJWTToken(claims)
		if err != nil {
			zlog.Error("Error while generating jwt", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		// Store the token into cache
		err = redis.RDB().Set(ctx, cacheTokenKey, tokenStr, 24*time.Hour).Err()
		if err != nil {
			zlog.Error("Fail to store token into cache", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}
	}

	var birthdayStr string
	if user.Birthday != nil {
		birthdayStr = user.Birthday.Format(time.RFC822Z)
	}

	return &sdto.AuthenticateOutput{
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

	organisers, err := dao.GetAllOrganisers(ctx)
	if err != nil {
		zlog.Error("Failed to retrieve all organisers", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// Check if a user is an organizer
	organiserMap := make(map[string]bool)
	for _, organiser := range organisers {
		organiserMap[organiser.UserID] = true
	}

	userDtos := make([]*sdto.GetAllOutput, len(users))
	for i, user := range users {
		var birthday string
		if user.Birthday != nil {
			birthday = user.Birthday.Format(time.RFC822Z)
		}

		// get avatar url from minio
		avatarURL := ""
		if user.AvatarURL != nil && !util.IsEmpty(*user.AvatarURL) {
			avatarURL, err = minio.GetUserAvatarUrl(ctx, *user.AvatarURL)
			if err != nil {
				zlog.Error("Error while get user avatar URL", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		organiserID := ""
		if _, exists := organiserMap[user.UserID]; exists {
			organiserID = user.UserID
		}

		userDtos[i] = &sdto.GetAllOutput{
			UserID:         user.UserID,
			Username:       user.Username,
			Gender:         user.Gender,
			Birthday:       birthday,
			Region:         user.Region,
			MembershipTime: user.MembershipTime,
			AvatarURL:      avatarURL,
			OrganiserID:    organiserID,
			MembershipType: user.MembershipType,
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
		birthday = user.Birthday.Format(time.RFC822Z)
	}

	// get avatar url from minio
	avatarURL := ""
	if user.AvatarURL != nil || !util.IsEmpty(user.AvatarURL) {
		avatarURL, err = minio.GetUserAvatarUrl(ctx, *user.AvatarURL)
		if err != nil {
			return nil, errorx.NewInternalErr()
		}
	}

	organiserID := ""
	organiser, err := dao.GetOrganiserByID(ctx, userID)
	if err == nil && organiser != nil {
		organiserID = organiser.UserID
	}

	userDto := &sdto.GetByIDOutput{
		UserID:         user.UserID,
		Username:       user.Username,
		Gender:         user.Gender,
		Birthday:       birthday,
		Region:         user.Region,
		MembershipTime: user.MembershipTime,
		AvatarURL:      avatarURL,
		OrganiserID:    organiserID,
		MembershipType: user.MembershipType,
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
		zlog.Warn("All specified users not found", zap.Strings("not_found_ids", notFoundIDs))
		return errorx.NewServicerErr(errorx.ErrExternal, "All specified users not found", map[string]any{"not_found_ids": notFoundIDs})
	}

	zlog.Info("Specified users deleted", zap.Strings("deleted_user_ids", deletedIDs))
	// Part of specified users were not found
	if len(notFoundIDs) > 0 {
		zlog.Warn("Some specified users not found", zap.Strings("not_found_ids", notFoundIDs))
	}

	return nil
}

func (s *UserService) BanByID(ctx context.Context, userIDs string) *errorx.ServiceErr {
	ids := strings.Split(userIDs, "|")
	var bannedIDs []string
	var notFoundIDs []string
	var alreadyBannedIDs []string

	for _, id := range ids {
		_, err := dao.GetUserByID(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundIDs = append(notFoundIDs, id)
				continue
			} else {
				zlog.Error("Failed to retrieve user by ID", zap.String("userID", id), zap.Error(err))
				return errorx.NewInternalErr()
			}
		}

		if s.IsBanned(ctx, id) {
			alreadyBannedIDs = append(alreadyBannedIDs, id)
			continue
		}

		banKey := fmt.Sprintf("ban:%s", id)
		// Disable never expire
		err = redis.RDB().Set(ctx, banKey, "banned", 0).Err()
		if err != nil {
			zlog.Error("Failed to ban user", zap.String("userID", id), zap.Error(err))
			return errorx.NewInternalErr()
		}
		bannedIDs = append(bannedIDs, id)
	}

	// All specified users were not found
	if len(notFoundIDs) == len(ids) {
		zlog.Warn("All specified users not found", zap.Strings("not_found_ids", notFoundIDs))
		return errorx.NewServicerErr(errorx.ErrExternal, "All specified users not found", map[string]any{"not_found_ids": notFoundIDs})
	}

	// All specified users were already banned
	if len(alreadyBannedIDs) == len(ids) {
		zlog.Warn("All specified users already banned", zap.Strings("already_banned_ids", alreadyBannedIDs))
		return errorx.NewServicerErr(errorx.ErrExternal, "All specified users already banned", map[string]any{"already_banned_ids": alreadyBannedIDs})
	}

	zlog.Info("Specified users banned", zap.Strings("banned_user_ids", bannedIDs))
	// Part of specified users were not found
	if len(notFoundIDs) > 0 {
		zlog.Warn("Some specified users not found", zap.Strings("not_found_ids", notFoundIDs))
	}

	// Part of specified users were already banned
	if len(alreadyBannedIDs) > 0 {
		zlog.Warn("Some specified users already banned", zap.Strings("already_banned_ids", alreadyBannedIDs))
	}

	return nil
}

func (s *UserService) UnbanByID(ctx context.Context, userIDs string) *errorx.ServiceErr {
	ids := strings.Split(userIDs, "|")
	var unbannedIDs []string
	var notFoundIDs []string
	var notBannedIDs []string

	for _, id := range ids {
		_, err := dao.GetUserByID(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundIDs = append(notFoundIDs, id)
				continue
			} else {
				zlog.Error("Failed to retrieve user by ID", zap.String("userID", id), zap.Error(err))
				return errorx.NewInternalErr()
			}
		}

		if !s.IsBanned(ctx, id) {
			notBannedIDs = append(notBannedIDs, id)
			continue
		}

		banKey := fmt.Sprintf("ban:%s", id)
		err = redis.RDB().Del(ctx, banKey).Err()
		if err != nil {
			zlog.Error("Failed to unban user", zap.String("userID", id), zap.Error(err))
			return errorx.NewInternalErr()
		}
		unbannedIDs = append(unbannedIDs, id)
	}

	// All specified users were not found
	if len(notFoundIDs) == len(ids) {
		zlog.Warn("All specified users not found", zap.Strings("not_found_ids", notFoundIDs))
		return errorx.NewServicerErr(errorx.ErrExternal, "All specified users not found", map[string]any{"not_found_ids": notFoundIDs})
	}

	// All specified users were not banned
	if len(notBannedIDs) == len(ids) {
		zlog.Warn("All specified users were not banned", zap.Strings("not_banned_ids", notBannedIDs))
		return errorx.NewServicerErr(errorx.ErrExternal, "All specified users were not banned", map[string]any{"not_banned_ids": notBannedIDs})
	}

	zlog.Info("Specified users unbanned", zap.Strings("unbanned_user_ids", unbannedIDs))
	// Part of specified users were not found
	if len(notFoundIDs) > 0 {
		zlog.Warn("Some specified users not found", zap.Strings("not_found_ids", notFoundIDs))
	}

	// Part of specified users were not banned
	if len(notBannedIDs) > 0 {
		zlog.Warn("Some specified users were not banned", zap.Strings("not_banned_ids", notBannedIDs))
	}

	return nil
}

func (s *UserService) IsBanned(ctx context.Context, userID string) bool {
	_, err := dao.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("User not found", zap.String("userID", userID))
		} else {
			zlog.Error("Failed to retrieve user by ID", zap.String("userID", userID), zap.Error(err))
		}
		return false
	}

	banKey := fmt.Sprintf("ban:%s", userID)
	exists, err := redis.RDB().Exists(ctx, banKey).Result()
	if err != nil || exists == 0 {
		return false
	}
	return true
}

func (s *UserService) GetAllStatus(ctx context.Context) ([]*sdto.GetAllStatusOutput, *errorx.ServiceErr) {
	users, err := dao.GetAllUsers(ctx)
	if err != nil {
		return nil, errorx.NewInternalErr()
	}

	var statusList []*sdto.GetAllStatusOutput
	for _, user := range users {
		isBanned := s.IsBanned(ctx, user.UserID)
		statusList = append(statusList, &sdto.GetAllStatusOutput{
			UserID:   user.UserID,
			IsBanned: isBanned,
		})
	}

	return statusList, nil
}

func (s *UserService) UpdateByID(ctx context.Context, userID string, input sdto.UpdateUserInput) *errorx.ServiceErr {
	updates := make(map[string]interface{})

	addUpdate := func(field string, value interface{}) {
		// Use reflection to check if the value is a pointer
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			updates[field] = v.Elem().Interface()
		} else if v.Kind() != reflect.Ptr {
			updates[field] = value
		}
	}

	addUpdate("username", input.Username)
	addUpdate("gender", input.Gender)
	addUpdate("region", input.Region)
	if input.Birthday != nil {
		if *input.Birthday == "" {
			addUpdate("birthday", nil)
		} else {
			_, err := time.Parse(time.RFC822Z, *input.Birthday)
			if err != nil {
				zlog.Error("Error while parsing birthday", zap.String("birthday", *input.Birthday), zap.Error(err))
				return errorx.NewServicerErr(errorx.ErrExternal, "Invalid birthday format", nil)
			}
			addUpdate("birthday", *input.Birthday)
		}
	}

	if len(updates) == 0 {
		zlog.Warn("All update fields were not provided")
		return errorx.NewServicerErr(errorx.ErrExternal, "All update fields were not provided", nil)
	}

	err := dao.UpdateUserByID(ctx, userID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("User not found", zap.String("userID", userID))
			return errorx.NewServicerErr(errorx.ErrExternal, "User not found", nil)
		} else {
			zlog.Error("Failed to update user", zap.String("userID", userID), zap.Any("updates", updates), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	return nil
}

func (s *UserService) Subscribe(ctx context.Context, userID string, membershipType int) *errorx.ServiceErr {
	user, err := dao.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("User not found", zap.String("userID", userID))
			return errorx.NewServicerErr(errorx.ErrExternal, "User not found", nil)
		} else {
			zlog.Error("Failed to retrieve user by ID", zap.String("userID", userID), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	if user.MembershipType != 0 {
		zlog.Warn("User has already subscribed", zap.String("userID", userID))
		return errorx.NewServicerErr(errorx.ErrExternal, "User has already subscribed", nil)
	}

	newExpiration := time.Now().Unix() + 30*24*60*60

	updates := map[string]interface{}{
		"membershipTime": newExpiration,
		"membershipType": membershipType,
	}

	err = dao.UpdateUserByID(ctx, userID, updates)
	if err != nil {
		zlog.Error("Failed to update user", zap.String("userID", userID), zap.Any("updates", updates), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (s *UserService) CancelByID(ctx context.Context, userID string) *errorx.ServiceErr {
	user, err := dao.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warn("User not found", zap.String("userID", userID))
			return errorx.NewServicerErr(errorx.ErrExternal, "User not found", nil)
		} else {
			zlog.Error("Failed to retrieve user by ID", zap.String("userID", userID), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	if user.MembershipType == 0 {
		zlog.Error("User has not subscribed", zap.String("userID", userID))
		return errorx.NewServicerErr(errorx.ErrExternal, "User has not subscribed", nil)
	}

	// No-reason refund (7 days after subscription start date)
	cancellationDeadline := user.MembershipTime - 23*24*60*60

	// Check if the current time is before the cancellation deadline
	if time.Now().Unix() > cancellationDeadline {
		zlog.Error("Cancellation period has expired", zap.String("userID", userID))
		return errorx.NewServicerErr(errorx.ErrExternal, "Cancellation period has expired", nil)
	}

	updates := map[string]interface{}{
		"membershipType": 0,
		"membershipTime": 0,
	}

	err = dao.UpdateUserByID(ctx, userID, updates)
	if err != nil {
		zlog.Error("Failed to update user", zap.String("userID", userID), zap.Any("updates", updates), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (s *UserService) UploadAvatar(ctx context.Context, in sdto.UploadAvatarInput) *errorx.ServiceErr {
	_, err := dao.GetUserByID(ctx, in.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.NewServicerErr(
				errorx.ErrExternal,
				"User not found",
				nil,
			)
		} else {
			zlog.Error("Error while finding user by userId", zap.String("userId", in.UserId), zap.Error(err))
			return errorx.NewInternalErr()
		}
	}

	// generate a uuid for user avatar
	avatarName, err := uuid.NewUUID()
	if err != nil {
		zlog.Error("Error while generate uuid: " + err.Error())
		return errorx.NewInternalErr()
	}

	err = minio.UploadUserAvatar(ctx, avatarName.String(), in.AvatarData)
	if err != nil {
		zlog.Error("Error while store user avatar into minio", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// modify user avatar
	avatarNameStr := avatarName.String()
	err = dao.UpdateUserByID(ctx, in.UserId, map[string]interface{}{
		"avatarUrl": avatarNameStr,
	})
	if err != nil {
		zlog.Error("Error while update user avatar", zap.String("userId", in.UserId), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (s *UserService) RefreshToken(ctx context.Context, userID string) (*sdto.RefreshTokenOutput, *errorx.ServiceErr) {
	user, err := dao.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewServicerErr(400, "User not found", nil)
		}

		zlog.Error("Error while find user by userID", zap.Error(err), zap.String("userID", userID))
		return nil, errorx.NewInternalErr()
	}

	// delete cache
	cacheKey := fmt.Sprintf("jwt:%v", user.Username)
	err = redis.RDB().Del(ctx, cacheKey).Err()
	if err != nil {
		zlog.Error("Error while delete jwt cache", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// generate new token
	organiser, err := dao.GetOrganiserByID(ctx, user.UserID)
	isOrganiser := false
	if err == nil && organiser != nil {
		isOrganiser = true
	}

	claims := jwt.MapClaims{
		"userID":         user.UserID,
		"isAdmin":        false,
		"exp":            time.Now().Add(24 * time.Hour).Unix(),
		"isOrganiser":    isOrganiser,
		"membershipType": user.MembershipType,
	}
	newToken, err := util.GenerateJWTToken(claims)
	if err != nil {
		zlog.Error("Error while generate new jwt token", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	err = redis.RDB().Set(ctx, cacheKey, newToken, 24*time.Hour).Err()
	if err != nil {
		zlog.Error("Error while set jwt cache", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	return &sdto.RefreshTokenOutput{
		NewToken: newToken,
	}, nil
}
