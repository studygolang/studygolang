// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"db"
	"logic"
	"model"
	"net/http"
	"net/url"
	"strings"
	"util"

	. "http"

	"github.com/gorilla/context"
	"github.com/labstack/echo"
)

// AutoLogin 用于 echo 框架的自动登录和通过 cookie 获取用户信息
func AutoLogin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// github.com/gorilla/sessions 要求必须 Clear
			defer context.Clear(Request(ctx))

			session := GetCookieSession(ctx)
			username, ok := session.Values["username"]
			if ok {
				if db.MasterDB != nil {
					// TODO: 考虑缓存，或延迟查询，避免每次都查询
					user := logic.DefaultUser.FindCurrentUser(ctx, username)
					if user.Uid != 0 {
						ctx.Set("user", user)
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

// NeedLogin 用于 echo 框架的验证必须登录的请求
func NeedLogin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user, ok := ctx.Get("user").(*model.Me)
			if !ok || user.Status != model.UserStatusAudit {
				method := ctx.Request().Method()
				if util.IsAjax(ctx) {
					return ctx.JSON(http.StatusForbidden, `{"ok":0,"error":"403 Forbidden"}`)
				} else {
					if method == "POST" {
						return ctx.HTML(http.StatusForbidden, `403 Forbidden`)
					}

					if !ok {
						reqURL := ctx.Request().URL()
						uri := reqURL.Path()
						if reqURL.QueryString() != "" {
							uri += "?" + reqURL.QueryString()
						}
						return ctx.Redirect(http.StatusSeeOther, "/account/login?redirect_uri="+url.QueryEscape(uri))
					} else {
						// 未激活可以查看账号信息
						if !strings.HasPrefix(ctx.Path(), "/account") {
							return echo.NewHTTPError(http.StatusForbidden, `您的邮箱未激活，<a href="/account/edit">去激活</a>`)
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
