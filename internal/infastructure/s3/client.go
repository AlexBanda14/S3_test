package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	client     *minio.Client
	bucketName string
}

func NewS3Client(endpoint, accessKeyID, secretAccessKey, bucketName string, ctx context.Context) (*S3Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	ok, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket exists: %w", err)
	}

	if !ok {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &S3Client{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *S3Client) PutObject(key string, data *bytes.Buffer, ctx context.Context) error {
	_, err := s.client.PutObject(ctx, s.bucketName, key, data, int64(data.Len()), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})

	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (s *S3Client) GetObject(key string, ctx context.Context) (*bytes.Buffer, error) {
	obj, err := s.client.GetObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj)
	if err != nil {
		return nil, fmt.Errorf("read object: %w", err)
	}
	return buf, nil
}

func (s *S3Client) ListObject(ctx context.Context) ([]string, error) {
	var keys []string
	objects := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{})
	for obj := range objects {
		if obj.Err != nil {
			return nil, fmt.Errorf("list objects: %w", obj.Err)
		}
		keys = append(keys, obj.Key)
	}
	return keys, nil
}
