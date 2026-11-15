// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	stdcontext "context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/http/middleware"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
)

const GoStoragePrefix = "https://golang.google.cn/dl/"

type DownloadController struct{}

// 注册路由
func (self DownloadController) RegisterRoute(g *echo.Group) {
	g.GET("/dl", self.GoDl)
	g.Match([]string{"GET", "HEAD"}, "/dl/golang/:filename", self.FetchGoInstallPackage)
	g.GET("/dl/add_new_version", self.AddNewDownload, middleware.NeedLogin())
}

// GoDl Go 语言安装包下载
func (DownloadController) GoDl(ctx echo.Context) error {
	downloads := logic.DefaultDownload.FindAll(context.EchoContext(ctx))

	featured := make([]*model.Download, 0, 5)
	stables := make(map[string][]*model.Download)
	stableVersions := make([]string, 0, 2)
	unstables := make(map[string][]*model.Download)
	archiveds := make(map[string][]*model.Download)
	archivedVersions := make([]string, 0, 20)

	for _, download := range downloads {
		version := download.Version
		if download.Category == model.DLStable || download.Category == model.DLFeatured {
			if _, ok := stables[version]; !ok {
				stableVersions = append(stableVersions, version)
				stables[version] = make([]*model.Download, 0, 15)
			}
			stables[version] = append(stables[version], download)

			if download.IsRecommend && len(featured) < 5 {
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
var safeFilenameReg = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func (self DownloadController) FetchGoInstallPackage(ctx echo.Context) error {
	filename := ctx.Param("filename")
	if filename == "" {
		return fail(ctx, 1, "文件名不能为空")
	}
	if !safeFilenameReg.MatchString(filename) || strings.Contains(filename, "..") {
		return fail(ctx, 1, "非法文件名")
	}

	go logic.DefaultDownload.RecordDLTimes(stdcontext.Background(), filename)

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
	if goVersion == "" {
		return fail(ctx, 1, "无法解析 Go 版本号")
	}
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

func (DownloadController) AddNewDownload(ctx echo.Context) error {
	me, ok := ctx.Get("user").(*model.Me)
	if !ok || !me.IsRoot {
		return fail(ctx, 1, "无权操作")
	}

	version := ctx.QueryParam("version")
	selector := ctx.QueryParam("selector")

	if version == "" {
		return fail(ctx, 1, "version is empty")
	}

	if selector == "" {
		selector = ".toggleVisible"
	}

	err := logic.DefaultDownload.AddNewDownload(context.EchoContext(ctx), version, selector)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}

var downloadHTTPClient = &http.Client{Timeout: 5 * time.Second}

func (DownloadController) headWithTimeout(dlUrl string) (*http.Response, error) {
	return downloadHTTPClient.Head(dlUrl)
}
