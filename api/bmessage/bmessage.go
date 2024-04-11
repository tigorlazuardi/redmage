package bmessage

import (
	"github.com/alecthomas/units"
)

type ImageMetadata struct {
	Kind   ImageKind
	URL    string
	Height int
	Width  int
}

type ImageKind int

const (
	KindImage ImageKind = iota
	KindThumbnail
)

type ImageDownloadMessage struct {
	Metadata      ImageMetadata
	ContentLength units.MetricBytes
	Downloaded    units.MetricBytes
	Subreddit     string
	PostURL       string
	PostID        string
	Error         error
}
