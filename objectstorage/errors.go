package objectstorage

import "fmt"

// InvalidBucketNameError is returned when a bucket name is invalid or empty.
type InvalidBucketNameError struct {
	Name string
}

// Error returns a string representation of the error.
func (e *InvalidBucketNameError) Error() string {
	return fmt.Sprintf("invalid bucket name: %s", e.Name)
}

// InvalidObjectKeyError is returned when an object key is invalid or empty.
type InvalidObjectKeyError struct {
	Key string
}

// Error returns a string representation of the error.
func (e *InvalidObjectKeyError) Error() string {
	return fmt.Sprintf("invalid object key: %s", e.Key)
}

// InvalidObjectDataError is returned when object data is invalid.
type InvalidObjectDataError struct {
	Message string
}

// Error returns a string representation of the error.
func (e *InvalidObjectDataError) Error() string {
	return fmt.Sprintf("invalid object data: %s", e.Message)
}

// InvalidPolicyError is returned when a bucket policy is invalid.
type InvalidPolicyError struct {
	Message string
}

// Error returns a string representation of the error.
func (e *InvalidPolicyError) Error() string {
	return fmt.Sprintf("invalid policy: %s", e.Message)
}

// BucketError represents an error that occurred during a bucket operation.
type BucketError struct {
	Operation string
	Bucket    string
	Message   string
}

// Error returns a string representation of the error.
func (e *BucketError) Error() string {
	return fmt.Sprintf("bucket operation %s on %s failed: %s", e.Operation, e.Bucket, e.Message)
}

// ObjectError represents an error that occurred during an object operation.
type ObjectError struct {
	Operation string
	Bucket    string
	Key       string
	Message   string
}

// Error returns a string representation of the error.
func (e *ObjectError) Error() string {
	return fmt.Sprintf("object operation %s on %s/%s failed: %s", e.Operation, e.Bucket, e.Key, e.Message)
}
