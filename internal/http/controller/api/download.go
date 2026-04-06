// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
)

const GoStoragePrefix = "https://golang.google.cn/dl/"

type DownloadController struct{}

func (self DownloadController) RegisterRoute(g *echo.Group) {
	g.GET("/downloads", self.List)
	g.HEAD("/downloads/golang/:filename", self.FetchPackage)
	g.GET("/downloads/golang/:filename", self.FetchPackage)
}

// List Go 安装包下载列表
func (DownloadController) List(ctx echo.Context) error {
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

	return success(ctx, map[string]interface{}{
		"featured":          featured,
		"stables":           stables,
		"stable_versions":   stableVersions,
		"unstables":         unstables,
		"archiveds":         archiveds,
		"archived_versions": archivedVersions,
	})
}

var filenameReg = regexp.MustCompile(`\d+\.\d[a-z\.]*\d+`)

// FetchPackage 下载 Go 安装包（重定向到实际下载地址）
func (DownloadController) FetchPackage(ctx echo.Context) error {
	filename := ctx.Param("filename")
	if filename == "" {
		return fail(ctx, "文件名不能为空")
	}

	// 异步记录下载次数，带 panic 恢复
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 静默恢复，避免 goroutine panic 导致进程崩溃
			}
		}()
		logic.DefaultDownload.RecordDLTimes(context.EchoContext(ctx), filename)
	}()

	// 优先从官方 CDN 下载
	officialUrl := GoStoragePrefix + filename
	if redirectUrl, ok := checkDownloadUrl(officialUrl); ok {
		return ctx.Redirect(http.StatusSeeOther, redirectUrl)
	}

	// 从配置的镜像下载
	goVersion := filenameReg.FindString(filename)
	if goVersion == "" {
		return fail(ctx, "无法解析 Go 版本号")
	}
	filePath := fmt.Sprintf("go/%s/%s", goVersion, filename)

	dlUrls := strings.Split(config.ConfigFile.MustValue("download", "dl_urls"), ",")
	for _, dlUrl := range dlUrls {
		fullUrl := dlUrl + filePath
		if redirectUrl, ok := checkDownloadUrl(fullUrl); ok {
			return ctx.Redirect(http.StatusSeeOther, redirectUrl)
		}
	}

	// 兜底：从站内静态目录下载
	return ctx.Redirect(http.StatusSeeOther, "/static/"+filePath)
}

// httpClient 复用连接池，避免每次请求创建新 Client
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// checkDownloadUrl 检查下载 URL 是否可用
// 正确关闭 Response.Body 防止连接泄漏
func checkDownloadUrl(dlUrl string) (string, bool) {
	resp, err := httpClient.Head(dlUrl)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", false
	}
	if resp.StatusCode == http.StatusOK {
		return dlUrl, true
	}
	return "", false
}
