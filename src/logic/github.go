// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"fmt"
	"io/ioutil"
	"model"
	"net/http"
	"strings"
	"time"

	"github.com/polaris1119/logger"
	"github.com/tidwall/gjson"
)

type GithubLogic struct{}

var DefaultGithub = GithubLogic{}

func (self GithubLogic) PullPR(repo string) error {
	prURL := fmt.Sprintf("%s/repos/%s/pulls?state=all&per_page=30", GithubAPIBaseUrl, repo)
	resp, err := http.Get(prURL)
	if err != nil {
		logger.Errorln("GithubLogic PullPR get error:", err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		logger.Errorln("GithubLogic PullPR read all error:", err)
		return err
	}

	var outErr error

	result := gjson.ParseBytes(body)
	result.ForEach(func(key, val gjson.Result) bool {
		thePRURL := val.Get("url").String()

		// 没有 merge 的忽略
		if val.Get("merged_at").Type == gjson.Null {
			logger.Infoln("The pr", val.Get("number"), "title:", val.Get("title"), "is not merged", "url:", thePRURL)
			return true
		}

		prTime := val.Get("created_at").Time()
		username := val.Get("user.login").String()
		avatar := val.Get("user.avatar_url").String()

		filesURL := thePRURL + "/files"
		resp, err = http.Get(filesURL)
		if err != nil {
			outErr = err
			logger.Errorln("github fetch files error:", err, "url:", filesURL)
			return true
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			outErr = err
			logger.Errorln("github read files resp error:", err)
			return true
		}
		filesResult := gjson.ParseBytes(body)

		err = self.dealFiles(username, avatar, prTime, filesResult)
		if err != nil {
			outErr = err
		}

		return true
	})

	// stat gctt user time
	self.statUserTime()

	return outErr
}

func (self GithubLogic) dealFiles(username, avatar string, prTime time.Time, filesResult gjson.Result) error {
	var outErr error
	filesResult.ForEach(func(key, val gjson.Result) bool {
		status := val.Get("status").String()
		filename := val.Get("filename").String()
		filenames := strings.SplitN(filename, "/", 3)
		fmt.Println(filename)
		if len(filenames) < 3 {
			return true
		}
		title := strings.Split(filenames[2], ".")[0]
		fmt.Println(title)

		// 对原文的改动
		if strings.HasPrefix(filename, "sources") {
			if status == "modified" {
				// 认为是开始翻译
				err := self.insertOrUpdateGCCT(username, avatar, title, prTime.Unix(), 0)
				if err != nil {
					outErr = err
				}
			}
			return true
		}

		// 翻译完成
		if strings.HasPrefix(filename, "translated") {
			if status == "added" {
				err := self.insertOrUpdateGCCT(username, avatar, title, 0, prTime.Unix())
				if err != nil {
					outErr = err
				}
			}
			return true
		}

		return true
	})

	return outErr
}

func (GithubLogic) insertOrUpdateGCCT(username, avatar, title string, translating, translated int64) error {
	gcttGit := &model.GCTTGit{}
	_, err := MasterDB.Where("username=? AND title=?", username, title).Get(gcttGit)
	if err != nil {
		logger.Errorln("GithubLogic insertOrUpdateGCCT get error:", err)
		return err
	}

	gcttUser := DefaultGCTT.FindOne(nil, username)

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	if gcttUser.Id == 0 {
		gcttUser.Username = username
		gcttUser.Avatar = avatar
		gcttUser.JoinedAt = translating
		_, err = session.Insert(gcttUser)
		if err != nil {
			session.Rollback()
			logger.Errorln("GithubLogic insertOrUpdateGCCT insert gctt_user error:", err)
			return err
		}
	}

	// 已经存在
	if gcttGit.Id > 0 {
		if gcttGit.TranslatedAt == 0 && translated != 0 {
			gcttGit.TranslatedAt = translated
			_, err = MasterDB.Id(gcttGit.Id).Update(gcttGit)
			if err != nil {
				session.Rollback()
				logger.Errorln("GithubLogic insertOrUpdateGCCT update error:", err)
				return err
			}
		}

		session.Commit()
		return nil
	}

	gcttGit.Username = username
	gcttGit.Title = title
	gcttGit.TranslatingAt = translating
	gcttGit.TranslatedAt = translated
	_, err = MasterDB.Insert(gcttGit)
	if err != nil {
		session.Rollback()
		logger.Errorln("GithubLogic insertOrUpdateGCCTGit insert error:", err)
		return err
	}

	session.Commit()
	return nil
}

func (GithubLogic) statUserTime() {
	gcttUsers := make([]*model.GCTTUser, 0)
	err := MasterDB.Find(&gcttUsers)
	if err != nil {
		logger.Errorln("GithubLogic statUserTime find error:", err)
		return
	}

	for _, gcttUser := range gcttUsers {
		gcttGits := make([]*model.GCTTGit, 0)
		err = MasterDB.Where("username=?", gcttUser.Username).OrderBy("id ASC").Find(&gcttGits)
		if err != nil {
			logger.Errorln("GithubLogic find gctt git error:", err)
			continue
		}

		var avgTime, lastAt int64
		for _, gcttGit := range gcttGits {
			if gcttGit.TranslatingAt != 0 && gcttGit.TranslatedAt != 0 {
				avgTime += gcttGit.TranslatedAt - gcttGit.TranslatingAt
			}

			if gcttGit.TranslatedAt > lastAt {
				lastAt = gcttGit.TranslatedAt
			}
		}

		gcttUser.Num = len(gcttGits)
		gcttUser.AvgTime = int(avgTime)
		gcttUser.LastAt = lastAt
		_, err = MasterDB.Id(gcttUser.Id).Update(gcttUser)
		if err != nil {
			logger.Errorln("GithubLogic update gctt user error:", err)
		}
	}
}
