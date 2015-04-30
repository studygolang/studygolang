// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
	//"path/filepath"

	"config"
	"github.com/studygolang/mux"
	//"process"
	"api"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	var (
		needAll          bool
		crawConfFilename string
		whichSite        string
	)
	flag.BoolVar(&needAll, "all", false, "是否需要全量抓取，默认否")
	flag.StringVar(&crawConfFilename, "config", "conf/auto_crawl_conf.json", "自动抓取配置文件")
	flag.StringVar(&whichSite, "site", "", "抓取配置中哪个站点（空表示所有配置站点）")
	flag.Parse()

	go autocrawl(needAll, crawConfFilename, whichSite)

	router := initRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(config.Config["crawl_host"], nil))
}

// 保存PID
func SavePid() {
	/*
		pidFile := Config["pid"]
		if !filepath.IsAbs(Config["pid"]) {
			pidFile = ROOT + "/" + pidFile
		}
		// TODO：错误不处理
		process.SavePidTo(pidFile)
	*/
}

func initRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", api.AddArticleHandler)
	router.HandleFunc("/reddit", api.AddRedditResourceHandler)
	return router
}
