package reddit

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/log"
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
		resp, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			log.New(r.Context()).
				Err(err).
				Errorf("reddit: %s %s", r.Method, r.URL)
			return resp, err
		}
		if resp.StatusCode >= 400 {
			log.New(r.Context()).
				Errorf("reddit: %s %s %d", r.Method, r.URL, resp.StatusCode)
			return resp, err
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.New(r.Context()).
				Infof("reddit: %s %s %d", r.Method, r.URL, resp.StatusCode)
		}
		return resp, err
	}
}
