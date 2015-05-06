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

// 代码片段分享
type CodeShare struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Remark  string `json:"remark"`
	Code    string `json:"code"`
	Viewnum int    `json:"viewnum"`
	Cmtnum  int    `json:"cmtnum"`
	Likenum int    `json:"likenum"`
	OpUser  string `json:"op_user"`
	Ctime   string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewCodeShare() *CodeShare {
	return &CodeShare{
		Dao: &Dao{tablename: "code_share"},
	}
}

func (this *CodeShare) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *CodeShare) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *CodeShare) FindAll(selectCol ...string) ([]*CodeShare, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	codeShareList := make([]*CodeShare, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		codeShare := NewCodeShare()
		err = this.Scan(rows, colNum, codeShare.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("CodeShare FindAll Scan Error:", err)
			continue
		}
		codeShareList = append(codeShareList, codeShare)
	}
	return codeShareList, nil
}

// 为了支持连写
func (this *CodeShare) Where(condition string, args ...interface{}) *CodeShare {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *CodeShare) Set(clause string, args ...interface{}) *CodeShare {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *CodeShare) Limit(limit string) *CodeShare {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *CodeShare) Order(order string) *CodeShare {
	this.Dao.Order(order)
	return this
}

func (this *CodeShare) prepareInsertData() {
	this.columns = []string{"title", "remark", "code", "op_user"}
	this.colValues = []interface{}{this.Title, this.Remark, this.Code, this.OpUser}
}

func (this *CodeShare) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"title":   &this.Title,
		"remark":  &this.Remark,
		"code":    &this.Code,
		"viewnum": &this.Viewnum,
		"cmtnum":  &this.Cmtnum,
		"likenum": &this.Likenum,
		"op_user": &this.OpUser,
		"ctime":   &this.Ctime,
	}
}
