// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

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
		Title       string     `xml:"title"`
		Link        string     `xml:"link"`
		Description string     `xml:"description"`
		Language    string     `xml:"language"`
		Items       []rssItem  `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
}

// RSS 生成 RSS 2.0 格式的 Feed
func (FeedController) RSS(ctx echo.Context) error {
	paginator := logic.NewPaginatorWithPerPage(1, 20)
	articles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "top DESC, ctime DESC", "")

	feed := rssFeed{Version: "2.0"}
	feed.Channel.Title = "Go语言中文网"
	feed.Channel.Link = "https://studygolang.com"
	feed.Channel.Description = "Go语言中文网，打造最棒的Go语言社区"
	feed.Channel.Language = "zh-CN"

	for _, article := range articles {
		description := article.Content
		if description == "" {
			description = article.Txt
		}
		if len(description) > 200 {
			description = description[:200]
		}

		item := rssItem{
			Title:       article.Title,
			Link:        "https://studygolang.com/articles/" + strconv.Itoa(article.Id),
			Description: description,
			PubDate:     time.Time(article.Ctime).Format(time.RFC1123),
		}
		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	ctx.Response().Header().Set("Content-Type", "application/xml; charset=utf-8")
	return xml.NewEncoder(ctx.Response().Writer).Encode(feed)
}
