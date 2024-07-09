package reddit

import (
	"net/http"
	"time"

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
		start := time.Now()
		r.Header.Set("User-Agent", cfg.String("download.useragent"))
		resp, err := http.DefaultTransport.RoundTrip(r)
		end := time.Now()
		duration := end.Sub(start)
		durStr := (duration / time.Millisecond).String()
		if err != nil {
			log.New(r.Context()).
				Err(err).
				Errorf("reddit: %s %s %sms", r.Method, r.URL, durStr)
			return resp, err
		}
		if resp.StatusCode >= 400 {
			log.New(r.Context()).
				Errorf("reddit: %s %s %sms", r.Method, r.URL, durStr)
			return resp, err
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.New(r.Context()).
				Infof("reddit: %s %s %d %sms", r.Method, r.URL, resp.StatusCode, durStr)
		}
		return resp, err
	}
}
