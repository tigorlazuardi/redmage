package api

import (
	"net/http"
)

func HealthCheck(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	_, _ = rw.Write([]byte(`{"status":"ok"}`))
}
