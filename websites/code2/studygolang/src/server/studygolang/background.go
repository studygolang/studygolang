// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"db"
	"global"
	"logic"
	"time"

	"github.com/polaris1119/logger"
	"github.com/robfig/cron"
)

// 后台运行的任务
func ServeBackGround() {

	if db.MasterDB == nil {
		return
	}

	// 初始化 七牛云存储
	logic.DefaultUploader.InitQiniu()

	// 常驻内存的数据
	go loadData()

	c := cron.New()

	// 每天对非活跃用户降频
	c.AddFunc("@daily", decrUserActiveWeight)

	// 两分钟刷一次浏览数（TODO：重启丢失问题？信号控制重启？）
	c.AddFunc("@every 2m", logic.Views.Flush)

	if global.OnlineEnv() {
		// 每天生成 sitemap 文件
		c.AddFunc("@daily", logic.GenSitemap)

		// 给用户发邮件，如通知网站最近的动态，每周的晨读汇总等
		c.AddFunc("0 0 4 * * 1", logic.DefaultEmail.EmailNotice)
	}

	c.Start()
}

func loadData() {
	logic.LoadAuthorities()
	logic.LoadRoles()
	logic.LoadRoleAuthorities()
	logic.LoadNodes()
	logic.LoadCategories()

	for {
		select {
		case <-global.AuthorityChan:
			logic.LoadAuthorities()
		case <-global.RoleChan:
			logic.LoadRoles()
		case <-global.RoleAuthChan:
			logic.LoadRoleAuthorities()
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
