package www

import (
	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/db/queries/subreddits"
)

type WWW struct {
	Subreddits *subreddits.Queries
}

func (www *WWW) Register(router chi.Router) {}
