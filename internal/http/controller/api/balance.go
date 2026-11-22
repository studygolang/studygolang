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

type BalanceController struct{}

func (self BalanceController) RegisterRoute(g *echo.Group) {
	g.GET("/balance", self.MyBalance)
	g.GET("/balance/add", self.Add)
}

// MyBalance 当前登录用户的积分余额明细（需要登录）
// 响应字段对齐前端 balanceAPI.getMyBalance 契约：
//
//	{ balance, incomes[], expenses[], total, page, has_more }
//
// incomes  = Num > 0 的明细（充值/奖励等）
// expenses = Num < 0 的明细（消费/处罚等）
// details  = 全部明细（保留以便老客户端兼容）
func (BalanceController) MyBalance(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	p := goutils.MustInt(ctx.QueryParam("p"), 1)

	details := logic.DefaultUserRich.FindBalanceDetail(context.EchoContext(ctx), me, p)
	total := logic.DefaultUserRich.Total(context.EchoContext(ctx), me.Uid)

	// 计算是否有更多（FindBalanceDetail 用 logic.CommentPerNum 作为每页大小）
	hasMore := int64(p*logic.CommentPerNum) < total

	if details == nil {
		details = make([]*model.UserBalanceDetail, 0)
	}

	// 按收入/支出拆分（Num 正负）
	incomes := make([]*model.UserBalanceDetail, 0, len(details))
	expenses := make([]*model.UserBalanceDetail, 0, len(details))
	for _, d := range details {
		if d.Num >= 0 {
			incomes = append(incomes, d)
		} else {
			expenses = append(expenses, d)
		}
	}

	// 没有任何积分记录 = 还没领过初始资本。旧版 balance.html 用「明细为空」
	// 渲染领取卡片，前后端分离后把这个判断所需的数据下发给前端。
	var initialMission *model.Mission
	if total == 0 {
		initialMission = logic.DefaultMission.InitialMission(context.EchoContext(ctx))
	}

	return success(ctx, map[string]interface{}{
		"balance":  me.Balance,
		"incomes":  incomes,
		"expenses": expenses,
		// 兼容字段（旧客户端或调试用）
		"details":         details,
		"total":           total,
		"page":            p,
		"has_more":        hasMore,
		"initial_mission": initialMission,
	})
}

// Add 充值记录（需要登录）
func (BalanceController) Add(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	p := goutils.MustInt(ctx.QueryParam("p"), 1)

	details := logic.DefaultUserRich.FindBalanceDetail(context.EchoContext(ctx), me, p, model.MissionTypeAdd)
	rechargeAmount := logic.DefaultUserRich.FindRecharge(context.EchoContext(ctx), me)

	if details == nil {
		details = make([]*model.UserBalanceDetail, 0)
	}

	return success(ctx, map[string]interface{}{
		"balance":         me.Balance,
		"details":         details,
		"recharge_amount": rechargeAmount,
		"page":            p,
	})
}
