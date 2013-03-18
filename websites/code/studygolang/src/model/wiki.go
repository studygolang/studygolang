// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"logger"
	"time"
	"util"
)

// 角色信息
type Wiki struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Uri     string `json:"uri"`
	Uid     int    `json:"uid"`
	Cuid    string `json:"cuid"`
	Ctime   string `json:"ctime"`
	Mtime   string `json:"mtime"`

	// 数据库访问对象
	*Dao
}

func NewWiki() *Wiki {
	return &Wiki{
		Dao: &Dao{tablename: "wiki"},
	}
}

func (this *Wiki) Insert() (int64, error) {
	this.Ctime = time.Now().Format("2006-01-02 15:04:05")
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Wiki) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Wiki) FindAll(selectCol ...string) ([]*Wiki, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	wikiList := make([]*Wiki, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		wiki := NewWiki()
		err = this.Scan(rows, colNum, wiki.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Wiki FindAll Scan Error:", err)
			continue
		}
		wikiList = append(wikiList, wiki)
	}
	return wikiList, nil
}

// 为了支持连写
func (this *Wiki) Where(condition string) *Wiki {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *Wiki) Limit(limit string) *Wiki {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Wiki) Order(order string) *Wiki {
	this.Dao.Order(order)
	return this
}

func (this *Wiki) prepareInsertData() {
	this.columns = []string{"title", "content", "uri", "uid", "ctime"}
	this.colValues = []interface{}{this.Title, this.Content, this.Uri, this.Uid, this.Ctime}
}

func (this *Wiki) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"title":   &this.Title,
		"content": &this.Content,
		"uri":     &this.Uri,
		"uid":     &this.Uid,
		"cuid":    &this.Cuid,
		"ctime":   &this.Ctime,
		"mtime":   &this.Mtime,
	}
}
