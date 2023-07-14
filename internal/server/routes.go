package server

func (s *server) routes() {
	s.router.HandleFunc("/metrics", s.metricsHandler()).Methods("GET")
	s.router.HandleFunc("/probe", s.probeHandler()).Methods("GET").Queries("target", "{target}", "credential", "{credential}")
}
