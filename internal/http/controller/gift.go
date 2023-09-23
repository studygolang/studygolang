// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/http/middleware"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type GiftController struct{}

// 注册路由
func (self GiftController) RegisterRoute(g *echo.Group) {
	g.GET("/gift", self.GiftList)
	g.POST("/gift/exchange", self.Exchange, middleware.NeedLogin())
	g.GET("/gift/mine", self.MyGift, middleware.NeedLogin())
}

func (GiftController) GiftList(ctx echo.Context) error {
	gifts := logic.DefaultGift.FindAllOnline(context.EchoContext(ctx))

	if len(gifts) > 0 {
		user, ok := ctx.Get("user").(*model.Me)
		if ok {
			logic.DefaultGift.UserCanExchange(context.EchoContext(ctx), user, gifts)
		}
	}

	data := map[string]interface{}{
		"gifts": gifts,
	}

	return render(ctx, "gift/list.html", data)
}

func (GiftController) Exchange(ctx echo.Context) error {
	giftId := goutils.MustInt(ctx.FormValue("gift_id"))
	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultGift.Exchange(context.EchoContext(ctx), me, giftId)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}

func (GiftController) MyGift(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)

	exchangeRecords := logic.DefaultGift.FindExchangeRecords(context.EchoContext(ctx), me)

	data := map[string]interface{}{
		"records": exchangeRecords,
	}

	return render(ctx, "gift/mine.html", data)
}
