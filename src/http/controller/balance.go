// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"http/middleware"
	"logic"
	"model"

	"github.com/labstack/echo"
)

type UserRichController struct{}

// 注册路由
func (self UserRichController) RegisterRoute(g *echo.Group) {
	g.Get("/balance", self.MyBalance, middleware.NeedLogin())
}

func (UserRichController) MyBalance(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	balanceDetails := logic.DefaultUserRich.FindBalanceDetail(ctx, me)

	data := map[string]interface{}{
		"details": balanceDetails,
	}
	return render(ctx, "rich/balance.html", data)
}
