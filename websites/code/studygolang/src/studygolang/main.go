package main

import (
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	router := initRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
