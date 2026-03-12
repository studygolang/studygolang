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
)

// getAllowedOrigins 从环境变量 ALLOWED_ORIGINS 读取允许的跨域来源，
// 格式为逗号分隔的 URL 列表，默认为开发环境 http://localhost:3000
func getAllowedOrigins() []string {
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw == "" {
		return []string{"http://localhost:3000"}
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if o := strings.TrimSpace(p); o != "" {
			result = append(result, o)
		}
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
}
