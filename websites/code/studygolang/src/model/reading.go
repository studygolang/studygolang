// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"strings"

	"logger"
	"util"
)

const (
	RtypeGo   = iota // Go技术晨读
	RtypeComp        // 综合技术晨读
)

// Go技术晨读
type MorningReading struct {
	Id       int    `json:"id" pk:"1"`
	Content  string `json:"content"`
	Rtype    int    `json:"rtype"`
	Inner    int    `json:"inner"`
	Url      string `json:"url"`
	Moreurls string `json:"moreurls"`
	Username string `json:"username"`
	Clicknum int    `json:"clicknum,omitempty"`
	Ctime    string `json:"ctime,omitempty"`

	// 晨读日期，从 ctime 中提取
	Rdate string `json:"rdate"`

	Urls []string `json:"urls"`

	// 数据库访问对象
	*Dao
}

func NewMorningReading() *MorningReading {
	return &MorningReading{
		Dao: &Dao{tablename: "morning_reading"},
	}
}

func (this *MorningReading) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *MorningReading) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *MorningReading) FindAll(selectCol ...string) ([]*MorningReading, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	readingList := make([]*MorningReading, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		reading := NewMorningReading()
		err = this.Scan(rows, colNum, reading.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("MorningReading FindAll Scan Error:", err)
			continue
		}
		reading.Rdate = reading.Ctime[:10]
		if reading.Moreurls != "" {
			reading.Urls = strings.Split(reading.Moreurls, ",")
		}
		readingList = append(readingList, reading)
	}
	return readingList, nil
}

// 为了支持连写
func (this *MorningReading) Where(condition string, args ...interface{}) *MorningReading {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *MorningReading) Set(clause string, args ...interface{}) *MorningReading {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *MorningReading) Limit(limit string) *MorningReading {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *MorningReading) Order(order string) *MorningReading {
	this.Dao.Order(order)
	return this
}

func (this *MorningReading) prepareInsertData() {
	this.columns = []string{"content", "rtype", "inner", "url", "moreurls", "username"}
	this.colValues = []interface{}{this.Content, this.Rtype, this.Inner, this.Url, this.Moreurls, this.Username}
}

func (this *MorningReading) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       &this.Id,
		"content":  &this.Content,
		"rtype":    &this.Rtype,
		"inner":    &this.Inner,
		"url":      &this.Url,
		"moreurls": &this.Moreurls,
		"clicknum": &this.Clicknum,
		"username": &this.Username,
		"ctime":    &this.Ctime,
	}
}
