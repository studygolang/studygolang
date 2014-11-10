// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"config"
	"logger"
	"util"
)

const (
	WsMsgNotify = iota // 通知消息
	WsMsgOnline        // 发送在线用户数（和需要时也发历史最高）
)

type Message struct {
	Type int         `json:"type"`
	Body interface{} `json:"body"`
}

func NewMessage(msgType int, msgBody interface{}) *Message {
	return &Message{
		Type: msgType,
		Body: msgBody,
	}
}

type UserData struct {
	// 该用户收到的消息（key为serverId）
	serverMsgQueue map[int]chan *Message
	lastAccessTime time.Time
	onlineDuartion time.Duration

	rwMutex sync.RWMutex
}

func (this *UserData) Len() int {
	this.rwMutex.RLock()
	defer this.rwMutex.RUnlock()
	return len(this.serverMsgQueue)
}

func (this *UserData) MessageQueue(serverId int) chan *Message {
	return this.serverMsgQueue[serverId]
}

var Book = &book{users: make(map[int]*UserData)}

type book struct {
	users   map[int]*UserData
	rwMutex sync.RWMutex
}

// 增加一个用户到book中（有可能是用户的另一个请求）
// user为UID或IP地址的int表示
func (this *book) AddUser(user, serverId int) *UserData {
	var userData *UserData
	var ok bool
	this.rwMutex.Lock()
	defer this.rwMutex.Unlock()
	if userData, ok = this.users[user]; ok {
		userData.serverMsgQueue[serverId] = make(chan *Message, 1)
		userData.onlineDuartion += time.Now().Sub(userData.lastAccessTime)
		userData.lastAccessTime = time.Now()
	} else {
		userData = &UserData{
			serverMsgQueue: map[int]chan *Message{serverId: make(chan *Message, 1)},
			lastAccessTime: time.Now(),
		}
		this.users[user] = userData

		onlineInfo := map[string]int{"online": len(this.users) + 50}
		// 在线人数超过历史最高
		if len(this.users) > MaxOnlineNum() {
			maxRwMu.Lock()
			maxOnlineNum = len(this.users)
			onlineInfo["maxonline"] = maxOnlineNum
			maxRwMu.Unlock()
			saveMaxOnlineNum()
		}
		// 广播给其他人：有新用户进来，包括可能的新历史最高
		message := NewMessage(WsMsgOnline, onlineInfo)
		go this.BroadcastToOthersMessage(message, user)
	}

	return userData
}

// 删除用户
func (this *book) DelUser(user, serverId int) {
	// 自己只有一个页面建立websocket连接
	if this.users[user].Len() == 1 {
		delete(this.users, user)
	} else {
		delete(this.users[user].serverMsgQueue, serverId)
	}
}

// 判断用户是否还在线
func (this *book) UserIsOnline(user int) bool {
	this.rwMutex.RLock()
	defer this.rwMutex.RUnlock()
	if _, ok := this.users[user]; ok {
		return true
	}
	return false
}

// 在线用户数
func (this *book) Len() int {
	this.rwMutex.RLock()
	defer this.rwMutex.RUnlock()
	return len(this.users)
}

// 给某个用户发送一条消息
func (this *book) PostMessage(uid int, message *Message) {
	if userData, ok := this.users[uid]; ok {
		logger.Infoln("post message to", uid, message)
		for _, messageQueue := range userData.serverMsgQueue {
			messageQueue <- message
		}
	}
}

// 给所有用户广播消息
func (this *book) BroadcastAllUsersMessage(message *Message) {
	logger.Infoln("BroadcastAllUsersMessage message", message)
	for _, userData := range this.users {
		for _, messageQueue := range userData.serverMsgQueue {
			messageQueue <- message
		}
	}
}

// 给除了自己的其他用户广播消息
func (this *book) BroadcastToOthersMessage(message *Message, myself int) {
	logger.Infoln("BroadcastToOthersMessage message", message)
	for uid, userData := range this.users {
		if uid == myself {
			continue
		}
		for _, messageQueue := range userData.serverMsgQueue {
			messageQueue <- message
		}
	}
}

var (
	// 保存历史最大在线用户数
	maxOnlineNum int
	maxRwMu      sync.RWMutex
)

func initMaxOnlineNum() {
	maxRwMu.Lock()
	defer maxRwMu.Unlock()
	if maxOnlineNum == 0 {
		data, err := ioutil.ReadFile(getDataFile())
		if err != nil {
			logger.Errorln("read data file error:", err)
			return
		}
		maxOnlineNum = util.MustInt(strings.TrimSpace(string(data)))
	}
}

// 获得历史最高在线人数
func MaxOnlineNum() int {
	initMaxOnlineNum()
	maxRwMu.RLock()
	defer maxRwMu.RUnlock()
	return maxOnlineNum
}

func saveMaxOnlineNum() {
	data := []byte(strconv.Itoa(MaxOnlineNum()))
	err := ioutil.WriteFile(getDataFile(), data, 0777)
	if err != nil {
		logger.Errorln("write data file error:", err)
		return
	}
}

var dataFile string

func getDataFile() string {
	if dataFile != "" {
		return dataFile
	}
	dataFile = config.Config["data"]
	if !filepath.IsAbs(dataFile) {
		dataFile = config.ROOT + "/" + dataFile
	}
	// 文件夹不存在，则创建
	dataPath := filepath.Dir(dataFile)
	if err := os.MkdirAll(dataPath, 0777); err != nil {
		logger.Errorln("MkdirAll error:", err)
	}
	return dataFile
}
