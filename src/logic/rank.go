// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"fmt"
	"time"

	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
	"github.com/polaris1119/times"
)

type RankLogic struct{}

var DefaultRank = RankLogic{}

func (self RankLogic) GenDayRank(objtype, objid, num int) {
	redisClient := nosql.NewRedisClient()
	key := self.getDayRankKey(objtype, times.Format("ymd"))
	err := redisClient.ZINCRBY(key, num, objid)
	if err != nil {
		logger.Errorln("view redis ZINCRBY error:", err)
	}
	redisClient.EXPIRE(key, 2*30*86400)
}

// GenWeekRank 过去 7 天排行榜
func (self RankLogic) GenWeekRank(objtype int) {
	redisClient := nosql.NewRedisClient()
	dest := self.getWeekRankKey(objtype)

	keys := self.getMultiKey(objtype, 7)

	err := redisClient.ZUNIONSTORE(dest, 7, keys, nil)
	if err != nil {
		logger.Errorln("GenWeekRank ZUNIONSTORE error:", err)
	}
}

// GenMonthRank 过去 30 天排行榜
func (self RankLogic) GenMonthRank(objtype int) {
	redisClient := nosql.NewRedisClient()
	dest := self.getMonthRankKey(objtype)

	keys := self.getMultiKey(objtype, 30)

	err := redisClient.ZUNIONSTORE(dest, 30, keys, nil)
	if err != nil {
		logger.Errorln("GenMonthRank ZUNIONSTORE error:", err)
	}
}

func (self RankLogic) getMultiKey(objtype, num int) []string {
	today := time.Now()

	keys := make([]string, num)
	for i := 0; i < num; i++ {
		ymd := times.Format("ymd", today.Add(time.Duration(-(i+1)*86400)*time.Second))
		keys[i] = self.getDayRankKey(objtype, ymd)
	}

	return keys
}

func (RankLogic) getDayRankKey(objtype int, ymd string) string {
	return fmt.Sprintf("view:type-%d:rank:%s", objtype, ymd)
}

func (RankLogic) getWeekRankKey(objtype int) string {
	return fmt.Sprintf("view:type-%d:rank:last-week", objtype)
}

func (RankLogic) getMonthRankKey(objtype int) string {
	return fmt.Sprintf("view:type-%d:rank:last-month", objtype)
}
