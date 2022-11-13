package httpclient

import (
	"net/http"

	"gopkg.in/h2non/gentleman.v2"
)

//HTTPClient implements the IHTTPClientProvider
var _ IHTTPClientProvider = &HTTPClient{}

//HTTPClient ..
type HTTPClient struct {
	client *gentleman.Client
}

//New ...
func New() *HTTPClient {
	client := gentleman.New()
	return &HTTPClient{client: client}
}

//Get ...
func (c *HTTPClient) Get(url string) (*http.Response, error) {
	c.client.URL(url)
	c.client.Method("GET")
	response, err := c.client.Request().Send()
	return response.RawResponse, err
}
