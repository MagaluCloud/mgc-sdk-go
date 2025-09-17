package blockstorage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func TestSchedulerService_List(t *testing.T) {
	tests := []struct {
		name       string
		opts       SchedulerListOptions
		response   string
		statusCode int
		want       int
		wantErr    bool
	}{
		{
			name: "basic list",
			opts: SchedulerListOptions{},
			response: `{
				"meta": {
					"page": {
						"offset": 0,
						"limit": 50,
						"count": 2,
						"total": 2,
						"max_limit": 100
					}
				},
				"schedulers": [
					{
						"id": "scheduler1",
						"name": "test-scheduler-1",
						"state": "available",
						"policy": {
							"retention_in_days": 7,
							"frequency": {
								"daily": {
									"start_time": "02:00:00"
								}
							}
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					{
						"id": "scheduler2",
						"name": "test-scheduler-2",
						"state": "available",
						"policy": {
							"retention_in_days": 30,
							"frequency": {
								"daily": {
									"start_time": "03:00:00"
								}
							}
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want:       2,
			wantErr:    false,
		},
		{
			name: "with pagination",
			opts: SchedulerListOptions{
				Limit:  helpers.IntPtr(1),
				Offset: helpers.IntPtr(1),
			},
			response: `{
				"meta": {
					"page": {
						"offset": 1,
						"limit": 1,
						"count": 1,
						"total": 2,
						"max_limit": 100
					}
				},
				"schedulers": [
					{
						"id": "scheduler2",
						"name": "test-scheduler-2",
						"state": "available",
						"policy": {
							"retention_in_days": 30,
							"frequency": {
								"daily": {
									"start_time": "03:00:00"
								}
							}
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name: "with expansion",
			opts: SchedulerListOptions{
				Expand: []ExpandSchedulers{ExpandSchedulersVolume},
			},
			response: `{
				"meta": {
					"page": {
						"offset": 0,
						"limit": 50,
						"count": 1,
						"total": 1,
						"max_limit": 100
					}
				},
				"schedulers": [
					{
						"id": "scheduler1",
						"name": "test-scheduler-1",
						"state": "available",
						"volumes": ["volume1", "volume2"],
						"policy": {
							"retention_in_days": 7,
							"frequency": {
								"daily": {
									"start_time": "02:00:00"
								}
							}
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name: "with sorting",
			opts: SchedulerListOptions{
				Sort: helpers.StrPtr("created_at:desc"),
			},
			response: `{
				"meta": {
					"page": {
						"offset": 0,
						"limit": 50,
						"count": 1,
						"total": 1,
						"max_limit": 100
					}
				},
				"schedulers": [
					{
						"id": "scheduler1",
						"name": "test-scheduler-1",
						"state": "available",
						"policy": {
							"retention_in_days": 7,
							"frequency": {
								"daily": {
									"start_time": "02:00:00"
								}
							}
						},
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				]
			}`,
			statusCode: http.StatusOK,
			want:       1,
			wantErr:    false,
		},
		{
			name:       "server error",
			response:   `{"error": "internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/schedulers", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSchedulerClient(server.URL)
			result, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want, len(result.Schedulers))
		})
	}
}

func TestSchedulerService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    SchedulerPayload
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			request: SchedulerPayload{
				Name:        "test-scheduler",
				Description: helpers.StrPtr("Test scheduler description"),
				Snapshot: SnapshotConfig{
					Type: "instant",
				},
				Policy: Policy{
					RetentionInDays: 7,
					Frequency: Frequency{
						Daily: DailyFrequency{
							StartTime: "02:00:00",
						},
					},
				},
			},
			response:   `{"id": "scheduler1"}`,
			statusCode: http.StatusCreated,
			wantID:     "scheduler1",
			wantErr:    false,
		},
		{
			name: "invalid retention days",
			request: SchedulerPayload{
				Name: "test-scheduler",
				Policy: Policy{
					RetentionInDays: 0,
				},
			},
			response:   `{"error": "invalid retention days"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid time format",
			request: SchedulerPayload{
				Name: "test-scheduler",
				Policy: Policy{
					RetentionInDays: 7,
					Frequency: Frequency{
						Daily: DailyFrequency{
							StartTime: "invalid-time",
						},
					},
				},
			},
			response:   `{"error": "invalid time format"}`,
			statusCode: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/schedulers", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req SchedulerPayload
				json.NewDecoder(r.Body).Decode(&req)
				assertEqual(t, tt.request.Name, req.Name)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSchedulerClient(server.URL)
			id, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.wantID, id)
		})
	}
}

func TestSchedulerService_Get(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")

	tests := []struct {
		name       string
		id         string
		expand     []ExpandSchedulers
		response   string
		statusCode int
		want       *SchedulerResponse
		wantErr    bool
	}{
		{
			name: "existing scheduler",
			id:   "scheduler1",
			response: `{
				"id": "scheduler1",
				"name": "test-scheduler",
				"description": "Test description",
				"state": "available",
				"policy": {
					"retention_in_days": 7,
					"frequency": {
						"daily": {
							"start_time": "02:00:00"
						}
					}
				},
				"snapshot": {
					"type": "instant"
				},
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			want: &SchedulerResponse{
				ID:          "scheduler1",
				Name:        "test-scheduler",
				Description: helpers.StrPtr("Test description"),
				State:       SchedulerStateAvailable,
				Policy: Policy{
					RetentionInDays: 7,
					Frequency: Frequency{
						Daily: DailyFrequency{
							StartTime: "02:00:00",
						},
					},
				},
				Snapshot: &SnapshotConfig{
					Type: "instant",
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			wantErr: false,
		},
		{
			name:       "not found",
			id:         "invalid",
			response:   `{"error": "not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:   "with expansion",
			id:     "scheduler1",
			expand: []ExpandSchedulers{ExpandSchedulersVolume},
			response: `{
				"id": "scheduler1",
				"name": "test-scheduler",
				"state": "available",
				"volumes": ["volume1", "volume2"],
				"policy": {
					"retention_in_days": 7,
					"frequency": {
						"daily": {
							"start_time": "02:00:00"
						}
					}
				},
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}`,
			statusCode: http.StatusOK,
			want: &SchedulerResponse{
				ID:      "scheduler1",
				Name:    "test-scheduler",
				State:   SchedulerStateAvailable,
				Volumes: []string{"volume1", "volume2"},
				Policy: Policy{
					RetentionInDays: 7,
					Frequency: Frequency{
						Daily: DailyFrequency{
							StartTime: "02:00:00",
						},
					},
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/schedulers/"+tt.id, r.URL.Path)
				assertEqual(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSchedulerClient(server.URL)
			scheduler, err := client.Get(context.Background(), tt.id, tt.expand)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.want.ID, scheduler.ID)
			assertEqual(t, tt.want.Name, scheduler.Name)
			assertEqual(t, tt.want.State, scheduler.State)
			if len(tt.want.Volumes) > 0 {
				assertEqual(t, len(tt.want.Volumes), len(scheduler.Volumes))
			}
		})
	}
}

func TestSchedulerService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "successful delete",
			id:         "scheduler1",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "not found",
			id:         "invalid",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "scheduler in use",
			id:         "scheduler-in-use",
			statusCode: http.StatusConflict,
			response:   `{"error": "scheduler has attached volumes"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/schedulers/"+tt.id, r.URL.Path)
				assertEqual(t, http.MethodDelete, r.Method)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSchedulerClient(server.URL)
			err := client.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestSchedulerService_AttachVolume(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    SchedulerVolumeIdentifierPayload
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "successful attach by ID",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					ID: helpers.StrPtr("volume1"),
				},
			},
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name: "successful attach by name",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					Name: helpers.StrPtr("test-volume"),
				},
			},
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name: "volume not found",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					ID: helpers.StrPtr("invalid-volume"),
				},
			},
			statusCode: http.StatusNotFound,
			response:   `{"error": "volume not found"}`,
			wantErr:    true,
		},
		{
			name: "volume already attached",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					ID: helpers.StrPtr("attached-volume"),
				},
			},
			statusCode: http.StatusConflict,
			response:   `{"error": "volume already attached"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/schedulers/"+tt.id+"/attach", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req SchedulerVolumeIdentifierPayload
				json.NewDecoder(r.Body).Decode(&req)
				if tt.request.Volume.ID != nil && req.Volume.ID != nil {
					assertEqual(t, *tt.request.Volume.ID, *req.Volume.ID)
				}
				if tt.request.Volume.Name != nil && req.Volume.Name != nil {
					assertEqual(t, *tt.request.Volume.Name, *req.Volume.Name)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSchedulerClient(server.URL)
			err := client.AttachVolume(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestSchedulerService_DetachVolume(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		request    SchedulerVolumeIdentifierPayload
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "successful detach by ID",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					ID: helpers.StrPtr("volume1"),
				},
			},
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name: "successful detach by name",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					Name: helpers.StrPtr("test-volume"),
				},
			},
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name: "volume not found",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					ID: helpers.StrPtr("invalid-volume"),
				},
			},
			statusCode: http.StatusNotFound,
			response:   `{"error": "volume not found"}`,
			wantErr:    true,
		},
		{
			name: "volume not attached",
			id:   "scheduler1",
			request: SchedulerVolumeIdentifierPayload{
				Volume: IDOrName{
					ID: helpers.StrPtr("unattached-volume"),
				},
			},
			statusCode: http.StatusConflict,
			response:   `{"error": "volume not attached to scheduler"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "/volume/v1/schedulers/"+tt.id+"/detach", r.URL.Path)
				assertEqual(t, http.MethodPost, r.Method)

				var req SchedulerVolumeIdentifierPayload
				json.NewDecoder(r.Body).Decode(&req)
				if tt.request.Volume.ID != nil && req.Volume.ID != nil {
					assertEqual(t, *tt.request.Volume.ID, *req.Volume.ID)
				}
				if tt.request.Volume.Name != nil && req.Volume.Name != nil {
					assertEqual(t, *tt.request.Volume.Name, *req.Volume.Name)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testSchedulerClient(server.URL)
			err := client.DetachVolume(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assertError(t, err)
				assertEqual(t, true, strings.Contains(err.Error(), strconv.Itoa(tt.statusCode)))
				return
			}

			assertNoError(t, err)
		})
	}
}

// Helper function to create a test scheduler client
func testSchedulerClient(baseURL string) SchedulerService {
	httpClient := &http.Client{}
	core := client.NewMgcClient("test-api",
		client.WithBaseURL(client.MgcUrl(baseURL)),
		client.WithHTTPClient(httpClient))
	return New(core).Schedulers()
}
