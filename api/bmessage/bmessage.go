package bmessage

import (
	"github.com/alecthomas/units"
)

type ImageMetadata struct {
	Kind   ImageKind
	URL    string
	Height int64
	Width  int64
}

type ImageKind int

const (
	KindImage ImageKind = iota
	KindThumbnail
)

type DownloadEvent int

func (do DownloadEvent) MarshalJSON() ([]byte, error) {
	return []byte(`"` + do.String() + `"`), nil
}

func (do DownloadEvent) String() string {
	switch do {
	case DownloadStart:
		return "DownloadStart"
	case DownloadProgress:
		return "DownloadProgress"
	case DownloadEnd:
		return "DownloadEnd"
	case DownloadError:
		return "DownloadError"
	default:
		return "Unknown"
	}
}

const (
	DownloadStart DownloadEvent = iota
	DownloadProgress
	DownloadEnd
	DownloadError
)

type ImageDownloadMessage struct {
	Event         DownloadEvent
	Metadata      ImageMetadata
	ContentLength units.MetricBytes
	Downloaded    units.MetricBytes
	Subreddit     string
	PostURL       string
	PostID        string
	Error         error
}
