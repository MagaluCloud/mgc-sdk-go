package objectstorage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/cors"
)

// minioClientInterface defines the interface for MinIO client operations
// This allows for mocking in tests while using the real client in production
type minioClientInterface interface {
	// Bucket operations
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	RemoveBucket(ctx context.Context, bucketName string) error
	GetBucketPolicy(ctx context.Context, bucketName string) (string, error)
	SetBucketPolicy(ctx context.Context, bucketName string, policy string) error
	GetObjectLockConfig(ctx context.Context, bucketName string) (string, *minio.RetentionMode, *uint, *minio.ValidityUnit, error)
	SetObjectLockConfig(ctx context.Context, bucketName string, mode *minio.RetentionMode, validity *uint, unit *minio.ValidityUnit) error
	GetBucketCors(ctx context.Context, bucketName string) (*cors.Config, error)
	SetBucketCors(ctx context.Context, bucketName string, corsConfig *cors.Config) error
	GetBucketVersioning(ctx context.Context, bucketName string) (minio.BucketVersioningConfiguration, error)
	EnableVersioning(ctx context.Context, bucketName string) error
	SuspendVersioning(ctx context.Context, bucketName string) error

	// Object operations
	PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error
	RemoveObjects(ctx context.Context, bucketName string, objectsCh <-chan minio.ObjectInfo, opts minio.RemoveObjectsOptions) <-chan minio.RemoveObjectError
	StatObject(ctx context.Context, bucketName string, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	PutObjectRetention(ctx context.Context, bucketName string, objectName string, opts minio.PutObjectRetentionOptions) error
	GetObjectRetention(ctx context.Context, bucketName string, objectName string, versionID string) (*minio.RetentionMode, *time.Time, error)
	SetAppInfo(appName string, appVersion string)
	PresignedGetObject(ctx context.Context, bucketName string, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error)
	PresignedPutObject(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (*url.URL, error)
	CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (minio.UploadInfo, error)
}

// Ensure *minio.Client implements minioClientInterface
var _ minioClientInterface = (*minio.Client)(nil)
