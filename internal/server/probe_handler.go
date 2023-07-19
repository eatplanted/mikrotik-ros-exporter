package server

import (
	"net/http"
	"strconv"

	"github.com/eatplanted/mikrotik-ros-exporter/internal/config"
	"github.com/eatplanted/mikrotik-ros-exporter/internal/metrics"
	"github.com/eatplanted/mikrotik-ros-exporter/internal/mikrotik"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	defaultTimeout = 120
	timeoutOffset  = 0.5
)

func (s *server) probeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("target")
		credentialName := r.URL.Query().Get("credential")
		skipTLSVerify := r.URL.Query().Get("skip_tls_verify") == "true"

		credential, err := s.config.FindCredential(credentialName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithFields(log.Fields{
				"target":     target,
				"credential": credentialName,
			}).WithError(err).Error("failed to find credential")
			return
		}

		timeout, err := getTimeout(r, s.config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"target":     target,
				"credential": credentialName,
			}).WithError(err).Error("failed to get timeout")
			return
		}

		client := mikrotik.NewClient(mikrotik.Configuration{
			Timeout:       timeout,
			Address:       target,
			SkipTLSVerify: skipTLSVerify,
			Username:      credential.Username,
			Password:      credential.Password,
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

func getTimeout(r *http.Request, configuration config.Configuration) (timeout float64, err error) {
	if value := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); value != "" {
		var err error
		timeout, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err
		}
	}

	if timeout == 0 {
		timeout = defaultTimeout
	}

	adjustedTimeout := timeout - timeoutOffset

	if configuration.Timeout < adjustedTimeout && configuration.Timeout > 0 || adjustedTimeout < 0 {
		return configuration.Timeout, nil
	}

	return adjustedTimeout, nil
}
