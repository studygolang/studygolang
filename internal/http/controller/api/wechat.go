// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"io"
	"net/http"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
)

type WechatController struct{}

func (self WechatController) RegisterRoute(g *echo.Group) {
	g.Any("/wechat/autoreply", self.AutoReply)
	g.POST("/wechat/bind", self.Bind)
}

// AutoReply 微信自动回复回调
// 微信服务器验证: GET 请求带 echostr 参数时原样返回
// 消息接收: POST 请求解析 XML 并自动回复
func (WechatController) AutoReply(ctx echo.Context) error {
	// 微信服务器配置验证
	if ctx.QueryParam("echostr") != "" {
		return ctx.String(http.StatusOK, ctx.QueryParam("echostr"))
	}

	body, err := io.ReadAll(io.LimitReader(ctx.Request().Body, 1<<20))
	if err != nil {
		return ctx.String(http.StatusOK, "")
	}

	if len(body) == 0 {
		return ctx.String(http.StatusOK, "")
	}

	wechatReply, err := logic.DefaultWechat.AutoReply(context.EchoContext(ctx), body)
	if err != nil {
		return ctx.String(http.StatusOK, "")
	}

	return ctx.XML(http.StatusOK, wechatReply)
}

// Bind 绑定微信公众号（需要登录 + 验证码）
func (WechatController) Bind(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	captcha := ctx.FormValue("captcha")
	if captcha == "" {
		return fail(ctx, "验证码不能为空")
	}

	err = logic.DefaultWechat.CheckCaptchaAndBind(context.EchoContext(ctx), me, captcha)
	if err != nil {
		return fail(ctx, "验证码错误，请确认获取了或没填错！")
	}

	return success(ctx, nil)
}
