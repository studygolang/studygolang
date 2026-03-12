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

type BookController struct{}

func (self BookController) RegisterRoute(g *echo.Group) {
	g.GET("/books", self.List)
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
