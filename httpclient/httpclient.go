package httpclient

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	defaultClient *http.Client
	timeout       time.Duration
	proxy         *url.URL
}

func New(defaultClient *http.Client, proxy *url.URL, timeout time.Duration) *Client {
	return &Client{
		defaultClient: defaultClient,
		proxy:         proxy,
		timeout:       timeout,
	}
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	c.defaultClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(c.proxy),
	}

	c.defaultClient.Timeout = c.timeout

	return c.defaultClient.Do(r)
}
