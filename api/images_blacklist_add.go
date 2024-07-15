package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/mattn/go-sqlite3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ImagesBlacklistAddParams struct {
	Subreddit string
	PostName  string
	Device    string
}

func (iabp *ImagesBlacklistAddParams) FillFromQuery(query Queryable) {
	iabp.Subreddit = query.Get("subreddit")
	iabp.PostName = query.Get("post_name")
	iabp.Device = query.Get("device")
}

func (iabp *ImagesBlacklistAddParams) FillFromCSV(values string) {
	v := strings.Split(values, ",")
	for len(v) < 3 {
		v = append(v, "")
	}
	iabp.Device = strings.TrimSpace(v[0])
	iabp.Subreddit = strings.TrimSpace(v[1])
	iabp.PostName = strings.TrimSpace(v[2])
}

func (iabp *ImagesBlacklistAddParams) Validate() error {
	var e []error
	if iabp.Subreddit == "" {
		e = append(e, errors.New("subreddit is required"))
	}
	if iabp.PostName == "" {
		e = append(e, errors.New("post_name is required"))
	}
	if err := errors.Join(e...); err != nil {
		return errs.Wrap(err).Code(http.StatusBadRequest).Details("params", iabp)
	}
	return nil
}

func (api *API) ImagesBlacklistAdd(ctx context.Context, params []ImagesBlacklistAddParams) (blacklists models.BlacklistSlice, err error) {
	ctx, span := tracer.Start(ctx, "*API.ImageAddToBlacklist")
	defer span.End()

	now := time.Now()

	sets := make([]*models.BlacklistSetter, 0, len(params))
	for _, param := range params {
		sets = append(sets, &models.BlacklistSetter{
			Device:    omit.From(param.Device),
			Subreddit: omit.From(param.Subreddit),
			PostName:  omit.From(param.PostName),
			CreatedAt: omit.From(now.Unix()),
		})
	}

	api.lockf(func() {
		blacklists, err = models.Blacklists.InsertMany(ctx, api.db, sets...)
	})

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return blacklists, errs.Wrapw(err, "blacklist already exists", "params", params).Code(http.StatusConflict)
			}
		}
		return blacklists, errs.Wrapw(err, "failed to insert blacklist", "params", params)
	}

	return blacklists, err
}
