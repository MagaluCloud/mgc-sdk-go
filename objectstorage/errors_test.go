package objectstorage

import (
	"testing"
)

func TestInvalidBucketNameError(t *testing.T) {
	t.Parallel()

	err := &InvalidBucketNameError{Name: ""}
	expectedMsg := "invalid bucket name: "
	if err.Error() != expectedMsg {
		t.Errorf("InvalidBucketNameError.Error() expected %q, got %q", expectedMsg, err.Error())
	}

	err = &InvalidBucketNameError{Name: "my-bucket"}
	expectedMsg = "invalid bucket name: my-bucket"
	if err.Error() != expectedMsg {
		t.Errorf("InvalidBucketNameError.Error() expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestInvalidObjectKeyError(t *testing.T) {
	t.Parallel()

	err := &InvalidObjectKeyError{Key: ""}
	expectedMsg := "invalid object key: "
	if err.Error() != expectedMsg {
		t.Errorf("InvalidObjectKeyError.Error() expected %q, got %q", expectedMsg, err.Error())
	}

	err = &InvalidObjectKeyError{Key: "my-key"}
	expectedMsg = "invalid object key: my-key"
	if err.Error() != expectedMsg {
		t.Errorf("InvalidObjectKeyError.Error() expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestInvalidObjectDataError(t *testing.T) {
	t.Parallel()

	err := &InvalidObjectDataError{Message: "data is empty"}
	expectedMsg := "invalid object data: data is empty"
	if err.Error() != expectedMsg {
		t.Errorf("InvalidObjectDataError.Error() expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestInvalidPolicyError(t *testing.T) {
	t.Parallel()

	err := &InvalidPolicyError{Message: "invalid statement"}
	expectedMsg := "invalid policy: invalid statement"
	if err.Error() != expectedMsg {
		t.Errorf("InvalidPolicyError.Error() expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestBucketError(t *testing.T) {
	t.Parallel()

	err := &BucketError{
		Operation: "create",
		Bucket:    "test-bucket",
		Message:   "bucket already exists",
	}
	expectedMsg := "bucket operation create on test-bucket failed: bucket already exists"
	if err.Error() != expectedMsg {
		t.Errorf("BucketError.Error() expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestObjectError(t *testing.T) {
	t.Parallel()

	err := &ObjectError{
		Operation: "put",
		Bucket:    "test-bucket",
		Key:       "test-key",
		Message:   "access denied",
	}
	expectedMsg := "object operation put on test-bucket/test-key failed: access denied"
	if err.Error() != expectedMsg {
		t.Errorf("ObjectError.Error() expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestErrorImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ error = (*InvalidBucketNameError)(nil)
	var _ error = (*InvalidObjectKeyError)(nil)
	var _ error = (*InvalidObjectDataError)(nil)
	var _ error = (*BucketError)(nil)
	var _ error = (*ObjectError)(nil)
}
