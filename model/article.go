// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/polaris1119/logger"
)

const (
	ArticleStatusNew = iota
	ArticleStatusOnline
	ArticleStatusOffline
)

var LangSlice = []string{"中文", "英文"}
var ArticleStatusSlice = []string{"未上线", "已上线", "已下线"}

// 抓取的文章信息
type Article struct {
	Id            int       `json:"id" xorm:"pk autoincr"`
	Domain        string    `json:"domain"`
	Name          string    `json:"name"`
	Title         string    `json:"title"`
	Cover         string    `json:"cover"`
	Author        string    `json:"author"`
	AuthorTxt     string    `json:"author_txt"`
	Lang          int       `json:"lang"`
	PubDate       string    `json:"pub_date"`
	Url           string    `json:"url"`
	Content       string    `json:"content"`
	Txt           string    `json:"txt"`
	Tags          string    `json:"tags"`
	Css           string    `json:"css"`
	Viewnum       int       `json:"viewnum"`
	Cmtnum        int       `json:"cmtnum"`
	Likenum       int       `json:"likenum"`
	Lastreplyuid  int       `json:"lastreplyuid"`
	Lastreplytime OftenTime `json:"lastreplytime"`
	Top           uint8     `json:"top"`
	Markdown      bool      `json:"markdown"`
	GCTT          bool      `json:"gctt" xorm:"gctt"`
	Status        int       `json:"status"`
	OpUser        string    `json:"op_user"`
	Ctime         OftenTime `json:"ctime" xorm:"created"`
	Mtime         OftenTime `json:"mtime" xorm:"<-"`

	IsSelf bool  `json:"is_self" xorm:"-"`
	User   *User `json:"-" xorm:"-"`
	// 排行榜阅读量
	RankView      int   `json:"rank_view" xorm:"-"`
	LastReplyUser *User `json:"last_reply_user" xorm:"-"`
}

func (this *Article) AfterSet(name string, cell xorm.Cell) {
	if name == "id" {
		this.IsSelf = strconv.Itoa(this.Id) == this.Url
	}
}

func (this *Article) BeforeInsert() {
	if this.Tags == "" {
		this.Tags = AutoTag(this.Title, this.Txt, 4)
	}
	this.Lastreplytime = NewOftenTime()
}

func (this *Article) AfterInsert() {
	go func() {
		// AfterInsert 时，自增 ID 还未赋值，这里 sleep 一会，确保自增 ID 有值
		for {
			if this.Id > 0 {
				PublishFeed(this, nil)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (*Article) TableName() string {
	return "articles"
}

type ArticleGCTT struct {
	ArticleID  int `xorm:"article_id pk"`
	Author     string
	AuthorURL  string `xorm:"author_url"`
	Translator string
	Checker    string
	URL        string `xorm:"url"`

	Avatar   string   `xorm:"-"`
	Checkers []string `xorm:"-"`
}

func (*ArticleGCTT) TableName() string {
	return "article_gctt"
}

func (this *ArticleGCTT) AfterSet(name string, cell xorm.Cell) {
	if name == "checker" {
		this.Checkers = strings.Split(this.Checker, ",")
	}
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
	Ext     string `json:"ext"`
	OpUser  string `json:"op_user"`
	Ctime   string `json:"ctime" xorm:"<-"`
}

func (this *CrawlRule) ParseExt() map[string]string {
	if this.Ext == "" {
		return nil
	}

	extMap := make(map[string]string)
	err := json.Unmarshal([]byte(this.Ext), &extMap)
	if err != nil {
		logger.Errorln("parse crawl rule ext error:", err)
		return nil
	}

	return extMap
}

const (
	AutoCrawlOn = 0
	AutoCrawOff = 1
)

// 网站自动抓取规则
type AutoCrawlRule struct {
	Id             int    `json:"id" xorm:"pk autoincr"`
	Website        string `json:"website"`
	AllUrl         string `json:"all_url"`
	IncrUrl        string `json:"incr_url"`
	Keywords       string `json:"keywords"`
	ListSelector   string `json:"list_selector"`
	ResultSelector string `json:"result_selector"`
	PageField      string `json:"page_field"`
	MaxPage        int    `json:"max_page"`
	Ext            string `json:"ext"`
	OpUser         string `json:"op_user"`
	Mtime          string `json:"mtime" xorm:"<-"`

	ExtMap map[string]string `json:"-" xorm:"-"`
}

func (this *AutoCrawlRule) AfterSet(name string, cell xorm.Cell) {
	if name == "ext" {
		if this.Ext == "" {
			return
		}

		this.ExtMap = make(map[string]string)
		err := json.Unmarshal([]byte(this.Ext), &this.ExtMap)
		if err != nil {
			logger.Errorln("parse auto crawl rule ext error:", err)
			return
		}
	}
}
