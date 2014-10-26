// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"strconv"

	"logger"
	"model"
)

// 获取晨读列表（分页）
func FindReadings(lastId, limit string) []*model.MorningReading {
	reading := model.NewMorningReading()

	cond := ""
	if lastId != "0" {
		cond = " AND id<" + lastId
	}

	readingList, err := reading.Where(cond).Order("id DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("reading service FindReadings Error:", err)
		return nil
	}

	return readingList
}

// 【我要晨读】
func IReading(id string) string {
	_, err := strconv.Atoi(id)
	if err != nil {
		return "/readings"
	}

	reading := model.NewMorningReading()
	err = reading.Where("id=?", id).Find()

	if err != nil {
		logger.Errorln("reading service IReading error:", err)
		return "/readings"
	}

	if reading.Id == 0 {
		return "/readings"
	}

	reading.Where("id=?", id).Increment("clicknum", 1)

	return reading.Url
}
