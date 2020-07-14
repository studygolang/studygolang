// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "encoding/xml"

const (
	WeMsgTypeText       = "text"
	WeMsgTypeImage      = "image"
	WeMsgTypeVoice      = "voice"
	WeMsgTypeVideo      = "video"
	WeMsgTypeShortVideo = "shortvideo"
	WeMsgTypeLocation   = "location"
	WeMsgTypeLink       = "link"
	WeMsgTypeEvent      = "event"

	WeEventSubscribe   = "subscribe"
	WeEventUnsubscribe = "unsubscribe"
)

type WechatMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64

	// 图片消息
	PicUrl  string
	MediaId string

	// 音频消息
	Format string

	// 视频或短视频
	ThumbMediaId string

	// 地理位置消息
	Location_X float64
	Location_Y float64
	Scale      int
	Label      string

	// 链接消息
	Title       string
	Description string
	Url         string

	// 事件
	Event string
}

type CData struct {
	Val string `xml:",cdata"`
}

type WechatReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   *CData
	FromUserName *CData
	CreateTime   int64
	MsgType      *CData
	Content      *CData       `xml:",omitempty"`
	Image        *WechatImage `xml:",omitempty"`
}

type WechatImage struct {
	MediaId *CData
}
