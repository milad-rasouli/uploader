package file

import (
	"context"
	"io"
)

type FileRepository interface {
	UploadPublicFile(ctx context.Context, bucketName, objectName, contentType string, file io.Reader) error
	GeneratePublicURL(bucketName string, objectName string) string
	DeleteFile(ctx context.Context, bucketName string, objectName string) error
}
