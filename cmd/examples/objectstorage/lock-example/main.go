package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/objectstorage"
)

const (
	testBucketName = "lock-example-bucket"
	testObjectKey  = "protected-document.txt"
	testObjectData = "This document is protected by Object Lock for compliance purposes."
)

func main() {
	// Get credentials from environment
	apiToken := os.Getenv("MGC_API_KEY")
	if apiToken == "" {
		log.Fatal("❌ MGC_API_KEY environment variable is not set")
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
	fmt.Println("║  MagaluCloud Object Storage - Object Lock Example         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Display configuration
	fmt.Printf("🔧 Configuration:\n")
	fmt.Printf("   Region: %s\n", region)
	fmt.Printf("   Bucket: %s\n", testBucketName)
	fmt.Printf("   Object: %s\n\n", testObjectKey)

	// Initialize the client
	coreClient := client.NewMgcClient(apiToken)

	// Create Object Storage client with selected region
	var opts []objectstorage.ClientOption
	if strings.ToLower(region) == "br-ne1" {
		opts = append(opts, objectstorage.WithEndpoint(objectstorage.BrNe1))
	}

	osClient, err := objectstorage.New(coreClient, accessKey, secretKey, opts...)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Object Storage client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Step 1: Create bucket (if not exists)
	fmt.Println("📍 Step 1: Create bucket")
	fmt.Printf("   Creating bucket '%s'...\n", testBucketName)
	err = osClient.Buckets().Create(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ⚠️  Bucket creation failed or already exists: %v\n", err)
	} else {
		fmt.Println("   ✓ Bucket created successfully")
	}
	fmt.Println()
	pause()

	// Step 2: Lock the bucket
	fmt.Println("📍 Step 2: Enable Object Lock on bucket")
	fmt.Printf("   Locking bucket '%s'...\n", testBucketName)
	err = osClient.Buckets().LockBucket(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to lock bucket: %v\n", err)
	} else {
		fmt.Println("   ✓ Bucket locked successfully")
	}
	fmt.Println()
	pause()

	// Step 3: Check bucket lock status
	fmt.Println("📍 Step 3: Check bucket lock status")
	isLocked, err := osClient.Buckets().GetBucketLockStatus(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to get bucket lock status: %v\n", err)
	} else {
		if isLocked {
			fmt.Println("   ✓ Bucket is locked (Object Lock enabled)")
		} else {
			fmt.Println("   ✗ Bucket is not locked")
		}
	}
	fmt.Println()
	pause()

	// Step 4: Upload an object
	fmt.Println("📍 Step 4: Upload object to locked bucket")
	fmt.Printf("   Uploading '%s'...\n", testObjectKey)
	err = osClient.Objects().Upload(ctx, testBucketName, testObjectKey, []byte(testObjectData), "text/plain")
	if err != nil {
		fmt.Printf("   ❌ Failed to upload object: %v\n", err)
	} else {
		fmt.Println("   ✓ Object uploaded successfully")
	}
	fmt.Println()
	pause()

	// Step 5: Lock the object with retention period
	fmt.Println("📍 Step 5: Apply retention lock to object")
	retentionDays := 7
	retainUntil := time.Now().UTC().AddDate(0, 0, retentionDays)
	fmt.Printf("   Locking object for %d days (until %s)...\n", retentionDays, retainUntil.Format("2006-01-02 15:04:05"))
	err = osClient.Objects().LockObject(ctx, testBucketName, testObjectKey, retainUntil)
	if err != nil {
		fmt.Printf("   ❌ Failed to lock object: %v\n", err)
	} else {
		fmt.Println("   ✓ Object locked successfully")
	}
	fmt.Println()
	pause()

	// Step 6: Check object lock status
	fmt.Println("📍 Step 6: Check object lock status")
	objIsLocked, err := osClient.Objects().GetObjectLockStatus(ctx, testBucketName, testObjectKey)
	if err != nil {
		fmt.Printf("   ❌ Failed to get object lock status: %v\n", err)
	} else {
		if objIsLocked {
			fmt.Printf("   ✓ Object is locked\n")
			fmt.Printf("   📅 Retain until: %s\n", retainUntil.Format("2006-01-02 15:04:05"))
			remaining := time.Until(retainUntil)
			fmt.Printf("   ⏳ Time remaining: %d days, %d hours\n", int(remaining.Hours())/24, int(remaining.Hours())%24)
		} else {
			fmt.Println("   ✗ Object is not locked")
		}
	}
	fmt.Println()
	pause()

	// Step 7: Try to delete the locked object (should fail)
	fmt.Println("📍 Step 7: Attempt to delete locked object (should fail)")
	fmt.Printf("   Attempting to delete '%s'...\n", testObjectKey)
	err = osClient.Objects().Delete(ctx, testBucketName, testObjectKey, nil)
	if err != nil {
		fmt.Printf("   ✓ Deletion blocked as expected: %v\n", err)
	} else {
		fmt.Println("   ⚠️  Object was deleted (lock may not be active)")
	}
	fmt.Println()
	pause()

	// Step 8: Download the object to verify it still exists
	fmt.Println("📍 Step 8: Download object to verify it's still protected")
	fmt.Printf("   Downloading '%s'...\n", testObjectKey)
	data, err := osClient.Objects().Download(ctx, testBucketName, testObjectKey, nil)
	if err != nil {
		fmt.Printf("   ❌ Failed to download object: %v\n", err)
	} else {
		fmt.Println("   ✓ Object downloaded successfully")
		fmt.Printf("   📄 Content: %s\n", string(data))
	}
	fmt.Println()
	pause()

	// Step 9: Get object metadata
	fmt.Println("📍 Step 9: Get object metadata")
	fmt.Printf("   Retrieving metadata for '%s'...\n", testObjectKey)
	metadata, err := osClient.Objects().Metadata(ctx, testBucketName, testObjectKey)
	if err != nil {
		fmt.Printf("   ❌ Failed to get metadata: %v\n", err)
	} else {
		fmt.Println("   ✓ Metadata retrieved successfully")
		fmt.Printf("   📊 Details:\n")
		fmt.Printf("      - Key: %s\n", metadata.Key)
		fmt.Printf("      - Size: %d bytes\n", metadata.Size)
		fmt.Printf("      - Content-Type: %s\n", metadata.ContentType)
		fmt.Printf("      - Last Modified: %s\n", metadata.LastModified.Format("2006-01-02 15:04:05"))
		fmt.Printf("      - ETag: %s\n", metadata.ETag)
	}
	fmt.Println()
	pause()

	// Step 10: Unlock the object (requires governance bypass)
	fmt.Println("📍 Step 10: Unlock object (remove retention)")
	fmt.Printf("   Unlocking '%s'...\n", testObjectKey)
	err = osClient.Objects().UnlockObject(ctx, testBucketName, testObjectKey)
	if err != nil {
		fmt.Printf("   ❌ Failed to unlock object: %v\n", err)
	} else {
		fmt.Println("   ✓ Object unlocked successfully")
	}
	fmt.Println()
	pause()

	// Step 11: Verify object is unlocked
	fmt.Println("📍 Step 11: Verify object is now unlocked")
	objIsLocked, err = osClient.Objects().GetObjectLockStatus(ctx, testBucketName, testObjectKey)
	if err != nil {
		fmt.Printf("   ❌ Failed to get object lock status: %v\n", err)
	} else {
		if !objIsLocked {
			fmt.Println("   ✓ Object is no longer locked")
		} else {
			fmt.Println("   ⚠️  Object is still locked")
		}
	}
	fmt.Println()
	pause()

	// Step 12: Delete the object (should succeed now)
	fmt.Println("📍 Step 12: Delete unlocked object")
	fmt.Printf("   Deleting '%s'...\n", testObjectKey)
	err = osClient.Objects().Delete(ctx, testBucketName, testObjectKey, nil)
	if err != nil {
		fmt.Printf("   ❌ Failed to delete object: %v\n", err)
	} else {
		fmt.Println("   ✓ Object deleted successfully")
	}
	fmt.Println()
	pause()

	// Step 13: Unlock the bucket
	fmt.Println("📍 Step 13: Disable Object Lock on bucket")
	fmt.Printf("   Unlocking bucket '%s'...\n", testBucketName)
	err = osClient.Buckets().UnlockBucket(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to unlock bucket: %v\n", err)
	} else {
		fmt.Println("   ✓ Bucket unlocked successfully")
	}
	fmt.Println()
	pause()

	// Step 14: Verify bucket is unlocked
	fmt.Println("📍 Step 14: Verify bucket is now unlocked")
	isLocked, err = osClient.Buckets().GetBucketLockStatus(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to get bucket lock status: %v\n", err)
	} else {
		if !isLocked {
			fmt.Println("   ✓ Bucket is no longer locked")
		} else {
			fmt.Println("   ⚠️  Bucket is still locked")
		}
	}
	fmt.Println()
	pause()

	// Step 15: Clean up - delete the bucket
	fmt.Println("📍 Step 15: Clean up - delete bucket")
	fmt.Printf("   Deleting bucket '%s'...\n", testBucketName)
	err = osClient.Buckets().Delete(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to delete bucket: %v\n", err)
	} else {
		fmt.Println("   ✓ Bucket deleted successfully")
	}
	fmt.Println()

	// Summary
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  ✓ Object Lock Example Completed Successfully!            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("📚 Key Takeaways:")
	fmt.Println("   • Bucket-level locking enables Object Lock for new objects")
	fmt.Println("   • Object-level locking applies retention to specific objects")
	fmt.Println("   • Locked objects cannot be deleted until unlock is called")
	fmt.Println("   • Retention periods help ensure compliance and data protection")
	fmt.Println()
}

func pause() {
	fmt.Println("   ⏸️  Press Enter to continue...")
	fmt.Scanln()
}
