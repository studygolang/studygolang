package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tylerb/graceful"
)

func gracefulRun(server *http.Server) {
	log.Fatal(graceful.ListenAndServe(server, 5*time.Second))
}
