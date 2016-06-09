// +build !windows,!plan9

package main

import (
	"log"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo/engine/standard"
)

func gracefulRun(std *standard.Server) {
	log.Fatal(gracehttp.Serve(std.Server))
}
