package objectstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

// ObjectService provides operations for managing objects.
type ObjectService interface {
	Upload(ctx context.Context, bucketName string, objectKey string, data []byte, contentType string) error
	Download(ctx context.Context, bucketName string, objectKey string, opts *DownloadOptions) ([]byte, error)
	DownloadStream(ctx context.Context, bucketName string, objectKey string, opts *DownloadStreamOptions) (io.Reader, error)
	List(ctx context.Context, bucketName string, opts ObjectListOptions) ([]Object, error)
	ListAll(ctx context.Context, bucketName string, opts ObjectFilterOptions) ([]Object, error)
	ListVersions(ctx context.Context, bucketName string, objectKey string, opts *ListVersionsOptions) ([]ObjectVersion, error)
	Delete(ctx context.Context, bucketName string, objectKey string, opts *DeleteOptions) error
	DeleteMany(ctx context.Context, bucketName string, objectKeys []string, opts minio.RemoveObjectsOptions) error
	Metadata(ctx context.Context, bucketName string, objectKey string) (*Object, error)
	LockObject(ctx context.Context, bucketName string, objectKey string, retainUntilDate time.Time) error
	UnlockObject(ctx context.Context, bucketName string, objectKey string) error
	GetObjectLockStatus(ctx context.Context, bucketName string, objectKey string) (bool, error)
}

// objectService implements the ObjectService interface.
type objectService struct {
	client *ObjectStorageClient
}

// Upload uploads an object to a bucket.
func (s *objectService) Upload(ctx context.Context, bucketName string, objectKey string, data []byte, contentType string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return &InvalidObjectKeyError{Key: objectKey}
	}

	if len(data) == 0 {
		return &InvalidObjectDataError{Message: "object data cannot be empty"}
	}

	_, err := s.client.minioClient.PutObject(ctx, bucketName, objectKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

// Download retrieves an object from a bucket and returns its content as bytes.
func (s *objectService) Download(ctx context.Context, bucketName string, objectKey string, opts *DownloadOptions) ([]byte, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	getOpts := minio.GetObjectOptions{}
	if opts != nil && opts.VersionID != "" {
		getOpts.VersionID = opts.VersionID
	}

	object, err := s.client.minioClient.GetObject(ctx, bucketName, objectKey, getOpts)
	if err != nil {
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// DownloadStream retrieves an object from a bucket and returns a reader for streaming.
func (s *objectService) DownloadStream(ctx context.Context, bucketName string, objectKey string, opts *DownloadStreamOptions) (io.Reader, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	getOpts := minio.GetObjectOptions{}
	if opts != nil && opts.VersionID != "" {
		getOpts.VersionID = opts.VersionID
	}

	object, err := s.client.minioClient.GetObject(ctx, bucketName, objectKey, getOpts)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// List retrieves a list of objects in a bucket with pagination.
func (s *objectService) List(ctx context.Context, bucketName string, opts ObjectListOptions) ([]Object, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	result := make([]Object, 0)
	objectCh := s.client.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    opts.Prefix,
		Recursive: opts.Delimiter == "",
	})

	limit := 50
	offset := 0

	if opts.Limit != nil {
		limit = *opts.Limit
	}

	if opts.Offset != nil {
		offset = *opts.Offset
	}

	count := 0
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		if count >= offset && count < offset+limit {
			result = append(result, Object{
				Key:          object.Key,
				Size:         object.Size,
				LastModified: object.LastModified,
				ETag:         object.ETag,
			})
		}

		count++

		if opts.Limit != nil && len(result) >= limit {
			break
		}
	}

	return result, nil
}

// ListAll retrieves all objects in a bucket without pagination.
func (s *objectService) ListAll(ctx context.Context, bucketName string, opts ObjectFilterOptions) ([]Object, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	result := make([]Object, 0)
	objectCh := s.client.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    opts.Prefix,
		Recursive: opts.Delimiter == "",
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		result = append(result, Object{
			Key:          object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ETag:         object.ETag,
		})
	}

	return result, nil
}

// Delete removes an object from a bucket.
func (s *objectService) Delete(ctx context.Context, bucketName string, objectKey string, opts *DeleteOptions) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return &InvalidObjectKeyError{Key: objectKey}
	}

	removeOpts := minio.RemoveObjectOptions{}
	if opts != nil && opts.VersionID != "" {
		removeOpts.VersionID = opts.VersionID
	}

	return s.client.minioClient.RemoveObject(ctx, bucketName, objectKey, removeOpts)
}

// DeleteMany removes objects from a bucket.
func (s *objectService) DeleteMany(ctx context.Context, bucketName string, objectKeys []string, opts minio.RemoveObjectsOptions) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if len(objectKeys) == 0 {
		return &InvalidObjectKeyError{Key: ""}
	}

	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)

		for _, key := range objectKeys {
			if key == "" {
				continue
			}

			objectsCh <- minio.ObjectInfo{
				Key: key,
			}
		}
	}()

	errorCh := s.client.minioClient.RemoveObjects(ctx, bucketName, objectsCh, opts)

	var errs []string

	for err := range errorCh {
		errs = append(errs, fmt.Sprintf("%s: %v", err.ObjectName, err.Err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete some objects: %s", strings.Join(errs, "; "))
	}

	return nil
}

// Metadata returns metadata about an object.
func (s *objectService) Metadata(ctx context.Context, bucketName string, objectKey string) (*Object, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	info, err := s.client.minioClient.StatObject(ctx, bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &Object{
		Key:          info.Key,
		Size:         info.Size,
		LastModified: info.LastModified,
		ETag:         info.ETag,
		ContentType:  info.ContentType,
	}, nil
}

// LockObject applies a retention lock to an object until the specified date.
func (s *objectService) LockObject(ctx context.Context, bucketName string, objectKey string, retainUntilDate time.Time) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return &InvalidObjectKeyError{Key: objectKey}
	}

	if retainUntilDate.IsZero() {
		return &InvalidObjectDataError{Message: "retain until date cannot be zero"}
	}

	// Use COMPLIANCE mode for object locking
	complianceMode := minio.Compliance

	opts := minio.PutObjectRetentionOptions{
		Mode:             &complianceMode,
		RetainUntilDate:  &retainUntilDate,
		GovernanceBypass: false,
	}

	return s.client.minioClient.PutObjectRetention(ctx, bucketName, objectKey, opts)
}

// UnlockObject removes the retention lock from an object.
func (s *objectService) UnlockObject(ctx context.Context, bucketName string, objectKey string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return &InvalidObjectKeyError{Key: objectKey}
	}

	// Set empty retention to remove lock
	opts := minio.PutObjectRetentionOptions{
		Mode:             nil,
		RetainUntilDate:  nil,
		GovernanceBypass: true,
	}

	return s.client.minioClient.PutObjectRetention(ctx, bucketName, objectKey, opts)
}

// GetObjectLockStatus retrieves the lock status of an object.
func (s *objectService) GetObjectLockStatus(ctx context.Context, bucketName string, objectKey string) (bool, error) {
	if bucketName == "" {
		return false, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return false, &InvalidObjectKeyError{Key: objectKey}
	}

	mode, _, err := s.client.minioClient.GetObjectRetention(ctx, bucketName, objectKey, "")
	if err != nil {
		return false, err
	}

	// Object is locked if mode is set
	isLocked := mode != nil

	return isLocked, nil
}

// ListVersions retrieves all versions of an object from a versioned bucket.
func (s *objectService) ListVersions(ctx context.Context, bucketName string, objectKey string, opts *ListVersionsOptions) ([]ObjectVersion, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	result := make([]ObjectVersion, 0)
	objectVersionCh := s.client.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    objectKey,
		Recursive: true,
	})

	limit := 50
	offset := 0

	if opts != nil {
		if opts.Limit != nil {
			limit = *opts.Limit
		}
		if opts.Offset != nil {
			offset = *opts.Offset
		}
	}

	count := 0
	for objectInfo := range objectVersionCh {
		if objectInfo.Err != nil {
			return nil, objectInfo.Err
		}

		// Only include versions for the exact object key (not prefixes)
		if objectInfo.Key == objectKey {
			if count >= offset && count < offset+limit {
				result = append(result, ObjectVersion{
					Key:          objectInfo.Key,
					VersionID:    objectInfo.VersionID,
					Size:         objectInfo.Size,
					LastModified: objectInfo.LastModified,
					ETag:         objectInfo.ETag,
				})
			}
			count++
		}
	}

	return result, nil
}
