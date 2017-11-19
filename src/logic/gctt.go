// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	"model"

	. "db"
)

type GCTTLogic struct{}

var DefaultGCTT = GCTTLogic{}

func (self GCTTLogic) FindTranslator(ctx context.Context, me *model.Me) *model.GCTTUser {
	objLog := GetLogger(ctx)

	gcttUser := &model.GCTTUser{}
	_, err := MasterDB.Where("uid=?", me.Uid).Get(gcttUser)
	if err != nil {
		objLog.Errorln("GCTTLogic FindTranslator error:", err)
		return nil
	}

	return gcttUser
}

func (self GCTTLogic) FindOne(ctx context.Context, username string) *model.GCTTUser {
	objLog := GetLogger(ctx)

	gcttUser := &model.GCTTUser{}
	_, err := MasterDB.Where("username=?", username).Get(gcttUser)
	if err != nil {
		objLog.Errorln("GCTTLogic FindOne error:", err)
		return nil
	}

	return gcttUser
}

func (self GCTTLogic) BindUser(ctx context.Context, gcttUser *model.GCTTUser, uid int, githubUser *model.BindUser) error {
	objLog := GetLogger(ctx)

	var err error

	if gcttUser.Id > 0 {
		gcttUser.Uid = uid
		_, err = MasterDB.Id(gcttUser.Id).Update(gcttUser)
	} else {
		gcttUser = &model.GCTTUser{
			Username: githubUser.Username,
			Avatar:   githubUser.Avatar,
			Uid:      uid,
		}
		_, err = MasterDB.Insert(gcttUser)
	}

	if err != nil {
		objLog.Errorln("GCTTLogic BindUser error:", err)
	}

	return err
}
