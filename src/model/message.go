// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"encoding/json"

	"github.com/polaris1119/logger"
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
	Id      int       `json:"id" xorm:"pk autoincr"`
	Content string    `json:"content"`
	Hasread string    `json:"hasread"`
	From    int       `json:"from"`
	Fdel    string    `json:"fdel"`
	To      int       `json:"to"`
	Tdel    string    `json:"tdel"`
	Ctime   OftenTime `json:"ctime" xorm:"created"`
}

const (
	// 和comment中objtype保持一致（除了@）
	MsgtypeTopicReply      = iota // 回复我的主题
	MsgtypeArticleComment         // 评论我的博文
	MsgtypeResourceComment        // 评论我的资源
	MsgtypeWikiComment            // 评论我的Wiki页
	MsgtypeProjectComment         // 评论我的项目

	MsgtypeAtMe        = 10 // 评论 @提到我
	MsgtypePublishAtMe = 11 // 发布时提到我

	MsgtypeSubjectContribute = 12 //专栏投稿
)

// 系统消息
type SystemMessage struct {
	Id      int       `json:"id" xorm:"pk autoincr"`
	Msgtype int       `json:"msgtype"`
	Hasread string    `json:"hasread"`
	To      int       `json:"to"`
	Ctime   OftenTime `json:"ctime" xorm:"created"`

	// 扩展信息，json格式
	Ext string
}

func (this *SystemMessage) GetExt() map[string]interface{} {
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(this.Ext), &result); err != nil {
		logger.Errorln("SystemMessage Ext JsonUnmarshal Error:", err)
		return nil
	}
	return result
}

func (this *SystemMessage) SetExt(ext map[string]interface{}) {
	if extBytes, err := json.Marshal(ext); err != nil {
		logger.Errorln("SystemMessage SetExt JsonMarshal Error:", err)
	} else {
		this.Ext = string(extBytes)
	}
}
