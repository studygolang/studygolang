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

// wordpress文章信息
type Blog struct {
	Id          int    `json:"ID"`
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
	PostStatus  string `json:"post_status"` // 只查=publish的
	PostName    string `json:"post_name"`   // 链接后缀
	PostDate    string `json:"post_date"`

	// 链接
	PostUri string
	// 数据库访问对象
	*Dao
}

func NewBlog() *Blog {
	return &Blog{
		Dao: &Dao{tablename: "go_posts"},
	}
}

func (this *Blog) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Blog) FindAll(selectCol ...string) ([]*Blog, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	blogList := make([]*Blog, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		blog := NewBlog()
		err = this.Scan(rows, colNum, blog.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Blog FindAll Scan Error:", err)
			continue
		}
		blogList = append(blogList, blog)
	}
	return blogList, nil
}

// 为了支持连写
func (this *Blog) Where(condition string, args ...interface{}) *Blog {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *Blog) Limit(limit string) *Blog {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Blog) Order(order string) *Blog {
	this.Dao.Order(order)
	return this
}

func (this *Blog) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":           &this.Id,
		"post_title":   &this.PostTitle,
		"post_content": &this.PostContent,
		"post_status":  &this.PostStatus,
		"post_name":    &this.PostName,
		"post_date":    &this.PostDate,
	}
}
