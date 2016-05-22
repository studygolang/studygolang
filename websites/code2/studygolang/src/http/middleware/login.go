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
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(ctx echo.Context) error {
			session := GetCookieSession(ctx)
			username, ok := session.Values["username"]
			if ok {
				// TODO: 考虑缓存，或延迟查询，避免每次都查询
				user := logic.DefaultUser.FindCurrentUser(ctx, username)
				if user.Uid != 0 {
					ctx.Set("user", user)
				}
			}

			if err := next.Handle(ctx); err != nil {
				return err
			}

			return nil
		})
	}
}

// NeedLogin 用于 echo 框架的验证必须登录的请求
func NeedLogin() echo.MiddlewareFunc {
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(ctx echo.Context) error {
			_, ok := ctx.Get("user").(*model.Me)
			if !ok {
				req := Request(ctx)
				method := req.Method
				if util.IsAjax(ctx) {
					return ctx.JSON(http.StatusForbidden, `{"ok":0,"error":"403 Forbidden"}`)
				} else {
					if method == "POST" {
						return ctx.HTML(http.StatusForbidden, `403 Forbidden`)
					}
					reqURL := req.URL
					uri := reqURL.Path
					if reqURL.RawQuery != "" {
						uri += "?" + reqURL.RawQuery
					}
					return ctx.Redirect(http.StatusSeeOther, "/account/login?redirect_uri="+url.QueryEscape(uri))
				}
			}

			if err := next.Handle(ctx); err != nil {
				return err
			}

			return nil
		})
	}
}
