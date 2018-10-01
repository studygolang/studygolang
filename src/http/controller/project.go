// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"http/middleware"
	"logic"
	"net/http"
	"util"

	"github.com/dchest/captcha"
	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"

	"html/template"
	. "http"
	"model"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeProject, logic.ProjectComment{})
	logic.RegisterLikeObject(model.TypeProject, logic.ProjectLike{})
}

type ProjectController struct{}

// 注册路由
func (self ProjectController) RegisterRoute(g *echo.Group) {
	g.GET("/projects", self.ReadList)
	g.Match([]string{"GET", "POST"}, "/project/new", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck(), middleware.PublishNotice(), middleware.CheckCaptcha())
	g.Match([]string{"GET", "POST"}, "/project/modify", self.Modify, middleware.NeedLogin(), middleware.Sensivite())
	g.GET("/p/:uri", self.Detail)
	g.GET("/project/uri", self.CheckExist)
}

// ReadList 开源项目列表页
func (ProjectController) ReadList(ctx echo.Context) error {
	limit := 20

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)
	paginator.SetPerPage(limit)
	total := logic.DefaultProject.Count(ctx, "")
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())
	pageInfo := template.HTML(pageHtml)

	projects := logic.DefaultProject.FindAll(ctx, paginator, "id DESC", "status IN(?,?)", model.ProjectStatusNew, model.ProjectStatusOnline)

	num := len(projects)
	if num == 0 {
		if curPage == int(total) {
			return render(ctx, "projects/list.html", map[string]interface{}{"projects": projects, "activeProjects": "active"})
		} else {
			return ctx.Redirect(http.StatusSeeOther, "/projects")
		}
	}

	// 获取当前用户喜欢对象信息
	me, ok := ctx.Get("user").(*model.Me)
	var likeFlags map[int]int
	if ok {
		likeFlags, _ = logic.DefaultLike.FindUserLikeObjects(ctx, me.Uid, model.TypeProject, projects[0].Id, projects[num-1].Id)
	}

	return render(ctx, "projects/list.html", map[string]interface{}{"projects": projects, "activeProjects": "active", "page": pageInfo, "likeflags": likeFlags})
}

// Create 新建项目
func (ProjectController) Create(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)

	name := ctx.FormValue("name")
	// 请求新建项目页面
	if name == "" || ctx.Request().Method() != "POST" {
		project := &model.OpenProject{}

		data := map[string]interface{}{"project": project, "activeProjects": "active"}
		if logic.NeedCaptcha(me) {
			data["captchaId"] = captcha.NewLen(util.CaptchaLen)
		}

		return render(ctx, "projects/new.html", data)
	}

	err := logic.DefaultProject.Publish(ctx, me, ctx.FormParams())
	if err != nil {
		return fail(ctx, 1, "内部服务错误！")
	}
	return success(ctx, nil)
}

// Modify 修改项目
func (ProjectController) Modify(ctx echo.Context) error {
	id := goutils.MustInt(ctx.FormValue("id"))
	if id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/projects")
	}

	// 请求编辑项目页面
	if ctx.Request().Method() != "POST" {
		project := logic.DefaultProject.FindOne(ctx, id)
		return render(ctx, "projects/new.html", map[string]interface{}{"project": project, "activeProjects": "active"})
	}

	user := ctx.Get("user").(*model.Me)
	err := logic.DefaultProject.Publish(ctx, user, ctx.FormParams())
	if err != nil {
		if err == logic.NotModifyAuthorityErr {
			return ctx.String(http.StatusForbidden, "没有权限")
		}
		return fail(ctx, 1, "内部服务错误！")
	}
	return success(ctx, nil)
}

// Detail 项目详情
func (ProjectController) Detail(ctx echo.Context) error {
	project := logic.DefaultProject.FindOne(ctx, ctx.Param("uri"))
	if project == nil || project.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/projects")
	}

	data := map[string]interface{}{
		"activeProjects": "active",
		"project":        project,
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		data["likeflag"] = logic.DefaultLike.HadLike(ctx, me.Uid, project.Id, model.TypeProject)
		data["hadcollect"] = logic.DefaultFavorite.HadFavorite(ctx, me.Uid, project.Id, model.TypeProject)

		logic.Views.Incr(Request(ctx), model.TypeProject, project.Id, me.Uid)

		if me.Uid != project.User.Uid {
			go logic.DefaultViewRecord.Record(project.Id, model.TypeProject, me.Uid)
		}

		if me.IsRoot || me.Uid == project.User.Uid {
			data["view_user_num"] = logic.DefaultViewRecord.FindUserNum(ctx, project.Id, model.TypeProject)
			data["view_source"] = logic.DefaultViewSource.FindOne(ctx, project.Id, model.TypeProject)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeProject, project.Id)
	}

	// 为了阅读数即时看到
	project.Viewnum++

	return render(ctx, "projects/detail.html,common/comment.html", data)
}

// CheckExist 检测 uri 对应的项目是否存在(验证，true表示不存在；false表示存在)
func (ProjectController) CheckExist(ctx echo.Context) error {
	uri := ctx.QueryParam("uri")
	if uri == "" {
		return ctx.JSON(http.StatusOK, `true`)
	}

	if logic.DefaultProject.UriExists(ctx, uri) {
		return ctx.JSON(http.StatusOK, `false`)
	}
	return ctx.JSON(http.StatusOK, `true`)
}
