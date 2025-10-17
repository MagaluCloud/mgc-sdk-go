package containerregistry

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestCreateImageQueryParams(t *testing.T) {
	tests := []struct {
		name string
		opts ImageListOptions
		want url.Values
	}{
		{
			name: "empty options",
			opts: ImageListOptions{},
			want: url.Values{},
		},
		{
			name: "with limit",
			opts: ImageListOptions{
				Limit: helpers.IntPtr(10),
			},
			want: url.Values{
				"_limit": []string{"10"},
			},
		},
		{
			name: "with offset",
			opts: ImageListOptions{
				Offset: helpers.IntPtr(5),
			},
			want: url.Values{
				"_offset": []string{"5"},
			},
		},
		{
			name: "with sort",
			opts: ImageListOptions{
				ImageFilterOptions: ImageFilterOptions{
					Sort: helpers.StrPtr("name:asc"),
				},
			},
			want: url.Values{
				"_sort": []string{"name:asc"},
			},
		},
		{
			name: "with single expand",
			opts: ImageListOptions{
				ImageFilterOptions: ImageFilterOptions{
					Expand: []ImageExpand{ImageTagsDetailsExpand},
				},
			},
			want: url.Values{
				"_expand": []string{"tags_details"},
			},
		},
		{
			name: "with multiple expand",
			opts: ImageListOptions{
				ImageFilterOptions: ImageFilterOptions{
					Expand: []ImageExpand{ImageTagsDetailsExpand, ImageExtraAttrExpand, ImageMediaTypeExpand},
				},
			},
			want: url.Values{
				"_expand": []string{"tags_details,extra_attr,media_type"},
			},
		},
		{
			name: "with empty expand slice",
			opts: ImageListOptions{
				ImageFilterOptions: ImageFilterOptions{
					Expand: []ImageExpand{},
				},
			},
			want: url.Values{},
		},
		{
			name: "with all options",
			opts: ImageListOptions{
				Limit:  helpers.IntPtr(20),
				Offset: helpers.IntPtr(10),
				ImageFilterOptions: ImageFilterOptions{
					Sort:   helpers.StrPtr("created_at:desc"),
					Expand: []ImageExpand{ImageTagsDetailsExpand, ImageExtraAttrExpand},
				},
			},
			want: url.Values{
				"_limit":  []string{"20"},
				"_offset": []string{"10"},
				"_sort":   []string{"created_at:desc"},
				"_expand": []string{"tags_details,extra_attr"},
			},
		},
		{
			name: "with zero limit and offset",
			opts: ImageListOptions{
				Limit:  helpers.IntPtr(0),
				Offset: helpers.IntPtr(0),
			},
			want: url.Values{
				"_limit":  []string{"0"},
				"_offset": []string{"0"},
			},
		},
		{
			name: "with empty sort string",
			opts: ImageListOptions{
				ImageFilterOptions: ImageFilterOptions{
					Sort: helpers.StrPtr(""),
				},
			},
			want: url.Values{
				"_sort": []string{""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &imagesService{}
			got := service.createImageQueryParams(tt.opts)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createImageQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
