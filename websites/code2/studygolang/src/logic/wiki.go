// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"net/url"
	"strings"

	. "db"
	"model"

	"golang.org/x/net/context"

	"github.com/polaris1119/goutils"
	"github.com/polaris1119/set"
)

type WikiLogic struct{}

var DefaultWiki = WikiLogic{}

// Create 创建一个wiki页面
func (WikiLogic) Create(ctx context.Context, me *model.Me, form url.Values) error {
	objLog := GetLogger(ctx)

	wiki := &model.Wiki{}
	err := schemaDecoder.Decode(wiki, form)
	if err != nil {
		objLog.Errorln("Create Wiki schema decode error:", err)
		return err
	}

	wiki.Uid = me.Uid
	if _, err = MasterDB.Insert(wiki); err != nil {
		objLog.Errorln("Create Wiki error:", err)
		return err
	}

	// 创建一个wiki页面，活跃度+10
	go DefaultUser.IncrUserWeight("uid", me.Uid, 10)

	return nil
}

// FindBy 获取 wiki 列表（分页）
func (WikiLogic) FindBy(ctx context.Context, limit int, lastIds ...int) []*model.Wiki {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.OrderBy("id DESC")

	if len(lastIds) > 0 && lastIds[0] > 0 {
		dbSession.Where("id<?", lastIds[0])
	}

	wikis := make([]*model.Wiki, 0)
	err := dbSession.Limit(limit).Find(&wikis)
	if err != nil {
		objLog.Errorln("WikiLogic FindBy Error:", err)
		return nil
	}

	uidSet := set.New(set.NonThreadSafe)
	for _, wiki := range wikis {
		uidSet.Add(wiki.Uid)
	}
	usersMap := DefaultUser.FindUserInfos(ctx, set.Int64Slice(uidSet))
	for _, wiki := range wikis {
		wiki.Users = map[int]*model.User{wiki.Uid: usersMap[wiki.Uid]}
	}

	return wikis
}

// FindOne 某个wiki页面详细信息
func (WikiLogic) FindOne(ctx context.Context, uri string) *model.Wiki {
	objLog := GetLogger(ctx)

	wiki := &model.Wiki{}
	if _, err := MasterDB.Where("uri=?", uri).Get(wiki); err != nil {
		objLog.Errorln("wiki logic FindOne error:", err)
		return nil
	}

	if wiki.Id == 0 {
		return nil
	}

	uidSet := set.New(set.NonThreadSafe)
	uidSet.Add(wiki.Uid)
	if wiki.Cuid != "" {
		cuids := strings.Split(wiki.Cuid, ",")
		for _, cuid := range cuids {
			uidSet.Add(goutils.MustInt(cuid))
		}
	}
	wiki.Users = DefaultUser.FindUserInfos(ctx, set.Int64Slice(uidSet))

	return wiki
}

// // 获取多个wiki页面详细信息
// func FindWikisByIds(ids []int) []*model.Wiki {
// 	if len(ids) == 0 {
// 		return nil
// 	}
// 	inIds := util.Join(ids, ",")
// 	wikis, err := model.NewWiki().Where("id in(" + inIds + ")").FindAll()
// 	if err != nil {
// 		logger.Errorln("wiki service FindWikisByIds error:", err)
// 		return nil
// 	}
// 	return wikis
// }
