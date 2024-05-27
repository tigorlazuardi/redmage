package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tigorlazuardi/redmage/pkg/caller"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type CheckSubredditParams struct {
	Subreddit     string        `json:"subreddit"`
	SubredditType SubredditType `json:"subreddit_type"`
}

// CheckSubreddit checks a subreddit existence and will return error if subreddit not found.
//
// The actual is the subreddit with proper capitalization if no error is returned.
func (reddit *Reddit) CheckSubreddit(ctx context.Context, params CheckSubredditParams) (actual string, err error) {
	ctx, span := tracer.Start(ctx, "*Reddit.CheckSubreddit")
	defer span.End()

	ctx = caller.WithContext(ctx, caller.New(2))

	url := fmt.Sprintf("https://reddit.com/%s/%s.json?limit=1", params.SubredditType.Code(), params.Subreddit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return actual, errs.Wrapw(err, "failed to create request", "url", url, "params", params)
	}
	req.Header.Set("User-Agent", reddit.Config.String("download.useragent"))

	resp, err := reddit.Client.Do(req)
	if err != nil {
		return actual, errs.Wrapw(err, "failed to execute request", "url", url, "params", params)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		var msg string
		if params.SubredditType == SubredditTypeUser {
			msg = "user not found"
		}
		if params.SubredditType == SubredditTypeSub {
			msg = "subreddit not found"
		}
		return actual, errs.Wrapw(err, msg, "url", url, "params", params).Code(http.StatusNotFound)
	}

	if resp.StatusCode == http.StatusForbidden {
		var msg string
		if params.SubredditType == SubredditTypeUser {
			msg = "user has set their profile to private"
		}
		if params.SubredditType == SubredditTypeSub {
			msg = "subreddit is private"
		}
		return actual, errs.Wrapw(err, msg, "url", url, "params", params).Code(http.StatusForbidden)
	}

	if params.SubredditType == SubredditTypeUser && resp.StatusCode == http.StatusOK {
		return params.Subreddit, nil
	}

	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("unexpected %d status code from reddit", resp.StatusCode)
		return actual, errs.
			Fail(msg, "url", url, "params", params, "response.status", resp.StatusCode).
			Code(http.StatusFailedDependency)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		var msg string
		dur, _ := time.ParseDuration(resp.Header.Get("Retry-After") + "s")
		if dur > 0 {
			msg = fmt.Sprintf("too many requests. Please retry after %s", dur)
		} else {
			msg = "too many requests. Please try again later"
		}
		return actual, errs.Fail(msg,
			"params", params,
			"url", url,
			"response.location", resp.Request.URL.String(),
		).Code(http.StatusTooManyRequests)
	}

	if resp.Request.URL.Path == "/subreddits/search.json" {
		return actual, errs.Fail("subreddit not found",
			"params", params,
			"url", url,
			"response.location", resp.Request.URL.String(),
		).Code(http.StatusNotFound)
	}

	var body Listing
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return actual, errs.Wrapw(err, "failed to decode json body")
	}
	sub := body.GetSubreddit()
	if sub == "" {
		return actual, errs.Fail("subreddit not found",
			"params", params,
			"url", url,
			"response.location", resp.Request.URL.String(),
		).Code(http.StatusNotFound)
	}

	return sub, nil
}
