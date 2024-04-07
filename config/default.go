package config

var DefaultConfig = map[string]any{
	"log.enable": true,
	"log.source": true,
	"log.format": "pretty",
	"log.level":  "info",
	"log.output": "stderr",

	"db.driver":      "sqlite3",
	"db.string":      "data.db",
	"db.automigrate": true,

	"http.port":             "8080",
	"http.host":             "0.0.0.0",
	"http.shutdown_timeout": "5s",
}
