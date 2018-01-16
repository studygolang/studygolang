// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"errors"
	"fmt"
	"io/ioutil"
	"model"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/polaris1119/goutils"

	"github.com/polaris1119/logger"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
)

type GithubLogic struct{}

var DefaultGithub = GithubLogic{}

type prInfo struct {
	prURL    string
	username string
	avatar   string
	prTime   time.Time
	hadMerge bool
	number   int
}

var noMoreDataErr = errors.New("pull request: no more data")

func (self GithubLogic) PullRequestEvent(ctx context.Context, body []byte) error {
	objLog := GetLogger(ctx)

	result := gjson.ParseBytes(body)

	thePRURL := result.Get("pull_request.url").String()
	objLog.Infoln("GithubLogic PullRequestEvent, url:", thePRURL)

	_prInfo := &prInfo{
		prURL:    thePRURL,
		username: result.Get("pull_request.user.login").String(),
		avatar:   result.Get("pull_request.user.avatar_url").String(),
		prTime:   result.Get("pull_request.created_at").Time(),
		hadMerge: result.Get("pull_request.merged").Bool(),
	}

	err := self.dealFiles(_prInfo)

	objLog.Infoln("pull request deal successfully!")

	go self.statUserTime()

	return err
}

func (self GithubLogic) PullPR(repo string, isAll bool) error {
	if !isAll {
		err := self.pullPR(repo, 1)

		// stat gctt user time
		self.statUserTime()

		return err
	}

	var (
		err  error
		page = 1
	)

	for {
		err = self.pullPR(repo, page, "asc")
		if err == noMoreDataErr {
			break
		}

		page++
	}

	// stat gctt user time
	self.statUserTime()

	return err
}

func (self GithubLogic) pullPR(repo string, page int, directions ...string) error {
	prListURL := fmt.Sprintf("%s/repos/%s/pulls?state=all&per_page=30&page=%d", GithubAPIBaseUrl, repo, page)

	if len(directions) > 0 {
		prListURL += "&direction=" + directions[0]
	}

	resp, err := http.Get(prListURL)
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

	result := gjson.ParseBytes(body)

	if len(result.Array()) == 0 {
		return noMoreDataErr
	}

	var outErr error

	result.ForEach(func(key, val gjson.Result) bool {
		_prInfo := &prInfo{
			prURL:    val.Get("url").String(),
			username: val.Get("user.login").String(),
			avatar:   val.Get("user.avatar_url").String(),
			prTime:   val.Get("created_at").Time(),
			hadMerge: val.Get("merged_at").Type != gjson.Null,
			number:   int(val.Get("number").Int()),
		}

		err = self.dealFiles(_prInfo)
		if err != nil {
			outErr = err
		}

		return true
	})

	return outErr
}

func (self GithubLogic) dealFiles(_prInfo *prInfo) error {
	if _prInfo.prURL == "" {
		return nil
	}

	filesURL := _prInfo.prURL + "/files"
	resp, err := http.Get(filesURL)
	if err != nil {
		logger.Errorln("github fetch files error:", err, "url:", filesURL)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		logger.Errorln("github read files resp error:", err)
		return err
	}
	filesResult := gjson.ParseBytes(body)

	// 1. 领取翻译任务时，只是改变一个文件，且是 sources 目录下的，文件修改；
	// 2. 任务完成时，删除一个文件，创建一个新文件，删除的文件是 sources 目录下的，创建的文件是 translated 目录下的
	// 3. 翻译完成一篇，同时又领取新的一篇

	length := len(filesResult.Array())
	if length == 1 {
		err = self.translating(filesResult, _prInfo)
	} else if length == 2 {
		err = self.translated(filesResult, _prInfo)
	} else if length == 3 {
		err = self.translateSilmu(filesResult, _prInfo)
	}

	return err
}

func (self GithubLogic) translating(filesResult gjson.Result, _prInfo *prInfo) error {
	var outErr error
	filesResult.ForEach(func(key, val gjson.Result) bool {
		filename := val.Get("filename").String()
		// 是否对原文的改动
		if !strings.HasPrefix(filename, "sources") {

			// 目前改为采用 issue 的方式选题，不再有 sources
			if strings.HasPrefix(filename, "translated") {
				filenames := strings.SplitN(filename, "/", 3)
				if len(filenames) < 3 {
					return true
				}
				title := filenames[2]
				if title == "" {
					return true
				}

				err := self.issueTranslated(_prInfo, title)
				if err != nil {
					outErr = err
				}
			}

			return true
		}

		filenames := strings.SplitN(filename, "/", 3)
		if len(filenames) < 3 {
			return true
		}
		title := filenames[2]
		if title == "" {
			return true
		}

		// 认为是开始翻译
		status := val.Get("status").String()
		if status == "modified" && _prInfo.hadMerge {
			err := self.insertOrUpdateGCCT(_prInfo, title, false)
			if err != nil {
				outErr = err
			}
		}
		return true
	})

	return outErr
}

func (self GithubLogic) issueTranslated(_prInfo *prInfo, title string) error {
	md5 := goutils.Md5(title)
	gcttGit := &model.GCTTGit{}
	_, err := MasterDB.Where("md5=?", md5).Get(gcttGit)
	if err != nil {
		logger.Errorln("GithubLogic insertOrUpdateGCCT get error:", err)
		return err
	}

	if gcttGit.Id > 0 {
		return nil
	}

	gcttUser := DefaultGCTT.FindOne(nil, _prInfo.username)

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	if gcttUser.Id == 0 {
		gcttUser.Username = _prInfo.username
		gcttUser.Avatar = _prInfo.avatar
		gcttUser.JoinedAt = _prInfo.prTime.Unix()
		_, err = session.Insert(gcttUser)
		if err != nil {
			session.Rollback()
			logger.Errorln("GithubLogic issueTranslated insert gctt_user error:", err)
			return err
		}
	}

	gcttGit.Username = _prInfo.username
	gcttGit.Title = title
	gcttGit.Md5 = md5
	gcttGit.PR = _prInfo.number
	gcttGit.TranslatedAt = _prInfo.prTime.Unix()
	_, err = MasterDB.Insert(gcttGit)
	if err != nil {
		session.Rollback()
		logger.Errorln("GithubLogic issueTranslated insert error:", err)
		return err
	}

	session.Commit()
	return nil
}

func (self GithubLogic) translated(filesResult gjson.Result, _prInfo *prInfo) error {
	var (
		sourceTitle  string
		isTranslated = true
	)

	// 校验是否一个包含删除 sources 的操作，一个包含增加 translated 的操作
	filesResult.ForEach(func(key, val gjson.Result) bool {
		if !isTranslated {
			return false
		}

		status := val.Get("status").String()
		filename := val.Get("filename").String()

		if status == "removed" {
			if strings.HasPrefix(filename, "sources") {
				filenames := strings.SplitN(filename, "/", 3)
				if len(filenames) < 3 {
					return true
				}
				sourceTitle = filenames[2]
			} else {
				isTranslated = false
			}
		} else if status == "added" {
			if !strings.HasPrefix(filename, "translated") {
				isTranslated = false
			}
		}

		return true
	})

	if !isTranslated || sourceTitle == "" {
		return nil
	}

	return self.insertOrUpdateGCCT(_prInfo, sourceTitle, true)
}

func (self GithubLogic) translateSilmu(filesResult gjson.Result, _prInfo *prInfo) error {
	var (
		sourceTitle  string
		isTranslated = true
	)

	filesResult.ForEach(func(key, val gjson.Result) bool {
		if !isTranslated {
			return false
		}

		status := val.Get("status").String()
		filename := val.Get("filename").String()

		if status == "removed" {
			if strings.HasPrefix(filename, "sources") {
				filenames := strings.SplitN(filename, "/", 3)
				if len(filenames) < 3 {
					return true
				}
				sourceTitle = filenames[2]
			} else {
				isTranslated = false
			}
		} else if status == "added" {
			if !strings.HasPrefix(filename, "translated") {
				isTranslated = false
			}
		} else if status == "modified" {
			// 提交完成，之后又领取了新的一篇
			if strings.HasPrefix(filename, "sources") {
				filenames := strings.SplitN(filename, "/", 3)
				if len(filenames) < 3 {
					return true
				}
				title := filenames[2]
				if title == "" {
					return true
				}

				self.insertOrUpdateGCCT(_prInfo, title, false)
			}
		}

		return true
	})

	if !isTranslated || sourceTitle == "" {
		return nil
	}

	return self.insertOrUpdateGCCT(_prInfo, sourceTitle, true)
}

func (GithubLogic) insertOrUpdateGCCT(_prInfo *prInfo, title string, isTranslated bool) error {
	md5 := goutils.Md5(title)
	gcttGit := &model.GCTTGit{}
	_, err := MasterDB.Where("md5=?", md5).Get(gcttGit)
	if err != nil {
		logger.Errorln("GithubLogic insertOrUpdateGCCT get error:", err)
		return err
	}
	if gcttGit.Id > 0 {
		if gcttGit.Username != _prInfo.username {
			return nil
		}
	}

	gcttUser := DefaultGCTT.FindOne(nil, _prInfo.username)

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	if gcttUser.Id == 0 {
		gcttUser.Username = _prInfo.username
		gcttUser.Avatar = _prInfo.avatar
		gcttUser.JoinedAt = _prInfo.prTime.Unix()
		_, err = session.Insert(gcttUser)
		if err != nil {
			session.Rollback()
			logger.Errorln("GithubLogic insertOrUpdateGCCT insert gctt_user error:", err)
			return err
		}
	}

	// 已经存在
	if gcttGit.Id > 0 {
		if gcttGit.TranslatedAt == 0 && isTranslated {
			gcttGit.TranslatedAt = _prInfo.prTime.Unix()
			gcttGit.PR = _prInfo.number
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

	gcttGit.Username = _prInfo.username
	gcttGit.Title = title
	gcttGit.Md5 = md5
	gcttGit.TranslatingAt = _prInfo.prTime.Unix()
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
		var words int
		for _, gcttGit := range gcttGits {
			if gcttGit.TranslatingAt != 0 && gcttGit.TranslatedAt != 0 {
				avgTime += gcttGit.TranslatedAt - gcttGit.TranslatingAt
			}

			if gcttGit.TranslatedAt > lastAt {
				lastAt = gcttGit.TranslatedAt
			}

			if gcttGit.Words == 0 && gcttGit.ArticleId > 0 {
				article, _ := DefaultArticle.FindById(nil, gcttGit.ArticleId)
				gcttGit.Words = utf8.RuneCountInString(article.Content)
			}

			words += gcttGit.Words

			MasterDB.Id(gcttGit.Id).Update(gcttGit)
		}

		// 查询是否绑定了本站账号
		uid := DefaultThirdUser.findUid(gcttUser.Username, model.BindTypeGithub)

		gcttUser.Num = len(gcttGits)
		gcttUser.Words = words
		if gcttUser.Num > 0 {
			gcttUser.AvgTime = int(avgTime) / gcttUser.Num
		}
		gcttUser.LastAt = lastAt
		gcttUser.Uid = uid
		_, err = MasterDB.Id(gcttUser.Id).Update(gcttUser)
		if err != nil {
			logger.Errorln("GithubLogic update gctt user error:", err)
		}
	}
}
