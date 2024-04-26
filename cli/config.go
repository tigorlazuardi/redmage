package cli

import (
	"fmt"

	"github.com/adrg/xdg"
	"github.com/joho/godotenv"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

var cfg *config.Config

func initConfig() {
	_ = godotenv.Load()

	xdgJson, _ := xdg.ConfigFile("redmage/config.json")
	xdgYaml, _ := xdg.ConfigFile("redmage/config.yaml")

	cfg = config.NewConfigBuilder().
		LoadDefault().
		LoadJSONFile("/etc/redmage/config.json").
		LoadYamlFile("/etc/redmage/config.yaml").
		LoadJSONFile(xdgJson).
		LoadYamlFile(xdgYaml).
		LoadJSONFile("config.json").
		LoadYamlFile("config.yaml").
		LoadEnv().
		LoadFlags(RootCmd.PersistentFlags()).
		Build()

	fmt.Println("download.concurrency.subreddits", cfg.Get("download.concurrency.subreddits"))
	fmt.Println("download.concurrency.images", cfg.Get("download.concurrency.images"))

	handler := log.NewHandler(cfg)
	log.SetDefault(handler)
}
