package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// Minio basic operations

// CreateBucket() will create a new bucket in minio, if bucket already exist, nothing will happen
func CreateBucket(bucketName string) error {
	if len(bucketName) <= 0 {
		return fmt.Errorf("invalid bucket name")
	}

	ctx := context.Background()

	// check if bucket already exist or not
	isFound, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	} else if isFound {
		return nil
	}

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "beijing"})
	if err != nil {
		return err
	}

	return nil
}

func UploadUserAvatar(ctx context.Context, avatarName string, avatarData []byte) error {
	avatarReader := bytes.NewReader(avatarData)

	err := UploadFile(ctx, AVATAR_BUCKET, avatarName, avatarReader, int64(len(avatarData)), "image/jpeg")
	if err != nil {
		return err
	}

	return nil
}

func UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := minioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetObjectUrl(ctx context.Context, bucketName, objectName string, exp time.Duration) (*url.URL, error) {
	if exp < 1 {
		exp = 24 * time.Hour
	}

	preSignedUrl, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, exp, nil)
	if err != nil {
		return nil, err
	}

	return preSignedUrl, nil
}
