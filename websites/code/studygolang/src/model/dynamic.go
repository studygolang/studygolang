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

// 动态（go动态；本站动态等）
type Dynamic struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Dmtype  int    `json:"dmtype"`
	Url     string `json:"url"`
	Seq     int    `json:"seq"`
	Ctime   string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewDynamic() *Dynamic {
	return &Dynamic{
		Dao: &Dao{tablename: "dynamic"},
	}
}

func (this *Dynamic) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Dynamic) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Dynamic) FindAll(selectCol ...string) ([]*Dynamic, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	dynamicList := make([]*Dynamic, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		dynamic := NewDynamic()
		err = this.Scan(rows, colNum, dynamic.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Dynamic FindAll Scan Error:", err)
			continue
		}
		dynamicList = append(dynamicList, dynamic)
	}
	return dynamicList, nil
}

// 为了支持连写
func (this *Dynamic) Where(condition string) *Dynamic {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *Dynamic) Set(clause string, args ...interface{}) *Dynamic {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *Dynamic) Limit(limit string) *Dynamic {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Dynamic) Order(order string) *Dynamic {
	this.Dao.Order(order)
	return this
}

func (this *Dynamic) prepareInsertData() {
	this.columns = []string{"content", "dmtype", "url", "seq"}
	this.colValues = []interface{}{this.Content, this.Dmtype, this.Url, this.Seq}
}

func (this *Dynamic) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"content": &this.Content,
		"dmtype":  &this.Dmtype,
		"url":     &this.Url,
		"seq":     &this.Seq,
		"ctime":   &this.Ctime,
	}
}
