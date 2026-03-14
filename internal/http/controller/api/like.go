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
}

// Toggle 点赞/取消点赞
func (LikeController) Toggle(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	objid := goutils.MustInt(ctx.Param("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))
	flag := goutils.MustInt(ctx.QueryParam("flag"))

	if objid == 0 || objtype == 0 {
		return fail(ctx, "参数错误")
	}

	// flag: 0=取消点赞, 1=点赞
	likeFlag := model.FlagCancel
	if flag == 1 {
		likeFlag = model.FlagLike
	}

	err := logic.DefaultLike.LikeObject(context.EchoContext(ctx), uid, objid, objtype, likeFlag)
	if err != nil {
		return fail(ctx, "操作失败")
	}

	return success(ctx, map[string]interface{}{
		"liked": flag == 1,
	})
}
