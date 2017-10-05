// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"model"
	"net/http"
	"time"

	. "http"

	"github.com/gorilla/feeds"
	"github.com/labstack/echo"
)

type FeedController struct{}

// 注册路由
func (self FeedController) RegisterRoute(g *echo.Group) {
	g.Get("/feed", self.List)
}

func (self FeedController) List(ctx echo.Context) error {
	link := "http://" + logic.WebsiteSetting.Domain
	isHttps := CheckIsHttps(ctx)
	if isHttps {
		link = "https://" + logic.WebsiteSetting.Domain
	}

	now := time.Now()

	feed := &feeds.Feed{
		Title:       logic.WebsiteSetting.Name,
		Link:        &feeds.Link{Href: link},
		Description: logic.WebsiteSetting.Slogan,
		Author:      &feeds.Author{Name: "polaris", Email: "polaris@studygolang.com"},
		Created:     now,
		Updated:     now,
	}

	siteFeeds := logic.DefaultFeed.FindRecent(ctx, 50)

	feed.Items = make([]*feeds.Item, len(siteFeeds))

	for i, siteFeed := range siteFeeds {
		strObjtype := ""
		switch siteFeed.Objtype {
		case model.TypeTopic:
			strObjtype = "主题"
		case model.TypeResource:
			strObjtype = "资源"
		case model.TypeArticle:
			strObjtype = "文章"
		case model.TypeProject:
			strObjtype = "开源项目"
		case model.TypeBook:
			strObjtype = "图书"
		case model.TypeWiki:
			strObjtype = "Wiki"
		}
		feed.Items[i] = &feeds.Item{
			Title:       siteFeed.Title,
			Link:        &feeds.Link{Href: link + siteFeed.Uri},
			Description: "这是" + strObjtype,
			Created:     time.Time(siteFeed.CreatedAt),
			Updated:     time.Time(siteFeed.UpdatedAt),
		}
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}

	return self.responseXML(ctx, atom)
}

func (FeedController) responseXML(ctx echo.Context, data string) (err error) {
	response := ctx.Response()
	response.Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	response.WriteHeader(http.StatusOK)
	_, err = response.Write([]byte(data))
	return
}
