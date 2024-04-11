package config

var DefaultConfig = map[string]any{
	"flags.containerized": false,

	"log.enable":      true,
	"log.source":      true,
	"log.format":      "pretty",
	"log.level":       "info",
	"log.output":      "stderr",
	"log.file.enable": true,
	"log.file.path":   "redmage.log",

	"db.driver":      "sqlite3",
	"db.string":      "data.db",
	"db.automigrate": true,

	"download.concurrency":       5,
	"download.directory":         "",
	"download.timeout.headers":   "10s",
	"download.timeout.idle":      "5s",
	"download.timeout.idlespeed": "10KB",

	"http.port":             "8080",
	"http.host":             "0.0.0.0",
	"http.shutdown_timeout": "5s",
	"http.hotreload":        false,
}
