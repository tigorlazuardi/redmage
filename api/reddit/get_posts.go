package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type SubredditType int

func (su *SubredditType) UnmarshalJSON(b []byte) error {
	if len(b) == 4 && string(b) == "null" {
		return nil
	}
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return errs.Wrapw(err, "failed to unquote string json value").Code(http.StatusBadRequest)
	}
	return su.Parse(s)
}

func (su *SubredditType) Parse(s string) error {
	switch s {
	case `"user"`, `"u"`, "1":
		*su = SubredditTypeUser
		return nil
	case `"r"`, `"subreddit"`, "0", "":
		*su = SubredditTypeSub
		return nil
	}
	return errs.
		Fail("subreddit type not recognized. Valid values are '' (empty), 'user', 'u', 'r', 'subreddit', 0, 1, and null",
			"got", s,
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

	req.Header.Set("User-Agent", reddit.Config.String("download.useragent"))

	res, err := reddit.Client.Do(req)
	if err != nil {
		return posts, errs.Wrapw(err, "reddit: failed to execute http request", "url", url, "params", params)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusTooManyRequests {
		return posts, errs.Fail("reddit: too many requests",
			"retry_after", res.Header.Get("Retry-After"),
			"url", url,
		)
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
