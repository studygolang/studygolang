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

// 开源项目信息
type OpenProject struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Uri      string `json:"uri"`
	Home     string `json:"home"`
	Doc      string `json:"doc"`
	Download string `json:"download"`
	Src      string `json:"src"`
	Logo     string `json:"logo"`
	Desc     string `json:"desc"`
	Repo     string `json:"repo"`
	Author   string `json:"author"`
	Licence  string `json:"licence"`
	Lang     string `json:"lang"`
	Os       string `json:"os"`
	Tags     string `json:"tags"`
	Username string `json:"username"`
	Viewnum  int    `json:"viewnum"`
	Cmtnum   int    `json:"cmtnum"`
	Likenum  int    `json:"likenum"`
	Status   int    `json:"status"`
	Ctime    string `json:"ctime"`
	Mtime    string `json:"mtime"`

	// 数据库访问对象
	*Dao
}

func NewOpenProject() *OpenProject {
	return &OpenProject{
		Dao: &Dao{tablename: "open_project"},
	}
}

func (this *OpenProject) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *OpenProject) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *OpenProject) FindAll(selectCol ...string) ([]*OpenProject, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	projectList := make([]*OpenProject, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		project := NewOpenProject()
		err = this.Scan(rows, colNum, project.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("OpenProject FindAll Scan Error:", err)
			continue
		}
		projectList = append(projectList, project)
	}
	return projectList, nil
}

// 为了支持连写
func (this *OpenProject) Where(condition string, args ...interface{}) *OpenProject {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *OpenProject) Set(clause string, args ...interface{}) *OpenProject {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *OpenProject) Limit(limit string) *OpenProject {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *OpenProject) Order(order string) *OpenProject {
	this.Dao.Order(order)
	return this
}

func (this *OpenProject) prepareInsertData() {
	this.columns = []string{"name", "category", "uri", "home", "doc", "download", "src", "logo", "desc", "repo", "author", "licence", "lang", "os", "tags", "username", "ctime"}
	this.colValues = []interface{}{this.Name, this.Category, this.Uri, this.Home, this.Doc, this.Download, this.Src, this.Logo, this.Desc, this.Repo, this.Author, this.Licence, this.Lang, this.Os, this.Tags, this.Username, this.Ctime}
}

func (this *OpenProject) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       &this.Id,
		"name":     &this.Name,
		"category": &this.Category,
		"uri":      &this.Uri,
		"home":     &this.Home,
		"doc":      &this.Doc,
		"download": &this.Download,
		"src":      &this.Src,
		"logo":     &this.Logo,
		"desc":     &this.Desc,
		"repo":     &this.Repo,
		"author":   &this.Author,
		"licence":  &this.Licence,
		"lang":     &this.Lang,
		"os":       &this.Os,
		"tags":     &this.Tags,
		"username": &this.Username,
		"viewnum":  &this.Viewnum,
		"cmtnum":   &this.Cmtnum,
		"likenum":  &this.Likenum,
		"status":   &this.Status,
		"ctime":    &this.Ctime,
		"mtime":    &this.Mtime,
	}
}
