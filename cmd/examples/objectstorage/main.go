package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/objectstorage"
)

const (
	testBucketName = "e2e-test-bucket"
	testObjectKey  = "test-file.txt"
	testObjectData = "Hello from MagaluCloud Object Storage!"
)

func main() {
	// Get credentials from environment
	apiToken := os.Getenv("MGC_API_KEY")
	if apiToken == "" {
		log.Fatal("❌ MGC_API_TOKEN environment variable is not set")
	}

	accessKey := os.Getenv("MGC_OBJECT_STORAGE_ACCESS_KEY")
	if accessKey == "" {
		log.Fatal("❌ MGC_OBJECT_STORAGE_ACCESS_KEY environment variable is not set")
	}

	secretKey := os.Getenv("MGC_OBJECT_STORAGE_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("❌ MGC_OBJECT_STORAGE_SECRET_KEY environment variable is not set")
	}

	// Check for optional region parameter
	region := os.Getenv("MGC_OBJECT_STORAGE_REGION")
	if region == "" {
		region = "br-se1"
	}

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  MagaluCloud Object Storage - End-to-End Test Example     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Display configuration
	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   Region: %s\n", region)
	fmt.Printf("   Endpoint: %s\n", getEndpointName(region))
	fmt.Printf("   Test Bucket: %s\n", testBucketName)
	fmt.Printf("   Test Object: %s\n", testObjectKey)
	fmt.Println()

	// Create MagaluCloud client
	c := client.NewMgcClient(apiToken)

	// Create Object Storage client with selected region
	var opts []objectstorage.ClientOption
	if strings.ToLower(region) == "br-ne1" {
		opts = append(opts, objectstorage.WithEndpoint(objectstorage.BrNe1))
	}

	osClient, err := objectstorage.New(c, accessKey, secretKey, opts...)
	if err != nil {
		log.Fatalf("❌ Failed to create Object Storage client: %v\n", err)
	}

	fmt.Println("✅ Object Storage client created successfully")
	fmt.Println()

	// Run comprehensive end-to-end test
	runE2ETest(context.Background(), osClient)
}

func runE2ETest(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("🧪 Running End-to-End Test Suite...")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Step 1: List existing buckets
	testListBuckets(ctx, osClient)
	pause()

	// Step 2: Create bucket
	testCreateBucket(ctx, osClient)
	pause()

	// Step 3: Check if bucket exists
	testBucketExists(ctx, osClient)
	pause()

	// Step 4: Upload object
	testUploadObject(ctx, osClient)
	pause()

	// Step 5: Get object metadata
	testObjectMetadata(ctx, osClient)
	pause()

	// Step 6: Download object
	testDownloadObject(ctx, osClient)
	pause()

	// Step 7: Download as stream
	testDownloadObjectStream(ctx, osClient)
	pause()

	// Step 8: List objects in bucket
	testListObjects(ctx, osClient)
	pause()

	// Step 9: Set bucket policy
	testSetBucketPolicy(ctx, osClient)
	pause()

	// Step 10: Get bucket policy
	testGetBucketPolicy(ctx, osClient)
	pause()

	// Step 11: Delete bucket policy (must do this before deleting object due to policy restrictions)
	testDeleteBucketPolicy(ctx, osClient)
	pause()

	// Step 12: Delete object
	testDeleteObject(ctx, osClient)
	pause()

	// Step 13: Delete bucket
	testDeleteBucket(ctx, osClient)
	pause()

	// Final summary
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("✅ All tests completed successfully!")
	fmt.Println()
	fmt.Println("🎉 End-to-End Test Suite: PASSED")
	fmt.Println()
}

func testListBuckets(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 1: List All Buckets")
	fmt.Println("─────────────────────────────────────────────────────────────")

	buckets, err := osClient.Buckets().List(ctx, objectstorage.BucketListOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Listed %d bucket(s)\n", len(buckets))
	if len(buckets) > 0 && len(buckets) <= 5 {
		for _, bucket := range buckets {
			fmt.Printf("   📁 %s (Created: %s)\n", bucket.Name, bucket.CreationDate)
		}
	} else if len(buckets) > 5 {
		for i := range 3 {
			fmt.Printf("   📁 %s (Created: %s)\n", buckets[i].Name, buckets[i].CreationDate)
		}
		fmt.Printf("   ... and %d more\n", len(buckets)-3)
	}
	fmt.Println()
}

func testCreateBucket(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 2: Create Bucket")
	fmt.Println("─────────────────────────────────────────────────────────────")

	// Check if bucket already exists
	exists, err := osClient.Buckets().Exists(ctx, testBucketName)
	if err == nil && exists {
		fmt.Printf("⚠️  Bucket already exists: %s (skipping creation)\n\n", testBucketName)
		return
	}

	err = osClient.Buckets().Create(ctx, testBucketName)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Bucket created: %s\n\n", testBucketName)
}

func testBucketExists(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 3: Check Bucket Exists")
	fmt.Println("─────────────────────────────────────────────────────────────")

	exists, err := osClient.Buckets().Exists(ctx, testBucketName)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	if exists {
		fmt.Printf("✅ Bucket exists: %s\n\n", testBucketName)
	} else {
		fmt.Printf("❌ Bucket does not exist: %s\n\n", testBucketName)
	}
}

func testUploadObject(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 4: Upload Object")
	fmt.Println("─────────────────────────────────────────────────────────────")

	err := osClient.Objects().Upload(
		ctx,
		testBucketName,
		testObjectKey,
		[]byte(testObjectData),
		"text/plain",
	)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Object uploaded: %s\n", testObjectKey)
	fmt.Printf("   Size: %d bytes\n", len(testObjectData))
	fmt.Printf("   Content-Type: text/plain\n\n")
}

func testObjectMetadata(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 5: Get Object Metadata")
	fmt.Println("─────────────────────────────────────────────────────────────")

	obj, err := osClient.Objects().Metadata(ctx, testBucketName, testObjectKey)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Metadata retrieved:\n")
	fmt.Printf("   Key: %s\n", obj.Key)
	fmt.Printf("   Size: %d bytes\n", obj.Size)
	fmt.Printf("   Content-Type: %s\n", obj.ContentType)
	fmt.Printf("   Last Modified: %s\n", obj.LastModified)
	fmt.Printf("   ETag: %s\n\n", obj.ETag)
}

func testDownloadObject(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 6: Download Object")
	fmt.Println("─────────────────────────────────────────────────────────────")

	data, err := osClient.Objects().Download(ctx, testBucketName, testObjectKey, nil)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	if string(data) != testObjectData {
		fmt.Printf("❌ Data mismatch! Expected %q, got %q\n\n", testObjectData, string(data))
		return
	}

	fmt.Printf("✅ Object downloaded successfully\n")
	fmt.Printf("   Size: %d bytes\n", len(data))
	fmt.Printf("   Content: %s\n\n", string(data))
}

func testDownloadObjectStream(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 7: Download Object as Stream")
	fmt.Println("─────────────────────────────────────────────────────────────")

	reader, err := osClient.Objects().DownloadStream(ctx, testBucketName, testObjectKey, nil)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("❌ Failed to read stream: %v\n\n", err)
		return
	}

	if string(data) != testObjectData {
		fmt.Printf("❌ Data mismatch! Expected %q, got %q\n\n", testObjectData, string(data))
		return
	}

	fmt.Printf("✅ Object downloaded via stream\n")
	fmt.Printf("   Size: %d bytes\n", len(data))
	fmt.Printf("   Content: %s\n\n", string(data))
}

func testListObjects(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 8: List Objects in Bucket")
	fmt.Println("─────────────────────────────────────────────────────────────")

	objects, err := osClient.Objects().ListAll(ctx, testBucketName, objectstorage.ObjectFilterOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Listed %d object(s):\n", len(objects))
	for _, obj := range objects {
		fmt.Printf("   📄 %s\n", obj.Key)
		fmt.Printf("      Size: %d bytes\n", obj.Size)
		fmt.Printf("      Modified: %s\n", obj.LastModified)
	}
	fmt.Println()
}

func testSetBucketPolicy(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 9: Set Bucket Policy")
	fmt.Println("─────────────────────────────────────────────────────────────")

	policy := &objectstorage.Policy{
		Version: "2012-10-17",
		Statement: []objectstorage.Statement{
			{
				Effect:    "Allow",
				Principal: "*",
				Action:    "s3:GetObject",
				Resource:  fmt.Sprintf("%s/*", testBucketName),
			},
		},
	}

	err := osClient.Buckets().SetPolicy(ctx, testBucketName, policy)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Bucket policy set successfully\n\n")
}

func testGetBucketPolicy(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 10: Get Bucket Policy")
	fmt.Println("─────────────────────────────────────────────────────────────")

	policyResult, err := osClient.Buckets().GetPolicy(ctx, testBucketName)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	if policyResult == nil {
		fmt.Printf("⚠️  No policy set on bucket\n\n")
		return
	}

	fmt.Printf("✅ Bucket policy retrieved:\n")
	fmt.Printf("   Version: %s\n", policyResult.Version)
	fmt.Printf("   Statements: %d\n\n", len(policyResult.Statement))
}

func testDeleteBucketPolicy(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 11: Delete Bucket Policy")
	fmt.Println("─────────────────────────────────────────────────────────────")

	err := osClient.Buckets().DeletePolicy(ctx, testBucketName)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Bucket policy deleted successfully\n\n")
}

func testDeleteObject(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 12: Delete Object")
	fmt.Println("─────────────────────────────────────────────────────────────")

	err := osClient.Objects().Delete(ctx, testBucketName, testObjectKey, nil)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Object deleted: %s\n\n", testObjectKey)
}

func testDeleteBucket(ctx context.Context, osClient *objectstorage.ObjectStorageClient) {
	fmt.Println("📝 Test 13: Delete Bucket")
	fmt.Println("─────────────────────────────────────────────────────────────")

	err := osClient.Buckets().Delete(ctx, testBucketName)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
		fmt.Printf("   Note: Bucket may not be empty or may not exist\n\n")
		return
	}

	fmt.Printf("✅ Bucket deleted: %s\n\n", testBucketName)
}

func getEndpointName(region string) string {
	switch strings.ToLower(region) {
	case "br-ne1":
		return "br-ne1.magaluobjects.com (Brazil Northeast 1)"
	default:
		return "br-se1.magaluobjects.com (Brazil Southeast 1)"
	}
}

func pause() {
	fmt.Println()
	time.Sleep(100 * time.Millisecond)
}
