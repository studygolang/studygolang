// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"flag"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"), "crawl")

	var (
		needAll           bool
		crawlConfFilename string
		whichSite         string
	)
	flag.BoolVar(&needAll, "all", false, "是否需要全量抓取，默认否")
	flag.StringVar(&crawlConfFilename, "config", "config/auto_crawl.json", "自动抓取配置文件")
	flag.StringVar(&whichSite, "site", "", "抓取配置中哪个站点（空表示所有配置站点）")
	flag.Parse()

	if !filepath.IsAbs(crawlConfFilename) {
		crawlConfFilename = config.ROOT + "/" + crawlConfFilename
	}

	go autocrawl(needAll, crawlConfFilename, whichSite)

	select {}
}
