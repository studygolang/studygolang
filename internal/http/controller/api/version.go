// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
)

// 版本切换相关常量
const (
	// versionCookieName 是记录用户版本偏好的 Cookie 名。
	// Nginx 通过 $cookie_sg_version 读取此值决定路由到 Next.js 或 Go 模板。
	// 设为 HttpOnly=false 以便前端 JS 读取当前版本（用于渲染切换横幅）。
	versionCookieName = "sg_version"
	versionCookieTTL  = 30 * 24 * 3600 // 30 天，过期后回到灰度默认分配

	versionNew = "new" // 新版（Next.js）
	versionOld = "old" // 旧版（Go 模板）
)

type VersionController struct{}

// RegisterRoute 注册版本切换路由
func (VersionController) RegisterRoute(g *echo.Group) {
	g.GET("/version", VersionInfo)
	g.POST("/version/switch", SwitchVersion)
}

// VersionInfo 返回当前用户的版本状态与灰度配置
//
// 前端（新旧两端均可）在页面加载时调用此接口，用于：
//  1. 渲染顶部版本切换横幅（"体验新版" / "返回旧版"）
//  2. 显示当前命中的版本
//
// 返回字段：
//   - current:  用户当前版本偏好（new/old/auto）
//   - grayscale: 无偏好时的灰度比例（0-100）
//   - canary_target: 灰度命中的后端（nextjs/go）
func VersionInfo(ctx echo.Context) error {
	cookie, err := ctx.Cookie(versionCookieName)
	current := "auto"
	if err == nil && cookie.Value != "" {
		current = cookie.Value
	}

	return success(ctx, map[string]interface{}{
		"current":       current,
		"grayscale":     getGrayscalePercent(),
		"canary_target": getCanaryTarget(),
		"cookie_name":   versionCookieName,
	})
}

// SwitchVersion 设置版本偏好 Cookie
//
// 请求体（form 或 JSON）：
//
//	version: "new" 切换到 Next.js 新版，"old" 切换到 Go 模板旧版，"auto" 清除偏好回到灰度默认
//
// 设置 Cookie 后，浏览器下一次页面请求会带上 sg_version，
// Nginx 据此路由到对应后端。前端收到响应后自行执行页面刷新。
func SwitchVersion(ctx echo.Context) error {
	version := strings.TrimSpace(ctx.FormValue("version"))
	if version == "" {
		// 兼容 JSON body
		if ct := ctx.Request().Header.Get("Content-Type"); strings.Contains(ct, "application/json") {
			var body struct {
				Version string `json:"version"`
			}
			if err := ctx.Bind(&body); err == nil {
				version = strings.TrimSpace(body.Version)
			}
		}
	}

	switch version {
	case versionNew, versionOld:
		setVersionCookie(ctx, version)
		return success(ctx, map[string]interface{}{
			"version": version,
			"action":  "switched",
		})
	case "auto":
		clearVersionCookie(ctx)
		return success(ctx, map[string]interface{}{
			"version": "auto",
			"action":  "cleared",
		})
	default:
		return fail(ctx, "version 参数无效，允许值：new / old / auto")
	}
}

// setVersionCookie 写入版本偏好 Cookie
func setVersionCookie(ctx echo.Context, version string) {
	cookie := new(http.Cookie)
	cookie.Name = versionCookieName
	cookie.Value = version
	cookie.Path = "/"
	cookie.MaxAge = versionCookieTTL
	// 不设 HttpOnly，前端 JS 需读取以渲染横幅
	cookie.HttpOnly = false
	cookie.Secure = isSecure(ctx)
	// Lax：版本切换后浏览器导航跳转需要带上 Cookie
	cookie.SameSite = http.SameSiteLaxMode
	ctx.SetCookie(cookie)
}

// clearVersionCookie 清除版本偏好 Cookie，用户回到灰度默认分配
func clearVersionCookie(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = versionCookieName
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.HttpOnly = false
	cookie.Secure = isSecure(ctx)
	cookie.SameSite = http.SameSiteLaxMode
	ctx.SetCookie(cookie)
}

// getGrayscalePercent 返回当前灰度比例（0-100）
// 由环境变量 SG_GRAY_PERCENT 控制，默认 0（全量旧版）
func getGrayscalePercent() int {
	if v := os.Getenv("SG_GRAY_PERCENT"); v != "" {
		p, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			p = 0
		}
		if p < 0 {
			p = 0
		}
		if p > 100 {
			p = 100
		}
		return p
	}
	// 生产环境默认 0，开发环境默认 100（方便本地验证新版）
	env := config.ConfigFile.MustValue("global", "env", "dev")
	if env == "prod" {
		return 0
	}
	return 100
}

// getCanaryTarget 返回灰度命中的后端标识
func getCanaryTarget() string {
	if getGrayscalePercent() >= 100 {
		return "nextjs"
	}
	return "go"
}
