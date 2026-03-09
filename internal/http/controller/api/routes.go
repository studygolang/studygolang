// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	echo "github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes 注册面向 Next.js 前端的 Web API 路由（/api/v1/...）
func RegisterRoutes(g *echo.Group) {
	// 允许跨域（开发环境 Next.js dev server 需要）
	g.Use(mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type", "Authorization", "X-Token"},
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
}
