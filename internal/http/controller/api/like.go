// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type LikeController struct{}

func (self *LikeController) RegisterRoute(g *echo.Group) {
	g.POST("/likes/:objid", self.Toggle)
	g.GET("/likes/:objid/status", self.Status)
}

// Status 查询当前用户是否已点赞
func (LikeController) Status(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return success(ctx, map[string]interface{}{"has_like": false})
	}

	uid, _, valid := ValidateTokenAuto(token)
	if !valid || uid == 0 {
		return success(ctx, map[string]interface{}{"has_like": false})
	}

	objid := goutils.MustInt(ctx.Param("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))

	if objid == 0 || !isValidObjType(objtype) {
		return success(ctx, map[string]interface{}{"has_like": false})
	}

	likeFlags, err := logic.DefaultLike.FindUserLikeObjects(
		context.EchoContext(ctx), uid, objtype, objid, objid,
	)
	if err != nil {
		return success(ctx, map[string]interface{}{"has_like": false})
	}

	hasLike := likeFlags[objid] == model.FlagLike
	return success(ctx, map[string]interface{}{"has_like": hasLike})
}

// Toggle 点赞/取消点赞
func (LikeController) Toggle(ctx echo.Context) error {
	// 写操作必须校验用户状态（与 master NeedLogin 一致）
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
	}

	objid := goutils.MustInt(ctx.Param("objid"))
	objtype := goutils.MustInt(ctx.FormValue("objtype"))
	flag := goutils.MustInt(ctx.FormValue("flag"))

	if objid == 0 || !isValidObjType(objtype) {
		return fail(ctx, "参数错误")
	}

	// flag: 0=取消点赞, 1=点赞
	likeFlag := model.FlagCancel
	if flag == 1 {
		likeFlag = model.FlagLike
	}

	err = logic.DefaultLike.LikeObject(context.EchoContext(ctx), uid, objid, objtype, likeFlag)
	if err != nil {
		return fail(ctx, "操作失败")
	}

	return success(ctx, map[string]interface{}{
		"liked": flag == 1,
	})
}
