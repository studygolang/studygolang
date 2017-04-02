// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"

	. "http"
	"model"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeBook, logic.BookComment{})
	logic.RegisterLikeObject(model.TypeBook, logic.BookLike{})
}

type BookController struct{}

// 注册路由
func (self BookController) RegisterRoute(g *echo.Group) {
	g.Get("/books", self.ReadList)

	g.Get("/book/:id", self.Detail)
}

// ReadList 图书列表页
func (BookController) ReadList(ctx echo.Context) error {
	limit := 20

	lastId := goutils.MustInt(ctx.QueryParam("lastid"))
	books := logic.DefaultGoBook.FindBy(ctx, limit+1, lastId)
	if books == nil {
		logger.Errorln("book controller: find book error")
		return ctx.Redirect(http.StatusSeeOther, "/books")
	}

	num := len(books)
	if num == 0 {
		if lastId == 0 {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		return ctx.Redirect(http.StatusSeeOther, "/books")
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId > 0 {
		prevId = lastId

		if prevId-books[0].Id > 1 {
			hasPrev = false
		} else {
			prevId += limit
			hasPrev = true
		}
	}

	if num > limit {
		hasNext = true
		books = books[:limit]
		nextId = books[limit-1].Id
	} else {
		nextId = books[num-1].Id
	}

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}

	return render(ctx, "books/list.html", map[string]interface{}{"books": books, "activeBooks": "active", "page": pageInfo})
}

// Detail 图书详细页
func (BookController) Detail(ctx echo.Context) error {
	book, err := logic.DefaultGoBook.FindById(ctx, ctx.Param("id"))
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/books")
	}

	if book == nil || book.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/books")
	}

	likeFlag := 0
	hadCollect := 0
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		likeFlag = logic.DefaultLike.HadLike(ctx, me.Uid, book.Id, model.TypeBook)
		hadCollect = logic.DefaultFavorite.HadFavorite(ctx, me.Uid, book.Id, model.TypeBook)
	}

	logic.Views.Incr(Request(ctx), model.TypeBook, book.Id)

	// 为了阅读数即时看到
	book.Viewnum++

	return render(ctx, "books/detail.html,common/comment.html", map[string]interface{}{"activeBooks": "active", "book": book, "likeflag": likeFlag, "hadcollect": hadCollect})
}
