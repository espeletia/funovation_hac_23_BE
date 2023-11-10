package fileStorage

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type FileStorageInterface interface {
	UploadFile(ctx context.Context, fileSrc string, dest string, contentType string) error
}

func NewFileS3Storage(s3client *s3.Client) *FileS3Storage {
	return &FileS3Storage{
		s3client: s3client,
		tracer:   otel.GetTracerProvider().Tracer("S3FileManager"),
	}
}

type FileS3Storage struct {
	tracer   trace.Tracer
	s3client *s3.Client
}

func (fs *FileS3Storage) UploadFile(ctx context.Context, fileSrc string, dest string, contentType string) error {
	u, err := url.Parse(dest)
	if err != nil {
		return err
	}
	data, err := os.Open(filepath.Clean(fileSrc))
	defer func() {
		err = data.Close()
		if err != nil {
			zap.L().Error("error closing file", zap.Error(err))
		}
	}()
	if err != nil {
		return err
	}
	key := strings.TrimPrefix(u.Path, "/")
	uploader := manager.NewUploader(fs.s3client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.Host),
		Key:         aws.String(key),
		Body:        data,
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return err
	}

	return nil
}
