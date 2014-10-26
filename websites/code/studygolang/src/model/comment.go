// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"logger"
	"util"
)

// 不要修改常量的顺序
const (
	TYPE_TOPIC    = iota // 帖子
	TYPE_ARTICLE         // 博文
	TYPE_RESOURCE        // 资源
	TYPE_WIKI            // WIKI
	TYPE_PROJECT         // 开源项目
)

var PathUrlMap = map[int]string{
	TYPE_TOPIC:    "/topics/",
	TYPE_ARTICLE:  "/articles/",
	TYPE_RESOURCE: "/resources/",
	TYPE_WIKI:     "/wiki/",
	TYPE_PROJECT:  "/p/",
}

// 评论信息（通用）
type Comment struct {
	Cid     int    `json:"cid"`
	Objid   int    `json:"objid"`
	Objtype int    `json:"objtype"`
	Content string `json:"content"`
	Uid     int    `json:"uid"`
	Floor   int    `json:"floor"`
	Flag    int    `json:"flag"`
	Ctime   string `json:"ctime"`

	Objinfo map[string]interface{} `json:"objinfo"`

	// 数据库访问对象
	*Dao
}

func NewComment() *Comment {
	return &Comment{
		Dao: &Dao{tablename: "comments"},
	}
}

func (this *Comment) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// 为了方便返回对象本身
func (this *Comment) Find(selectCol ...string) (*Comment, error) {
	return this, this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Comment) FindAll(selectCol ...string) ([]*Comment, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	commentList := make([]*Comment, 0, 10)
	colNum := len(selectCol)
	for rows.Next() {
		comment := NewComment()
		err = this.Scan(rows, colNum, comment.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Comment FindAll Scan Error:", err)
			continue
		}
		commentList = append(commentList, comment)
	}
	return commentList, nil
}

// 为了支持连写
func (this *Comment) Where(condition string, args ...interface{}) *Comment {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *Comment) Set(clause string, args ...interface{}) *Comment {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *Comment) Limit(limit string) *Comment {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Comment) Order(order string) *Comment {
	this.Dao.Order(order)
	return this
}

func (this *Comment) prepareInsertData() {
	this.columns = []string{"objid", "content", "objtype", "uid", "floor"}
	this.colValues = []interface{}{this.Objid, this.Content, this.Objtype, this.Uid, this.Floor}
}

func (this *Comment) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"cid":     &this.Cid,
		"objid":   &this.Objid,
		"content": &this.Content,
		"objtype": &this.Objtype,
		"uid":     &this.Uid,
		"flag":    &this.Flag,
		"floor":   &this.Floor,
		"ctime":   &this.Ctime,
	}
}
