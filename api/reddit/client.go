package reddit

import (
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}
