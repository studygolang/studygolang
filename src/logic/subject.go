// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"model"
	"util"

	"github.com/polaris1119/slices"
	"golang.org/x/net/context"

	. "db"
)

type SubjectLogic struct{}

var DefaultSubject = SubjectLogic{}

func (self SubjectLogic) FindOne(ctx context.Context, sid int) *model.Subject {
	objLog := GetLogger(ctx)

	subject := &model.Subject{}
	_, err := MasterDB.Id(sid).Get(subject)
	if err != nil {
		objLog.Errorln("SubjectLogic FindOne get error:", err)
	}

	if subject.Uid > 0 {
		subject.User = DefaultUser.findUser(ctx, subject.Uid)
	}

	return subject
}

func (self SubjectLogic) FindArticles(ctx context.Context, sid int, orderBy string) []*model.Article {
	objLog := GetLogger(ctx)

	order := "subject_article.created_at DESC"
	if orderBy == "commented_at" {
		order = "articles.lastreplytime DESC"
	}

	subjectArticles := make([]*model.SubjectArticles, 0)
	err := MasterDB.Join("INNER", "subject_article", "subject_article.article_id = articles.id").
		Where("subject_article.sid=?", sid).OrderBy(order).Find(&subjectArticles)
	if err != nil {
		objLog.Errorln("SubjectLogic FindArticles Find subject_article error:", err)
		return nil
	}

	articles := make([]*model.Article, 0, len(subjectArticles))
	for _, subjectArticle := range subjectArticles {
		if subjectArticle.Status == model.ArticleStatusOffline {
			continue
		}

		articles = append(articles, &subjectArticle.Article)
	}

	DefaultArticle.fillUser(articles)
	return articles
}

// FindArticleTotal 专题收录的文章数
func (self SubjectLogic) FindArticleTotal(ctx context.Context, sid int) int64 {
	objLog := GetLogger(ctx)

	total, err := MasterDB.Where("sid=?", sid).Count(new(model.SubjectArticle))
	if err != nil {
		objLog.Errorln("SubjectLogic FindArticleTotal error:", err)
	}

	return total
}

// FindFollowers 专题关注的用户
func (self SubjectLogic) FindFollowers(ctx context.Context, sid int) []*model.SubjectFollower {
	objLog := GetLogger(ctx)

	followers := make([]*model.SubjectFollower, 0)
	err := MasterDB.Where("sid=?", sid).OrderBy("id DESC").Limit(8).Find(&followers)
	if err != nil {
		objLog.Errorln("SubjectLogic FindFollowers error:", err)
	}

	if len(followers) == 0 {
		return followers
	}

	uids := slices.StructsIntSlice(followers, "Uid")
	usersMap := DefaultUser.FindUserInfos(ctx, uids)
	for _, follower := range followers {
		follower.User = usersMap[follower.Uid]
		follower.TimeAgo = util.TimeAgo(follower.CreatedAt)
	}

	return followers
}

// FindFollowerTotal 专题关注的用户数
func (self SubjectLogic) FindFollowerTotal(ctx context.Context, sid int) int64 {
	objLog := GetLogger(ctx)

	total, err := MasterDB.Where("sid=?", sid).Count(new(model.SubjectFollower))
	if err != nil {
		objLog.Errorln("SubjectLogic FindFollowerTotal error:", err)
	}

	return total
}
