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

// 权限信息
type Authority struct {
	Aid    int    `json:"aid" pk:"1"`
	Name   string `json:"name"`
	Menu1  int    `json:"menu1"`
	Menu2  int    `json:"menu2"`
	Route  string `json:"route"`
	OpUser string `json:"op_user"`
	Ctime  string `json:"ctime,omitempty"`
	Mtime  string `json:"mtime,omitempty"`

	// 内嵌
	*Dao
}

func NewAuthority() *Authority {
	return &Authority{
		Dao: &Dao{tablename: "authority"},
	}
}

func (this *Authority) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Authority) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Authority) FindAll(selectCol ...string) ([]*Authority, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	authorityList := make([]*Authority, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		authority := NewAuthority()
		err = this.Scan(rows, colNum, authority.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Authority FindAll Scan Error:", err)
			continue
		}
		authorityList = append(authorityList, authority)
	}
	return authorityList, nil
}

// 为了支持连写
func (this *Authority) Where(condition string) *Authority {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *Authority) Limit(limit string) *Authority {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Authority) Set(clause string, args ...interface{}) *Authority {
	this.Dao.Set(clause, args...)
	return this
}

func (this *Authority) prepareInsertData() {
	this.columns = []string{"name", "menu1", "menu2", "route", "op_user", "ctime"}
	this.colValues = []interface{}{this.Name, this.Menu1, this.Menu2, this.Route, this.OpUser, this.Ctime}
}

func (this *Authority) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"aid":     &this.Aid,
		"name":    &this.Name,
		"menu1":   &this.Menu1,
		"menu2":   &this.Menu2,
		"route":   &this.Route,
		"op_user": &this.OpUser,
		"ctime":   &this.Ctime,
		"mtime":   &this.Mtime,
	}
}
