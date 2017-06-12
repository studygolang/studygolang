// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"flag"
	"math/rand"
	"time"

	"logic"

	"github.com/polaris1119/config"
	"github.com/polaris1119/keyword"
	"github.com/polaris1119/logger"
	"github.com/robfig/cron"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"), "crawl")
	go keyword.Extractor.Init(keyword.DefaultProps, true, config.ROOT+"/data/programming.txt,"+config.ROOT+"/data/dictionary.txt")

	var (
		needAll   bool
		whichSite string
	)
	flag.BoolVar(&needAll, "all", false, "是否需要全量抓取，默认否")
	flag.StringVar(&whichSite, "site", "", "抓取哪个站点（空表示所有站点）")
	flag.Parse()

	go autocrawl(needAll, whichSite)

	select {}
}

func autocrawl(needAll bool, whichSite string) {
	if needAll {
		if whichSite != "" {
			go logic.DefaultAutoCrawl.CrawlWebsite(whichSite, needAll)
		} else {
			go logic.DefaultAutoCrawl.DoCrawl(needAll)
		}
	}

	// 定时增量
	c := cron.New()
	c.AddFunc(config.ConfigFile.MustValue("crawl", "spec", "0 0 */1 * * ?"), func() {
		// 抓取 reddit
		go logic.DefaultReddit.Parse("")

		projectUrl := config.ConfigFile.MustValue("crawl", "project_url")
		if projectUrl != "" {
			// 抓取 project
			go logic.DefaultProject.ParseProjectList(projectUrl)
		}

		// 抓取 article
		go logic.DefaultAutoCrawl.DoCrawl(false)
	})
	c.Start()
}
