// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"errors"
	"http/middleware"
	"logic"
	"model"
	"net/http"
	"strconv"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/echoutils"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/slices"
)

type CommentController struct{}

func (self CommentController) RegisterRoute(g *echo.Group) {
	g.Get("/at/users", self.AtUsers)
	g.Post("/comment/:objid", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck(), middleware.PublishNotice())
	g.Get("/object/comments", self.CommentList)
	g.Post("/object/comments/:cid", self.Modify, middleware.NeedLogin(), middleware.Sensivite())

	g.Get("/topics/:objid/comment/:cid", self.TopicDetail)
	g.Get("/articles/:objid/comment/:cid", self.ArticleDetail)
}

// AtUsers 评论或回复 @ 某人 suggest
func (CommentController) AtUsers(ctx echo.Context) error {
	term := ctx.QueryParam("term")
	isHttps := CheckIsHttps(ctx)
	users := logic.DefaultUser.GetUserMentions(term, 10, isHttps)
	return ctx.JSON(http.StatusOK, users)
}

// Create 评论（或回复）
func (CommentController) Create(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Me)

	// 入库
	objid := goutils.MustInt(ctx.Param("objid"))
	if objid == 0 {
		return fail(ctx, 1, "参数有误，请刷新后重试！")
	}
	comment, err := logic.DefaultComment.Publish(ctx, user.Uid, objid, ctx.FormParams())
	if err != nil {
		return fail(ctx, 2, "服务器内部错误")
	}

	return success(ctx, comment)
}

// 修改评论
func (CommentController) Modify(ctx echo.Context) error {
	cid := goutils.MustInt(ctx.Param("cid"))
	content := ctx.FormValue("content")
	comment, err := logic.DefaultComment.FindById(cid)

	if err != nil {
		return fail(ctx, 2, "评论不存在")
	}

	if content == "" {
		return fail(ctx, 1, "内容不能为空")
	}

	me := ctx.Get("user").(*model.Me)
	// CanEdit 已包含修改时间限制
	if !logic.CanEdit(me, comment) {
		return fail(ctx, 3, "没有修改权限")
	}

	errMsg, err := logic.DefaultComment.Modify(echoutils.WrapEchoContext(ctx), cid, content)
	if err != nil {
		return fail(ctx, 4, errMsg)
	}

	return success(ctx, map[string]interface{}{"cid": cid})
}

// CommentList 获取某对象的评论信息
func (CommentController) CommentList(ctx echo.Context) error {
	objid := goutils.MustInt(ctx.QueryParam("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))
	p := goutils.MustInt(ctx.QueryParam("p"))

	commentList, replyComments, pageNum, err := logic.DefaultComment.FindObjectComments(ctx, objid, objtype, p)
	if err != nil {
		return fail(ctx, 1, "服务器内部错误")
	}

	uids := slices.StructsIntSlice(commentList, "Uid")
	if len(replyComments) > 0 {
		replyUids := slices.StructsIntSlice(replyComments, "Uid")
		uids = append(uids, replyUids...)
	}
	users := logic.DefaultUser.FindUserInfos(ctx, uids)

	result := map[string]interface{}{
		"comments": commentList,
		"page_num": pageNum,
	}

	// json encode 不支持 map[int]...
	for uid, user := range users {
		result[strconv.Itoa(uid)] = user
	}

	replyMap := make(map[string]interface{})
	for _, comment := range replyComments {
		replyMap[strconv.Itoa(comment.Floor)] = comment
	}
	result["reply_comments"] = replyMap

	return success(ctx, result)
}

func (self CommentController) TopicDetail(ctx echo.Context) error {
	objid := goutils.MustInt(ctx.Param("objid"))
	cid := goutils.MustInt(ctx.Param("cid"))

	topicMaps := logic.DefaultTopic.FindFullinfoByTids([]int{objid})
	if len(topicMaps) < 1 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	topic := topicMaps[0]
	topic["node"] = logic.GetNode(topic["nid"].(int))

	data := map[string]interface{}{
		"topic": topic,
	}
	data["appends"] = logic.DefaultTopic.FindAppend(ctx, objid)

	err := self.fillCommentAndUser(ctx, data, cid, objid, model.TypeTopic)

	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/topics/"+strconv.Itoa(objid))
	}

	return render(ctx, "topics/comment.html", data)
}

func (self CommentController) ArticleDetail(ctx echo.Context) error {
	objid := goutils.MustInt(ctx.Param("objid"))
	cid := goutils.MustInt(ctx.Param("cid"))

	article, err := logic.DefaultArticle.FindById(ctx, objid)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}
	articleGCTT := logic.DefaultArticle.FindArticleGCTT(ctx, article)

	data := map[string]interface{}{
		"article":      article,
		"article_gctt": articleGCTT,
	}

	err = self.fillCommentAndUser(ctx, data, cid, objid, model.TypeArticle)

	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/articles/"+strconv.Itoa(objid))
	}

	return render(ctx, "articles/comment.html", data)
}

func (CommentController) fillCommentAndUser(ctx echo.Context, data map[string]interface{}, cid, objid, objtype int) error {
	comment, comments := logic.DefaultComment.FindComment(ctx, cid, objid, objtype)

	if comment.Cid == 0 {
		return errors.New("comment not exists!")
	}

	uids := make([]int, 1+len(comments))
	uids[0] = comment.Uid
	for i, comment := range comments {
		uids[i+1] = comment.Uid
	}
	users := logic.DefaultUser.FindUserInfos(ctx, uids)

	data["comment"] = comment
	data["comments"] = comments
	data["users"] = users

	return nil
}
