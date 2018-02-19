// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"fmt"
	"logic"
	"model"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/polaris1119/config"
)

const GoStoragePrefix = "https://dl.google.com/go/"

type DownloadController struct{}

// 注册路由
func (self DownloadController) RegisterRoute(g *echo.Group) {
	g.Get("/dl", self.GoDl)
	g.Get("/dl/golang/:filename", self.FetchGoInstallPackage)
}

// GoDl Go 语言安装包下载
func (DownloadController) GoDl(ctx echo.Context) error {
	downloads := logic.DefaultDownload.FindAll(ctx)

	featured := make([]*model.Download, 0, 4)
	stables := make(map[string][]*model.Download)
	stableVersions := make([]string, 0, 2)
	unstables := make(map[string][]*model.Download)
	archiveds := make(map[string][]*model.Download)
	archivedVersions := make([]string, 0, 20)

	for _, download := range downloads {
		version := download.Version
		if download.Category == model.DLStable {
			if _, ok := stables[version]; !ok {
				stableVersions = append(stableVersions, version)
				stables[version] = make([]*model.Download, 0, 15)
			}
			stables[version] = append(stables[version], download)

			if download.IsRecommend && len(featured) < 4 {
				featured = append(featured, download)
			}
		} else if download.Category == model.DLUnstable {
			if _, ok := unstables[version]; !ok {
				unstables[version] = make([]*model.Download, 0, 15)
			}
			unstables[version] = append(unstables[version], download)
		} else if download.Category == model.DLArchived {
			if _, ok := archiveds[version]; !ok {
				archivedVersions = append(archivedVersions, version)
				archiveds[version] = make([]*model.Download, 0, 15)
			}
			archiveds[version] = append(archiveds[version], download)
		}
	}

	data := map[string]interface{}{
		"activeDl":          "active",
		"featured":          featured,
		"stables":           stables,
		"stable_versions":   stableVersions,
		"unstables":         unstables,
		"archiveds":         archiveds,
		"archived_versions": archivedVersions,
	}

	return render(ctx, "download/go.html", data)
}

var filenameReg = regexp.MustCompile(`\d+\.\d[a-z\.]*\d+`)

func (self DownloadController) FetchGoInstallPackage(ctx echo.Context) error {
	filename := ctx.Param("filename")

	officalUrl := GoStoragePrefix + filename
	resp, err := self.headWithTimeout(officalUrl)
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		return ctx.Redirect(http.StatusSeeOther, officalUrl)
	}
	if err == nil {
		resp.Body.Close()
	}

	goVersion := filenameReg.FindString(filename)
	filePath := fmt.Sprintf("go/%s/%s", goVersion, filename)

	dlUrls := strings.Split(config.ConfigFile.MustValue("download", "dl_urls"), ",")
	for _, dlUrl := range dlUrls {
		dlUrl += filePath
		resp, err = self.headWithTimeout(dlUrl)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return ctx.Redirect(http.StatusSeeOther, dlUrl)
		}
		if err == nil {
			resp.Body.Close()
		}
	}

	getLogger(ctx).Infoln("download:", filename, "from the site static directory")

	return ctx.Redirect(http.StatusSeeOther, "/static/"+filePath)
}

func (DownloadController) headWithTimeout(dlUrl string) (*http.Response, error) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	return client.Head(dlUrl)
}
