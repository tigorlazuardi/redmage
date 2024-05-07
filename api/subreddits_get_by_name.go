package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) SubredditsGetByName(ctx context.Context, name string) (subreddit *models.Subreddit, err error) {
	ctx, span := tracer.Start(ctx, "*API.SubredditsGetByName")
	defer span.End()

	subreddit, err = models.FindSubreddit(ctx, api.db, name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errs.Wrapw(err, "subreddit not found", "name", name).Code(http.StatusNotFound)
		}
		return nil, errs.Wrapw(err, "failed to get subreddit by name")
	}

	return subreddit, nil
}

type SubredditGetByNameImageParams struct {
	Q       string
	Limit   int64
	Offset  int64
	OrderBy string
	Sort    string
	SFW     int
	After   time.Time
	Device  string
}

func (sgb SubredditGetByNameImageParams) IntoQuery() url.Values {
	queries := make(url.Values)

	if sgb.Q != "" {
		queries.Set("q", sgb.Q)
	}
	if sgb.Limit > 0 {
		queries.Set("limit", strconv.FormatInt(sgb.Limit, 10))
	}
	if sgb.Offset > 0 {
		queries.Set("offset", strconv.FormatInt(sgb.Offset, 10))
	}
	if sgb.OrderBy != "" {
		queries.Set("order_by", sgb.OrderBy)
	}
	if sgb.Sort != "" {
		queries.Set("sort", sgb.Sort)
	}
	if !sgb.After.IsZero() {
		queries.Set("after", strconv.FormatInt(sgb.After.Unix(), 10))
	}
	if sgb.Device != "" {
		queries.Set("device", sgb.Device)
	}

	return queries
}

func (sgb SubredditGetByNameImageParams) IntoQueryWith(keyValue ...string) url.Values {
	queries := sgb.IntoQuery()
	for i := 0; i < len(keyValue); i += 2 {
		queries.Set(keyValue[i], keyValue[i+1])
	}
	return queries
}

func (sgb *SubredditGetByNameImageParams) FillFromQuery(query Queryable) {
	sgb.Q = query.Get("q")
	sgb.Device = query.Get("device")
	sgb.Limit, _ = strconv.ParseInt(query.Get("limit"), 10, 64)
	if sgb.Limit < 1 {
		sgb.Limit = 25
	} else if sgb.Limit > 100 {
		sgb.Limit = 100
	}

	sgb.Offset, _ = strconv.ParseInt(query.Get("offset"), 10, 64)
	if sgb.Offset < 0 {
		sgb.Offset = 0
	}

	sgb.OrderBy = query.Get("order_by")
	sgb.Sort = strings.ToLower(query.Get("sort"))

	afterint, _ := strconv.ParseInt(query.Get("after"), 10, 64)
	if afterint > 0 {
		sgb.After = time.Unix(afterint, 0)
	} else if afterint < 0 {
		sgb.After = time.Now().Add(time.Duration(afterint) * time.Second)
	}

	sgb.SFW, _ = strconv.Atoi(query.Get("sfw"))
	if sgb.SFW < 0 {
		sgb.SFW = 0
	} else if sgb.SFW > 1 {
		sgb.SFW = 1
	}
}

func (sgb *SubredditGetByNameImageParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if sgb.Q != "" {
		arg := sqlite.Arg("%" + sgb.Q + "%")
		expr = append(expr,
			sm.Where(
				models.ImageColumns.PostTitle.Like(arg).
					Or(models.ImageColumns.PostURL.Like(arg)).
					Or(models.ImageColumns.ImageRelativePath.Like(arg)).
					Or(models.ImageColumns.PostAuthor.Like(arg)),
			),
		)
	}

	if sgb.Device != "" {
		expr = append(expr, models.SelectWhere.Images.Device.EQ(sgb.Device))
	}

	if !sgb.After.IsZero() {
		expr = append(expr, models.SelectWhere.Images.CreatedAt.GTE(sgb.After.Unix()))
	}

	if sgb.SFW == 1 {
		expr = append(expr, models.SelectWhere.Images.NSFW.EQ(0))
	}

	return expr
}

func (sgb *SubredditGetByNameImageParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	expr = append(expr, sgb.CountQuery()...)

	if sgb.Limit > 0 {
		expr = append(expr, sm.Limit(sgb.Limit))
	}

	if sgb.Offset > 0 {
		expr = append(expr, sm.Offset(sgb.Offset))
	}

	if sgb.OrderBy != "" {
		order := sm.OrderBy(sqlite.Quote(sgb.OrderBy))
		if sgb.Sort == "desc" {
			expr = append(expr, order.Desc())
		} else {
			expr = append(expr, order.Asc())
		}
	} else {
		expr = append(expr, sm.OrderBy(models.ImageColumns.CreatedAt).Desc())
	}

	return expr
}

type SubredditGetByNameImageResult struct {
	Subreddit *models.Subreddit
	Images    models.ImageSlice
	Total     int64
}

func (api *API) SubredditGetByNameWithImages(ctx context.Context, name string, imageParams SubredditGetByNameImageParams) (result SubredditGetByNameImageResult, err error) {
	ctx, span := tracer.Start(ctx, "*API.SubredditsGetByNameWithImages")
	defer span.End()

	result.Subreddit, err = api.SubredditsGetByName(ctx, name)
	if err != nil {
		return result, err
	}

	result.Images, err = models.Images.
		Query(ctx, api.db, append(imageParams.Query(), models.SelectWhere.Images.Subreddit.EQ(result.Subreddit.Name))...).
		All()
	if err != nil {
		return result, errs.Wrapw(err, "failed to get images by subreddit", "subreddit", result.Subreddit.Name, "params", imageParams)
	}

	if err := result.Images.LoadImageDevice(ctx, api.db); err != nil {
		return result, errs.Wrapw(err, "failed to get device by images")
	}

	result.Total, err = models.Images.
		Query(ctx, api.db, append(imageParams.CountQuery(), models.SelectWhere.Images.Subreddit.EQ(result.Subreddit.Name))...).
		Count()
	if err != nil {
		return result, errs.Wrapw(err, "failed to count images by subreddit", "subreddit", result.Subreddit.Name, "params", imageParams)
	}

	return result, nil
}
