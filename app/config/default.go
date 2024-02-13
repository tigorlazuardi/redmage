package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/apparentlymart/go-userdirs/userdirs"
	"github.com/inhies/go-bytesize"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

func init() {
	_ = godotenv.Load()
}

func Default() *Config {
	dirs := userdirs.ForApp("Redmage", "Tigor", "id.web.tigor.redmage")
	downloadDir := os.Getenv("REDMAGE_DOWNLOAD_DIRECTORY")
	if downloadDir == "" {
		downloadDir = dirs.DataHome()
	}
	configFile := os.Getenv("REDMAGE_CONFIG")
	if configFile == "" {
		configFile = filepath.Join(dirs.ConfigHome(), "redmage.yaml")
	}
	k := koanf.New(".")
	c := &Config{
		Server: Server{
			Host: "0.0.0.0",
			Port: 8090,
		},
		Devices:    make(map[string]Device),
		Subreddits: make(map[string]SubredditConfig),
		Download: Download{
			Directory:             downloadDir,
			Concurrency:           runtime.NumCPU(),
			ConnectionTimeout:     30 * time.Second,
			DownloadIdleTimeout:   time.Minute,
			DownloadIdleThreshold: 5 * bytesize.KB,
		},
		HotReload:  strings.HasPrefix(os.Args[0], os.TempDir()), // check if command is executed using "go run"
		ConfigFile: configFile,
		Koanf:      k,
	}

	if err := k.Load(structs.Provider(c, "koanf"), nil); err != nil {
		panic(err)
	}

	return c
}
