package objectstorage

import (
	"net/http"
)

type forceDeleteTransport struct {
	base http.RoundTripper
}

func (t *forceDeleteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodDelete && HasForceDelete(req.Context()) {
		req.Header.Set("X-Force-Container-Delete", "true")
	}

	return t.base.RoundTrip(req)
}
