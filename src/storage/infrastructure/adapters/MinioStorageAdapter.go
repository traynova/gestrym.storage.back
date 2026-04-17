package adapters

import (
	"context"
	"fmt"
	"gestrym-storage/src/storage/domain"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type minioStorageAdapter struct {
	client *minio.Client
	bucket string
}

func NewMinioStorageAdapter() (domain.IStorageAdapter, error) {
	endpoint := viper.GetString("MINIO_ENDPOINT")
	accessKey := viper.GetString("MINIO_ACCESS_KEY")
	secretKey := viper.GetString("MINIO_SECRET_KEY")
	bucket := viper.GetString("MINIO_BUCKET")
	useSSL := viper.GetBool("MINIO_USE_SSL")

	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	if bucket == "" {
		bucket = "gestrym-storage"
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	// Make bucket if not exists
	ctx := context.Background()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
	if errBucketExists == nil && !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &minioStorageAdapter{
		client: minioClient,
		bucket: bucket,
	}, nil
}

func (m *minioStorageAdapter) UploadFile(fileStream io.Reader, objectSize int64, contentType, fileName, collectionID string) (string, error) {
	ctx := context.Background()
	objectName := fmt.Sprintf("%s/%s", collectionID, fileName)
	_, err := m.client.PutObject(ctx, m.bucket, objectName, fileStream, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (m *minioStorageAdapter) DeleteFile(fileName, collectionID string) error {
	ctx := context.Background()
	objectName := fmt.Sprintf("%s/%s", collectionID, fileName)
	return m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
}

func (m *minioStorageAdapter) GetFileURL(fileName, collectionID string) (string, error) {
	ctx := context.Background()
	objectName := fmt.Sprintf("%s/%s", collectionID, fileName)
	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, time.Hour*24, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
