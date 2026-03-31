// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"os"
	"strings"

	echo "github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/polaris1119/config"
)

// getAllowedOrigins 从环境变量 ALLOWED_ORIGINS 读取允许的跨域来源，
// 格式为逗号分隔的 URL 列表。
// 生产环境必须配置，开发环境默认为 http://localhost:3000
func getAllowedOrigins() []string {
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw == "" {
		// 检查环境
		env := config.ConfigFile.MustValue("global", "env", "dev")
		if env == "prod" {
			panic("生产环境必须配置 ALLOWED_ORIGINS 环境变量！示例：ALLOWED_ORIGINS=https://studygolang.com,https://www.studygolang.com")
		}
		// 开发环境使用默认值
		return []string{"http://localhost:3000"}
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if o := strings.TrimSpace(p); o != "" {
			result = append(result, o)
		}
	}

	// 生产环境校验：至少配置一个来源
	if len(result) == 0 {
		env := config.ConfigFile.MustValue("global", "env", "dev")
		if env == "prod" {
			panic("ALLOWED_ORIGINS 配置为空，生产环境必须至少指定一个允许的来源！")
		}
		return []string{"http://localhost:3000"}
	}

	return result
}

// RegisterRoutes 注册面向 Next.js 前端的 Web API 路由（/api/v1/...）
func RegisterRoutes(g *echo.Group) {
	// CORS：允许来源通过 ALLOWED_ORIGINS 环境变量配置
	// 开发默认 http://localhost:3000；生产设置为实际域名（如 https://studygolang.com）
	// AllowCredentials=true 以支持 HttpOnly Cookie 认证
	g.Use(mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Token"},
		AllowCredentials: true,
	}))
	// CSRF 防护：写操作检查 Origin/Referer
	g.Use(originCheck)

	new(IndexController).RegisterRoute(g)
	new(ArticleController).RegisterRoute(g)
	new(TopicController).RegisterRoute(g)
	new(ProjectController).RegisterRoute(g)
	new(ResourceController).RegisterRoute(g)
	new(ReadingController).RegisterRoute(g)
	new(BookController).RegisterRoute(g)
	new(UserController).RegisterRoute(g)
	new(CommentController).RegisterRoute(g)
	new(SearchController).RegisterRoute(g)
	new(SidebarController).RegisterRoute(g)
	new(WikiController).RegisterRoute(g)
	new(InterviewController).RegisterRoute(g)
	new(JobController).RegisterRoute(g)
	new(MessageController).RegisterRoute(g)
	new(FavoriteController).RegisterRoute(g)
	new(LikeController).RegisterRoute(g)
}
