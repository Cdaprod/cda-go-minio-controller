package main

import (
    "context"
    "fmt"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOOperations represents MinIO operations
type MinIOOperations struct {
    client *minio.Client
}

// NewMinIOOperations creates a new instance of MinIOOperations
func NewMinIOOperations(endpoint, accessKey, secretKey string, secure bool) (*MinIOOperations, error) {
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: secure,
    })
    if err != nil {
        return nil, err
    }
    return &MinIOOperations{client: client}, nil
}

// CreateBucket creates a new bucket in MinIO
func (m *MinIOOperations) CreateBucket(bucketName string) error {
    ctx := context.Background()
    err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
    if err != nil {
        return fmt.Errorf("failed to create bucket %s: %v", bucketName, err)
    }
    return nil
}

// ... (Add other MinIO operations as needed)