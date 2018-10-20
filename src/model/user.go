// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/polaris1119/goutils"
)

// 用户登录信息
type UserLogin struct {
	Uid       int       `json:"uid" xorm:"pk"`
	Username  string    `json:"username"`
	Passcode  string    `json:"passcode"` // 加密随机串
	Passwd    string    `json:"passwd"`
	Email     string    `json:"email"`
	LoginIp   string    `json:"login_ip"`
	LoginTime time.Time `json:"login_time" xorm:"<-"`
}

func (this *UserLogin) TableName() string {
	return "user_login"
}

// 生成加密密码
func (this *UserLogin) GenMd5Passwd() error {
	if this.Passwd == "" {
		return errors.New("password is empty!")
	}
	this.Passcode = fmt.Sprintf("%x", rand.Int31())
	// 密码经过md5(passwd+passcode)加密保存
	this.Passwd = goutils.Md5(this.Passwd + this.Passcode)
	return nil
}

const (
	UserStatusNoAudit = iota
	UserStatusAudit   // 已激活
	UserStatusRefuse
	UserStatusFreeze // 冻结
	UserStatusOutage // 停用
)

const (
	// 用户拥有的权限设置
	DauAuthTopic = 1 << iota
	DauAuthArticle
	DauAuthResource
	DauAuthWiki
	DauAuthProject
	DauAuthBook
	DauAuthComment // 评论
	DauAuthTop     // 置顶
)

// 置顶
const (
	TypeComment = 100
	TypeTop     = 101
)

const DefaultAuth = DauAuthTopic | DauAuthArticle | DauAuthResource | DauAuthProject | DauAuthComment

// 用户基本信息
type User struct {
	Uid         int       `json:"uid" xorm:"pk autoincr"`
	Username    string    `json:"username" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
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
	Balance     int       `json:"balance"`
	IsThird     int       `json:"is_third"`
	DauAuth     int       `json:"dau_auth"`
	Status      int       `json:"status"`
	IsRoot      bool      `json:"is_root"`
	Ctime       OftenTime `json:"ctime" xorm:"created"`
	Mtime       time.Time `json:"mtime" xorm:"<-"`

	// 非用户表中的信息，为了方便放在这里
	Roleids   []int    `xorm:"-"`
	Rolenames []string `xorm:"-"`

	// 活跃度
	Weight int `json:"weight" xorm:"-"`
	Gold   int `json:"gold" xorm:"-"`
	Silver int `json:"silver" xorm:"-"`
	Copper int `json:"copper" xorm:"-"`

	IsOnline bool `json:"is_online" xorm:"-"`
}

func (this *User) TableName() string {
	return "user_info"
}

func (this *User) String() string {
	buffer := goutils.NewBuffer()
	buffer.Append(this.Username).Append(" ").
		Append(this.Email).Append(" ").
		Append(this.Uid).Append(" ").
		Append(this.Mtime)

	return buffer.String()
}

func (this *User) AfterSet(name string, cell xorm.Cell) {
	if name == "balance" {
		this.Gold = this.Balance / 10000
		balance := this.Balance % 10000

		this.Silver = balance / 100
		this.Copper = balance % 100
	}
}

// Me 代表当前用户
type Me struct {
	Uid       int       `json:"uid"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Monlog    string    `json:"monlog"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Status    int       `json:"status"`
	MsgNum    int       `json:"msgnum"`
	IsAdmin   bool      `json:"isadmin"`
	IsRoot    bool      `json:"is_root"`
	DauAuth   int       `json:"dau_auth"`
	CreatedAt time.Time `json:"created_at"`

	Balance int `json:"balance"`
	Gold    int `json:"gold"`
	Silver  int `json:"silver"`
	Copper  int `json:"copper"`

	RoleIds []int `json:"-"`
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
	Uid      int       `json:"uid" xorm:"pk"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Weight   int       `json:"weight"`
	Mtime    time.Time `json:"mtime" xorm:"<-"`
}

// 用户角色信息
type UserRole struct {
	Uid    int    `json:"uid"`
	Roleid int    `json:"roleid"`
	ctime  string `xorm:"-"`
}

const (
	BindTypeGithub = iota
)

type BindUser struct {
	Id           int       `json:"id" xorm:"pk autoincr"`
	Uid          int       `json:"uid"`
	Type         int       `json:"type"`
	Email        string    `json:"email"`
	Tuid         int       `json:"tuid"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expire       int       `json:"expire"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"created_at" xorm:"<-"`
}
