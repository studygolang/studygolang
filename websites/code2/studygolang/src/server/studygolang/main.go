// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package main

import (
	"http/controller"
	"http/controller/admin"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/polaris1119/config"

	pwm "http/middleware"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/fatih/structs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"github.com/polaris1119/logger"
	thirdmw "github.com/polaris1119/middleware"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())

	structs.DefaultTagName = "json"
}

func main() {
	savePid()

	logger.Init(ROOT+"/log", ConfigFile.MustValue("global", "log_level", "DEBUG"))

	go ServeBackGround()

	e := echo.New()

	e.Use(thirdmw.EchoLogger())
	e.Use(mw.Recover())
	e.Use(pwm.HTTPError())
	e.Use(pwm.AutoLogin())

	e.Static("/static/", ROOT+"/static")
	e.File("/favicon.ico", ROOT+"/static/img/go.ico")
	// 服务 sitemap 文件
	e.Static("/sitemap/", ROOT+"/sitemap")

	frontG := e.Group("", thirdmw.EchoCache())
	controller.RegisterRoutes(frontG)

	frontG.Get("/admin", echo.HandlerFunc(admin.AdminIndex), pwm.NeedLogin(), pwm.AdminAuth())
	adminG := e.Group("/admin", pwm.NeedLogin(), pwm.AdminAuth())
	admin.RegisterRoutes(adminG)

	addr := ConfigFile.MustValue("listen", "host", "") + ":" + ConfigFile.MustValue("listen.http", "port", "8080")
	std := standard.New(addr)
	std.SetHandler(e)

	log.Fatal(gracehttp.Serve(std.Server))
}

const (
	IfNoneMatch = "IF-NONE-MATCH"
	Etag        = "Etag"
)

func savePid() {
	pidFilename := ROOT + "/pid/" + filepath.Base(os.Args[0]) + ".pid"
	pid := os.Getpid()

	ioutil.WriteFile(pidFilename, []byte(strconv.Itoa(pid)), 0755)
}
