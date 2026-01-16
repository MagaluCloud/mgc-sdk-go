package objectstorage

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

// ObjectService provides operations for managing objects.
type ObjectService interface {
	Upload(ctx context.Context, bucketName string, objectKey string, data []byte, contentType string, storageClass *string) error
	Download(ctx context.Context, bucketName string, objectKey string, opts *DownloadOptions) ([]byte, error)
	DownloadStream(ctx context.Context, bucketName string, objectKey string, opts *DownloadStreamOptions) (io.Reader, error)
	List(ctx context.Context, bucketName string, opts ObjectListOptions) ([]Object, error)
	ListAll(ctx context.Context, bucketName string, opts ObjectFilterOptions) ([]Object, error)
	ListVersions(ctx context.Context, bucketName string, objectKey string, opts *ListVersionsOptions) ([]ObjectVersion, error)
	ListAllVersions(ctx context.Context, bucketName string, objectKey string) ([]ObjectVersion, error)
	Delete(ctx context.Context, bucketName string, objectKey string, opts *DeleteOptions) error
	DeleteAll(ctx context.Context, bucketName string, opts *DeleteAllOptions) (*DeleteAllResult, error)
	Metadata(ctx context.Context, bucketName string, objectKey string, opts *MetadataOptions) (*Object, error)
	LockObject(ctx context.Context, bucketName string, objectKey string, retainUntilDate time.Time) error
	UnlockObject(ctx context.Context, bucketName string, objectKey string) error
	GetObjectLockStatus(ctx context.Context, bucketName string, objectKey string) (bool, error)
	GetObjectLockInfo(ctx context.Context, bucketName string, objectKey string) (*ObjectLockInfo, error)
	GetPresignedURL(ctx context.Context, bucketName string, objectKey string, opts GetPresignedURLOptions) (*PresignedURL, error)
	Copy(ctx context.Context, src CopySrcConfig, dst CopyDstConfig) error
}

// objectService implements the ObjectService interface.
type objectService struct {
	client *ObjectStorageClient
}

// Upload uploads an object to a bucket.
func (s *objectService) Upload(ctx context.Context, bucketName string, objectKey string, data []byte, contentType string, storageClass *string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return &InvalidObjectKeyError{Key: objectKey}
	}

	if len(data) == 0 {
		return &InvalidObjectDataError{Message: "object data cannot be empty"}
	}

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	if storageClass != nil && *storageClass != "" {
		err := storageClassIsValid(*storageClass)
		if err != nil {
			return err
		}

		opts.StorageClass = *storageClass
	}

	_, err := s.client.minioClient.PutObject(ctx, bucketName, objectKey, bytes.NewReader(data), int64(len(data)), opts)

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
			addObj := true

			if opts.Filter != nil {
				for _, filter := range *opts.Filter {
					if filter.Include != "" && !matchesPattern(object.Key, filter.Include) {
						addObj = false
						continue
					}
					if filter.Exclude != "" && matchesPattern(object.Key, filter.Exclude) {
						addObj = false
						continue
					}
				}
			}

			if addObj {
				result = append(result, Object{
					Key:          object.Key,
					Size:         object.Size,
					LastModified: object.LastModified,
					ETag:         object.ETag,
					StorageClass: object.StorageClass,
				})
			}
		}

		count++

		if opts.Limit != nil && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func matchesPattern(key, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(key)
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

// DeleteAll removes all objects from a bucket in batches based on filter criteria.
func (s *objectService) DeleteAll(ctx context.Context, bucketName string, opts *DeleteAllOptions) (*DeleteAllResult, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	result := &DeleteAllResult{
		DeletedCount: 0,
		ErrorCount:   0,
		Errors:       make([]DeleteError, 0),
	}

	batchSize := 1000
	if opts != nil && opts.BatchSize != nil && *opts.BatchSize > 0 && *opts.BatchSize <= 1000 {
		batchSize = *opts.BatchSize
	}

	listOpts := ObjectFilterOptions{
		Prefix:    "",
		Delimiter: "",
	}

	objects, err := s.ListAll(ctx, bucketName, listOpts)
	if err != nil {
		return nil, err
	}

	var objectsToDelete []Object
	if opts != nil && opts.Filter != nil {
		for _, obj := range objects {
			shouldDelete := true

			for _, filter := range *opts.Filter {
				if filter.Include != "" && !matchesPattern(obj.Key, filter.Include) {
					shouldDelete = false
					break
				}
				if filter.Exclude != "" && matchesPattern(obj.Key, filter.Exclude) {
					shouldDelete = false
					break
				}
			}

			if shouldDelete {
				objectsToDelete = append(objectsToDelete, obj)
			}
		}
	} else {
		objectsToDelete = objects
	}

	for i := 0; i < len(objectsToDelete); i += batchSize {
		end := i + batchSize
		if end > len(objectsToDelete) {
			end = len(objectsToDelete)
		}

		batch := objectsToDelete[i:end]

		objectsCh := make(chan minio.ObjectInfo, len(batch))
		for _, obj := range batch {
			objectsCh <- minio.ObjectInfo{Key: obj.Key}
		}
		close(objectsCh)

		removeResultsCh := s.client.minioClient.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{})

		for removeResult := range removeResultsCh {
			if removeResult.Err != nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, DeleteError{
					ObjectKey: removeResult.ObjectName,
					Error:     removeResult.Err,
				})
			} else {
				result.DeletedCount++
			}
		}
	}

	return result, nil
}

// Metadata returns metadata about an object.
func (s *objectService) Metadata(ctx context.Context, bucketName string, objectKey string, opts *MetadataOptions) (*Object, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	metadataOpts := minio.StatObjectOptions{}
	if opts != nil && opts.VersionID != "" {
		metadataOpts.VersionID = opts.VersionID
	}

	info, err := s.client.minioClient.StatObject(ctx, bucketName, objectKey, metadataOpts)
	if err != nil {
		return nil, err
	}

	storageClass := ""
	if val, ok := info.Metadata["X-Amz-Storage-Class"]; ok && len(val) > 0 {
		storageClass = val[0]
	}

	return &Object{
		Key:          info.Key,
		Size:         info.Size,
		LastModified: info.LastModified,
		ETag:         info.ETag,
		ContentType:  info.ContentType,
		StorageClass: storageClass,
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

	ctx = WithFixRetentionTime(ctx)

	mode, _, err := s.client.minioClient.GetObjectRetention(ctx, bucketName, objectKey, "")
	if err != nil {
		return false, err
	}

	// Object is locked if mode is set
	isLocked := mode != nil

	return isLocked, nil
}

// GetObjectLockInfo retrieves the lock information of an object.
func (s *objectService) GetObjectLockInfo(ctx context.Context, bucketName string, objectKey string) (*ObjectLockInfo, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	ctx = WithFixRetentionTime(ctx)

	mode, retentionUntilDate, err := s.client.minioClient.GetObjectRetention(ctx, bucketName, objectKey, "")
	if err != nil {
		return nil, err
	}

	if mode == nil {
		return &ObjectLockInfo{
			Locked: false,
		}, nil
	}

	return &ObjectLockInfo{
		Locked:          true,
		Mode:            mode.String(),
		RetainUntilDate: retentionUntilDate,
	}, nil
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

func (s *objectService) ListAllVersions(ctx context.Context, bucketName string, objectKey string) ([]ObjectVersion, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	result := make([]ObjectVersion, 0)
	objectVersionCh := s.client.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:       objectKey,
		Recursive:    true,
		WithVersions: true,
	})

	for objectInfo := range objectVersionCh {
		if objectInfo.Err != nil {
			return nil, objectInfo.Err
		}

		result = append(result, ObjectVersion{
			Key:            objectInfo.Key,
			VersionID:      objectInfo.VersionID,
			Size:           objectInfo.Size,
			LastModified:   objectInfo.LastModified,
			ETag:           objectInfo.ETag,
			IsDeleteMarker: objectInfo.IsDeleteMarker,
			IsLatest:       objectInfo.IsLatest,
			Owner: ObjectOwner{
				DisplayName: objectInfo.Owner.DisplayName,
				ID:          objectInfo.Owner.ID,
			},
			StorageClass: objectInfo.StorageClass,
		})
	}

	return result, nil
}

func (s *objectService) GetPresignedURL(ctx context.Context, bucketName string, objectKey string, opts GetPresignedURLOptions) (*PresignedURL, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if objectKey == "" {
		return nil, &InvalidObjectKeyError{Key: objectKey}
	}

	if opts.Method != http.MethodGet && opts.Method != http.MethodPut {
		return nil, &InvalidObjectDataError{Message: "Invalid HTTP method"}
	}

	var presignedURL *url.URL
	var err error

	expiryInSeconds := 5 * time.Minute

	if opts.ExpiryInSeconds != nil {
		expiryInSeconds = *opts.ExpiryInSeconds
	}

	switch opts.Method {
	case http.MethodGet:
		presignedURL, err = s.client.minioClient.PresignedGetObject(ctx, bucketName, objectKey, expiryInSeconds, url.Values{})
	case http.MethodPut:
		presignedURL, err = s.client.minioClient.PresignedPutObject(ctx, bucketName, objectKey, expiryInSeconds)
	}

	if err != nil {
		return nil, err
	}

	return &PresignedURL{URL: presignedURL.String()}, nil
}

func (s *objectService) Copy(ctx context.Context, src CopySrcConfig, dst CopyDstConfig) error {
	if src.BucketName == "" {
		return &InvalidBucketNameError{Name: src.BucketName}
	}
	if src.ObjectKey == "" {
		return &InvalidObjectKeyError{Key: src.ObjectKey}
	}

	if dst.BucketName == "" {
		return &InvalidBucketNameError{Name: dst.BucketName}
	}
	if dst.ObjectKey == "" {
		return &InvalidObjectKeyError{Key: dst.ObjectKey}
	}

	copyDst := minio.CopyDestOptions{
		Bucket: dst.BucketName,
		Object: dst.ObjectKey,
	}

	copySrc := minio.CopySrcOptions{
		Bucket: src.BucketName,
		Object: src.ObjectKey,
	}

	if src.VersionID != "" {
		copySrc.VersionID = src.VersionID
	}

	if dst.StorageClass != "" {
		err := storageClassIsValid(dst.StorageClass)
		if err != nil {
			return err
		}

		ctx = WithStorageClass(ctx, dst.StorageClass)
	}

	_, err := s.client.minioClient.CopyObject(
		ctx,
		copyDst,
		copySrc,
	)

	return err
}

func storageClassIsValid(storageClass string) error {
	validStorageClasses := []string{"standard", "cold_instant"}

	if !slices.Contains(validStorageClasses, strings.ToLower(storageClass)) {
		return &InvalidObjectDataError{Message: "invalid storage class. Valid options are 'standard' and 'cold_instant'"}
	}

	return nil
}
