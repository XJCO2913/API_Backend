package minio

import (
	"api.backend.xjco2913/util/config"
	"api.backend.xjco2913/util/zlog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// bucket
const (
	AVATAR_BUCKET = "user-avatar"
	MOMENT_BUCKET = "user-moment"
	ACTIVITY_BUCKET = "activity"
)

var (
	minioClient *minio.Client
)

func init() {
	endpoint := config.Get("database.minio.endpoint")
	accessKeyId := config.Get("database.minio.accessKeyId")
	secretAccessKey := config.Get("database.minio.secretAccessKey")

	client, err := minio.New(endpoint, &minio.Options{
		Secure: false,
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
	})
	if err != nil {
		panic(err)
	}

	minioClient = client

	// init bucket
	err = CreateBucket(AVATAR_BUCKET)
	if err != nil {
		panic(err)
	}
	err = CreateBucket(MOMENT_BUCKET)
	if err != nil {
		panic(err)
	}
	err = CreateBucket(ACTIVITY_BUCKET)
	if err != nil {
		panic(err)
	}

	zlog.Info("minio client init successfully")
}
