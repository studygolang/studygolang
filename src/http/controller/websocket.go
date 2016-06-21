// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"

	"golang.org/x/net/websocket"
)

type WebsocketController struct {
	ServerId int
	mutex    sync.Mutex
}

func (this *WebsocketController) RegisterRoute(g *echo.Group) {
	g.GET("/ws", standard.WrapHandler(websocket.Handler(this.Ws)))
}

// websocket，统计在线用户数
// uri: /ws
func (this *WebsocketController) Ws(wsConn *websocket.Conn) {
	defer wsConn.Close()

	this.mutex.Lock()
	this.ServerId++
	serverId := this.ServerId
	this.mutex.Unlock()
	req := wsConn.Request()
	user := goutils.MustInt(req.FormValue("uid"))
	if user == 0 {
		user = int(goutils.Ip2long(goutils.RemoteIp(req)))
	}
	userData := logic.Book.AddUser(user, serverId)
	// 给自己发送消息，告诉当前在线用户数、历史最高在线人数
	onlineInfo := map[string]int{"online": logic.Book.Len(), "maxonline": logic.MaxOnlineNum()}
	message := logic.NewMessage(logic.WsMsgOnline, onlineInfo)
	err := websocket.JSON.Send(wsConn, message)
	if err != nil {
		logger.Errorln("Sending onlineusers error:", err)
	}
	var clientClosed = false
	for {
		select {
		case message := <-userData.MessageQueue(serverId):
			if err := websocket.JSON.Send(wsConn, message); err != nil {
				clientClosed = true
			}
			// 心跳
		case <-time.After(30e9):
			if err := websocket.JSON.Send(wsConn, ""); err != nil {
				clientClosed = true
			}
		}
		if clientClosed {
			logic.Book.DelUser(user, serverId)
			logger.Infoln("user:", user, "client close")
			break
		}
	}
	// 用户退出时需要变更其他用户看到的在线用户数
	if !logic.Book.UserIsOnline(user) {
		message := logic.NewMessage(logic.WsMsgOnline, map[string]int{"online": logic.Book.Len()})
		go logic.Book.BroadcastAllUsersMessage(message)
	}
}
