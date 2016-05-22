// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"
	"time"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
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
func (ArticleLogic) FindBy(ctx context.Context, limit int, lastIds ...int) []*model.Article {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.Where("status IN(?,?)", model.ArticleStatusNew, model.ArticleStatusOnline)

	if len(lastIds) > 0 && lastIds[0] > 0 {
		dbSession.And("id<?", lastIds[0])
	}

	articles := make([]*model.Article, 0)
	err := dbSession.OrderBy("id DESC").Limit(limit).Find(&articles)
	if err != nil {
		objLog.Errorln("ArticleLogic FindBy Error:", err)
		return nil
	}

	topArticles := make([]*model.Article, 0)
	err = MasterDB.Where("top=?", 1).OrderBy("id DESC").Find(&topArticles)
	if err != nil {
		objLog.Errorln("ArticleLogic Find Top Articles Error:", err)
		return nil
	}
	if len(topArticles) > 0 {
		articles = append(topArticles, articles...)
	}

	return articles
}

// 获取抓取的文章列表（分页）：后台用
func (ArticleLogic) FindArticleByPage(ctx context.Context, conds map[string]string, curPage, limit int) ([]*model.Article, int) {
	objLog := GetLogger(ctx)

	offset := (curPage - 1) * limit
	session := MasterDB.Limit(limit, offset)
	for k, v := range conds {
		session.And(k+"=?", v)
	}

	totalSession := session.Clone()

	articleList := make([]*model.Article, 0)
	err := session.OrderBy("id DESC").Find(&articleList)
	if err != nil {
		objLog.Errorln("find error:", err)
		return nil, 0
	}

	total, err := totalSession.Count(new(model.Article))
	if err != nil {
		objLog.Errorln("find count error:", err)
		return nil, 0
	}

	return articleList, int(total)
}

// FindByIds 获取多个文章详细信息
func (ArticleLogic) FindByIds(ids []int) []*model.Article {
	if len(ids) == 0 {
		return nil
	}
	articles := make([]*model.Article, 0)
	err := MasterDB.In("id", ids).Find(&articles)
	if err != nil {
		logger.Errorln("ArticleLogic FindByIds error:", err)
		return nil
	}
	return articles
}

// findByIds 获取多个文章详细信息 包内使用
func (ArticleLogic) findByIds(ids []int) map[int]*model.Article {
	if len(ids) == 0 {
		return nil
	}
	articles := make(map[int]*model.Article)
	err := MasterDB.In("id", ids).Find(&articles)
	if err != nil {
		logger.Errorln("ArticleLogic findByIds error:", err)
		return nil
	}
	return articles
}

// FindByIdAndPreNext 获取当前(id)博文以及前后博文
func (ArticleLogic) FindByIdAndPreNext(ctx context.Context, id int) (curArticle *model.Article, prevNext []*model.Article, err error) {
	objLog := GetLogger(ctx)

	articles := make([]*model.Article, 0)

	err = MasterDB.Where("id BETWEEN ? AND ? AND status!=?", id-5, id+5, model.ArticleStatusOffline).Find(&articles)
	if err != nil {
		objLog.Errorln("ArticleLogic FindByIdAndPreNext Error:", err)
		return
	}

	if len(articles) == 0 {
		objLog.Errorln("ArticleLogic FindByIdAndPreNext not find articles, id:", id)
		return
	}

	prevNext = make([]*model.Article, 2)
	prevId, nextId := articles[0].Id, articles[len(articles)-1].Id
	for _, article := range articles {
		if article.Id < id && article.Id > prevId {
			prevId = article.Id
			prevNext[0] = article
		} else if article.Id > id && article.Id < nextId {
			nextId = article.Id
			prevNext[1] = article
		} else if article.Id == id {
			curArticle = article
		}
	}

	if prevId == id {
		prevNext[0] = nil
	}

	if nextId == id {
		prevNext[1] = nil
	}

	return
}

// 博文评论
type ArticleComment struct{}

// UpdateComment 更新该文章的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ArticleComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新评论数（TODO：暂时每次都更新表）
	_, err := MasterDB.Id(objid).Incr("cmtnum", 1).Update(new(model.Article))
	if err != nil {
		logger.Errorln("更新文章评论数失败：", err)
	}
}

func (self ArticleComment) String() string {
	return "article"
}

// SetObjinfo 实现 CommentObjecter 接口
func (self ArticleComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	articles := DefaultArticle.FindByIds(ids)
	if len(articles) == 0 {
		return
	}

	for _, article := range articles {
		objinfo := make(map[string]interface{})
		objinfo["title"] = article.Title
		objinfo["uri"] = model.PathUrlMap[model.TypeArticle]
		objinfo["type_name"] = model.TypeNameMap[model.TypeArticle]

		for _, comment := range commentMap[article.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 博文喜欢
type ArticleLike struct{}

// 更新该文章的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self ArticleLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Where("id=?", objid).Incr("likenum", num).Update(new(model.Article))
	if err != nil {
		logger.Errorln("更新文章喜欢数失败：", err)
	}
}

func (self ArticleLike) String() string {
	return "article"
}
