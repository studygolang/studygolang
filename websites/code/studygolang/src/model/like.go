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
	FLAG_CANCEL = iota
	FLAG_LIKE   // 喜欢
	FLAG_UNLIKE // 不喜欢（暂时不支持）
)

// 评论信息（通用）
type Like struct {
	Uid     int    `json:"uid"`
	Objid   int    `json:"objid"`
	Objtype int    `json:"objtype"`
	Flag    int    `json:"flag"`
	Ctime   string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewLike() *Like {
	return &Like{
		Dao: &Dao{tablename: "Likes"},
	}
}

func (this *Like) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Like) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Like) FindAll(selectCol ...string) ([]*Like, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	likeList := make([]*Like, 0, 10)
	colNum := len(selectCol)
	for rows.Next() {
		like := NewLike()
		err = this.Scan(rows, colNum, like.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Like FindAll Scan Error:", err)
			continue
		}
		likeList = append(likeList, like)
	}
	return likeList, nil
}

// 为了支持连写
func (this *Like) Where(condition string, args ...interface{}) *Like {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *Like) Set(clause string, args ...interface{}) *Like {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *Like) Limit(limit string) *Like {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Like) Order(order string) *Like {
	this.Dao.Order(order)
	return this
}

func (this *Like) prepareInsertData() {
	this.columns = []string{"objid", "objtype", "uid", "flag"}
	this.colValues = []interface{}{this.Objid, this.Objtype, this.Uid, this.Flag}
}

func (this *Like) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"uid":     &this.Uid,
		"objid":   &this.Objid,
		"objtype": &this.Objtype,
		"flag":    &this.Flag,
		"ctime":   &this.Ctime,
	}
}
