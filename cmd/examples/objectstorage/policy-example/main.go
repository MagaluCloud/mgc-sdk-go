package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/objectstorage"
)

const (
	testBucketName = "policy-example-bucket"
	testObjectKey  = "public-document.txt"
	testObjectData = "This is a publicly readable object"
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
	fmt.Println("║  MagaluCloud Object Storage - Bucket Policy Example       ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Display configuration
	fmt.Printf("🔧 Configuration:\n")
	fmt.Printf("   Region: %s\n", region)
	fmt.Printf("   Bucket: %s\n", testBucketName)
	fmt.Printf("   Object: %s\n\n", testObjectKey)

	// Initialize the client
	coreClient := client.NewMgcClient(client.WithAPIKey(apiToken))

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

	// Step 1: Create bucket
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

	// Step 2: Upload an object
	fmt.Println("📍 Step 2: Upload object to bucket")
	fmt.Printf("   Uploading '%s'...\n", testObjectKey)
	err = osClient.Objects().Upload(ctx, testBucketName, testObjectKey, []byte(testObjectData), "text/plain", helpers.StrPtr("standard"))
	if err != nil {
		fmt.Printf("   ❌ Failed to upload object: %v\n", err)
	} else {
		fmt.Println("   ✓ Object uploaded successfully")
	}
	fmt.Println()
	pause()

	// Step 3: Create a public read policy
	fmt.Println("📍 Step 3: Create public read policy")
	fmt.Println("   Creating policy that allows public read access...")

	policy := &objectstorage.Policy{
		Version: "2012-10-17",
		Id:      "PublicReadPolicy",
		Statement: []objectstorage.Statement{
			{
				Sid:       "PublicReadGetObject",
				Effect:    "Allow",
				Principal: "*",
				Action:    "s3:GetObject",
				Resource:  fmt.Sprintf("%s/*", testBucketName),
			},
		},
	}

	fmt.Printf("   Policy structure:\n")
	prettyPrintPolicy(policy)
	fmt.Println()
	pause()

	// Step 4: Set the bucket policy
	fmt.Println("📍 Step 4: Apply policy to bucket")
	fmt.Printf("   Applying policy to '%s'...\n", testBucketName)
	err = osClient.Buckets().SetPolicy(ctx, testBucketName, policy)
	if err != nil {
		fmt.Printf("   ❌ Failed to set policy: %v\n", err)
	} else {
		fmt.Println("   ✓ Policy applied successfully")
	}
	fmt.Println()
	pause()

	// Step 5: Get the bucket policy
	fmt.Println("📍 Step 5: Retrieve bucket policy")
	fmt.Printf("   Getting policy from '%s'...\n", testBucketName)
	retrievedPolicy, err := osClient.Buckets().GetPolicy(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to get policy: %v\n", err)
	} else if retrievedPolicy == nil {
		fmt.Println("   ℹ️  No policy found on bucket")
	} else {
		fmt.Println("   ✓ Policy retrieved successfully")
		fmt.Printf("   Retrieved policy structure:\n")
		prettyPrintPolicy(retrievedPolicy)
	}
	fmt.Println()
	pause()

	// Step 6: Create a more complex policy with multiple statements
	fmt.Println("📍 Step 6: Create complex policy with multiple statements")
	fmt.Println("   Creating policy with Allow and Deny statements...")

	complexPolicy := &objectstorage.Policy{
		Version: "2012-10-17",
		Id:      "ComplexPolicy",
		Statement: []objectstorage.Statement{
			{
				Sid:       "AllowPublicRead",
				Effect:    "Allow",
				Principal: "*",
				Action:    "s3:GetObject",
				Resource:  fmt.Sprintf("%s/public/*", testBucketName),
			},
			{
				Sid:       "DenyPrivateDelete",
				Effect:    "Deny",
				Principal: "*",
				Action:    "s3:DeleteObject",
				Resource:  fmt.Sprintf("%s/private/*", testBucketName),
			},
		},
	}

	fmt.Printf("   Complex policy structure:\n")
	prettyPrintPolicy(complexPolicy)
	fmt.Println()
	pause()

	// Step 7: Update bucket policy with complex policy
	fmt.Println("📍 Step 7: Update bucket with complex policy")
	fmt.Printf("   Applying complex policy to '%s'...\n", testBucketName)
	err = osClient.Buckets().SetPolicy(ctx, testBucketName, complexPolicy)
	if err != nil {
		fmt.Printf("   ❌ Failed to set complex policy: %v\n", err)
	} else {
		fmt.Println("   ✓ Complex policy applied successfully")
	}
	fmt.Println()
	pause()

	// Step 8: Verify updated policy
	fmt.Println("📍 Step 8: Verify updated policy")
	fmt.Printf("   Getting updated policy from '%s'...\n", testBucketName)
	updatedPolicy, err := osClient.Buckets().GetPolicy(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to get updated policy: %v\n", err)
	} else if updatedPolicy == nil {
		fmt.Println("   ℹ️  No policy found on bucket")
	} else {
		fmt.Println("   ✓ Updated policy retrieved successfully")
		fmt.Printf("   Number of statements: %d\n", len(updatedPolicy.Statement))
		for i, stmt := range updatedPolicy.Statement {
			fmt.Printf("      Statement %d: %s (%s)\n", i+1, stmt.Sid, stmt.Effect)
		}
	}
	fmt.Println()
	pause()

	// Step 9: Delete the bucket policy
	fmt.Println("📍 Step 9: Delete bucket policy")
	fmt.Printf("   Deleting policy from '%s'...\n", testBucketName)
	err = osClient.Buckets().DeletePolicy(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ❌ Failed to delete policy: %v\n", err)
	} else {
		fmt.Println("   ✓ Policy deleted successfully")
	}
	fmt.Println()
	pause()

	// Step 10: Verify policy is deleted
	fmt.Println("📍 Step 10: Verify policy is deleted")
	fmt.Printf("   Checking policy on '%s'...\n", testBucketName)
	deletedPolicy, err := osClient.Buckets().GetPolicy(ctx, testBucketName)
	if err != nil {
		fmt.Printf("   ℹ️  Error getting policy (may be expected): %v\n", err)
	} else if deletedPolicy == nil {
		fmt.Println("   ✓ Policy successfully deleted (bucket has no policy)")
	} else {
		fmt.Printf("   ⚠️  Policy still exists: %v\n", deletedPolicy)
	}
	fmt.Println()
	pause()

	// Step 11: Clean up - delete the object
	fmt.Println("📍 Step 11: Clean up - delete object")
	fmt.Printf("   Deleting '%s'...\n", testObjectKey)
	err = osClient.Objects().Delete(ctx, testBucketName, testObjectKey, nil)
	if err != nil {
		fmt.Printf("   ❌ Failed to delete object: %v\n", err)
	} else {
		fmt.Println("   ✓ Object deleted successfully")
	}
	fmt.Println()
	pause()

	// Step 12: Clean up - delete the bucket
	fmt.Println("📍 Step 12: Clean up - delete bucket")
	fmt.Printf("   Deleting bucket '%s'...\n", testBucketName)
	err = osClient.Buckets().Delete(ctx, testBucketName, false)
	if err != nil {
		fmt.Printf("   ❌ Failed to delete bucket: %v\n", err)
	} else {
		fmt.Println("   ✓ Bucket deleted successfully")
	}
	fmt.Println()

	// Summary
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  ✓ Bucket Policy Example Completed Successfully!          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("📚 Key Takeaways:")
	fmt.Println("   • Policies control access to buckets and objects")
	fmt.Println("   • Support Allow and Deny statements")
	fmt.Println("   • Can target specific principals and resources")
	fmt.Println("   • Policies are returned as structured objects")
	fmt.Println("   • DeletePolicy removes all policies from a bucket")
	fmt.Println()
}

func prettyPrintPolicy(policy *objectstorage.Policy) {
	data, err := json.MarshalIndent(policy, "      ", "  ")
	if err != nil {
		fmt.Printf("      Error marshaling policy: %v\n", err)
		return
	}
	fmt.Printf("      %s\n", string(data))
}

func pause() {
	fmt.Println("   ⏸️  Press Enter to continue...")
	fmt.Scanln()
}
