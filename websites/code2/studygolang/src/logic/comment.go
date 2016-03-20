// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"fmt"
	"model"
	"regexp"
	"text/template"
	"util"

	. "db"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

type CommentLogic struct{}

var DefaultComment = CommentLogic{}

// Total 评论总数(objtypes[0] 取某一类型的评论总数)
func (CommentLogic) Total(objtypes ...int) int64 {
	var (
		total int64
		err   error
	)
	if len(objtypes) > 0 {
		total, err = MasterDB.Where("objtype=?", objtypes[0]).Count(new(model.Comment))
	} else {

		total, err = MasterDB.Count(new(model.Comment))
	}
	if err != nil {
		logger.Errorln("CommentLogic Total error:", err)
	}
	return total
}

// FindRecent 获得最近的评论
// 如果 uid!=0，表示获取某人的评论；
// 如果 objtype!=-1，表示获取某类型的评论；
func (self CommentLogic) FindRecent(ctx context.Context, uid, objtype, limit int) []*model.Comment {
	dbSession := MasterDB.OrderBy("cid DESC").Limit(limit)

	if uid != 0 {
		dbSession.And("uid=?", uid)
	}
	if objtype != -1 {
		dbSession.And("objtype=?", objtype)
	}

	comments := make([]*model.Comment, 0)
	err := dbSession.Find(&comments)
	if err != nil {
		logger.Errorln("CommentLogic FindRecent error:", err)
		return nil
	}

	cmtMap := make(map[int][]*model.Comment, len(model.PathUrlMap))
	for _, comment := range comments {
		self.decodeCmtContent(ctx, comment)

		if _, ok := cmtMap[comment.Objtype]; !ok {
			cmtMap[comment.Objtype] = make([]*model.Comment, 0, 10)
		}

		cmtMap[comment.Objtype] = append(cmtMap[comment.Objtype], comment)
	}

	cmtObjs := []CommentObjecter{
		model.TypeTopic:    TopicComment{},
		model.TypeArticle:  ArticleComment{},
		model.TypeResource: ResourceComment{},
		model.TypeWiki:     nil,
		model.TypeProject:  ProjectComment{},
	}
	for cmtType, cmts := range cmtMap {
		FillCommentObjs(cmts, cmtObjs[cmtType])
	}

	return comments
}

// fillObjinfos 填充评论对应的主体信息
func (CommentLogic) fillObjinfos(comments []*model.Comment, cmtObj CommentObjecter) {
	if len(comments) == 0 {
		return
	}
	count := len(comments)
	commentMap := make(map[int][]*model.Comment, count)
	idMap := make(map[int]int, count)
	for _, comment := range comments {
		if _, ok := commentMap[comment.Objid]; !ok {
			commentMap[comment.Objid] = make([]*model.Comment, 0, count)
		}
		commentMap[comment.Objid] = append(commentMap[comment.Objid], comment)
		idMap[comment.Objid] = 1
	}
	ids := util.MapIntKeys(idMap)
	cmtObj.SetObjinfo(ids, commentMap)
}

func (CommentLogic) decodeCmtContent(ctx context.Context, comment *model.Comment) string {
	// 安全过滤
	content := template.HTMLEscapeString(comment.Content)
	// @别人
	content = parseAtUser(ctx, content)

	// 回复某一楼层
	reg := regexp.MustCompile(`#(\d+)楼`)
	url := fmt.Sprintf("%s%d#comment", model.PathUrlMap[comment.Objtype], comment.Objid)
	content = reg.ReplaceAllString(content, `<a href="`+url+`$1" title="$1">#$1<span>楼</span></a>`)

	comment.Content = content

	return content
}

// 填充 Comment 对象的 Objinfo 成员接口
// 评论属主应该实现该接口（以便填充 Objinfo 成员）
type CommentObjecter interface {
	// ids 是属主的主键 slice （comment 中的 objid）
	// commentMap 中的 key 是属主 id
	SetObjinfo(ids []int, commentMap map[int][]*model.Comment)
}
