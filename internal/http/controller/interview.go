// Copyright 2022 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/http/middleware"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeInterview, logic.InterviewComment{})
	logic.RegisterLikeObject(model.TypeInterview, logic.InterviewLike{})
}

type InterviewController struct{}

// RegisterRoute 注册路由
func (self InterviewController) RegisterRoute(g *echo.Group) {
	g.GET("/interview/question", self.TodayQuestion)
	g.GET("/interview/question/:show_sn", self.Find)

	g.Match([]string{"GET", "POST"}, "/interview/new", self.Create, middleware.NeedLogin(), middleware.AdminAuth())
}

func (InterviewController) Create(ctx echo.Context) error {
	question := ctx.FormValue("question")
	// 请求新建面试题页面
	if question == "" || ctx.Request().Method != "POST" {
		interview := &model.InterviewQuestion{}
		return render(ctx, "interview/new.html", map[string]interface{}{"interview": interview})
	}

	forms, _ := ctx.FormParams()
	interview, err := logic.DefaultInterview.Publish(context.EchoContext(ctx), forms)
	if err != nil {
		return fail(ctx, 1, "内部服务错误！")
	}
	return success(ctx, interview)
}

// TodayQuestion 今日题目
func (ic InterviewController) TodayQuestion(ctx echo.Context) error {
	question := logic.DefaultInterview.TodayQuestion(context.EchoContext(ctx))

	data := map[string]interface{}{
		"title": "Go每日一题 今日（" + time.Now().Format("2006-01-02") + "）",
	}
	return ic.detail(ctx, question, data)
}

// Find 某个题目的详情
func (ic InterviewController) Find(ctx echo.Context) error {
	showSn := ctx.Param("show_sn")
	sn, err := strconv.ParseInt(showSn, 32, 64)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/interview/question?"+err.Error())
	}

	question, err := logic.DefaultInterview.FindOne(context.EchoContext(ctx), sn)
	if err != nil || question.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/interview/question")
	}

	data := map[string]interface{}{
		"title": "Go每日一题（" + strconv.Itoa(question.Id) + "）",
	}

	return ic.detail(ctx, question, data)
}

func (InterviewController) detail(ctx echo.Context, question *model.InterviewQuestion, data map[string]interface{}) error {
	data["question"] = question
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		data["likeflag"] = logic.DefaultLike.HadLike(context.EchoContext(ctx), me.Uid, question.Id, model.TypeInterview)
		// data["hadcollect"] = logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), me.Uid, question.Id, model.TypeInterview)

		logic.Views.Incr(Request(ctx), model.TypeInterview, question.Id, me.Uid)

		go logic.DefaultViewRecord.Record(question.Id, model.TypeInterview, me.Uid)

		if me.IsRoot {
			data["view_user_num"] = logic.DefaultViewRecord.FindUserNum(context.EchoContext(ctx), question.Id, model.TypeInterview)
			data["view_source"] = logic.DefaultViewSource.FindOne(context.EchoContext(ctx), question.Id, model.TypeInterview)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeInterview, question.Id)
	}

	return render(ctx, "interview/question.html,common/comment.html", data)
}
