// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	"fmt"
	"model"
	"time"

	. "db"

	"github.com/garyburd/redigo/redis"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
	"github.com/polaris1119/times"
)

type RankLogic struct{}

var DefaultRank = RankLogic{}

func (self RankLogic) GenDayRank(objtype, objid, num int) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()
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
	defer redisClient.Close()
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
	defer redisClient.Close()
	dest := self.getMonthRankKey(objtype)

	keys := self.getMultiKey(objtype, 30)

	err := redisClient.ZUNIONSTORE(dest, 30, keys, nil)
	if err != nil {
		logger.Errorln("GenMonthRank ZUNIONSTORE error:", err)
	}
}

// GenDAURank 生成日活跃用户排行
func (self RankLogic) GenDAURank(uid, weight int) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()
	key := self.getDAURankKey(times.Format("ymd"))
	err := redisClient.ZINCRBY(key, weight, uid)
	if err != nil {
		logger.Errorln("dau redis ZINCRBY error:", err)
	}
	redisClient.EXPIRE(key, 2*30*86400)
}

// FindDayRank needExt 是否需要扩展数据
func (self RankLogic) FindDayRank(ctx context.Context, objtype int, ymd string, num int, needExt ...bool) (result interface{}) {
	objLog := GetLogger(ctx)

	redisClient := nosql.NewRedisClient()
	key := self.getDayRankKey(objtype, ymd)
	resultSlice, err := redisClient.ZREVRANGE(key, 0, num-1, true)
	redisClient.Close()
	if err != nil {
		objLog.Errorln("FindDayRank ZREVRANGE error:", err)
		return nil
	}

	return self.findModelsByRank(resultSlice, objtype, num, needExt...)
}

func (self RankLogic) FindWeekRank(ctx context.Context, objtype, num int, needExt ...bool) (result interface{}) {
	objLog := GetLogger(ctx)

	redisClient := nosql.NewRedisClient()
	key := self.getWeekRankKey(objtype)
	resultSlice, err := redisClient.ZREVRANGE(key, 0, num-1, true)
	redisClient.Close()
	if err != nil {
		objLog.Errorln("FindWeekRank ZREVRANGE error:", err)
		return nil
	}

	return self.findModelsByRank(resultSlice, objtype, num, needExt...)
}

func (self RankLogic) FindMonthRank(ctx context.Context, objtype, num int, needExt ...bool) (result interface{}) {
	objLog := GetLogger(ctx)

	redisClient := nosql.NewRedisClient()
	key := self.getMonthRankKey(objtype)
	resultSlice, err := redisClient.ZREVRANGE(key, 0, num-1, true)
	redisClient.Close()
	if err != nil {
		objLog.Errorln("FindMonthRank ZREVRANGE error:", err)
		return nil
	}

	return self.findModelsByRank(resultSlice, objtype, num, needExt...)
}

// FindDAURank DAU 排名，默认获取当天的
func (self RankLogic) FindDAURank(ctx context.Context, num int, ymds ...string) []*model.User {
	objLog := GetLogger(ctx)

	ymd := times.Format("ymd")
	if len(ymds) > 0 {
		ymd = ymds[0]
	}

	redisClient := nosql.NewRedisClient()
	key := self.getDAURankKey(ymd)
	resultSlice, err := redisClient.ZREVRANGE(key, 0, num-1, true)
	redisClient.Close()
	if err != nil {
		objLog.Errorln("FindDAURank ZREVRANGE error:", err)
		return nil
	}

	uids := make([]int, 0, num)
	weights := make([]int, 0, num)

	for len(resultSlice) > 0 {
		var (
			uid, weight int
			err         error
		)
		resultSlice, err = redis.Scan(resultSlice, &uid, &weight)
		if err != nil {
			logger.Errorln("FindDAURank redis Scan error:", err)
			return nil
		}

		uids = append(uids, uid)
		weights = append(weights, weight)
	}

	if len(uids) == 0 {
		return nil
	}

	userMap := DefaultUser.FindDAUUsers(ctx, uids)
	users := make([]*model.User, len(userMap))
	for i, uid := range uids {
		user := userMap[uid]
		user.Weight = weights[i]
		users[i] = user
	}

	return users
}

// TotalDAUUser 今日活跃用户数
func (self RankLogic) TotalDAUUser(ctx context.Context) int {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	key := self.getDAURankKey(times.Format("ymd"))
	return redisClient.ZCARD(key)
}

func (self RankLogic) UserDAURank(ctx context.Context, uid int) int {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	key := self.getDAURankKey(times.Format("ymd"))
	return redisClient.ZREVRANK(key, uid)
}

// FindRichRank 社区财富排行榜
func (self RankLogic) FindRichRank(ctx context.Context) []*model.User {
	objLog := GetLogger(ctx)

	userList := make([]*model.User, 0)
	err := MasterDB.Where("balance>?", 0).Desc("balance").Limit(25).Find(&userList)
	if err != nil {
		objLog.Errorln("find rich rank error:", err)
		return nil
	}

	return userList
}

func (RankLogic) findModelsByRank(resultSlice []interface{}, objtype, num int, needExt ...bool) (result interface{}) {
	objids := make([]int, 0, num)
	viewNums := make([]int, 0, num)

	for len(resultSlice) > 0 {
		var (
			objid, viewNum int
			err            error
		)
		resultSlice, err = redis.Scan(resultSlice, &objid, &viewNum)
		if err != nil {
			logger.Errorln("findModelsByRank redis Scan error:", err)
			return nil
		}

		objids = append(objids, objid)
		viewNums = append(viewNums, viewNum)
	}

	switch objtype {
	case model.TypeTopic:
		if len(needExt) > 0 && needExt[0] {
			topics := DefaultTopic.FindFullinfoByTids(objids)
			for i, topic := range topics {
				topic["rank_view"] = viewNums[i]
			}
			result = topics
		} else {
			topics := DefaultTopic.FindByTids(objids)
			for i, topic := range topics {
				topic.RankView = viewNums[i]
			}
			result = topics
		}
	case model.TypeResource:
		resources := DefaultResource.FindByIds(objids)
		for i, resource := range resources {
			resource.RankView = viewNums[i]
		}
		result = resources
	case model.TypeArticle:
		articles := DefaultArticle.FindByIds(objids)
		for i, article := range articles {
			article.RankView = viewNums[i]
		}
		result = articles
	case model.TypeProject:
		projects := DefaultProject.FindByIds(objids)
		for i, project := range projects {
			project.RankView = viewNums[i]
		}
		result = projects
	case model.TypeBook:
		books := DefaultGoBook.FindByIds(objids)
		for i, book := range books {
			book.RankView = viewNums[i]
		}
		result = books
	}

	return
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

func (RankLogic) getDAURankKey(ymd string) string {
	return fmt.Sprintf("dau:rank:%s", ymd)
}
