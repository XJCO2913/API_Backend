package minio

import (
	"context"
	"fmt"

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
