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
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
)

// ObjectService provides operations for managing objects.
type ObjectService interface {
	Upload(ctx context.Context, bucketName string, objectKey string, data []byte, contentType string, storageClass *string) error
	UploadDir(ctx context.Context, bucketName string, objectKey string, srcDir string, opts *UploadDirOptions) (*UploadAllResult, error)
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

const (
	defaultBatchSize = 1000
	maxParallel      = 10
)

func validateBucket(bucket string) error {
	if bucket == "" {
		return &InvalidBucketNameError{Name: bucket}
	}
	return nil
}

func validateObjectKey(key string) error {
	if key == "" {
		return &InvalidObjectKeyError{Key: key}
	}
	return nil
}

func resolveBatchSize(size *int) int {
	if size != nil && *size > 0 {
		return *size
	}
	return defaultBatchSize
}

// Upload uploads an object to a bucket.
func (s *objectService) Upload(ctx context.Context, bucketName string, objectKey string, data []byte, contentType string, storageClass *string) error {
	if err := validateBucket(bucketName); err != nil {
		return err
	}

	if err := validateObjectKey(objectKey); err != nil {
		return err
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

	p := GetProgress(ctx)
	size := int64(len(data))

	var reader io.Reader = bytes.NewReader(data)

	var mu sync.Mutex

	if p != nil {
		mu.Lock()
		p.Start(size)
		mu.Unlock()
		reader = &ProgressReader{
			r: reader,
			p: p,
		}
		defer p.Finish()
	}

	_, err := s.client.minioClient.PutObject(ctx, bucketName, objectKey, reader, int64(len(data)), opts)

	return err
}

func (s *objectService) UploadDir(ctx context.Context, bucketName string, objectKey string, srcDir string, opts *UploadDirOptions) (*UploadAllResult, error) {
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}
	if srcDir == "" {
		return nil, &InvalidObjectDataError{Message: "srcDir is empty"}
	}

	batchSize := defaultBatchSize
	if opts != nil {
		batchSize = resolveBatchSize(&opts.BatchSize)
	}

	if opts != nil && opts.StorageClass != "" {
		err := storageClassIsValid(opts.StorageClass)
		if err != nil {
			return nil, err
		}

		ctx = WithStorageClass(ctx, opts.StorageClass)
	}

	var files []string

	err := filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if opts != nil && opts.Shallow && path != srcDir {
				return filepath.SkipDir
			}
			return nil
		}

		if opts != nil && !shouldProcessObject(opts.Filter, path) {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	fileCh := make(chan string, batchSize*2)

	go func() {
		defer close(fileCh)
		for _, f := range files {
			select {
			case <-ctx.Done():
				return
			case fileCh <- f:
			}
		}
	}()

	result := &UploadAllResult{
		Errors: make([]UploadError, 0),
	}

	var mu sync.Mutex

	handler := func(ctx context.Context, path string) error {
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstKey := filepath.ToSlash(filepath.Join(objectKey, relPath))

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		putOpts := minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		}

		if opts != nil && opts.StorageClass != "" {
			if err := storageClassIsValid(opts.StorageClass); err != nil {
				return err
			}
			putOpts.StorageClass = opts.StorageClass
		}

		_, err = s.client.minioClient.PutObject(
			ctx,
			bucketName,
			dstKey,
			file,
			-1,
			putOpts,
		)
		return err
	}

	p := GetProgress(ctx)

	if p != nil {
		mu.Lock()
		p.Start(int64(len(files)))
		mu.Unlock()
		defer p.Finish()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	processStreamInBatches(
		ctx,
		fileCh,
		batchSize,
		maxParallel,
		handler,
		func() {
			mu.Lock()
			result.UploadedCount++
			mu.Unlock()

			if p != nil {
				mu.Lock()
				p.Add(1)
				mu.Unlock()
			}
		},
		func(err error) {
			mu.Lock()
			result.ErrorCount++
			result.Errors = append(result.Errors, UploadError{Error: err})
			mu.Unlock()
		},
	)

	return result, nil
}

func (s *objectService) getObject(ctx context.Context, bucketName, objectKey string, versionID string) (*minio.Object, error) {
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}
	if err := validateObjectKey(objectKey); err != nil {
		return nil, err
	}

	opts := minio.GetObjectOptions{}
	if versionID != "" {
		opts.VersionID = versionID
	}

	return s.client.minioClient.GetObject(ctx, bucketName, objectKey, opts)
}

// Download retrieves an object from a bucket and returns its content as bytes.
func (s *objectService) Download(ctx context.Context, bucketName string, objectKey string, opts *DownloadOptions) ([]byte, error) {
	versionID := ""
	if opts != nil {
		versionID = opts.VersionID
	}

	object, err := s.getObject(ctx, bucketName, objectKey, versionID)
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return io.ReadAll(object)
}

// DownloadStream retrieves an object from a bucket and returns a reader for streaming.
func (s *objectService) DownloadStream(ctx context.Context, bucketName string, objectKey string, opts *DownloadStreamOptions) (io.Reader, error) {
	versionID := ""
	if opts != nil {
		versionID = opts.VersionID
	}

	return s.getObject(ctx, bucketName, objectKey, versionID)
}

// DownloadAll downloads all objects from a bucket to a destination directory.
func (s *objectService) DownloadAll(ctx context.Context, bucketName string, dstPath string, opts *DownloadAllOptions) (*DownloadAllResult, error) {
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}
	if dstPath == "" {
		return nil, &InvalidObjectDataError{Message: "dstPath is empty"}
	}

	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return nil, err
	}

	batchSize := defaultBatchSize
	if opts != nil {
		batchSize = resolveBatchSize(&opts.BatchSize)
	}

	listOpts := minio.ListObjectsOptions{Recursive: true}
	if opts != nil && opts.Prefix != "" {
		listOpts.Prefix = opts.Prefix
	}

	var objects []minio.ObjectInfo

	for obj := range s.client.minioClient.ListObjects(ctx, bucketName, listOpts) {
		if obj.Err != nil {
			return nil, obj.Err
		}

		if opts != nil && !shouldProcessObject(opts.Filter, obj.Key) {
			continue
		}

		objects = append(objects, obj)
	}

	total := int64(len(objects))

	var mu sync.Mutex

	p := GetProgress(ctx)
	if p != nil {
		mu.Lock()
		p.Start(total)
		mu.Unlock()
		defer p.Finish()
	}

	objectCh := make(chan minio.ObjectInfo, batchSize*2)

	go func() {
		defer close(objectCh)
		for _, obj := range objects {
			select {
			case <-ctx.Done():
				return
			case objectCh <- obj:
			}
		}
	}()

	result := &DownloadAllResult{
		Errors: make([]DownloadError, 0),
	}

	handler := func(ctx context.Context, obj minio.ObjectInfo) error {
		if obj.Size == 0 && strings.HasSuffix(obj.Key, "/") {
			return os.MkdirAll(filepath.Join(dstPath, obj.Key), 0755)
		}

		filePath := filepath.Join(dstPath, obj.Key)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		object, err := s.client.minioClient.GetObject(
			ctx,
			bucketName,
			obj.Key,
			minio.GetObjectOptions{},
		)
		if err != nil {
			return err
		}
		defer object.Close()

		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, object)
		return err
	}

	processStreamInBatches(
		ctx,
		objectCh,
		batchSize,
		maxParallel,
		handler,
		func() {
			mu.Lock()
			result.DownloadedCount++
			mu.Unlock()

			if p != nil {
				mu.Lock()
				p.Add(1)
				mu.Unlock()
			}
		},
		func(err error) {
			mu.Lock()
			result.ErrorCount++
			result.Errors = append(result.Errors, DownloadError{Error: err})
			mu.Unlock()
		},
	)

	return result, nil
}

// List retrieves a list of objects in a bucket with pagination.
func (s *objectService) List(ctx context.Context, bucketName string, opts ObjectListOptions) ([]Object, error) {
	if err := validateBucket(bucketName); err != nil {
		return nil, err
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
			shouldAddObj := shouldProcessObject(opts.Filter, object.Key)

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
	if err := validateBucket(bucketName); err != nil {
		return nil, err
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
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}

	batchSize := defaultBatchSize
	if opts != nil {
		batchSize = resolveBatchSize(opts.BatchSize)
	}

	listOpts := minio.ListObjectsOptions{Recursive: true}
	if opts != nil && opts.ObjectKey != "" {
		listOpts.Prefix = opts.ObjectKey
	}

	var objects []minio.ObjectInfo

	for obj := range s.client.minioClient.ListObjects(ctx, bucketName, listOpts) {
		if obj.Err != nil {
			return nil, obj.Err
		}

		if opts != nil && !shouldProcessObject(opts.Filter, obj.Key) {
			continue
		}

		objects = append(objects, obj)
	}

	total := int64(len(objects))

	var mu sync.Mutex

	p := GetProgress(ctx)
	if p != nil {
		mu.Lock()
		p.Start(total)
		mu.Unlock()
		defer p.Finish()
	}

	objectCh := make(chan minio.ObjectInfo, batchSize*2)

	go func() {
		defer close(objectCh)
		for _, obj := range objects {
			select {
			case <-ctx.Done():
				return
			case objectCh <- obj:
			}
		}
	}()

	result := &DeleteAllResult{
		Errors: make([]DeleteError, 0),
	}

	handler := func(ctx context.Context, obj minio.ObjectInfo) error {
		return s.client.minioClient.RemoveObject(
			ctx,
			bucketName,
			obj.Key,
			minio.RemoveObjectOptions{},
		)
	}

	processStreamInBatches(
		ctx,
		objectCh,
		batchSize,
		maxParallel,
		handler,
		func() {
			mu.Lock()
			result.DeletedCount++
			mu.Unlock()

			if p != nil {
				mu.Lock()
				p.Add(1)
				mu.Unlock()
			}
		},
		func(err error) {
			mu.Lock()
			result.ErrorCount++
			result.Errors = append(result.Errors, DeleteError{Error: err})
			mu.Unlock()
		},
	)

	return result, nil
}

// Metadata returns metadata about an object.
func (s *objectService) Metadata(ctx context.Context, bucketName string, objectKey string, opts *MetadataOptions) (*Object, error) {
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}

	if err := validateObjectKey(objectKey); err != nil {
		return nil, err
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
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}

	if err := validateObjectKey(objectKey); err != nil {
		return nil, err
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
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}

	if err := validateObjectKey(objectKey); err != nil {
		return nil, err
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
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}

	if err := validateObjectKey(objectKey); err != nil {
		return nil, err
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
	if err := validateBucket(bucketName); err != nil {
		return nil, err
	}

	if err := validateObjectKey(objectKey); err != nil {
		return nil, err
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
	if err := validateBucket(src.BucketName); err != nil {
		return err
	}
	if err := validateObjectKey(src.ObjectKey); err != nil {
		return err
	}

	if err := validateBucket(dst.BucketName); err != nil {
		return err
	}
	if err := validateObjectKey(dst.ObjectKey); err != nil {
		return err
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

func (s *objectService) CopyAll(
	ctx context.Context,
	src CopyPath,
	dst CopyPath,
	opts *CopyAllOptions,
) (*CopyAllResult, error) {
	if err := validateBucket(src.BucketName); err != nil {
		return nil, err
	}
	if err := validateBucket(dst.BucketName); err != nil {
		return nil, err
	}

	batchSize := defaultBatchSize
	if opts != nil {
		batchSize = resolveBatchSize(&opts.BatchSize)
	}

	if opts != nil && opts.StorageClass != "" {
		err := storageClassIsValid(opts.StorageClass)

		if err != nil {
			return nil, err
		}

		ctx = WithStorageClass(ctx, opts.StorageClass)
	}

	listOpts := minio.ListObjectsOptions{Recursive: true}
	if src.ObjectKey != "" {
		listOpts.Prefix = src.ObjectKey
	}

	var objects []minio.ObjectInfo

	for obj := range s.client.minioClient.ListObjects(ctx, src.BucketName, listOpts) {
		if obj.Err != nil {
			return nil, obj.Err
		}

		if opts != nil && !shouldProcessObject(opts.Filter, obj.Key) {
			continue
		}

		if obj.Size == 0 && strings.HasSuffix(obj.Key, "/") {
			continue
		}

		objects = append(objects, obj)
	}

	total := int64(len(objects))

	var mu sync.Mutex

	p := GetProgress(ctx)
	if p != nil {
		mu.Lock()
		p.Start(total)
		mu.Unlock()
		defer p.Finish()
	}

	objectCh := make(chan minio.ObjectInfo, batchSize*2)

	go func() {
		defer close(objectCh)
		for _, obj := range objects {
			select {
			case <-ctx.Done():
				return
			case objectCh <- obj:
			}
		}
	}()

	result := &CopyAllResult{
		Errors: make([]CopyError, 0),
	}

	handler := func(ctx context.Context, obj minio.ObjectInfo) error {
		dstKey := obj.Key
		if dst.ObjectKey != "" {
			dstKey = filepath.ToSlash(filepath.Join(dst.ObjectKey, obj.Key))
		}

		_, err := s.client.minioClient.CopyObject(
			ctx,
			minio.CopyDestOptions{
				Bucket: dst.BucketName,
				Object: dstKey,
			},
			minio.CopySrcOptions{
				Bucket: src.BucketName,
				Object: obj.Key,
			},
		)
		return err
	}

	processStreamInBatches(
		ctx,
		objectCh,
		batchSize,
		maxParallel,
		handler,
		func() {
			mu.Lock()
			result.CopiedCount++
			mu.Unlock()

			if p != nil {
				mu.Lock()
				p.Add(1)
				mu.Unlock()
			}
		},
		func(err error) {
			mu.Lock()
			result.ErrorCount++
			result.Errors = append(result.Errors, CopyError{Error: err})
			mu.Unlock()
		},
	)

	return result, nil
}

func storageClassIsValid(storageClass string) error {
	switch strings.ToLower(storageClass) {
	case "standard", "cold_instant":
		return nil
	default:
		return &InvalidObjectDataError{
			Message: "invalid storage class. Valid options are 'standard' and 'cold_instant'",
		}
	}
}

func shouldProcessObject(filters *[]FilterOptions, objectKey string) bool {
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

func processStreamInBatches[T any](
	ctx context.Context,
	input <-chan T,
	batchSize int,
	maxParallel int,
	handler func(context.Context, T) error,
	onSuccess func(),
	onError func(error),
) {
	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, maxParallel)
	)

	batch := make([]T, 0, batchSize)

	flush := func(items []T) {
		for _, item := range items {
			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			sem <- struct{}{}

			go func(it T) {
				defer wg.Done()
				defer func() { <-sem }()

				if err := handler(ctx, it); err != nil {
					onError(err)
					return
				}
				onSuccess()
			}(item)
		}

		wg.Wait()
	}

	for {
		select {
		case <-ctx.Done():
			return

		case item, ok := <-input:
			if !ok {
				if len(batch) > 0 {
					flush(batch)
				}
				return
			}

			batch = append(batch, item)

			if len(batch) == batchSize {
				flush(batch)
				batch = batch[:0]
			}
		}
	}
}

func (pr *ProgressReader) Read(b []byte) (int, error) {
	n, err := pr.r.Read(b)
	if n > 0 {
		pr.p.Add(int64(n))
	}
	return n, err
}
