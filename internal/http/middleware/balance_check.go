// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"net/http"

	"github.com/studygolang/studygolang/internal/model"
	"github.com/studygolang/studygolang/util"

	echo "github.com/labstack/echo/v4"
)

// BalanceCheck 用于 echo 框架，用户发布内容校验余额是否足够
func BalanceCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			if util.IsAjax(ctx) {

				curUser := ctx.Get("user").(*model.Me)

				title := ctx.FormValue("title")
				content := ctx.FormValue("content")
				if ctx.Request().Method == "POST" && (title != "" || content != "") {
					if ctx.Path() == "/comment/:objid" {
						if curUser.Balance < 5 {
							return ctx.String(http.StatusOK, `{"ok":0,"error":"对不起，您的账号余额不足，可以领取初始资本！"}`)
						}
					} else {
						if curUser.Balance < 20 {
							return ctx.String(http.StatusOK, `{"ok":0,"error":"对不起，您的账号余额不足，可以领取初始资本！"}`)
						}
					}
				}
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}
