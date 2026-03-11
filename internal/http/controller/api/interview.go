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

	echo "github.com/labstack/echo/v4"
)

type InterviewController struct{}

func (self InterviewController) RegisterRoute(g *echo.Group) {
	g.GET("/interviews", self.List)
	g.GET("/interviews/today", self.Today)
	g.GET("/interviews/question/:sn", self.Question)
}

// List 获取面试题列表（分页）
func (InterviewController) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("p"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	level, _ := strconv.Atoi(ctx.QueryParam("level"))
	if level < 0 || level > 2 {
		level = -1
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

// Question 根据 show_sn（base32 字符串）获取面试题详情
func (InterviewController) Question(ctx echo.Context) error {
	showSn := ctx.Param("sn")
	sn, err := strconv.ParseInt(showSn, 32, 64)
	if err != nil {
		return fail(ctx, "无效的面试题 ID")
	}

	question, err := logic.DefaultInterview.FindOne(context.EchoContext(ctx), sn)
	if err != nil || question.Id == 0 {
		return fail(ctx, "面试题不存在")
	}
	return success(ctx, map[string]interface{}{"question": question})
}
