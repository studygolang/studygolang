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
	g.GET("/messages/unread-count", self.UnreadCount)
}

// List 消息列表（支持 Cookie 和 X-Token header）
// 查询参数：type=system|inbox|outbox, p=页码
func (MessageController) List(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	uid, _, valid := ValidateTokenAuto(token)
	if !valid || uid == 0 {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
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
	// 写操作必须校验用户状态（与 master NeedLogin 一致）
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
	}

	var req sendRequest
	if err := ctx.Bind(&req); err != nil {
		return fail(ctx, "请求参数错误")
	}

	if req.To == 0 {
		return fail(ctx, "收件人不能为空")
	}

	// 禁止自发自收（避免系统通知环路 / 自我刷屏）
	if req.To == uid {
		return fail(ctx, "不能给自己发送私信")
	}

	if req.Content == "" {
		return fail(ctx, "消息内容不能为空")
	}

	// 内容长度限制（防止超大 payload 拖垮 DB 与邮件队列）
	// 5000 字符约等于一条长私信，足以覆盖正常场景
	const maxMessageLen = 5000
	if len([]rune(req.Content)) > maxMessageLen {
		return fail(ctx, "消息内容过长（最多 5000 字符）")
	}

	// 检查收件人是否存在
	toUser := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", req.To)
	if toUser == nil || toUser.Uid == 0 {
		return fail(ctx, "收件人不存在")
	}

	if !logic.DefaultMessage.SendMessageTo(context.EchoContext(ctx), uid, req.To, req.Content) {
		return fail(ctx, "发送失败，请稍后重试")
	}

	return success(ctx, map[string]interface{}{
		"message": "发送成功",
	})
}

// Delete 删除消息（支持 Cookie 和 X-Token header）
// 查询参数：type=system|inbox|outbox
//
// CSRF 校验由全局 originCheck 中间件统一处理（routes.go: g.Use(originCheck)），
// 这里不再重复实现。鉴权用 parseActiveAuthUID 与 master NeedLogin 对齐：
// 冻结/未激活用户不允许删除消息。
func (MessageController) Delete(ctx echo.Context) error {
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
	}

	id := ctx.Param("id")
	if id == "" {
		return fail(ctx, "消息 ID 不能为空")
	}

	msgtype := ctx.QueryParam("type")

	// 验证 msgtype 参数（必填，避免误删 system vs inbox 路径上的消息）
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

// UnreadCount 获取未读消息计数（需要登录）
// GET /api/v1/messages/unread-count
// 返回 { code, data: { system, inbox }, message }
func (MessageController) UnreadCount(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
	if err != nil {
		return err
	}

	// 系统消息未读数
	sysUnread := logic.DefaultMessage.SysMsgUnreadCount(context.EchoContext(ctx), uid)
	// 私信未读数
	inboxUnread := logic.DefaultMessage.ToMsgUnreadCount(context.EchoContext(ctx), uid)

	return success(ctx, map[string]interface{}{
		"system": sysUnread,
		"inbox":  inboxUnread,
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
			// 严格匹配：必须以 ".domain" 结尾（避免 evilstudygolang.com 匹配 studygolang.com）
			// 使用 Hostname() 去掉端口（如 "host:port" → "host"）
			host := parsed.Hostname()
			if host != "" && strings.HasSuffix(host, "."+domain) {
				return true
			}
		}
	}

	return false
}
