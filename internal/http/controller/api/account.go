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

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/http/internal/helper"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type AccountController struct{}

// RegisterRoute 注册账号相关路由
func (self AccountController) RegisterRoute(g *echo.Group) {
	// 公开路由
	g.GET("/account/activate", self.Activate)
	g.POST("/account/send-activate-email", self.SendActivateEmail)
	g.GET("/account/email/unsubscribe", self.UnsubscribePage)
	g.POST("/account/email/unsubscribe", self.Unsubscribe)
	// 需要登录
	g.GET("/account/bind_users", self.BindUsers)
	g.POST("/account/social/unbind", self.SocialUnbind)
}

// SendActivateEmail 发送激活邮件（需要登录）
// POST /api/v1/account/send-activate-email
func (AccountController) SendActivateEmail(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", me.Uid)
	if user == nil || user.Email == "" {
		return fail(ctx, "用户不存在或未设置邮箱")
	}

	// 生成激活 UUID 并发送邮件
	email := strings.TrimSpace(user.Email)
	isHttps := CheckIsHttps(ctx)
	// 生成真实 uuid 并登记 email 映射（master 用 RegActivateCode.GenUUID）
	uuid := helper.RegActivateCode.GenUUID(email)
	logic.DefaultEmail.SendActivateMail(email, uuid, isHttps)

	return success(ctx, map[string]interface{}{
		"message": "激活邮件已发送",
	})
}

// Activate 处理账号激活链接
// GET /api/v1/account/activate?param=xxx
func (AccountController) Activate(ctx echo.Context) error {
	param := ctx.QueryParam("param")
	if param == "" {
		return fail(ctx, "缺少激活参数")
	}

	// Base64 解码参数
	decoded := goutils.Base64Decode(param)
	values, err := url.ParseQuery(decoded)
	if err != nil {
		return fail(ctx, "参数格式错误")
	}

	uuid := values.Get("uuid")
	timestamp := goutils.MustInt64(values.Get("timestamp"))
	sign := values.Get("sign")

	if uuid == "" {
		return fail(ctx, "激活链接不完整")
	}

	// 通过 uuid 反查 email（master 用 RegActivateCode.GetEmail），链接里不含 email
	email, ok := helper.RegActivateCode.GetEmail(uuid)
	if !ok {
		return fail(ctx, "非法请求")
	}

	// 激活链接有效期 4 小时（master 同此校验）
	if timestamp < time.Now().Add(-4*time.Hour).Unix() {
		helper.RegActivateCode.DelUUID(uuid)
		return fail(ctx, "激活链接不完整")
	}

	user, err := logic.DefaultUser.Activate(context.EchoContext(ctx), email, uuid, timestamp, sign)
	if err != nil {
		return fail(ctx, "激活失败："+err.Error())
	}

	// 激活成功后清除 uuid 映射
	helper.RegActivateCode.DelUUID(uuid)

	return success(ctx, map[string]interface{}{
		"message": "激活成功",
		"uid":     user.Uid,
	})
}

// socialUnbindRequest 解绑社交账号请求体
type socialUnbindRequest struct {
	BindID   int    `json:"bind_id"`
	Platform string `json:"platform"` // github 或 gitea
}

// SocialUnbind 解绑社交账号（需要登录）
// POST /api/v1/account/social/unbind
func (AccountController) SocialUnbind(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	var req socialUnbindRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	if req.BindID == 0 {
		return fail(ctx, "绑定 ID 不合法")
	}

	platform := strings.TrimSpace(req.Platform)
	if platform != "github" && platform != "gitea" {
		return fail(ctx, "不支持该平台")
	}

	if err := logic.DefaultThirdUser.UnBindUser(context.EchoContext(ctx), req.BindID, me); err != nil {
		return fail(ctx, "解绑失败："+err.Error())
	}

	return success(ctx, map[string]interface{}{
		"message": "解绑成功",
	})
}

// UnsubscribePage 邮件退订页面数据
// GET /api/v1/account/email/unsubscribe?token=xxx&email=xxx
func (AccountController) UnsubscribePage(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	email := ctx.QueryParam("email")

	if token == "" || email == "" {
		return fail(ctx, "参数不完整")
	}

	// 校验 token 的合法性
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "email", email)
	if user == nil || user.Email == "" {
		return fail(ctx, "用户不存在")
	}

	realToken := logic.DefaultEmail.GenUnsubscribeToken(user)
	if token != realToken {
		return fail(ctx, "验证失败")
	}

	return success(ctx, map[string]interface{}{
		"email":       email,
		"token":       token,
		"unsubscribe": user.Unsubscribe,
	})
}

// Unsubscribe 执行邮件退订
// POST /api/v1/account/email/unsubscribe
func (AccountController) Unsubscribe(ctx echo.Context) error {
	token := ctx.FormValue("token")
	email := ctx.FormValue("email")

	if token == "" || email == "" {
		return fail(ctx, "参数不完整")
	}

	// 校验 token 的合法性
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "email", email)
	if user == nil || user.Email == "" {
		return fail(ctx, "用户不存在")
	}

	realToken := logic.DefaultEmail.GenUnsubscribeToken(user)
	if token != realToken {
		return fail(ctx, "验证失败")
	}

	unsubscribe := goutils.MustInt(ctx.FormValue("unsubscribe"))
	logic.DefaultUser.EmailSubscribe(context.EchoContext(ctx), user.Uid, unsubscribe)

	return success(ctx, map[string]interface{}{
		"message": "操作成功",
	})
}

// bindUserDTO 绑定账号响应（过滤敏感字段，不暴露 access_token/refresh_token）
type bindUserDTO struct {
	ID       int    `json:"id"`
	Platform string `json:"platform"` // "github" 或 "gitea"
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

// bindPlatformName 将绑定类型常量转为平台名称
func bindPlatformName(t int) string {
	switch t {
	case model.BindTypeGithub:
		return "github"
	case model.BindTypeGitea:
		return "gitea"
	default:
		return "unknown"
	}
}

// BindUsers 获取已绑定的社交账号列表（需要登录）
// GET /api/v1/account/bind_users
func (AccountController) BindUsers(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	bindUsers := logic.DefaultUser.FindBindUsers(context.EchoContext(ctx), me.Uid)

	dtos := make([]bindUserDTO, 0, len(bindUsers))
	for _, bu := range bindUsers {
		dtos = append(dtos, bindUserDTO{
			ID:       bu.Id,
			Platform: bindPlatformName(bu.Type),
			Username: bu.Username,
			Name:     bu.Name,
			Avatar:   bu.Avatar,
		})
	}

	return success(ctx, map[string]interface{}{
		"bind_users": dtos,
	})
}
