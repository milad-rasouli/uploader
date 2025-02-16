package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"time"
)

type Option func(m *Minio)

type Config struct {
	MinioHost        string
	MinioAccessToken string
	MinioSecret      string
	Secure           bool
	Buckets          []string
	DisableMultiPart bool
	ExpirationDays   uint16 //amount of days to expire and delete objects
}

func WithConfig(config *Config) Option {
	return func(m *Minio) {
		m.conf = config
		m.DisableMultiPart = config.DisableMultiPart
	}
}

type Minio struct {
	conf             *Config
	M                *minio.Client
	DisableMultiPart bool
}

func NewMinio(options ...Option) *Minio {
	m := &Minio{}
	for i := 0; i < len(options); i++ {
		options[i](m)
	}
	return m
}

func (m *Minio) Setup(ctx context.Context) error {
	minioClient, err := minio.New(m.conf.MinioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(m.conf.MinioAccessToken, m.conf.MinioSecret, ""),
		Secure: m.conf.Secure,
	})
	if err != nil {
		return err
	}

	for _, bucket := range m.conf.Buckets {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
			if errBucketExists != nil && !exists {
				return err
			}
		}
		if m.conf.ExpirationDays > 0 {
			config := lifecycle.NewConfiguration()
			config.Rules = []lifecycle.Rule{
				{
					ID:     bucket,
					Status: "Enabled",
					Expiration: lifecycle.Expiration{
						Days: lifecycle.ExpirationDays(m.conf.ExpirationDays),
					},
				},
			}

			err = minioClient.SetBucketLifecycle(context.Background(), bucket, config)
			if err != nil {
				return err
			}
		}
	}

	m.M = minioClient
	return nil
}

func (m *Minio) GeneratePublicURL(bucketName, objectName string) string {
	scheme := "http://"
	if m.conf.Secure {
		scheme = "https://"
	}
	return GeneratePublicURL(scheme, m.conf.MinioHost, bucketName, objectName)
}

// ReadinessCheck verifies that the Minio client can interact with the Minio server.
func (m *Minio) ReadinessCheck() error {
	// Check if the connection to the Minio server is healthy by listing Buckets
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := m.M.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("minio readiness check failed: %w", err)
	}
	return nil
}

func GeneratePublicURL(scheme, minioHost, bucketName, objectName string) string {
	return fmt.Sprintf("%s%s/%s/%s", scheme, minioHost, bucketName, objectName)
}
