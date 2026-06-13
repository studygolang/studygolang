// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"strings"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type ArticleController struct{}

// RegisterRoute 注册路由
func (self *ArticleController) RegisterRoute(g *echo.Group) {
	g.GET("/articles", self.List)
	g.GET("/articles/crawl", self.Crawl)
	g.GET("/articles/:id", self.Detail)
	g.GET("/articles/:id/edit", self.Edit)
	g.PUT("/articles/:id", self.Update)
	g.POST("/articles", self.Create)
}

// List 文章列表，支持 p 分页、tag 过滤和 sort 排序参数
func (ArticleController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	// 解析 tag 过滤参数
	tag := strings.TrimSpace(ctx.QueryParam("tag"))
	var tagCondition string
	var tagArgs []interface{}
	if tag != "" {
		tagCondition = "tags LIKE ?"
		tagArgs = []interface{}{"%" + tag + "%"}
	}

	// 解析 sort 排序参数：hot(热门), latest(最新), noreply(无回复)
	sort := strings.ToLower(strings.TrimSpace(ctx.QueryParam("sort")))
	orderBy := getArticleSortOrder(sort)

	// 置顶文章单独查询，不受 sort 参数影响（置顶优先）
	topArticles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "top=1")
	articles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, orderBy, tagCondition, tagArgs...)

	total := logic.DefaultArticle.Count(context.EchoContext(ctx), tagCondition, tagArgs...)
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"list":     append(topArticles, articles...),
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// getArticleSortOrder 根据 sort 参数返回对应的排序 SQL
func getArticleSortOrder(sort string) string {
	switch sort {
	case "hot":
		// 热门：评论数 > 点赞数 > 浏览数 > ID
		return "cmtnum DESC, likenum DESC, viewnum DESC, id DESC"
	case "noreply":
		// 无回复：优先展示评论数为 0 的文章，按 ID 倒序
		return "cmtnum ASC, id DESC"
	case "latest":
		fallthrough
	default:
		// 最新：按 ID 倒序
		return "id DESC"
	}
}

// Detail 文章详情，增加浏览量
func (ArticleController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "文章 id 非法")
	}

	article, prevNext, err := logic.DefaultArticle.FindByIdAndPreNext(context.EchoContext(ctx), id)
	if err != nil {
		return fail(ctx, "获取文章失败")
	}

	if article == nil || article.Id == 0 || article.Status == model.ArticleStatusOffline {
		return success(ctx, map[string]interface{}{"article": map[string]interface{}{"id": 0}})
	}

	article.Viewnum++

	articleGCTT := logic.DefaultArticle.FindArticleGCTT(context.EchoContext(ctx), article)

	replies, _, lastReplyUser := logic.DefaultComment.FindObjComments(
		context.EchoContext(ctx), article.Id, model.TypeArticle, 0, article.Lastreplyuid,
	)
	if article.Lastreplyuid != 0 {
		article.LastReplyUser = lastReplyUser
	}

	article.Txt = ""

	result := map[string]interface{}{
		"article":      article,
		"article_gctt": articleGCTT,
		"replies":      normalizeReplies(replies),
		"prev_next":    prevNext,
		"subjects":     logic.DefaultSubject.FindArticleSubjects(context.EchoContext(ctx), article.Id),
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		result["likeflag"] = logic.DefaultLike.HadLike(context.EchoContext(ctx), me.Uid, article.Id, model.TypeArticle)
		result["hadcollect"] = logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), me.Uid, article.Id, model.TypeArticle)

		logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id, me.Uid)

		if !article.IsSelf || me.Uid != article.User.Uid {
			go logic.DefaultViewRecord.Record(article.Id, model.TypeArticle, me.Uid)
		}

		if me.IsRoot || (article.IsSelf && me.Uid == article.User.Uid) {
			result["view_user_num"] = logic.DefaultViewRecord.FindUserNum(context.EchoContext(ctx), article.Id, model.TypeArticle)
			result["view_source"] = logic.DefaultViewSource.FindOne(context.EchoContext(ctx), article.Id, model.TypeArticle)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id)
	}

	return success(ctx, result)
}

// Edit 获取文章编辑数据（需要登录，验证权限）
func (ArticleController) Edit(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := ctx.Param("id")
	article, err := logic.DefaultArticle.FindById(context.EchoContext(ctx), id)
	if err != nil || article == nil || article.Id == 0 {
		return fail(ctx, "文章不存在")
	}

	// 验证权限：只能编辑自己的文章，或管理员可以编辑所有文章
	if !logic.CanEdit(me, article) {
		return fail(ctx, "没有编辑权限")
	}

	return success(ctx, map[string]interface{}{
		"article": article,
	})
}

// Update 更新文章（需要登录，验证权限）
func (ArticleController) Update(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := ctx.Param("id")
	article, err := logic.DefaultArticle.FindById(context.EchoContext(ctx), id)
	if err != nil || article == nil || article.Id == 0 {
		return fail(ctx, "文章不存在")
	}

	// 验证权限
	if !logic.CanEdit(me, article) {
		return fail(ctx, "没有编辑权限")
	}

	// 敏感词检查
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
	}

	forms, _ := ctx.FormParams()
	forms.Set("id", id)
	errMsg, err := logic.DefaultArticle.Modify(context.EchoContext(ctx), me, forms)
	if err != nil {
		if errMsg != "" {
			return fail(ctx, errMsg)
		}
		return fail(ctx, "更新失败")
	}

	return success(ctx, map[string]interface{}{"id": article.Id})
}

// Create 创建文章（需要登录）
func (ArticleController) Create(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	// 敏感词检查
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
	}
	// 余额检查
	if !balanceCheck(me, false) {
		return failBalance(ctx)
	}
	// 验证码检查（新注册用户或频繁发布时需要）
	if !captchaCheck(ctx, me) {
		return failCaptcha(ctx)
	}

	// 获取表单数据
	forms, _ := ctx.FormParams()

	// 验证必填字段
	title := ctx.FormValue("title")
	if title == "" {
		return fail(ctx, "标题不能为空")
	}

	content := ctx.FormValue("content")
	if content == "" {
		return fail(ctx, "内容不能为空")
	}

	// 调用业务逻辑发布文章
	id, err := logic.DefaultArticle.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "发布失败: "+err.Error())
	}

	// 发布后邮件通知站长
	publishNotice(ctx, me)

	return success(ctx, map[string]interface{}{"id": id})
}

// Crawl 抓取外部文章内容（需要登录）
func (ArticleController) Crawl(ctx echo.Context) error {
	if _, err := requireAuth(ctx); err != nil {
		return err
	}

	strUrl := strings.TrimSpace(ctx.QueryParam("url"))
	if strUrl == "" {
		return fail(ctx, "url 参数不能为空")
	}

	article, err := logic.DefaultArticle.ParseArticle(context.EchoContext(ctx), strUrl, false)
	if err != nil {
		return fail(ctx, "抓取文章失败: "+err.Error())
	}

	return success(ctx, map[string]interface{}{
		"title":   article.Title,
		"content": article.Content,
		"author":  article.Author,
	})
}
