// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"html/template"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/http"
	"github.com/studygolang/studygolang/http/middleware"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"
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
	g.GET("/books", self.ReadList)

	g.GET("/book/:id", self.Detail)

	g.Match([]string{"GET", "POST"}, "/book/new", self.Create, middleware.NeedLogin(), middleware.BalanceCheck(), middleware.PublishNotice())
}

// ReadList 图书列表页
func (BookController) ReadList(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	books := logic.DefaultGoBook.FindAll(context.EchoContext(ctx), paginator, "likenum DESC,id DESC")
	total := logic.DefaultGoBook.Count(context.EchoContext(ctx))
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL.Path)

	data := map[string]interface{}{
		"books":       books,
		"activeBooks": "active",
		"page":        template.HTML(pageHtml),
	}

	return render(ctx, "books/list.html", data)
}

// Create 发布新书
func (BookController) Create(ctx echo.Context) error {
	name := ctx.FormValue("name")
	// 请求新建图书页面
	if name == "" || ctx.Request().Method != "POST" {
		book := &model.Book{}
		return render(ctx, "books/new.html", map[string]interface{}{"book": book, "activeBooks": "active"})
	}

	user := ctx.Get("user").(*model.Me)
	forms, _ := ctx.FormParams()
	err := logic.DefaultGoBook.Publish(context.EchoContext(ctx), user, forms)
	if err != nil {
		return fail(ctx, 1, "内部服务错误！")
	}
	return success(ctx, nil)
}

// Detail 图书详细页
func (BookController) Detail(ctx echo.Context) error {
	book, err := logic.DefaultGoBook.FindById(context.EchoContext(ctx), ctx.Param("id"))
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/books")
	}

	if book == nil || book.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/books")
	}

	data := map[string]interface{}{
		"activeBooks": "active",
		"book":        book,
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		data["likeflag"] = logic.DefaultLike.HadLike(context.EchoContext(ctx), me.Uid, book.Id, model.TypeBook)
		data["hadcollect"] = logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), me.Uid, book.Id, model.TypeBook)

		logic.Views.Incr(Request(ctx), model.TypeBook, book.Id, me.Uid)

		if me.Uid != book.Uid {
			go logic.DefaultViewRecord.Record(book.Id, model.TypeBook, me.Uid)
		}

		if me.IsRoot || me.Uid == book.Uid {
			data["view_user_num"] = logic.DefaultViewRecord.FindUserNum(context.EchoContext(ctx), book.Id, model.TypeBook)
			data["view_source"] = logic.DefaultViewSource.FindOne(context.EchoContext(ctx), book.Id, model.TypeBook)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeBook, book.Id)
	}

	// 为了阅读数即时看到
	book.Viewnum++

	return render(ctx, "books/detail.html,common/comment.html", data)
}
