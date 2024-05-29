package api

import (
	"context"
	"encoding/json"
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

type ImageListParams struct {
	Q         string
	NSFW      int32
	OrderBy   string
	Sort      string
	Offset    int64
	Limit     int64
	Device    string
	Subreddit string
	CreatedAt time.Time
}

func (ilp *ImageListParams) FillFromQuery(query Queryable) {
	ilp.Q = query.Get("q")
	if nsfw, err := strconv.Atoi(query.Get("nsfw")); err == nil {
		ilp.NSFW = int32(nsfw)
	} else {
		ilp.NSFW = -1
	}
	ilp.OrderBy = query.Get("order_by")
	ilp.Sort = strings.ToLower(query.Get("sort"))
	ilp.Offset, _ = strconv.ParseInt(query.Get("offset"), 10, 64)
	limit := query.Get("limit")
	var err error
	ilp.Limit, err = strconv.ParseInt(limit, 10, 64)
	if err != nil || ilp.Limit < 0 {
		ilp.Limit = 25
	}
	ilp.Device = query.Get("device")
	ilp.Subreddit = query.Get("subreddit")

	createdAtint, _ := strconv.ParseInt(query.Get("created_at"), 10, 64)
	if createdAtint > 0 {
		ilp.CreatedAt = time.Unix(createdAtint, 0)
	} else if createdAtint < 0 {
		// Negative value means relative to now.
		ilp.CreatedAt = time.Now().Add(time.Duration(createdAtint) * time.Second)
	}
}

func (ilp ImageListParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if ilp.Q != "" {
		arg := sqlite.Arg("%" + ilp.Q + "%")
		expr = append(expr,
			sm.Where(
				models.ImageColumns.PostTitle.Like(arg).
					Or(models.ImageColumns.PostAuthor.Like(arg)).
					Or(models.ImageColumns.ImageRelativePath.Like(arg)),
			),
		)
	}

	if ilp.NSFW >= 0 {
		expr = append(expr, models.SelectWhere.Images.NSFW.EQ(ilp.NSFW))
	}

	if len(ilp.Device) > 0 {
		expr = append(expr, models.SelectWhere.Images.Device.EQ(ilp.Device))
	}

	if len(ilp.Subreddit) > 0 {
		expr = append(expr, models.SelectWhere.Images.Subreddit.EQ(ilp.Subreddit))
	}

	if !ilp.CreatedAt.IsZero() {
		expr = append(expr, models.SelectWhere.Images.CreatedAt.GTE(ilp.CreatedAt.Unix()))
	}

	return expr
}

func (ilp ImageListParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	expr = append(expr, ilp.CountQuery()...)
	if ilp.Limit > 0 {
		expr = append(expr, sm.Limit(ilp.Limit))
	} else if ilp.Offset > 0 {
		expr = append(expr, sm.Limit(-1))
	}

	if ilp.Offset > 0 {
		expr = append(expr, sm.Offset(ilp.Offset))
	}

	if ilp.OrderBy != "" {
		order := sm.OrderBy(sqlite.Quote(ilp.OrderBy))
		if ilp.Sort == "desc" {
			expr = append(expr, order.Desc())
		} else {
			expr = append(expr, order.Asc())
		}
	} else {
		expr = append(expr, sm.OrderBy(models.ImageColumns.CreatedAt).Desc())
	}
	return expr
}

type ImageListResult struct {
	Total  int64
	Images models.ImageSlice
}

func (im ImageListResult) MarshalJSON() ([]byte, error) {
	type I struct {
		*models.Image
		Device    *models.Device    `json:"device,omitempty"`
		Subreddit *models.Subreddit `json:"subreddit,omitempty"`
	}
	type A struct {
		Total  int64 `json:"total"`
		Images []I   `json:"images"`
	}

	a := A{Total: im.Total}
	a.Images = make([]I, len(im.Images))
	for i := 0; i < len(a.Images); i++ {
		a.Images[i].Image = im.Images[i]
		a.Images[i].Device = im.Images[i].R.Device
		a.Images[i].Subreddit = im.Images[i].R.Subreddit
	}

	return json.Marshal(a)
}

func (api *API) ImagesList(ctx context.Context, params ImageListParams) (result ImageListResult, err error) {
	ctx, span := tracer.Start(ctx, "*API.ImagesList")
	defer span.End()

	result.Images, err = models.Images.Query(ctx, api.db, params.Query()...).All()
	if err != nil {
		return result, errs.Wrapw(err, "failed to query for images", "params", params)
	}

	result.Total, err = models.Images.Query(ctx, api.db, params.CountQuery()...).Count()
	if err != nil {
		return result, errs.Wrapw(err, "failed to query for images", "params", params)
	}

	return result, err
}

func (api *API) ImagesListWithDevicesAndSubreddits(ctx context.Context, params ImageListParams) (result ImageListResult, err error) {
	ctx, span := tracer.Start(ctx, "*API.ImagesListWithDevicesAndSubreddits")
	defer span.End()

	result, err = api.ImagesList(ctx, params)
	if err != nil {
		return result, err
	}

	if err := result.Images.LoadImageDevice(ctx, api.db); err != nil {
		return result, errs.Wrapw(err, "failed to load image devices", "params", params)
	}

	if err := result.Images.LoadImageSubreddit(ctx, api.db); err != nil {
		return result, errs.Wrapw(err, "failed to load image subreddits", "params", params)
	}

	return result, err
}
