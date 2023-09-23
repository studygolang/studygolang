// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"sync/atomic"
	"time"

	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"golang.org/x/net/websocket"
)

type WebsocketController struct {
	ServerId uint32
}

func (this *WebsocketController) RegisterRoute(g *echo.Group) {
	g.GET("/ws", echo.WrapHandler(websocket.Handler(this.Ws)))
}

// websocket，统计在线用户数
// uri: /ws
func (this *WebsocketController) Ws(wsConn *websocket.Conn) {
	defer wsConn.Close()

	serverId := int(atomic.AddUint32(&this.ServerId, 1))

	isUid := true
	req := wsConn.Request()
	user := goutils.MustInt(req.FormValue("uid"))
	if user == 0 {
		user = int(goutils.Ip2long(goutils.RemoteIp(req)))
		isUid = false
	}
	userData := logic.Book.AddUser(user, serverId, isUid)
	// 给自己发送消息，告诉当前在线用户数、历史最高在线人数
	onlineInfo := map[string]int{"online": logic.Book.Len(), "maxonline": logic.MaxOnlineNum()}
	message := logic.NewMessage(logic.WsMsgOnline, onlineInfo)
	err := websocket.JSON.Send(wsConn, message)
	if err != nil {
		logger.Errorln("Sending onlineusers error:", err)
		return
	}

	messageChan := userData.MessageQueue(serverId)

	ticker := time.NewTicker(15e9)
	defer ticker.Stop()

	var clientClosed = false
	for {
		select {
		case message := <-messageChan:
			if err := websocket.JSON.Send(wsConn, message); err != nil {
				// logger.Errorln("Send message", message, " to user:", user, "server_id:", serverId, "error:", err)
				clientClosed = true
			}
			// 心跳
		case <-ticker.C:
			if err := websocket.JSON.Send(wsConn, ""); err != nil {
				// logger.Errorln("Send heart message to user:", user, "server_id:", serverId, "error:", err)
				clientClosed = true
			}
		}
		if clientClosed {
			logic.Book.DelUser(user, serverId, isUid)
			logger.Infoln("user:", user, "client close")
			break
		}
	}
	// 用户退出时需要变更其他用户看到的在线用户数
	if !logic.Book.UserIsOnline(user) {
		logger.Infoln("user:", user, "had leave")

		message := logic.NewMessage(logic.WsMsgOnline, map[string]int{"online": logic.Book.Len()})
		go logic.Book.BroadcastAllUsersMessage(message)
	}
}
