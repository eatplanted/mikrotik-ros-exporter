package server

import (
	"github.com/eatplanted/mikrotik-ros-exporter/internal/metrics"
	"github.com/eatplanted/mikrotik-ros-exporter/internal/mikrotik"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (s *server) probeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("target")
		credentialName := r.URL.Query().Get("credential")

		credential, err := s.config.FindCredential(credentialName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithFields(log.Fields{
				"target":     target,
				"credential": credentialName,
			}).WithError(err).Error("failed to find credential")
			return
		}

		client := mikrotik.NewClient(mikrotik.Configuration{
			Address:  target,
			Username: credential.Username,
			Password: credential.Password,
		})

		registry, err := metrics.CreateRegistryWithMetrics(client)
		if err != nil {
			log.WithFields(log.Fields{
				"target":     target,
				"credential": credentialName,
			}).WithError(err).Error("failed to get health")

			// We don't return here because the registry is still valid
			// and contains a mikrotik_probe_success metric which is
			// set to 0.
		}

		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		handler.ServeHTTP(w, r)
	}
}
