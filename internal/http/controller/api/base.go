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
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
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

// setAuthCookie 设置认证 Cookie（HttpOnly, SameSite=Lax）
func setAuthCookie(ctx echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = authCookieName
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode
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

// requireAuth 验证用户登录，返回 Me 对象
func requireAuth(ctx echo.Context) (*model.Me, error) {
	token := getAuthToken(ctx)
	if token == "" {
		return nil, fail(ctx, "请先登录", NeedReLoginCode)
	}
	if !ValidateToken(token) {
		return nil, fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}
	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return nil, fail(ctx, "无效的 token", NeedReLoginCode)
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return nil, fail(ctx, "用户不存在")
	}

	return &model.Me{
		Uid:      user.Uid,
		Username: user.Username,
		IsRoot:   user.IsRoot,
		IsAdmin:  user.IsRoot, // Root 用户即为管理员
		IsVip:    user.IsVip,
	}, nil
}
