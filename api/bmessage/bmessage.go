package bmessage

import (
	"github.com/alecthomas/units"
)

type ImageMetadata struct {
	URL             string
	Height          int
	Width           int
	ThumbnailURL    string
	ThumbnailHeight int
	ThumbnailWidth  int
}

type DownloadStatusMessage struct {
	Metadata      ImageMetadata
	ContantLength units.MetricBytes
	Downloaded    units.MetricBytes
	Subreddit     string
	PostURL       string
	PostID        string
	Error         error
}
