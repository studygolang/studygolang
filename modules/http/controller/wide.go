// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import "github.com/labstack/echo"

type WideController struct{}

// 注册路由
func (self WideController) RegisterRoute(g *echo.Group) {
	g.GET("/wide/playground", self.Playground)
}

// Playground Wide 的内嵌 iframe 的 playground
func (WideController) Playground(ctx echo.Context) error {
	return render(ctx, "wide/playground.html", nil)
}
