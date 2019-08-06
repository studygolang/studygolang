// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"net/http"
	"strings"

	"github.com/studygolang/studygolang/modules/context"
	"github.com/studygolang/studygolang/modules/global"
	. "github.com/studygolang/studygolang/modules/http"
	"github.com/studygolang/studygolang/modules/http/middleware"
	"github.com/studygolang/studygolang/modules/logic"
	"github.com/studygolang/studygolang/modules/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type SubjectController struct{}

// 注册路由
func (self SubjectController) RegisterRoute(g *echo.Group) {
	g.GET("/subject/:id", self.Index)
	g.POST("/subject/follow", self.Follow, middleware.NeedLogin())
	g.GET("/subject/my_articles", self.MyArticles, middleware.NeedLogin())
	g.POST("/subject/contribute", self.Contribute, middleware.NeedLogin())
	g.POST("/subject/remove_contribute", self.RemoveContribute, middleware.NeedLogin())
	g.GET("/subject/mine", self.Mine, middleware.NeedLogin())

	g.Match([]string{"GET", "POST"}, "/subject/new", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck(), middleware.PublishNotice())
	g.Match([]string{"GET", "POST"}, "/subject/modify", self.Modify, middleware.NeedLogin(), middleware.Sensivite())
}

func (SubjectController) Index(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	subject := logic.DefaultSubject.FindOne(context.EchoContext(ctx), id)
	if subject.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}
	if subject.Cover != "" && !strings.HasPrefix(subject.Cover, "http") {
		cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))
		subject.Cover = cdnDomain + subject.Cover
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	orderBy := ctx.QueryParam("order_by")
	articles := logic.DefaultSubject.FindArticles(context.EchoContext(ctx), id, paginator, orderBy)
	if orderBy == "" {
		orderBy = "added_at"
	}

	articleNum := logic.DefaultSubject.FindArticleTotal(context.EchoContext(ctx), id)

	pageHtml := paginator.SetTotal(articleNum).GetPageHtml(ctx.Request().URL.Path)

	followers := logic.DefaultSubject.FindFollowers(context.EchoContext(ctx), id)
	followerNum := logic.DefaultSubject.FindFollowerTotal(context.EchoContext(ctx), id)

	// 是否已关注
	followed := false
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		followed = logic.DefaultSubject.HadFollow(context.EchoContext(ctx), id, me)
	}

	data := map[string]interface{}{
		"subject":      subject,
		"articles":     articles,
		"article_num":  articleNum,
		"followers":    followers,
		"follower_num": followerNum,
		"order_by":     orderBy,
		"followed":     followed,
		"page":         pageHtml,
	}

	return render(ctx, "subject/index.html", data)
}

func (self SubjectController) Follow(ctx echo.Context) error {
	sid := goutils.MustInt(ctx.FormValue("sid"))

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultSubject.Follow(context.EchoContext(ctx), sid, me)
	if err != nil {
		return fail(ctx, 1, "关注失败！")
	}

	return success(ctx, nil)
}

func (self SubjectController) MyArticles(ctx echo.Context) error {
	kw := ctx.QueryParam("kw")
	sid := goutils.MustInt(ctx.FormValue("sid"))

	me := ctx.Get("user").(*model.Me)

	articles := logic.DefaultArticle.SearchMyArticles(context.EchoContext(ctx), me, sid, kw)

	return success(ctx, map[string]interface{}{
		"articles": articles,
	})
}

// Contribute 投稿
func (self SubjectController) Contribute(ctx echo.Context) error {
	sid := goutils.MustInt(ctx.FormValue("sid"))
	articleId := goutils.MustInt(ctx.FormValue("article_id"))

	me := ctx.Get("user").(*model.Me)

	err := logic.DefaultSubject.Contribute(context.EchoContext(ctx), me, sid, articleId)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}

// RemoveContribute 删除投稿
func (self SubjectController) RemoveContribute(ctx echo.Context) error {
	sid := goutils.MustInt(ctx.FormValue("sid"))
	articleId := goutils.MustInt(ctx.FormValue("article_id"))

	err := logic.DefaultSubject.RemoveContribute(context.EchoContext(ctx), sid, articleId)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}

// Mine 我管理的专栏
func (self SubjectController) Mine(ctx echo.Context) error {
	kw := ctx.QueryParam("kw")
	articleId := goutils.MustInt(ctx.FormValue("article_id"))
	me := ctx.Get("user").(*model.Me)

	subjects := logic.DefaultSubject.FindMine(context.EchoContext(ctx), me, articleId, kw)

	return success(ctx, map[string]interface{}{"subjects": subjects})
}

// Create 新建专栏
func (SubjectController) Create(ctx echo.Context) error {

	name := ctx.FormValue("name")
	// 请求新建专栏页面
	if name == "" || ctx.Request().Method != "POST" {
		data := map[string]interface{}{}
		return render(ctx, "subject/new.html", data)
	}

	exist := logic.DefaultSubject.ExistByName(name)
	if exist {
		return fail(ctx, 1, "专栏已经存在 : "+name)
	}

	me := ctx.Get("user").(*model.Me)
	forms, _ := ctx.FormParams()
	sid, err := logic.DefaultSubject.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, 1, "内部服务错误:"+err.Error())
	}

	return success(ctx, map[string]interface{}{"sid": sid})
}

// Modify 修改专栏
func (SubjectController) Modify(ctx echo.Context) error {
	sid := goutils.MustInt(ctx.FormValue("sid"))
	if sid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/subjects")
	}

	if ctx.Request().Method != "POST" {
		subject := logic.DefaultSubject.FindOne(context.EchoContext(ctx), sid)
		if subject == nil {
			return ctx.Redirect(http.StatusSeeOther, "/subjects")
		}

		data := map[string]interface{}{
			"subject": subject,
		}

		return render(ctx, "subject/new.html", data)
	}

	me := ctx.Get("user").(*model.Me)
	forms, _ := ctx.FormParams()
	_, err := logic.DefaultSubject.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, 2, "服务错误，请稍后重试！")
	}
	return success(ctx, map[string]interface{}{"sid": sid})
}
