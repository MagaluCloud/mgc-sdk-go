package helpers

import (
	"encoding/json"
	"testing"
)

func TestPaginatedResponse(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    *PaginatedResponse[map[string]interface{}]
		wantErr bool
	}{
		{
			name: "valid paginated response",
			json: `{
				"meta": {
					"page": {
						"count": 2,
						"limit": 20,
						"offset": 0,
						"total": 42
					}
				},
				"results": [
					{"id": "1", "name": "item1"},
					{"id": "2", "name": "item2"}
				]
			}`,
			want: &PaginatedResponse[map[string]interface{}]{
				Meta: PaginatedMeta{
					Page: PaginatedPage{
						Count:  2,
						Limit:  20,
						Offset: 0,
						Total:  42,
					},
				},
				Results: []map[string]interface{}{
					{"id": "1", "name": "item1"},
					{"id": "2", "name": "item2"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty results",
			json: `{
				"meta": {
					"page": {
						"count": 0,
						"limit": 20,
						"offset": 0,
						"total": 0
					}
				},
				"results": []
			}`,
			want: &PaginatedResponse[map[string]interface{}]{
				Meta: PaginatedMeta{
					Page: PaginatedPage{
						Count:  0,
						Limit:  20,
						Offset: 0,
						Total:  0,
					},
				},
				Results: []map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			json:    `{"meta": {"page": {"count": 1}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *PaginatedResponse[map[string]interface{}]
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Meta.Page.Count != tt.want.Meta.Page.Count {
					t.Errorf("Count = %d, want %d", got.Meta.Page.Count, tt.want.Meta.Page.Count)
				}
				if got.Meta.Page.Total != tt.want.Meta.Page.Total {
					t.Errorf("Total = %d, want %d", got.Meta.Page.Total, tt.want.Meta.Page.Total)
				}
				if len(got.Results) != len(tt.want.Results) {
					t.Errorf("Results length = %d, want %d", len(got.Results), len(tt.want.Results))
				}
			}
		})
	}
}
