package subreddits

import (
	subredditsDB "github.com/tigorlazuardi/redmage/db/queries/subreddits"
)

type API struct {
	SubredditDB *subredditsDB.Queries
}
