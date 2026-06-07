// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/url"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gorilla/sessions"
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	guuid "github.com/twinj/uuid"
)

type AuthController struct{}

// loginRequest 登录请求体（支持 JSON）
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Passwd   string `json:"passwd"` // 兼容两种字段名
}

// registerRequest 注册请求体（支持 JSON）
type registerRequest struct {
	Username        string `json:"username" form:"username"`
	Email           string `json:"email" form:"email"`
	Passwd          string `json:"passwd" form:"passwd"`
	Password        string `json:"password" form:"password"`        // 兼容两种字段名
	CaptchaID       string `json:"captcha_id" form:"captcha_id"`   // 验证码 ID（可选）
	CaptchaSolution string `json:"captcha_solution" form:"captcha_solution"` // 验证码答案（可选）
}

// forgotPasswordRequest 忘记密码请求体
type forgotPasswordRequest struct {
	Email string `json:"email"`
}

// resetPasswordRequest 重置密码请求体
type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (self AuthController) RegisterRoute(g *echo.Group) {
	g.POST("/user/login", self.Login)
	g.POST("/user/logout", self.Logout)
	g.POST("/user/register", self.Register)
	g.POST("/user/forgot-password", self.ForgotPassword)
	g.POST("/user/reset-password", self.ResetPassword)
}

// Login 用户登录，返回 token
func (AuthController) Login(ctx echo.Context) error {
	// 速率限制：每 IP 每分钟最多 10 次登录尝试
	if !loginLimiter.check(ctx.RealIP(), 10, time.Minute) {
		return fail(ctx, "登录尝试过于频繁，请稍后再试")
	}

	var req loginRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	username := req.Username
	if username == "" {
		return fail(ctx, "用户名不能为空")
	}

	passwd := req.Passwd
	if passwd == "" {
		passwd = req.Password
	}
	if passwd == "" {
		return fail(ctx, "密码不能为空")
	}

	userLogin, err := logic.DefaultUser.Login(context.EchoContext(ctx), username, passwd)
	if err != nil {
		return fail(ctx, err.Error())
	}

	// 使用新的 JWT Token
	token, err := GenJWTToken(userLogin.Uid, userLogin.Username)
	if err != nil {
		getLogger(ctx).Errorln("failed to generate JWT token:", err)
		// 回退到旧 Token
		token = GenToken(userLogin.Uid)
	}
	// 设置 HttpOnly Cookie，防止 XSS 窃取 token
	setAuthCookie(ctx, token)
	// 同时设置旧的 session（用于 /admin 等旧路由）
	SetLoginCookie(ctx, userLogin.Username)

	return success(ctx, map[string]interface{}{
		"uid":      userLogin.Uid,
		"username": userLogin.Username,
	})
}

// Logout 退出登录，清除认证 Cookie 和 session
func (AuthController) Logout(ctx echo.Context) error {
	clearAuthCookie(ctx)
	// 同时清除旧的 session（用于 /admin 等旧路由）
	session := GetCookieSession(ctx)
	session.Options = &sessions.Options{Path: "/", MaxAge: -1}
	session.Save(Request(ctx), ResponseWriter(ctx))
	return success(ctx, nil)
}

// Register 用户注册
func (AuthController) Register(ctx echo.Context) error {
	// 速率限制：每 IP 每小时最多 5 次注册尝试
	if !registerLimiter.check(ctx.RealIP(), 5, time.Hour) {
		return fail(ctx, "注册尝试过于频繁，请稍后再试")
	}

	var req registerRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	// 验证码校验（强制要求，防止绕过）
	if req.CaptchaID == "" || req.CaptchaSolution == "" {
		return fail(ctx, "请完成验证码")
	}
	if !captcha.VerifyString(req.CaptchaID, req.CaptchaSolution) {
		return fail(ctx, "验证码错误")
	}

	username := req.Username
	if username == "" {
		return fail(ctx, "用户名不能为空")
	}

	email := req.Email
	if email == "" {
		return fail(ctx, "邮箱不能为空")
	}

	passwd := req.Passwd
	if passwd == "" {
		passwd = req.Password
	}
	if passwd == "" {
		return fail(ctx, "密码不能为空")
	}

	form := url.Values{}
	form.Set("username", username)
	form.Set("email", email)
	form.Set("passwd", passwd)

	errMsg, err := logic.DefaultUser.CreateUser(context.EchoContext(ctx), form)
	if err != nil {
		if errMsg == "" {
			errMsg = "注册失败，请稍后重试"
		}
		return fail(ctx, errMsg)
	}

	// 注册成功后自动登录，设置 Cookie
	userLogin, err := logic.DefaultUser.Login(context.EchoContext(ctx), username, passwd)
	if err != nil {
		return success(ctx, map[string]interface{}{
			"username": username,
		})
	}

	// 使用新的 JWT Token
	token, err := GenJWTToken(userLogin.Uid, userLogin.Username)
	if err != nil {
		getLogger(ctx).Errorln("failed to generate JWT token:", err)
		// 回退到旧 Token
		token = GenToken(userLogin.Uid)
	}
	setAuthCookie(ctx, token)
	// 同时设置旧的 session（用于 /admin 等旧路由）
	SetLoginCookie(ctx, userLogin.Username)

	return success(ctx, map[string]interface{}{
		"uid":      userLogin.Uid,
		"username": userLogin.Username,
	})
}

// ForgotPassword 发送重置密码邮件
func (AuthController) ForgotPassword(ctx echo.Context) error {
	var req forgotPasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "参数错误")
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		return fail(ctx, "邮箱不能为空")
	}

	// 检查邮箱是否注册
	if !logic.DefaultUser.UserExists(context.EchoContext(ctx), "email", email) {
		// 为防止邮箱枚举攻击，返回相同提示
		return success(ctx, map[string]interface{}{
			"message": "如果邮箱已注册，重置链接已发送",
		})
	}

	// 生成唯一 token
	var token string
	for {
		token = guuid.NewV4().String()
		if _, loaded := resetPwdMap.LoadOrStore(token, &pwdResetEntry{
			email:    email,
			expireAt: time.Now().Add(30 * time.Minute),
		}); !loaded {
			break
		}
	}

	// 异步发送邮件
	go logic.DefaultEmail.SendResetpwdMail(email, token, CheckIsHttps(ctx))

	return success(ctx, map[string]interface{}{
		"message": "重置密码邮件已发送，请查收",
	})
}

// ResetPassword 通过 token 重置密码
func (AuthController) ResetPassword(ctx echo.Context) error {
	var req resetPasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "参数错误")
	}

	token := strings.TrimSpace(req.Token)
	newPassword := strings.TrimSpace(req.NewPassword)

	if token == "" {
		return fail(ctx, "重置链接无效")
	}
	if len(newPassword) < 6 || len(newPassword) > 32 {
		return fail(ctx, "密码长度必须在6到32个字符之间")
	}

	// 验证 token（含过期检查）
	val, ok := resetPwdMap.Load(token)
	if !ok {
		return fail(ctx, "重置链接已过期或无效，请重新申请")
	}
	entry, ok := val.(*pwdResetEntry)
	if !ok || time.Now().After(entry.expireAt) {
		resetPwdMap.Delete(token)
		return fail(ctx, "重置链接已过期，请重新申请")
	}
	email := entry.email

	// 重置密码
	_, err := logic.DefaultUser.ResetPasswd(context.EchoContext(ctx), email, newPassword)
	if err != nil {
		return fail(ctx, "重置密码失败，请重试")
	}

	// 清除 token
	resetPwdMap.Delete(token)

	return success(ctx, map[string]interface{}{
		"message": "密码重置成功，请重新登录",
	})
}
