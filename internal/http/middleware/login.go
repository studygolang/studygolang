// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mycontext "github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/db"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"
	"github.com/studygolang/studygolang/util"

	"github.com/gorilla/context"
	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

// AutoLogin 用于 echo 框架的自动登录和通过 cookie 获取用户信息
func AutoLogin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// github.com/gorilla/sessions 要求必须 Clear
			defer context.Clear(Request(ctx))

			ctx.Set("req_start_time", time.Now())

			var getCurrentUser = func(usernameOrId interface{}) {
				ip := goutils.RemoteIp(Request(ctx))
				// IP 黑名单，不让登录
				if logic.DefaultRisk.IsBlackIP(ip) {
					return
				}

				if db.MasterDB != nil {
					ctx.Set("ip", ip)
					// TODO: 考虑缓存，或延迟查询，避免每次都查询
					user := logic.DefaultUser.FindCurrentUser(mycontext.EchoContext(ctx), usernameOrId)
					if user.Uid != 0 {
						ctx.Set("user", user)

						if !util.IsAjax(ctx) && ctx.Path() != "/ws" {
							go logic.ViewObservable.NotifyObservers(user.Uid, 0, 0)
						}
					}
				}
			}

			session := GetCookieSession(ctx)
			username, ok := session.Values["username"]
			if ok {
				getCurrentUser(username)
			} else {
				// App（手机） 登录
				uid, ok := ParseToken(ctx.FormValue("token"))
				if ok {
					getCurrentUser(uid)
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
				method := ctx.Request().Method
				if util.IsAjax(ctx) {
					if !strings.HasPrefix(ctx.Path(), "/account") {
						return ctx.JSON(http.StatusForbidden, map[string]interface{}{"ok": 0, "error": "403 Forbidden"})
					}
				} else {
					if method == "POST" {
						return ctx.HTML(http.StatusForbidden, `403 Forbidden`)
					}

					if !ok {
						reqURL := ctx.Request().URL
						uri := reqURL.Path
						if reqURL.RawQuery != "" {
							uri += "?" + reqURL.RawQuery
						}
						return ctx.Redirect(http.StatusSeeOther, "/account/login?redirect_uri="+url.QueryEscape(uri))
					} else {
						// 未激活可以查看账号信息
						if !strings.HasPrefix(ctx.Path(), "/account") {
							return echo.NewHTTPError(http.StatusForbidden, `您的邮箱未激活，<a href="/account/edit">去激活</a>`)
						}
					}
				}
			} else {
				newUserWait := time.Duration(logic.UserSetting[model.KeyNewUserWait]) * time.Second
				if newUserWait > 0 {
					elapse := time.Now().Sub(user.CreatedAt)
					if elapse <= newUserWait {
						return echo.NewHTTPError(http.StatusForbidden, `您需要再等待`+(newUserWait-elapse).String()+"才能进行此操作")
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

// AppNeedLogin 用于 echo 框架的验证必须登录的请求（APP 专用）
func AppNeedLogin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user, ok := ctx.Get("user").(*model.Me)
			if ok {
				// 校验 token 是否有效
				if !ValidateToken(ctx.QueryParam("token")) {
					return outputAppJSON(ctx, NeedReLoginCode, "token无效，请重新登录！")
				}

				if user.Status != model.UserStatusAudit {
					return outputAppJSON(ctx, 1, "账号未审核通过、被冻结或被停号，请联系我们")
				}
			} else {
				return outputAppJSON(ctx, NeedReLoginCode, "请先登录！")
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}

func outputAppJSON(ctx echo.Context, code int, msg string) error {
	AccessControl(ctx)
	return ctx.JSON(http.StatusForbidden, map[string]interface{}{"code": strconv.Itoa(code), "msg": msg})
}
