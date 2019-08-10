// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"html/template"
	"net/http"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/http"
	"github.com/studygolang/studygolang/http/middleware"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"
	"github.com/studygolang/studygolang/util"

	"github.com/dchest/captcha"
	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeResource, logic.ResourceComment{})
	logic.RegisterLikeObject(model.TypeResource, logic.ResourceLike{})
}

type ResourceController struct{}

// 注册路由
func (self ResourceController) RegisterRoute(g *echo.Group) {
	g.GET("/resources", self.ReadList)
	g.GET("/resources/cat/:catid", self.ReadCatResources)
	g.GET("/resources/:id", self.Detail)
	g.Match([]string{"GET", "POST"}, "/resources/new", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck(), middleware.PublishNotice(), middleware.CheckCaptcha())
	g.Match([]string{"GET", "POST"}, "/resources/modify", self.Modify, middleware.NeedLogin(), middleware.Sensivite())
}

// ReadList 资源索引页
func (ResourceController) ReadList(ctx echo.Context) error {
	return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
}

// ReadCatResources 某个分类的资源列表
func (ResourceController) ReadCatResources(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)
	catid := goutils.MustInt(ctx.Param("catid"))

	resources, total := logic.DefaultResource.FindByCatid(context.EchoContext(ctx), paginator, catid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL.Path)

	return render(ctx, "resources/index.html", map[string]interface{}{"activeResources": "active", "resources": resources, "categories": logic.AllCategory, "page": template.HTML(pageHtml), "curCatid": catid})
}

// Detail 某个资源详细页
func (ResourceController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
	}
	resource, comments := logic.DefaultResource.FindById(context.EchoContext(ctx), id)
	if len(resource) == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
	}

	data := map[string]interface{}{
		"activeResources": "active",
		"resource":        resource,
		"comments":        comments,
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		id := resource["id"].(int)
		data["likeflag"] = logic.DefaultLike.HadLike(context.EchoContext(ctx), me.Uid, id, model.TypeResource)
		data["hadcollect"] = logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), me.Uid, id, model.TypeResource)

		logic.Views.Incr(Request(ctx), model.TypeResource, id, me.Uid)

		if me.Uid != resource["uid"].(int) {
			go logic.DefaultViewRecord.Record(id, model.TypeResource, me.Uid)
		}

		if me.IsRoot || me.Uid == resource["uid"].(int) {
			data["view_user_num"] = logic.DefaultViewRecord.FindUserNum(context.EchoContext(ctx), id, model.TypeResource)
			data["view_source"] = logic.DefaultViewSource.FindOne(context.EchoContext(ctx), id, model.TypeResource)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeResource, id)
	}

	return render(ctx, "resources/detail.html,common/comment.html", data)
}

// Create 发布新资源
func (ResourceController) Create(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)

	title := ctx.FormValue("title")
	// 请求新建资源页面
	if title == "" || ctx.Request().Method != "POST" {
		data := map[string]interface{}{"activeResources": "active", "categories": logic.AllCategory}
		if logic.NeedCaptcha(me) {
			data["captchaId"] = captcha.NewLen(util.CaptchaLen)
		}
		return render(ctx, "resources/new.html", data)
	}

	errMsg := ""
	resForm := ctx.FormValue("form")
	if resForm == model.LinkForm {
		if ctx.FormValue("url") == "" {
			errMsg = "url不能为空"
		}
	} else {
		if ctx.FormValue("content") == "" {
			errMsg = "内容不能为空"
		}
	}
	if errMsg != "" {
		return fail(ctx, 1, errMsg)
	}

	forms, _ := ctx.FormParams()
	err := logic.DefaultResource.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, 2, "内部服务错误，请稍候再试！")
	}

	return success(ctx, nil)
}

// Modify 修改資源
func (ResourceController) Modify(ctx echo.Context) error {
	id := goutils.MustInt(ctx.FormValue("id"))
	if id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
	}

	// 请求编辑資源页面
	if ctx.Request().Method != "POST" {
		resource := logic.DefaultResource.FindResource(context.EchoContext(ctx), id)
		return render(ctx, "resources/new.html", map[string]interface{}{"resource": resource, "activeResources": "active", "categories": logic.AllCategory})
	}

	me := ctx.Get("user").(*model.Me)
	forms, _ := ctx.FormParams()
	err := logic.DefaultResource.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		if err == logic.NotModifyAuthorityErr {
			return ctx.String(http.StatusForbidden, "没有权限修改")
		}
		return fail(ctx, 2, "内部服务错误，请稍候再试！")
	}

	return success(ctx, nil)
}
