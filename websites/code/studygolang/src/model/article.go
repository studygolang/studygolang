// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"logger"
	"util"
)

const (
	StatusNew = iota
	StatusOnline
	StatusOffline
)

var LangSlice = []string{"中文", "英文"}
var StatusSlice = []string{"未上线", "已上线", "已下线"}

// 抓取的文章信息
type Article struct {
	Id        int    `json:"id"`
	Domain    string `json:"domain"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Author    string `json:"author"`
	AuthorTxt string `json:"author_txt"`
	Lang      string `json:"lang"`
	PubDate   string `json:"pub_date"`
	Url       string `json:"url"`
	Content   string `json:"content"`
	Txt       string `json:"txt"`
	Tags      string `json:"tags"`
	Status    int    `json:"status"`
	OpUser    string `json:"op_user"`
	Ctime     string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewArticle() *Article {
	return &Article{
		Dao: &Dao{tablename: "articles"},
	}
}

func (this *Article) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Article) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Article) FindAll(selectCol ...string) ([]*Article, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	articleList := make([]*Article, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		article := NewArticle()
		err = this.Scan(rows, colNum, article.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Article FindAll Scan Error:", err)
			continue
		}
		articleList = append(articleList, article)
	}
	return articleList, nil
}

// 为了支持连写
func (this *Article) Where(condition string) *Article {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *Article) Set(clause string, args ...interface{}) *Article {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *Article) Limit(limit string) *Article {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Article) Order(order string) *Article {
	this.Dao.Order(order)
	return this
}

func (this *Article) prepareInsertData() {
	this.columns = []string{"domain", "name", "title", "author", "author_txt", "lang", "pub_date", "url", "content", "txt", "tags"}
	this.colValues = []interface{}{this.Domain, this.Name, this.Title, this.Author, this.AuthorTxt, this.Lang, this.PubDate, this.Url, this.Content, this.Txt, this.Tags}
}

func (this *Article) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         &this.Id,
		"domain":     &this.Domain,
		"name":       &this.Name,
		"title":      &this.Title,
		"cover":      &this.Cover,
		"author":     &this.Author,
		"author_txt": &this.AuthorTxt,
		"lang":       &this.Lang,
		"pub_date":   &this.PubDate,
		"url":        &this.Url,
		"content":    &this.Content,
		"txt":        &this.Txt,
		"tags":       &this.Tags,
		"status":     &this.Status,
		"op_user":    &this.OpUser,
		"ctime":      &this.Ctime,
	}
}

// 抓取网站文章的规则
type CrawlRule struct {
	Id      int    `json:"id"`
	Domain  string `json:"domain"`
	Subpath string `json:"subpath"`
	Lang    string `json:"lang"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	InUrl   bool   `json:"in_url"`
	PubDate string `json:"pub_date"`
	Content string `json:"content"`
	OpUser  string `json:"op_user"`
	Ctime   string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewCrawlRule() *CrawlRule {
	return &CrawlRule{
		Dao: &Dao{tablename: "crawl_rule"},
	}
}

func (this *CrawlRule) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *CrawlRule) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *CrawlRule) FindAll(selectCol ...string) ([]*CrawlRule, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	ruleList := make([]*CrawlRule, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		rule := NewCrawlRule()
		err = this.Scan(rows, colNum, rule.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("CrawlRule FindAll Scan Error:", err)
			continue
		}
		ruleList = append(ruleList, rule)
	}
	return ruleList, nil
}

// 为了支持连写
func (this *CrawlRule) Where(condition string) *CrawlRule {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *CrawlRule) Set(clause string, args ...interface{}) *CrawlRule {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *CrawlRule) Limit(limit string) *CrawlRule {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *CrawlRule) Order(order string) *CrawlRule {
	this.Dao.Order(order)
	return this
}

func (this *CrawlRule) prepareInsertData() {
	this.columns = []string{"domain", "subpath", "lang", "name", "title", "author", "in_url", "pub_date", "content", "op_user"}
	this.colValues = []interface{}{this.Domain, this.Subpath, this.Lang, this.Name, this.Title, this.Author, this.InUrl, this.PubDate, this.Content, this.OpUser}
}

func (this *CrawlRule) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       &this.Id,
		"domain":   &this.Domain,
		"subpath":  &this.Subpath,
		"lang":     &this.Lang,
		"name":     &this.Name,
		"title":    &this.Title,
		"author":   &this.Author,
		"in_url":   &this.InUrl,
		"pub_date": &this.PubDate,
		"content":  &this.Content,
		"op_user":  &this.OpUser,
		"ctime":    &this.Ctime,
	}
}
