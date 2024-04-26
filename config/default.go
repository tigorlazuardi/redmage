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

	"download.concurrency.images":     5,
	"download.concurrency.subreddits": 3,

	"download.directory":         "",
	"download.timeout.headers":   "10s",
	"download.timeout.idle":      "5s",
	"download.timeout.idlespeed": "10KB",
	"download.useragent":         "redmage",

	"download.pubsub.ack.deadline": "3h",

	"http.port":             "8080",
	"http.host":             "0.0.0.0",
	"http.shutdown_timeout": "5s",
	"http.hotreload":        false,

	"telemetry.openobserve.enable":             false,
	"telemetry.openobserve.log.enable":         true,
	"telemetry.openobserve.log.level":          "info",
	"telemetry.openobserve.log.source":         true,
	"telemetry.openobserve.log.endpoint":       "http://localhost:5080/api/default/default/_json",
	"telemetry.openobserve.log.concurrency":    4,
	"telemetry.openobserve.log.buffer.size":    2 * 1024, // 2kb
	"telemetry.openobserve.log.buffer.timeout": "2s",
	"telemetry.openobserve.log.username":       "root@example.com",
	"telemetry.openobserve.log.password":       "Complexpass#123",

	"telemetry.openobserve.trace.enable": true,
	"telemetry.openobserve.trace.url":    "http://localhost:5080/api/default/v1/traces",
	"telemetry.openobserve.trace.auth":   "Basic AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",

	"telemetry.trace.ratio": 1,

	"runtime.version":     "0.0.1",
	"runtime.environment": "development",
}
