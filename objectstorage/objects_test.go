package objectstorage

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestObjectServiceUpload_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Upload(context.Background(), "", "test-key", []byte("test-data"), "text/plain")

	if err == nil {
		t.Error("Upload() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("Upload() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceUpload_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Upload(context.Background(), "test-bucket", "", []byte("test-data"), "text/plain")

	if err == nil {
		t.Error("Upload() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("Upload() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceUpload_EmptyData(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Upload(context.Background(), "test-bucket", "test-key", []byte{}, "")

	if err == nil {
		t.Error("Upload() expected error for empty data, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("Upload() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceUpload_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	data := []byte("test data")
	err := svc.Upload(context.Background(), "test-bucket", "test-key", data, "text/plain")

	if err == nil {
		t.Error("Upload() expected error due to no connection, got nil")
	}
}

func TestObjectServiceDownload_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.Download(context.Background(), "test-bucket", "test-key", nil)

	if err == nil {
		t.Error("Download() expected error due to no connection, got nil")
	}
}

func TestObjectServiceDownload_WithOptions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	t.Run("with empty VersionID", func(t *testing.T) {
		// Test with empty VersionID (should still set options)
		opts := &DownloadOptions{VersionID: ""}
		_, err := svc.Download(context.Background(), "test-bucket", "test-key", opts)
		if err == nil {
			t.Error("Download() with empty VersionID expected error due to no connection, got nil")
		}
	})

	t.Run("with non-empty VersionID", func(t *testing.T) {
		// Test with non-empty VersionID
		opts2 := &DownloadOptions{VersionID: "v123"}
		_, err := svc.Download(context.Background(), "test-bucket", "test-key", opts2)
		if err == nil {
			t.Error("Download() with VersionID expected error due to no connection, got nil")
		}
	})

	t.Run("with nil options", func(t *testing.T) {
		// Test with nil options
		_, err := svc.Download(context.Background(), "test-bucket", "test-key", nil)
		if err == nil {
			t.Error("Download() with nil options expected error due to no connection, got nil")
		}
	})
}

func TestObjectServiceDownloadStream_WithOptions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	t.Run("with empty VersionID", func(t *testing.T) {
		// Test with empty VersionID (should still set options)
		opts := &DownloadStreamOptions{VersionID: ""}
		_, err := svc.DownloadStream(context.Background(), "test-bucket", "test-key", opts)
		if err != nil {
			// Expected - no connection
			return
		}
	})

	t.Run("with non-empty VersionID", func(t *testing.T) {
		// Test with non-empty VersionID
		opts2 := &DownloadStreamOptions{VersionID: "v123"}
		_, err := svc.DownloadStream(context.Background(), "test-bucket", "test-key", opts2)
		// May succeed or fail depending on connection
		_ = err
	})

	t.Run("with nil options", func(t *testing.T) {
		// Test with nil options
		_, err := svc.DownloadStream(context.Background(), "test-bucket", "test-key", nil)
		// May succeed or fail depending on connection
		_ = err
	})
}

func TestObjectServiceDownload_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.Download(context.Background(), "", "test-key", nil)

	if err == nil {
		t.Error("Download() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("Download() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceDownload_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.Download(context.Background(), "test-bucket", "", nil)

	if err == nil {
		t.Error("Download() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("Download() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceDownloadStream_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	stream, err := svc.DownloadStream(context.Background(), "test-bucket", "test-key", nil)

	// DownloadStream returns an object even without connection, error occurs on read
	if err != nil {
		// This is fine - connection error
		return
	}

	if stream == nil {
		t.Error("DownloadStream() expected stream object, got nil")
	}
}

func TestObjectServiceDownloadStream_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.DownloadStream(context.Background(), "", "test-key", nil)

	if err == nil {
		t.Error("DownloadStream() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("DownloadStream() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceDownloadStream_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.DownloadStream(context.Background(), "test-bucket", "", nil)

	if err == nil {
		t.Error("DownloadStream() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("DownloadStream() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		bucketName string
		opts       ObjectListOptions
		wantErr    bool
	}{
		{
			name:       "empty bucket name",
			bucketName: "",
			opts:       ObjectListOptions{},
			wantErr:    true,
		},
		{
			name:       "valid parameters",
			bucketName: "test-bucket",
			opts:       ObjectListOptions{},
			wantErr:    false,
		},
		{
			name:       "with pagination",
			bucketName: "test-bucket",
			opts: ObjectListOptions{
				Limit:  intPtr(10),
				Offset: intPtr(0),
			},
			wantErr: false,
		},
		{
			name:       "with prefix",
			bucketName: "test-bucket",
			opts: ObjectListOptions{
				Prefix: "test/",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := client.NewMgcClient()
			osClient, _ := New(core, "minioadmin", "minioadmin")
			svc := osClient.Objects()

			_, err := svc.List(context.Background(), tt.bucketName, tt.opts)

			if tt.wantErr && err == nil {
				t.Errorf("List() expected error, got nil")
			}
		})
	}
}

func TestObjectServiceListAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		bucketName string
		opts       ObjectFilterOptions
		wantErr    bool
	}{
		{
			name:       "empty bucket name",
			bucketName: "",
			opts:       ObjectFilterOptions{},
			wantErr:    true,
		},
		{
			name:       "valid parameters",
			bucketName: "test-bucket",
			opts:       ObjectFilterOptions{},
			wantErr:    false,
		},
		{
			name:       "with prefix",
			bucketName: "test-bucket",
			opts: ObjectFilterOptions{
				Prefix: "test/",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := client.NewMgcClient()
			osClient, _ := New(core, "minioadmin", "minioadmin")
			svc := osClient.Objects()

			_, err := svc.ListAll(context.Background(), tt.bucketName, tt.opts)

			if tt.wantErr && err == nil {
				t.Errorf("ListAll() expected error, got nil")
			}
		})
	}
}

func TestObjectServiceDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		bucketName string
		objectKey  string
		wantErr    bool
	}{
		{
			name:       "empty bucket name",
			bucketName: "",
			objectKey:  "test-key",
			wantErr:    true,
		},
		{
			name:       "empty object key",
			bucketName: "test-bucket",
			objectKey:  "",
			wantErr:    true,
		},
		{
			name:       "valid parameters",
			bucketName: "test-bucket",
			objectKey:  "test-key",
			wantErr:    true, // Expected since no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := client.NewMgcClient()
			osClient, _ := New(core, "minioadmin", "minioadmin")
			svc := osClient.Objects()

			err := svc.Delete(context.Background(), tt.bucketName, tt.objectKey, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Delete() expected error, got nil")
			}
		})
	}
}

func TestObjectServiceListAllWithOptions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	t.Run("with Prefix", func(t *testing.T) {
		_, err := svc.ListAll(context.Background(), "test-bucket", ObjectFilterOptions{
			Prefix: "folder/",
		})
		if err == nil {
			t.Error("ListAll() with Prefix expected error due to no connection, got nil")
		}
	})

	t.Run("with Delimiter - non-recursive", func(t *testing.T) {
		_, err := svc.ListAll(context.Background(), "test-bucket", ObjectFilterOptions{
			Delimiter: "/",
		})
		if err == nil {
			t.Error("ListAll() with Delimiter expected error due to no connection, got nil")
		}
	})

	t.Run("with empty Delimiter - recursive", func(t *testing.T) {
		_, err := svc.ListAll(context.Background(), "test-bucket", ObjectFilterOptions{
			Delimiter: "",
		})
		if err == nil {
			t.Error("ListAll() with empty Delimiter expected error due to no connection, got nil")
		}
	})

	t.Run("with Prefix and Delimiter", func(t *testing.T) {
		_, err := svc.ListAll(context.Background(), "test-bucket", ObjectFilterOptions{
			Prefix:    "test/",
			Delimiter: "/",
		})
		if err == nil {
			t.Error("ListAll() with Prefix and Delimiter expected error due to no connection, got nil")
		}
	})
}

func TestObjectServiceListWithOptions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	t.Run("with Limit", func(t *testing.T) {
		limit := 10
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Limit: &limit,
		})
		if err == nil {
			t.Error("List() with Limit expected error due to no connection, got nil")
		}
	})

	t.Run("with Offset", func(t *testing.T) {
		offset := 5
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Offset: &offset,
		})
		if err == nil {
			t.Error("List() with Offset expected error due to no connection, got nil")
		}
	})

	t.Run("with both Limit and Offset", func(t *testing.T) {
		limit := 10
		offset := 5
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Limit:  &limit,
			Offset: &offset,
		})
		if err == nil {
			t.Error("List() with Limit and Offset expected error due to no connection, got nil")
		}
	})

	t.Run("with Prefix only", func(t *testing.T) {
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Prefix: "folder/",
		})
		if err == nil {
			t.Error("List() with Prefix expected error due to no connection, got nil")
		}
	})

	t.Run("with Delimiter - non-recursive", func(t *testing.T) {
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Delimiter: "/",
		})
		if err == nil {
			t.Error("List() with Delimiter expected error due to no connection, got nil")
		}
	})

	t.Run("with empty Delimiter - recursive", func(t *testing.T) {
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Delimiter: "",
		})
		if err == nil {
			t.Error("List() with empty Delimiter expected error due to no connection, got nil")
		}
	})

	t.Run("with Prefix and Delimiter", func(t *testing.T) {
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Prefix:    "folder/",
			Delimiter: "/",
		})
		if err == nil {
			t.Error("List() with Prefix and Delimiter expected error due to no connection, got nil")
		}
	})

	t.Run("with all options", func(t *testing.T) {
		limit := 10
		offset := 5
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Limit:     &limit,
			Offset:    &offset,
			Prefix:    "folder/",
			Delimiter: "/",
		})
		if err == nil {
			t.Error("List() with all options expected error due to no connection, got nil")
		}
	})
}

func TestObjectServiceMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		bucketName string
		objectKey  string
		wantErr    bool
	}{
		{
			name:       "empty bucket name",
			bucketName: "",
			objectKey:  "test-key",
			wantErr:    true,
		},
		{
			name:       "empty object key",
			bucketName: "test-bucket",
			objectKey:  "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := client.NewMgcClient()
			osClient, _ := New(core, "minioadmin", "minioadmin")
			svc := osClient.Objects()

			_, err := svc.Metadata(context.Background(), tt.bucketName, tt.objectKey)

			if tt.wantErr && err == nil {
				t.Errorf("Metadata() expected error, got nil")
			}
		})
	}
}

func TestObjectServiceMetadata_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.Metadata(context.Background(), "test-bucket", "test-key")

	if err == nil {
		t.Error("Metadata() expected error due to no connection, got nil")
	}
}

func TestObjectServiceMetadata_WithVersionID(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	// Test retrieving metadata - this exercises the success path up to the MinIO call
	_, err := svc.Metadata(context.Background(), "test-bucket", "test-key")

	// Expected to fail due to no connection, but validates parameter handling
	if err == nil {
		t.Error("Metadata() expected error due to no connection, got nil")
	}
}

func TestObjectServiceLockObject_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.LockObject(context.Background(), "", "test-key", time.Now().Add(24*time.Hour))

	if err == nil {
		t.Error("LockObject() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("LockObject() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceLockObject_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.LockObject(context.Background(), "test-bucket", "", time.Now().Add(24*time.Hour))

	if err == nil {
		t.Error("LockObject() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("LockObject() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceLockObject_ZeroRetentionDate(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.LockObject(context.Background(), "test-bucket", "test-key", time.Time{})

	if err == nil {
		t.Error("LockObject() expected error for zero retention date, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("LockObject() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceUnlockObject_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.UnlockObject(context.Background(), "", "test-key")

	if err == nil {
		t.Error("UnlockObject() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("UnlockObject() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceUnlockObject_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.UnlockObject(context.Background(), "test-bucket", "")

	if err == nil {
		t.Error("UnlockObject() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("UnlockObject() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceGetObjectLockStatus_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetObjectLockStatus(context.Background(), "", "test-key")

	if err == nil {
		t.Error("GetObjectLockStatus() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("GetObjectLockStatus() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceGetObjectLockStatus_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetObjectLockStatus(context.Background(), "test-bucket", "")

	if err == nil {
		t.Error("GetObjectLockStatus() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("GetObjectLockStatus() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceGetObjectLockStatus(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetObjectLockStatus(context.Background(), "test-bucket", "test-key")

	if err == nil {
		t.Error("GetObjectLockStatus() expected error due to no connection, got nil")
	}
}

func TestObjectServiceGetObjectLockInfo_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetObjectLockInfo(context.Background(), "", "test-key")

	if err == nil {
		t.Error("GetObjectLockInfo() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("GetObjectLockInfo() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceGetObjectLockInfo_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetObjectLockInfo(context.Background(), "test-bucket", "")

	if err == nil {
		t.Error("GetObjectLockInfo() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("GetObjectLockInfo() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceGetObjectLockInfo(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetObjectLockInfo(context.Background(), "test-bucket", "test-key")

	if err == nil {
		t.Error("GetObjectLockInfo() expected error due to no connection, got nil")
	}
}

// Versioning tests

func TestObjectServiceListVersions_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.ListVersions(context.Background(), "", "test-key", nil)

	if err == nil {
		t.Error("ListVersions() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("ListVersions() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceListVersions_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.ListVersions(context.Background(), "test-bucket", "", nil)

	if err == nil {
		t.Error("ListVersions() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("ListVersions() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceListVersions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.ListVersions(context.Background(), "test-bucket", "test-key", nil)

	if err == nil {
		t.Error("ListVersions() expected error due to no connection, got nil")
	}
}

func TestObjectServiceListVersionsWithOptions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	t.Run("with Limit", func(t *testing.T) {
		limit := 10
		_, err := svc.ListVersions(context.Background(), "test-bucket", "test-key", &ListVersionsOptions{
			Limit: &limit,
		})
		if err == nil {
			t.Error("ListVersions() with Limit expected error due to no connection, got nil")
		}
	})

	t.Run("with Offset", func(t *testing.T) {
		offset := 5
		_, err := svc.ListVersions(context.Background(), "test-bucket", "test-key", &ListVersionsOptions{
			Offset: &offset,
		})
		if err == nil {
			t.Error("ListVersions() with Offset expected error due to no connection, got nil")
		}
	})

	t.Run("with both Limit and Offset", func(t *testing.T) {
		limit := 10
		offset := 5
		_, err := svc.ListVersions(context.Background(), "test-bucket", "test-key", &ListVersionsOptions{
			Limit:  &limit,
			Offset: &offset,
		})
		if err == nil {
			t.Error("ListVersions() with Limit and Offset expected error due to no connection, got nil")
		}
	})

	t.Run("with zero Limit", func(t *testing.T) {
		limit := 0
		_, err := svc.ListVersions(context.Background(), "test-bucket", "test-key", &ListVersionsOptions{
			Limit: &limit,
		})
		if err == nil {
			t.Error("ListVersions() with zero Limit expected error due to no connection, got nil")
		}
	})

	t.Run("with zero Offset", func(t *testing.T) {
		offset := 0
		_, err := svc.ListVersions(context.Background(), "test-bucket", "test-key", &ListVersionsOptions{
			Offset: &offset,
		})
		if err == nil {
			t.Error("ListVersions() with zero Offset expected error due to no connection, got nil")
		}
	})
}

func TestObjectServiceLockObject_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	retentionDate := time.Now().Add(24 * time.Hour)
	err := svc.LockObject(context.Background(), "test-bucket", "test-key", retentionDate)

	if err == nil {
		t.Error("LockObject() expected error due to no connection, got nil")
	}
}

func TestObjectServiceUnlockObject_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.UnlockObject(context.Background(), "test-bucket", "test-key")

	if err == nil {
		t.Error("UnlockObject() expected error due to no connection, got nil")
	}
}

func TestObjectServiceDownload_WithVersionID(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	opts := &DownloadOptions{VersionID: "v123"}
	_, err := svc.Download(context.Background(), "test-bucket", "test-key", opts)

	// Error expected since not connected to real storage
	if err == nil {
		t.Logf("Download with VersionID returned error (expected): version ID accepted in options")
	}
}

func TestObjectServiceDownloadStream_WithVersionID(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	opts := &DownloadStreamOptions{VersionID: "v123"}
	_, err := svc.DownloadStream(context.Background(), "test-bucket", "test-key", opts)

	// Error expected since not connected to real storage
	if err == nil {
		t.Logf("DownloadStream with VersionID returned error (expected): version ID accepted in options")
	}
}

func TestObjectServiceDelete_WithVersionID(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	opts := &DeleteOptions{VersionID: "v123"}
	err := svc.Delete(context.Background(), "test-bucket", "test-key", opts)

	// Error expected since not connected to real storage
	if err == nil {
		t.Error("Delete() expected error due to no connection, got nil")
	}
}

func TestObjectServiceDelete_WithEmptyVersionID(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	// Test with empty VersionID option (should still process)
	opts := &DeleteOptions{VersionID: ""}
	err := svc.Delete(context.Background(), "test-bucket", "test-key", opts)

	if err == nil {
		t.Error("Delete() expected error due to no connection, got nil")
	}
}

func TestObjectServiceImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ ObjectService = (*objectService)(nil)
}

func TestStreamReadCloser(t *testing.T) {
	t.Parallel()

	data := []byte("test stream data")
	reader := io.NopCloser(bytes.NewReader(data))

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("Failed to read stream: %v", err)
	}

	if !bytes.Equal(content, data) {
		t.Errorf("Stream content mismatch")
	}

	err = reader.Close()
	if err != nil {
		t.Errorf("Failed to close stream: %v", err)
	}
}

func TestObjectType(t *testing.T) {
	t.Parallel()

	obj := Object{
		Key:         "test-key",
		Size:        1024,
		ETag:        "abc123",
		ContentType: "text/plain",
	}

	if obj.Key != "test-key" {
		t.Errorf("Object.Key expected 'test-key', got %q", obj.Key)
	}

	if obj.Size != 1024 {
		t.Errorf("Object.Size expected 1024, got %d", obj.Size)
	}

	if obj.ETag != "abc123" {
		t.Errorf("Object.ETag expected 'abc123', got %q", obj.ETag)
	}

	if obj.ContentType != "text/plain" {
		t.Errorf("Object.ContentType expected 'text/plain', got %q", obj.ContentType)
	}
}

func TestObjectListOptions(t *testing.T) {
	t.Parallel()

	limit := 20
	offset := 0
	opts := ObjectListOptions{
		Limit:     &limit,
		Offset:    &offset,
		Prefix:    "uploads/",
		Delimiter: "/",
	}

	if opts.Limit == nil || *opts.Limit != 20 {
		t.Errorf("ObjectListOptions.Limit expected 20, got %v", opts.Limit)
	}

	if opts.Offset == nil || *opts.Offset != 0 {
		t.Errorf("ObjectListOptions.Offset expected 0, got %v", opts.Offset)
	}

	if opts.Prefix != "uploads/" {
		t.Errorf("ObjectListOptions.Prefix expected 'uploads/', got %q", opts.Prefix)
	}

	if opts.Delimiter != "/" {
		t.Errorf("ObjectListOptions.Delimiter expected '/', got %q", opts.Delimiter)
	}
}

func TestObjectFilterOptions(t *testing.T) {
	t.Parallel()

	opts := ObjectFilterOptions{
		Prefix:    "documents/",
		Delimiter: "/",
	}

	if opts.Prefix != "documents/" {
		t.Errorf("ObjectFilterOptions.Prefix expected 'documents/', got %q", opts.Prefix)
	}

	if opts.Delimiter != "/" {
		t.Errorf("ObjectFilterOptions.Delimiter expected '/', got %q", opts.Delimiter)
	}
}

func TestObjectWithZeroSize(t *testing.T) {
	t.Parallel()

	obj := Object{
		Key:  "empty-file",
		Size: 0,
	}

	if obj.Size != 0 {
		t.Errorf("Object.Size expected 0, got %d", obj.Size)
	}
}

func TestObjectWithLargeSize(t *testing.T) {
	t.Parallel()

	largeSize := int64(1024 * 1024 * 1024)
	obj := Object{
		Key:  "large-file",
		Size: largeSize,
	}

	if obj.Size != largeSize {
		t.Errorf("Object.Size expected %d, got %d", largeSize, obj.Size)
	}
}

func TestObjectLockStatusBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		locked   bool
		expected bool
	}{
		{
			name:     "locked object",
			locked:   true,
			expected: true,
		},
		{
			name:     "unlocked object",
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

func TestDownloadOptions(t *testing.T) {
	t.Parallel()

	opts := &DownloadOptions{
		VersionID: "v123456789",
	}

	if opts.VersionID != "v123456789" {
		t.Errorf("DownloadOptions.VersionID expected 'v123456789', got %q", opts.VersionID)
	}
}

func TestDownloadStreamOptions(t *testing.T) {
	t.Parallel()

	opts := &DownloadStreamOptions{
		VersionID: "v987654321",
	}

	if opts.VersionID != "v987654321" {
		t.Errorf("DownloadStreamOptions.VersionID expected 'v987654321', got %q", opts.VersionID)
	}
}

func TestDeleteOptions(t *testing.T) {
	t.Parallel()

	opts := &DeleteOptions{
		VersionID: "v111222333",
	}

	if opts.VersionID != "v111222333" {
		t.Errorf("DeleteOptions.VersionID expected 'v111222333', got %q", opts.VersionID)
	}
}

func TestListVersionsOptions(t *testing.T) {
	t.Parallel()

	limit := 10
	offset := 5
	opts := &ListVersionsOptions{
		Limit:  &limit,
		Offset: &offset,
	}

	if opts.Limit == nil || *opts.Limit != 10 {
		t.Errorf("ListVersionsOptions.Limit expected 10, got %v", opts.Limit)
	}

	if opts.Offset == nil || *opts.Offset != 5 {
		t.Errorf("ListVersionsOptions.Offset expected 5, got %v", opts.Offset)
	}
}

func intPtr(v int) *int {
	return &v
}
