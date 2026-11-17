// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

type FeedController struct{}

func (self FeedController) RegisterRoute(g *echo.Group) {
	g.GET("/feed", self.RSS)
}

// rssFeed 定义 RSS 2.0 的 XML 结构
type rssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Author      string   `xml:"author"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
}

// RSS 生成 RSS 2.0 格式的 Feed
func (FeedController) RSS(ctx echo.Context) error {
	respBody, err := logic.DefaultSearcher.FindAtomFeeds(50)
	if err != nil {
		return err
	}

	feed := rssFeed{Version: "2.0"}
	feed.Channel.Title = logic.WebsiteSetting.Name
	scheme := "http"
	if logic.WebsiteSetting.OnlyHttps {
		scheme = "https"
	}
	feed.Channel.Link = scheme + "://" + logic.WebsiteSetting.Domain + "/"
	feed.Channel.Description = logic.WebsiteSetting.Slogan
	feed.Channel.Language = "zh-CN"

	for _, doc := range respBody.Docs {
		url := ""
		switch doc.Objtype {
		case model.TypeTopic:
			url = fmt.Sprintf("%stopics/%d", feed.Channel.Link, doc.Objid)
		case model.TypeArticle:
			url = fmt.Sprintf("%sarticles/%d", feed.Channel.Link, doc.Objid)
		case model.TypeResource:
			url = fmt.Sprintf("%sresources/%d", feed.Channel.Link, doc.Objid)
		case model.TypeProject:
			url = fmt.Sprintf("%sp/%d", feed.Channel.Link, doc.Objid)
		case model.TypeWiki:
			url = fmt.Sprintf("%swiki/%d", feed.Channel.Link, doc.Objid)
		case model.TypeBook:
			url = fmt.Sprintf("%sbook/%d", feed.Channel.Link, doc.Objid)
		case model.TypeInterview:
			url = fmt.Sprintf("%sinterview/%d", feed.Channel.Link, doc.Objid)
		}

		item := rssItem{
			Title:       doc.Title,
			Link:        url,
			GUID:        url,
			Author:      doc.Author,
			Description: doc.Content,
			PubDate:     time.Time(doc.CreatedAt).UTC().Format(time.RFC1123Z),
		}
		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	data, err := xml.Marshal(feed)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, echo.MIMEApplicationXMLCharsetUTF8, append([]byte(xml.Header), data...))
}
