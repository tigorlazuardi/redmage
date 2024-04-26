package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type SubredditType int

func (su *SubredditType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "null":
		return nil
	case `"user"`, `"u"`, "1":
		*su = SubredditTypeUser
		return nil
	case `"r"`, `"subreddit"`, "0":
		*su = SubredditTypeSub
		return nil
	}
	return errs.
		Fail("subreddit type not recognized. Valid values are 'user', 'u', 'r', 'subreddit', 0, 1, and null",
			"got", string(b),
		).
		Code(http.StatusBadRequest)
}

const (
	SubredditTypeSub SubredditType = iota
	SubredditTypeUser
)

func (s SubredditType) Code() string {
	switch s {
	case SubredditTypeUser:
		return "user"
	default:
		return "r"
	}
}

func (s SubredditType) String() string {
	return s.Code()
}

type GetPostsParam struct {
	Subreddit     string
	Limit         int
	After         string
	SubredditType SubredditType
}

func (reddit *Reddit) GetPosts(ctx context.Context, params GetPostsParam) (posts Listing, err error) {
	ctx, span := tracer.Start(ctx, "*Reddit.GetPosts")
	defer span.End()

	url := fmt.Sprintf("https://reddit.com/%s/%s.json?limit=%d&after=%s", params.SubredditType.Code(), params.Subreddit, params.Limit, params.After)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return posts, errs.Wrapw(err, "reddit: failed to create http request instance", "url", url, "params", params)
	}

	res, err := reddit.Client.Do(req)
	if err != nil {
		return posts, errs.Wrapw(err, "reddit: failed to execute http request", "url", url, "params", params)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusTooManyRequests {
		retryAfter, _ := time.ParseDuration(res.Header.Get("Retry-After"))
		return posts, errs.Fail("reddit: too many requests", "retry_after", retryAfter.String())
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return posts, errs.Fail("reddit: unexpected status code when executing GetPosts",
			slog.Group("request", "url", url, "params", params),
			slog.Group("response", "status_code", res.StatusCode, "body", formatLogBody(res, body)),
		).Code(res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		return posts, errs.Wrapw(err, "reddit: failed to parse response body when getting posts from reddit", "url", url, "params", params)
	}

	return posts, nil
}

func formatLogBody(res *http.Response, body []byte) any {
	if strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		return json.RawMessage(body)
	}
	return string(body)
}
