// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"logger"
	"util"
)

// 角色分界点：roleid小于该值，则没有管理权限
const AdminMinRoleId int = 6

// 角色信息
type Role struct {
	Roleid int    `json:"roleid" pk:"1"`
	Name   string `json:"name"`
	OpUser string `json:"op_user"`
	Ctime  string `json:"ctime,omitempty"`
	Mtime  string `json:"mtime,omitempty"`

	// 数据库访问对象
	*Dao
}

func NewRole() *Role {
	return &Role{
		Dao: &Dao{tablename: "role"},
	}
}

func (this *Role) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Role) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Role) FindAll(selectCol ...string) ([]*Role, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	roleList := make([]*Role, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		role := NewRole()
		err = this.Scan(rows, colNum, role.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Role FindAll Scan Error:", err)
			continue
		}
		roleList = append(roleList, role)
	}
	return roleList, nil
}

// 为了支持连写
func (this *Role) Where(condition string, args ...interface{}) *Role {
	this.Dao.Where(condition, args...)
	return this
}

// 为了支持连写
func (this *Role) Order(order string) *Role {
	this.Dao.Order(order)
	return this
}

// 为了支持连写
func (this *Role) Limit(limit string) *Role {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Role) Set(clause string, args ...interface{}) *Role {
	this.Dao.Set(clause, args...)
	return this
}

func (this *Role) prepareInsertData() {
	this.columns = []string{"name", "op_user", "ctime"}
	this.colValues = []interface{}{this.Name, this.OpUser, this.Ctime}
}

func (this *Role) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"roleid":  &this.Roleid,
		"name":    &this.Name,
		"op_user": &this.OpUser,
		"ctime":   &this.Ctime,
		"mtime":   &this.Mtime,
	}
}

// 角色权限信息
type RoleAuthority struct {
	Roleid int    `json:"roleid"`
	Aid    int    `json:"aid"`
	OpUser string `json:"op_user"`
	Ctime  string `json:"ctime"`

	// 内嵌
	*Dao
}

func NewRoleAuthority() *RoleAuthority {
	return &RoleAuthority{
		Dao: &Dao{tablename: "role_authority"},
	}
}

func (this *RoleAuthority) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *RoleAuthority) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *RoleAuthority) FindAll(selectCol ...string) ([]*RoleAuthority, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	roleAuthList := make([]*RoleAuthority, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		roleAuth := NewRoleAuthority()
		err = this.Scan(rows, colNum, roleAuth.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("RoleAuthority FindAll Scan Error:", err)
			continue
		}
		roleAuthList = append(roleAuthList, roleAuth)
	}
	return roleAuthList, nil
}

// 为了支持连写
func (this *RoleAuthority) Where(condition string, args ...interface{}) *RoleAuthority {
	this.Dao.Where(condition, args...)
	return this
}

func (this *RoleAuthority) prepareInsertData() {
	this.columns = []string{"roleid", "aid", "op_user"}
	this.colValues = []interface{}{this.Roleid, this.Aid, this.OpUser}
}

func (this *RoleAuthority) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"roleid":  &this.Roleid,
		"aid":     &this.Aid,
		"op_user": &this.OpUser,
		"ctime":   &this.Ctime,
	}
}
