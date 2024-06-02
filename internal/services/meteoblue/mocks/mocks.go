package mocks

import "net/http"

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return m.DoFunc(req)
}
