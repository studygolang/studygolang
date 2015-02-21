// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/websocket"
	"logger"
	"service"
	"util"
)

var ServerId int
var mutex sync.Mutex

// websocket，统计在线用户数
// uri: /ws
func WsHandler(wsConn *websocket.Conn) {
	mutex.Lock()
	ServerId++
	serverId := ServerId
	mutex.Unlock()
	req := wsConn.Request()
	user, err := strconv.Atoi(req.FormValue("uid"))
	if err != nil || user == 0 {
		ip := util.Ip(req)
		logger.Debugln("user ip:", ip)
		user = int(util.Ip2long(ip))
	}
	userData := service.Book.AddUser(user, serverId)
	// 给自己发送消息，告诉当前在线用户数、历史最高在线人数
	onlineInfo := map[string]int{"online": service.Book.Len() + 50, "maxonline": service.MaxOnlineNum()}
	message := service.NewMessage(service.WsMsgOnline, onlineInfo)
	err = websocket.JSON.Send(wsConn, message)
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
			service.Book.DelUser(user, serverId)
			break
		}
	}
	// 用户退出时需要变更其他用户看到的在线用户数
	if !service.Book.UserIsOnline(user) {
		message := service.NewMessage(service.WsMsgOnline, map[string]int{"online": service.Book.Len() + 50})
		go service.Book.BroadcastAllUsersMessage(message)
	}
}
