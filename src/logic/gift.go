// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	"errors"
	"fmt"
	"model"
	"time"

	. "db"

	"github.com/go-xorm/xorm"
)

type GiftLogic struct{}

var DefaultGift = GiftLogic{}

func (self GiftLogic) FindAllOnline(ctx context.Context) []*model.Gift {
	objLog := GetLogger(ctx)

	gifts := make([]*model.Gift, 0)
	err := MasterDB.Where("state=?", model.GiftStateOnline).Find(&gifts)
	if err != nil {
		objLog.Errorln("GiftLogic FindAllOnline error:", err)
		return nil
	}

	for _, gift := range gifts {
		if gift.ExpireTime.Before(time.Now()) {
			gift.State = model.GiftStateExpired
			go self.doExpire(gift)
		}
	}

	return gifts
}

func (self GiftLogic) Exchange(ctx context.Context, me *model.Me, giftId int) error {
	objLog := GetLogger(ctx)

	gift := &model.Gift{}
	_, err := MasterDB.Id(giftId).Get(gift)
	if err != nil {
		objLog.Errorln("GiftLogic Exchange error:", err)
		return err
	}

	if gift.RemainNum == 0 {
		return errors.New("已兑完")
	}

	total, err := MasterDB.Where("gift_id=? AND uid=?", giftId, me.Uid).Count(new(model.UserExchangeRecord))
	if err != nil {
		objLog.Errorln("GiftLogic Count UserExchangeRecord error:", err)
		return err
	}

	if gift.BuyLimit <= int(total) {
		return errors.New("已兑换过")
	}

	if gift.Typ == model.GiftTypRedeem {
		return self.exchangeRedeem(gift, me)
	} else if gift.Typ == model.GiftTypDiscount {
		return self.exchangeDiscount(gift, me)
	}

	return nil
}

func (self GiftLogic) FindExchangeRecords(ctx context.Context, me *model.Me) []*model.UserExchangeRecord {
	objLog := GetLogger(ctx)

	records := make([]*model.UserExchangeRecord, 0)
	err := MasterDB.Where("uid=?", me.Uid).Desc("id").Find(&records)
	if err != nil {
		objLog.Errorln("GiftLogic FindExchangeRecords error:", err)
		return nil
	}

	return records
}

func (self GiftLogic) UserCanExchange(ctx context.Context, me *model.Me, gifts []*model.Gift) {
	num := len(gifts)
	if num == 0 {
		return
	}
	objLog := GetLogger(ctx)

	giftIds := make([]int, num)
	for i, gift := range gifts {
		giftIds[i] = gift.Id
	}

	exchangeRecords := make([]*model.UserExchangeRecord, 0)
	err := MasterDB.In("gift_id", giftIds).And("uid=?", me.Uid).Find(&exchangeRecords)
	if err != nil {
		objLog.Errorln("GiftLogic FindUserGifts error:", err)
		return
	}
	for _, record := range exchangeRecords {
		for _, gift := range gifts {
			if record.GiftId == gift.Id {
				gift.BuyLimit--
				break
			}
		}
	}
}

func (self GiftLogic) exchangeRedeem(gift *model.Gift, me *model.Me) error {
	giftRedeem := &model.GiftRedeem{}
	_, err := MasterDB.Where("gift_id=? AND exchange=0", gift.Id).Get(giftRedeem)
	if err != nil {
		return err
	}

	if giftRedeem.Id == 0 {
		return errors.New("no more gift redeem")
	}

	return self.doExchange(gift, me, "兑换码："+giftRedeem.Code, func(session *xorm.Session) error {
		_, err := session.Table(giftRedeem).Where("id=? AND exchange=0", giftRedeem.Id).
			Update(map[string]interface{}{"exchange": 1, "uid": me.Uid})

		return err
	})
}

func (self GiftLogic) exchangeDiscount(gift *model.Gift, me *model.Me) error {
	return self.doExchange(gift, me, "已兑换，我们会尽快联系合作方处理", nil)
}

func (self GiftLogic) doExchange(gift *model.Gift, me *model.Me, remark string, moreOp func(session *xorm.Session) error) error {
	if me.Balance < gift.Price {
		return errors.New("兑换失败：铜币不够！")
	}

	session := MasterDB.NewSession()
	defer session.Close()

	session.Begin()

	exchangeRecord := &model.UserExchangeRecord{
		GiftId:     gift.Id,
		Uid:        me.Uid,
		Remark:     remark,
		ExpireTime: gift.ExpireTime,
	}
	_, err := MasterDB.Insert(exchangeRecord)
	if err != nil {
		session.Rollback()
		return err
	}

	if moreOp != nil {
		err = moreOp(session)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	_, err = session.Id(gift.Id).Decr("remain_num", 1).Update(new(model.Gift))
	if err != nil {
		session.Rollback()
		return err
	}

	desc := fmt.Sprintf("兑换 %s 消费 %d 铜币", gift.Name, gift.Price)
	err = DefaultMission.changeUserBalance(session, me, model.MissionTypeGift, -gift.Price, desc)
	if err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

func (self GiftLogic) doExpire(gift *model.Gift) {
	MasterDB.Table(gift).Where("id=?", gift.Id).Update(map[string]interface{}{"state": gift.State})
}
