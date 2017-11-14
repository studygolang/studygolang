// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

// 角色分界点：roleid 大于该值，则没有管理权限
const AdminMinRoleId = 7 // 晨读管理员

const (
	// Master 站长
	Master = iota + 1
	AssistantMaster
	Administrator
	TopicAdmin
	ResourceAdmin
	ArticleAdmin
	ReadingAdmin
)

// 角色信息
type Role struct {
	Roleid int    `json:"roleid" xorm:"pk autoincr"`
	Name   string `json:"name"`
	OpUser string `json:"op_user"`
	Ctime  string `json:"ctime,omitempty" xorm:"created"`
	Mtime  string `json:"mtime,omitempty" xorm:"<-"`
}

// 角色权限信息
type RoleAuthority struct {
	Roleid int    `json:"roleid" xorm:"pk autoincr"`
	Aid    int    `json:"aid"`
	OpUser string `json:"op_user"`
	Ctime  string `json:"ctime" xorm:"<-"`
}
