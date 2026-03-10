// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
)

type InterviewController struct{}

func (self InterviewController) RegisterRoute(g *echo.Group) {
	g.GET("/interviews/today", self.Today)
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
