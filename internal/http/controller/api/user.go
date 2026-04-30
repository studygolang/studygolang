// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dchest/captcha"
	"github.com/gorilla/sessions"
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
	guuid "github.com/twinj/uuid"
)

// pwdResetEntry 重置密码条目，带过期时间
type pwdResetEntry struct {
	email    string
	expireAt time.Time
}

// resetPwdMap 存储重置密码 token -> *pwdResetEntry 的映射（内存存储，重启失效）
var resetPwdMap = sync.Map{}

func init() {
	// 后台定期清理过期 token（每 10 分钟）
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			now := time.Now()
			resetPwdMap.Range(func(key, value interface{}) bool {
				if entry, ok := value.(*pwdResetEntry); ok && now.After(entry.expireAt) {
					resetPwdMap.Delete(key)
				}
				return true
			})
		}
	}()
}

type UserController struct{}

// ======================== 速率限制器 ========================

// rateLimiter 简单的内存速率限制器，按 IP 维度计数
type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*attemptInfo
}

// attemptInfo 记录某个 key 在时间窗口内的尝试次数
type attemptInfo struct {
	count    int
	expireAt time.Time
}

// check 检查是否超过限制，返回 true 表示允许，false 表示被限流
func (rl *rateLimiter) check(key string, maxAttempts int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.attempts[key]
	if !exists || now.After(info.expireAt) {
		rl.attempts[key] = &attemptInfo{count: 1, expireAt: now.Add(window)}
		return true
	}

	info.count++
	if info.count > maxAttempts {
		return false
	}
	return true
}

var (
	loginLimiter    = &rateLimiter{attempts: make(map[string]*attemptInfo)}
	registerLimiter = &rateLimiter{attempts: make(map[string]*attemptInfo)}
)

func (self UserController) RegisterRoute(g *echo.Group) {
	g.POST("/user/login", self.Login)
	g.POST("/user/logout", self.Logout)
	g.POST("/user/register", self.Register)
	g.GET("/user/me", self.Me)
	g.POST("/user/sync-session", self.SyncSession)
	g.GET("/users", self.List)
	g.GET("/user/:username", self.Home)
	// 用户设置相关（需要登录）
	g.GET("/user/profile", self.GetProfile)
	g.PUT("/user/profile", self.UpdateProfile)
	g.PUT("/user/avatar", self.UploadAvatar)
	g.PUT("/user/password", self.ChangePassword)
	// 用户评论列表
	g.GET("/users/:username/comments", self.UserComments)
	g.GET("/users/:username/topics", self.UserTopics)
	g.GET("/users/:username/articles", self.UserArticles)
	g.GET("/users/:username/resources", self.UserResources)
	g.GET("/users/:username/projects", self.UserProjects)
	// 密码找回
	g.POST("/user/forgot-password", self.ForgotPassword)
	g.POST("/user/reset-password", self.ResetPassword)
}

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

// Login 用户登录，返回 token
func (UserController) Login(ctx echo.Context) error {
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

// Me 当前登录用户信息（支持 Cookie 和 X-Token header）
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

func (UserController) Me(ctx echo.Context) error {
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
func (UserController) SyncSession(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
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
	uid, err := parseAuthUID(ctx)
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
	if len(req.Introduce) > 500 {
		return fail(ctx, "个人简介不能超过500个字符")
	}
	if len(req.Website) > 200 {
		return fail(ctx, "个人网站地址过长")
	}

	// 邮箱格式校验
	if req.Email != "" {
		if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
			return fail(ctx, "邮箱格式不正确")
		}
	}

	// 个人网站 URL 格式校验
	if req.Website != "" {
		website := strings.TrimSpace(req.Website)
		if !strings.HasPrefix(website, "http://") && !strings.HasPrefix(website, "https://") {
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

// changePasswordRequest 修改密码请求体
type changePasswordRequest struct {
	CurPasswd string `json:"cur_passwd"`
	NewPasswd string `json:"new_passwd"`
}

// ChangePassword 修改密码（需要登录）
func (UserController) ChangePassword(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
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
	uid, err := parseAuthUID(ctx)
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
		if strings.HasPrefix(lowerAvatar, "javascript:") {
			return fail(ctx, "非法的头像地址")
		}
	}

	err = logic.DefaultUser.ChangeAvatar(context.EchoContext(ctx), uid, avatar)
	if err != nil {
		return fail(ctx, "更换头像失败")
	}

	return success(ctx, nil)
}

// UserComments 获取指定用户的评论列表
// GET /api/v1/users/:username/comments?p=1
func (UserController) UserComments(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	comments := logic.DefaultComment.FindAll(context.EchoContext(ctx), paginator, "cid DESC", "uid=?", user.Uid)
	total := logic.DefaultComment.Count(context.EchoContext(ctx), "uid=?", user.Uid)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"comments": comments,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// UserTopics 获取指定用户的话题列表
// GET /api/v1/users/:username/topics?p=1
func (UserController) UserTopics(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	topics := logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "tid DESC", "uid=?", user.Uid)
	total := logic.DefaultTopic.Count(context.EchoContext(ctx), "uid=?", user.Uid)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"topics":  topics,
		"total":  total,
		"page":   curPage,
		"has_more": hasMore,
	})
}

// UserArticles 获取指定用户的文章列表
// GET /api/v1/users/:username/articles?p=1
func (UserController) UserArticles(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	articles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "author_txt=?", username)
	total := logic.DefaultArticle.Count(context.EchoContext(ctx), "author_txt=?", username)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"articles": articles,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// UserResources 获取指定用户的资源列表
// GET /api/v1/users/:username/resources?p=1
func (UserController) UserResources(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	resources, total := logic.DefaultResource.FindAll(context.EchoContext(ctx), paginator, "id DESC", "uid=?", user.Uid)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"resources": resources,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// UserProjects 获取指定用户的开源项目列表
// GET /api/v1/users/:username/projects?p=1
func (UserController) UserProjects(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	projects := logic.DefaultProject.FindAll(context.EchoContext(ctx), paginator, "id DESC", "username=?", username)
	total := logic.DefaultProject.Count(context.EchoContext(ctx), "username=?", username)

	hasMore := int64(curPage*paginator.PerPage()) < total

	return success(ctx, map[string]interface{}{
		"projects": projects,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
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

// ForgotPassword 发送重置密码邮件
func (UserController) ForgotPassword(ctx echo.Context) error {
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
	go logic.DefaultEmail.SendResetpwdMail(email, token)

	return success(ctx, map[string]interface{}{
		"message": "重置密码邮件已发送，请查收",
	})
}

// ResetPassword 通过 token 重置密码
func (UserController) ResetPassword(ctx echo.Context) error {
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
