package openmeteo

import "net/http"

type httpClient interface {
	Get(string) (*http.Response, error)
}
