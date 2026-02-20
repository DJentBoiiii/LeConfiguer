package storage

import (
	"context"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOStorage is a MinIO bucket-based implementation of the Storage interface.
type MinIOStorage struct {
	mu         sync.RWMutex
	client     *minio.Client
	bucketName string
	ctx        context.Context
}

// NewMinIOStorage creates a new MinIO-based storage instance.
func NewMinIOStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinIOStorage, error) {
	ctx := context.Background()

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Check if bucket exists, if not create it
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIOStorage{
		client:     client,
		bucketName: bucketName,
		ctx:        ctx,
	}, nil
}
