// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"encoding/xml"
	"fmt"
	"model"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/polaris1119/config"
)

type WechatLogic struct{}

var DefaultWechat = WechatLogic{}

func (self WechatLogic) AutoReply(ctx context.Context, reqData []byte) (*model.WechatReply, error) {
	objLog := GetLogger(ctx)

	wechatMsg := &model.WechatMsg{}
	err := xml.Unmarshal(reqData, wechatMsg)
	if err != nil {
		objLog.Errorln("wechat autoreply xml unmarshal error:", err)
		return nil, err
	}

	switch wechatMsg.MsgType {
	case model.WeMsgTypeText:
		if strings.Contains(wechatMsg.Content, "晨读") {
			return self.readingContent(ctx, wechatMsg)
		} else {
			return self.searchContent(ctx, wechatMsg)
		}
	case model.WeMsgTypeEvent:
		switch wechatMsg.Event {
		case model.WeEventSubscribe:
			return self.wechatResponse(ctx, config.ConfigFile.MustValue("wechat", "subscribe"), wechatMsg)
		}
	}

	return self.wechatResponse(ctx, "success", wechatMsg)
}

func (self WechatLogic) readingContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	var formatContent = func(reading *model.MorningReading) string {
		if reading.Inner == 0 {
			return fmt.Sprintf("%s\n%s", reading.Content, reading.Url)
		}

		host := "http://"
		if WebsiteSetting.OnlyHttps {
			host = "https://"
		}
		return fmt.Sprintf("%s\n%s%s/articles/%d", reading.Content, host, WebsiteSetting.Domain, reading.Inner)
	}

	var readings []*model.MorningReading
	if wechatMsg.Content == "最新晨读" {
		readings = DefaultReading.FindBy(ctx, 1, model.RtypeGo)
		if len(readings) == 0 {
			return self.wechatResponse(ctx, "没有找到您要的内容", wechatMsg)
		}

		return self.wechatResponse(ctx, formatContent(readings[0]), wechatMsg)
	}

	readings = DefaultReading.FindBy(ctx, 5, model.RtypeGo)

	respContentSlice := make([]string, len(readings))
	for i, reading := range readings {
		respContentSlice[i] = formatContent(reading)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n\n"), wechatMsg)
}

func (self WechatLogic) searchContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {
	objLog := GetLogger(ctx)

	respBody, err := DefaultSearcher.SearchByField("title", wechatMsg.Content, 0, 5)
	if err != nil {
		objLog.Errorln("wechat search by field error:", err)
		return nil, err
	}

	if respBody.NumFound == 0 {
		return self.wechatResponse(ctx, "没有找到您要的内容", wechatMsg)
	}

	host := WebsiteSetting.Domain
	if WebsiteSetting.OnlyHttps {
		host = "https://" + host
	} else {
		host = "http://" + host
	}

	respContentSlice := make([]string, len(respBody.Docs))
	for i, doc := range respBody.Docs {
		url := ""

		switch doc.Objtype {
		case model.TypeTopic:
			url = fmt.Sprintf("%s/topics/%d", host, doc.Objid)
		case model.TypeArticle:
			url = fmt.Sprintf("%s/articles/%d", host, doc.Objid)
		case model.TypeProject:
			url = fmt.Sprintf("%s/p/%d", host, doc.Objid)
		case model.TypeWiki:
			url = fmt.Sprintf("%s/wiki/%d", host, doc.Objid)
		case model.TypeBook:
			url = fmt.Sprintf("%s/book/%d", host, doc.Objid)
		}
		respContentSlice[i] = fmt.Sprintf("《%s》 %s", doc.Title, url)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) wechatResponse(ctx context.Context, respContent string, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {
	wechatReply := &model.WechatReply{
		ToUserName:   &model.CData{wechatMsg.ToUserName},
		FromUserName: &model.CData{wechatMsg.FromUserName},
		MsgType:      &model.CData{wechatMsg.MsgType},
		CreateTime:   time.Now().Unix(),
	}
	switch wechatMsg.MsgType {
	case model.WeMsgTypeText:
		wechatReply.Content = &model.CData{respContent}
	default:
		wechatReply.Content = &model.CData{"没有找到您要的内容"}
	}

	return wechatReply, nil
}
