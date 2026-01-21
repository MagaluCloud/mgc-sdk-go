package objectstorage

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
)

// ObjectService provides operations for managing objects.
type ObjectService interface {
	Upload(ctx context.Context, bucketName string, objectKey string, data []byte, contentType string, storageClass *string) error
	Download(ctx context.Context, bucketName string, objectKey string, opts *DownloadOptions) ([]byte, error)
	DownloadStream(ctx context.Context, bucketName string, objectKey string, opts *DownloadStreamOptions) (io.Reader, error)
	DownloadAll(ctx context.Context, bucketName string, dst string, opts *DownloadAllOptions) (*DownloadAllResult, error)
	List(ctx context.Context, bucketName string, opts ObjectListOptions) ([]Object, error)
	ListAll(ctx context.Context, bucketName string, opts ObjectFilterOptions) ([]Object, error)
	ListVersions(ctx context.Context, bucketName string, objectKey string, opts *ListVersionsOptions) ([]ObjectVersion, error)
	ListAllVersions(ctx context.Context, bucketName string, objectKey string) ([]ObjectVersion, error)
	Delete(ctx context.Context, bucketName string, objectKey string, opts *DeleteOptions) error
	DeleteAll(ctx context.Context, bucketName string, opts *DeleteAllOptions) (*DeleteAllResult, error)
	Copy(ctx context.Context, src CopySrcConfig, dst CopyDstConfig) error
	CopyAll(ctx context.Context, src CopyPath, dst CopyPath, opts *CopyAllOptions) (*CopyAllResult, error)
	Metadata(ctx context.Context, bucketName string, objectKey string, opts *MetadataOptions) (*Object, error)
	LockObject(ctx context.Context, bucketName string, objectKey string, retainUntilDate time.Time) error
	UnlockObject(ctx context.Context, bucketName string, objectKey string) error
	GetObjectLockStatus(ctx context.Context, bucketName string, objectKey string) (bool, error)
	GetObjectLockInfo(ctx context.Context, bucketName string, objectKey string) (*ObjectLockInfo, error)
	GetPresignedURL(ctx context.Context, bucketName string, objectKey string, opts GetPresignedURLOptions) (*PresignedURL, error)
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

// DownloadAll downloads all objects from a bucket to a destination directory.
func (s *objectService) DownloadAll(
	ctx context.Context,
	bucketName string,
	dstPath string,
	opts *DownloadAllOptions,
) (*DownloadAllResult, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	if dstPath == "" {
		return nil, &InvalidObjectDataError{Message: "dstPath is empty"}
	}

	result := &DownloadAllResult{
		DownloadedCount: 0,
		ErrorCount:      0,
		Errors:          make([]DownloadError, 0),
	}

	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return result, err
	}

	listOpts := minio.ListObjectsOptions{
		Recursive: true,
	}

	if opts != nil && opts.Prefix != "" {
		listOpts.Prefix = opts.Prefix
	}

	maxParallel := 10

	objectCh := s.client.minioClient.ListObjects(ctx, bucketName, listOpts)

	workers := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for objectInfo := range objectCh {
		if objectInfo.Err != nil {
			mu.Lock()
			result.ErrorCount++
			result.Errors = append(result.Errors, DownloadError{
				ObjectKey: objectInfo.Key,
				Error:     objectInfo.Err,
			})
			mu.Unlock()
			continue
		}

		if opts != nil {
			shouldDownload := shouldProccessObject(opts.Filter, objectInfo.Key)

			if !shouldDownload {
				continue
			}
		}

		wg.Add(1)
		workers <- struct{}{}

		go func(obj minio.ObjectInfo) {
			defer wg.Done()
			defer func() { <-workers }()

			filePath := filepath.Join(dstPath, obj.Key)
			fileDir := filepath.Dir(filePath)

			if obj.Size == 0 && strings.HasSuffix(obj.Key, "/") {
				dirPath := filepath.Join(dstPath, obj.Key)

				if err := os.MkdirAll(dirPath, 0755); err != nil {
					mu.Lock()
					result.ErrorCount++
					result.Errors = append(result.Errors, DownloadError{
						ObjectKey: obj.Key,
						Error:     err,
					})
					mu.Unlock()
					return
				}

				mu.Lock()
				result.DownloadedCount++
				mu.Unlock()
				return
			}

			if err := os.MkdirAll(fileDir, 0755); err != nil {
				mu.Lock()
				result.ErrorCount++
				result.Errors = append(result.Errors, DownloadError{
					ObjectKey: obj.Key,
					Error:     err,
				})
				mu.Unlock()
				return
			}

			object, err := s.client.minioClient.GetObject(
				ctx,
				bucketName,
				obj.Key,
				minio.GetObjectOptions{},
			)
			if err != nil {
				mu.Lock()
				result.ErrorCount++
				result.Errors = append(result.Errors, DownloadError{
					ObjectKey: obj.Key,
					Error:     err,
				})
				mu.Unlock()
				return
			}
			defer object.Close()

			file, err := os.Create(filePath)
			if err != nil {
				mu.Lock()
				result.ErrorCount++
				result.Errors = append(result.Errors, DownloadError{
					ObjectKey: obj.Key,
					Error:     err,
				})
				mu.Unlock()
				return
			}

			defer file.Close()

			success := false
			defer func() {
				file.Close()
				object.Close()
				if !success {
					_ = os.Remove(filePath)
				}
			}()

			if _, err := io.Copy(file, object); err != nil {
				mu.Lock()
				result.ErrorCount++
				result.Errors = append(result.Errors, DownloadError{
					ObjectKey: obj.Key,
					Error:     err,
				})
				mu.Unlock()
				return
			}

			success = true

			mu.Lock()
			result.DownloadedCount++
			mu.Unlock()
		}(objectInfo)
	}

	wg.Wait()
	return result, nil
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
			shouldAddObj := shouldProccessObject(opts.Filter, object.Key)

			if shouldAddObj {
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

	listOpts := minio.ListObjectsOptions{
		Recursive: true,
	}

	if opts != nil && opts.ObjectKey != "" {
		listOpts.Prefix = opts.ObjectKey
	}

	listCh := s.client.minioClient.ListObjects(
		ctx,
		bucketName,
		listOpts,
	)

	batch := make([]minio.ObjectInfo, 0, batchSize)

	for obj := range listCh {
		if obj.Err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, DeleteError{
				ObjectKey: obj.Key,
				Error:     obj.Err,
			})
			continue
		}

		if opts != nil {
			shouldDeleteObject := shouldProccessObject(opts.Filter, obj.Key)

			if !shouldDeleteObject {
				continue
			}
		}

		batch = append(batch, minio.ObjectInfo{Key: obj.Key})

		if len(batch) == batchSize {
			s.deleteBatch(ctx, bucketName, batch, result)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		s.deleteBatch(ctx, bucketName, batch, result)
	}

	return result, nil
}

func (s *objectService) deleteBatch(
	ctx context.Context,
	bucketName string,
	batch []minio.ObjectInfo,
	result *DeleteAllResult,
) {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, obj := range batch {
			objectsCh <- obj
		}
	}()

	removeResultsCh := s.client.minioClient.RemoveObjects(
		ctx,
		bucketName,
		objectsCh,
		minio.RemoveObjectsOptions{},
	)

	for res := range removeResultsCh {
		if res.Err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, DeleteError{
				ObjectKey: res.ObjectName,
				Error:     res.Err,
			})
		} else {
			result.DeletedCount++
		}
	}
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

func (s *objectService) CopyAll(ctx context.Context, src CopyPath, dst CopyPath, opts *CopyAllOptions) (*CopyAllResult, error) {
	if src.BucketName == "" {
		return nil, &InvalidBucketNameError{Name: src.BucketName}
	}

	if dst.BucketName == "" {
		return nil, &InvalidBucketNameError{Name: dst.BucketName}
	}

	if opts != nil && opts.StorageClass != "" {
		err := storageClassIsValid(opts.StorageClass)

		if err != nil {
			return nil, err
		}

		ctx = WithStorageClass(ctx, opts.StorageClass)
	}

	result := &CopyAllResult{
		CopiedCount: 0,
		ErrorCount:  0,
		Errors:      make([]CopyError, 0),
	}

	maxParallel := 20

	listOpts := minio.ListObjectsOptions{
		Recursive: true,
	}

	if src.ObjectKey != "" {
		listOpts.Prefix = src.ObjectKey
	}

	objectCh := s.client.minioClient.ListObjects(ctx, src.BucketName, listOpts)

	workers := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for objectInfo := range objectCh {
		if objectInfo.Err != nil {
			mu.Lock()
			result.ErrorCount++
			result.Errors = append(result.Errors, CopyError{
				ObjectKey: objectInfo.Key,
				Error:     objectInfo.Err,
			})
			mu.Unlock()
			continue
		}

		if opts != nil {
			shouldCopy := shouldProccessObject(opts.Filter, objectInfo.Key)

			if !shouldCopy {
				continue
			}
		}

		wg.Add(1)
		workers <- struct{}{}

		go func(obj minio.ObjectInfo) {
			defer wg.Done()
			defer func() { <-workers }()

			if obj.Size == 0 && strings.HasSuffix(obj.Key, "/") {
				mu.Lock()
				result.CopiedCount++
				mu.Unlock()
				return
			}

			copySrc := minio.CopySrcOptions{
				Bucket: src.BucketName,
				Object: obj.Key,
			}

			dstObject := obj.Key

			if dst.ObjectKey != "" {
				dstObject = dst.ObjectKey + "/" + dstObject
			}

			copyDst := minio.CopyDestOptions{
				Bucket: dst.BucketName,
				Object: dstObject,
			}

			_, err := s.client.minioClient.CopyObject(ctx, copyDst, copySrc)
			if err != nil {
				mu.Lock()
				result.ErrorCount++
				result.Errors = append(result.Errors, CopyError{
					ObjectKey: obj.Key,
					Error:     err,
				})
				mu.Unlock()
				return
			}

			mu.Lock()
			result.CopiedCount++
			mu.Unlock()
		}(objectInfo)
	}

	wg.Wait()
	return result, nil
}

func storageClassIsValid(storageClass string) error {
	validStorageClasses := []string{"standard", "cold_instant"}

	if !slices.Contains(validStorageClasses, strings.ToLower(storageClass)) {
		return &InvalidObjectDataError{Message: "invalid storage class. Valid options are 'standard' and 'cold_instant'"}
	}

	return nil
}

func shouldProccessObject(filters *[]FilterOptions, objectKey string) bool {
	shouldDownload := true

	if filters != nil {
		for _, filter := range *filters {
			if filter.Include != "" && !matchesPattern(objectKey, filter.Include) {
				shouldDownload = false
				break
			}
			if filter.Exclude != "" && matchesPattern(objectKey, filter.Exclude) {
				shouldDownload = false
				break
			}
		}
	}

	return shouldDownload
}
