// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	"math/rand"
	"runtime"
	"time"
	"flag"
	//"path/filepath"

	"github.com/robfig/cron"
	"logger"
	//"process"
	"service"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	var manualIndex bool
	flag.BoolVar(&manualIndex, "manual", false, "do manual index once or not")
	flag.Parse()

	if manualIndex {
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

	service.Indexing(isAll)
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
