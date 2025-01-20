package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"os"
)

type MinIOFileRepository struct {
	m *Minio
}

func NewMinIOFileRepository(m *Minio) *MinIOFileRepository {
	return &MinIOFileRepository{m: m}
}

// UploadPublicFile uploads the file to MinIO, ensures the bucket exists, and returns the file's public URL
// UploadPublicFile uploads the file to MinIO, ensures the bucket exists, and returns the file's public URL
func (r *MinIOFileRepository) UploadPublicFile(ctx context.Context, bucketName, objectName, contentType string, file io.Reader) error {
	// If file is of type *os.File, we can get its size directly
	if f, ok := file.(*os.File); ok {
		fileInfo, err := f.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info: %w", err)
		}

		// Upload with known size
		_, err = r.m.M.PutObject(ctx, bucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
			ContentType: contentType,
		})
		return err
	}

	// For other types of readers (like multipart.File), we need to buffer to get size
	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		return fmt.Errorf("failed to copy file to buffer: %w", err)
	}

	// Upload from buffer with known size
	_, err = r.m.M.PutObject(ctx, bucketName, objectName, &buf, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (r *MinIOFileRepository) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	return r.m.M.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}
