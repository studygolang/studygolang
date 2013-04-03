// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"encoding/json"
	"logger"
	"util"
)

const (
	FdelNotDel = "未删"
	FdelHasDel = "已删"

	TdelNotDel = "未删"
	TdelHasDel = "已删"

	HasRead = "已读"
	NotRead = "未读"
)

// 短消息
type Message struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Hasread string `json:"hasread"`
	From    int    `json:"from"`
	Fdel    string `json:"fdel"`
	To      int    `json:"to"`
	Tdel    string `json:"tdel"`
	Ctime   string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewMessage() *Message {
	return &Message{
		Dao: &Dao{tablename: "message"},
	}
}

func (this *Message) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// 为了方便返回对象本身
func (this *Message) Find(selectCol ...string) (*Message, error) {
	return this, this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Message) FindAll(selectCol ...string) ([]*Message, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	messageList := make([]*Message, 0, 10)
	colNum := len(selectCol)
	for rows.Next() {
		message := NewMessage()
		err = this.Scan(rows, colNum, message.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Message FindAll Scan Error:", err)
			continue
		}
		messageList = append(messageList, message)
	}
	return messageList, nil
}

// 为了支持连写
func (this *Message) Where(condition string) *Message {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *Message) Set(clause string) *Message {
	this.Dao.Set(clause)
	return this
}

// 为了支持连写
func (this *Message) Limit(limit string) *Message {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Message) Order(order string) *Message {
	this.Dao.Order(order)
	return this
}

func (this *Message) prepareInsertData() {
	this.columns = []string{"content", "from", "to"}
	this.colValues = []interface{}{this.Content, this.From, this.To}
}

func (this *Message) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"content": &this.Content,
		"hasread": &this.Hasread,
		"from":    &this.From,
		"fdel":    &this.Fdel,
		"to":      &this.To,
		"tdel":    &this.Tdel,
		"ctime":   &this.Ctime,
	}
}

const (
	// 和comment中objtype保持一致（除了@）
	MsgtypeTopicReply      = iota // 回复我的主题
	MsgtypeBlogComment            // 评论我的博客
	MsgtypeResourceComment        // 评论我的资源
	MsgtypeWikiComment            // 评论我的Wiki页
	MsgtypeAtMe                   // @提到我
)

// 短消息
type SystemMessage struct {
	Id      int    `json:"id"`
	Msgtype int    `json:"msgtype"`
	Hasread string `json:"hasread"`
	To      int    `json:"to"`
	Ctime   string `json:"ctime"`

	// 扩展信息，json格式
	ext string

	// 数据库访问对象
	*Dao
}

func NewSystemMessage() *SystemMessage {
	return &SystemMessage{
		Dao: &Dao{tablename: "system_message"},
	}
}

func (this *SystemMessage) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// 为了方便返回对象本身
func (this *SystemMessage) Find(selectCol ...string) (*SystemMessage, error) {
	return this, this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *SystemMessage) FindAll(selectCol ...string) ([]*SystemMessage, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	messageList := make([]*SystemMessage, 0, 10)
	colNum := len(selectCol)
	for rows.Next() {
		message := NewSystemMessage()
		err = this.Scan(rows, colNum, message.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("SystemMessage FindAll Scan Error:", err)
			continue
		}
		messageList = append(messageList, message)
	}
	return messageList, nil
}

func (this *SystemMessage) Ext() map[string]interface{} {
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(this.ext), &result); err != nil {
		logger.Errorln("SystemMessage Ext JsonUnmarshal Error:", err)
		return nil
	}
	return result
}

func (this *SystemMessage) SetExt(ext map[string]interface{}) {
	if extBytes, err := json.Marshal(ext); err != nil {
		logger.Errorln("SystemMessage SetExt JsonMarshal Error:", err)
	} else {
		this.ext = string(extBytes)
	}
}

// 为了支持连写
func (this *SystemMessage) Set(clause string) *SystemMessage {
	this.Dao.Set(clause)
	return this
}

// 为了支持连写
func (this *SystemMessage) Where(condition string) *SystemMessage {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *SystemMessage) Limit(limit string) *SystemMessage {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *SystemMessage) Order(order string) *SystemMessage {
	this.Dao.Order(order)
	return this
}

func (this *SystemMessage) prepareInsertData() {
	this.columns = []string{"msgtype", "to", "ext"}
	this.colValues = []interface{}{this.Msgtype, this.To, this.ext}
}

func (this *SystemMessage) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"msgtype": &this.Msgtype,
		"hasread": &this.Hasread,
		"to":      &this.To,
		"ext":     &this.ext,
		"ctime":   &this.Ctime,
	}
}
