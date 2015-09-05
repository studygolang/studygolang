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
	FLAG_NO_AUDIT = iota
	FLAG_NORMAL
	FLAG_AUDIT_DELETE
	FLAG_USER_DELETE
)

// 帖子信息
type Topic struct {
	Tid           int    `json:"tid"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Nid           int    `json:"nid"`
	Uid           int    `json:"uid"`
	Flag          int    `json:"flag"`
	Lastreplyuid  int    `json:"lastreplyuid"`
	Lastreplytime string `json:"lastreplytime"`
	EditorUid     int    `json:"editor_uid"`
	Top           uint8  `json:"top"`
	Ctime         string `json:"ctime"`
	Mtime         string `json:"mtime"`

	// 为了方便，加上Node（节点名称，数据表没有）
	Node string
	// 数据库访问对象
	*Dao
}

func NewTopic() *Topic {
	return &Topic{
		Dao: &Dao{tablename: "topics"},
	}
}

func (this *Topic) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (this *Topic) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Topic) FindAll(selectCol ...string) ([]*Topic, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	topicList := make([]*Topic, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		topic := NewTopic()
		err = this.Scan(rows, colNum, topic.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Topic FindAll Scan Error:", err)
			continue
		}
		topicList = append(topicList, topic)
	}
	return topicList, nil
}

// 为了支持连写
func (this *Topic) Set(clause string, args ...interface{}) *Topic {
	this.Dao.Set(clause, args...)
	return this
}

// 为了支持连写
func (this *Topic) Where(condition string, args ...interface{}) *Topic {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *Topic) Limit(limit string) *Topic {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Topic) Order(order string) *Topic {
	this.Dao.Order(order)
	return this
}

func (this *Topic) prepareInsertData() {
	this.columns = []string{"title", "content", "nid", "uid", "ctime"}
	this.colValues = []interface{}{this.Title, this.Content, this.Nid, this.Uid, this.Ctime}
}

func (this *Topic) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"tid":           &this.Tid,
		"title":         &this.Title,
		"content":       &this.Content,
		"nid":           &this.Nid,
		"uid":           &this.Uid,
		"flag":          &this.Flag,
		"lastreplyuid":  &this.Lastreplyuid,
		"lastreplytime": &this.Lastreplytime,
		"editor_uid":    &this.EditorUid,
		"top":           &this.Top,
		"ctime":         &this.Ctime,
		"mtime":         &this.Mtime,
	}
}

// 帖子扩展（计数）信息
type TopicEx struct {
	Tid   int    `json:"tid"`
	View  int    `json:"view"`
	Reply int    `json:"reply"`
	Like  int    `json:"like"`
	Mtime string `json:"mtime"`

	// 数据库访问对象
	*Dao
}

func NewTopicEx() *TopicEx {
	return &TopicEx{
		Dao: &Dao{tablename: "topics_ex"},
	}
}

func (this *TopicEx) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	num, err := result.RowsAffected()
	return int(num), err
}

func (this *TopicEx) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *TopicEx) FindAll(selectCol ...string) ([]*TopicEx, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	topicExList := make([]*TopicEx, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		topicEx := NewTopicEx()
		err = this.Scan(rows, colNum, topicEx.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("TopicEx FindAll Scan Error:", err)
			continue
		}
		topicExList = append(topicExList, topicEx)
	}
	return topicExList, nil
}

// 为了支持连写
func (this *TopicEx) Where(condition string, args ...interface{}) *TopicEx {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *TopicEx) Limit(limit string) *TopicEx {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *TopicEx) Order(order string) *TopicEx {
	this.Dao.Order(order)
	return this
}

func (this *TopicEx) prepareInsertData() {
	this.columns = []string{"tid", "view", "reply", "like"}
	this.colValues = []interface{}{this.Tid, this.View, this.Reply, this.Like}
}

func (this *TopicEx) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"tid":   &this.Tid,
		"view":  &this.View,
		"reply": &this.Reply,
		"like":  &this.Like,
		"mtime": &this.Mtime,
	}
}

// 帖子节点信息
type TopicNode struct {
	Nid    int    `json:"nid"`
	Parent int    `json:"parent"`
	Name   string `json:"name"`
	Intro  string `json:"intro"`
	Ctime  string `json:"ctime"`

	// 数据库访问对象
	*Dao
}

func NewTopicNode() *TopicNode {
	return &TopicNode{
		Dao: &Dao{tablename: "topics_node"},
	}
}

func (this *TopicNode) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (this *TopicNode) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *TopicNode) FindAll(selectCol ...string) ([]*TopicNode, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	nodeList := make([]*TopicNode, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		node := NewTopicNode()
		err = this.Scan(rows, colNum, node.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("TopicNode FindAll Scan Error:", err)
			continue
		}
		nodeList = append(nodeList, node)
	}
	return nodeList, nil
}

func (this *TopicNode) prepareInsertData() {
	this.columns = []string{"parent", "name", "intro"}
	this.colValues = []interface{}{this.Parent, this.Name, this.Intro}
}

func (this *TopicNode) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"nid":    &this.Nid,
		"parent": &this.Parent,
		"name":   &this.Name,
		"intro":  &this.Intro,
		"ctime":  &this.Ctime,
	}
}
