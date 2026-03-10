// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/url"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

type UserController struct{}

func (self UserController) RegisterRoute(g *echo.Group) {
	g.POST("/user/login", self.Login)
	g.POST("/user/register", self.Register)
	g.GET("/user/me", self.Me)
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

	return success(ctx, map[string]interface{}{
		"token":    GenToken(userLogin.Uid),
		"uid":      userLogin.Uid,
		"username": userLogin.Username,
	})
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
			errMsg = err.Error()
		}
		return fail(ctx, errMsg)
	}

	// 注册成功后自动登录，返回 token
	userLogin, err := logic.DefaultUser.Login(context.EchoContext(ctx), username, passwd)
	if err != nil {
		return success(ctx, map[string]interface{}{
			"username": username,
		})
	}

	return success(ctx, map[string]interface{}{
		"token":    GenToken(userLogin.Uid),
		"uid":      userLogin.Uid,
		"username": userLogin.Username,
	})
}

// Me 当前登录用户信息（需要 X-Token header）
func (UserController) Me(ctx echo.Context) error {
	token := ctx.Request().Header.Get("X-Token")
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
