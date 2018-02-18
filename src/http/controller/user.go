// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"html/template"
	"logic"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/slices"
)

type UserController struct{}

// 注册路由
func (self UserController) RegisterRoute(g *echo.Group) {
	g.GET("/user/:username", self.Home)
	g.GET("/user/:username/topics", self.Topics)
	g.GET("/user/:username/articles", self.Articles)
	g.GET("/user/:username/resources", self.Resources)
	g.GET("/user/:username/projects", self.Projects)
	g.GET("/user/:username/comments", self.Comments)
	g.GET("/users", self.ReadList)
	g.Match([]string{"GET", "POST"}, "/user/email/unsubscribe", self.EmailUnsub)
}

// Home 用户个人首页
func (UserController) Home(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(ctx, "username", username)
	if user == nil || user.Uid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/users")
	}

	user.Weight = logic.DefaultRank.UserDAURank(ctx, user.Uid)

	topics := logic.DefaultTopic.FindRecent(5, user.Uid)

	articles := logic.DefaultArticle.FindByUser(ctx, user.Username, 5)

	resources := logic.DefaultResource.FindRecent(ctx, user.Uid)
	for _, resource := range resources {
		resource.CatName = logic.GetCategoryName(resource.Catid)
	}

	projects := logic.DefaultProject.FindRecent(ctx, user.Username)
	comments := logic.DefaultComment.FindRecent(ctx, user.Uid, -1, 5)

	user.IsOnline = logic.Book.RegUserIsOnline(user.Uid)

	return render(ctx, "user/profile.html", map[string]interface{}{
		"activeUsers": "active",
		"topics":      topics,
		"articles":    articles,
		"resources":   resources,
		"projects":    projects,
		"comments":    comments,
		"user":        user,
	})
}

// ReadList 会员列表
func (UserController) ReadList(ctx echo.Context) error {
	// 获取活跃会员
	// activeUsers := logic.DefaultUser.FindActiveUsers(ctx, 36)
	activeUsers := logic.DefaultRank.FindDAURank(ctx, 36)
	// 获取最新加入会员
	newUsers := logic.DefaultUser.FindNewUsers(ctx, 36)
	// 获取会员总数
	total := logic.DefaultUser.Total()

	return render(ctx, "user/users.html", map[string]interface{}{"activeUsers": "active", "actives": activeUsers, "news": newUsers, "total": total})
}

// EmailUnsub 邮件订阅/退订页面
func (UserController) EmailUnsub(ctx echo.Context) error {
	token := ctx.FormValue("u")
	if token == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	// 校验 token 的合法性
	email := ctx.FormValue("email")
	user := logic.DefaultUser.FindOne(ctx, "email", email)
	if user.Email == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	realToken := logic.DefaultEmail.GenUnsubscribeToken(user)
	if token != realToken {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	if ctx.Request().Method() != "POST" {
		data := map[string]interface{}{
			"email":       email,
			"token":       token,
			"unsubscribe": user.Unsubscribe,
		}

		return render(ctx, "user/email_unsub.html", data)
	}

	logic.DefaultUser.EmailSubscribe(ctx, user.Uid, goutils.MustInt(ctx.FormValue("unsubscribe")))

	return success(ctx, nil)
}

func (UserController) Topics(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(ctx, "username", username)
	if user == nil || user.Uid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/users")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	querystring := "uid=?"
	topics := logic.DefaultTopic.FindAll(ctx, paginator, "topics.tid DESC", querystring, user.Uid)
	total := logic.DefaultTopic.Count(ctx, querystring, user.Uid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	return render(ctx, "user/topics.html", map[string]interface{}{
		"user":         user,
		"activeTopics": "active",
		"topics":       topics,
		"page":         template.HTML(pageHtml),
		"total":        total,
	})
}

func (UserController) Articles(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(ctx, "username", username)
	if user == nil || user.Uid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/users")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	querystring := "author_txt=?"
	articles := logic.DefaultArticle.FindAll(ctx, paginator, "id DESC", querystring, user.Username)
	total := logic.DefaultArticle.Count(ctx, querystring, user.Username)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	return render(ctx, "user/articles.html", map[string]interface{}{
		"user":           user,
		"activeArticles": "active",
		"articles":       articles,
		"page":           template.HTML(pageHtml),
		"total":          total,
	})
}

func (UserController) Resources(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(ctx, "username", username)
	if user == nil || user.Uid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/users")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	querystring := "uid=?"
	resources, total := logic.DefaultResource.FindAll(ctx, paginator, "resource.id DESC", querystring, user.Uid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	return render(ctx, "user/resources.html", map[string]interface{}{
		"user":            user,
		"activeResources": "active",
		"resources":       resources,
		"page":            template.HTML(pageHtml),
		"total":           total,
	})
}

func (UserController) Projects(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(ctx, "username", username)
	if user == nil || user.Uid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/users")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	querystring := "username=?"
	projects := logic.DefaultProject.FindAll(ctx, paginator, "id DESC", querystring, user.Username)
	total := logic.DefaultProject.Count(ctx, querystring, user.Username)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	return render(ctx, "user/projects.html", map[string]interface{}{
		"user":           user,
		"activeProjects": "active",
		"projects":       projects,
		"page":           template.HTML(pageHtml),
		"total":          total,
	})
}
func (UserController) Comments(ctx echo.Context) error {

	username := ctx.Param("username")

	userid := 0
	querystring := ""

	if username != "0" {
		user := logic.DefaultUser.FindOne(ctx, "username", username)
		if user == nil || user.Uid == 0 {
			return ctx.Redirect(http.StatusSeeOther, "/users")
		}
		querystring = "uid=?"
		userid = user.Uid
		username = user.Username
	} else {
		username = ""
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	comments := logic.DefaultComment.FindAll(ctx, paginator, "cid DESC", querystring, userid)

	total := logic.DefaultComment.Count(ctx, querystring, userid)

	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	data := map[string]interface{}{
		"comments": comments,
		"page":     template.HTML(pageHtml),
		"total":    total,
	}

	if username == "" {
		uids := slices.StructsIntSlice(comments, "Uid")
		data["users"] = logic.DefaultUser.FindUserInfos(ctx, uids)
	}

	return render(ctx, "user/comments.html", data)

}
