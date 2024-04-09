package routes

import (
	"net/http"
)

func (routes *Routes) HealthCheck(rw http.ResponseWriter, req *http.Request) {
	_, _ = rw.Write([]byte{})
}
