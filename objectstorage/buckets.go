package objectstorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/cors"
)

// BucketService provides operations for managing buckets.
type BucketService interface {
	Create(ctx context.Context, bucketName string) error
	List(ctx context.Context) ([]Bucket, error)
	Exists(ctx context.Context, bucketName string) (bool, error)
	Delete(ctx context.Context, bucketName string, recursive bool) error
	GetPolicy(ctx context.Context, bucketName string) (*Policy, error)
	SetPolicy(ctx context.Context, bucketName string, policy *Policy) error
	DeletePolicy(ctx context.Context, bucketName string) error
	LockBucket(ctx context.Context, bucketName string) error
	UnlockBucket(ctx context.Context, bucketName string) error
	GetBucketLockStatus(ctx context.Context, bucketName string) (bool, error)
	SetCORS(ctx context.Context, bucketName string, corsConfig *CORSConfiguration) error
	GetCORS(ctx context.Context, bucketName string) (*CORSConfiguration, error)
	DeleteCORS(ctx context.Context, bucketName string) error
	EnableVersioning(ctx context.Context, bucketName string) error
	SuspendVersioning(ctx context.Context, bucketName string) error
	GetVersioningStatus(ctx context.Context, bucketName string) (*BucketVersioningConfiguration, error)
}

// bucketService implements the BucketService interface.
type bucketService struct {
	client *ObjectStorageClient
}

// Create creates a new bucket.
func (s *bucketService) Create(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	return s.client.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

// List retrieves all buckets.
func (s *bucketService) List(ctx context.Context) ([]Bucket, error) {
	buckets, err := s.client.minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Bucket, len(buckets))
	for i, b := range buckets {
		result[i] = Bucket{
			Name:         b.Name,
			CreationDate: b.CreationDate,
		}
	}

	return result, nil
}

// Exists checks if a bucket exists.
func (s *bucketService) Exists(ctx context.Context, bucketName string) (bool, error) {
	if bucketName == "" {
		return false, &InvalidBucketNameError{Name: bucketName}
	}

	return s.client.minioClient.BucketExists(ctx, bucketName)
}

// Delete deletes a bucket.
func (s *bucketService) Delete(ctx context.Context, bucketName string, recursive bool) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if recursive {
		objects, err := s.client.Objects().ListAll(ctx, bucketName, ObjectFilterOptions{})
		if err != nil {
			return fmt.Errorf("error to get all objects: %w", err)
		}

		if len(objects) != 0 {
			objectKeys := []string{}

			for _, object := range objects {
				objectKeys = append(objectKeys, object.Key)
			}

			err = s.client.Objects().DeleteMany(ctx, bucketName, objectKeys, minio.RemoveObjectsOptions{})
			if err != nil {
				return fmt.Errorf("error to delete the objects: %w", err)
			}
		}
	}

	return s.client.minioClient.RemoveBucket(ctx, bucketName)
}

// GetPolicy retrieves the policy of a bucket.
func (s *bucketService) GetPolicy(ctx context.Context, bucketName string) (*Policy, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	policyStr, err := s.client.minioClient.GetBucketPolicy(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if policyStr == "" {
		return nil, nil
	}

	var policy Policy
	err = unmarshalPolicy(policyStr, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// SetPolicy sets the policy of a bucket.
func (s *bucketService) SetPolicy(ctx context.Context, bucketName string, policy *Policy) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if policy == nil {
		return &InvalidPolicyError{Message: "policy cannot be nil"}
	}

	if len(policy.Statement) == 0 {
		return &InvalidPolicyError{Message: "policy must have at least one statement"}
	}

	policyStr, err := marshalPolicy(policy)
	if err != nil {
		return err
	}

	return s.client.minioClient.SetBucketPolicy(ctx, bucketName, policyStr)
}

// DeletePolicy removes the policy from a bucket.
func (s *bucketService) DeletePolicy(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	return s.client.minioClient.SetBucketPolicy(ctx, bucketName, "")
}

// marshalPolicy converts a Policy struct to a JSON string.
func marshalPolicy(policy *Policy) (string, error) {
	data, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// unmarshalPolicy converts a JSON policy string to a Policy struct.
func unmarshalPolicy(policyStr string, policy *Policy) error {
	return json.Unmarshal([]byte(policyStr), policy)
}

// LockBucket enables Object Lock for a bucket.
func (s *bucketService) LockBucket(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	// Use COMPLIANCE mode as the default retention mode for bucket locking
	complianceMode := minio.Compliance
	validity := uint(1)
	unit := minio.Days

	return s.client.minioClient.SetObjectLockConfig(ctx, bucketName, &complianceMode, &validity, &unit)
}

// UnlockBucket disables Object Lock for a bucket by removing the configuration.
func (s *bucketService) UnlockBucket(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	// Remove lock configuration by setting all parameters to nil
	return s.client.minioClient.SetObjectLockConfig(ctx, bucketName, nil, nil, nil)
}

// GetBucketLockStatus retrieves the lock status of a bucket.
func (s *bucketService) GetBucketLockStatus(ctx context.Context, bucketName string) (bool, error) {
	if bucketName == "" {
		return false, &InvalidBucketNameError{Name: bucketName}
	}

	objectLock, mode, validity, unit, err := s.client.minioClient.GetObjectLockConfig(ctx, bucketName)
	if err != nil {
		return false, err
	}

	// Bucket is locked if objectLock string is not empty and mode is set
	isLocked := objectLock != "" && mode != nil && validity != nil && unit != nil

	return isLocked, nil
}

// SetCORS sets the CORS configuration for a bucket.
func (s *bucketService) SetCORS(ctx context.Context, bucketName string, corsConfig *CORSConfiguration) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	if corsConfig == nil {
		return &InvalidPolicyError{Message: "CORS configuration cannot be nil"}
	}

	if len(corsConfig.CORSRules) == 0 {
		return &InvalidPolicyError{Message: "CORS configuration must have at least one rule"}
	}

	// Convert to MinIO CORS config
	minioCORSConfig := &cors.Config{}
	for _, rule := range corsConfig.CORSRules {
		minioCORSConfig.CORSRules = append(minioCORSConfig.CORSRules, cors.Rule{
			AllowedOrigin: rule.AllowedOrigins,
			AllowedMethod: rule.AllowedMethods,
			AllowedHeader: rule.AllowedHeaders,
			ExposeHeader:  rule.ExposeHeaders,
			MaxAgeSeconds: rule.MaxAgeSeconds,
		})
	}

	return s.client.minioClient.SetBucketCors(ctx, bucketName, minioCORSConfig)
}

// GetCORS retrieves the CORS configuration for a bucket.
func (s *bucketService) GetCORS(ctx context.Context, bucketName string) (*CORSConfiguration, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	minioCORSConfig, err := s.client.minioClient.GetBucketCors(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if minioCORSConfig == nil || len(minioCORSConfig.CORSRules) == 0 {
		return nil, nil
	}

	// Convert from MinIO CORS config
	corsRules := make([]CORSRule, len(minioCORSConfig.CORSRules))
	for i, rule := range minioCORSConfig.CORSRules {
		corsRules[i] = CORSRule{
			AllowedOrigins: rule.AllowedOrigin,
			AllowedMethods: rule.AllowedMethod,
			AllowedHeaders: rule.AllowedHeader,
			ExposeHeaders:  rule.ExposeHeader,
			MaxAgeSeconds:  rule.MaxAgeSeconds,
		}
	}

	return &CORSConfiguration{
		CORSRules: corsRules,
	}, nil
}

// DeleteCORS removes the CORS configuration from a bucket.
func (s *bucketService) DeleteCORS(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	// Set empty CORS config to delete
	return s.client.minioClient.SetBucketCors(ctx, bucketName, nil)
}

// EnableVersioning enables versioning for a bucket.
func (s *bucketService) EnableVersioning(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	return s.client.minioClient.EnableVersioning(ctx, bucketName)
}

// SuspendVersioning suspends versioning for a bucket.
func (s *bucketService) SuspendVersioning(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return &InvalidBucketNameError{Name: bucketName}
	}

	return s.client.minioClient.SuspendVersioning(ctx, bucketName)
}

// GetVersioningStatus retrieves the versioning status of a bucket.
func (s *bucketService) GetVersioningStatus(ctx context.Context, bucketName string) (*BucketVersioningConfiguration, error) {
	if bucketName == "" {
		return nil, &InvalidBucketNameError{Name: bucketName}
	}

	minioConfig, err := s.client.minioClient.GetBucketVersioning(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	config := &BucketVersioningConfiguration{
		Status: VersioningStatus(minioConfig.Status),
	}

	return config, nil
}
