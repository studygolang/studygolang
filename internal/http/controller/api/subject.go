// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"strings"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/global"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type SubjectController struct{}

func (self SubjectController) RegisterRoute(g *echo.Group) {
	g.GET("/subjects", self.List)
	g.GET("/subject/:id", self.Index)
	g.POST("/subject/follow", self.Follow)
	g.GET("/subject/my_articles", self.MyArticles)
	g.POST("/subject/contribute", self.Contribute)
	g.POST("/subject/remove_contribute", self.RemoveContribute)
	g.GET("/subject/mine", self.Mine)
	g.POST("/subject/new", self.Create)
	g.POST("/subject/modify", self.Modify)
}

// Index 专题详情（含文章列表和关注者）
func (SubjectController) Index(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "专题 ID 不能为空")
	}

	subject := logic.DefaultSubject.FindOne(context.EchoContext(ctx), id)
	if subject.Id == 0 {
		return fail(ctx, "专题不存在")
	}

	// 补全封面 CDN URL
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
	hasMore := paginator.SetTotal(articleNum).HasMorePage()

	followers := logic.DefaultSubject.FindFollowers(context.EchoContext(ctx), id)
	followerNum := logic.DefaultSubject.FindFollowerTotal(context.EchoContext(ctx), id)

	// 检查当前用户是否已关注
	followed := false
	me, _ := requireAuth(ctx)
	if me != nil {
		followed = logic.DefaultSubject.HadFollow(context.EchoContext(ctx), id, me)
	}

	return success(ctx, map[string]interface{}{
		"subject":      subject,
		"articles":     articles,
		"article_num":  articleNum,
		"followers":    followers,
		"follower_num": followerNum,
		"order_by":     orderBy,
		"followed":     followed,
		"page":         curPage,
		"has_more":     hasMore,
	})
}

// Follow 关注/取消关注专题（需登录，form: sid）
func (SubjectController) Follow(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	sid := goutils.MustInt(ctx.FormValue("sid"))
	if sid == 0 {
		return fail(ctx, "专题 ID 不能为空")
	}

	err2 := logic.DefaultSubject.Follow(context.EchoContext(ctx), sid, me)
	if err2 != nil {
		return fail(ctx, "关注失败")
	}

	return success(ctx, nil)
}

// MyArticles 搜索当前用户的文章（需登录，用于投稿选择）
func (SubjectController) MyArticles(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	kw := ctx.QueryParam("kw")
	sid := goutils.MustInt(ctx.QueryParam("sid"))

	articles := logic.DefaultArticle.SearchMyArticles(context.EchoContext(ctx), me, sid, kw)

	return success(ctx, map[string]interface{}{
		"articles": articles,
	})
}

// Contribute 向专题投稿（需登录，form: sid, article_id）
func (SubjectController) Contribute(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	sid := goutils.MustInt(ctx.FormValue("sid"))
	articleId := goutils.MustInt(ctx.FormValue("article_id"))

	err2 := logic.DefaultSubject.Contribute(context.EchoContext(ctx), me, sid, articleId)
	if err2 != nil {
		return fail(ctx, err2.Error())
	}

	return success(ctx, nil)
}

// RemoveContribute 删除投稿（需登录 + 权限校验，form: sid, article_id）
func (SubjectController) RemoveContribute(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	sid := goutils.MustInt(ctx.FormValue("sid"))
	articleId := goutils.MustInt(ctx.FormValue("article_id"))

	// 权限校验：管理员或专题所有者
	if !me.IsRoot {
		subject := logic.DefaultSubject.FindOne(context.EchoContext(ctx), sid)
		if subject.Id == 0 {
			return fail(ctx, "专题不存在")
		}
		if subject.Uid != me.Uid {
			return fail(ctx, "无权删除此投稿")
		}
	}

	err2 := logic.DefaultSubject.RemoveContribute(context.EchoContext(ctx), sid, articleId)
	if err2 != nil {
		return fail(ctx, err2.Error())
	}

	return success(ctx, nil)
}

// Mine 我管理的专栏（需登录）
func (SubjectController) Mine(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	kw := ctx.QueryParam("kw")
	articleId := goutils.MustInt(ctx.QueryParam("article_id"))

	subjects := logic.DefaultSubject.FindMine(context.EchoContext(ctx), me, articleId, kw)

	return success(ctx, map[string]interface{}{
		"subjects": subjects,
	})
}

// Create 创建专栏（需登录 + 权限校验，form: name, cover, desc 等）
func (SubjectController) Create(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	if !me.IsRoot {
		return fail(ctx, "无权创建专栏")
	}

	name := ctx.FormValue("name")
	if name == "" {
		return fail(ctx, "专栏名称不能为空")
	}

	exist := logic.DefaultSubject.ExistByName(name)
	if exist {
		return fail(ctx, "专栏名称已存在")
	}

	forms, _ := ctx.FormParams()
	sid, err2 := logic.DefaultSubject.Publish(context.EchoContext(ctx), me, forms)
	if err2 != nil {
		getLogger(ctx).Errorln("subject create failed:", err2)
		return fail(ctx, "创建失败，请稍后重试")
	}

	return success(ctx, map[string]interface{}{
		"sid": sid,
	})
}

// Modify 修改专栏（需登录 + 所有者/管理员校验，form: sid + 其他字段）
func (SubjectController) Modify(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	sid := goutils.MustInt(ctx.FormValue("sid"))
	if sid == 0 {
		return fail(ctx, "专栏 ID 不能为空")
	}

	// 校验权限：管理员或专栏所有者
	if !me.IsRoot {
		subject := logic.DefaultSubject.FindOne(context.EchoContext(ctx), sid)
		if subject.Id == 0 {
			return fail(ctx, "专栏不存在")
		}
		if subject.Uid != me.Uid {
			return fail(ctx, "无权修改此专栏")
		}
	}

	forms, _ := ctx.FormParams()
	_, err2 := logic.DefaultSubject.Publish(context.EchoContext(ctx), me, forms)
	if err2 != nil {
		getLogger(ctx).Errorln("subject modify failed:", err2)
		return fail(ctx, "修改失败，请稍后重试")
	}

	return success(ctx, map[string]interface{}{
		"sid": sid,
	})
}

// List 专栏列表（分页）
func (SubjectController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	subjects := logic.DefaultSubject.FindBy(context.EchoContext(ctx), paginator)

	// 通过返回数量判断是否有更多
	hasMore := len(subjects) >= perPage

	return success(ctx, map[string]interface{}{
		"subjects": subjects,
		"page":     curPage,
		"has_more": hasMore,
	})
}
