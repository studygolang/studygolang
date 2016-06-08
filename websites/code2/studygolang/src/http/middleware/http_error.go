// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"net/http"
	"util"

	. "http"

	"github.com/labstack/echo"
)

// EchoLogger 用于 echo 框架的日志中间件
func HTTPError() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if err := next(ctx); err != nil {

				if !ctx.Response().Committed() {
					if he, ok := err.(*echo.HTTPError); ok {
						switch he.Code {
						case http.StatusNotFound:
							if util.IsAjax(ctx) {
								return ctx.String(http.StatusOK, `{"ok":0,"error":"接口不存在"}`)
							}
							return Render(ctx, "404.html", nil)
						case http.StatusForbidden:
							if util.IsAjax(ctx) {
								return ctx.String(http.StatusOK, `{"ok":0,"error":"没有权限访问"}`)
							}
							return Render(ctx, "403.html", map[string]interface{}{"msg": he.Message})
						case http.StatusInternalServerError:
							if util.IsAjax(ctx) {
								return ctx.String(http.StatusOK, `{"ok":0,"error":"接口服务器错误"}`)
							}
							return Render(ctx, "500.html", nil)
						}
					}
				}
			}
			return nil
		}
	}
}
