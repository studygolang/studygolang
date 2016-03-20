// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"
	"strconv"

	"github.com/polaris1119/logger"
)

type ArticleLogic struct{}

var DefaultArticle = ArticleLogic{}

func (ArticleLogic) FindLastList(beginTime string, limit int) ([]*model.Article, error) {
	articles := make([]*model.Article, 0)
	err := MasterDB.Where("ctime>? AND status!=?", beginTime, model.ArticleStatusOffline).
		OrderBy("cmtnum DESC, likenum DESC, viewnum DESC").Limit(limit).Find(&articles)

	return articles, err
}

// Total 博文总数
func (ArticleLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.Article))
	if err != nil {
		logger.Errorln("ArticleLogic Total error:", err)
	}
	return total
}

// FindBy 获取抓取的文章列表（分页）
func (ArticleLogic) FindBy(limit int, lastIds ...int) []*model.Article {
	dbSession := MasterDB.Where("status IN(?,?)", model.ArticleStatusNew, model.ArticleStatusOnline)

	if len(lastIds) > 0 {
		dbSession.And("id<?", lastIds[0])
	}

	articles := make([]*model.Article, 0)
	err := dbSession.OrderBy("id DESC").Limit(limit).Find(&articles)
	if err != nil {
		logger.Errorln("ArticleLogic FindBy Error:", err)
		return nil
	}

	topArticles := make([]*model.Article, 0)
	err = MasterDB.Where("top=?", 1).OrderBy("id DESC").Find(&topArticles)
	if err != nil {
		logger.Errorln("ArticleLogic Find Top Articles Error:", err)
		return nil
	}
	if len(topArticles) > 0 {
		articles = append(topArticles, articles...)
	}

	return articles
}

// 博文评论
type ArticleComment struct{}

// 更新该文章的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ArticleComment) UpdateComment(cid, objid, uid int, cmttime string) {
	id := strconv.Itoa(objid)

	// 更新评论数（TODO：暂时每次都更新表）
	err := model.NewArticle().Where("id="+id).Increment("cmtnum", 1)
	if err != nil {
		logger.Errorln("更新文章评论数失败：", err)
	}
}

func (self ArticleComment) String() string {
	return "article"
}

// 实现 CommentObjecter 接口
func (self ArticleComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	articles := FindArticlesByIds(ids)
	if len(articles) == 0 {
		return
	}

	for _, article := range articles {
		objinfo := make(map[string]interface{})
		objinfo["title"] = article.Title
		objinfo["uri"] = model.PathUrlMap[model.TYPE_ARTICLE]
		objinfo["type_name"] = model.TypeNameMap[model.TYPE_ARTICLE]

		for _, comment := range commentMap[article.Id] {
			comment.Objinfo = objinfo
		}
	}
}
