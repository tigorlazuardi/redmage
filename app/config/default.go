package config

import (
	"os"
	"runtime"
	"time"

	"github.com/apparentlymart/go-userdirs/userdirs"
	"github.com/inhies/go-bytesize"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

func Default() *Config {
	downloadDir := os.Getenv("REDMAGE_DEFAULT_DOWNLOAD_DIRECTORY")
	if downloadDir == "" {
		dirs := userdirs.ForApp("Redmage", "Tigor", "id.web.tigor.redmage")
		downloadDir = dirs.DataHome()
	}
	k := koanf.New(".")
	c := &Config{
		Profiles:   make(map[string]Profile),
		Subreddits: make(map[string]SubredditConfig),
		Download: Download{
			Directory:             downloadDir,
			Concurrency:           runtime.NumCPU(),
			ConnectionTimeout:     30 * time.Second,
			DownloadIdleTimeout:   time.Minute,
			DownloadIdleThreshold: 5 * bytesize.KB,
		},
		Koanf: k,
	}

	if err := k.Load(structs.Provider(c, "koanf"), nil); err != nil {
		panic(err)
	}

	return c
}
