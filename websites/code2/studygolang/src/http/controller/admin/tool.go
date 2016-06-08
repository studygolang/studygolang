// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"logic"

	"github.com/labstack/echo"
)

type ToolController struct{}

// 注册路由
func (self ToolController) RegisterRoute(g *echo.Group) {
	g.GET("/tool/sitemap", self.GenSitemap)
}

// GenSitemap
func (ToolController) GenSitemap(ctx echo.Context) error {
	logic.GenSitemap()
	return render(ctx, "tool/sitemap.html", nil)
}
