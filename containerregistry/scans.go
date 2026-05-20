package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

type (
	// ScansService manages the image scan lifecycle: scheduling, listing,
	// inspecting, stopping and reading vulnerabilities.
	ScansService interface {
		Schedule(ctx context.Context, in ScheduleScanInput) (*ImageScanScheduleResponse, error)
		ListByImage(ctx context.Context, in ListByImageInput, opts ListImageScansOptions) (*ListImageScansResponse, error)
		ListByImageAll(ctx context.Context, in ListByImageInput, filterOpts ImageScansFilterOptions) ([]ImageScanResponse, error)
		Get(ctx context.Context, scanID string) (*ImageScanResponse, error)
		ListVulnerabilities(ctx context.Context, scanID string, opts ListVulnerabilitiesOptions) (*ListVulnerabilitiesResponse, error)
		ListVulnerabilitiesAll(ctx context.Context, scanID string, filterOpts VulnerabilitiesFilterOptions) ([]VulnerabilityResponse, error)
		Stop(ctx context.Context, scanID string) error
	}

	// ScheduleScanInput identifies the image whose scan should be scheduled.
	ScheduleScanInput struct {
		RegistryID   string
		RepositoryID string
		DigestOrTag  string
	}

	// ListByImageInput identifies the image whose scans should be listed.
	ListByImageInput struct {
		RegistryID   string
		RepositoryID string
		DigestOrTag  string
	}

	// ListImageScansOptions provides pagination and filters for image scans.
	ListImageScansOptions struct {
		Offset *int
		Limit  *int
		ImageScansFilterOptions
	}

	// ImageScansFilterOptions provides filters that persist across paginated calls.
	ImageScansFilterOptions struct {
		Sort   *string
		Status *string
	}

	// ListVulnerabilitiesOptions provides pagination and filters for vulnerabilities.
	ListVulnerabilitiesOptions struct {
		Offset *int
		Limit  *int
		VulnerabilitiesFilterOptions
	}

	// VulnerabilitiesFilterOptions provides filters that persist across paginated calls.
	VulnerabilitiesFilterOptions struct {
		Sort        *string
		Severity    []string
		PackageName *string
		CveID       *string
		Fixable     *bool
	}

	// SeverityResponse summarizes vulnerability counts per severity.
	SeverityResponse struct {
		Total    int `json:"total"`
		Low      int `json:"low"`
		Medium   int `json:"medium"`
		High     int `json:"high"`
		Critical int `json:"critical"`
		Fixable  int `json:"fixable"`
	}

	// ImageScanScheduleResponse is returned when a scan is scheduled.
	ImageScanScheduleResponse struct {
		ID           string  `json:"id"`
		Digest       string  `json:"digest"`
		RegistryID   string  `json:"registry_id"`
		RepositoryID string  `json:"repository_id"`
		InputTag     *string `json:"input_tag,omitempty"`
		Status       string  `json:"status"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    string  `json:"updated_at"`
	}

	// ChildImageScanResponse describes a scan for a specific architecture.
	ChildImageScanResponse struct {
		ID              string            `json:"id"`
		Digest          string            `json:"digest"`
		InputTag        *string           `json:"input_tag,omitempty"`
		Status          string            `json:"status"`
		OS              *string           `json:"os,omitempty"`
		Architecture    *string           `json:"architecture,omitempty"`
		SeveritySummary *SeverityResponse `json:"severity_summary,omitempty"`
		CreatedAt       string            `json:"created_at"`
		UpdatedAt       string            `json:"updated_at"`
		FinishedAt      *string           `json:"finished_at,omitempty"`
	}

	// ImageScanResponse describes a scan and (optionally) its per-arch child scans.
	ImageScanResponse struct {
		ID              string                   `json:"id"`
		Digest          string                   `json:"digest"`
		InputTag        *string                  `json:"input_tag,omitempty"`
		Status          string                   `json:"status"`
		OS              *string                  `json:"os,omitempty"`
		Architecture    *string                  `json:"architecture,omitempty"`
		SeveritySummary *SeverityResponse        `json:"severity_summary,omitempty"`
		CreatedAt       string                   `json:"created_at"`
		UpdatedAt       string                   `json:"updated_at"`
		FinishedAt      *string                  `json:"finished_at,omitempty"`
		ChildScans      []ChildImageScanResponse `json:"child_scans,omitempty"`
	}

	// VulnerabilityCvssSourceResponse holds CVSS data from a single source.
	VulnerabilityCvssSourceResponse struct {
		Score  *float64 `json:"score"`
		Vector *string  `json:"vector"`
	}

	// VulnerabilityCvssMaxResponse holds the highest CVSS data among sources.
	VulnerabilityCvssMaxResponse struct {
		Source *string  `json:"source"`
		Score  *float64 `json:"score"`
		Vector *string  `json:"vector"`
	}

	// VulnerabilityCvssResponse aggregates CVSS data across sources.
	VulnerabilityCvssResponse struct {
		Preferred *VulnerabilityCvssSourceResponse           `json:"preferred"`
		Max       *VulnerabilityCvssMaxResponse              `json:"max"`
		BySource  map[string]VulnerabilityCvssSourceResponse `json:"by_source"`
	}

	// VulnerabilityResponse describes a single vulnerability found in a scan.
	VulnerabilityResponse struct {
		CveID          string                    `json:"cve_id"`
		Severity       string                    `json:"severity"`
		Cvss           VulnerabilityCvssResponse `json:"cvss"`
		IsAllowlisted  bool                      `json:"is_allowlisted"`
		PackageName    string                    `json:"package_name"`
		CurrentVersion string                    `json:"current_version"`
		FixedVersion   *string                   `json:"fixed_version,omitempty"`
		Fixable        bool                      `json:"fixable"`
		TargetPath     string                    `json:"target_path,omitempty"`
		Description    string                    `json:"description"`
		Links          []string                  `json:"links"`
	}

	// ListImageScansResponse is the paginated response of image scans.
	ListImageScansResponse = helpers.PaginatedResponse[ImageScanResponse]

	// ListVulnerabilitiesResponse is the paginated response of vulnerabilities.
	ListVulnerabilitiesResponse = helpers.PaginatedResponse[VulnerabilityResponse]

	scansService struct {
		client *ContainerRegistryClient
	}
)

// Schedule queues a scan for the image identified by digest or tag.
func (s *scansService) Schedule(ctx context.Context, in ScheduleScanInput) (*ImageScanScheduleResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ImageScanScheduleResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodPost, scansByImagePath(in.RegistryID, in.RepositoryID, in.DigestOrTag), nil, nil,
	)
}

// ListByImage returns a page of scans for the given image.
func (s *scansService) ListByImage(ctx context.Context, in ListByImageInput, opts ListImageScansOptions) (*ListImageScansResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ListImageScansResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodGet, scansByImagePath(in.RegistryID, in.RepositoryID, in.DigestOrTag), nil,
		imageScansQuery(opts),
	)
}

// ListByImageAll walks every page until all scans for the image are retrieved.
func (s *scansService) ListByImageAll(ctx context.Context, in ListByImageInput, filterOpts ImageScansFilterOptions) ([]ImageScanResponse, error) {
	var all []ImageScanResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		page, err := s.ListByImage(ctx, in, ListImageScansOptions{
			Offset:                  &currentOffset,
			Limit:                   &currentLimit,
			ImageScansFilterOptions: filterOpts,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, page.Results...)
		if len(page.Results) < limit {
			break
		}
		offset += limit
	}
	return all, nil
}

// Get returns the details of a specific scan.
func (s *scansService) Get(ctx context.Context, scanID string) (*ImageScanResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ImageScanResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodGet, fmt.Sprintf("/v1/scans/%s", scanID), nil, nil,
	)
}

// ListVulnerabilities returns a page of vulnerabilities found in a scan.
func (s *scansService) ListVulnerabilities(ctx context.Context, scanID string, opts ListVulnerabilitiesOptions) (*ListVulnerabilitiesResponse, error) {
	return mgc_http.ExecuteSimpleRequestWithRespBody[ListVulnerabilitiesResponse](
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodGet, fmt.Sprintf("/v1/scans/%s/vulnerabilities", scanID), nil,
		vulnerabilitiesQuery(opts),
	)
}

// ListVulnerabilitiesAll walks every page until all vulnerabilities are retrieved.
func (s *scansService) ListVulnerabilitiesAll(ctx context.Context, scanID string, filterOpts VulnerabilitiesFilterOptions) ([]VulnerabilityResponse, error) {
	var all []VulnerabilityResponse
	offset := 0
	limit := 50

	for {
		currentOffset := offset
		currentLimit := limit
		page, err := s.ListVulnerabilities(ctx, scanID, ListVulnerabilitiesOptions{
			Offset:                       &currentOffset,
			Limit:                        &currentLimit,
			VulnerabilitiesFilterOptions: filterOpts,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, page.Results...)
		if len(page.Results) < limit {
			break
		}
		offset += limit
	}
	return all, nil
}

// Stop cancels an ongoing scan.
func (s *scansService) Stop(ctx context.Context, scanID string) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx, s.client.newRequest, s.client.GetConfig(),
		http.MethodPost, fmt.Sprintf("/v1/scans/%s/stop", scanID), nil, nil,
	)
}

func scansByImagePath(registryID, repositoryID, digestOrTag string) string {
	return fmt.Sprintf("/v1/registries/%s/repositories/%s/images/%s/scans",
		registryID, repositoryID, digestOrTag)
}

func imageScansQuery(opts ListImageScansOptions) url.Values {
	q := make(url.Values)
	if opts.Limit != nil {
		q.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Set("_sort", *opts.Sort)
	}
	if opts.Status != nil {
		q.Set("status", *opts.Status)
	}
	return q
}

func vulnerabilitiesQuery(opts ListVulnerabilitiesOptions) url.Values {
	q := make(url.Values)
	if opts.Limit != nil {
		q.Set("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Set("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Set("_sort", *opts.Sort)
	}
	for _, sev := range opts.Severity {
		q.Add("severity", sev)
	}
	if opts.PackageName != nil {
		q.Set("package_name", *opts.PackageName)
	}
	if opts.CveID != nil {
		q.Set("cve_id", *opts.CveID)
	}
	if opts.Fixable != nil {
		q.Set("fixable", strconv.FormatBool(*opts.Fixable))
	}
	return q
}
