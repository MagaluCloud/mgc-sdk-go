package objectstorage

import "time"

// Bucket represents an object storage bucket.
type Bucket struct {
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
}

// Object represents an object stored in a bucket.
type Object struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag,omitempty"`
	ContentType  string    `json:"content_type,omitempty"`
	StorageClass string    `json:"storage_class,omitempty"`
}

// BucketListOptions defines parameters for filtering and pagination of bucket lists.
type BucketListOptions struct {
	Limit  *int `json:"_limit,omitempty"`
	Offset *int `json:"_offset,omitempty"`
}

// ObjectListOptions defines parameters for filtering and pagination of object lists.
type ObjectListOptions struct {
	Limit     *int   `json:"_limit,omitempty"`
	Offset    *int   `json:"_offset,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Delimiter string `json:"delimiter,omitempty"`
}

// ObjectFilterOptions defines filtering options for ListAll (without pagination).
type ObjectFilterOptions struct {
	Prefix    string `json:"prefix,omitempty"`
	Delimiter string `json:"delimiter,omitempty"`
}

// Statement represents a single statement in an S3 bucket policy.
type Statement struct {
	Sid       string `json:"Sid,omitempty"`
	Effect    string `json:"Effect"`
	Principal any    `json:"Principal"`
	Action    any    `json:"Action"`
	Resource  any    `json:"Resource"`
}

// Policy represents an S3 bucket policy with version and statements.
type Policy struct {
	Version   string      `json:"Version"`
	Id        string      `json:"Id,omitempty"`
	Statement []Statement `json:"Statement"`
}

// CORSRule represents a single CORS rule for a bucket.
type CORSRule struct {
	AllowedOrigins []string `json:"AllowedOrigins"`
	AllowedMethods []string `json:"AllowedMethods"`
	AllowedHeaders []string `json:"AllowedHeaders"`
	ExposeHeaders  []string `json:"ExposeHeaders,omitempty"`
	MaxAgeSeconds  int      `json:"MaxAgeSeconds"`
}

// CORSConfiguration represents CORS configuration for a bucket.
type CORSConfiguration struct {
	CORSRules []CORSRule `json:"CORSRules"`
}

// VersioningStatus represents the status of bucket versioning.
type VersioningStatus string

const (
	// VersioningStatusEnabled means versioning is enabled for the bucket.
	VersioningStatusEnabled VersioningStatus = "Enabled"
	// VersioningStatusSuspended means versioning is suspended for the bucket.
	VersioningStatusSuspended VersioningStatus = "Suspended"
	// VersioningStatusOff means versioning is not enabled for the bucket.
	VersioningStatusOff VersioningStatus = ""
)

// BucketVersioningConfiguration represents the versioning configuration of a bucket.
type BucketVersioningConfiguration struct {
	Status VersioningStatus `json:"Status,omitempty"`
}

// ObjectVersion represents a version of an object in a versioned bucket.
type ObjectVersion struct {
	Key            string    `json:"key"`
	VersionID      string    `json:"version_id"`
	Size           int64     `json:"size"`
	LastModified   time.Time `json:"last_modified"`
	IsDeleteMarker bool      `json:"is_delete_marker"`
	ETag           string    `json:"etag,omitempty"`
}

// DownloadOptions defines optional parameters for downloading objects.
type DownloadOptions struct {
	VersionID string `json:"version_id,omitempty"`
}

// DownloadStreamOptions defines optional parameters for streaming object downloads.
type DownloadStreamOptions struct {
	VersionID string `json:"version_id,omitempty"`
}

// DeleteOptions defines optional parameters for deleting objects.
type DeleteOptions struct {
	VersionID string `json:"version_id,omitempty"`
}

// ListVersionsOptions defines parameters for listing object versions.
type ListVersionsOptions struct {
	Limit  *int `json:"_limit,omitempty"`
	Offset *int `json:"_offset,omitempty"`
}

type ObjectLockInfo struct {
	Mode            string     `json:"mode,omitempty"`
	RetainUntilDate *time.Time `json:"retain_until_date,omitempty"`
	Locked          bool       `json:"locked"`
}

type MetadataOptions struct {
	VersionID string `json:"version_id,omitempty"`
}

type GetPresignedURLOptions struct {
	Method          string         `json:"method,omitempty"`
	ExpiryInSeconds *time.Duration `json:"expiry_in_seconds,omitempty"`
}

type PresignedURL struct {
	URL string `json:"url"`
}
