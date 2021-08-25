//go:build !windows && !plan9
// +build !windows,!plan9

package main

import (
	"log"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
)

func gracefulRun(server *http.Server) {
	log.Fatal(gracehttp.Serve(server))
}
