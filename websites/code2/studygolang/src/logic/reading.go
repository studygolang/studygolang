// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"
	"strconv"

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

// // 获取晨读列表（分页）
// func FindReadingByPage(conds map[string]string, curPage, limit int) ([]*model.MorningReading, int) {
// 	conditions := make([]string, 0, len(conds))
// 	for k, v := range conds {
// 		conditions = append(conditions, k+"="+v)
// 	}

// 	reading := model.NewMorningReading()

// 	limitStr := strconv.Itoa((curPage-1)*limit) + "," + strconv.Itoa(limit)
// 	readingList, err := reading.Where(strings.Join(conditions, " AND ")).Order("id DESC").Limit(limitStr).
// 		FindAll()
// 	if err != nil {
// 		logger.Errorln("reading service FindArticleByPage Error:", err)
// 		return nil, 0
// 	}

// 	total, err := reading.Count()
// 	if err != nil {
// 		logger.Errorln("reading service FindReadingByPage COUNT Error:", err)
// 		return nil, 0
// 	}

// 	return readingList, total
// }

// // 保存晨读
// func SaveReading(form url.Values, username string) (errMsg string, err error) {
// 	reading := model.NewMorningReading()
// 	err = util.ConvertAssign(reading, form)
// 	if err != nil {
// 		logger.Errorln("reading SaveReading error", err)
// 		errMsg = err.Error()
// 		return
// 	}

// 	reading.Moreurls = strings.TrimSpace(reading.Moreurls)
// 	if strings.Contains(reading.Moreurls, "\n") {
// 		reading.Moreurls = strings.Join(strings.Split(reading.Moreurls, "\n"), ",")
// 	}

// 	reading.Username = username

// 	logger.Infoln(reading.Rtype, "id=", reading.Id)
// 	if reading.Id != 0 {
// 		err = reading.Persist(reading)
// 	} else {
// 		_, err = reading.Insert()
// 	}

// 	if err != nil {
// 		errMsg = "内部服务器错误"
// 		logger.Errorln("reading save:", errMsg, ":", err)
// 		return
// 	}

// 	return
// }

// // 获取单条晨读
// func FindReadingById(id int) (*model.MorningReading, error) {
// 	reading := model.NewMorningReading()
// 	err := reading.Where("id=?", id).Find()
// 	if err != nil {
// 		logger.Errorln("reading service FindReadingById Error:", err)
// 	}

// 	return reading, err
// }
