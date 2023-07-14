package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/eatplanted/mikrotik-ros-exporter/internal/config"
	"github.com/eatplanted/mikrotik-ros-exporter/internal/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	portPtr := flag.Int("port", 9142, "Listening Port")
	configFilePtr := flag.String("config", "", "ConfigFile")
	flag.Parse()

	if configFilePtr == nil {
		panic("Config file is required")
	}

	configuration, err := config.NewConfiguration(*configFilePtr)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	listeningAddr := fmt.Sprintf("0.0.0.0:%d", *portPtr)

	s := server.NewServer(configuration)
	log.Fatal(http.ListenAndServe(listeningAddr, s))
}
