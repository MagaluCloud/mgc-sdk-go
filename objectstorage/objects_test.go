package objectstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/minio/minio-go/v7"
)

func TestObjectServiceUpload_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Upload(context.Background(), "", "test-key", []byte("test-data"), "text/plain", helpers.StrPtr("standard"))

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

	err := svc.Upload(context.Background(), "test-bucket", "", []byte("test-data"), "text/plain", helpers.StrPtr("cold_instant"))

	if err == nil {
		t.Error("Upload() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("Upload() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceUpload_InvalidStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	data := []byte("test data")
	err := svc.Upload(context.Background(), "test-bucket", "test-key", data, "text/plain", helpers.StrPtr("test"))

	if err == nil {
		t.Error("Upload() expected error for invalid storage class, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("Upload() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceUpload_EmptyData(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Upload(context.Background(), "test-bucket", "test-key", []byte{}, "", helpers.StrPtr("standard"))

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
	err := svc.Upload(context.Background(), "test-bucket", "test-key", data, "text/plain", helpers.StrPtr("standard"))

	if err == nil {
		t.Error("Upload() expected error due to no connection, got nil")
	}
}

func TestObjectServiceUpload_ValidStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	data := []byte("test data")
	err := svc.Upload(context.Background(), "test-bucket", "test-key", data, "text/plain", helpers.StrPtr("cold_instant"))

	if err == nil {
		t.Error("Upload() expected error due to no connection, got nil")
	}
}

func TestObjectServiceUploadDir_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.UploadDir(context.Background(), "", "test-key", "src", &UploadDirOptions{})

	if err == nil {
		t.Error("UploadDir() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("UploadDir() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceUploadDir_EmptySrcDir(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.UploadDir(
		context.Background(),
		"bucket",
		"key",
		"",
		nil,
	)

	if err == nil {
		t.Error("UploadDir() expected error for empty src dir, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("UploadDir() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceUploadDir_InvalidStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.UploadDir(context.Background(), "bucket-name", "test-key", "src", &UploadDirOptions{
		StorageClass: "invalid",
	})

	if err == nil {
		t.Error("UploadDir() expected error for invalid storage class, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("UploadDir() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceUploadDir_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.UploadDir(context.Background(), "bucket-name", "test-key", "src", &UploadDirOptions{
		Shallow:      false,
		StorageClass: "standard",
		BatchSize:    100,
	})

	if err == nil {
		t.Error("UploadDir() expected error, got nil")
	}
}

func TestObjectServiceUploadDir_ValidStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.UploadDir(context.Background(), "", "test-key", "src", &UploadDirOptions{
		Shallow:      false,
		StorageClass: "cold_instant",
		BatchSize:    100,
	})

	if err == nil {
		t.Error("UploadDir() expected error due to no connection, got nil")
	}
}

func TestObjectServiceUploadDir_BatchSizeZero(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.UploadDir(context.Background(), "bucket", "key", "src", &UploadDirOptions{
		BatchSize: 0,
	})

	if err == nil {
		t.Error("expected error for batch size zero")
	}
}

func TestObjectServiceUploadDir_WalkDirCollectsFiles(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "file.txt")
	err := os.WriteFile(filePath, []byte("data"), 0644)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = svc.UploadDir(
		context.Background(),
		"bucket",
		"key",
		tmpDir,
		nil,
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}
}

func TestObjectServiceUploadDir_ShallowSkipsSubDirs(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	tmpDir := t.TempDir()

	subDir := filepath.Join(tmpDir, "sub")

	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("failed to create subDir: %v", err)
	}

	err = os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("x"), 0644)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	result, err := svc.UploadDir(
		context.Background(),
		"bucket",
		"key",
		tmpDir,
		&UploadDirOptions{
			Shallow: true,
		},
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}

	if result.UploadedCount != 0 {
		t.Errorf("expected UploadedCount to be %v, got %v", 0, result.UploadedCount)
	}
}

func TestObjectServiceUploadDir_FilterSkipsFile(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("x"), 0644)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	filters := []FilterOptions{
		{Include: "no-match"},
	}

	result, err := svc.UploadDir(
		context.Background(),
		"bucket",
		"key",
		tmpDir,
		&UploadDirOptions{
			Filter: &filters,
		},
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}
	if result.UploadedCount != 0 {
		t.Errorf("expected UploadedCount to be %v, got %v", 0, result.UploadedCount)
	}
	if result.ErrorCount != 0 {
		t.Errorf("expected ErrorCount to be %v, got %v", 0, result.ErrorCount)
	}
}

func TestObjectServiceUploadDir_WithoutFilter(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("x"), 0644)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	result, err := svc.UploadDir(
		context.Background(),
		"bucket",
		"key",
		tmpDir,
		nil,
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}
	if result.ErrorCount != 1 {
		t.Errorf("expected ErrorCount to be %v, got %v", 1, result.ErrorCount)
	}
}

func TestObjectServiceUploadDir_WithStorageClass(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	srcDir := t.TempDir()

	filePath := filepath.Join(srcDir, "file.txt")
	err := os.WriteFile(filePath, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	mock := newMockMinioClient()

	var receivedStorageClass string

	mock.putObjectFunc = func(
		ctx context.Context,
		bucketName string,
		objectName string,
		reader io.Reader,
		objectSize int64,
		opts minio.PutObjectOptions,
	) (minio.UploadInfo, error) {
		receivedStorageClass = opts.StorageClass
		return minio.UploadInfo{}, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	storageClass := "STANDARD"

	res, err := svc.UploadDir(
		ctx,
		"bucket-name",
		"dst",
		srcDir,
		&UploadDirOptions{
			StorageClass: storageClass,
		},
	)

	if err != nil {
		t.Fatalf("UploadDir() unexpected error: %v", err)
	}

	if res.UploadedCount != 1 {
		t.Fatalf("expected UploadedCount 1, got %v", res.UploadedCount)
	}

	if receivedStorageClass != storageClass {
		t.Fatalf(
			"expected StorageClass %q, got %q",
			storageClass,
			receivedStorageClass,
		)
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

func TestObjectServiceDownloadAll_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.DownloadAll(context.Background(), "", "test-key", nil)

	if err == nil {
		t.Error("DownloadAll() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("DownloadAll() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceDownloadAll_InvalidDstPath(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.DownloadAll(context.Background(), "test-bucket", "", nil)

	if err == nil {
		t.Error("DownloadAll() expected error for empty dst path, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("DownloadAll() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceDownloadAll_ValidParameters(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.DownloadAll(context.Background(), "test-bucket", "test-key", nil)

	if err == nil {
		t.Error("DownloadAll() expected error, got nil")
	}
}

func TestObjectServiceDownloadAll_WithOptions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	t.Run("with empty Filter", func(t *testing.T) {
		opts := &DownloadAllOptions{Filter: nil}
		_, err := svc.DownloadAll(context.Background(), "test-bucket", "test-key", opts)
		if err == nil {
			t.Error("DownloadAll() expected error for empty data, got nil")
		}
	})

	t.Run("with non-empty Filter", func(t *testing.T) {
		opts := &DownloadAllOptions{Filter: &[]FilterOptions{{Exclude: "test", Include: "new"}}}
		_, err := svc.DownloadAll(context.Background(), "test-bucket", "test-key", opts)
		if err == nil {
			t.Error("DownloadAll() expected error for empty data, got nil")
		}
	})

	t.Run("with empty Prefix", func(t *testing.T) {
		opts := &DownloadAllOptions{Prefix: ""}
		_, err := svc.DownloadAll(context.Background(), "test-bucket", "test-key", opts)
		if err == nil {
			t.Error("DownloadAll() expected error for empty data, got nil")
		}
	})

	t.Run("with non-empty Prefix", func(t *testing.T) {
		opts := &DownloadAllOptions{Prefix: "test"}
		_, err := svc.DownloadAll(context.Background(), "test-bucket", "test-key", opts)
		if err == nil {
			t.Error("DownloadAll() expected error for empty data, got nil")
		}
	})
}

func TestObjectServiceDownloadAll_SuccessEmpty(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tmpDir := t.TempDir()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucket string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		close(ch)
		return ch
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	res, err := svc.DownloadAll(ctx, "bucket", tmpDir, nil)

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.DownloadedCount != 0 {
		t.Errorf("expected DownloadedCount 0, got %d", res.DownloadedCount)
	}
	if res.ErrorCount != 0 {
		t.Errorf("expected ErrorCount 0, got %d", res.ErrorCount)
	}
}

func TestObjectServiceDownloadAll_HandlerError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tmpDir := t.TempDir()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucket string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo, 1)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{
				Key:  "file.txt",
				Size: 10,
			}
		}()
		return ch
	}

	mock.getObjectFunc = func(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.GetObjectOptions,
	) (*minio.Object, error) {
		return nil, fmt.Errorf("boom")
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	res, err := svc.DownloadAll(ctx, "bucket", tmpDir, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.DownloadedCount != 0 {
		t.Errorf("expected DownloadedCount 0, got %d", res.DownloadedCount)
	}
	if res.ErrorCount != 1 {
		t.Errorf("expected ErrorCount 1, got %d", res.ErrorCount)
	}
	if len(res.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(res.Errors))
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
		{
			name:       "with filter",
			bucketName: "test-bucket",
			opts: ObjectListOptions{
				Filter: &[]FilterOptions{
					{Include: "images"},
				},
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

func TestObjectServiceList_WithOffset(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(ctx context.Context, bucket string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "a.txt"}
			ch <- minio.ObjectInfo{Key: "b.txt"}
			ch <- minio.ObjectInfo{Key: "c.txt"}
		}()
		return ch
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	offset := 1

	objs, err := svc.List(ctx, "bucket", ObjectListOptions{
		Offset: &offset,
	})

	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	gotKeys := make([]string, 0)
	for _, o := range objs {
		gotKeys = append(gotKeys, o.Key)
	}

	expected := []string{"b.txt", "c.txt"}

	if !reflect.DeepEqual(gotKeys, expected) {
		t.Errorf("expected %v, got %v", expected, gotKeys)
	}
}

func TestObjectServiceList_WithLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(ctx context.Context, bucket string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "a.txt"}
			ch <- minio.ObjectInfo{Key: "b.txt"}
			ch <- minio.ObjectInfo{Key: "c.txt"}
		}()
		return ch
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	limit := 2

	objs, err := svc.List(ctx, "bucket", ObjectListOptions{
		Limit: &limit,
	})

	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	gotKeys := make([]string, 0)
	for _, o := range objs {
		gotKeys = append(gotKeys, o.Key)
	}

	expected := []string{"a.txt", "b.txt"}

	if !reflect.DeepEqual(gotKeys, expected) {
		t.Errorf("expected %v, got %v", expected, gotKeys)
	}
}

func TestObjectServiceList_WithOffsetAndLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(ctx context.Context, bucket string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "a.txt"}
			ch <- minio.ObjectInfo{Key: "b.txt"}
			ch <- minio.ObjectInfo{Key: "c.txt"}
			ch <- minio.ObjectInfo{Key: "d.txt"}
		}()
		return ch
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	offset := 1
	limit := 2

	objs, err := svc.List(ctx, "bucket", ObjectListOptions{
		Offset: &offset,
		Limit:  &limit,
	})

	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	gotKeys := make([]string, 0)
	for _, o := range objs {
		gotKeys = append(gotKeys, o.Key)
	}

	expected := []string{"b.txt", "c.txt"}

	if !reflect.DeepEqual(gotKeys, expected) {
		t.Errorf("expected %v, got %v", expected, gotKeys)
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

func TestObjectServiceListAll_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	lastModified := time.Now()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucket string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "a.txt", Size: 10, LastModified: lastModified, ETag: "etag-a"}
			ch <- minio.ObjectInfo{Key: "b.txt", Size: 20, LastModified: lastModified.Add(time.Minute), ETag: "etag-b"}
		}()
		return ch
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	objects, err := svc.ListAll(ctx, "bucket-name", ObjectFilterOptions{})

	if err != nil {
		t.Errorf("ListAll() expected success, got %v", err)
	}

	if len(objects) != 2 {
		t.Errorf("ListAll() expected 2 objects, got %v", len(objects))
	}

	expected := []Object{
		{
			Key:          "a.txt",
			Size:         10,
			LastModified: lastModified,
			ETag:         "etag-a",
		},
		{
			Key:          "b.txt",
			Size:         20,
			LastModified: lastModified.Add(time.Minute),
			ETag:         "etag-b",
		},
	}

	if !reflect.DeepEqual(objects, expected) {
		t.Errorf("ListAll() mismatch\nexpected: %+v\ngot: %+v", expected, objects)
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

	t.Run("with Filter", func(t *testing.T) {
		_, err := svc.List(context.Background(), "test-bucket", ObjectListOptions{
			Filter: &[]FilterOptions{
				{Exclude: "folder/"},
			},
		})
		if err == nil {
			t.Error("List() with Filter expected error due to no connection, got nil")
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

			_, err := svc.Metadata(context.Background(), tt.bucketName, tt.objectKey, nil)

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

	_, err := svc.Metadata(context.Background(), "test-bucket", "test-key", nil)

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
	_, err := svc.Metadata(context.Background(), "test-bucket", "test-key", &MetadataOptions{
		VersionID: "version-id",
	})

	// Expected to fail due to no connection, but validates parameter handling
	if err == nil {
		t.Error("Metadata() expected error due to no connection, got nil")
	}
}

func TestObjectServiceMetadata_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lastModified := time.Now()

	mock := newMockMinioClient()

	mock.statObjectFunc = func(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.StatObjectOptions,
	) (minio.ObjectInfo, error) {
		return minio.ObjectInfo{
			Key:          objectName,
			Size:         123,
			LastModified: lastModified,
			ETag:         "etag-123",
			ContentType:  "text/plain",
			Metadata: map[string][]string{
				"X-Amz-Storage-Class": {"STANDARD"},
			},
		}, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	obj, err := svc.Metadata(ctx, "bucket-name", "file.txt", nil)

	if err != nil {
		t.Fatalf("Metadata() unexpected error: %v", err)
	}

	expected := &Object{
		Key:          "file.txt",
		Size:         123,
		LastModified: lastModified,
		ETag:         "etag-123",
		ContentType:  "text/plain",
		StorageClass: "STANDARD",
	}

	if !reflect.DeepEqual(obj, expected) {
		t.Fatalf("Metadata() mismatch\nexpected: %+v\ngot: %+v", expected, obj)
	}
}

func TestObjectServiceMetadata_NoStorageClass(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.statObjectFunc = func(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.StatObjectOptions,
	) (minio.ObjectInfo, error) {
		return minio.ObjectInfo{
			Key:         objectName,
			Size:        50,
			ContentType: "application/json",
			Metadata:    map[string][]string{},
		}, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	obj, err := svc.Metadata(ctx, "bucket-name", "data.json", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if obj.StorageClass != "" {
		t.Fatalf("expected empty StorageClass, got %q", obj.StorageClass)
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

func TestObjectServiceGetObjectLockStatus_Unlocked(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.getObjectRetentionFunc = func(
		ctx context.Context,
		bucketName string,
		objectKey string,
		versionID string,
	) (*minio.RetentionMode, *time.Time, error) {
		return nil, nil, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)

	isLocked, err := osClient.Objects().GetObjectLockStatus(
		ctx,
		"bucket-name",
		"file.txt",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if isLocked {
		t.Errorf("expected object to be unlocked, got locked")
	}
}

func TestObjectServiceGetObjectLockStatus_Locked(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.getObjectRetentionFunc = func(
		ctx context.Context,
		bucketName string,
		objectKey string,
		versionID string,
	) (*minio.RetentionMode, *time.Time, error) {
		mode := minio.RetentionMode("COMPLIANCE")
		until := time.Now().Add(time.Hour)

		return &mode, &until, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)

	isLocked, err := osClient.Objects().GetObjectLockStatus(
		ctx,
		"bucket-name",
		"file.txt",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !isLocked {
		t.Errorf("expected object to be locked, got unlocked")
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

func TestObjectServiceGetObjectLockInfo_Unlocked(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.getObjectRetentionFunc = func(
		ctx context.Context,
		bucketName string,
		objectKey string,
		versionID string,
	) (*minio.RetentionMode, *time.Time, error) {
		return nil, nil, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)

	info, err := osClient.Objects().GetObjectLockInfo(
		ctx,
		"bucket-name",
		"file.txt",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if info == nil {
		t.Fatalf("expected ObjectLockInfo, got nil")
	}

	if info.Locked {
		t.Errorf("expected Locked=false, got true")
	}

	if info.Mode != "" {
		t.Errorf("expected empty Mode, got %v", info.Mode)
	}

	if info.RetainUntilDate != nil {
		t.Errorf("expected RetainUntilDate to be nil, got %v", info.RetainUntilDate)
	}
}

func TestObjectServiceGetObjectLockInfo_Locked(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	retainUntil := time.Now().Add(24 * time.Hour)
	mode := minio.RetentionMode("COMPLIANCE")

	mock.getObjectRetentionFunc = func(
		ctx context.Context,
		bucketName string,
		objectKey string,
		versionID string,
	) (*minio.RetentionMode, *time.Time, error) {
		return &mode, &retainUntil, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)

	info, err := osClient.Objects().GetObjectLockInfo(
		ctx,
		"bucket-name",
		"file.txt",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info == nil {
		t.Fatalf("expected ObjectLockInfo, got nil")
	}
	if !info.Locked {
		t.Errorf("expected Locked=true, got false")
	}
	if info.Mode != mode.String() {
		t.Errorf("expected Mode=%v, got %v", mode.String(), info.Mode)
	}
	if !info.RetainUntilDate.Equal(retainUntil) {
		t.Errorf(
			"expected RetainUntilDate=%v, got %v",
			retainUntil,
			*info.RetainUntilDate,
		)
	}
}

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

func TestObjectServiceListVersions_WithLimitAndOffset(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lastModified := time.Now()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucketName string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)

		go func() {
			defer close(ch)

			ch <- minio.ObjectInfo{
				Key:          "file.txt",
				VersionID:    "v1",
				Size:         10,
				LastModified: lastModified,
				ETag:         "etag-v1",
			}
			ch <- minio.ObjectInfo{
				Key:          "file.txt",
				VersionID:    "v2",
				Size:         20,
				LastModified: lastModified.Add(time.Minute),
				ETag:         "etag-v2",
			}
			ch <- minio.ObjectInfo{
				Key:          "file.txt",
				VersionID:    "v3",
				Size:         30,
				LastModified: lastModified.Add(2 * time.Minute),
				ETag:         "etag-v3",
			}
			ch <- minio.ObjectInfo{
				Key:          "other.txt",
				VersionID:    "x1",
				Size:         999,
				LastModified: lastModified,
				ETag:         "etag-x",
			}
		}()

		return ch
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	limit := 1
	offset := 1

	versions, err := svc.ListVersions(
		ctx,
		"bucket-name",
		"file.txt",
		&ListVersionsOptions{
			Limit:  &limit,
			Offset: &offset,
		},
	)

	if err != nil {
		t.Errorf("ListVersions() expected success, got %v", err)
	}

	expected := []ObjectVersion{
		{
			Key:          "file.txt",
			VersionID:    "v2",
			Size:         20,
			LastModified: lastModified.Add(time.Minute),
			ETag:         "etag-v2",
		},
	}

	if len(versions) != len(expected) {
		t.Fatalf("expected %v versions, got %v", len(expected), len(versions))
	}

	if !reflect.DeepEqual(versions, expected) {
		t.Errorf("ListVersions() mismatch\nexpected: %+v\ngot: %+v", expected, versions)
	}
}

func TestObjectServiceListAllVersions_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.ListAllVersions(context.Background(), "", "test-key")

	if err == nil {
		t.Error("ListAllVersions() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("ListAllVersions() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceListAllVersions_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.ListAllVersions(context.Background(), "test-bucket", "")

	if err == nil {
		t.Error("ListAllVersions() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("ListAllVersions() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceListAllVersions(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.ListAllVersions(context.Background(), "test-bucket", "test-key")

	if err == nil {
		t.Error("ListAllVersions() expected error due to no connection, got nil")
	}
}

func TestObjectServiceListAllVersions_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	lastModified := time.Now()

	mock := newMockMinioClient()

	mock.buckets["bucket-name"] = &mockBucket{
		name: "bucket-name",
		objects: map[string]*mockObject{
			"file.txt": {
				key:          "file.txt",
				size:         123,
				lastModified: lastModified,
				etag:         "etag-123",
			},
		},
	}

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucketName string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo, 2)

		go func() {
			defer close(ch)

			ch <- minio.ObjectInfo{
				Key:            "file.txt",
				VersionID:      "v1",
				Size:           123,
				LastModified:   lastModified,
				ETag:           "etag-v1",
				IsDeleteMarker: false,
				IsLatest:       false,
				StorageClass:   "STANDARD",
				Owner: minio.Owner{
					DisplayName: "user-1",
					ID:          "owner-1",
				},
			}

			ch <- minio.ObjectInfo{
				Key:            "file.txt",
				VersionID:      "v2",
				Size:           456,
				LastModified:   lastModified.Add(time.Minute),
				ETag:           "etag-v2",
				IsDeleteMarker: true,
				IsLatest:       true,
				StorageClass:   "COLD",
				Owner: minio.Owner{
					DisplayName: "user-2",
					ID:          "owner-2",
				},
			}
		}()

		return ch
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	versions, err := svc.ListAllVersions(
		ctx,
		"bucket-name",
		"file.txt",
	)

	if err != nil {
		t.Errorf("ListAllVersions() expected success, got %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %v", len(versions))
	}

	expected := []ObjectVersion{
		{
			Key:            "file.txt",
			VersionID:      "v1",
			Size:           123,
			LastModified:   lastModified,
			ETag:           "etag-v1",
			IsDeleteMarker: false,
			IsLatest:       false,
			Owner: ObjectOwner{
				DisplayName: "user-1",
				ID:          "owner-1",
			},
			StorageClass: "STANDARD",
		},
		{
			Key:            "file.txt",
			VersionID:      "v2",
			Size:           456,
			LastModified:   lastModified.Add(time.Minute),
			ETag:           "etag-v2",
			IsDeleteMarker: true,
			IsLatest:       true,
			Owner: ObjectOwner{
				DisplayName: "user-2",
				ID:          "owner-2",
			},
			StorageClass: "COLD",
		},
	}

	if !reflect.DeepEqual(versions, expected) {
		t.Errorf("ListAllVersions() mismatch\nexpected: %+v\ngot: %+v", expected, versions)
	}
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

func TestObjectServiceGetPresignedURL_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetPresignedURL(context.Background(), "", "test-key", GetPresignedURLOptions{
		Method: http.MethodGet,
	})

	if err == nil {
		t.Error("GetPresignedURL() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("GetPresignedURL() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceGetPresignedURL_InvalidObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetPresignedURL(context.Background(), "test-bucket", "", GetPresignedURLOptions{
		Method: http.MethodGet,
	})

	if err == nil {
		t.Error("GetPresignedURL() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("GetPresignedURL() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceGetPresignedURL_InvalidMethod(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetPresignedURL(context.Background(), "test-bucket", "test-object", GetPresignedURLOptions{
		Method: http.MethodPost,
	})

	if err == nil {
		t.Error("GetPresignedURL() expected error for invalid method, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("GetPresignedURL() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceGetPresignedURL_PUTMethod(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetPresignedURL(context.Background(), "test-bucket", "test-key", GetPresignedURLOptions{
		Method: http.MethodPut,
	})

	if err != nil {
		t.Error("GetPresignedURL() expected presigned URL, got nil")
	}
}

func TestObjectServiceGetPresignedURL_GETMethod(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.GetPresignedURL(context.Background(), "test-bucket", "test-key", GetPresignedURLOptions{
		Method: http.MethodGet,
	})

	if err != nil {
		t.Error("GetPresignedURL() expected presigned URL, got nil")
	}
}

func TestObjectServiceGetPresignedURL_WithExpiry(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	expire := 5 * time.Minute

	_, err := svc.GetPresignedURL(context.Background(), "test-bucket", "test-key", GetPresignedURLOptions{
		Method:          http.MethodPut,
		ExpiryInSeconds: &expire,
	})

	if err != nil {
		t.Error("GetPresignedURL() expected presigned URL, got nil")
	}
}

func TestObjectServiceCopy_InvalidSrcBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "",
		ObjectKey:  "object-key",
	}, CopyDstConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	})

	if err == nil {
		t.Error("Copy() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("Copy() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceCopy_InvalidSrcObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "bucket-name",
		ObjectKey:  "",
	}, CopyDstConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	})

	if err == nil {
		t.Error("Copy() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("Copy() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceCopy_InvalidDstBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
		VersionID:  "version-id",
	}, CopyDstConfig{
		BucketName:   "",
		ObjectKey:    "object-key",
		StorageClass: "standard",
	})

	if err == nil {
		t.Error("Copy() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("Copy() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceCopy_InvalidDstObjectKey(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyDstConfig{
		BucketName: "bucket-name",
		ObjectKey:  "",
	})

	if err == nil {
		t.Error("Copy() expected error for empty object key, got nil")
	}

	if _, ok := err.(*InvalidObjectKeyError); !ok {
		t.Errorf("Copy() expected InvalidObjectKeyError, got %T", err)
	}
}

func TestObjectServiceCopy_InvalidStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyDstConfig{
		BucketName:   "bucket-name",
		ObjectKey:    "object-key",
		StorageClass: "invalid-class",
	})

	if err == nil {
		t.Error("Copy() expected error for invalid storage class, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("Copy() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceCopy(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyDstConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	})

	if err == nil {
		t.Error("Copy() expected error due to no connection, got nil")
	}
}

func TestObjectServiceCopy_WithStandardStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyDstConfig{
		BucketName:   "bucket-name",
		ObjectKey:    "object-key",
		StorageClass: "standard",
	})

	if err == nil {
		t.Error("Copy() expected error due to no connection, got nil")
	}
}

func TestObjectServiceCopy_WithColdInstantStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	err := svc.Copy(context.Background(), CopySrcConfig{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyDstConfig{
		BucketName:   "bucket-name",
		ObjectKey:    "object-key",
		StorageClass: "cold_instant",
	})

	if err == nil {
		t.Error("Copy() expected error due to no connection, got nil")
	}
}

func TestObjectServiceCopyAll_InvalidSrcBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.CopyAll(context.Background(), CopyPath{
		BucketName: "",
		ObjectKey:  "object-key",
	}, CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, nil)

	if err == nil {
		t.Error("CopyAll() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("CopyAll() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceCopyAll_InvalidDstBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.CopyAll(context.Background(), CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyPath{
		BucketName: "",
		ObjectKey:  "object-key",
	}, nil)

	if err == nil {
		t.Error("CopyAll() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("CopyAll() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceCopyAll_InvalidStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.CopyAll(context.Background(), CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, &CopyAllOptions{StorageClass: "invalid-class"})

	if err == nil {
		t.Error("CopyAll() expected error for invalid storage class, got nil")
	}

	if _, ok := err.(*InvalidObjectDataError); !ok {
		t.Errorf("CopyAll() expected InvalidObjectDataError, got %T", err)
	}
}

func TestObjectServiceCopyAll(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.CopyAll(context.Background(), CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, nil)

	if err == nil {
		t.Error("CopyAll() expected error, got nil")
	}
}

func TestObjectServiceCopyAll_WithStandardStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.CopyAll(context.Background(), CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, &CopyAllOptions{
		StorageClass: "standard",
	})

	if err == nil {
		t.Error("CopyAll() expected error, got nil")
	}
}

func TestObjectServiceCopyAll_WithColdInstantStorageClass(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.CopyAll(context.Background(), CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, CopyPath{
		BucketName: "bucket-name",
		ObjectKey:  "object-key",
	}, &CopyAllOptions{
		StorageClass: "cold_instant",
	})

	if err == nil {
		t.Error("CopyAll() expected error due to no connection, got nil")
	}
}

func TestObjectServiceCopyAll_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()
	mock.buckets["src"] = &mockBucket{
		name: "src",
		objects: map[string]*mockObject{
			"a.txt": {key: "a.txt", size: 10},
			"b.txt": {key: "b.txt", size: 20},
		},
	}
	mock.buckets["dst"] = &mockBucket{
		name:    "dst",
		objects: map[string]*mockObject{},
	}

	var mu sync.Mutex
	copied := make([]string, 0)

	mock.copyObjectFunc = func(
		ctx context.Context,
		dst minio.CopyDestOptions,
		src minio.CopySrcOptions,
	) (minio.UploadInfo, error) {
		mu.Lock()
		copied = append(copied, dst.Object)
		mu.Unlock()
		return minio.UploadInfo{}, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	res, err := svc.CopyAll(
		ctx,
		CopyPath{BucketName: "src"},
		CopyPath{BucketName: "dst"},
		nil,
	)

	if err != nil {
		t.Errorf("CopyAll() expected success, got %v", err)
	}
	if res.CopiedCount != 2 {
		t.Errorf("CopiedCount expected 2, got %v", res.CopiedCount)
	}
	if res.ErrorCount != 0 {
		t.Errorf("ErrorCount expected 0, got %v", res.ErrorCount)
	}
	if len(copied) != 2 {
		t.Errorf("copied expected 2, got %v", len(copied))
	}
}

func TestObjectServiceCopyAll_SkipsDirectories(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()
	mock.buckets["src"] = &mockBucket{
		name: "src",
		objects: map[string]*mockObject{
			"dir/":     {key: "dir/", size: 0},
			"dir/a":    {key: "dir/a", size: 5},
			"file.txt": {key: "file.txt", size: 3},
		},
	}
	mock.buckets["dst"] = &mockBucket{
		name:    "dst",
		objects: map[string]*mockObject{},
	}

	var mu sync.Mutex
	count := 0
	mock.copyObjectFunc = func(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (minio.UploadInfo, error) {
		mu.Lock()
		count++
		mu.Unlock()
		return minio.UploadInfo{}, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	res, err := svc.CopyAll(ctx,
		CopyPath{BucketName: "src"},
		CopyPath{BucketName: "dst"},
		nil,
	)

	if err != nil {
		t.Errorf("CopyAll() expected success, got %v", err)
	}
	if res.CopiedCount != 2 {
		t.Errorf("CopiedCount expected 2, got %v", res.CopiedCount)
	}
	if res.ErrorCount != 0 {
		t.Errorf("ErrorCount expected 0, got %v", res.ErrorCount)
	}
	if count != 2 {
		t.Errorf("count expected 2, got %v", count)
	}
}

func TestObjectServiceCopyAll_Filter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()
	mock.buckets["src"] = &mockBucket{
		name: "src",
		objects: map[string]*mockObject{
			"a.txt": {key: "a.txt", size: 1},
			"b.log": {key: "b.log", size: 1},
		},
	}
	mock.buckets["dst"] = &mockBucket{
		name:    "dst",
		objects: map[string]*mockObject{},
	}

	var mu sync.Mutex
	var copied []string
	mock.copyObjectFunc = func(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (minio.UploadInfo, error) {
		mu.Lock()
		copied = append(copied, src.Object)
		mu.Unlock()
		return minio.UploadInfo{}, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	filter := []FilterOptions{
		{Include: ".txt"},
	}

	opts := &CopyAllOptions{
		Filter: &filter,
	}

	res, err := svc.CopyAll(ctx,
		CopyPath{BucketName: "src"},
		CopyPath{BucketName: "dst"},
		opts,
	)

	if err != nil {
		t.Errorf("CopyAll() expected success, got %v", err)
	}
	if res.CopiedCount != 1 {
		t.Errorf("CopiedCount expected 1, got %v", res.CopiedCount)
	}
	if res.ErrorCount != 0 {
		t.Errorf("ErrorCount expected 0, got %v", res.ErrorCount)
	}
	if slices.Compare([]string{"a.txt"}, copied) != 0 {
		t.Errorf("copied expected ['a.txt'], got %v", copied)
	}

}

func TestCopyAll_CopyError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()
	mock.buckets["src"] = &mockBucket{
		name: "src",
		objects: map[string]*mockObject{
			"a.txt": {key: "a.txt", size: 1},
			"b.txt": {key: "b.txt", size: 1},
		},
	}
	mock.buckets["dst"] = &mockBucket{
		name:    "dst",
		objects: map[string]*mockObject{},
	}

	var mu sync.Mutex
	mock.copyObjectFunc = func(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (minio.UploadInfo, error) {
		mu.Lock()
		if src.Object == "b.txt" {
			mu.Unlock()
			return minio.UploadInfo{}, fmt.Errorf("boom")
		}
		mu.Unlock()
		return minio.UploadInfo{}, nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	svc := osClient.Objects()

	res, err := svc.CopyAll(ctx,
		CopyPath{BucketName: "src"},
		CopyPath{BucketName: "dst"},
		nil,
	)

	if err != nil {
		t.Errorf("CopyAll() expected success, got %v", err)
	}
	if res.CopiedCount != 1 {
		t.Errorf("CopiedCount expected 1, got %v", res.CopiedCount)
	}
	if res.ErrorCount != 1 {
		t.Errorf("ErrorCount expected 1, got %v", res.ErrorCount)
	}
	if len(res.Errors) != 1 {
		t.Errorf("len(Errors) expected 1, got %v", len(res.Errors))
	}
}

func TestObjectServiceDeleteAll_InvalidBucketName(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.DeleteAll(context.Background(), "", nil)

	if err == nil {
		t.Error("DeleteAll() expected error for empty bucket name, got nil")
	}

	if _, ok := err.(*InvalidBucketNameError); !ok {
		t.Errorf("DeleteAll() expected InvalidBucketNameError, got %T", err)
	}
}

func TestObjectServiceDeleteAll(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	osClient, _ := New(core, "minioadmin", "minioadmin")
	svc := osClient.Objects()

	_, err := svc.DeleteAll(context.Background(), "bucket-name", nil)

	if err == nil {
		t.Error("DeleteAll() expected error, got nil")
	}
}

func TestObjectServiceDeleteAll_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucket string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "a.txt"}
			ch <- minio.ObjectInfo{Key: "b.txt"}
			ch <- minio.ObjectInfo{Key: "c.txt"}
		}()
		return ch
	}

	removed := make([]string, 0)

	var mu sync.Mutex

	mock.removeObjectFunc = func(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.RemoveObjectOptions,
	) error {
		mu.Lock()
		removed = append(removed, objectName)
		mu.Unlock()
		return nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	res, err := svc.DeleteAll(ctx, "bucket-name", nil)

	if err != nil {
		t.Fatalf("DeleteAll() unexpected error: %v", err)
	}

	if res.DeletedCount != 3 {
		t.Fatalf("DeletedCount expected 3, got %d", res.DeletedCount)
	}

	if res.ErrorCount != 0 {
		t.Fatalf("ErrorCount expected 0, got %d", res.ErrorCount)
	}
}

func TestObjectServiceDeleteAll_PartialError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucket string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "a.txt"}
			ch <- minio.ObjectInfo{Key: "b.txt"}
			ch <- minio.ObjectInfo{Key: "c.txt"}
		}()
		return ch
	}

	var mu sync.Mutex

	mock.removeObjectFunc = func(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.RemoveObjectOptions,
	) error {
		mu.Lock()
		if objectName == "b.txt" {
			mu.Unlock()
			return fmt.Errorf("delete failed")
		}
		mu.Unlock()
		return nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	res, err := svc.DeleteAll(ctx, "bucket-name", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.DeletedCount != 2 {
		t.Fatalf("DeletedCount expected 2, got %d", res.DeletedCount)
	}

	if res.ErrorCount != 1 {
		t.Fatalf("ErrorCount expected 1, got %d", res.ErrorCount)
	}

	if len(res.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(res.Errors))
	}
}

func TestObjectServiceDeleteAll_WithFilter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucket string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "a.txt"}
			ch <- minio.ObjectInfo{Key: "b.log"}
			ch <- minio.ObjectInfo{Key: "c.txt"}
		}()
		return ch
	}

	removed := make([]string, 0)

	var mu sync.Mutex

	mock.removeObjectFunc = func(
		ctx context.Context,
		bucketName string,
		objectName string,
		opts minio.RemoveObjectOptions,
	) error {
		mu.Lock()
		removed = append(removed, objectName)
		mu.Unlock()
		return nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	res, err := svc.DeleteAll(ctx, "bucket-name", &DeleteAllOptions{
		Filter: &[]FilterOptions{
			{Include: ".txt"},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.DeletedCount != 2 {
		t.Fatalf("DeletedCount expected 2, got %d", res.DeletedCount)
	}
}

func TestObjectServiceDeleteAll_WithObjectKeyPrefix(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock := newMockMinioClient()

	var receivedPrefix string

	mock.listObjectsFunc = func(
		ctx context.Context,
		bucket string,
		opts minio.ListObjectsOptions,
	) <-chan minio.ObjectInfo {
		receivedPrefix = opts.Prefix

		ch := make(chan minio.ObjectInfo)
		go func() {
			defer close(ch)
			ch <- minio.ObjectInfo{Key: "photos/a.jpg"}
			ch <- minio.ObjectInfo{Key: "photos/b.jpg"}
		}()
		return ch
	}

	deleted := make([]string, 0)

	var mu sync.Mutex

	mock.removeObjectFunc = func(
		ctx context.Context,
		bucket string,
		object string,
		opts minio.RemoveObjectOptions,
	) error {
		mu.Lock()
		deleted = append(deleted, object)
		mu.Unlock()
		return nil
	}

	core := client.NewMgcClient()
	osClient, _ := New(
		core,
		"minioadmin",
		"minioadmin",
		WithMinioClientInterface(mock),
	)
	svc := osClient.Objects()

	res, err := svc.DeleteAll(
		ctx,
		"bucket-name",
		&DeleteAllOptions{
			ObjectKey: "photos/",
		},
	)

	if err != nil {
		t.Fatalf("DeleteAll() unexpected error: %v", err)
	}

	if receivedPrefix != "photos/" {
		t.Fatalf("expected prefix %q, got %q", "photos/", receivedPrefix)
	}

	if res.DeletedCount != 2 {
		t.Fatalf("expected DeletedCount 2, got %v", res.DeletedCount)
	}

	expectedDeleted := map[string]bool{
		"photos/a.jpg": true,
		"photos/b.jpg": true,
	}

	for _, obj := range deleted {
		if !expectedDeleted[obj] {
			t.Fatalf("unexpected object deleted: %s", obj)
		}
	}

	if res.ErrorCount != 0 {
		t.Fatalf("expected ErrorCount 0, got %v", res.ErrorCount)
	}
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
		Filter: &[]FilterOptions{
			{Include: "text", Exclude: "image"},
		},
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

	if (*opts.Filter)[0].Include != "text" {
		t.Errorf("ObjectListOptions.Filter.Include expected 'text', got %q", (*opts.Filter)[0].Include)
	}

	if (*opts.Filter)[0].Exclude != "image" {
		t.Errorf("ObjectListOptions.Filter.Exclude expected 'image', got %q", (*opts.Filter)[0].Exclude)
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

func TestGetPresignedURLOptions(t *testing.T) {
	t.Parallel()

	expires := time.Duration(300)

	opts := &GetPresignedURLOptions{
		Method:          http.MethodGet,
		ExpiryInSeconds: &expires,
	}

	if opts.Method != http.MethodGet {
		t.Errorf("GetPresignedURLOptions.Method expected http.MethodGet, got %q", opts.Method)
	}

	if opts.ExpiryInSeconds == nil || *opts.ExpiryInSeconds != 300 {
		t.Errorf("GetPresignedURLOptions.ExpiryInSeconds expected 300, got %v", opts.ExpiryInSeconds)
	}
}

func TestObjectDeleteOptions(t *testing.T) {
	t.Parallel()

	opts := DeleteAllOptions{
		Filter:    &[]FilterOptions{{Include: "test", Exclude: "bucket"}},
		BatchSize: 100,
	}

	if opts.BatchSize != 100 {
		t.Errorf("DeleteAllOptions.BatchSize expected 100, got %d", opts.BatchSize)
	}

	if len(*opts.Filter) != 1 || (*opts.Filter)[0].Include != "test" || (*opts.Filter)[0].Exclude != "bucket" {
		t.Errorf("DeleteAllOptions.Filter expected [{'Include': 'test', 'Exclude': 'bucket'}], got %q", *opts.Filter)
	}
}

func TestStorageClassIsValid(t *testing.T) {

	tests := []struct {
		name  string
		class string
		valid bool
	}{
		{"standard", "standard", true},
		{"cold instant", "cold_instant", true},
		{"empty", "", false},
		{"invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := storageClassIsValid(tt.class)

			if got != nil && tt.valid {
				t.Errorf("expected storage class is valid, got %v", got)
			}

			if got == nil && !tt.valid {
				t.Errorf("expected error, got storage class is valid")
			}
		})
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		pattern string
		match   bool
	}{
		{
			name:    "matches pattern",
			key:     "images/photo.jpg",
			pattern: "images/.*",
			match:   true,
		},
		{
			name:    "does not match pattern",
			key:     "docs/readme.md",
			pattern: "images/.*",
			match:   false,
		},
		{
			name:    "invalid regex",
			key:     "file.txt",
			pattern: "[",
			match:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesPattern(tt.key, tt.pattern)
			if got != tt.match {
				t.Fatalf("expected %v, got %v", tt.match, got)
			}
		})
	}
}

func TestShouldProcessObject(t *testing.T) {
	filters := []FilterOptions{
		{Include: "images/.*"},
	}

	if shouldProcessObject(&filters, "docs/readme.md") {
		t.Error("expected object to be skipped")
	}

	if !shouldProcessObject(&filters, "images/photo.jpg") {
		t.Error("expected object to be processed")
	}
}

func TestResolveBatchSize(t *testing.T) {
	t.Run("nil size returns defaultBatchSize", func(t *testing.T) {
		got := resolveBatchSize(nil)

		if got != defaultBatchSize {
			t.Fatalf("expected %d, got %d", defaultBatchSize, got)
		}
	})

	t.Run("zero size returns defaultBatchSize", func(t *testing.T) {
		size := 0

		got := resolveBatchSize(&size)

		if got != defaultBatchSize {
			t.Fatalf("expected %d, got %d", defaultBatchSize, got)
		}
	})

	t.Run("negative size returns defaultBatchSize", func(t *testing.T) {
		size := -10

		got := resolveBatchSize(&size)

		if got != defaultBatchSize {
			t.Fatalf("expected %d, got %d", defaultBatchSize, got)
		}
	})

	t.Run("positive size returns provided value", func(t *testing.T) {
		size := 42

		got := resolveBatchSize(&size)

		if got != size {
			t.Fatalf("expected %d, got %d", size, got)
		}
	})
}

func TestProcessStreamInBatches_AllSuccess(t *testing.T) {
	ctx := context.Background()
	input := make(chan int)

	var handled int32
	var success int32
	var errors int32

	handler := func(ctx context.Context, v int) error {
		atomic.AddInt32(&handled, 1)
		return nil
	}

	onSuccess := func() {
		atomic.AddInt32(&success, 1)
	}

	onError := func(err error) {
		atomic.AddInt32(&errors, 1)
	}

	go func() {
		input <- 1
		input <- 2
		input <- 3
		close(input)
	}()

	processStreamInBatches(
		ctx,
		input,
		2,
		2,
		handler,
		onSuccess,
		onError,
	)

	if handled != 3 {
		t.Fatalf("expected 3 handled items, got %d", handled)
	}
	if success != 3 {
		t.Fatalf("expected 3 successes, got %d", success)
	}
	if errors != 0 {
		t.Fatalf("expected 0 errors, got %d", errors)
	}
}

func TestProcessStreamInBatches_HandlerError(t *testing.T) {
	ctx := context.Background()
	input := make(chan int)

	var success int32
	var errors int32

	handler := func(ctx context.Context, v int) error {
		if v == 2 {
			return fmt.Errorf("boom")
		}
		return nil
	}

	onSuccess := func() {
		atomic.AddInt32(&success, 1)
	}

	onError := func(err error) {
		atomic.AddInt32(&errors, 1)
	}

	go func() {
		input <- 1
		input <- 2
		input <- 3
		close(input)
	}()

	processStreamInBatches(
		ctx,
		input,
		3,
		2,
		handler,
		onSuccess,
		onError,
	)

	if success != 2 {
		t.Fatalf("expected 2 successes, got %d", success)
	}
	if errors != 1 {
		t.Fatalf("expected 1 error, got %d", errors)
	}
}

func TestProcessStreamInBatches_BatchFlush(t *testing.T) {
	ctx := context.Background()
	input := make(chan int)

	var calls int32

	handler := func(ctx context.Context, v int) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}

	go func() {
		input <- 1
		input <- 2
		input <- 3
		input <- 4
		close(input)
	}()

	processStreamInBatches(
		ctx,
		input,
		2,
		1,
		handler,
		func() {},
		func(error) {},
	)

	if calls != 4 {
		t.Fatalf("expected 4 handler calls, got %d", calls)
	}
}

func TestProcessStreamInBatches_FlushOnChannelClose(t *testing.T) {
	ctx := context.Background()

	input := make(chan int)
	go func() {
		defer close(input)
		input <- 1
		input <- 2
	}()

	var success int32

	processStreamInBatches(
		ctx,
		input,
		3,
		1,
		func(ctx context.Context, v int) error {
			return nil
		},
		func() {
			atomic.AddInt32(&success, 1)
		},
		func(err error) {
			t.Fatalf("unexpected error: %v", err)
		},
	)

	if success != 2 {
		t.Fatalf("expected 2 success calls, got %d", success)
	}
}

func intPtr(v int) *int {
	return &v
}
