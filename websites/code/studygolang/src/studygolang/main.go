package main

import (
	"log"
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	router := initRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}