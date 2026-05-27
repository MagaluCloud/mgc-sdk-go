package objectstorage

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type objectStorageTransport struct {
	base http.RoundTripper
}

func (t *objectStorageTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodDelete && HasForceDelete(req.Context()) {
		req.Header.Set("X-Force-Container-Delete", "true")
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK || resp.Body == nil {
		return resp, nil
	}

	if req.Method == http.MethodGet && strings.Contains(req.URL.RawQuery, "retention") && HasFixRetentionTime(req.Context()) {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return resp, nil
		}

		resp.Body.Close()

		fixed := fixRetentionTime(body)

		resp.Body = io.NopCloser(bytes.NewReader(fixed))
		resp.ContentLength = int64(len(fixed))
		resp.Header.Set("Content-Length", strconv.Itoa(len(fixed)))
	}

	return resp, nil
}

var tzFix = regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})([+-]\d{2})(\d{2})`)

func fixRetentionTime(body []byte) []byte {

	return tzFix.ReplaceAll(body, []byte("$1$2:$3"))
}
