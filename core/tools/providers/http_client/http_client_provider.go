package httpclient

import "net/http"

//IHTTPClientProvider ...
type IHTTPClientProvider interface {
	Get(url string) (*http.Response, error)
}
