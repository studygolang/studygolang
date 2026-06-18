// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type GiftController struct{}

func (self GiftController) RegisterRoute(g *echo.Group) {
	g.GET("/gift", self.GiftList)
	g.POST("/gift/exchange", self.Exchange)
	g.GET("/gift/mine", self.MyGift)
}

// GiftList 在线礼品列表（无需登录，登录用户带兑换状态）
func (GiftController) GiftList(ctx echo.Context) error {
	gifts := logic.DefaultGift.FindAllOnline(context.EchoContext(ctx))
	if gifts == nil {
		gifts = make([]*model.Gift, 0)
	}

	// 已登录用户检查兑换状态（登录可选；不能用 requireAuth，会触发响应双写）
	if me := optionalAuth(ctx); me != nil && len(gifts) > 0 {
		logic.DefaultGift.UserCanExchange(context.EchoContext(ctx), me, gifts)
	}

	return success(ctx, map[string]interface{}{
		"gifts": gifts,
	})
}

// Exchange 兑换礼品（需登录，form: gift_id）
func (GiftController) Exchange(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	giftId := goutils.MustInt(ctx.FormValue("gift_id"))
	if giftId == 0 {
		return fail(ctx, "礼品 ID 不能为空")
	}

	err2 := logic.DefaultGift.Exchange(context.EchoContext(ctx), me, giftId)
	if err2 != nil {
		getLogger(ctx).Errorln("gift exchange failed:", err2)
		return fail(ctx, "兑换失败，请稍后重试")
	}

	return success(ctx, nil)
}

// MyGift 我的兑换记录（需登录）
func (GiftController) MyGift(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	records := logic.DefaultGift.FindExchangeRecords(context.EchoContext(ctx), me)
	if records == nil {
		records = make([]*model.UserExchangeRecord, 0)
	}

	return success(ctx, map[string]interface{}{
		"records": records,
	})
}
