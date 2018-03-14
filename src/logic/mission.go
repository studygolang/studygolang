// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"model"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/times"
	"golang.org/x/net/context"
)

type MissionLogic struct{}

var DefaultMission = MissionLogic{}

// HasLoginMission 是否有今日登录奖励
func (self MissionLogic) HasLoginMission(ctx context.Context, me *model.Me) bool {
	// 还没有铜币，当然有可能是消耗尽了
	if me.Balance == 0 {
		// 初始资本没有领取，必须先领取
		if DefaultUserRich.Total(ctx, me.Uid) == 0 {
			return false
		}
	}

	userLoginMission := self.FindLoginMission(ctx, me)
	if userLoginMission == nil {
		return false
	}

	if userLoginMission.Uid == 0 {
		return true
	}

	// 今日是否领取
	if times.Format("Ymd") == strconv.Itoa(userLoginMission.Date) {
		return false
	}
	return true
}

// RedeemLoginAward 领取登录奖励
func (self MissionLogic) RedeemLoginAward(ctx context.Context, me *model.Me) error {
	objLog := GetLogger(ctx)

	mission := self.findMission(ctx, model.MissionTypeLogin)
	if mission.Id == 0 {
		objLog.Errorln("每日登录任务不存在")
		return errors.New("任务不存在")
	}

	userLoginMission := self.FindLoginMission(ctx, me)
	if userLoginMission == nil {
		objLog.Errorln("查询数据库失败")
		return errors.New("服务内部错误")
	}

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	if userLoginMission.Uid == 0 {
		userLoginMission.Date = goutils.MustInt(time.Now().Format("20060102"))
		userLoginMission.Days = 1
		userLoginMission.TotalDays = 1
		userLoginMission.Award = mission.Min
		userLoginMission.Uid = me.Uid

		_, err := session.Insert(userLoginMission)
		if err != nil {
			session.Rollback()
			objLog.Errorln("insert user_login_mission error:", err)
			return errors.New("服务内部错误")
		}

	} else {
		today := goutils.MustInt(times.Format("Ymd"))
		if today == userLoginMission.Date {
			session.Rollback()
			return errors.New("今日已领取")
		}
		// 昨日是否领取了
		yesterday := goutils.MustInt(times.Format("Ymd", time.Now().Add(-86400*time.Second)))
		if yesterday != userLoginMission.Date {
			userLoginMission.Award = mission.Min
			userLoginMission.Days = 1
		} else {
			userLoginMission.Days++
			if userLoginMission.Award == mission.Max {
				userLoginMission.Award = mission.Min
			} else {
				award := userLoginMission.Award + rand.Intn(mission.Incr) + 1
				userLoginMission.Award = int(math.Min(float64(award), float64(mission.Max)))
			}
		}

		userLoginMission.Date = today
		userLoginMission.TotalDays++
		userLoginMission.UpdatedAt = time.Now()

		_, err := session.Where("uid=?", userLoginMission.Uid).Update(userLoginMission)
		if err != nil {
			session.Rollback()
			objLog.Errorln("update user_login_mission error:", err)
			return errors.New("服务内部错误")
		}
	}

	desc := times.Format("Ymd") + " 的每日登录奖励 " + strconv.Itoa(userLoginMission.Award) + " 铜币"
	err := self.changeUserBalance(session, me, model.MissionTypeLogin, userLoginMission.Award, desc)
	if err != nil {
		session.Rollback()
		objLog.Errorln("changeUserBalance error:", err)
		return errors.New("服务内部错误")
	}

	session.Commit()

	return nil
}

func (MissionLogic) FindLoginMission(ctx context.Context, me *model.Me) *model.UserLoginMission {
	objLog := GetLogger(ctx)

	userLoginMission := &model.UserLoginMission{}
	_, err := MasterDB.Where("uid=?", me.Uid).Get(userLoginMission)
	if err != nil {
		objLog.Errorln("MissionLogic FindLoginMission error:", err)
		return nil
	}

	return userLoginMission
}

// Complete 完成任务（非每日任务）
func (MissionLogic) Complete(ctx context.Context, me *model.Me, id string) error {
	objLog := GetLogger(ctx)

	mission := &model.Mission{}
	_, err := MasterDB.Id(id).Get(mission)
	if err != nil {
		objLog.Errorln("MissionLogic FindLoginMission error:", err)
		return err
	}

	if mission.Id == 0 || mission.State != 0 {
		return errors.New("任务不存在或已过期")
	}

	user := DefaultUser.FindOne(ctx, "uid", me.Uid)

	// 初始任务，不允许重复提交
	if id == strconv.Itoa(model.InitialMissionId) {
		if user.Balance > 0 {
			objLog.Errorln("repeat claim init award", user.Username)
			return nil
		}

		details := DefaultUserRich.FindBalanceDetail(ctx, me, mission.Type)
		if len(details) > 0 {
			return nil
		}
	}

	desc := fmt.Sprintf("获得%s %d 铜币", model.BalanceTypeMap[mission.Type], mission.Fixed)
	DefaultUserRich.IncrUserRich(user, mission.Type, mission.Fixed, desc)

	return nil
}

func (MissionLogic) findMission(ctx context.Context, typ int) *model.Mission {
	mission := &model.Mission{}
	MasterDB.Where("type=?", typ).Get(mission)
	return mission
}

func (self MissionLogic) changeUserBalance(session *xorm.Session, me *model.Me, typ, award int, desc string) error {
	_, err := session.Where("uid=?", me.Uid).Incr("balance", award).Update(new(model.User))
	if err != nil {
		return errors.New("服务内部错误")
	}

	balanceDetail := &model.UserBalanceDetail{
		Uid:     me.Uid,
		Type:    typ,
		Num:     award,
		Balance: me.Balance + award,
		Desc:    desc,
	}
	return DefaultUserRich.add(session, balanceDetail)
}
