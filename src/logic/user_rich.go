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
	"net/url"
	"time"
	"util"

	. "db"

	"github.com/garyburd/redigo/redis"
	"github.com/go-xorm/xorm"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
	"github.com/polaris1119/times"
	"golang.org/x/net/context"
)

var (
	beginAwardWeight = 50
)

type UserRichLogic struct{}

var DefaultUserRich = UserRichLogic{}

func (self UserRichLogic) AwardCooper() {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()
	ymd := times.Format("ymd", time.Now().Add(-86400*time.Second))
	key := DefaultRank.getDAURankKey(ymd)

	var (
		cursor      uint64
		err         error
		resultSlice []interface{}
		count       = 20
	)

	for {
		cursor, resultSlice, err = redisClient.ZSCAN(key, cursor, "COUNT", count)
		if err != nil {
			logger.Errorln("AwardCooper ZSCAN error:", err)
			break
		}

		for len(resultSlice) > 0 {
			var (
				uid, weight int
				err         error
			)
			resultSlice, err = redis.Scan(resultSlice, &uid, &weight)
			if err != nil {
				logger.Errorln("AwardCooper redis Scan error:", err)
				continue
			}

			if weight < beginAwardWeight {
				continue
			}

			award := util.Max((weight-500)*5, 0) +
				util.UMin((weight-400), 100)*4 +
				util.UMin((weight-300), 100)*3 +
				util.UMin((weight-200), 100)*2 +
				util.UMin((weight-100), 100) +
				int(float64(util.UMin((weight-beginAwardWeight), beginAwardWeight))*0.5)

			userRank := redisClient.ZREVRANK(key, uid)
			desc := fmt.Sprintf("%s 的活跃度为 %d，排名第 %d，奖励 %d 铜币", ymd, weight, userRank, award)

			user := DefaultUser.FindOne(nil, "uid", uid)
			self.IncrUserRich(user, model.MissionTypeActive, award, desc)
		}

		if cursor == 0 {
			break
		}
	}
}

// IncrUserRich 增加或减少用户财富
func (self UserRichLogic) IncrUserRich(user *model.User, typ, award int, desc string) {
	if award == 0 {
		logger.Errorln("IncrUserRich, but award is empty!")
		return
	}

	var (
		total int64 = -1
		err   error
	)

	if award > 0 && (typ == model.MissionTypeReplied || typ == model.MissionTypeActive) {
		// 老用户，因为之前的主题被人回复而增加财富，自动帮其领取初始资本
		// 因为活跃奖励铜币，自动帮其领取初始资本
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
	if user.Balance < 0 {
		user.Balance = 0
	}
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

func (UserRichLogic) FindBalanceDetail(ctx context.Context, me *model.Me, types ...int) []*model.UserBalanceDetail {
	objLog := GetLogger(ctx)

	balanceDetails := make([]*model.UserBalanceDetail, 0)
	session := MasterDB.Where("uid=?", me.Uid)
	if len(types) > 0 {
		session.And("type=?", types[0])
	}

	err := session.Desc("id").Find(&balanceDetails)
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

func (self UserRichLogic) FindRecharge(ctx context.Context, me *model.Me) int {
	objLog := GetLogger(ctx)

	total, err := MasterDB.Where("uid=?", me.Uid).SumInt(new(model.UserRecharge), "amount")
	if err != nil {
		objLog.Errorln("UserRichLogic FindRecharge error:", err)
		return 0
	}

	return int(total)
}

// Recharge 用户充值
func (self UserRichLogic) Recharge(ctx context.Context, uid string, form url.Values) {
	objLog := GetLogger(ctx)

	createdAt, _ := time.ParseInLocation("2006-01-02 15:04:05", form.Get("time"), time.Local)
	userRecharge := &model.UserRecharge{
		Uid:       goutils.MustInt(uid),
		Amount:    goutils.MustInt(form.Get("amount")),
		Channel:   form.Get("channel"),
		CreatedAt: createdAt,
	}

	session := MasterDB.NewSession()
	session.Begin()

	_, err := session.Insert(userRecharge)
	if err != nil {
		session.Rollback()
		objLog.Errorln("UserRichLogic Recharge error:", err)
		return
	}

	user := DefaultUser.FindOne(ctx, "uid", uid)
	me := &model.Me{
		Uid:     user.Uid,
		Balance: user.Balance,
	}

	award := goutils.MustInt(form.Get("copper"))
	desc := fmt.Sprintf("%s 充值 ￥%d，获得 %d 个铜币", times.Format("Ymd"), userRecharge.Amount, award)
	err = DefaultMission.changeUserBalance(session, me, model.MissionTypeAdd, award, desc)
	if err != nil {
		session.Rollback()
		objLog.Errorln("UserRichLogic changeUserBalance error:", err)
		return
	}
	session.Commit()
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
