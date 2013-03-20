// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	. "config"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"process"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	SavePid()
	// 服务静态文件
	http.Handle("/static/", http.FileServer(http.Dir(ROOT)))

	router := initRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(Config["host"], nil))
}

// 保存PID
func SavePid() {
	pidFile := Config["pid"]
	if !filepath.IsAbs(Config["pid"]) {
		pidFile = ROOT + "/" + pidFile
	}
	// TODO：错误不处理
	process.SavePidTo(pidFile)
}
