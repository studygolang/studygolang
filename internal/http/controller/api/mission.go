// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"strconv"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/polaris1119/times"

	echo "github.com/labstack/echo/v4"
)

type MissionController struct{}

func (self MissionController) RegisterRoute(g *echo.Group) {
	g.GET("/mission/daily", self.Daily)
	g.POST("/mission/daily/redeem", self.DailyRedeem)
	g.POST("/mission/complete/:id", self.Complete)
}

// Daily 获取每日登录任务状态（需登录）
func (MissionController) Daily(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	userLoginMission := logic.DefaultMission.FindLoginMission(context.EchoContext(ctx), me)

	hadRedeem := false
	if userLoginMission != nil {
		userLoginMission.Uid = me.Uid
		if times.Format("Ymd") == strconv.Itoa(userLoginMission.Date) {
			hadRedeem = true
		}
	}

	return success(ctx, map[string]interface{}{
		"login_mission": userLoginMission,
		"had_redeem":    hadRedeem,
	})
}

// DailyRedeem 领取每日登录奖励（需登录）
func (MissionController) DailyRedeem(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	err = logic.DefaultMission.RedeemLoginAward(context.EchoContext(ctx), me)
	if err != nil {
		return fail(ctx, err.Error())
	}

	return success(ctx, map[string]interface{}{
		"message": "领取成功",
	})
}

// Complete 完成指定任务（需登录）
func (MissionController) Complete(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := ctx.Param("id")
	if id == "" {
		return fail(ctx, "任务 ID 不能为空")
	}

	err = logic.DefaultMission.Complete(context.EchoContext(ctx), me, id)
	if err != nil {
		return fail(ctx, err.Error())
	}

	return success(ctx, map[string]interface{}{
		"message": "任务完成",
	})
}
