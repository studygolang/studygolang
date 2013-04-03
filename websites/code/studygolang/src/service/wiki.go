// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"model"
	"net/url"
	"strconv"
	"strings"
	"util"
)

// 创建一个wiki页面
func CreateWiki(uid int, form url.Values) bool {
	wiki := model.NewWiki()
	err := util.ConvertAssign(wiki, form)
	if err != nil {
		logger.Errorln("wiki ConvertAssign error", err)
		return false
	}
	wiki.Uid = uid
	if _, err := wiki.Insert(); err != nil {
		logger.Errorln("wiki service CreateWiki error:", err)
		return false
	}

	// 创建一个wiki页面，活跃度+10
	go IncUserWeight("uid="+strconv.Itoa(uid), 10)

	return true
}

// 获得wiki列表（TODO：暂时不分页）
func FindWikiList() []*model.Wiki {
	wikiList, err := model.NewWiki().Order("mtime DESC").FindAll("title", "uri")
	if err != nil {
		logger.Errorln("wiki service FindWikiList error:", err)
		return nil
	}
	return wikiList
}

// 某个wiki页面详细信息
func FindWiki(uri string) map[string]interface{} {
	wiki := model.NewWiki()
	if err := wiki.Where("uri=" + uri).Find(); err != nil {
		logger.Errorln("wiki service FindWiki error:", err)
		return nil
	}
	uids := make(map[int]int)
	uids[wiki.Uid] = wiki.Uid
	if wiki.Cuid != "" {
		cuids := strings.Split(wiki.Cuid, ",")
		for _, cuid := range cuids {
			tmpUid := util.MustInt(cuid)
			uids[tmpUid] = tmpUid
		}
	}
	userMap := getUserInfos(uids)
	result := make(map[string]interface{})
	util.Struct2Map(result, wiki)
	result["user"] = userMap[wiki.Uid]
	if wiki.Cuid != "" {
		cuids := strings.Split(wiki.Cuid, ",")
		cusers := make([]*model.User, len(cuids))
		for i, cuid := range cuids {
			cusers[i] = userMap[util.MustInt(cuid)]
		}
		result["cuser"] = cusers
	}
	return result
}

// 通过id获得wiki的所有者
func getWikiOwner(id int) int {
	wiki := model.NewWiki()
	err := wiki.Where("id=" + strconv.Itoa(id)).Find()
	if err != nil {
		logger.Errorln("wiki service getWikiOwner Error:", err)
		return 0
	}
	return wiki.Uid
}

// 获取多个wiki页面详细信息
func FindWikisByIds(ids []int) []*model.Wiki {
	if len(ids) == 0 {
		return nil
	}
	inIds := util.Join(ids, ",")
	wikis, err := model.NewWiki().Where("id in(" + inIds + ")").FindAll()
	if err != nil {
		logger.Errorln("wiki service FindWikisByIds error:", err)
		return nil
	}
	return wikis
}

// 提供给其他service调用（包内）
func getWikis(ids map[int]int) map[int]*model.Wiki {
	wikis := FindWikisByIds(util.MapIntKeys(ids))
	wikiMap := make(map[int]*model.Wiki, len(wikis))
	for _, wiki := range wikis {
		wikiMap[wiki.Id] = wiki
	}
	return wikiMap
}
