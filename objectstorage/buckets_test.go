package objectstorage

import (
	"context"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestBucketServiceCreate_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.Create(context.Background(), "")

	if err == nil {
		t.Error("Create() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("Create() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceExists_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.Exists(context.Background(), "")

	if err == nil {
		t.Error("Exists() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("Exists() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceDelete_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.Delete(context.Background(), "")

	if err == nil {
		t.Error("Delete() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("Delete() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceGetPolicy_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetPolicy(context.Background(), "")

	if err == nil {
		t.Error("GetPolicy() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("GetPolicy() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceList(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.List(context.Background(), BucketListOptions{})

	if err == nil {
		t.Error("List() expected error due to no connection, got nil")
	}
}

func TestBucketServiceListWithOptions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	// Test with Limit
	limit := 10
	_, err := svc.List(context.Background(), BucketListOptions{
		Limit: &limit,
	})
	if err == nil {
		t.Error("List() with Limit expected error due to no connection, got nil")
	}

	// Test with Offset
	offset := 5
	_, err = svc.List(context.Background(), BucketListOptions{
		Offset: &offset,
	})
	if err == nil {
		t.Error("List() with Offset expected error due to no connection, got nil")
	}

	// Test with both Limit and Offset
	_, err = svc.List(context.Background(), BucketListOptions{
		Limit:  &limit,
		Offset: &offset,
	})
	if err == nil {
		t.Error("List() with Limit and Offset expected error due to no connection, got nil")
	}
}

func TestBucketServiceGetPolicy(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetPolicy(context.Background(), "test-bucket")

	if err == nil {
		t.Error("GetPolicy() expected error due to no connection, got nil")
	}
}

func TestBucketServiceSetPolicy_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	policy := &Policy{
		Version: "2012-10-17",
		Statement: []Statement{
			{
				Effect:   "Allow",
				Action:   "s3:GetObject",
				Resource: "arn:aws:s3:::bucket/*",
			},
		},
	}

	err := svc.SetPolicy(context.Background(), "", policy)

	if err == nil {
		t.Error("SetPolicy() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("SetPolicy() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceSetPolicy_NilPolicy(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.SetPolicy(context.Background(), "test-bucket", nil)

	if err == nil {
		t.Error("SetPolicy() expected error for nil policy, got nil")
	}

	if _, ok := err.(*InvalidPolicyError); !ok {
		t.Errorf("SetPolicy() expected InvalidPolicyError, got %T", err)
	}
}

func TestBucketServiceSetPolicy_EmptyStatements(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	policy := &Policy{
		Version:   "2012-10-17",
		Statement: []Statement{},
	}

	err := svc.SetPolicy(context.Background(), "test-bucket", policy)

	if err == nil {
		t.Error("SetPolicy() expected error for empty statements, got nil")
	}

	if _, ok := err.(*InvalidPolicyError); !ok {
		t.Errorf("SetPolicy() expected InvalidPolicyError, got %T", err)
	}
}

func TestBucketServiceDeletePolicy_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.DeletePolicy(context.Background(), "")

	if err == nil {
		t.Error("DeletePolicy() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("DeletePolicy() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceLockBucket_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.LockBucket(context.Background(), "")

	if err == nil {
		t.Error("LockBucket() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("LockBucket() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceUnlockBucket_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.UnlockBucket(context.Background(), "")

	if err == nil {
		t.Error("UnlockBucket() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("UnlockBucket() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceGetBucketLockStatus_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetBucketLockStatus(context.Background(), "")

	if err == nil {
		t.Error("GetBucketLockStatus() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("GetBucketLockStatus() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceGetBucketLockStatus(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetBucketLockStatus(context.Background(), "test-bucket")

	if err == nil {
		t.Error("GetBucketLockStatus() expected error due to no connection, got nil")
	}
}

func TestBucketServiceLockBucket_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.LockBucket(context.Background(), "test-bucket")

	if err == nil {
		t.Error("LockBucket() expected error due to no connection, got nil")
	}
}

func TestBucketServiceUnlockBucket_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.UnlockBucket(context.Background(), "test-bucket")

	if err == nil {
		t.Error("UnlockBucket() expected error due to no connection, got nil")
	}
}

func TestBucketServiceCreate_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.Create(context.Background(), "test-bucket")

	if err == nil {
		t.Error("Create() expected error due to no connection, got nil")
	}
}

func TestBucketServiceExists_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.Exists(context.Background(), "test-bucket")

	if err == nil {
		t.Error("Exists() expected error due to no connection, got nil")
	}
}

func TestBucketServiceDelete_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.Delete(context.Background(), "test-bucket")

	if err == nil {
		t.Error("Delete() expected error due to no connection, got nil")
	}
}

func TestBucketServiceDeletePolicy_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.DeletePolicy(context.Background(), "test-bucket")

	if err == nil {
		t.Error("DeletePolicy() expected error due to no connection, got nil")
	}
}

func TestBucketServiceSetPolicy_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	policy := &Policy{
		Version: "2012-10-17",
		Statement: []Statement{
			{
				Effect:   "Allow",
				Action:   []string{"s3:GetObject"},
				Resource: "arn:aws:s3:::test-bucket/*",
			},
		},
	}

	err := svc.SetPolicy(context.Background(), "test-bucket", policy)

	if err == nil {
		t.Error("SetPolicy() expected error due to no connection, got nil")
	}
}

func TestBucketServiceSetPolicy_UnmarshalablePolicy(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	// Create a policy with a channel (which can't be marshaled to JSON)
	policy := &Policy{
		Version: "2012-10-17",
		Statement: []Statement{
			{
				Effect:   "Allow",
				Action:   make(chan int), // channels can't be marshaled
				Resource: "arn:aws:s3:::test-bucket/*",
			},
		},
	}

	err := svc.SetPolicy(context.Background(), "test-bucket", policy)

	if err == nil {
		t.Error("SetPolicy() expected error for unmarshalable policy, got nil")
	}
}

func TestBucketServiceEnableVersioning_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.EnableVersioning(context.Background(), "test-bucket")

	if err == nil {
		t.Error("EnableVersioning() expected error due to no connection, got nil")
	}
}

func TestBucketServiceSuspendVersioning_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.SuspendVersioning(context.Background(), "test-bucket")

	if err == nil {
		t.Error("SuspendVersioning() expected error due to no connection, got nil")
	}
}

func TestBucketServiceDeleteCORS_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.DeleteCORS(context.Background(), "test-bucket")

	if err == nil {
		t.Error("DeleteCORS() expected error due to no connection, got nil")
	}
}

func TestBucketServiceSetCORS_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	corsConfig := &CORSConfiguration{
		CORSRules: []CORSRule{
			{
				AllowedOrigins: []string{"https://example.com"},
				AllowedMethods: []string{"GET", "PUT"},
				AllowedHeaders: []string{"*"},
			},
		},
	}

	err := svc.SetCORS(context.Background(), "test-bucket", corsConfig)

	if err == nil {
		t.Error("SetCORS() expected error due to no connection, got nil")
	}
}

func TestBucketServiceImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ BucketService = (*bucketService)(nil)
}

func TestMarshalPolicy(t *testing.T) {
	t.Parallel()

	t.Run("valid policy", func(t *testing.T) {
		policy := &Policy{
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Effect:   "Allow",
					Action:   []string{"s3:GetObject"},
					Resource: "arn:aws:s3:::bucket/*",
				},
			},
		}

		policyStr, err := marshalPolicy(policy)
		if err != nil {
			t.Errorf("marshalPolicy() error = %v", err)
		}

		if policyStr == "" {
			t.Error("marshalPolicy() returned empty string")
		}

		// Verify it's valid JSON by unmarshaling it back
		var result Policy
		err = unmarshalPolicy(policyStr, &result)
		if err != nil {
			t.Errorf("unmarshalPolicy() error = %v", err)
		}

		if result.Version != policy.Version {
			t.Errorf("unmarshalPolicy() version = %s, want %s", result.Version, policy.Version)
		}
	})

	t.Run("policy with multiple statements", func(t *testing.T) {
		policy := &Policy{
			Version: "2012-10-17",
			Id:      "PolicyId",
			Statement: []Statement{
				{
					Sid:       "Statement1",
					Effect:    "Allow",
					Action:    []string{"s3:GetObject", "s3:PutObject"},
					Resource:  "arn:aws:s3:::bucket/*",
					Principal: map[string]interface{}{"AWS": "*"},
				},
				{
					Sid:      "Statement2",
					Effect:   "Deny",
					Action:   []string{"s3:DeleteObject"},
					Resource: "arn:aws:s3:::bucket/protected/*",
				},
			},
		}

		policyStr, err := marshalPolicy(policy)
		if err != nil {
			t.Errorf("marshalPolicy() error = %v", err)
		}

		if policyStr == "" {
			t.Error("marshalPolicy() returned empty string")
		}

		var result Policy
		err = unmarshalPolicy(policyStr, &result)
		if err != nil {
			t.Errorf("unmarshalPolicy() error = %v", err)
		}

		if len(result.Statement) != 2 {
			t.Errorf("unmarshalPolicy() len(Statement) = %d, want 2", len(result.Statement))
		}
	})

	t.Run("policy with unmarshalable value", func(t *testing.T) {
		// Create a policy with a channel (which can't be marshaled to JSON)
		policy := &Policy{
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Effect:   "Allow",
					Action:   make(chan int), // channels can't be marshaled
					Resource: "arn:aws:s3:::bucket/*",
				},
			},
		}

		_, err := marshalPolicy(policy)
		if err == nil {
			t.Error("marshalPolicy() expected error for unmarshalable value, got nil")
		}
	})
}

func TestUnmarshalPolicy(t *testing.T) {
	t.Parallel()

	policyStr := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:GetObject"],"Resource":"arn:aws:s3:::bucket/*"}]}`

	var policy Policy
	err := unmarshalPolicy(policyStr, &policy)
	if err != nil {
		t.Errorf("unmarshalPolicy() error = %v", err)
	}

	if policy.Version != "2012-10-17" {
		t.Errorf("unmarshalPolicy() version = %s, want 2012-10-17", policy.Version)
	}

	if len(policy.Statement) != 1 {
		t.Errorf("unmarshalPolicy() len(Statement) = %d, want 1", len(policy.Statement))
	}
}

func TestUnmarshalPolicy_InvalidJSON(t *testing.T) {
	t.Parallel()

	invalidJSON := `{"Version":"2012-10-17","Statement":[invalid json}`

	var policy Policy
	err := unmarshalPolicy(invalidJSON, &policy)
	if err == nil {
		t.Error("unmarshalPolicy() expected error for invalid JSON, got nil")
	}
}

func TestBucketType(t *testing.T) {
	t.Parallel()

	bucket := Bucket{
		Name: "test-bucket",
	}

	if bucket.Name != "test-bucket" {
		t.Errorf("Bucket.Name expected 'test-bucket', got %q", bucket.Name)
	}

}

func TestBucketListOptions(t *testing.T) {
	t.Parallel()

	limit := 10
	offset := 5
	opts := BucketListOptions{
		Limit:  &limit,
		Offset: &offset,
	}

	if opts.Limit == nil || *opts.Limit != 10 {
		t.Errorf("BucketListOptions.Limit expected 10, got %v", opts.Limit)
	}

	if opts.Offset == nil || *opts.Offset != 5 {
		t.Errorf("BucketListOptions.Offset expected 5, got %v", opts.Offset)
	}
}

func TestPolicy(t *testing.T) {
	t.Parallel()

	policy := &Policy{
		Version: "2012-10-17",
		Id:      "policy-1",
		Statement: []Statement{
			{
				Sid:       "AllowGetObject",
				Effect:    "Allow",
				Principal: map[string]interface{}{"AWS": "*"},
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/*",
			},
		},
	}

	if policy.Version != "2012-10-17" {
		t.Errorf("Policy.Version expected '2012-10-17', got %q", policy.Version)
	}

	if policy.Id != "policy-1" {
		t.Errorf("Policy.Id expected 'policy-1', got %q", policy.Id)
	}

	if len(policy.Statement) != 1 {
		t.Errorf("Policy.Statement expected 1 statement, got %d", len(policy.Statement))
	}

	if policy.Statement[0].Effect != "Allow" {
		t.Errorf("Statement.Effect expected 'Allow', got %q", policy.Statement[0].Effect)
	}
}

func TestStatement(t *testing.T) {
	t.Parallel()

	statement := Statement{
		Sid:       "TestStatement",
		Effect:    "Deny",
		Principal: map[string]interface{}{"AWS": "arn:aws:iam::123456789012:user/testuser"},
		Action:    []string{"s3:DeleteObject", "s3:PutObject"},
		Resource:  "arn:aws:s3:::bucket/*",
	}

	if statement.Sid != "TestStatement" {
		t.Errorf("Statement.Sid expected 'TestStatement', got %q", statement.Sid)
	}

	if statement.Effect != "Deny" {
		t.Errorf("Statement.Effect expected 'Deny', got %q", statement.Effect)
	}
}

func TestBucketLockStatusBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		locked   bool
		expected bool
	}{
		{
			name:     "locked bucket",
			locked:   true,
			expected: true,
		},
		{
			name:     "unlocked bucket",
			locked:   false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.locked != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.locked)
			}
		})
	}
}

func TestBucketServiceSetCORS_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	corsConfig := &CORSConfiguration{
		CORSRules: []CORSRule{
			{
				AllowedOrigins: []string{"https://example.com"},
				AllowedMethods: []string{"GET"},
				AllowedHeaders: []string{"*"},
				MaxAgeSeconds:  3600,
			},
		},
	}

	err := svc.SetCORS(context.Background(), "", corsConfig)

	if err == nil {
		t.Error("SetCORS() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("SetCORS() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceSetCORS_NilConfig(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.SetCORS(context.Background(), "test-bucket", nil)

	if err == nil {
		t.Error("SetCORS() expected error for nil config, got nil")
	}

	if _, ok := err.(*InvalidPolicyError); !ok {
		t.Errorf("SetCORS() expected InvalidPolicyError, got %T", err)
	}
}

func TestBucketServiceSetCORS_EmptyRules(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	corsConfig := &CORSConfiguration{
		CORSRules: []CORSRule{},
	}

	err := svc.SetCORS(context.Background(), "test-bucket", corsConfig)

	if err == nil {
		t.Error("SetCORS() expected error for empty rules, got nil")
	}

	if _, ok := err.(*InvalidPolicyError); !ok {
		t.Errorf("SetCORS() expected InvalidPolicyError, got %T", err)
	}
}

func TestBucketServiceGetCORS_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetCORS(context.Background(), "")

	if err == nil {
		t.Error("GetCORS() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("GetCORS() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceGetCORS(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetCORS(context.Background(), "test-bucket")

	if err == nil {
		t.Error("GetCORS() expected error due to no connection, got nil")
	}
}

func TestBucketServiceDeleteCORS_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.DeleteCORS(context.Background(), "")

	if err == nil {
		t.Error("DeleteCORS() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("DeleteCORS() expected InvalidBucketNameError, got %T", err)
	}
}

func TestCORSRule(t *testing.T) {
	t.Parallel()

	rule := CORSRule{
		AllowedOrigins: []string{"https://example.com", "https://another.com"},
		AllowedMethods: []string{"GET", "PUT"},
		AllowedHeaders: []string{"*"},
		ExposeHeaders:  []string{"ETag", "x-amz-request-id"},
		MaxAgeSeconds:  3600,
	}

	if len(rule.AllowedOrigins) != 2 {
		t.Errorf("CORSRule.AllowedOrigins expected 2 origins, got %d", len(rule.AllowedOrigins))
	}

	if len(rule.AllowedMethods) != 2 {
		t.Errorf("CORSRule.AllowedMethods expected 2 methods, got %d", len(rule.AllowedMethods))
	}

	if rule.MaxAgeSeconds != 3600 {
		t.Errorf("CORSRule.MaxAgeSeconds expected 3600, got %d", rule.MaxAgeSeconds)
	}
}

func TestCORSConfiguration(t *testing.T) {
	t.Parallel()

	config := &CORSConfiguration{
		CORSRules: []CORSRule{
			{
				AllowedOrigins: []string{"https://example.com"},
				AllowedMethods: []string{"GET"},
				AllowedHeaders: []string{"*"},
				MaxAgeSeconds:  3600,
			},
			{
				AllowedOrigins: []string{"https://another.com"},
				AllowedMethods: []string{"POST", "PUT"},
				AllowedHeaders: []string{"Content-Type"},
				MaxAgeSeconds:  7200,
			},
		},
	}

	if len(config.CORSRules) != 2 {
		t.Errorf("CORSConfiguration.CORSRules expected 2 rules, got %d", len(config.CORSRules))
	}

	if config.CORSRules[0].AllowedOrigins[0] != "https://example.com" {
		t.Errorf("CORSConfiguration.CORSRules[0].AllowedOrigins mismatch")
	}
}

func TestBucketServiceEnableVersioning_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.EnableVersioning(context.Background(), "")

	if err == nil {
		t.Error("EnableVersioning() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("EnableVersioning() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceSuspendVersioning_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	err := svc.SuspendVersioning(context.Background(), "")

	if err == nil {
		t.Error("SuspendVersioning() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("SuspendVersioning() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceGetVersioningStatus_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetVersioningStatus(context.Background(), "")

	if err == nil {
		t.Error("GetVersioningStatus() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("GetVersioningStatus() expected InvalidBucketNameError, got %T", err)
	}
}

func TestBucketServiceGetVersioningStatus(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Buckets()

	_, err := svc.GetVersioningStatus(context.Background(), "test-bucket")

	if err == nil {
		t.Error("GetVersioningStatus() expected error due to no connection, got nil")
	}
}

func TestVersioningStatus(t *testing.T) {
	t.Parallel()

	if VersioningStatusEnabled != "Enabled" {
		t.Errorf("VersioningStatusEnabled expected 'Enabled', got %q", VersioningStatusEnabled)
	}

	if VersioningStatusSuspended != "Suspended" {
		t.Errorf("VersioningStatusSuspended expected 'Suspended', got %q", VersioningStatusSuspended)
	}

	if VersioningStatusOff != "" {
		t.Errorf("VersioningStatusOff expected empty string, got %q", VersioningStatusOff)
	}
}

func TestBucketVersioningConfiguration(t *testing.T) {
	t.Parallel()

	config := &BucketVersioningConfiguration{
		Status: VersioningStatusEnabled,
	}

	if config.Status != VersioningStatusEnabled {
		t.Errorf("BucketVersioningConfiguration.Status expected 'Enabled', got %q", config.Status)
	}
}

func TestObjectVersion(t *testing.T) {
	t.Parallel()

	version := ObjectVersion{
		Key:            "test-object",
		VersionID:      "version-123",
		Size:           1024,
		IsDeleteMarker: false,
		ETag:           "abc123",
	}

	if version.Key != "test-object" {
		t.Errorf("ObjectVersion.Key expected 'test-object', got %q", version.Key)
	}

	if version.VersionID != "version-123" {
		t.Errorf("ObjectVersion.VersionID expected 'version-123', got %q", version.VersionID)
	}

	if version.Size != 1024 {
		t.Errorf("ObjectVersion.Size expected 1024, got %d", version.Size)
	}

	if version.IsDeleteMarker {
		t.Error("ObjectVersion.IsDeleteMarker expected false, got true")
	}
}
