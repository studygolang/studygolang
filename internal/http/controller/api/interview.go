// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"math"
	"strconv"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type InterviewController struct{}

func (self InterviewController) RegisterRoute(g *echo.Group) {
	g.GET("/interviews", self.List)
	g.GET("/interviews/today", self.Today)
	g.GET("/interviews/question/:sn", self.Question)
	g.POST("/interviews", self.Create)
}

// List 获取面试题列表（分页）
func (InterviewController) List(ctx echo.Context) error {
	page := goutils.MustInt(ctx.QueryParam("p"), 1)
	pageSize := 20
	levelStr := ctx.QueryParam("level")
	level := -1 // -1 表示全部，不过滤 level
	if levelStr != "" {
		if parsed, err := strconv.Atoi(levelStr); err == nil && parsed >= 0 && parsed <= 2 {
			level = parsed
		}
	}

	questions, total, err := logic.DefaultInterview.FindAll(context.EchoContext(ctx), page, pageSize, level)
	if err != nil {
		return fail(ctx, "获取面试题列表失败")
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	return success(ctx, map[string]interface{}{
		"questions":   questions,
		"total":       total,
		"page":        page,
		"total_pages": totalPages,
		"has_more":    page < totalPages,
	})
}

// Today 获取今日面试题
func (InterviewController) Today(ctx echo.Context) error {
	question := logic.DefaultInterview.TodayQuestion(context.EchoContext(ctx))
	if question == nil || question.Id == 0 {
		return success(ctx, nil)
	}
	return success(ctx, map[string]interface{}{
		"question": question,
	})
}

// Create 创建面试题（需要登录 + 管理员权限）
// POST /api/v1/interviews
func (InterviewController) Create(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}
	if !me.IsRoot && !me.IsAdmin {
		return fail(ctx, "无权创建面试题")
	}

	forms, _ := ctx.FormParams()
	question := forms.Get("question")
	if question == "" {
		return fail(ctx, "面试题内容不能为空")
	}

	// 敏感词检查：管理员账号也可能被盗用，所有写操作都应过敏感词
	// master 的 /interview/new 路由同样挂了敏感词中间件
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
	}

	result, err := logic.DefaultInterview.Publish(context.EchoContext(ctx), forms)
	if err != nil {
		return fail(ctx, "创建失败："+err.Error())
	}

	return success(ctx, map[string]interface{}{
		"question": result,
	})
}

// Question 根据 show_sn（base32 字符串）获取面试题详情
func (InterviewController) Question(ctx echo.Context) error {
	showSn := ctx.Param("sn")
	sn, err := strconv.ParseInt(showSn, 32, 64)
	if err != nil {
		return fail(ctx, "无效的面试题 ID")
	}

	question, err := logic.DefaultInterview.FindOne(context.EchoContext(ctx), sn)
	if err != nil {
		return fail(ctx, "面试题不存在")
	}
	if question == nil || question.Id == 0 {
		return fail(ctx, "面试题不存在")
	}

	// 浏览计数（与 master interview.detail 一致）：登录用户记录 uid 并异步留痕，匿名用户仅计数
	if me := optionalAuth(ctx); me != nil {
		logic.Views.Incr(ctx.Request(), model.TypeInterview, question.Id, me.Uid)
		go logic.DefaultViewRecord.Record(question.Id, model.TypeInterview, me.Uid)
	} else {
		logic.Views.Incr(ctx.Request(), model.TypeInterview, question.Id)
	}

	return success(ctx, map[string]interface{}{"question": question})
}
