// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"db"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// 文档对象（供solr使用）
type Document struct {
	Id      string `json:"id"`
	Objid   int    `json:"objid"`
	Objtype int    `json:"objtype"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	PubTime string `json:"pub_time"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
	Viewnum int    `json:"viewnum"`
	Cmtnum  int    `json:"cmtnum"`
	Likenum int    `json:"likenum"`

	HlTitle   string `json:",omitempty"` // 高亮的标题
	HlContent string `json:",omitempty"` // 高亮的内容
}

func NewDocument(object interface{}, objectExt interface{}) *Document {
	var document *Document
	switch objdoc := object.(type) {
	case *Topic:
		viewnum, cmtnum, likenum := 0, 0, 0
		if objectExt != nil {
			// 传递过来的是一个 *TopicEx 对象，类型是有的，即时值是 nil，这里也和 nil 是不等
			topicEx := objectExt.(*TopicEx)
			if topicEx != nil {
				viewnum = topicEx.View
				cmtnum = topicEx.Reply
				likenum = topicEx.Like
			}
		}

		userLogin := &UserLogin{}
		db.MasterDB.Id(objdoc.Uid).Get(userLogin)
		document = &Document{
			Id:      fmt.Sprintf("%d%d", TypeTopic, objdoc.Tid),
			Objid:   objdoc.Tid,
			Objtype: TypeTopic,
			Title:   objdoc.Title,
			Author:  userLogin.Username,
			PubTime: objdoc.Ctime.String(),
			Content: objdoc.Content,
			Tags:    "",
			Viewnum: viewnum,
			Cmtnum:  cmtnum,
			Likenum: likenum,
		}
	case *Article:
		document = &Document{
			Id:      fmt.Sprintf("%d%d", TypeArticle, objdoc.Id),
			Objid:   objdoc.Id,
			Objtype: TypeArticle,
			Title:   filterTxt(objdoc.Title),
			Author:  objdoc.AuthorTxt,
			PubTime: objdoc.PubDate,
			Content: filterTxt(objdoc.Txt),
			Tags:    objdoc.Tags,
			Viewnum: objdoc.Viewnum,
			Cmtnum:  objdoc.Cmtnum,
			Likenum: objdoc.Likenum,
		}
	case *Resource:
		viewnum, cmtnum, likenum := 0, 0, 0
		if objectExt != nil {
			resourceEx := objectExt.(*ResourceEx)
			if resourceEx != nil {
				viewnum = resourceEx.Viewnum
				cmtnum = resourceEx.Cmtnum
			}
		}

		userLogin := &UserLogin{}
		db.MasterDB.Id(objdoc.Uid).Get(userLogin)
		document = &Document{
			Id:      fmt.Sprintf("%d%d", TypeResource, objdoc.Id),
			Objid:   objdoc.Id,
			Objtype: TypeResource,
			Title:   objdoc.Title,
			Author:  userLogin.Username,
			PubTime: objdoc.Ctime.String(),
			Content: template.HTMLEscapeString(objdoc.Content),
			Tags:    "",
			Viewnum: viewnum,
			Cmtnum:  cmtnum,
			Likenum: likenum,
		}
	case *OpenProject:
		document = &Document{
			Id:      fmt.Sprintf("%d%d", TypeProject, objdoc.Id),
			Objid:   objdoc.Id,
			Objtype: TypeProject,
			Title:   objdoc.Category + objdoc.Name,
			Author:  objdoc.Author,
			PubTime: objdoc.Ctime.String(),
			Content: objdoc.Desc,
			Tags:    objdoc.Tags,
			Viewnum: objdoc.Viewnum,
			Cmtnum:  objdoc.Cmtnum,
			Likenum: objdoc.Likenum,
		}
	}

	return document
}

var docRe = regexp.MustCompile("[\r　\n  \t\v]+")
var docSpaceRe = regexp.MustCompile("[ ]+")

// 文本过滤（预处理）
func filterTxt(txt string) string {
	txt = strings.TrimSpace(strings.TrimPrefix(txt, "原"))
	txt = strings.TrimSpace(strings.TrimPrefix(txt, "荐"))
	txt = strings.TrimSpace(strings.TrimPrefix(txt, "顶"))
	txt = strings.TrimSpace(strings.TrimPrefix(txt, "转"))

	txt = docRe.ReplaceAllLiteralString(txt, " ")
	return docSpaceRe.ReplaceAllLiteralString(txt, " ")
}

type AddCommand struct {
	Doc          *Document `json:"doc"`
	Boost        float64   `json:"boost,omitempty"`
	Overwrite    bool      `json:"overwrite"`
	CommitWithin int       `json:"commitWithin,omitempty"`
}

func NewDefaultArgsAddCommand(doc *Document) *AddCommand {
	return NewAddCommand(doc, 0.0, true, 0)
}

func NewAddCommand(doc *Document, boost float64, overwrite bool, commitWithin int) *AddCommand {
	return &AddCommand{
		Doc:          doc,
		Boost:        boost,
		Overwrite:    overwrite,
		CommitWithin: commitWithin,
	}
}

type DelCommand struct {
	Id string `json:"id"`
}

func NewDelCommand(doc *Document) *DelCommand {
	return &DelCommand{Id: doc.Id}
}

type ResponseBody struct {
	NumFound int         `json:"numFound"`
	Start    int         `json:"start"`
	Docs     []*Document `json:"docs"`
}

type Highlighting struct {
	Title   []string `json:"title"`
	Content []string `json:"content"`
}

type SearchResponse struct {
	RespHeader map[string]interface{}   `json:"responseHeader"`
	RespBody   *ResponseBody            `json:"response"`
	Highlight  map[string]*Highlighting `json:"highlighting"`
}
