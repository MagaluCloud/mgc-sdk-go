package objectstorage

import (
	"context"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/cors"
)

// TestBucketServiceList_WithMockSuccess tests List with mock MinIO returning buckets
func TestBucketServiceList_WithMockSuccess(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["bucket1"] = &mockBucket{
		name:         "bucket1",
		creationDate: time.Now().Add(-24 * time.Hour),
		objects:      make(map[string]*mockObject),
	}
	mock.buckets["bucket2"] = &mockBucket{
		name:         "bucket2",
		creationDate: time.Now(),
		objects:      make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	buckets, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(buckets) != 2 {
		t.Errorf("List() returned %d buckets, want 2", len(buckets))
	}

	// Verify bucket data is properly mapped
	for _, b := range buckets {
		if b.Name == "" {
			t.Error("List() bucket has empty name")
		}
		if b.CreationDate.IsZero() {
			t.Error("List() bucket has zero creation date")
		}
	}
}

// TestBucketServiceGetPolicy_WithMockSuccess tests GetPolicy with mock returning policy
func TestBucketServiceGetPolicy_WithMockSuccess(t *testing.T) {
	t.Parallel()

	policyJSON := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:GetObject"],"Resource":"arn:aws:s3:::test-bucket/*"}]}`

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		policy:       policyJSON,
		objects:      make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	policy, err := svc.GetPolicy(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetPolicy() error = %v", err)
	}

	if policy == nil {
		t.Fatal("GetPolicy() returned nil policy")
	}

	if policy.Version != "2012-10-17" {
		t.Errorf("GetPolicy() version = %s, want 2012-10-17", policy.Version)
	}

	if len(policy.Statement) != 1 {
		t.Errorf("GetPolicy() len(Statement) = %d, want 1", len(policy.Statement))
	}
}

// TestBucketServiceGetPolicy_EmptyPolicy tests GetPolicy when policy is empty
func TestBucketServiceGetPolicy_EmptyPolicy(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		policy:       "",
		objects:      make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	policy, err := svc.GetPolicy(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetPolicy() error = %v", err)
	}

	if policy != nil {
		t.Errorf("GetPolicy() expected nil for empty policy, got %+v", policy)
	}
}

// TestBucketServiceGetPolicy_InvalidJSON tests GetPolicy with invalid JSON
func TestBucketServiceGetPolicy_InvalidJSON(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		policy:       "{invalid json}",
		objects:      make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	_, err := svc.GetPolicy(context.Background(), "test-bucket")
	if err == nil {
		t.Error("GetPolicy() expected error for invalid JSON, got nil")
	}
}

// TestBucketServiceGetBucketLockStatus_Locked tests GetBucketLockStatus when bucket is locked
func TestBucketServiceGetBucketLockStatus_Locked(t *testing.T) {
	t.Parallel()

	complianceMode := minio.Compliance
	validity := uint(1)
	unit := minio.Days

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		lockConfig: &mockLockConfig{
			objectLock: "Enabled",
			mode:       &complianceMode,
			validity:   &validity,
			unit:       &unit,
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	isLocked, err := svc.GetBucketLockStatus(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetBucketLockStatus() error = %v", err)
	}

	if !isLocked {
		t.Error("GetBucketLockStatus() expected true for locked bucket, got false")
	}
}

// TestBucketServiceGetBucketLockStatus_Unlocked tests GetBucketLockStatus when bucket is not locked
func TestBucketServiceGetBucketLockStatus_Unlocked(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		lockConfig:   nil,
		objects:      make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	isLocked, err := svc.GetBucketLockStatus(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetBucketLockStatus() error = %v", err)
	}

	if isLocked {
		t.Error("GetBucketLockStatus() expected false for unlocked bucket, got true")
	}
}

// TestBucketServiceGetBucketLockStatus_PartialConfig tests GetBucketLockStatus with partial lock config
func TestBucketServiceGetBucketLockStatus_PartialConfig(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		lockConfig: &mockLockConfig{
			objectLock: "Enabled",
			mode:       nil, // missing mode
			validity:   nil,
			unit:       nil,
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	isLocked, err := svc.GetBucketLockStatus(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetBucketLockStatus() error = %v", err)
	}

	if isLocked {
		t.Error("GetBucketLockStatus() expected false for partial config, got true")
	}
}

// TestBucketServiceGetCORS_WithMockSuccess tests GetCORS with mock returning CORS config
func TestBucketServiceGetCORS_WithMockSuccess(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		corsConfig: &cors.Config{
			CORSRules: []cors.Rule{
				{
					AllowedOrigin: []string{"*"},
					AllowedMethod: []string{"GET", "PUT"},
					AllowedHeader: []string{"*"},
					ExposeHeader:  []string{"ETag"},
					MaxAgeSeconds: 3600,
				},
			},
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	corsConfig, err := svc.GetCORS(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetCORS() error = %v", err)
	}

	if corsConfig == nil {
		t.Fatal("GetCORS() returned nil")
	}

	if len(corsConfig.CORSRules) != 1 {
		t.Errorf("GetCORS() len(CORSRules) = %d, want 1", len(corsConfig.CORSRules))
	}

	rule := corsConfig.CORSRules[0]
	if len(rule.AllowedOrigins) != 1 || rule.AllowedOrigins[0] != "*" {
		t.Errorf("GetCORS() AllowedOrigins = %v, want [*]", rule.AllowedOrigins)
	}

	if len(rule.AllowedMethods) != 2 {
		t.Errorf("GetCORS() len(AllowedMethods) = %d, want 2", len(rule.AllowedMethods))
	}

	if rule.MaxAgeSeconds != 3600 {
		t.Errorf("GetCORS() MaxAgeSeconds = %d, want 3600", rule.MaxAgeSeconds)
	}
}

// TestBucketServiceGetCORS_NilConfig tests GetCORS when CORS config is nil
func TestBucketServiceGetCORS_NilConfig(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		corsConfig:   nil,
		objects:      make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	corsConfig, err := svc.GetCORS(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetCORS() error = %v", err)
	}

	if corsConfig != nil {
		t.Errorf("GetCORS() expected nil for no CORS config, got %+v", corsConfig)
	}
}

// TestBucketServiceGetCORS_EmptyRules tests GetCORS when CORS has no rules
func TestBucketServiceGetCORS_EmptyRules(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		corsConfig: &cors.Config{
			CORSRules: []cors.Rule{},
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	corsConfig, err := svc.GetCORS(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetCORS() error = %v", err)
	}

	if corsConfig != nil {
		t.Errorf("GetCORS() expected nil for empty CORS rules, got %+v", corsConfig)
	}
}

// TestBucketServiceGetCORS_MultipleRules tests GetCORS with multiple CORS rules
func TestBucketServiceGetCORS_MultipleRules(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		corsConfig: &cors.Config{
			CORSRules: []cors.Rule{
				{
					AllowedOrigin: []string{"https://example.com"},
					AllowedMethod: []string{"GET"},
					AllowedHeader: []string{"Authorization"},
					MaxAgeSeconds: 3600,
				},
				{
					AllowedOrigin: []string{"https://api.example.com"},
					AllowedMethod: []string{"POST", "PUT"},
					AllowedHeader: []string{"Content-Type"},
					ExposeHeader:  []string{"ETag", "Content-Length"},
					MaxAgeSeconds: 7200,
				},
			},
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	corsConfig, err := svc.GetCORS(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetCORS() error = %v", err)
	}

	if corsConfig == nil {
		t.Fatal("GetCORS() returned nil")
	}

	if len(corsConfig.CORSRules) != 2 {
		t.Errorf("GetCORS() len(CORSRules) = %d, want 2", len(corsConfig.CORSRules))
	}

	// Verify first rule
	rule1 := corsConfig.CORSRules[0]
	if len(rule1.AllowedOrigins) != 1 || rule1.AllowedOrigins[0] != "https://example.com" {
		t.Errorf("GetCORS() rule1 AllowedOrigins = %v", rule1.AllowedOrigins)
	}

	// Verify second rule
	rule2 := corsConfig.CORSRules[1]
	if len(rule2.AllowedMethods) != 2 {
		t.Errorf("GetCORS() rule2 len(AllowedMethods) = %d, want 2", len(rule2.AllowedMethods))
	}
	if len(rule2.ExposeHeaders) != 2 {
		t.Errorf("GetCORS() rule2 len(ExposeHeaders) = %d, want 2", len(rule2.ExposeHeaders))
	}
}

// TestBucketServiceGetVersioningStatus_Enabled tests GetVersioningStatus when enabled
func TestBucketServiceGetVersioningStatus_Enabled(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		versioning: minio.BucketVersioningConfiguration{
			Status: "Enabled",
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	config, err := svc.GetVersioningStatus(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetVersioningStatus() error = %v", err)
	}

	if config == nil {
		t.Fatal("GetVersioningStatus() returned nil")
	}

	if config.Status != VersioningStatusEnabled {
		t.Errorf("GetVersioningStatus() status = %s, want %s", config.Status, VersioningStatusEnabled)
	}
}

// TestBucketServiceGetVersioningStatus_Suspended tests GetVersioningStatus when suspended
func TestBucketServiceGetVersioningStatus_Suspended(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		versioning: minio.BucketVersioningConfiguration{
			Status: "Suspended",
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	config, err := svc.GetVersioningStatus(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetVersioningStatus() error = %v", err)
	}

	if config == nil {
		t.Fatal("GetVersioningStatus() returned nil")
	}

	if config.Status != VersioningStatusSuspended {
		t.Errorf("GetVersioningStatus() status = %s, want %s", config.Status, VersioningStatusSuspended)
	}
}

// TestBucketServiceGetVersioningStatus_Off tests GetVersioningStatus when off
func TestBucketServiceGetVersioningStatus_Off(t *testing.T) {
	t.Parallel()

	mock := newMockMinioClient()
	mock.buckets["test-bucket"] = &mockBucket{
		name:         "test-bucket",
		creationDate: time.Now(),
		versioning: minio.BucketVersioningConfiguration{
			Status: "",
		},
		objects: make(map[string]*mockObject),
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Buckets()

	config, err := svc.GetVersioningStatus(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("GetVersioningStatus() error = %v", err)
	}

	if config == nil {
		t.Fatal("GetVersioningStatus() returned nil")
	}

	if config.Status != VersioningStatusOff {
		t.Errorf("GetVersioningStatus() status = %s, want %s", config.Status, VersioningStatusOff)
	}
}
