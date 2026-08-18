package tag

import (
	"net/url"
	"strconv"
)

// TagKind describes what a tag is used for.
// The API may return kinds beyond the constants below.
type TagKind string

const (
	TagKindFinops TagKind = "finops"
)

// Product identifies the product that owns a resource type, as in "virtual-machine".
// The products that support tags are the ones returned by ResourceTypes().List.
type Product string

// ResourceTypeName identifies a type of resource that supports tagging, prefixed by
// the product that owns it, as in "vm.instance".
// The types that support tags are the ones returned by ResourceTypes().List, and grow
// as products adopt tagging.
type ResourceTypeName string

// ResourceTypeUnknown is the type of a resource the API does not recognize.
const ResourceTypeUnknown ResourceTypeName = "un.unknown"

// listOptions holds the pagination parameters accepted by every list endpoint.
type listOptions struct {
	Limit  *int
	Offset *int
	Sort   *string
}

// makeListQuery builds the pagination query shared by every list endpoint.
func makeListQuery(opts listOptions) url.Values {
	query := make(url.Values)
	if opts.Limit != nil {
		query.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		query.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		query.Set("_sort", *opts.Sort)
	}
	return query
}
