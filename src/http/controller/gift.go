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
	"github.com/polaris1119/goutils"
)

type GiftController struct{}

// 注册路由
func (self GiftController) RegisterRoute(g *echo.Group) {
	g.Get("/gift", self.GiftList)
	g.Post("/gift/exchange", self.Exchange, middleware.NeedLogin())
	g.Get("/gift/mine", self.MyGift, middleware.NeedLogin())
}

func (GiftController) GiftList(ctx echo.Context) error {
	gifts := logic.DefaultGift.FindAllOnline(ctx)

	if len(gifts) > 0 {
		user, ok := ctx.Get("user").(*model.Me)
		if ok {
			logic.DefaultGift.UserCanExchange(ctx, user, gifts)
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
	err := logic.DefaultGift.Exchange(ctx, me, giftId)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}

func (GiftController) MyGift(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)

	exchangeRecords := logic.DefaultGift.FindExchangeRecords(ctx, me)

	data := map[string]interface{}{
		"records": exchangeRecords,
	}

	return render(ctx, "gift/mine.html", data)
}
