// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/times"
)

type TopController struct{}

// 注册路由
func (self TopController) RegisterRoute(g *echo.Group) {
	g.GET("/top/dau", self.TopDAU)
	g.GET("/top/rich", self.TopRich)
}

func (TopController) TopDAU(ctx echo.Context) error {
	data := map[string]interface{}{
		"today": times.Format("Ymd"),
	}

	data["users"] = logic.DefaultRank.FindDAURank(context.EchoContext(ctx), 10)
	data["active_num"] = logic.DefaultRank.TotalDAUUser(context.EchoContext(ctx))

	return render(ctx, "top/dau.html", data)
}

func (TopController) TopRich(ctx echo.Context) error {
	data := map[string]interface{}{
		"users": logic.DefaultRank.FindRichRank(context.EchoContext(ctx)),
	}

	return render(ctx, "top/rich.html", data)
}
