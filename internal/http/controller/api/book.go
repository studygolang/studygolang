// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"strconv"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type BookController struct{}

func (self BookController) RegisterRoute(g *echo.Group) {
	g.GET("/books", self.List)
	g.POST("/books", self.Publish)
	g.GET("/books/:id/edit", self.Edit)
	g.PUT("/books/:id", self.Update)
	g.GET("/books/:id", self.Detail)
}

// List 书籍列表
func (BookController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	books := logic.DefaultGoBook.FindAll(context.EchoContext(ctx), paginator, "likenum DESC,id DESC")
	total := logic.DefaultGoBook.Count(context.EchoContext(ctx))
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"books":    books,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// Detail 书籍详情
func (BookController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	book, err := logic.DefaultGoBook.FindById(context.EchoContext(ctx), id)
	if err != nil || book == nil || book.Id == 0 {
		return fail(ctx, "书籍不存在")
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		logic.Views.Incr(Request(ctx), model.TypeBook, book.Id, me.Uid)
	} else {
		logic.Views.Incr(Request(ctx), model.TypeBook, book.Id)
	}

	// 为了阅读数即时看到
	book.Viewnum++

	return success(ctx, map[string]interface{}{
		"book": book,
	})
}

// Publish 发布新图书（需要登录，支持 Cookie 和 X-Token header）
func (BookController) Publish(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	forms, _ := ctx.FormParams()

	// 基本字段验证
	name := forms.Get("name")
	if name == "" {
		return fail(ctx, "书名不能为空")
	}

	err = logic.DefaultGoBook.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "发布失败："+err.Error())
	}

	// 获取刚发布的图书 ID
	id := forms.Get("id")

	return success(ctx, map[string]interface{}{
		"id": id,
	})
}

// Edit 获取图书编辑数据（需要登录，验证权限）
func (BookController) Edit(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "图书 ID 非法")
	}

	book, err := logic.DefaultGoBook.FindById(context.EchoContext(ctx), id)
	if err != nil || book == nil || book.Id == 0 {
		return fail(ctx, "图书不存在")
	}

	// 验证权限
	if !logic.CanEdit(me, book) {
		return fail(ctx, "没有编辑权限")
	}

	return success(ctx, map[string]interface{}{
		"book": book,
	})
}

// Update 更新图书（需要登录，验证权限）
func (BookController) Update(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "图书 ID 非法")
	}

	book, err := logic.DefaultGoBook.FindById(context.EchoContext(ctx), id)
	if err != nil || book == nil || book.Id == 0 {
		return fail(ctx, "图书不存在")
	}

	// 验证权限
	if !logic.CanEdit(me, book) {
		return fail(ctx, "没有编辑权限")
	}

	forms, _ := ctx.FormParams()
	forms.Set("id", strconv.Itoa(book.Id))

	err = logic.DefaultGoBook.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "更新失败："+err.Error())
	}

	return success(ctx, map[string]interface{}{"id": book.Id})
}
