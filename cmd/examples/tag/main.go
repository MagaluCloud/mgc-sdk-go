package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/tag"
)

func main() {
	// Tags
	ExampleListTags()
	tagName := ExampleCreateTag()
	ExampleGetTag(tagName)
	ExampleUpdateTag(tagName)

	// Tag values
	valueName := ExampleCreateTagValue(tagName)
	ExampleListTagValues(tagName)
	ExampleGetTagValue(tagName, valueName)
	ExampleUpdateTagValue(tagName, valueName)

	// Resource types
	ExampleListResourceTypes()

	// Tagged resources. Resources belong to other products, so set MGC_RESOURCE_ID
	// to the ID of one you already own to run the attach and detach examples.
	ExampleListResources()
	if externalID := os.Getenv("MGC_RESOURCE_ID"); externalID != "" {
		ExampleAttachTags(externalID, tagName, valueName)
		ExampleGetResource(externalID)
		ExampleDetachTag(externalID, tagName)
	}

	// Cleanup
	ExampleDeleteTagValue(tagName, valueName)
	ExampleDeleteTag(tagName)
}

// newTagClient builds a tag client from the MGC_API_TOKEN environment variable.
// Tags are a global service, so no region needs to be configured.
func newTagClient() *tag.TagClient {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}

	core := client.NewMgcClient(client.WithJWToken(apiToken))

	// The API accepts a tenant override for accounts with more than one tenant.
	var opts []tag.ClientOption
	if tenantID := os.Getenv("MGC_TENANT_ID"); tenantID != "" {
		opts = append(opts, tag.WithTenantID(tenantID))
	}

	return tag.New(core, opts...)
}

func ExampleListTags() {
	tagClient := newTagClient()

	tags, err := tagClient.Tags().List(context.Background(), tag.ListTagsOptions{
		Kinds: []tag.TagKind{tag.TagKindFinops},
		Limit: helpers.IntPtr(10),
		Sort:  helpers.StrPtr("name:asc"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d tags:\n", len(tags))
	for _, t := range tags {
		fmt.Printf("  Tag: %s (color: %s, values: %d)\n", t.Name, deref(t.Color), len(t.Values))
	}
}

func ExampleCreateTag() string {
	tagClient := newTagClient()

	t, err := tagClient.Tags().Create(context.Background(), tag.CreateTagRequest{
		Name:        "environment_" + randomString(),
		Description: helpers.StrPtr("Identifies the deployment environment"),
		Color:       helpers.StrPtr("F54927"),
		Kinds:       []tag.TagKind{tag.TagKindFinops},
		Values: []tag.CreateTagValueRequest{
			{Name: "production", Description: helpers.StrPtr("Production workloads")},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created tag: %s (color: %s)\n", t.Name, deref(t.Color))
	return t.Name
}

func ExampleGetTag(tagName string) {
	tagClient := newTagClient()

	t, err := tagClient.Tags().Get(context.Background(), tagName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tag Details:\n")
	fmt.Printf("  Name: %s\n", t.Name)
	fmt.Printf("  Description: %s\n", deref(t.Description))
	fmt.Printf("  Color: %s\n", deref(t.Color))
	fmt.Printf("  Created At: %s\n", time.Time(t.CreatedAt))
}

func ExampleUpdateTag(tagName string) {
	tagClient := newTagClient()

	t, err := tagClient.Tags().Update(context.Background(), tagName, tag.UpdateTagRequest{
		Description: helpers.StrPtr("Identifies the deployment environment (updated)"),
		Color:       helpers.StrPtr("F54927"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated tag: %s (color: %s)\n", t.Name, deref(t.Color))
}

func ExampleCreateTagValue(tagName string) string {
	tagClient := newTagClient()

	v, err := tagClient.Values().Create(context.Background(), tagName, tag.CreateTagValueRequest{
		Name:        "staging",
		Description: helpers.StrPtr("Staging workloads"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created value %q on tag %q\n", v.Name, tagName)
	return v.Name
}

func ExampleListTagValues(tagName string) {
	tagClient := newTagClient()

	values, err := tagClient.Values().List(context.Background(), tagName, tag.ListTagValuesOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tag %q has %d values:\n", tagName, len(values))
	for _, v := range values {
		fmt.Printf("  Value: %s - %s\n", v.Name, deref(v.Description))
	}
}

func ExampleGetTagValue(tagName, valueName string) {
	tagClient := newTagClient()

	v, err := tagClient.Values().Get(context.Background(), tagName, valueName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Value Details:\n")
	fmt.Printf("  Name: %s\n", v.Name)
	fmt.Printf("  Description: %s\n", deref(v.Description))
	if v.Tag != nil {
		fmt.Printf("  Tag: %s\n", v.Tag.Name)
	}
}

func ExampleUpdateTagValue(tagName, valueName string) {
	tagClient := newTagClient()

	v, err := tagClient.Values().Update(context.Background(), tagName, valueName, tag.UpdateTagValueRequest{
		Description: helpers.StrPtr("Staging workloads (updated)"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated value %q: %s\n", v.Name, deref(v.Description))
}

func ExampleListResourceTypes() {
	tagClient := newTagClient()

	product := tag.Product("network")
	types, err := tagClient.ResourceTypes().List(context.Background(), tag.ListResourceTypesOptions{
		Product: &product,
		Limit:   helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d resource types for product %q:\n", len(types), product)
	for _, rt := range types {
		fmt.Printf("  Name: %s (product: %s)\n", rt.Name, rt.Product)
	}
}

func ExampleListResources() {
	tagClient := newTagClient()

	resources, err := tagClient.Resources().List(context.Background(), tag.ListResourcesOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d tagged resources:\n", len(resources))
	for _, r := range resources {
		fmt.Printf("  %s (%s, %s) with %d tags\n", r.ExternalID, r.ResourceType.Name, r.Region, len(r.Tags))
	}
}

func ExampleGetResource(externalID string) {
	tagClient := newTagClient()

	resource, err := tagClient.Resources().Get(context.Background(), externalID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Resource Details:\n")
	fmt.Printf("  External ID: %s\n", resource.ExternalID)
	fmt.Printf("  Type: %s (product: %s)\n", resource.ResourceType.Name, resource.ResourceType.Product)
	fmt.Printf("  Region: %s\n", resource.Region)
	for _, t := range resource.Tags {
		fmt.Printf("  Tag: %s = %s\n", t.Name, t.Value)
	}
}

func ExampleAttachTags(externalID, tagName, valueName string) {
	tagClient := newTagClient()

	resource, err := tagClient.Resources().AttachTags(context.Background(), externalID, tag.AttachTagsRequest{
		Tags: []tag.AttachTag{
			{Name: tagName, Value: valueName},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Attached tag %q with value %q to resource %q, now with %d tags\n",
		tagName, valueName, externalID, len(resource.Tags))
}

func ExampleDetachTag(externalID, tagName string) {
	tagClient := newTagClient()

	if err := tagClient.Resources().DetachTag(context.Background(), externalID, tagName); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Detached tag %q from resource %q\n", tagName, externalID)
}

func ExampleDeleteTagValue(tagName, valueName string) {
	tagClient := newTagClient()

	if err := tagClient.Values().Delete(context.Background(), tagName, valueName); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted value %q from tag %q\n", valueName, tagName)
}

func ExampleDeleteTag(tagName string) {
	tagClient := newTagClient()

	if err := tagClient.Tags().Delete(context.Background(), tagName); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted tag %q\n", tagName)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func randomString() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(rand.Int64(), 10)
}
