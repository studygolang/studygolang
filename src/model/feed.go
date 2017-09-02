// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"db"

	"github.com/polaris1119/logger"
)

const FeedOffline = 1

type Feed struct {
	Id            int `xorm:"pk autoincr"`
	Title         string
	Objid         int
	Objtype       int
	Uid           int
	Author        string
	Nid           int
	Lastreplyuid  int
	Lastreplytime OftenTime
	Tags          string
	Cmtnum        int
	Top           uint8
	State         int
	CreatedAt     OftenTime `xorm:"created"`
	UpdatedAt     OftenTime `json:"updated_at" xorm:"<-"`

	User          *User                  `xorm:"-"`
	Lastreplyuser *User                  `xorm:"-"`
	Node          map[string]interface{} `xorm:"-"`
	Uri           string                 `xorm:"-"`
}

// PublishFeed 发布动态
func PublishFeed(object interface{}, objectExt interface{}) {
	var feed *Feed
	switch objdoc := object.(type) {
	case *Topic:
		node := &TopicNode{}
		_, err := db.MasterDB.Id(objdoc.Nid).Get(node)
		if err == nil && !node.ShowIndex {
			return
		}

		cmtnum := 0
		if objectExt != nil {
			// 传递过来的是一个 *TopicEx 对象，类型是有的，即时值是 nil，这里也和 nil 是不等
			topicEx := objectExt.(*TopicEx)
			if topicEx != nil {
				cmtnum = topicEx.Reply
			}
		}

		feed = &Feed{
			Objid:         objdoc.Tid,
			Objtype:       TypeTopic,
			Title:         objdoc.Title,
			Uid:           objdoc.Uid,
			Tags:          objdoc.Tags,
			Cmtnum:        cmtnum,
			Nid:           objdoc.Nid,
			Top:           objdoc.Top,
			Lastreplyuid:  objdoc.Lastreplyuid,
			Lastreplytime: objdoc.Lastreplytime,
			UpdatedAt:     objdoc.Mtime,
		}
	case *Article:
		var uid int
		if objdoc.Domain == WebsiteSetting.Domain {
			userLogin := &UserLogin{}
			db.MasterDB.Where("username=?", objdoc.AuthorTxt).Get(userLogin)
			uid = userLogin.Uid
		}
		feed = &Feed{
			Objid:         objdoc.Id,
			Objtype:       TypeArticle,
			Title:         FilterTxt(objdoc.Title),
			Author:        objdoc.AuthorTxt,
			Uid:           uid,
			Tags:          objdoc.Tags,
			Cmtnum:        objdoc.Cmtnum,
			Top:           objdoc.Top,
			Lastreplyuid:  objdoc.Lastreplyuid,
			Lastreplytime: objdoc.Lastreplytime,
			UpdatedAt:     objdoc.Mtime,
		}
	case *Resource:
		cmtnum := 0
		if objectExt != nil {
			resourceEx := objectExt.(*ResourceEx)
			if resourceEx != nil {
				cmtnum = resourceEx.Cmtnum
			}
		}

		feed = &Feed{
			Objid:         objdoc.Id,
			Objtype:       TypeResource,
			Title:         objdoc.Title,
			Uid:           objdoc.Uid,
			Tags:          objdoc.Tags,
			Cmtnum:        cmtnum,
			Nid:           objdoc.Catid,
			Lastreplyuid:  objdoc.Lastreplyuid,
			Lastreplytime: objdoc.Lastreplytime,
			UpdatedAt:     objdoc.Mtime,
		}
	case *OpenProject:
		userLogin := &UserLogin{}
		db.MasterDB.Where("username=?", objdoc.Username).Get(userLogin)
		feed = &Feed{
			Objid:         objdoc.Id,
			Objtype:       TypeProject,
			Title:         objdoc.Category + " " + objdoc.Name,
			Author:        objdoc.Author,
			Uid:           userLogin.Uid,
			Tags:          objdoc.Tags,
			Cmtnum:        objdoc.Cmtnum,
			Lastreplyuid:  objdoc.Lastreplyuid,
			Lastreplytime: objdoc.Lastreplytime,
			UpdatedAt:     objdoc.Mtime,
		}
	case *Book:
		feed = &Feed{
			Objid:         objdoc.Id,
			Objtype:       TypeBook,
			Title:         "分享一本图书《" + objdoc.Name + "》",
			Uid:           objdoc.Uid,
			Tags:          objdoc.Tags,
			Cmtnum:        objdoc.Cmtnum,
			Lastreplyuid:  objdoc.Lastreplyuid,
			Lastreplytime: objdoc.Lastreplytime,
			UpdatedAt:     objdoc.UpdatedAt,
		}
	}

	_, err := db.MasterDB.Insert(feed)
	if err != nil {
		logger.Errorln("publish feed:", object, " error:", err)
	}
}
