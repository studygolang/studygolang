// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

// 不要修改常量的顺序
const (
	TypeTopic    = iota // 主题
	TypeArticle         // 博文
	TypeResource        // 资源
	TypeWiki            // WIKI
	TypeProject         // 开源项目
	TypeBook            // 图书
)

const (
	TopicURI    = "topics"
	ArticleURI  = "articles"
	ResourceURI = "resources"
	WikiURI     = "wiki"
	ProjectURI  = "p"
	BookURI     = "book"
)

var PathUrlMap = map[int]string{
	TypeTopic:    "/topics/",
	TypeArticle:  "/articles/",
	TypeResource: "/resources/",
	TypeWiki:     "/wiki/",
	TypeProject:  "/p/",
	TypeBook:     "/book/",
}

var TypeNameMap = map[int]string{
	TypeTopic:    "主题",
	TypeArticle:  "博文",
	TypeResource: "资源",
	TypeWiki:     "Wiki",
	TypeProject:  "项目",
	TypeBook:     "图书",
}

// 评论信息（通用）
type Comment struct {
	Cid     int       `json:"cid" xorm:"pk autoincr"`
	Objid   int       `json:"objid"`
	Objtype int       `json:"objtype"`
	Content string    `json:"content"`
	Uid     int       `json:"uid"`
	Floor   int       `json:"floor"`
	Flag    int       `json:"flag"`
	Ctime   OftenTime `json:"ctime" xorm:"created"`

	Objinfo    map[string]interface{} `json:"objinfo" xorm:"-"`
	ReplyFloor int                    `json:"reply_floor" xorm:"-"` // 回复某一楼层
}

func (*Comment) TableName() string {
	return "comments"
}
