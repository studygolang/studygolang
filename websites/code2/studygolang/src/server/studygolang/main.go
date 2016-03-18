// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package main

import (
	"http/controller"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	. "github.com/polaris1119/config"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/fatih/structs"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/polaris1119/logger"
	thirdmw "github.com/polaris1119/middleware"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())

	structs.DefaultTagName = "json"
}

func main() {
	logger.Init(ROOT+"/log", ConfigFile.MustValue("global", "log_level", "DEBUG"))

	go ServeBackGround()

	router := echo.New()

	router.Use(thirdmw.EchoLogger())
	router.Use(mw.Recover())
	// router.Use(mw.Gzip())
	router.Use(thirdmw.EchoCache())

	router.Static("/static/", ROOT+"/static")

	controller.RegisterRoutes(router)

	router.Get("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
	})

	addr := ConfigFile.MustValue("listen", "host", "") + ":" + ConfigFile.MustValue("listen.http", "port", "8080")
	server := router.Server(addr)

	// HTTP2 is currently enabled by default in echo.New(). To override TLS handshake errors
	// you will need to override the TLSConfig for the server so it does not attempt to validate
	// the connection using TLS as required by HTTP2
	server.TLSConfig = nil

	log.Fatal(gracehttp.Serve(server))
}

const (
	IfNoneMatch = "IF-NONE-MATCH"
	Etag        = "Etag"
)
