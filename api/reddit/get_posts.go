package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type SubredditType int

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

type GetPostsParam struct {
	Subreddit     string
	Limit         int
	Page          int
	SubredditType SubredditType
}

func (reddit *Reddit) GetPosts(ctx context.Context, params GetPostsParam) (posts []Post, err error) {
	url := fmt.Sprintf("https://reddit.com/%s/%s.json?limit=%d&page=%d", params.SubredditType.Code(), params.Subreddit, params.Limit, params.Page)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return posts, errs.Wrapw(err, "reddit: failed to create http request instance", "url", url, "params", params)
	}

	res, err := reddit.Client.Do(req)
	if err != nil {
		return posts, errs.Wrapw(err, "reddit: failed to execute http request", "url", url, "params", params)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return posts, errs.Fail("reddit: unexpected status code when executing GetPosts",
			slog.Group("request", "url", url, "params", params),
			slog.Group("response", "status_code", res.StatusCode, "body", json.RawMessage(body)),
		)
	}

	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		return posts, errs.Wrapw(err, "reddit: failed to parse response body when getting posts from reddit", "url", url, "params", params)
	}

	return posts, nil
}
