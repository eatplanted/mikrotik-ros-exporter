package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func (s *server) metricsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}
