package api

import "github.com/tigorlazuardi/redmage/db/queries"

type DownloadStatusMessage struct {
	Data      queries.Image
	Progress  float64
	Subreddit string
}
