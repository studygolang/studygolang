package main

import (
	"log"
	"time"

	"github.com/labstack/echo/engine/standard"
	"github.com/tylerb/graceful"
)

func gracefulRun(std *standard.Server) {
	log.Fatal(graceful.ListenAndServe(std.Server, 5*time.Second))
}
