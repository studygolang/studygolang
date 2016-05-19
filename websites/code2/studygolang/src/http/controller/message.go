// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"net/http"

	"http/middleware"
	"logic"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"

	. "http"
	"model"
)

type MessageController struct{}

// 注册路由
func (self MessageController) RegisterRoute(e *echo.Echo) {
	e.Get("/message/:msgtype", echo.HandlerFunc(self.ReadList), middleware.NeedLogin())
	e.Get("/message/system", echo.HandlerFunc(self.ReadList), middleware.NeedLogin())
	e.Match([]string{"GET", "POST"}, "/message/send", echo.HandlerFunc(self.Send), middleware.NeedLogin())
	e.Post("/message/delete", echo.HandlerFunc(self.Delete), middleware.NeedLogin())
}

// Send 发短消息
func (MessageController) Send(ctx echo.Context) error {
	content := ctx.FormValue("content")
	// 请求发送消息页面
	if content == "" || Request(ctx).Method != "POST" {
		username := ctx.FormValue("username")
		if username == "" {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		user := logic.DefaultUser.FindOne(ctx, "username", username)
		return render(ctx, "messages/send.html", map[string]interface{}{"user": user})
	}

	user := ctx.Get("user").(*model.Me)
	to := goutils.MustInt(ctx.FormValue("to"))
	ok := logic.DefaultMessage.SendMessageTo(ctx, user.Uid, to, content)
	if !ok {
		return fail(ctx, 1, "对不起，发送失败，请稍候再试！")
	}

	return success(ctx, nil)
}

// 消息列表
func (MessageController) ReadList(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Me)
	msgtype := ctx.Param("msgtype")
	if msgtype == "" {
		msgtype = "system"
	}

	var messages []map[string]interface{}
	switch msgtype {
	case "system":
		messages = logic.DefaultMessage.FindSysMsgsByUid(ctx, user.Uid)
	case "inbox":
		messages = logic.DefaultMessage.FindToMsgsByUid(ctx, user.Uid)
	case "outbox":
		messages = logic.DefaultMessage.FindFromMsgsByUid(ctx, user.Uid)
	default:
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	return render(ctx, "messages/list.html", map[string]interface{}{"messages": messages, "msgtype": msgtype})
}

// 删除消息
func (MessageController) Delete(ctx echo.Context) error {
	id := ctx.FormValue("id")
	msgtype := ctx.FormValue("msgtype")
	if !logic.DefaultMessage.DeleteMessage(ctx, id, msgtype) {
		return fail(ctx, 1, "对不起，删除失败，请稍候再试！")
	}

	return success(ctx, nil)
}
