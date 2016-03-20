// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package middleware

import (
	"logic"
	"model"
	"net/http"
	"net/url"
	"util"

	. "http"

	"github.com/labstack/echo"
)

// AutoLogin 用于 echo 框架的自动登录和通过 cookie 获取用户信息
func AutoLogin() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			session := GetCookieSession(ctx)
			username, ok := session.Values["username"]
			if ok {
				user := logic.DefaultUser.FindOne(ctx, "username", username)
				if user.Uid != 0 {
					ctx.Set("user", user)
				}
			}

			if err := h(ctx); err != nil {
				ctx.Error(err)
			}

			return nil
		}
	}
}

// NeedLogin 用于 echo 框架的验证必须登录的请求
func NeedLogin() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			_, ok := ctx.Get("user").(*model.User)
			if !ok {
				req := ctx.Request()
				method := req.Method
				if util.IsAjax(req) {
					return ctx.JSON(http.StatusForbidden, `{"ok":0,"error":"403 Forbidden"}`)
				} else {
					if method == "POST" {
						return ctx.HTML(http.StatusForbidden, `403 Forbidden`)
					}
					reqURL := ctx.Request().URL
					uri := reqURL.Path
					if reqURL.RawQuery != "" {
						uri += "?" + reqURL.RawQuery
					}
					return ctx.Redirect(http.StatusSeeOther, "/account/login?redirect_uri="+url.QueryEscape(uri))
				}
			}

			if err := h(ctx); err != nil {
				ctx.Error(err)
			}

			return nil
		}
	}
}
