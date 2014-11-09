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

const (
	LinkForm    = "只是链接"
	ContentForm = "包括内容"
)

// 资源信息
type Resource struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Form    string `json:"form"`
	Content string `json:"content"`
	Url     string `json:"url"`
	Uid     int    `json:"uid"`
	Catid   int    `json:"catid"`
	Ctime   string `json:"ctime"`
	Mtime   string `json:"mtime"`

	// 数据库访问对象
	*Dao
}

func NewResource() *Resource {
	return &Resource{
		Dao: &Dao{tablename: "resource"},
	}
}

func (this *Resource) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Resource) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Resource) FindAll(selectCol ...string) ([]*Resource, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	resourceList := make([]*Resource, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		resource := NewResource()
		err = this.Scan(rows, colNum, resource.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Resource FindAll Scan Error:", err)
			continue
		}
		resourceList = append(resourceList, resource)
	}
	return resourceList, nil
}

// 为了支持连写
func (this *Resource) Set(clause string, args ...interface{}) *Resource {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *Resource) Where(condition string, args ...interface{}) *Resource {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *Resource) Limit(limit string) *Resource {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Resource) Order(order string) *Resource {
	this.Dao.Order(order)
	return this
}

func (this *Resource) prepareInsertData() {
	this.columns = []string{"title", "form", "content", "url", "uid", "catid", "ctime"}
	this.colValues = []interface{}{this.Title, this.Form, this.Content, this.Url, this.Uid, this.Catid, this.Ctime}
}

func (this *Resource) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"title":   &this.Title,
		"form":    &this.Form,
		"content": &this.Content,
		"url":     &this.Url,
		"uid":     &this.Uid,
		"catid":   &this.Catid,
		"ctime":   &this.Ctime,
		"mtime":   &this.Mtime,
	}
}

// 资源扩展（计数）信息
type ResourceEx struct {
	Id      int    `json:"id"`
	Viewnum int    `json:"viewnum"`
	Cmtnum  int    `json:"cmtnum"`
	Likenum int    `json:"likenum"`
	Mtime   string `json:"mtime"`

	// 数据库访问对象
	*Dao
}

func NewResourceEx() *ResourceEx {
	return &ResourceEx{
		Dao: &Dao{tablename: "resource_ex"},
	}
}

func (this *ResourceEx) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	num, err := result.RowsAffected()
	return int(num), err
}

func (this *ResourceEx) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *ResourceEx) FindAll(selectCol ...string) ([]*ResourceEx, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	resourceExList := make([]*ResourceEx, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		resourceEx := NewResourceEx()
		err = this.Scan(rows, colNum, resourceEx.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("TopicEx FindAll Scan Error:", err)
			continue
		}
		resourceExList = append(resourceExList, resourceEx)
	}
	return resourceExList, nil
}

// 为了支持连写
func (this *ResourceEx) Where(condition string, args ...interface{}) *ResourceEx {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *ResourceEx) Limit(limit string) *ResourceEx {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *ResourceEx) Order(order string) *ResourceEx {
	this.Dao.Order(order)
	return this
}

func (this *ResourceEx) prepareInsertData() {
	this.columns = []string{"id", "viewnum", "cmtnum", "likenum"}
	this.colValues = []interface{}{this.Id, this.Viewnum, this.Cmtnum, this.Likenum}
}

func (this *ResourceEx) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      &this.Id,
		"viewnum": &this.Viewnum,
		"cmtnum":  &this.Cmtnum,
		"likenum": &this.Likenum,
		"mtime":   &this.Mtime,
	}
}

// 资源分类信息
type ResourceCat struct {
	Catid int    `json:"catid"`
	Name  string `json:"name"`
	Intro string `json:"intro"`
	Ctime string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewResourceCat() *ResourceCat {
	return &ResourceCat{
		Dao: &Dao{tablename: "resource_category"},
	}
}

func (this *ResourceCat) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (this *ResourceCat) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *ResourceCat) FindAll(selectCol ...string) ([]*ResourceCat, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	catList := make([]*ResourceCat, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		category := NewResourceCat()
		err = this.Scan(rows, colNum, category.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("TopicNode FindAll Scan Error:", err)
			continue
		}
		catList = append(catList, category)
	}
	return catList, nil
}

func (this *ResourceCat) prepareInsertData() {
	this.columns = []string{"name", "intro"}
	this.colValues = []interface{}{this.Name, this.Intro}
}

func (this *ResourceCat) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"catid": &this.Catid,
		"name":  &this.Name,
		"intro": &this.Intro,
		"ctime": &this.Ctime,
	}
}
