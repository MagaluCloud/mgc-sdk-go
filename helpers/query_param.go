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
		valueOf := reflect.ValueOf(value)
		typeOf := reflect.TypeOf(value)

		// Handle pointer types by dereferencing them
		if typeOf.Kind() == reflect.Ptr && !valueOf.IsNil() {
			valueOf = valueOf.Elem()
			typeOf = valueOf.Type()
		}

		switch typeOf.Kind() {
		case reflect.String:
			q.query.Set(name, valueOf.String())
		case reflect.Int:
			q.query.Set(name, strconv.FormatInt(valueOf.Int(), 10))
		}
	}
}

func (q *queryParam) Encode() string {
	return q.query.Encode()
}
