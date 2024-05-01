package subredditsview

import (
	"github.com/tigorlazuardi/redmage/api"
)

type Data struct {
	Subreddits api.ListSubredditsResult
	Error      string
}
