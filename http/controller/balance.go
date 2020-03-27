// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"github.com/polaris1119/goutils"
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/http/middleware"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"

	echo "github.com/labstack/echo/v4"
)

type UserRichController struct{}

// 注册路由
func (self UserRichController) RegisterRoute(g *echo.Group) {
	g.GET("/balance", self.MyBalance, middleware.NeedLogin())
	g.GET("/balance/add", self.Add, middleware.NeedLogin())
}

func (UserRichController) MyBalance(ctx echo.Context) error {
	p := goutils.MustInt(ctx.QueryParam("p"), 1)
	me := ctx.Get("user").(*model.Me)
	balanceDetails := logic.DefaultUserRich.FindBalanceDetail(context.EchoContext(ctx), me, p)
	total := logic.DefaultUserRich.Total(context.EchoContext(ctx), me.Uid)

	data := map[string]interface{}{
		"details": balanceDetails,
		"total":   int(total),
		"cur_p":   p,
	}
	return render(ctx, "rich/balance.html", data)
}

func (UserRichController) Add(ctx echo.Context) error {
	p := goutils.MustInt(ctx.QueryParam("p"), 1)
	me := ctx.Get("user").(*model.Me)
	balanceDetails := logic.DefaultUserRich.FindBalanceDetail(context.EchoContext(ctx), me, p, model.MissionTypeAdd)

	rechargeAmount := logic.DefaultUserRich.FindRecharge(context.EchoContext(ctx), me)

	data := map[string]interface{}{
		"details":         balanceDetails,
		"recharge_amount": rechargeAmount,
	}
	return render(ctx, "rich/add.html", data)
}
