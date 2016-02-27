// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package model

import (
	"fmt"
	"math/rand"
	"time"

	"util"
)

// 用户登录信息
type UserLogin struct {
	Uid       int    `json:"uid" gorm:"primary_key"`
	Username  string `json:"username"`
	Passwd    string `json:"passwd"`
	Email     string `json:"email"`
	LoginTime string `json:"login_time"`
	passcode  string // 加密随机串
}

func (this *UserLogin) TableName() string {
	return "user_login"
}

// 生成加密密码
func (this *UserLogin) GenMd5Passwd(origPwd string) string {
	if origPwd == "" {
		origPwd = this.Passwd
	}
	this.passcode = fmt.Sprintf("%x", rand.Int31())
	// 密码经过md5(passwd+passcode)加密保存
	this.Passwd = util.Md5(origPwd + this.passcode)
	return this.Passwd
}

func (this *UserLogin) GetPasscode() string {
	return this.passcode
}

const (
	StatusNoAudit = iota
	StatusAudit
	StatusRefuse
	StatusFreeze // 冻结
	StatusOutage // 停用
)

// 用户基本信息
type User struct {
	Uid         int       `json:"uid" gorm:"primary_key"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Open        int       `json:"open"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
	City        string    `json:"city"`
	Company     string    `json:"company"`
	Github      string    `json:"github"`
	Weibo       string    `json:"weibo"`
	Website     string    `json:"website"`
	Monlog      string    `json:"monlog"`
	Introduce   string    `json:"introduce"`
	Unsubscribe int       `json:"unsubscribe"`
	Status      int       `json:"status"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`

	// 非用户表中的信息，为了方便放在这里
	//Roleids   []int
	//Rolenames []string
}

func (this *User) TableName() string {
	return "user_info"
}

// 活跃用户信息
// 活跃度规则：
//	1、注册成功后 +2
//	2、登录一次 +1
//	3、修改资料 +1
//	4、发帖子 + 10
//	5、评论 +5
//	6、创建Wiki页 +10
type UserActive struct {
	Uid      int       `json:"uid" gorm:"primary_key"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Weight   int       `json:"weight"`
	Mtime    time.Time `json:"mtime"`
}

func (this *UserActive) TableName() string {
	return "user_active"
}

// 用户角色信息
type UserRole struct {
	Uid    int `json:"uid"`
	Roleid int `json:"roleid"`
	ctime  string
}

func (this *UserRole) TableName() string {
	return "user_role"
}
