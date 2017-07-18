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

	"golang.org/x/net/context"

	"github.com/polaris1119/config"
)

type WechatLogic struct{}

var DefaultWechat = WechatLogic{}

func (self WechatLogic) AutoReply(ctx context.Context, reqData []byte) (string, error) {
	objLog := GetLogger(ctx)

	wechatMsg := &model.WechatMsg{}
	err := xml.Unmarshal(reqData, wechatMsg)
	if err != nil {
		objLog.Errorln("wechat autoreply xml unmarshal error:", err)
		return "", err
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
			return config.ConfigFile.MustValue("wechat", "subscribe"), nil
		}
	}

	return "", nil
}

func (self WechatLogic) readingContent(ctx context.Context, wechatMsg *model.WechatMsg) (string, error) {

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
			return "没有找到您要的内容", nil
		}

		return formatContent(readings[0]), nil
	}

	readings = DefaultReading.FindBy(ctx, 5, model.RtypeGo)

	respContentSlice := make([]string, len(readings))
	for i, reading := range readings {
		respContentSlice[i] = formatContent(reading)
	}

	return strings.Join(respContentSlice, "\n\n"), nil
}

func (self WechatLogic) searchContent(ctx context.Context, wechatMsg *model.WechatMsg) (string, error) {
	objLog := GetLogger(ctx)

	respBody, err := DefaultSearcher.SearchByField("title", wechatMsg.Content, 0, 5)
	if err != nil {
		objLog.Errorln("wechat search by field error:", err)
		return "", err
	}

	if respBody.NumFound == 0 {
		return "没有找到您要的内容", nil
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

	return strings.Join(respContentSlice, "\n"), nil
}
