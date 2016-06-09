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

	. "github.com/polaris1119/config"
	"github.com/polaris1119/logger"
	"github.com/robfig/cron"
)

var manualIndex = flag.Bool("manual", false, "do manual index once or not")

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	logger.Init(ROOT+"/log", ConfigFile.MustValue("global", "log_level", "DEBUG"))

	if *manualIndex {
		indexing(true)
	}

	c := cron.New()
	// 构建 solr 需要的索引数据
	// 一天一次全量
	c.AddFunc("@daily", func() {
		indexing(true)
	})

	c.Start()

	select {}
}

func indexing(isAll bool) {
	logger.Infoln("indexing start...")

	start := time.Now()
	defer func() {
		logger.Infoln("indexing spend time:", time.Now().Sub(start))
	}()

	logic.DefaultSearcher.Indexing(isAll)
}
