// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
	"github.com/polaris1119/config"
)

const GoStoragePrefix = "https://storage.googleapis.com/golang/"

type DownloadController struct{}

// 注册路由
func (self DownloadController) RegisterRoute(g *echo.Group) {
	g.Get("/dl", self.GoDl)
	g.Get("/dl/golang/:filename", self.FetchGoInstallPackage)
}

// GoDl Go 语言安装包下载
func (DownloadController) GoDl(ctx echo.Context) error {

	data := map[string]interface{}{
		"activeDl": "active",
	}

	return render(ctx, "download/go.html", data)
}

var filenameReg = regexp.MustCompile(`\d+\.\d[a-z\.]*\d+`)

func (DownloadController) FetchGoInstallPackage(ctx echo.Context) error {
	filename := ctx.Param("filename")

	officalUrl := GoStoragePrefix + filename
	resp, err := http.Head(officalUrl)
	if err == nil && resp.StatusCode == http.StatusOK {
		return ctx.Redirect(http.StatusSeeOther, officalUrl)
	}

	goVersion := filenameReg.FindString(filename)
	filePath := fmt.Sprintf("go/%s/%s", goVersion, filename)

	dlUrls := strings.Split(config.ConfigFile.MustValue("download", "dl_urls"), ",")
	for _, dlUrl := range dlUrls {
		dlUrl += filePath
		resp, err = http.Head(dlUrl)
		if err == nil && resp.StatusCode == http.StatusOK {
			return ctx.Redirect(http.StatusSeeOther, dlUrl)
		}
	}

	getLogger(ctx).Infoln("download:", filename, "from the site static directory")

	return ctx.Redirect(http.StatusSeeOther, "/static/"+filePath)
}
