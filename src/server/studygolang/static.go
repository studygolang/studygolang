// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"github.com/labstack/echo"

	. "github.com/polaris1119/config"
)

type staticRootConf struct {
	root   string
	isFile bool
}

var staticFileMap = map[string]staticRootConf{
	"/static/":     {"/static", false},
	"/favicon.ico": {"/static/img/go.ico", true},
	// 服务 sitemap 文件
	"/sitemap/": {"/sitemap", false},
}

var filterPrefixs = make([]string, 0, 3)

func serveStatic(e *echo.Echo) {
	for prefix, rootConf := range staticFileMap {
		filterPrefixs = append(filterPrefixs, prefix)

		if rootConf.isFile {
			e.File(prefix, ROOT+rootConf.root)
		} else {
			e.Static(prefix, ROOT+rootConf.root)
		}
	}
}
