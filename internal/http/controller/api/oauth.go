// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gorilla/sessions"
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
)

const oauthStateSessionKey = "oauth_state"

type OAuthController struct{}

func (self OAuthController) RegisterRoute(g *echo.Group) {
	g.GET("/oauth/github/url", self.GithubURL)
	g.GET("/oauth/github/callback", self.GithubCallback)
	g.GET("/oauth/gitea/url", self.GiteaURL)
	g.GET("/oauth/gitea/callback", self.GiteaCallback)
}

// generateOAuthState 生成随机 state 并存入 session，返回 state 值
func generateOAuthState(ctx echo.Context) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		getLogger(ctx).Errorln("generateOAuthState rand.Read failed:", err)
		return ""
	}
	state := hex.EncodeToString(b)

	session := GetCookieSession(ctx)
	session.Values[oauthStateSessionKey] = state
	if err := session.Save(Request(ctx), ResponseWriter(ctx)); err != nil {
		getLogger(ctx).Errorln("generateOAuthState session.Save failed:", err)
		return ""
	}

	return state
}

// validateOAuthState 验证回调中的 state 是否与 session 中存储的一致
func validateOAuthState(ctx echo.Context, state string) bool {
	if state == "" {
		return false
	}

	session := GetCookieSession(ctx)
	expected, ok := session.Values[oauthStateSessionKey].(string)
	if !ok || expected == "" || expected != state {
		return false
	}

	// 验证通过后清除 state（一次性使用）
	delete(session.Values, oauthStateSessionKey)
	session.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}
	if err := session.Save(Request(ctx), ResponseWriter(ctx)); err != nil {
		getLogger(ctx).Errorln("validateOAuthState session.Save failed:", err)
	}

	return true
}

// GithubURL 返回 GitHub OAuth 授权 URL（前端跳转用）
// 生成随机 state 参数防止 CSRF 攻击
func (OAuthController) GithubURL(ctx echo.Context) error {
	uri := ctx.QueryParam("uri")
	state := generateOAuthState(ctx)
	if state == "" {
		return fail(ctx, "生成授权状态失败，请重试")
	}
	url := logic.DefaultThirdUser.GithubAuthCodeUrl(context.EchoContext(ctx), uri, state)
	return success(ctx, map[string]interface{}{
		"url": url,
	})
}

// GithubCallback GitHub OAuth 回调（登录或绑定）
// 验证 state 参数防止 CSRF 攻击
func (OAuthController) GithubCallback(ctx echo.Context) error {
	// 验证 state 参数
	state := ctx.QueryParam("state")
	if !validateOAuthState(ctx, state) {
		return fail(ctx, "无效的授权请求，请重试")
	}

	// OAuth 标准回调使用 Query 参数传递 code
	code := ctx.QueryParam("code")
	if code == "" {
		return fail(ctx, "授权码不能为空")
	}

	// 检查是否已登录（绑定场景）
	me, err := requireAuth(ctx)
	if err == nil && me != nil {
		if bindErr := logic.DefaultThirdUser.BindGithub(context.EchoContext(ctx), code, me); bindErr != nil {
			getLogger(ctx).Errorln("OAuth GitHub bind failed:", bindErr)
			return fail(ctx, "绑定失败，请稍后重试")
		}
		return success(ctx, map[string]interface{}{
			"action":  "bind",
			"message": "GitHub 账号绑定成功",
		})
	}

	// 未登录用户走登录流程
	user, loginErr := logic.DefaultThirdUser.LoginFromGithub(context.EchoContext(ctx), code)
	if loginErr != nil || user.Uid == 0 {
		if loginErr != nil {
			getLogger(ctx).Errorln("OAuth GitHub login failed:", loginErr)
		}
		return fail(ctx, "登录失败，请稍后重试")
	}

	// 登录成功，种 cookie
	SetLoginCookie(ctx, user.Username)

	// 设置 JWT auth cookie（Next.js 前端使用）
	if jwtToken, err := GenJWTToken(user.Uid, user.Username); err == nil {
		setAuthCookie(ctx, jwtToken)
	}

	return success(ctx, map[string]interface{}{
		"action":   "login",
		"username": user.Username,
		"balance":  user.Balance,
	})
}

// GiteaURL 返回 Gitea OAuth 授权 URL（前端跳转用）
// 生成随机 state 参数防止 CSRF 攻击
func (OAuthController) GiteaURL(ctx echo.Context) error {
	uri := ctx.QueryParam("uri")
	state := generateOAuthState(ctx)
	if state == "" {
		return fail(ctx, "生成授权状态失败，请重试")
	}
	url := logic.DefaultThirdUser.GiteaAuthCodeUrl(context.EchoContext(ctx), uri, state)
	return success(ctx, map[string]interface{}{
		"url": url,
	})
}

// GiteaCallback Gitea OAuth 回调（登录或绑定）
// 验证 state 参数防止 CSRF 攻击
func (OAuthController) GiteaCallback(ctx echo.Context) error {
	// 验证 state 参数
	state := ctx.QueryParam("state")
	if !validateOAuthState(ctx, state) {
		return fail(ctx, "无效的授权请求，请重试")
	}

	code := ctx.QueryParam("code")
	if code == "" {
		return fail(ctx, "授权码不能为空")
	}

	me, err := requireAuth(ctx)
	if err == nil && me != nil {
		if bindErr := logic.DefaultThirdUser.BindGitea(context.EchoContext(ctx), code, me); bindErr != nil {
			getLogger(ctx).Errorln("OAuth Gitea bind failed:", bindErr)
			return fail(ctx, "绑定失败，请稍后重试")
		}
		return success(ctx, map[string]interface{}{
			"action":  "bind",
			"message": "Gitea 账号绑定成功",
		})
	}

	user, loginErr := logic.DefaultThirdUser.LoginFromGitea(context.EchoContext(ctx), code)
	if loginErr != nil || user.Uid == 0 {
		if loginErr != nil {
			getLogger(ctx).Errorln("OAuth Gitea login failed:", loginErr)
		}
		return fail(ctx, "登录失败，请稍后重试")
	}

	SetLoginCookie(ctx, user.Username)

	// 设置 JWT auth cookie（Next.js 前端使用）
	if jwtToken, err := GenJWTToken(user.Uid, user.Username); err == nil {
		setAuthCookie(ctx, jwtToken)
	}

	return success(ctx, map[string]interface{}{
		"action":   "login",
		"username": user.Username,
		"balance":  user.Balance,
	})
}
