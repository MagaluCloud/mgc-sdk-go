package containerregistry

import (
	"net/url"
	"strconv"
	"strings"
)

type (
	Meta struct {
		Page Page `json:"page"`
	}

	Page struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}
)

func CreatePaginationParams(opts ListOptions) url.Values {
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
	if len(opts.Expand) > 0 {
		query.Set("_expand", strings.Join(opts.Expand, ","))
	}
	return query
}
