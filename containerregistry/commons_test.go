package containerregistry

import (
	"net/url"
	"reflect"
	"testing"
)

func TestCreatePaginationParams(t *testing.T) {
	tests := []struct {
		name string
		opts ListOptions
		want url.Values
	}{
		{
			name: "empty options",
			opts: ListOptions{},
			want: url.Values{},
		},
		{
			name: "with limit",
			opts: ListOptions{
				Limit: intPtr(10),
			},
			want: url.Values{
				"_limit": []string{"10"},
			},
		},
		{
			name: "with offset",
			opts: ListOptions{
				Offset: intPtr(5),
			},
			want: url.Values{
				"_offset": []string{"5"},
			},
		},
		{
			name: "with sort",
			opts: ListOptions{
				Sort: strPtr("name:asc"),
			},
			want: url.Values{
				"_sort": []string{"name:asc"},
			},
		},
		{
			name: "with single expand",
			opts: ListOptions{
				Expand: []string{"tags"},
			},
			want: url.Values{
				"_expand": []string{"tags"},
			},
		},
		{
			name: "with multiple expand",
			opts: ListOptions{
				Expand: []string{"tags", "manifest", "layers"},
			},
			want: url.Values{
				"_expand": []string{"tags,manifest,layers"},
			},
		},
		{
			name: "with empty expand slice",
			opts: ListOptions{
				Expand: []string{},
			},
			want: url.Values{},
		},
		{
			name: "with all options",
			opts: ListOptions{
				Limit:  intPtr(20),
				Offset: intPtr(10),
				Sort:   strPtr("created_at:desc"),
				Expand: []string{"tags", "metadata"},
			},
			want: url.Values{
				"_limit":  []string{"20"},
				"_offset": []string{"10"},
				"_sort":   []string{"created_at:desc"},
				"_expand": []string{"tags,metadata"},
			},
		},
		{
			name: "with zero values",
			opts: ListOptions{
				Limit:  intPtr(0),
				Offset: intPtr(0),
				Sort:   strPtr(""),
				Expand: []string{""},
			},
			want: url.Values{
				"_limit":  []string{"0"},
				"_offset": []string{"0"},
				"_sort":   []string{""},
				"_expand": []string{""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreatePaginationParams(tt.opts)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePaginationParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
