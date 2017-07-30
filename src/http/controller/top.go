// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"

	"github.com/labstack/echo"

	"github.com/polaris1119/times"
)

type TopController struct{}

// 注册路由
func (self TopController) RegisterRoute(g *echo.Group) {
	g.Get("/top/dau", self.TopDAU)
	g.Get("/top/rich", self.TopRich)
}

func (TopController) TopDAU(ctx echo.Context) error {
	data := map[string]interface{}{
		"today": times.Format("Ymd"),
	}

	data["users"] = logic.DefaultRank.FindDAURank(ctx, 10)
	data["active_num"] = logic.DefaultRank.TotalDAUUser(ctx)

	return render(ctx, "top/dau.html", data)
}

func (TopController) TopRich(ctx echo.Context) error {
	data := map[string]interface{}{
		"users": logic.DefaultRank.FindRichRank(ctx),
	}

	return render(ctx, "top/rich.html", data)
}
