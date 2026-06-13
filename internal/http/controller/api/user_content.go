// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

func buildUserInfo(user *model.User) map[string]interface{} {
	return map[string]interface{}{
		"username": user.Username,
		"name":     user.Name,
		"avatar":   user.Avatar,
	}
}

type UserContentController struct{}

func (self UserContentController) RegisterRoute(g *echo.Group) {
	g.GET("/users/:username/comments", self.UserComments)
	g.GET("/users/:username/topics", self.UserTopics)
	g.GET("/users/:username/articles", self.UserArticles)
	g.GET("/users/:username/resources", self.UserResources)
	g.GET("/users/:username/projects", self.UserProjects)
}

// UserComments 获取指定用户的评论列表
// GET /api/v1/users/:username/comments?p=1
func (UserContentController) UserComments(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	comments := logic.DefaultComment.FindAll(context.EchoContext(ctx), paginator, "cid DESC", "uid=?", user.Uid)
	total := logic.DefaultComment.Count(context.EchoContext(ctx), "uid=?", user.Uid)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"user":     buildUserInfo(user),
		"comments": comments,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// UserTopics 获取指定用户的话题列表
// GET /api/v1/users/:username/topics?p=1
func (UserContentController) UserTopics(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	// FindAll JOIN topics_ex（同样含 tid），裸 tid 会 Ambiguous，需带表名前缀
	topics := logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.tid DESC", "uid=?", user.Uid)
	total := logic.DefaultTopic.Count(context.EchoContext(ctx), "uid=?", user.Uid)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"user":     buildUserInfo(user),
		"topics":   topics,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// UserArticles 获取指定用户的文章列表
// GET /api/v1/users/:username/articles?p=1
func (UserContentController) UserArticles(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	articles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "author_txt=?", username)
	total := logic.DefaultArticle.Count(context.EchoContext(ctx), "author_txt=?", username)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"user":     buildUserInfo(user),
		"articles": articles,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// UserResources 获取指定用户的资源列表
// GET /api/v1/users/:username/resources?p=1
func (UserContentController) UserResources(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	// FindAll JOIN resource_ex（同样含 id），裸 id 会 Ambiguous，需带表名前缀
	resources, total := logic.DefaultResource.FindAll(context.EchoContext(ctx), paginator, "resource.id DESC", "uid=?", user.Uid)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"user":      buildUserInfo(user),
		"resources": resources,
		"total":     total,
		"page":      curPage,
		"has_more":  hasMore,
	})
}

// UserProjects 获取指定用户的开源项目列表
// GET /api/v1/users/:username/projects?p=1
func (UserContentController) UserProjects(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	projects := logic.DefaultProject.FindAll(context.EchoContext(ctx), paginator, "id DESC", "username=?", username)
	total := logic.DefaultProject.Count(context.EchoContext(ctx), "username=?", username)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"user":     buildUserInfo(user),
		"projects": projects,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}
