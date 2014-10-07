// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	"time"

	"github.com/robfig/cron"
	"global"
	"logger"
	"service"
	"util"
)

// 后台运行的任务
func ServeBackGround() {

	// 初始化 七牛云存储
	service.InitQiniu()

	// 常驻内存的数据
	go loadData()

	c := cron.New()

	// 每天对非活跃用户降频
	c.AddFunc("@daily", decrUserActiveWeight)

	// 两分钟刷一次浏览数（TODO：重启丢失问题？信号控制重启？）
	c.AddFunc("@every 2m", service.Views.Flush)

	c.Start()
}

func loadData() {
	service.LoadAuthorities()
	service.LoadRoles()
	service.LoadRoleAuthorities()
	service.LoadNodes()
	service.LoadCategories()

	for {
		select {
		case <-global.AuthorityChan:
			service.LoadAuthorities()
		case <-global.RoleChan:
			service.LoadRoles()
		case <-global.RoleAuthChan:
			service.LoadRoleAuthorities()
		}
	}
}

func decrUserActiveWeight() {
	logger.Debugln("start decr user active weight...")

	loginTime := time.Now().Add(-72 * time.Hour)
	userList, err := service.FindNotLoginUsers(loginTime.Format(util.TIME_LAYOUT_OFTEN))
	if err != nil {
		logger.Errorln("获取最近未登录用户失败：", err)
		return
	}

	logger.Debugln("need dealing users:", len(userList))

	for _, user := range userList {
		divide := 5

		lastLoginTime, err := util.TimeParseOften(user.LoginTime)
		if err == nil {
			hours := (loginTime.Sub(lastLoginTime) / 24).Hours()
			if hours < 24 {
				divide = 2
			} else if hours < 48 {
				divide = 3
			} else if hours < 72 {
				divide = 4
			}
		}

		logger.Debugln("decr user weight, username:", user.Username, "divide:", divide)
		service.DecrUserWeight("username='"+user.Username+"'", divide)
	}

	logger.Debugln("end decr user active weight...")
}
