// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

type UserController struct{}

func (self UserController) RegisterRoute(g *echo.Group) {
	g.POST("/user/login", self.Login)
	g.POST("/user/logout", self.Logout)
	g.POST("/user/register", self.Register)
	g.GET("/user/me", self.Me)
	g.POST("/user/sync-session", self.SyncSession)
	g.GET("/users", self.List)
	g.GET("/user/:username", self.Home)
}

// loginRequest 登录请求体（支持 JSON）
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Passwd   string `json:"passwd"` // 兼容两种字段名
}

// registerRequest 注册请求体（支持 JSON）
type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
	Password string `json:"password"` // 兼容两种字段名
}

// Login 用户登录，返回 token
func (UserController) Login(ctx echo.Context) error {
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

	token := GenToken(userLogin.Uid)
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
func (UserController) Logout(ctx echo.Context) error {
	clearAuthCookie(ctx)
	// 同时清除旧的 session（用于 /admin 等旧路由）
	session := GetCookieSession(ctx)
	session.Options = &sessions.Options{Path: "/", MaxAge: -1}
	session.Save(Request(ctx), ResponseWriter(ctx))
	return success(ctx, nil)
}

// Register 用户注册
func (UserController) Register(ctx echo.Context) error {
	var req registerRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
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

	setAuthCookie(ctx, GenToken(userLogin.Uid))
	// 同时设置旧的 session（用于 /admin 等旧路由）
	SetLoginCookie(ctx, userLogin.Username)

	return success(ctx, map[string]interface{}{
		"uid":      userLogin.Uid,
		"username": userLogin.Username,
	})
}

// Me 当前登录用户信息（支持 Cookie 和 X-Token header）
func (UserController) Me(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	return success(ctx, map[string]interface{}{
		"user": user,
	})
}

// SyncSession 同步 session（用于前端跳转到后端管理页面前建立 session）
// 从 token 中获取用户信息，设置到 session 中
func (UserController) SyncSession(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	// 设置 session（用于后端管理页面）
	SetLoginCookie(ctx, user.Username)

	return success(ctx, map[string]interface{}{
		"username": user.Username,
	})
}

// List 会员列表（活跃会员 + 新加入会员）
func (UserController) List(ctx echo.Context) error {
	// 获取活跃会员（按 DAU 排名）
	activeUsers := logic.DefaultRank.FindDAURank(context.EchoContext(ctx), 36)
	// 获取最新加入会员
	newUsers := logic.DefaultUser.FindNewUsers(context.EchoContext(ctx), 36)
	// 获取会员总数
	total := logic.DefaultUser.Total()

	return success(ctx, map[string]interface{}{
		"active_users": activeUsers,
		"new_users":    newUsers,
		"total":        total,
	})
}

// Home 用户主页（话题、文章等）
func (UserController) Home(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 || user.Status == model.UserStatusOutage {
		return fail(ctx, "用户不存在")
	}

	topics := logic.DefaultTopic.FindRecent(5, user.Uid)
	articles := logic.DefaultArticle.FindByUser(context.EchoContext(ctx), user.Username, 5)
	resources := logic.DefaultResource.FindRecent(context.EchoContext(ctx), user.Uid)
	projects := logic.DefaultProject.FindRecent(context.EchoContext(ctx), user.Username)
	comments := logic.DefaultComment.FindRecent(context.EchoContext(ctx), user.Uid, -1, 5)

	return success(ctx, map[string]interface{}{
		"user":      user,
		"topics":    topics,
		"articles":  articles,
		"resources": resources,
		"projects":  projects,
		"comments":  comments,
	})
}

// ======================== 个人设置相关 API ========================

// GetProfile 获取个人信息（需要登录）
func (UserController) GetProfile(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	hasPasswd := logic.DefaultUser.HasPasswd(context.EchoContext(ctx), uid)

	return success(ctx, map[string]interface{}{
		"user":       user,
		"has_passwd": hasPasswd,
	})
}

// updateProfileRequest 更新个人信息请求体
type updateProfileRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	City      string `json:"city"`
	Company   string `json:"company"`
	Github    string `json:"github"`
	Website   string `json:"website"`
	Introduce string `json:"introduce"`
	Open      string `json:"open"` // "1" 或 "0"
}

// UpdateProfile 更新个人信息（需要登录）
func (UserController) UpdateProfile(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	var req updateProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	// 构造 url.Values 传递给 logic 层
	form := url.Values{}
	form.Set("name", req.Name)
	form.Set("email", req.Email)
	form.Set("city", req.City)
	form.Set("company", req.Company)
	form.Set("github", req.Github)
	form.Set("website", req.Website)
	form.Set("introduce", req.Introduce)
	form.Set("open", req.Open)

	me := &model.Me{
		Uid:      uid,
		Username: user.Username,
		Email:    user.Email,
	}

	errMsg, err := logic.DefaultUser.Update(context.EchoContext(ctx), me, form)
	if err != nil {
		return fail(ctx, errMsg)
	}

	return success(ctx, nil)
}

// changePasswordRequest 修改密码请求体
type changePasswordRequest struct {
	CurPasswd string `json:"cur_passwd"`
	NewPasswd string `json:"new_passwd"`
}

// ChangePassword 修改密码（需要登录）
func (UserController) ChangePassword(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	var req changePasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	if req.NewPasswd == "" {
		return fail(ctx, "新密码不能为空")
	}

	if len(req.NewPasswd) < 6 || len(req.NewPasswd) > 32 {
		return fail(ctx, "密码长度必须在6到32个字符之间")
	}

	errMsg, err := logic.DefaultUser.UpdatePasswd(context.EchoContext(ctx), user.Username, req.CurPasswd, req.NewPasswd)
	if err != nil {
		return fail(ctx, errMsg)
	}

	return success(ctx, nil)
}

// uploadAvatarRequest 上传头像请求体
type uploadAvatarRequest struct {
	Avatar string `json:"avatar"` // 头像 URL 或空字符串（使用 gravatar）
}

// UploadAvatar 更换头像（需要登录）
func (UserController) UploadAvatar(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	var req uploadAvatarRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	err := logic.DefaultUser.ChangeAvatar(context.EchoContext(ctx), uid, req.Avatar)
	if err != nil {
		return fail(ctx, "更换头像失败")
	}

	return success(ctx, nil)
}
