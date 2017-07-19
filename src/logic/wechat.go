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
		} else if strings.Contains(wechatMsg.Content, "主题") || strings.Contains(wechatMsg.Content, "帖子") {
			return self.topicContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "文章") {
			return self.articleContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "资源") {
			return self.resourceContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "项目") {
			return self.projectContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "图书") || strings.Contains(wechatMsg.Content, "book") {
			return self.bookContent(ctx, wechatMsg)
		} else {
			return self.searchContent(ctx, wechatMsg)
		}
	case model.WeMsgTypeEvent:
		switch wechatMsg.Event {
		case model.WeEventSubscribe:
			wechatMsg.MsgType = model.WeMsgTypeText
			return self.wechatResponse(ctx, config.ConfigFile.MustValue("wechat", "subscribe"), wechatMsg)
		}
	}

	return self.wechatResponse(ctx, "success", wechatMsg)
}

func (self WechatLogic) topicContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	topics := DefaultTopic.FindRecent(5)

	respContentSlice := make([]string, len(topics))
	for i, topic := range topics {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/topics/%d", i+1, topic.Title, website(), topic.Tid)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) articleContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	articles := DefaultArticle.FindBy(ctx, 5)

	respContentSlice := make([]string, len(articles))
	for i, article := range articles {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/articles/%d", i+1, article.Title, website(), article.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) resourceContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	resources := DefaultResource.FindBy(ctx, 5)

	respContentSlice := make([]string, len(resources))
	for i, resource := range resources {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/resources/%d", i+1, resource.Title, website(), resource.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) projectContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	projects := DefaultProject.FindBy(ctx, 5)

	respContentSlice := make([]string, len(projects))
	for i, project := range projects {
		respContentSlice[i] = fmt.Sprintf("%d.《%s%s》 %s/p/%d", i+1, project.Category, project.Name, website(), project.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) bookContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	books := DefaultGoBook.FindBy(ctx, 5)

	respContentSlice := make([]string, len(books))
	for i, book := range books {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/book/%d", i+1, book.Name, website(), book.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) readingContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	var formatContent = func(reading *model.MorningReading) string {
		if reading.Inner == 0 {
			return fmt.Sprintf("%s\n%s", reading.Content, reading.Url)
		}

		return fmt.Sprintf("%s\n%s/articles/%d", reading.Content, website(), reading.Inner)
	}

	var readings []*model.MorningReading
	if wechatMsg.Content == "最新晨读" {
		readings = DefaultReading.FindBy(ctx, 1, model.RtypeGo)
		if len(readings) == 0 {
			return self.wechatResponse(ctx, config.ConfigFile.MustValue("wechat", "not_found"), wechatMsg)
		}

		return self.wechatResponse(ctx, formatContent(readings[0]), wechatMsg)
	}

	readings = DefaultReading.FindBy(ctx, 5, model.RtypeGo)

	respContentSlice := make([]string, len(readings))
	for i, reading := range readings {
		respContentSlice[i] = fmt.Sprintf("%d. %s", i+1, formatContent(reading))
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
		return self.wechatResponse(ctx, config.ConfigFile.MustValue("wechat", "not_found"), wechatMsg)
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
		case model.TypeResource:
			url = fmt.Sprintf("%s/resources/%d", host, doc.Objid)
		case model.TypeProject:
			url = fmt.Sprintf("%s/p/%d", host, doc.Objid)
		case model.TypeWiki:
			url = fmt.Sprintf("%s/wiki/%d", host, doc.Objid)
		case model.TypeBook:
			url = fmt.Sprintf("%s/book/%d", host, doc.Objid)
		}
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s", i+1, doc.Title, url)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) wechatResponse(ctx context.Context, respContent string, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {
	wechatReply := &model.WechatReply{
		ToUserName:   &model.CData{Val: wechatMsg.FromUserName},
		FromUserName: &model.CData{Val: wechatMsg.ToUserName},
		MsgType:      &model.CData{Val: wechatMsg.MsgType},
		CreateTime:   time.Now().Unix(),
	}
	switch wechatMsg.MsgType {
	case model.WeMsgTypeText:
		wechatReply.Content = &model.CData{Val: respContent}
	default:
		wechatReply.Content = &model.CData{Val: config.ConfigFile.MustValue("wechat", "not_found")}
	}

	return wechatReply, nil
}
