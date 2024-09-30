package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

type MinIOFileRepository struct {
	m *Minio
}

func NewMinIOFileRepository(m *Minio) *MinIOFileRepository {
	return &MinIOFileRepository{m: m}
}

// UploadPublicFile uploads the file to MinIO, ensures the bucket exists, and returns the file's public URL
func (r *MinIOFileRepository) UploadPublicFile(ctx context.Context, bucketName, objectName, contentType string, file io.Reader) error {

	// Upload the file to MinIO
	_, err := r.m.M.PutObject(ctx, bucketName, objectName, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *MinIOFileRepository) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	return r.m.M.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}
