package client

import "net/http"

const ProjectIDHeader = "X-Project-Id"

type CallOptions struct {
	ProjectID string
}

func (o CallOptions) Apply(req *http.Request) {
	if req == nil {
		return
	}

	if o.ProjectID != "" {
		req.Header.Set(ProjectIDHeader, o.ProjectID)
	}
}
