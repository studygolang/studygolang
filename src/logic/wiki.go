// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	. "db"
	"model"

	"golang.org/x/net/context"

	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
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

	go publishObservable.NotifyObservers(me.Uid, model.TypeWiki, wiki.Id)

	return nil
}

func (self WikiLogic) Modify(ctx context.Context, me *model.Me, form url.Values) error {
	objLog := GetLogger(ctx)

	id := goutils.MustInt(form.Get("id"))
	wiki := self.FindById(ctx, id)
	if !CanEdit(me, wiki) {
		return errors.New("没有权限")
	}

	if wiki.Uid != me.Uid {
		hasExists := false
		cuids := strings.Split(wiki.Cuid, ",")
		for _, cuid := range cuids {
			if me.Uid == goutils.MustInt(cuid) {
				hasExists = true
				break
			}
		}

		if !hasExists {
			cuids = append(cuids, strconv.Itoa(me.Uid))
			wiki.Cuid = strings.Join(cuids, ",")
		}
	}

	wiki.Title = form.Get("title")
	wiki.Content = form.Get("content")

	_, err := MasterDB.Id(id).Update(wiki)
	if err != nil {
		objLog.Errorf("更新wiki 【%d】 信息失败：%s\n", id, err)
		return err
	}

	go modifyObservable.NotifyObservers(me.Uid, model.TypeWiki, wiki.Id)

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
	usersMap := DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))
	for _, wiki := range wikis {
		wiki.Users = map[int]*model.User{wiki.Uid: usersMap[wiki.Uid]}
	}

	return wikis
}

// FindById 通过ID获取Wiki
func (WikiLogic) FindById(ctx context.Context, id int) *model.Wiki {
	objLog := GetLogger(ctx)

	wiki := &model.Wiki{}
	if _, err := MasterDB.Where("id=?", id).Get(wiki); err != nil {
		objLog.Errorln("wiki logic FindById error:", err)
		return nil
	}
	return wiki
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
	wiki.Users = DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))

	return wiki
}

// getOwner 通过id获得wiki的所有者
func (WikiLogic) getOwner(id int) int {
	wiki := &model.Wiki{}
	_, err := MasterDB.Id(id).Get(wiki)
	if err != nil {
		logger.Errorln("wiki logic getOwner Error:", err)
		return 0
	}
	return wiki.Uid
}

// FindByIds 获取多个wiki页面详细信息
func (WikiLogic) FindByIds(ids []int) []*model.Wiki {
	if len(ids) == 0 {
		return nil
	}
	wikis := make([]*model.Wiki, 0)
	err := MasterDB.In("id", ids).Find(&wikis)
	if err != nil {
		logger.Errorln("wiki logic FindByIds error:", err)
		return nil
	}
	return wikis
}

// findByIds 获取多个wiki页面详细信息 包内使用
func (WikiLogic) findByIds(ids []int) map[int]*model.Wiki {
	if len(ids) == 0 {
		return nil
	}
	wikis := make(map[int]*model.Wiki)
	err := MasterDB.In("id", ids).Find(&wikis)
	if err != nil {
		logger.Errorln("wiki logic FindByIds error:", err)
		return nil
	}
	return wikis
}
