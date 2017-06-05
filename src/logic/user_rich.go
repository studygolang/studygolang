// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
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
func (UserRichLogic) IncrUserRich(user *model.User, typ, award int, desc string) {
	session := MasterDB.NewSession()
	defer session.Close()

	session.Begin()

	user.Balance += award
	_, err := session.Where("uid=?", user.Uid).Cols("balance").Update(user)
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
