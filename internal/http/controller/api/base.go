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

// normalizeReplies 将 FindObjComments 返回的评论 map（大写 Go 字段名 + 嵌套 user）
// 规整为前端期望的扁平、小写 key 结构：id/uid/floor/content/ctime/name/avatar。
// 旧的 Go 模板直接用大写 key 渲染没问题，但新的 JSON API 暴露给 Next.js 前端时，
// 前端类型是扁平小写 key，因此需要转换，否则 reply.id/uid/floor 全部 undefined。
func normalizeReplies(replies []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(replies))
	for _, reply := range replies {
		item := map[string]interface{}{
			"id":      reply["Cid"],
			"objid":   reply["Objid"],
			"objtype": reply["Objtype"],
			"uid":     reply["Uid"],
			"floor":   reply["Floor"],
			"flag":    reply["Flag"],
			"ctime":   reply["Ctime"],
			"content": reply["content"],
		}
		if user, ok := reply["user"].(*model.User); ok && user != nil {
			item["name"] = user.Username
			item["avatar"] = user.Avatar
			item["username"] = user.Username
		}
		result = append(result, item)
	}
	return result
}

// getAuthToken 读取认证 token：优先读 HttpOnly Cookie，回退到 X-Token header（兼容旧客户端）
func getAuthToken(ctx echo.Context) string {
	if cookie, err := ctx.Cookie(authCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return ctx.Request().Header.Get("X-Token")
}

// isSecure 判断当前请求是否为 HTTPS（根据反向代理 header 或环境变量）
func isSecure(ctx echo.Context) bool {
	if ctx.Request().Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}
	return config.ConfigFile.MustValue("global", "env", "prod") == "prod"
}

// setAuthCookie 设置认证 Cookie（HttpOnly, Secure, SameSite=Strict）
func setAuthCookie(ctx echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = authCookieName
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Secure = isSecure(ctx) // 根据 X-Forwarded-Proto 或环境动态决定

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
	// SameSite 与 setAuthCookie 保持一致
	env := config.ConfigFile.MustValue("global", "env", "prod")
	if env == "prod" {
		cookie.SameSite = http.SameSiteStrictMode
	} else {
		cookie.SameSite = http.SameSiteLaxMode
	}
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

		// Origin 和 Referer 都没有时，拒绝写操作（CSRF 防护）
		// 前端所有写请求（fetch/XHR）都会自动带上 Origin，表单提交也会带上 Referer
		if origin == "" && referer == "" {
			return fail(ctx, "缺少 Origin 或 Referer", 403)
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
			// Referer 存在但不在白名单中，拒绝（CSRF 防护）
			// 之前这里 fallthrough 到 next(ctx) 等于绕过校验
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    1,
				"message": "非法的请求来源",
			})
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

	// 优先从缓存获取用户信息
	userInfo := logic.GetOrFetchUserInfo(uid, func() *logic.UserInfoCache {
		user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
		if user == nil || user.Uid == 0 {
			return nil
		}
		// IsAdmin 必须基于 user_role 表判断（AdminMinRoleId=7），
		// 不能简化为 IsRoot，否则板块管理员/晨读管理员等角色会丢失权限。
		// FindOne 已填充 user.Roleids，UserLogic.IsAdmin 据此判断。
		return &logic.UserInfoCache{
			Uid:      user.Uid,
			Username: user.Username,
			Email:    user.Email,
			IsRoot:   user.IsRoot,
			IsAdmin:  logic.DefaultUser.IsAdmin(user),
			IsVip:    user.IsVip,
			Avatar:   user.Avatar,
			Balance:  user.Balance,
			Status:   user.Status,
		}
	})

	if userInfo == nil {
		return nil, fail(ctx, "用户不存在")
	}

	// 状态校验：与 master 的 NeedLogin 中间件保持一致。
	// 仅 UserStatusAudit（已激活）允许写操作；UserStatusOutage（冻结）等需重新登录。
	// master 在中间件层统一拦截；refactor 把 requireAuth 作为唯一认证入口，故在此检查。
	if userInfo.Status != model.UserStatusAudit {
		return nil, fail(ctx, "账号已被冻结或未激活，请重新登录", NeedReLoginCode)
	}

	return &model.Me{
		Uid:      userInfo.Uid,
		Username: userInfo.Username,
		IsRoot:   userInfo.IsRoot,
		IsAdmin:  userInfo.IsAdmin,
		IsVip:    userInfo.IsVip,
		Balance:  userInfo.Balance,
		Status:   userInfo.Status,
	}, nil
}

// isValidObjType 校验对象类型是否在合法枚举内。
// TypeTopic=0 是个陷阱（不能简单用 objtype==0 判断无效）。
func isValidObjType(objtype int) bool {
	switch objtype {
	case model.TypeTopic, model.TypeArticle, model.TypeResource,
		model.TypeWiki, model.TypeProject, model.TypeBook, model.TypeInterview:
		return true
	}
	return false
}
