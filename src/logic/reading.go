// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"errors"
	"model"
	"net/url"
	"strconv"
	"strings"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

type ReadingLogic struct{}

var DefaultReading = ReadingLogic{}

func (ReadingLogic) FindLastList(beginTime string) ([]*model.MorningReading, error) {
	readings := make([]*model.MorningReading, 0)
	err := MasterDB.Where("ctime>? AND rtype=0", beginTime).OrderBy("id DESC").Find(&readings)

	return readings, err
}

// 获取晨读列表（分页）
func (ReadingLogic) FindBy(ctx context.Context, limit, rtype int, lastIds ...int) []*model.MorningReading {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.Where("rtype=?", rtype)
	if len(lastIds) > 0 && lastIds[0] > 0 {
		dbSession.And("id<?", lastIds[0])
	}

	readingList := make([]*model.MorningReading, 0)
	err := dbSession.OrderBy("id DESC").Limit(limit).Find(&readingList)
	if err != nil {
		objLog.Errorln("ResourceLogic FindReadings Error:", err)
		return nil
	}

	return readingList
}

// 【我要晨读】
func (ReadingLogic) IReading(ctx context.Context, id int) string {
	objLog := GetLogger(ctx)

	reading := &model.MorningReading{}
	_, err := MasterDB.Id(id).Get(reading)
	if err != nil {
		objLog.Errorln("reading logic IReading error:", err)
		return "/readings"
	}

	if reading.Id == 0 {
		return "/readings"
	}

	go MasterDB.Id(id).Incr("clicknum", 1).Update(reading)

	if reading.Inner == 0 {
		return "/wr?u=" + reading.Url
	}

	return "/articles/" + strconv.Itoa(reading.Inner)
}

// FindReadingByPage 获取晨读列表（分页）
func (ReadingLogic) FindReadingByPage(ctx context.Context, conds map[string]string, curPage, limit int) ([]*model.MorningReading, int) {
	objLog := GetLogger(ctx)

	session := MasterDB.NewSession()

	for k, v := range conds {
		session.And(k+"=?", v)
	}

	totalSession := session.Clone()

	offset := (curPage - 1) * limit
	readingList := make([]*model.MorningReading, 0)
	err := session.OrderBy("id DESC").Limit(limit, offset).Find(&readingList)
	if err != nil {
		objLog.Errorln("reading find error:", err)
		return nil, 0
	}

	total, err := totalSession.Count(new(model.MorningReading))
	if err != nil {
		objLog.Errorln("reading find count error:", err)
		return nil, 0
	}

	return readingList, int(total)
}

// SaveReading 保存晨读
func (ReadingLogic) SaveReading(ctx context.Context, form url.Values, username string) (errMsg string, err error) {
	reading := &model.MorningReading{}
	err = schemaDecoder.Decode(reading, form)
	if err != nil {
		logger.Errorln("reading SaveReading error", err)
		errMsg = err.Error()
		return
	}

	readings := make([]*model.MorningReading, 0)
	if reading.Inner != 0 {
		reading.Url = ""
		err = MasterDB.Where("`inner`=?", reading.Inner).OrderBy("id DESC").Find(&readings)
	} else {
		err = MasterDB.Where("url=?", reading.Url).OrderBy("id DESC").Find(&readings)
	}
	if err != nil {
		logger.Errorln("reading SaveReading MasterDB.Where() error", err)
		errMsg = err.Error()
		return
	}

	reading.Moreurls = strings.TrimSpace(reading.Moreurls)
	if strings.Contains(reading.Moreurls, "\n") {
		reading.Moreurls = strings.Join(strings.Split(reading.Moreurls, "\n"), ",")
	}

	reading.Username = username

	logger.Debugln(reading.Rtype, "id=", reading.Id)
	if reading.Id != 0 {
		_, err = MasterDB.Id(reading.Id).Update(reading)
	} else {
		if len(readings) > 0 {
			logger.Errorln("reading report:", reading)
			errMsg, err = "已经存在了!!", errors.New("已经存在了!!")
			return
		}
		_, err = MasterDB.Insert(reading)
	}

	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln("reading save:", errMsg, ":", err)
		return
	}

	return
}

// FindById 获取单条晨读
func (ReadingLogic) FindById(ctx context.Context, id int) *model.MorningReading {
	reading := &model.MorningReading{}
	_, err := MasterDB.Id(id).Get(reading)
	if err != nil {
		logger.Errorln("reading logic FindReadingById Error:", err)
		return nil
	}

	if reading.Id == 0 {
		return nil
	}

	return reading
}
