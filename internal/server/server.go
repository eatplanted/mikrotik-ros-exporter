package server

import (
	"github.com/eatplanted/mikrotik-ros-exporter/internal/config"
	"github.com/gorilla/mux"
	"net/http"
)

type Server interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type server struct {
	config config.Configuration
	router *mux.Router
}

func NewServer(config config.Configuration) Server {
	server := &server{
		config: config,
		router: mux.NewRouter(),
	}

	server.routes()

	return server
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
