// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/times"
)

type TopController struct{}

func (self TopController) RegisterRoute(g *echo.Group) {
	g.GET("/top/dau", self.TopDAU)
	g.GET("/top/rich", self.TopRich)
}

// TopDAU 日活跃用户排行榜（当日 top10）
func (TopController) TopDAU(ctx echo.Context) error {
	users := logic.DefaultRank.FindDAURank(context.EchoContext(ctx), 10)
	activeNum := logic.DefaultRank.TotalDAUUser(context.EchoContext(ctx))

	return success(ctx, map[string]interface{}{
		"users":      users,
		"active_num": activeNum,
		"today":      times.Format("Ymd"),
	})
}

// TopRich 财富排行榜
func (TopController) TopRich(ctx echo.Context) error {
	users := logic.DefaultRank.FindRichRank(context.EchoContext(ctx))

	return success(ctx, map[string]interface{}{
		"users": users,
	})
}
