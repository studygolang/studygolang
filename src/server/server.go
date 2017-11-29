// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

// 可选择是否在启动主程序时，同时嵌入 indexer 和 crawler，减少内存占用

package server

import (
	"flag"
	"fmt"
	"os"
	"time"

	"logic"

	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
	"github.com/robfig/cron"
)

var usageStr = `
Usage: migrator [options]

Opthions:
	--changeVersion <version>    changeset version(1.0)
`

var (
	manualIndex   = flag.Bool("manual", false, "do manual index once or not")
	needAll       = flag.Bool("all", false, "是否需要全量抓取，默认否")
	whichSite     = flag.String("site", "", "抓取哪个站点（空表示所有站点）")
	changeVersion = flag.String("changeVersion", "", usageStr)
)

func IndexingServer() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if *manualIndex {
		indexing(true)
	}

	c := cron.New()
	// 构建 solr 需要的索引数据
	// 1 分钟一次增量
	c.AddFunc("@every 1m", func() {
		indexing(false)
	})
	// 一天一次全量
	c.AddFunc("@daily", func() {
		indexing(true)
	})

	c.Start()
}

func indexing(isAll bool) {
	logger.Infoln("indexing start...")

	start := time.Now()
	defer func() {
		logger.Infoln("indexing spend time:", time.Now().Sub(start))
	}()

	logic.DefaultSearcher.Indexing(isAll)
}

func CrawlServer() {
	if !flag.Parsed() {
		flag.Parse()
	}

	go autocrawl(*needAll, *whichSite)
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

func MigratorServer() {
	if !flag.Parsed() {
		flag.Parse()
	}
	if *changeVersion == "" {
		fmt.Printf("%s\n", usageStr)
		os.Exit(1)
	}
	logic.DefaultMigrator.Migrator(*changeVersion)
}
