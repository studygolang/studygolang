// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
)

const oauthStateSessionKey = "oauth_state"
const oauthRedirectSessionKey = "oauth_redirect"

type OAuthController struct{}

func (self OAuthController) RegisterRoute(g *echo.Group) {
	g.GET("/oauth/github/url", self.GithubURL)
	g.GET("/oauth/github/callback", self.GithubCallbackRedirect)
	g.GET("/oauth/gitea/url", self.GiteaURL)
	g.GET("/oauth/gitea/callback", self.GiteaCallbackRedirect)
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

// saveOAuthRedirect 将用户原始跳转目标保存到 session
// 仅接受相对路径（必须以 "/" 开头），拒绝绝对 URL 与协议相对 URL（"//evil.com"）
// 以避免 open redirect 漏洞
func saveOAuthRedirect(ctx echo.Context, redirect string) {
	if !isSafeRedirect(redirect) {
		redirect = "/"
	}
	session := GetCookieSession(ctx)
	session.Values[oauthRedirectSessionKey] = redirect
	if err := session.Save(Request(ctx), ResponseWriter(ctx)); err != nil {
		// 记录失败原因，便于排查 OAuth 跳转目标丢失的问题
		// （cookie secret 改变、Redis 故障、cookie 超长等）
		getLogger(ctx).Errorln("saveOAuthRedirect session.Save failed:", err)
	}
}

// isSafeRedirect 校验 redirect 字符串是否为安全的相对路径
func isSafeRedirect(redirect string) bool {
	if redirect == "" {
		return false
	}
	// 必须以 "/" 开头（相对路径）
	if !strings.HasPrefix(redirect, "/") {
		return false
	}
	// 拒绝协议相对 URL（"//evil.com" → 浏览器会按 https://evil.com 解析）
	if strings.HasPrefix(redirect, "//") {
		return false
	}
	// 拒绝含反斜杠的 URL（浏览器会将 "\" 解析为 "/"，"/\evil.com" → "//evil.com"）
	if strings.Contains(redirect, "\\") {
		return false
	}
	// 拒绝以 "/\" 开头（Windows 风格路径分隔符，避免协议相对 URL 绕过）
	if strings.HasPrefix(redirect, "/\\") {
		return false
	}
	return true
}

// getOAuthRedirect 从 session 取出跳转目标并清除
func getOAuthRedirect(ctx echo.Context) string {
	session := GetCookieSession(ctx)
	val, _ := session.Values[oauthRedirectSessionKey]
	delete(session.Values, oauthRedirectSessionKey)
	if err := session.Save(Request(ctx), ResponseWriter(ctx)); err != nil {
		getLogger(ctx).Errorln("getOAuthRedirect session.Save failed:", err)
	}
	if s, ok := val.(string); ok {
		return s
	}
	return "/"
}

// oauthCallbackURL 构造 OAuth 提供商回调地址（指向后端 /api/v1/oauth/:provider/callback）
func oauthCallbackURL(ctx echo.Context, provider string) string {
	scheme := "http"
	host := ctx.Request().Host
	if ctx.Request().Header.Get("X-Forwarded-Proto") == "https" || ctx.Request().TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + host + "/api/v1/oauth/" + provider + "/callback"
}

// appendOAuthFlag 在 redirect 上拼接 oauth 状态标记。
// 当 redirect 已带 query（例如 /topics?tab=hot）时，必须用 & 分隔，
// 否则 /topics?tab=hot?oauth=bind_success 会被解析为 tab=hot?oauth=bind_success，
// 前端无法读到 oauth 标志。
func appendOAuthFlag(redirect, flag string) string {
	sep := "?"
	if strings.Contains(redirect, "?") {
		sep = "&"
	}
	return redirect + sep + flag
}

// GithubURL 返回 GitHub OAuth 授权 URL（前端跳转用）
func (OAuthController) GithubURL(ctx echo.Context) error {
	// uri 是用户登录后要跳转的前端页面，保存到 session
	redirect := ctx.QueryParam("uri")
	if redirect == "" {
		redirect = "/"
	}
	saveOAuthRedirect(ctx, redirect)

	state := generateOAuthState(ctx)
	if state == "" {
		return fail(ctx, "生成授权状态失败，请重试")
	}
	// 回调地址指向后端，而非前端页面
	callbackURL := oauthCallbackURL(ctx, "github")
	authURL := logic.DefaultThirdUser.GithubAuthCodeUrl(context.EchoContext(ctx), callbackURL, state)
	return success(ctx, map[string]interface{}{
		"url": authURL,
	})
}

// GithubCallbackRedirect GitHub OAuth 回调（由 OAuth 提供商直接重定向到后端）
// 处理完后重定向回前端页面
func (OAuthController) GithubCallbackRedirect(ctx echo.Context) error {
	state := ctx.QueryParam("state")
	if !validateOAuthState(ctx, state) {
		return ctx.Redirect(http.StatusSeeOther, "/account/login?error=oauth_invalid")
	}

	code := ctx.QueryParam("code")
	if code == "" {
		return ctx.Redirect(http.StatusSeeOther, "/account/login?error=oauth_failed")
	}

	// 检查是否已登录（绑定场景）
	// 必须用 optionalAuth —— requireAuth 失败时会 fail() 写错误响应，
	// 之后走到 ctx.Redirect() 会触发 "HTTP response already written" 双写问题，
	// 且 token 过期的登录用户会看到 JSON 401 而不是绑定成功跳转。
	me := optionalAuth(ctx)
	if me != nil {
		if bindErr := logic.DefaultThirdUser.BindGithub(context.EchoContext(ctx), code, me); bindErr != nil {
			getLogger(ctx).Errorln("OAuth GitHub bind failed:", bindErr)
			return ctx.Redirect(http.StatusSeeOther, "/account/login?error=bind_failed")
		}
		redirect := getOAuthRedirect(ctx)
		return ctx.Redirect(http.StatusSeeOther, appendOAuthFlag(redirect, "oauth=bind_success"))
	}

	// 未登录用户走登录流程
	user, loginErr := logic.DefaultThirdUser.LoginFromGithub(context.EchoContext(ctx), code)
	if loginErr != nil || user.Uid == 0 {
		getLogger(ctx).Errorln("OAuth GitHub login failed:", loginErr)
		return ctx.Redirect(http.StatusSeeOther, "/account/login?error=oauth_login_failed")
	}

	SetLoginCookie(ctx, user.Username)
	if jwtToken, err := GenJWTToken(user.Uid, user.Username); err == nil {
		setAuthCookie(ctx, jwtToken)
	}

	// 新 OAuth 用户余额为 0，引导去领取初始资本（master 同此逻辑）
	if user.Balance == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/balance")
	}

	redirect := getOAuthRedirect(ctx)
	return ctx.Redirect(http.StatusSeeOther, redirect)
}

// GiteaURL 返回 Gitea OAuth 授权 URL（前端跳转用）
func (OAuthController) GiteaURL(ctx echo.Context) error {
	redirect := ctx.QueryParam("uri")
	if redirect == "" {
		redirect = "/"
	}
	saveOAuthRedirect(ctx, redirect)

	state := generateOAuthState(ctx)
	if state == "" {
		return fail(ctx, "生成授权状态失败，请重试")
	}
	callbackURL := oauthCallbackURL(ctx, "gitea")
	authURL := logic.DefaultThirdUser.GiteaAuthCodeUrl(context.EchoContext(ctx), callbackURL, state)
	return success(ctx, map[string]interface{}{
		"url": authURL,
	})
}

// GiteaCallbackRedirect Gitea OAuth 回调（由 OAuth 提供商直接重定向到后端）
func (OAuthController) GiteaCallbackRedirect(ctx echo.Context) error {
	state := ctx.QueryParam("state")
	if !validateOAuthState(ctx, state) {
		return ctx.Redirect(http.StatusSeeOther, "/account/login?error=oauth_invalid")
	}

	code := ctx.QueryParam("code")
	if code == "" {
		return ctx.Redirect(http.StatusSeeOther, "/account/login?error=oauth_failed")
	}

	// 同 GithubCallbackRedirect：必须用 optionalAuth 避免双写
	me := optionalAuth(ctx)
	if me != nil {
		if bindErr := logic.DefaultThirdUser.BindGitea(context.EchoContext(ctx), code, me); bindErr != nil {
			getLogger(ctx).Errorln("OAuth Gitea bind failed:", bindErr)
			return ctx.Redirect(http.StatusSeeOther, "/account/login?error=bind_failed")
		}
		redirect := getOAuthRedirect(ctx)
		return ctx.Redirect(http.StatusSeeOther, appendOAuthFlag(redirect, "oauth=bind_success"))
	}

	user, loginErr := logic.DefaultThirdUser.LoginFromGitea(context.EchoContext(ctx), code)
	if loginErr != nil || user.Uid == 0 {
		getLogger(ctx).Errorln("OAuth Gitea login failed:", loginErr)
		return ctx.Redirect(http.StatusSeeOther, "/account/login?error=oauth_login_failed")
	}

	SetLoginCookie(ctx, user.Username)
	if jwtToken, err := GenJWTToken(user.Uid, user.Username); err == nil {
		setAuthCookie(ctx, jwtToken)
	}

	// 新 OAuth 用户余额为 0，引导去领取初始资本（master 同此逻辑）
	if user.Balance == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/balance")
	}

	redirect := getOAuthRedirect(ctx)
	return ctx.Redirect(http.StatusSeeOther, redirect)
}
