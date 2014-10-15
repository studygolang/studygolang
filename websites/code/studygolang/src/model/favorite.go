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

// 用户收藏（用户可以收藏文章、话题、资源等）
type Favorite struct {
	Uid     int    `json:"uid"`
	Objtype int    `json:"objtype"`
	Objid   int    `json:"objid"`
	Ctime   string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewFavorite() *Favorite {
	return &Favorite{
		Dao: &Dao{tablename: "favorites"},
	}
}

func (this *Favorite) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Favorite) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Favorite) FindAll(selectCol ...string) ([]*Favorite, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	favoriteList := make([]*Favorite, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		favorite := NewFavorite()
		err = this.Scan(rows, colNum, favorite.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Favorite FindAll Scan Error:", err)
			continue
		}
		favoriteList = append(favoriteList, favorite)
	}
	return favoriteList, nil
}

// 为了支持连写
func (this *Favorite) Where(condition string, args ...interface{}) *Favorite {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *Favorite) Set(clause string, args ...interface{}) *Favorite {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *Favorite) Limit(limit string) *Favorite {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Favorite) Order(order string) *Favorite {
	this.Dao.Order(order)
	return this
}

func (this *Favorite) prepareInsertData() {
	this.columns = []string{"uid", "objtype", "objid"}
	this.colValues = []interface{}{this.Uid, this.Objtype, this.Objid}
}

func (this *Favorite) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"uid":     &this.Uid,
		"objtype": &this.Objtype,
		"objid":   &this.Objid,
		"ctime":   &this.Ctime,
	}
}
