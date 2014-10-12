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

// 搜索词统计
type SearchStat struct {
	Id      int    `json:"id"`
	Keyword string `json:"keyword"`
	Times   int    `json:"times"`
	Ctime   string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewSearchStat() *SearchStat {
	return &SearchStat{
		Dao: &Dao{tablename: "search_stat"},
	}
}

func (this *SearchStat) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *SearchStat) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *SearchStat) FindAll(selectCol ...string) ([]*SearchStat, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	searchStatList := make([]*SearchStat, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		searchStat := NewSearchStat()
		err = this.Scan(rows, colNum, searchStat.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("SearchStat FindAll Scan Error:", err)
			continue
		}
		searchStatList = append(searchStatList, searchStat)
	}
	return searchStatList, nil
}

// 为了支持连写
func (this *SearchStat) Where(condition string, args ...interface{}) *SearchStat {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *SearchStat) Increment(field string, num int) *SearchStat {
	this.Dao.Increment(field, num)
	return this
}

// 为了支持连写
func (this *SearchStat) Set(clause string, args ...interface{}) *SearchStat {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *SearchStat) Limit(limit string) *SearchStat {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *SearchStat) Order(order string) *SearchStat {
	this.Dao.Order(order)
	return this
}

func (this *SearchStat) prepareInsertData() {
	this.columns = []string{"keyword", "times"}
	this.colValues = []interface{}{this.Keyword, this.Times}
}

func (this *SearchStat) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"keyword": &this.Keyword,
		"times":   &this.Times,
		"ctime":   &this.Ctime,
	}
}
