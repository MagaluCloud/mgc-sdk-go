package objectstorage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/cors"
)

// mockMinioClient is a mock implementation of the MinIO client for testing
type mockMinioClient struct {
	// Storage for mock data
	buckets                map[string]*mockBucket
	listBucketsFunc        func(ctx context.Context) ([]minio.BucketInfo, error)
	makeBucketFunc         func(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
	bucketExistsFunc       func(ctx context.Context, bucketName string) (bool, error)
	removeBucketFunc       func(ctx context.Context, bucketName string) error
	getBucketPolicyFunc    func(ctx context.Context, bucketName string) (string, error)
	setBucketPolicyFunc    func(ctx context.Context, bucketName string, policy string) error
	getLockConfigFunc      func(ctx context.Context, bucketName string) (string, *minio.RetentionMode, *uint, *minio.ValidityUnit, error)
	setLockConfigFunc      func(ctx context.Context, bucketName string, mode *minio.RetentionMode, validity *uint, unit *minio.ValidityUnit) error
	getCorsFunc            func(ctx context.Context, bucketName string) (*cors.Config, error)
	setCorsFunc            func(ctx context.Context, bucketName string, corsConfig *cors.Config) error
	getVersioningFunc      func(ctx context.Context, bucketName string) (minio.BucketVersioningConfiguration, error)
	enableVersioningFunc   func(ctx context.Context, bucketName string) error
	suspendVersioningFunc  func(ctx context.Context, bucketName string) error
	putObjectFunc          func(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	getObjectFunc          func(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	listObjectsFunc        func(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	removeObjectFunc       func(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error
	removeObjectsFunc      func(ctx context.Context, bucketName string, objectsCh <-chan minio.ObjectInfo, opts minio.RemoveObjectsOptions) <-chan minio.RemoveObjectError
	statObjectFunc         func(ctx context.Context, bucketName string, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	putObjectRetentionFunc func(ctx context.Context, bucketName string, objectName string, opts minio.PutObjectRetentionOptions) error
	getObjectRetentionFunc func(ctx context.Context, bucketName string, objectName string, versionID string) (*minio.RetentionMode, *time.Time, error)
	presignedGetObjectFunc func(ctx context.Context, bucketName string, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error)
	presignedPutObjectFunc func(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (*url.URL, error)
	copyObjectFunc         func(context.Context, minio.CopyDestOptions, minio.CopySrcOptions) (minio.UploadInfo, error)

	setAppInfoCalls int
	lastAppName     string
	lastAppVersion  string
}

type mockBucket struct {
	name         string
	creationDate time.Time
	policy       string
	corsConfig   *cors.Config
	versioning   minio.BucketVersioningConfiguration
	lockConfig   *mockLockConfig
	objects      map[string]*mockObject
}

type mockLockConfig struct {
	objectLock string
	mode       *minio.RetentionMode
	validity   *uint
	unit       *minio.ValidityUnit
}

type mockObject struct {
	key          string
	size         int64
	lastModified time.Time
	etag         string
	contentType  string
	data         []byte
	retention    *mockObjectRetention
}

type mockObjectRetention struct {
	mode            *minio.RetentionMode
	retainUntilDate *time.Time
}

// newMockMinioClient creates a new mock MinIO client
func newMockMinioClient() *mockMinioClient {
	return &mockMinioClient{
		buckets: make(map[string]*mockBucket),
	}
}

// ListBuckets mocks the MinIO ListBuckets method
func (m *mockMinioClient) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	if m.listBucketsFunc != nil {
		return m.listBucketsFunc(ctx)
	}

	var buckets []minio.BucketInfo
	for _, b := range m.buckets {
		buckets = append(buckets, minio.BucketInfo{
			Name:         b.name,
			CreationDate: b.creationDate,
		})
	}
	return buckets, nil
}

// MakeBucket mocks the MinIO MakeBucket method
func (m *mockMinioClient) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	if m.makeBucketFunc != nil {
		return m.makeBucketFunc(ctx, bucketName, opts)
	}

	m.buckets[bucketName] = &mockBucket{
		name:         bucketName,
		creationDate: time.Now(),
		objects:      make(map[string]*mockObject),
	}
	return nil
}

// BucketExists mocks the MinIO BucketExists method
func (m *mockMinioClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	if m.bucketExistsFunc != nil {
		return m.bucketExistsFunc(ctx, bucketName)
	}

	_, exists := m.buckets[bucketName]
	return exists, nil
}

// RemoveBucket mocks the MinIO RemoveBucket method
func (m *mockMinioClient) RemoveBucket(ctx context.Context, bucketName string) error {
	if m.removeBucketFunc != nil {
		return m.removeBucketFunc(ctx, bucketName)
	}

	delete(m.buckets, bucketName)
	return nil
}

// GetBucketPolicy mocks the MinIO GetBucketPolicy method
func (m *mockMinioClient) GetBucketPolicy(ctx context.Context, bucketName string) (string, error) {
	if m.getBucketPolicyFunc != nil {
		return m.getBucketPolicyFunc(ctx, bucketName)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return "", nil
	}
	return bucket.policy, nil
}

// SetBucketPolicy mocks the MinIO SetBucketPolicy method
func (m *mockMinioClient) SetBucketPolicy(ctx context.Context, bucketName string, policy string) error {
	if m.setBucketPolicyFunc != nil {
		return m.setBucketPolicyFunc(ctx, bucketName, policy)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil
	}
	bucket.policy = policy
	return nil
}

// GetObjectLockConfig mocks the MinIO GetObjectLockConfig method
func (m *mockMinioClient) GetObjectLockConfig(ctx context.Context, bucketName string) (string, *minio.RetentionMode, *uint, *minio.ValidityUnit, error) {
	if m.getLockConfigFunc != nil {
		return m.getLockConfigFunc(ctx, bucketName)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists || bucket.lockConfig == nil {
		return "", nil, nil, nil, nil
	}

	return bucket.lockConfig.objectLock, bucket.lockConfig.mode, bucket.lockConfig.validity, bucket.lockConfig.unit, nil
}

// SetObjectLockConfig mocks the MinIO SetObjectLockConfig method
func (m *mockMinioClient) SetObjectLockConfig(ctx context.Context, bucketName string, mode *minio.RetentionMode, validity *uint, unit *minio.ValidityUnit) error {
	if m.setLockConfigFunc != nil {
		return m.setLockConfigFunc(ctx, bucketName, mode, validity, unit)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil
	}

	if mode == nil && validity == nil && unit == nil {
		bucket.lockConfig = nil
	} else {
		bucket.lockConfig = &mockLockConfig{
			objectLock: "Enabled",
			mode:       mode,
			validity:   validity,
			unit:       unit,
		}
	}
	return nil
}

// GetBucketCors mocks the MinIO GetBucketCors method
func (m *mockMinioClient) GetBucketCors(ctx context.Context, bucketName string) (*cors.Config, error) {
	if m.getCorsFunc != nil {
		return m.getCorsFunc(ctx, bucketName)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil, nil
	}
	return bucket.corsConfig, nil
}

// SetBucketCors mocks the MinIO SetBucketCors method
func (m *mockMinioClient) SetBucketCors(ctx context.Context, bucketName string, corsConfig *cors.Config) error {
	if m.setCorsFunc != nil {
		return m.setCorsFunc(ctx, bucketName, corsConfig)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil
	}
	bucket.corsConfig = corsConfig
	return nil
}

// GetBucketVersioning mocks the MinIO GetBucketVersioning method
func (m *mockMinioClient) GetBucketVersioning(ctx context.Context, bucketName string) (minio.BucketVersioningConfiguration, error) {
	if m.getVersioningFunc != nil {
		return m.getVersioningFunc(ctx, bucketName)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return minio.BucketVersioningConfiguration{}, nil
	}
	return bucket.versioning, nil
}

// EnableVersioning mocks the MinIO EnableVersioning method
func (m *mockMinioClient) EnableVersioning(ctx context.Context, bucketName string) error {
	if m.enableVersioningFunc != nil {
		return m.enableVersioningFunc(ctx, bucketName)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil
	}
	bucket.versioning.Status = "Enabled"
	return nil
}

// SuspendVersioning mocks the MinIO SuspendVersioning method
func (m *mockMinioClient) SuspendVersioning(ctx context.Context, bucketName string) error {
	if m.suspendVersioningFunc != nil {
		return m.suspendVersioningFunc(ctx, bucketName)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil
	}
	bucket.versioning.Status = "Suspended"
	return nil
}

// PutObject mocks the MinIO PutObject method
func (m *mockMinioClient) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	if m.putObjectFunc != nil {
		return m.putObjectFunc(ctx, bucketName, objectName, reader, objectSize, opts)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return minio.UploadInfo{}, nil
	}

	bucket.objects[objectName] = &mockObject{
		key:          objectName,
		size:         objectSize,
		lastModified: time.Now(),
		etag:         "mock-etag",
		contentType:  opts.ContentType,
	}

	return minio.UploadInfo{
		Bucket: bucketName,
		Key:    objectName,
		ETag:   "mock-etag",
		Size:   objectSize,
	}, nil
}

// GetObject mocks the MinIO GetObject method
func (m *mockMinioClient) GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	if m.getObjectFunc != nil {
		return m.getObjectFunc(ctx, bucketName, objectName, opts)
	}

	// Return nil for mock - actual object reading would need more complex mocking
	return nil, nil
}

// ListObjects mocks the MinIO ListObjects method
func (m *mockMinioClient) ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	if m.listObjectsFunc != nil {
		return m.listObjectsFunc(ctx, bucketName, opts)
	}

	ch := make(chan minio.ObjectInfo)
	go func() {
		defer close(ch)
		bucket, exists := m.buckets[bucketName]
		if !exists {
			return
		}

		for _, obj := range bucket.objects {
			ch <- minio.ObjectInfo{
				Key:          obj.key,
				Size:         obj.size,
				LastModified: obj.lastModified,
				ETag:         obj.etag,
				ContentType:  obj.contentType,
			}
		}
	}()
	return ch
}

// RemoveObject mocks the MinIO RemoveObject method
func (m *mockMinioClient) RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error {
	if m.removeObjectFunc != nil {
		return m.removeObjectFunc(ctx, bucketName, objectName, opts)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil
	}
	delete(bucket.objects, objectName)
	return nil
}

// StatObject mocks the MinIO StatObject method
func (m *mockMinioClient) StatObject(ctx context.Context, bucketName string, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	if m.statObjectFunc != nil {
		return m.statObjectFunc(ctx, bucketName, objectName, opts)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return minio.ObjectInfo{}, nil
	}

	obj, exists := bucket.objects[objectName]
	if !exists {
		return minio.ObjectInfo{}, nil
	}

	return minio.ObjectInfo{
		Key:          obj.key,
		Size:         obj.size,
		LastModified: obj.lastModified,
		ETag:         obj.etag,
		ContentType:  obj.contentType,
	}, nil
}

// PutObjectRetention mocks the MinIO PutObjectRetention method
func (m *mockMinioClient) PutObjectRetention(ctx context.Context, bucketName string, objectName string, opts minio.PutObjectRetentionOptions) error {
	if m.putObjectRetentionFunc != nil {
		return m.putObjectRetentionFunc(ctx, bucketName, objectName, opts)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil
	}

	obj, exists := bucket.objects[objectName]
	if !exists {
		return nil
	}

	obj.retention = &mockObjectRetention{
		mode:            opts.Mode,
		retainUntilDate: opts.RetainUntilDate,
	}
	return nil
}

// GetObjectRetention mocks the MinIO GetObjectRetention method
func (m *mockMinioClient) GetObjectRetention(ctx context.Context, bucketName string, objectName string, versionID string) (*minio.RetentionMode, *time.Time, error) {
	if m.getObjectRetentionFunc != nil {
		return m.getObjectRetentionFunc(ctx, bucketName, objectName, versionID)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil, nil, nil
	}

	obj, exists := bucket.objects[objectName]
	if !exists || obj.retention == nil {
		return nil, nil, nil
	}

	return obj.retention.mode, obj.retention.retainUntilDate, nil
}

func (m *mockMinioClient) PresignedGetObject(ctx context.Context, bucketName string, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	if m.presignedGetObjectFunc != nil {
		return m.presignedGetObjectFunc(ctx, bucketName, objectName, expiry, reqParams)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil, nil
	}

	obj, exists := bucket.objects[objectName]
	if !exists {
		return nil, nil
	}

	mockURL := "https://mock-minio/" + bucketName + "/" + obj.key + "?expiry=" + expiry.String()

	parsedURL, err := url.Parse(mockURL)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
}

func (m *mockMinioClient) PresignedPutObject(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (*url.URL, error) {
	if m.presignedPutObjectFunc != nil {
		return m.presignedPutObjectFunc(ctx, bucketName, objectName, expiry)
	}

	bucket, exists := m.buckets[bucketName]
	if !exists {
		return nil, nil
	}

	obj, exists := bucket.objects[objectName]
	if !exists {
		return nil, nil
	}

	mockURL := "https://mock-minio/" + bucketName + "/" + obj.key + "?expiry=" + expiry.String()

	parsedURL, err := url.Parse(mockURL)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
}

func (m *mockMinioClient) SetAppInfo(appName, appVersion string) {
	m.setAppInfoCalls++
	m.lastAppName = appName
	m.lastAppVersion = appVersion
}

func (m *mockMinioClient) CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (minio.UploadInfo, error) {
	if m.copyObjectFunc != nil {
		return m.copyObjectFunc(ctx, dst, src)
	}

	return minio.UploadInfo{
		Bucket: dst.Bucket,
		Key:    dst.Object,
	}, nil
}

func (m *mockMinioClient) RemoveObjects(ctx context.Context, bucketName string, objectsCh <-chan minio.ObjectInfo, opts minio.RemoveObjectsOptions) <-chan minio.RemoveObjectError {
	if m.removeObjectsFunc != nil {
		return m.removeObjectsFunc(ctx, bucketName, objectsCh, opts)
	}

	return nil
}
