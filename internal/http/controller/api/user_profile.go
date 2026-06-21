// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/mail"
	"net/url"
	"strings"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

type UserProfileController struct{}

// meResponse /api/v1/user/me 返回的安全用户信息（过滤敏感字段）
type meResponse struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	IsRoot   bool   `json:"is_root"`
	IsVip    bool   `json:"is_vip"`
	Balance  int    `json:"balance"`
	Gold     int    `json:"gold"`
	Silver   int    `json:"silver"`
	Copper   int    `json:"copper"`
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

// changePasswordRequest 修改密码请求体
type changePasswordRequest struct {
	CurPasswd string `json:"cur_passwd"`
	NewPasswd string `json:"new_passwd"`
}

// uploadAvatarRequest 上传头像请求体
type uploadAvatarRequest struct {
	Avatar string `json:"avatar"` // 头像 URL 或空字符串（使用 gravatar）
}

func (self UserProfileController) RegisterRoute(g *echo.Group) {
	g.GET("/user/me", self.Me)
	g.POST("/user/sync-session", self.SyncSession)
	g.GET("/users", self.List)
	g.GET("/user/:username", self.Home)
	g.GET("/user/profile", self.GetProfile)
	g.PUT("/user/profile", self.UpdateProfile)
	g.PUT("/user/avatar", self.UploadAvatar)
	g.PUT("/user/password", self.ChangePassword)
}

// Me 当前登录用户信息（支持 Cookie 和 X-Token header）
func (UserProfileController) Me(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
	if err != nil {
		return err
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	resp := &meResponse{
		Uid:      user.Uid,
		Username: user.Username,
		Avatar:   user.Avatar,
		Name:     user.Name,
		Status:   user.Status,
		IsRoot:   user.IsRoot,
		IsVip:    user.IsVip,
		Balance:  user.Balance,
		Gold:     user.Gold,
		Silver:   user.Silver,
		Copper:   user.Copper,
	}

	return success(ctx, map[string]interface{}{
		"user": resp,
	})
}

// SyncSession 同步 session（用于前端跳转到后端管理页面前建立 session）
// 从 token 中获取用户信息，设置到 session 中
func (UserProfileController) SyncSession(ctx echo.Context) error {
	// 用于跨子域跳进 /admin 等旧路由，必须要求账号已激活（冻结用户不允许建立 admin session）
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
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
func (UserProfileController) List(ctx echo.Context) error {
	// 获取活跃会员（按 DAU 排名）
	activeUsers := logic.DefaultRank.FindDAURank(context.EchoContext(ctx), 36)
	// 获取最新加入会员
	newUsers := logic.DefaultUser.FindNewUsers(context.EchoContext(ctx), 36)
	// 获取会员总数
	total := logic.DefaultUser.Total()

	// Email 脱敏：未勾选"公开 Email"的用户清空 Email 字段，
	// 与 master 模板 profile.html `{{if .user.Open}}` 语义对齐。
	sanitizeUsersForPublic(activeUsers)
	sanitizeUsersForPublic(newUsers)

	return success(ctx, map[string]interface{}{
		"active_users": activeUsers,
		"new_users":    newUsers,
		"total":        total,
	})
}

// Home 用户主页（话题、文章等）
func (UserProfileController) Home(ctx echo.Context) error {
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

	// Email 脱敏：访客访问他人主页时，未公开邮箱的用户清空 Email 字段，
	// 与 master 模板 profile.html `{{if .user.Open}}` 语义对齐。
	sanitizeUserForPublic(user)

	return success(ctx, map[string]interface{}{
		"user":      user,
		"topics":    topics,
		"articles":  articles,
		"resources": resources,
		"projects":  projects,
		"comments":  comments,
	})
}

// GetProfile 获取个人信息（需要登录）
func (UserProfileController) GetProfile(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
	if err != nil {
		return err
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

// UpdateProfile 更新个人信息（需要登录）
func (UserProfileController) UpdateProfile(ctx echo.Context) error {
	// 写操作：parseActiveAuthUID 同时校验用户状态（拦截冻结用户）
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	var req updateProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	// 输入校验：字段长度限制
	if len(req.Name) > 50 {
		return fail(ctx, "昵称不能超过50个字符")
	}
	if len(req.City) > 50 {
		return fail(ctx, "所在城市不能超过50个字符")
	}
	if len(req.Company) > 100 {
		return fail(ctx, "公司名称不能超过100个字符")
	}
	if len(req.Github) > 100 {
		return fail(ctx, "GitHub 账号过长")
	}
	if len(req.Introduce) > 500 {
		return fail(ctx, "个人简介不能超过500个字符")
	}
	if len(req.Website) > 200 {
		return fail(ctx, "个人网站地址过长")
	}

	// 邮箱格式校验（使用 net/mail 替代 Contains 判断，能拦截 a@.b / .@a 等垃圾）
	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return fail(ctx, "邮箱格式不正确")
		}
	}

	// 个人网站 URL 格式校验
	if req.Website != "" {
		website := strings.TrimSpace(req.Website)
		lower := strings.ToLower(website)
		if !strings.HasPrefix(lower, "http://") && !strings.HasPrefix(lower, "https://") {
			return fail(ctx, "个人网站地址必须以 http:// 或 https:// 开头")
		}
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

// ChangePassword 修改密码（需要登录）
func (UserProfileController) ChangePassword(ctx echo.Context) error {
	// 写操作：parseActiveAuthUID 同时校验用户状态（拦截冻结用户）
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	var req changePasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	if req.CurPasswd == "" {
		return fail(ctx, "当前密码不能为空")
	}
	if req.NewPasswd == "" {
		return fail(ctx, "新密码不能为空")
	}

	if len(req.NewPasswd) < 6 || len(req.NewPasswd) > 32 {
		return fail(ctx, "密码长度必须在6到32个字符之间")
	}

	if req.CurPasswd == req.NewPasswd {
		return fail(ctx, "新密码不能与当前密码相同")
	}

	errMsg, err := logic.DefaultUser.UpdatePasswd(context.EchoContext(ctx), user.Username, req.CurPasswd, req.NewPasswd)
	if err != nil {
		return fail(ctx, errMsg)
	}

	return success(ctx, nil)
}

// UploadAvatar 更换头像（需要登录）
func (UserProfileController) UploadAvatar(ctx echo.Context) error {
	// 写操作：parseActiveAuthUID 同时校验用户状态（拦截冻结用户）
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
	}

	var req uploadAvatarRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	// 头像 URL 安全校验：只允许 http/https 协议，防止 javascript: 等 XSS 攻击
	avatar := strings.TrimSpace(req.Avatar)
	if avatar != "" {
		lowerAvatar := strings.ToLower(avatar)
		if !strings.HasPrefix(lowerAvatar, "http://") && !strings.HasPrefix(lowerAvatar, "https://") {
			return fail(ctx, "头像地址必须以 http:// 或 https:// 开头")
		}
	}

	err = logic.DefaultUser.ChangeAvatar(context.EchoContext(ctx), uid, avatar)
	if err != nil {
		return fail(ctx, "更换头像失败")
	}

	return success(ctx, nil)
}
