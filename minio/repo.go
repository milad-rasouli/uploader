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
	m                *Minio
	disableMultiPart bool
}

func NewMinIOFileRepository(m *Minio, disableMultiPart bool) *MinIOFileRepository {
	return &MinIOFileRepository{m: m, disableMultiPart: disableMultiPart}
}

// UploadPublicFile uploads the file to MinIO, ensures the bucket exists, and returns the file's public URL
// UploadPublicFile uploads the file to MinIO, ensures the bucket exists, and returns the file's public URL
func (r *MinIOFileRepository) UploadPublicFile(ctx context.Context, bucketName, objectName, contentType string, file io.Reader) error {
	var size int64

	// Check if it's a bytes.Reader to get size directly
	if rs, ok := file.(*bytes.Reader); ok {
		size = int64(rs.Len())
	} else if f, ok := file.(*os.File); ok {
		// Handle regular files
		fileInfo, err := f.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info: %w", err)
		}
		size = fileInfo.Size()
	} else {
		// For other readers, buffer the content
		var buf bytes.Buffer
		var err error
		size, err = io.Copy(&buf, file)
		if err != nil {
			return fmt.Errorf("failed to copy file to buffer: %w", err)
		}
		file = &buf
	}

	_, err := r.m.M.PutObject(ctx, bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType:      contentType,
		DisableMultipart: r.disableMultiPart,
	})
	return err
}

func (r *MinIOFileRepository) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	return r.m.M.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}
