package reddit

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/config"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func NewRedditHTTPClient(cfg *config.Config) Client {
	return &http.Client{
		Transport: createRoundTripper(cfg),
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (ro roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return ro(req)
}

func createRoundTripper(cfg *config.Config) roundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		r.Header.Set("User-Agent", cfg.String("download.useragent"))
		return http.DefaultTransport.RoundTrip(r)
	}
}
