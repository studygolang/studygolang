// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"model"
	"net/url"
	"util"

	"github.com/polaris1119/slices"
	"golang.org/x/net/context"

	. "db"
	"github.com/polaris1119/goutils"
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
		Where("sid=? AND state=?", sid, model.ContributeStateOnline).OrderBy(order).Find(&subjectArticles)
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

// Follow 关注或取消关注
func (self SubjectLogic) Follow(ctx context.Context, sid int, me *model.Me) (err error) {
	objLog := GetLogger(ctx)

	follower := &model.SubjectFollower{}
	_, err = MasterDB.Where("sid=? AND uid=?", sid, me.Uid).Get(follower)
	if err != nil {
		objLog.Errorln("SubjectLogic Follow Get error:", err)
	}

	if follower.Id > 0 {
		_, err = MasterDB.Where("sid=? AND uid=?", sid, me.Uid).Delete(new(model.SubjectFollower))
		if err != nil {
			objLog.Errorln("SubjectLogic Follow Delete error:", err)
		}

		return
	}

	follower.Sid = sid
	follower.Uid = me.Uid
	_, err = MasterDB.Insert(follower)
	if err != nil {
		objLog.Errorln("SubjectLogic Follow insert error:", err)
	}
	return
}

func (self SubjectLogic) HadFollow(ctx context.Context, sid int, me *model.Me) bool {
	objLog := GetLogger(ctx)

	num, err := MasterDB.Where("sid=? AND uid=?", sid, me.Uid).Count(new(model.SubjectFollower))
	if err != nil {
		objLog.Errorln("SubjectLogic Follow insert error:", err)
	}

	return num > 0
}

// Contribute 投稿
func (self SubjectLogic) Contribute(ctx context.Context, me *model.Me, sid, articleId int) error {
	objLog := GetLogger(ctx)

	subject := self.FindOne(ctx, sid)
	if subject.Id == 0 {
		return errors.New("该专题不存在")
	}

	count, _ := MasterDB.Where("article_id=?", articleId).Count(new(model.SubjectArticle))
	if count >= 5 {
		return errors.New("该文超过 5 次投稿")
	}

	article, err := DefaultArticle.FindById(ctx, articleId)
	if article.AuthorTxt != me.Username {
		return errors.New("该文不是你的文章，不能投稿")
	}

	subjectArticle := &model.SubjectArticle{
		Sid:       sid,
		ArticleId: articleId,
		State:     model.ContributeStateNew,
	}
	if subject.Uid == me.Uid {
		subjectArticle.State = model.ContributeStateOnline
	}

	_, err = MasterDB.Insert(subjectArticle)
	if err != nil {
		objLog.Errorln("SubjectLogic Contribute insert error:", err)
		return errors.New("投稿失败:" + err.Error())
	}

	return nil
}

// RemoveContribute 删除投稿
func (self SubjectLogic) RemoveContribute(ctx context.Context, sid, articleId int) error {
	objLog := GetLogger(ctx)

	_, err := MasterDB.Where("sid=? AND article_id=?", sid, articleId).Delete(new(model.SubjectArticle))
	if err != nil {
		objLog.Errorln("SubjectLogic RemoveContribute delete error:", err)
		return errors.New("删除投稿失败:" + err.Error())
	}

	return nil
}

func (self SubjectLogic) ExistByName(name string) bool {
	exist, _ := MasterDB.Where("name=?", name).Exist(new(model.Subject))
	return exist
}

// Publish 发布专题。
func (self SubjectLogic) Publish(ctx context.Context, me *model.Me, form url.Values) (sid int, err error) {
	objLog := GetLogger(ctx)

	sid = goutils.MustInt(form.Get("sid"))
	if sid != 0 {
		subject := &model.Subject{}
		_, err = MasterDB.Id(sid).Get(subject)
		if err != nil {
			objLog.Errorln("Publish Subject find error:", err)
			return
		}

		_, err = self.Modify(ctx, me, form)
		if err != nil {
			objLog.Errorln("Publish Subject modify error:", err)
			return
		}

	} else {
		subject := &model.Subject{}
		err = schemaDecoder.Decode(subject, form)
		if err != nil {
			objLog.Errorln("SubjectLogic Publish decode error:", err)
			return
		}
		subject.Uid = me.Uid

		_, err = MasterDB.Insert(subject)
		if err != nil {
			objLog.Errorln("SubjectLogic Publish insert error:", err)
			return
		}
		sid = subject.Id
	}
	return
}

// Modify 修改专题
func (SubjectLogic) Modify(ctx context.Context, user *model.Me, form url.Values) (errMsg string, err error) {
	objLog := GetLogger(ctx)

	change := map[string]interface{}{}

	fields := []string{"name", "description", "cover", "contribute", "audit"}
	for _, field := range fields {
		change[field] = form.Get(field)
	}

	sid := form.Get("sid")
	_, err = MasterDB.Table(new(model.Subject)).Id(sid).Update(change)
	if err != nil {
		objLog.Errorf("更新专题 【%s】 信息失败：%s\n", sid, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	return
}
