package config

var defaultConfig = map[string]any{
	"log.enable": true,
	"log.source": true,
	"log.format": "pretty",
	"log.level":  "info",
	"log.output": "stderr",

	"db.driver":      "sqlite3",
	"db.string":      "data.db",
	"db.automigrate": true,
}
