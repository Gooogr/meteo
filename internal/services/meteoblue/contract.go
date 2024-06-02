package meteoblue

import "net/http"

type httpClient interface {
	Get(string) (*http.Response, error)
}
