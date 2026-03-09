// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

// Package api 提供面向 Next.js 前端的 REST API（Web 端）
package api

import (
	"net/http"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/logger"
)

const perPage = 20 // 每页默认条数

func getLogger(ctx echo.Context) *logger.Logger {
	return logic.GetLogger(context.EchoContext(ctx))
}

// success 返回成功响应
func success(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}

// fail 返回错误响应
func fail(ctx echo.Context, msg string, codes ...int) error {
	code := 1
	if len(codes) > 0 {
		code = codes[0]
	}

	getLogger(ctx).Errorln("api fail:", msg)

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
	})
}
