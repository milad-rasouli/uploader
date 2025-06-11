package uploader

import (
	"context"
	"io"
)

//go:generate mockgen -source=file.go -destination=./mock/file/file.go
type FileRepository interface {
	UploadPublicFile(ctx context.Context, bucketName, objectName, contentType string, file io.Reader, userMetadata map[string]string) error
	DeleteFile(ctx context.Context, bucketName string, objectName string) error
}
