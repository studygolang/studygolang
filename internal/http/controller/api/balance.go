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
func (BalanceController) MyBalance(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
	if err != nil {
		return err
	}

	// 查询用户信息以获取总余额
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	// 构造 Me 对象供 FindBalanceDetail 使用
	me := &model.Me{Uid: uid}

	p := goutils.MustInt(ctx.QueryParam("p"), 1)

	details := logic.DefaultUserRich.FindBalanceDetail(context.EchoContext(ctx), me, p)
	total := logic.DefaultUserRich.Total(context.EchoContext(ctx), uid)

	// 计算是否有更多
	const pageSize = 20 // CommentPerNum 默认值
	hasMore := int64(p*pageSize) < total

	if details == nil {
		details = make([]*model.UserBalanceDetail, 0)
	}

	return success(ctx, map[string]interface{}{
		"details":  details,
		"total":    user.Balance,
		"page":     p,
		"has_more": hasMore,
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
		"details":         details,
		"recharge_amount": rechargeAmount,
		"page":            p,
	})
}
