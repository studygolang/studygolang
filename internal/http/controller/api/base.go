// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

// Package api 提供面向 Next.js 前端的 REST API（Web 端）
package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

const (
	perPage        = 20            // 每页默认条数
	authCookieName = "sg_token"    // HttpOnly 认证 Cookie 名
	cookieMaxAge   = 7 * 24 * 3600 // Cookie 有效期 7 天
)

// getAuthToken 读取认证 token：优先读 HttpOnly Cookie，回退到 X-Token header（兼容旧客户端）
func getAuthToken(ctx echo.Context) string {
	if cookie, err := ctx.Cookie(authCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return ctx.Request().Header.Get("X-Token")
}

// setAuthCookie 设置认证 Cookie（HttpOnly, Secure, SameSite=Strict）
func setAuthCookie(ctx echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = authCookieName
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Secure = true // 强制 HTTPS

	// 生产环境使用 Strict，开发环境可以使用 Lax
	env := config.ConfigFile.MustValue("global", "env", "prod")
	if env == "prod" {
		cookie.SameSite = http.SameSiteStrictMode
	} else {
		cookie.SameSite = http.SameSiteLaxMode // 开发环境允许部分跨站请求
	}

	cookie.Path = "/"
	cookie.MaxAge = cookieMaxAge
	ctx.SetCookie(cookie)
}

// clearAuthCookie 清除认证 Cookie
func clearAuthCookie(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = authCookieName
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.MaxAge = -1
	ctx.SetCookie(cookie)
}

// originCheck 检查写操作请求来源是否合法（CSRF 防护）
// GET/HEAD/OPTIONS 请求直接放行；POST/PUT/DELETE 等写操作检查 Origin 或 Referer
func originCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		method := ctx.Request().Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return next(ctx)
		}

		origin := ctx.Request().Header.Get("Origin")
		referer := ctx.Request().Header.Get("Referer")

		// Origin 和 Referer 都没有，放行（兼容直接调用 API 的客户端）
		if origin == "" && referer == "" {
			return next(ctx)
		}

		allowed := getAllowedOrigins()

		// 优先检查 Origin
		if origin != "" {
			for _, a := range allowed {
				if origin == a {
					return next(ctx)
				}
			}
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    1,
				"message": "非法的请求来源",
			})
		}

		// Referer 检查
		if referer != "" {
			if refURL, err := url.Parse(referer); err == nil {
				refOrigin := refURL.Scheme + "://" + refURL.Host
				for _, a := range allowed {
					if refOrigin == a {
						return next(ctx)
					}
				}
			}
		}

		return next(ctx)
	}
}

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

// parseAuthUID 从请求中获取并验证 token，返回 uid。
// 支持 JWT 和旧版 MD5 Token 双轨验证。
func parseAuthUID(ctx echo.Context) (int, error) {
	token := getAuthToken(ctx)
	if token == "" {
		return 0, fail(ctx, "请先登录", NeedReLoginCode)
	}

	uid, _, valid := ValidateTokenAuto(token)
	if !valid {
		return 0, fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}
	if uid == 0 {
		return 0, fail(ctx, "无效的 token", NeedReLoginCode)
	}

	return uid, nil
}
func requireAuth(ctx echo.Context) (*model.Me, error) {
	token := getAuthToken(ctx)
	if token == "" {
		return nil, fail(ctx, "请先登录", NeedReLoginCode)
	}

	// 使用统一验证方法，自动识别 JWT 和旧 MD5 Token
	uid, _, valid := ValidateTokenAuto(token)
	if !valid {
		return nil, fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	// 记录 Token 类型（用于监控迁移进度）
	if strings.HasPrefix(token, "eyJ") {
		// JWT Token
	} else {
		// Legacy MD5 Token
	}

	// 优先从缓存获取用户信息
	userInfo := logic.GetOrFetchUserInfo(uid, func() *logic.UserInfoCache {
		user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
		if user == nil || user.Uid == 0 {
			return nil
		}
		return &logic.UserInfoCache{
			Uid:      user.Uid,
			Username: user.Username,
			Email:    user.Email,
			IsRoot:   user.IsRoot,
			IsVip:    user.IsVip,
			Avatar:   user.Avatar,
		}
	})

	if userInfo == nil {
		return nil, fail(ctx, "用户不存在")
	}

	return &model.Me{
		Uid:      userInfo.Uid,
		Username: userInfo.Username,
		IsRoot:   userInfo.IsRoot,
		IsAdmin:  userInfo.IsRoot, // Root 用户即为管理员
		IsVip:    userInfo.IsVip,
	}, nil
}
