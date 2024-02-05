package config

import (
	"os"
	"runtime"
	"time"

	"github.com/apparentlymart/go-userdirs/userdirs"
	"github.com/inhies/go-bytesize"
)

func Default() *Config {
	downloadDir := os.Getenv("REDMAGE_DEFAULT_DOWNLOAD_DIRECTORY")
	if downloadDir == "" {
		dirs := userdirs.ForApp("Redmage", "Tigor", "id.web.tigor.redmage")
		downloadDir = dirs.DataHome()
	}
	return &Config{
		Profiles:   []Profile{},
		Subreddits: []Subreddit{},
		Download: Download{
			Directory:             downloadDir,
			Concurrency:           runtime.NumCPU(),
			ConnectionTimeout:     30 * time.Second,
			DownloadIdleTimeout:   time.Minute,
			DownloadIdleThreshold: 2 * bytesize.KB,
		},
	}
}
