// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"http/middleware"
	"logic"
	"model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/slices"
)

type CommentController struct{}

func (self CommentController) RegisterRoute(g *echo.Group) {
	g.Get("/at/users", self.AtUsers)
	g.Post("/comment/:objid", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.PublishNotice())
	g.Get("/object/comments", self.CommentList)
}

// AtUsers 评论或回复 @ 某人 suggest
func (CommentController) AtUsers(ctx echo.Context) error {
	term := ctx.QueryParam("term")
	users := logic.DefaultUser.GetUserMentions(term, 10)
	return ctx.JSON(http.StatusOK, users)
}

// Create 评论（或回复）
func (CommentController) Create(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Me)

	// 入库
	objid := goutils.MustInt(ctx.Param("objid"))
	comment, err := logic.DefaultComment.Publish(ctx, user.Uid, objid, ctx.FormParams())
	if err != nil {
		return fail(ctx, 1, "服务器内部错误")
	}

	return success(ctx, comment)
}

// CommentList 获取某对象的评论信息
func (CommentController) CommentList(ctx echo.Context) error {
	objid := goutils.MustInt(ctx.QueryParam("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))

	commentList, err := logic.DefaultComment.FindObjectComments(ctx, objid, objtype)
	if err != nil {
		return fail(ctx, 1, "服务器内部错误")
	}

	uids := slices.StructsIntSlice(commentList, "Uid")
	users := logic.DefaultUser.FindUserInfos(ctx, uids)

	result := map[string]interface{}{
		"comments": commentList,
	}

	// json encode 不支持 map[int]...
	for uid, user := range users {
		result[strconv.Itoa(uid)] = user
	}

	return success(ctx, result)
}
