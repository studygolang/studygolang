// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"net/http"
	"service"
	"strconv"
	"util"
)

// 发短消息
// uri: /message/send{json:(|.json)}
func SendMessageHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	content := req.FormValue("content")
	// 请求新建帖子页面
	if content == "" || req.Method != "POST" || vars["json"] == "" {
		user := service.FindUserByUsername(req.FormValue("username"))
		filter.SetData(req, map[string]interface{}{"user": user})
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/messages/send.html")
		return
	}

	user, _ := filter.CurrentUser(req)
	to := util.MustInt(req.FormValue("to"))
	success := service.SendMessageTo(user["uid"].(int), to, content)
	if !success {
		fmt.Fprint(rw, `{"errno": 1, "error":"对不起，发送失败，请稍候再试！"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "error":""}`)
}

// 消息列表
// uri: /message/{msgtype:(system|inbox|outbox)}
func MessageHandler(rw http.ResponseWriter, req *http.Request) {
	user, _ := filter.CurrentUser(req)
	uid := user["uid"].(int)
	vars := mux.Vars(req)
	msgtype := vars["msgtype"]
	var messages []map[string]interface{}
	if msgtype == "system" {
		messages = service.FindSysMsgsByUid(strconv.Itoa(uid))
	} else if msgtype == "inbox" {
		messages = service.FindToMsgsByUid(strconv.Itoa(uid))
	} else {
		messages = service.FindFromMsgsByUid(strconv.Itoa(uid))
	}
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/messages/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"messages": messages, "msgtype": msgtype})
}

// 删除消息
// uri: /message/delete.json
func DeleteMessageHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprint(rw, `{"errno": 1, "error":"非法请求！"}`)
		return
	}
	id := req.FormValue("id")
	msgtype := req.FormValue("msgtype")
	if !service.DeleteMessage(id, msgtype) {
		fmt.Fprint(rw, `{"errno": 1, "error":"对不起，删除失败，请稍候再试！"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "error":""}`)
}
