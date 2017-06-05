// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"fmt"
	"model"

	. "db"

	"github.com/go-xorm/xorm"
	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

type UserRichLogic struct{}

var DefaultUserRich = UserRichLogic{}

func (self UserRichLogic) Add(ctx context.Context) error {
	return nil
}

// IncrUserRich 增加或减少用户财富
func (self UserRichLogic) IncrUserRich(user *model.User, typ, award int, desc string) {
	var (
		total int64 = -1
		err   error
	)

	if award > 0 && typ == model.MissionTypeReplied {
		// 老用户，因为之前的主题被人回复而增加财富，自动帮其领取初始资本
		total, err = MasterDB.Where("uid=?", user.Uid).Count(new(model.UserBalanceDetail))
		if err != nil {
			logger.Errorln("IncrUserRich count error:", err)
			return
		}
	}

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	var initialAward int
	if total == 0 {
		initialAward, err = self.autoCompleteInitial(session, user)
		if err != nil {
			logger.Errorln("IncrUserRich autoCompleteInitial error:", err)
			session.Rollback()
			return
		}
	}

	user.Balance += initialAward + award
	_, err = session.Where("uid=?", user.Uid).Cols("balance").Update(user)
	if err != nil {
		logger.Errorln("IncrUserRich update error:", err)
		session.Rollback()
		return
	}

	balanceDetail := &model.UserBalanceDetail{
		Uid:     user.Uid,
		Type:    typ,
		Num:     award,
		Balance: user.Balance,
		Desc:    desc,
	}
	_, err = session.Insert(balanceDetail)
	if err != nil {
		logger.Errorln("IncrUserRich insert error:", err)
		session.Rollback()
		return
	}

	session.Commit()
}

func (UserRichLogic) FindBalanceDetail(ctx context.Context, me *model.Me) []*model.UserBalanceDetail {
	objLog := GetLogger(ctx)

	balanceDetails := make([]*model.UserBalanceDetail, 0)
	err := MasterDB.Where("uid=?", me.Uid).Desc("id").Find(&balanceDetails)
	if err != nil {
		objLog.Errorln("UserRichLogic FindBalanceDetail error:", err)
		return nil
	}

	return balanceDetails
}

func (UserRichLogic) Total(ctx context.Context, uid int) int64 {
	total, err := MasterDB.Where("uid=?", uid).Count(new(model.UserBalanceDetail))
	if err != nil {
		logger.Errorln("UserRichLogic Total error:", err)
	}
	return total
}

func (UserRichLogic) add(session *xorm.Session, balanceDetail *model.UserBalanceDetail) error {
	_, err := session.Insert(balanceDetail)
	return err
}

func (UserRichLogic) autoCompleteInitial(session *xorm.Session, user *model.User) (int, error) {
	mission := &model.Mission{}
	_, err := session.Where("id=?", model.InitialMissionId).Get(mission)
	if err != nil {
		return 0, err
	}
	if mission.Id == 0 {
		return 0, errors.New("初始资本任务不存在！")
	}

	balanceDetail := &model.UserBalanceDetail{
		Uid:     user.Uid,
		Type:    model.MissionTypeInitial,
		Num:     mission.Fixed,
		Balance: mission.Fixed,
		Desc:    fmt.Sprintf("获得%s %d 铜币", model.BalanceTypeMap[mission.Type], mission.Fixed),
	}
	_, err = session.Insert(balanceDetail)

	return mission.Fixed, err
}
