// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type MessageController struct{}

func (self MessageController) RegisterRoute(g *echo.Group) {
	g.GET("/messages", self.List)
	g.POST("/messages", self.Send)
	g.DELETE("/messages/:id", self.Delete)
}

// List 消息列表（支持 Cookie 和 X-Token header）
// 查询参数：type=system|inbox|outbox, p=页码
func (MessageController) List(ctx echo.Context) error {
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

	msgtype := ctx.QueryParam("type")
	if msgtype == "" {
		msgtype = "system"
	}

	// 验证 msgtype 参数
	if msgtype != "system" && msgtype != "inbox" && msgtype != "outbox" {
		return fail(ctx, "参数有误：type 必须是 system、inbox 或 outbox")
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	var (
		messages []map[string]interface{}
		total    int64
	)

	switch msgtype {
	case "system":
		messages = logic.DefaultMessage.FindSysMsgsByUid(context.EchoContext(ctx), uid, paginator)
		total = logic.DefaultMessage.SysMsgCount(context.EchoContext(ctx), uid)
	case "inbox":
		messages = logic.DefaultMessage.FindToMsgsByUid(context.EchoContext(ctx), uid, paginator)
		total = logic.DefaultMessage.ToMsgCount(context.EchoContext(ctx), uid)
	case "outbox":
		messages = logic.DefaultMessage.FindFromMsgsByUid(context.EchoContext(ctx), uid, paginator)
		total = logic.DefaultMessage.FromMsgCount(context.EchoContext(ctx), uid)
	}

	return success(ctx, map[string]interface{}{
		"messages": messages,
		"page":     curPage,
		"total":    total,
		"has_more": int64(curPage*paginator.PerPage()) < total,
	})
}

// sendRequest 发送私信请求体
type sendRequest struct {
	To      int    `json:"to"`
	Content string `json:"content"`
}

// Send 发送私信（支持 Cookie 和 X-Token header）
func (MessageController) Send(ctx echo.Context) error {
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

	var req sendRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	if req.To == 0 {
		return fail(ctx, "收件人不能为空")
	}

	if req.Content == "" {
		return fail(ctx, "消息内容不能为空")
	}

	// 检查收件人是否存在
	toUser := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", req.To)
	if toUser == nil || toUser.Uid == 0 {
		return fail(ctx, "收件人不存在")
	}

	ok = logic.DefaultMessage.SendMessageTo(context.EchoContext(ctx), uid, req.To, req.Content)
	if !ok {
		return fail(ctx, "发送失败，请稍后重试")
	}

	return success(ctx, map[string]interface{}{
		"message": "发送成功",
	})
}

// Delete 删除消息（支持 Cookie 和 X-Token header）
// 查询参数：type=system|inbox|outbox
func (MessageController) Delete(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
	if err != nil {
		return err
	}

	// CSRF 保护：验证 Origin 或 Referer 头
	origin := ctx.Request().Header.Get("Origin")
	referer := ctx.Request().Header.Get("Referer")
	if origin == "" && referer == "" {
		return fail(ctx, "缺少 Origin 或 Referer 头", 403)
	}

	// 验证来源是否合法（对照 ALLOWED_ORIGINS 白名单）
	allowedOrigins := getAllowedOrigins()
	if origin != "" && !isOriginAllowed(origin, allowedOrigins) {
		return fail(ctx, "非法的跨域请求", 403)
	}
	if referer != "" && !isOriginAllowed(referer, allowedOrigins) {
		return fail(ctx, "非法的跨域请求", 403)
	}

	id := ctx.Param("id")
	if id == "" {
		return fail(ctx, "消息 ID 不能为空")
	}

	msgtype := ctx.QueryParam("type")
	if msgtype == "" {
		msgtype = "system"
	}

	// 验证 msgtype 参数
	if msgtype != "system" && msgtype != "inbox" && msgtype != "outbox" {
		return fail(ctx, "参数有误：type 必须是 system、inbox 或 outbox")
	}

	// 验证消息所有权（防止删除他人消息）
	if msgtype == "system" {
		// 系统消息：验证是否是发给当前用户的
		sysMsg := logic.DefaultMessage.FindSysMsgById(context.EchoContext(ctx), id)
		if sysMsg == nil || sysMsg.To != uid {
			return fail(ctx, "无权删除此消息")
		}
	} else {
		// 私信：验证是否是当前用户的收件箱或发件箱
		msg := logic.DefaultMessage.FindMsgById(context.EchoContext(ctx), id)
		if msg == nil {
			return fail(ctx, "消息不存在")
		}
		if msgtype == "inbox" && msg.To != uid {
			return fail(ctx, "无权删除此消息")
		}
		if msgtype == "outbox" && msg.From != uid {
			return fail(ctx, "无权删除此消息")
		}
	}

	ok := logic.DefaultMessage.DeleteMessage(context.EchoContext(ctx), id, msgtype)
	if !ok {
		return fail(ctx, "删除失败，请稍后重试")
	}

	return success(ctx, map[string]interface{}{
		"message": "删除成功",
	})
}

// isOriginAllowed 检查 origin 是否在白名单中
func isOriginAllowed(originOrReferer string, allowedOrigins []string) bool {
	parsed, err := url.Parse(originOrReferer)
	if err != nil {
		return false
	}

	// 提取 origin（scheme + host + port）
	origin := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)

	for _, allowed := range allowedOrigins {
		// 精确匹配
		if origin == allowed {
			return true
		}
		// 支持通配符子域名（如 *.studygolang.com）
		if strings.HasPrefix(allowed, "*.") {
			domain := allowed[2:] // 去掉 "*."
			if strings.HasSuffix(parsed.Host, domain) {
				return true
			}
		}
	}

	return false
}

// isSameOrigin 检查请求来源是否与目标主机同源（已废弃，使用 isOriginAllowed）
// Deprecated: 使用 isOriginAllowed 对照白名单验证
func isSameOrigin(originOrReferer, targetHost string) bool {
	// 简单检查：提取 origin/referer 中的 host 部分
	// 格式：http(s)://host:port/path
	parsed, err := url.Parse(originOrReferer)
	if err != nil {
		return false
	}

	// 比较主机名（忽略端口）
	originHost := parsed.Host
	if strings.Contains(originHost, ":") {
		originHost = strings.Split(originHost, ":")[0]
	}

	targetHostClean := targetHost
	if strings.Contains(targetHostClean, ":") {
		targetHostClean = strings.Split(targetHostClean, ":")[0]
	}

	return originHost == targetHostClean
}
