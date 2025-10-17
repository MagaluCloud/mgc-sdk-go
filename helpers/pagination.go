package helpers

// PaginatedMeta contains pagination metadata (nested format)
// Used by APIs that have meta.page structure
type PaginatedMeta struct {
	Page PaginatedPage `json:"page"`
}

// PaginatedPage contains page information
type PaginatedPage struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// PaginatedResponse is a generic type for paginated API responses
// Used by APIs that have meta.page structure
type PaginatedResponse[T any] struct {
	Meta    PaginatedMeta `json:"meta"`
	Results []T           `json:"results"`
}

// AuditPaginatedMeta contains pagination metadata (flat format)
// Used by audit APIs that have flat meta structure
type AuditPaginatedMeta struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// AuditPaginatedResponse is a generic type for audit paginated API responses
// Used by audit APIs that have flat meta structure
type AuditPaginatedResponse[T any] struct {
	Meta    AuditPaginatedMeta `json:"meta"`
	Results []T                `json:"results"`
}
