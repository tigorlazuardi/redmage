package config

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/apparentlymart/go-userdirs/userdirs"
	"github.com/inhies/go-bytesize"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

func init() {
	_ = godotenv.Load()
}

func Default() *Config {
	downloadDir := os.Getenv("REDMAGE_DOWNLOAD_DIRECTORY")
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
		HotReload: strings.HasPrefix(os.Args[0], os.TempDir()), // check if command is executed using "go run"
		Koanf:     k,
	}

	if err := k.Load(structs.Provider(c, "koanf"), nil); err != nil {
		panic(err)
	}

	_ = k.Load(env.Provider("REDMAGE_", ".", func(s string) string {
		key := strings.TrimPrefix(s, "REDMAGE_")
		key = strings.ToLower(key)
		key = strings.ReplaceAll(key, "__", ".")
		return key
	}), nil)

	_ = k.Unmarshal("", c)

	return c
}
