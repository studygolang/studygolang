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

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type CommentController struct{}

func (self CommentController) RegisterRoute(g *echo.Group) {
	g.GET("/comments", self.List)
	g.POST("/comments/:objid", self.Create)
}

// List 评论列表（objid、objtype 查询参数）
func (CommentController) List(ctx echo.Context) error {
	objid := goutils.MustInt(ctx.QueryParam("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))

	if objid == 0 {
		return fail(ctx, "参数有误")
	}

	comments, _, _ := logic.DefaultComment.FindObjComments(
		context.EchoContext(ctx), objid, objtype, 0, 0,
	)

	return success(ctx, map[string]interface{}{
		"comments": comments,
	})
}

// Create 创建评论（支持 Cookie 和 X-Token header）
func (CommentController) Create(ctx echo.Context) error {
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
	if objid == 0 {
		return fail(ctx, "参数有误")
	}

	forms, _ := ctx.FormParams()
	comment, err := logic.DefaultComment.Publish(context.EchoContext(ctx), uid, objid, forms)
	if err != nil {
		return fail(ctx, "发布评论失败", 2)
	}

	return success(ctx, map[string]interface{}{
		"comment": comment,
	})
}
