package main

import (
	"context"
	"fmt"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
)

func ClearJunkAvatarFile(ctx context.Context) {
	users, err := dao.GetAllUsers(ctx)
	if err != nil {
		panic(err)
	}

	inUsedAvatar := map[string]bool{}
	for _, user := range users {
		if user.AvatarURL == nil || *user.AvatarURL == "" {
			continue
		}

		inUsedAvatar[*user.AvatarURL] = true
	}

	// remove
	minio.RemoveUnusedAvatar(ctx, inUsedAvatar)
	fmt.Println("Clear junk avatar obj successfully")
}

func ClearJunkMomentMediaFile(ctx context.Context) {
	moments, err := dao.GetAllMoment(ctx)
	if err != nil {
		panic(err)
	}

	inUsedMomentMedia := map[string]bool{}
	for _, moment := range moments {
		if moment.ImageURL != nil && *moment.ImageURL != "" {
			inUsedMomentMedia[*moment.ImageURL] = true
		}
		if moment.VideoURL != nil && *moment.VideoURL != "" {
			inUsedMomentMedia[*moment.VideoURL] = true
		}
	}

	// remove
	minio.RemoveUnusedMomentMedia(ctx, inUsedMomentMedia)
	fmt.Println("Clear junk moment obj successfully")
}

func ClearJunkActivityImageFile(ctx context.Context) {
	activities, err := dao.GetAllActivities(ctx)
	if err != nil {
		panic(err)
	}

	inUsedImage := map[string]bool{}
	for _, activity := range activities {
		if activity.CoverURL == "" {
			continue
		}

		inUsedImage[activity.CoverURL] = true
	}

	// remove
	minio.RemoveUnusedActivityImage(ctx, inUsedImage)
	fmt.Println("Clear junk activity image obj successfully")
}