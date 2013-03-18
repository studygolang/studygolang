// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"fmt"
	"logger"
	"math/rand"
	"util"
)

// 用户登录信息
type UserLogin struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	Email    string `json:"email"`
	passcode string // 加密随机串

	// 数据库访问对象
	*Dao
}

func NewUserLogin() *UserLogin {
	return &UserLogin{
		Dao: &Dao{tablename: "user_login"},
	}
}

func (this *UserLogin) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *UserLogin) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

// 为了支持连写
func (this *UserLogin) Where(condition string) *UserLogin {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *UserLogin) Set(clause string) *UserLogin {
	this.Dao.Set(clause)
	return this
}

func (this *UserLogin) prepareInsertData() {
	this.columns = []string{"uid", "username", "passwd", "email", "passcode"}
	this.GenMd5Passwd("")
	this.colValues = []interface{}{this.Uid, this.Username, this.Passwd, this.Email, this.passcode}
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

// 由于在DAO中没法调用 具体 model 的方法，如果将该映射关系定义为 具体 model 字段，有些浪费
func (this *UserLogin) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"uid":      &this.Uid,
		"username": &this.Username,
		"passwd":   &this.Passwd,
		"email":    &this.Email,
		"passcode": &this.passcode,
	}
}

func (this *UserLogin) GetPasscode() string {
	return this.passcode
}

// 用户基本信息
type User struct {
	Uid       int    `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	City      string `json:"city"`
	Company   string `json:"company"`
	Github    string `json:"github"`
	Weibo     string `json:"weibo"`
	Website   string `json:"website"`
	Status    string `json:"status"`
	Introduce string `json:"introduce"`
	Ctime     string `json:"ctime"`
	Open      int    `json:"open"`

	// 非用户表中的信息，为了方便放在这里
	Roleids   []int
	Rolenames []string

	// 内嵌
	*Dao
}

func NewUser() *User {
	return &User{
		Dao: &Dao{tablename: "user_info"},
	}
}

func (this *User) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (this *User) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *User) FindAll(selectCol ...string) ([]*User, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	userList := make([]*User, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		user := NewUser()
		err = this.Scan(rows, colNum, user.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("User FindAll Scan Error:", err)
			continue
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func (this *User) Where(condition string) *User {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *User) Set(clause string) *User {
	this.Dao.Set(clause)
	return this
}

// 为了支持连写
func (this *User) Limit(limit string) *User {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *User) Order(order string) *User {
	this.Dao.Order(order)
	return this
}

func (this *User) prepareInsertData() {
	this.columns = []string{"username", "email", "name", "avatar", "city", "company", "github", "weibo", "website", "status", "introduce"}
	this.colValues = []interface{}{this.Username, this.Email, this.Name, this.Avatar, this.City, this.Company, this.Github, this.Weibo, this.Website, this.Status, this.Introduce}
}

func (this *User) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"uid":       &this.Uid,
		"username":  &this.Username,
		"email":     &this.Email,
		"name":      &this.Name,
		"avatar":    &this.Avatar,
		"city":      &this.City,
		"company":   &this.Company,
		"github":    &this.Github,
		"weibo":     &this.Weibo,
		"website":   &this.Website,
		"status":    &this.Status,
		"introduce": &this.Introduce,
		"open":      &this.Open,
		"ctime":     &this.Ctime,
	}
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
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Weight   int    `json:"weight"`
	Mtime    string `json:"mtime"`

	// 内嵌
	*Dao
}

func NewUserActive() *UserActive {
	return &UserActive{
		Dao: &Dao{tablename: "user_active"},
	}
}

func (this *UserActive) Insert() (int, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (this *UserActive) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *UserActive) FindAll(selectCol ...string) ([]*UserActive, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	userList := make([]*UserActive, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		user := NewUserActive()
		err = this.Scan(rows, colNum, user.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("User FindAll Scan Error:", err)
			continue
		}
		userList = append(userList, user)
	}
	return userList, nil
}

// 设置更新字段
func (this *UserActive) Set(clause string) *UserActive {
	this.Dao.Set(clause)
	return this
}

// 为了支持连写
func (this *UserActive) Where(condition string) *UserActive {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *UserActive) Limit(limit string) *UserActive {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *UserActive) Order(order string) *UserActive {
	this.Dao.Order(order)
	return this
}

func (this *UserActive) prepareInsertData() {
	this.columns = []string{"uid", "username", "email", "avatar", "weight"}
	this.colValues = []interface{}{this.Uid, this.Username, this.Email, this.Avatar, this.Weight}
}

func (this *UserActive) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"uid":      &this.Uid,
		"username": &this.Username,
		"email":    &this.Email,
		"avatar":   &this.Avatar,
		"weight":   &this.Weight,
		"mtime":    &this.Mtime,
	}
}

// 用户角色信息
type UserRole struct {
	Uid    int `json:"uid"`
	Roleid int `json:"roleid"`
	ctime  string

	//
	*Dao
}

func NewUserRole() *UserRole {
	return &UserRole{
		Dao: &Dao{tablename: "user_role"},
	}
}

func (this *UserRole) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *UserRole) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *UserRole) FindAll(selectCol ...string) ([]*UserRole, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		logger.Errorln("[UserRole.FindAll] error:", err)
		return nil, err
	}
	// TODO:
	userRoleList := make([]*UserRole, 0, 10)
	colNum := len(selectCol)
	for rows.Next() {
		userRole := NewUserRole()
		err = this.Scan(rows, colNum, userRole.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("UserRole FindAll Scan Error:", err)
			continue
		}
		userRoleList = append(userRoleList, userRole)
	}
	return userRoleList, nil
}

func (this *UserRole) Where(condition string) *UserRole {
	this.Dao.Where(condition)
	return this
}

func (this *UserRole) prepareInsertData() {
	this.columns = []string{"uid", "roleid"}
	this.colValues = []interface{}{this.Uid, this.Roleid}
}

func (this *UserRole) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"uid":    &this.Uid,
		"roleid": &this.Roleid,
		"ctime":  &this.ctime,
	}
}
