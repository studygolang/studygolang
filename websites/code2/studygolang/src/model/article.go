// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

const (
	ArticleStatusNew = iota
	ArticleStatusOnline
	ArticleStatusOffline
)

var LangSlice = []string{"中文", "英文"}
var ArticleStatusSlice = []string{"未上线", "已上线", "已下线"}

// 抓取的文章信息
type Article struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Domain    string    `json:"domain"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Author    string    `json:"author"`
	AuthorTxt string    `json:"author_txt"`
	Lang      int       `json:"lang"`
	PubDate   string    `json:"pub_date"`
	Url       string    `json:"url"`
	Content   string    `json:"content"`
	Txt       string    `json:"txt"`
	Tags      string    `json:"tags"`
	Css       string    `json:"css"`
	Viewnum   int       `json:"viewnum"`
	Cmtnum    int       `json:"cmtnum"`
	Likenum   int       `json:"likenum"`
	Top       uint8     `json:"top"`
	Status    int       `json:"status"`
	OpUser    string    `json:"op_user"`
	Ctime     OftenTime `json:"ctime" xorm:"created"`
	Mtime     OftenTime `json:"mtime" xorm:"<-"`
}

func (*Article) TableName() string {
	return "articles"
}

// 抓取网站文章的规则
type CrawlRule struct {
	Id      int    `json:"id" xorm:"pk autoincr"`
	Domain  string `json:"domain"`
	Subpath string `json:"subpath"`
	Lang    int    `json:"lang"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	InUrl   bool   `json:"in_url"`
	PubDate string `json:"pub_date"`
	Content string `json:"content"`
	OpUser  string `json:"op_user"`
	Ctime   string `json:"ctime" xorm:"<-"`
}
