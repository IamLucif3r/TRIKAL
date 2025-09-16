package httpx

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	*http.Client
}

func New(timeout time.Duration, ua string, maxIdlePerHost int) *Client {
	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: maxIdlePerHost,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
	}
	c := &http.Client{
		Transport: uaRoundTripper{base: tr, ua: ua},
		Timeout:   timeout,
	}
	return &Client{c}
}

type uaRoundTripper struct {
	base http.RoundTripper
	ua   string
}

func (u uaRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if u.ua != "" {
		req = req.Clone(req.Context())
		req.Header.Set("User-Agent", u.ua)
	}
	return u.base.RoundTrip(req)
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) CloseIdle() {
	if tr, ok := c.Transport.(*http.Transport); ok {
		tr.CloseIdleConnections()
	}
}
