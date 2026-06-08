// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
)

const (
	// wsWriteWait 写超时时间
	wsWriteWait = 10 * time.Second
	// wsPongWait Pong 响应超时时间
	wsPongWait = 60 * time.Second
	// wsPingInterval Ping 发送间隔（必须小于 wsPongWait）
	wsPingInterval = 30 * time.Second
	// wsMaxMessageSize 单条消息最大字节数
	wsMaxMessageSize = 512
)

// wsUpgrader WebSocket 升级器
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 校验 Origin 白名单，防止跨站 WebSocket 劫持（CSWSH）
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // 非浏览器客户端（如 curl、Go client）无 Origin
		}
		return isOriginAllowed(origin, getAllowedOrigins())
	},
}

// WSMessage WebSocket 消息格式
type WSMessage struct {
	Type string      `json:"type"` // 消息类型: unread_count, new_message, notification
	Data interface{} `json:"data"` // 消息数据
}

// WSClient WebSocket 客户端连接
type WSClient struct {
	hub  *WSHub
	uid  int
	conn *websocket.Conn
	send chan []byte
}

// sendToRequest 定向发送请求
type sendToRequest struct {
	uid  int
	data []byte
}

// WSHub 管理所有 WebSocket 连接和消息广播
type WSHub struct {
	// 已注册的客户端（按 uid 索引，每用户限制 1 个连接）
	clients map[int]*WSClient

	// 注册请求
	register chan *WSClient
	// 注销请求
	unregister chan *WSClient
	// 广播消息
	broadcast chan *WSMessage
	// 定向发送消息（避免 SendToUser 的 send-on-closed-channel 竞态）
	sendTo chan *sendToRequest

	mu sync.RWMutex
}

// globalWSHub 全局 WebSocket Hub 单例
var globalWSHub *WSHub

func init() {
	globalWSHub = NewWSHub()
	go globalWSHub.Run()
}

// NewWSHub 创建新的 WebSocket Hub
func NewWSHub() *WSHub {
	return &WSHub{
		clients:    make(map[int]*WSClient),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		broadcast:  make(chan *WSMessage, 64),
		sendTo:     make(chan *sendToRequest, 64),
	}
}

// Run 启动 Hub 事件循环
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// 如果该用户已有连接，关闭旧连接
			if old, ok := h.clients[client.uid]; ok {
				close(old.send)
				old.conn.Close()
			}
			h.clients[client.uid] = client
			h.mu.Unlock()
			// 注册完成后安全地推送初始未读计数（在 Hub 事件循环内，无竞态）
			go h.pushInitialUnread(client)

		case client := <-h.unregister:
			h.mu.Lock()
			if existing, ok := h.clients[client.uid]; ok && existing == client {
				delete(h.clients, client.uid)
				close(client.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			// 使用写锁遍历，避免 RLock→Lock 升级导致死锁
			h.mu.Lock()
			for uid, client := range h.clients {
				select {
				case client.send <- data:
				default:
					// 发送缓冲区满，关闭连接
					delete(h.clients, uid)
					close(client.send)
				}
			}
			h.mu.Unlock()

		case req := <-h.sendTo:
			// 定向发送：在 Hub 事件循环内处理，避免 send-on-closed-channel 竞态
			h.mu.RLock()
			client, ok := h.clients[req.uid]
			h.mu.RUnlock()
			if !ok {
				continue
			}
			select {
			case client.send <- req.data:
			default:
				// 缓冲区满，静默丢弃
			}
		}
	}
}

// SendToUser 向指定用户发送消息（通过 Hub 事件循环，避免竞态）
func (h *WSHub) SendToUser(uid int, msg *WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case h.sendTo <- &sendToRequest{uid: uid, data: data}:
	default:
		// sendTo 缓冲区满，静默丢弃
	}
}

// pushInitialUnread 推送初始未读消息计数（在 register case 中调用，确保 client 已注册）
func (h *WSHub) pushInitialUnread(client *WSClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sysUnread := logic.DefaultMessage.SysMsgUnreadCount(ctx, client.uid)
	inboxUnread := logic.DefaultMessage.ToMsgUnreadCount(ctx, client.uid)

	msg := &WSMessage{
		Type: "unread_count",
		Data: map[string]int64{
			"system": sysUnread,
			"inbox":  inboxUnread,
		},
	}
	h.SendToUser(client.uid, msg)
}

// OnlineCount 返回当前在线连接数
func (h *WSHub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// readPump 从 WebSocket 连接读取消息（保活 + 心跳检测）
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(wsMaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(wsPongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// writePump 向 WebSocket 连接写入消息（心跳 + 消息发送）
func (c *WSClient) writePump() {
	ticker := time.NewTicker(wsPingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 排队消息批量发送
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WSController WebSocket 控制器
type WSController struct{}

func (WSController) RegisterRoute(g *echo.Group) {
	g.GET("/ws", WSHandler)
}

// WSHandler 处理 WebSocket 连接升级
func WSHandler(ctx echo.Context) error {
	// 从 Cookie 解析 uid（复用 parseAuthUID 逻辑）
	token := getAuthToken(ctx)
	if token == "" {
		return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    1,
			"message": "请先登录",
		})
	}

	uid, _, valid := ValidateTokenAuto(token)
	if !valid || uid == 0 {
		return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    1,
			"message": "token 已过期，请重新登录",
		})
	}

	// 升级为 WebSocket 连接
	conn, err := wsUpgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return nil // Upgrade 已发送 HTTP 错误响应
	}

	client := &WSClient{
		hub:  globalWSHub,
		uid:  uid,
		conn: conn,
		send: make(chan []byte, 64),
	}

	// 注册到 Hub（Run 的 register case 会通过 pushInitialUnread 推送初始未读计数）
	client.hub.register <- client

	// 在独立 goroutine 中运行读写
	go client.writePump()
	go client.readPump()

	return nil
}
