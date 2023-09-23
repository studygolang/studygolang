// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"flag"
	"time"

	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
	"github.com/robfig/cron/v3"

	"github.com/studygolang/studygolang/cmd"
	"github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/global"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"
)

var (
	embedIndexing = flag.Bool("embed_indexing", false, "是否嵌入 indexer 的功能，默认否")
	embedCrawler  = flag.Bool("embed_crawler", false, "是否嵌入 crawler 的功能，默认否")
	syncAllGCTT   = flag.Bool("sync_gctt", false, "是否全量同步 GCTT PR 一次")
)

// 后台运行的任务
func ServeBackGround() {

	if db.MasterDB == nil {
		return
	}

	// 初始化 七牛云存储
	logic.DefaultUploader.InitQiniu()

	if *embedIndexing {
		cmd.IndexingServer()
	}
	if *embedCrawler {
		cmd.CrawlServer()
	}

	// 常驻内存的数据
	go loadData()

	c := cron.New()

	if config.ConfigFile.MustBool("global", "is_master", false) {
		// 每天对非活跃用户降频
		c.AddFunc("@daily", decrUserActiveWeight)

		// 生成阅读排行榜
		c.AddFunc("@daily", genViewRank)

		if global.OnlineEnv() {
			// 每天生成 sitemap 文件
			c.AddFunc("@daily", logic.GenSitemap)

			// 给用户发邮件，如通知网站最近的动态，每周的晨读汇总等
			c.AddFunc("0 0 0 * * *", logic.DefaultEmail.EmailNotice)

			// webhook 方式增量，每天补漏
			c.AddFunc("@daily", syncGCTTRepo)
		}

		// 取消置顶
		c.AddFunc("0 * * * * *", unsetTop)

		// 每天对活跃用户奖励铜币
		c.AddFunc("@daily", logic.DefaultUserRich.AwardCooper)

		// 首页推荐自动调整
		c.AddFunc("@every 5m", logic.DefaultFeed.AutoUpdateSeq)

		// 每日题目
		c.AddFunc("@daily", logic.DefaultInterview.UpdateTodayQuestionID)
	}

	// 两分钟刷一次浏览数（TODO：重启丢失问题？信号控制重启？）
	c.AddFunc("@every 2m", logic.Views.Flush)

	c.Start()
}

func loadData() {
	logic.LoadAuthorities()
	logic.LoadRoles()
	logic.LoadRoleAuthorities()
	logic.LoadNodes()
	logic.LoadCategories()
	logic.LoadWebsiteSetting()
	logic.LoadDefaultAvatar()
	logic.LoadUserSetting()

	for {
		select {
		case <-global.AuthorityChan:
			logic.LoadAuthorities()
		case <-global.RoleChan:
			logic.LoadRoles()
		case <-global.RoleAuthChan:
			logic.LoadRoleAuthorities()
		case <-global.UserSettingChan:
			logic.LoadUserSetting()
		case <-global.TopicNodeChan:
			logic.LoadNodes()
		}
	}
}

func decrUserActiveWeight() {
	logger.Debugln("start decr user active weight...")

	loginTime := time.Now().Add(-72 * time.Hour)
	userList, err := logic.DefaultUser.FindNotLoginUsers(loginTime)
	if err != nil {
		logger.Errorln("获取最近未登录用户失败：", err)
		return
	}

	logger.Debugln("need dealing users:", len(userList))

	for _, user := range userList {
		divide := 5

		if err == nil {
			hours := (loginTime.Sub(user.LoginTime) / 24).Hours()
			if hours < 24 {
				divide = 2
			} else if hours < 48 {
				divide = 3
			} else if hours < 72 {
				divide = 4
			}
		}

		logger.Debugln("decr user weight, username:", user.Username, "divide:", divide)
		logic.DefaultUser.DecrUserWeight("username", user.Username, divide)
	}

	logger.Debugln("end decr user active weight...")
}

func genViewRank() {
	needRankTypes := []int{
		model.TypeTopic,
		model.TypeResource,
		model.TypeProject,
		model.TypeArticle,
		model.TypeBook,
	}

	for _, objtype := range needRankTypes {
		logic.DefaultRank.GenWeekRank(objtype)
		logic.DefaultRank.GenMonthRank(objtype)
	}
}

func unsetTop() {
	logic.DefaultTopic.AutoUnsetTop()
}

func syncGCTTRepo() {
	repo := config.ConfigFile.MustValue("gctt", "repo")
	if repo == "" {
		return
	}

	logic.DefaultGithub.PullPR(repo, *syncAllGCTT)
	logic.DefaultGithub.SyncIssues(repo, *syncAllGCTT)
	*syncAllGCTT = false
}
