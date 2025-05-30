package lbaas

import (
	"context"
	"net/http"
	"strconv"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

const listeners = "listeners"

type (
	CreateNetworkListenerRequest struct {
		LoadBalancerID   string           `json:"-"`
		BackendID        string           `json:"-"`
		TLSCertificateID *string          `json:"tls_certificate_id,omitempty"`
		Name             string           `json:"name"`
		Description      *string          `json:"description,omitempty"`
		Protocol         ListenerProtocol `json:"protocol"`
		Port             int              `json:"port"`
	}

	DeleteNetworkListenerRequest struct {
		LoadBalancerID string `json:"-"`
		ListenerID     string `json:"-"`
	}

	GetNetworkListenerRequest struct {
		LoadBalancerID string `json:"-"`
		ListenerID     string `json:"-"`
	}

	ListNetworkListenerRequest struct {
		LoadBalancerID string  `json:"-"`
		Offset         *int    `json:"-"`
		Limit          *int    `json:"-"`
		Sort           *string `json:"-"`
	}

	UpdateNetworkListenerRequest struct {
		LoadBalancerID   string  `json:"-"`
		ListenerID       string  `json:"-"`
		TLSCertificateID *string `json:"tls_certificate_id,omitempty"`
	}

	NetworkListenerResponse struct {
		ID               string           `json:"id"`
		TLSCertificateID *string          `json:"tls_certificate_id,omitempty"`
		BackendID        string           `json:"backend_id"`
		Name             string           `json:"name"`
		Description      *string          `json:"description,omitempty"`
		Protocol         ListenerProtocol `json:"protocol"`
		Port             int              `json:"port"`
		CreatedAt        string           `json:"created_at"`
		UpdatedAt        string           `json:"updated_at"`
	}

	NetworkPaginatedListenerResponse struct {
		Meta    interface{}               `json:"meta"`
		Results []NetworkListenerResponse `json:"results"`
	}

	NetworkListenerService interface {
		Create(ctx context.Context, req CreateNetworkListenerRequest) (*NetworkListenerResponse, error)
		Delete(ctx context.Context, req DeleteNetworkListenerRequest) error
		Get(ctx context.Context, req GetNetworkListenerRequest) (*NetworkListenerResponse, error)
		List(ctx context.Context, req ListNetworkListenerRequest) ([]NetworkListenerResponse, error)
		Update(ctx context.Context, req UpdateNetworkListenerRequest) error
	}

	networkListenerService struct {
		client *LbaasClient
	}
)

func (s *networkListenerService) Create(ctx context.Context, req CreateNetworkListenerRequest) (*NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners)

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}

	// Adicionar backend_id como query parameter obrigatório
	query := httpReq.URL.Query()
	query.Set("backend_id", req.BackendID)
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *networkListenerService) Delete(ctx context.Context, req DeleteNetworkListenerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners, req.ListenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}

func (s *networkListenerService) Get(ctx context.Context, req GetNetworkListenerRequest) (*NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners, req.ListenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp NetworkListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *networkListenerService) List(ctx context.Context, req ListNetworkListenerRequest) ([]NetworkListenerResponse, error) {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners)

	httpReq, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Adicionar query parameters se fornecidos
	query := httpReq.URL.Query()
	if req.Offset != nil {
		query.Set("_offset", strconv.Itoa(*req.Offset))
	}
	if req.Limit != nil {
		query.Set("_limit", strconv.Itoa(*req.Limit))
	}
	if req.Sort != nil {
		query.Set("_sort", *req.Sort)
	}
	httpReq.URL.RawQuery = query.Encode()

	var resp NetworkPaginatedListenerResponse
	result, err := mgc_http.Do(s.client.GetConfig(), ctx, httpReq, &resp)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (s *networkListenerService) Update(ctx context.Context, req UpdateNetworkListenerRequest) error {
	path := urlNetworkLoadBalancer(&req.LoadBalancerID, listeners, req.ListenerID)

	httpReq, err := s.client.newRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return err
	}

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, httpReq, nil)
	return err
}
