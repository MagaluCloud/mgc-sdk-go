package helpers

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type queryParam struct {
	httpReq *http.Request
	query   *url.Values
}

type QueryParams interface {
	Add(name string, value *string)
	AddReflect(name string, value any)
	Encode() string
}

func NewQueryParams(httpReq *http.Request) QueryParams {
	query := httpReq.URL.Query()
	return &queryParam{
		httpReq: httpReq,
		query:   &query,
	}
}

func (q *queryParam) Add(name string, value *string) {
	if value != nil {
		q.query.Set(name, *value)
	}
}

func (q *queryParam) AddReflect(name string, value any) {
	if value != nil {
		typeOf := reflect.TypeOf(value)
		switch typeOf.Kind() {
		case reflect.String:
			q.query.Set(name, value.(string))
		case reflect.Int:
			q.query.Set(name, strconv.Itoa(value.(int)))
		}
	}
}

func (q *queryParam) Encode() string {
	return q.query.Encode()
}
