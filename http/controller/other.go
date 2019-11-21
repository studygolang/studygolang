// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/polaris1119/config"

	"github.com/studygolang/studygolang/util"
)

// OtherController 有些页面只是前端，因此通过这个页面统一控制
// 只需要创建模板文件就可以访问到
type OtherController struct{}

// RegisterRoute 注册路由
func (self OtherController) RegisterRoute(g *echo.Group) {
	g.GET("/*", self.Any)
}

func (OtherController) Any(ctx echo.Context) error {
	uri := ctx.Request().RequestURI
	tplFile := uri + ".html"
	if util.Exist(path.Clean(config.TemplateDir + tplFile)) {
		return render(ctx, tplFile, nil)
	}

	return echo.NewHTTPError(http.StatusNotFound)
}
