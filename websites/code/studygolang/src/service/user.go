// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"logger"
	"model"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"
)

func CreateUser(form url.Values) (errMsg string, err error) {
	if EmailExists(form.Get("email")) {
		err = errors.New("该邮箱已注册过")
		return
	}
	if UsernameExists(form.Get("username")) {
		err = errors.New("用户名已存在")
		return
	}
	// 存用户基本信息，产生自增长UID
	user := model.NewUser()
	err = util.ConvertAssign(user, form)
	if err != nil {
		logger.Errorln("user ConvertAssign error", err)
		errMsg = err.Error()
		return
	}
	uid, err := user.Insert()
	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln(errMsg, "：", err)
		return
	}

	// 存用户登录信息
	userLogin := model.NewUserLogin()
	err = util.ConvertAssign(userLogin, form)
	if err != nil {
		errMsg = err.Error()
		logger.Errorln("CreateUser error:", err)
		return
	}
	userLogin.Uid = uid
	_, err = userLogin.Insert()
	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln(errMsg, "：", err)
		return
	}

	// 存用户角色信息
	userRole := model.NewUserRole()
	// 默认为初级会员
	roleId := model.AllRoleId[len(model.AllRoleId)-1]
	if form.Get("roleid") != "" {
		tmpId, err := strconv.Atoi(form.Get("roleid"))
		if err == nil {
			roleId = tmpId
		}
	}
	userRole.Roleid = roleId
	userRole.Uid = uid
	if _, err = userRole.Insert(); err != nil {
		logger.Errorln("userRole insert Error:", err)
	}

	// 存用户活跃信息，初始活跃+2
	userActive := model.NewUserActive()
	userActive.Uid = uid
	userActive.Username = user.Username
	userActive.Email = user.Email
	userActive.Weight = 2
	if _, err = userActive.Insert(); err != nil {
		logger.Errorln("UserActive insert Error:", err)
	}
	return
}

// 修改用户资料
func UpdateUser(form url.Values) (errMsg string, err error) {
	fields := []string{"name", "open", "city", "company", "github", "weibo", "website", "status", "introduce"}
	setClause := GenSetClause(form, fields)
	username := form.Get("username")
	err = model.NewUser().Set(setClause).Where("username=" + username).Update()
	if err != nil {
		logger.Errorf("更新用户 【%s】 信息失败：%s", username, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	// 修改用户资料，活跃度+1
	go IncUserWeight("username="+username, 1)

	return
}

// 获取当前登录用户信息（常用信息）
func FindCurrentUser(username string) (user map[string]interface{}, err error) {
	userInfo := model.NewUser()
	err = userInfo.Where("username=" + username).Find()
	if err != nil {
		logger.Errorf("获取用户 %s 信息失败：%s", username, err)
		return
	}
	if userInfo.Uid == 0 {
		logger.Infof("用户 %s 不存在！", username)
		return
	}
	user = map[string]interface{}{
		"uid":      userInfo.Uid,
		"username": userInfo.Username,
		"email":    userInfo.Email,
	}

	// 获取未读消息数
	user["msgnum"] = FindNotReadMsgNum(userInfo.Uid)

	// 获取角色信息
	userRoleList, err := model.NewUserRole().Where("uid=" + strconv.Itoa(userInfo.Uid)).FindAll()
	if err != nil {
		logger.Errorf("获取用户 %s 角色 信息失败：%s", username, err)
		return
	}
	for _, userRole := range userRoleList {
		if userRole.Roleid <= model.AdminMinRoleId {
			// 是管理员
			user["isadmin"] = true
		}
	}

	RecordLoginTime(username)

	return
}

// 判断指定的用户名是否存在
func UsernameExists(username string) bool {
	userLogin := model.NewUserLogin()
	if err := userLogin.Where("username=" + username).Find("uid"); err != nil {
		logger.Errorln("service UsernameExists error:", err)
		return false
	}
	if userLogin.Uid != 0 {
		return true
	}
	return false
}

// 判断指定的邮箱（email）是否存在
func EmailExists(email string) bool {
	userLogin := model.NewUserLogin()
	if err := userLogin.Where("email=" + email).Find("uid"); err != nil {
		logger.Errorln("service EmailExists error:", err)
		return false
	}
	if userLogin.Uid != 0 {
		return true
	}
	return false
}

// 获取单个用户信息
func FindUserByUsername(username string) *model.User {
	user := model.NewUser()
	err := user.Where("username=" + username).Find()
	if err != nil {
		logger.Errorf("获取用户 %s 信息失败：%s", username, err)
		return nil
	}
	if user.Uid == 0 {
		return nil
	}

	// 获取角色信息
	userRoleList, err := model.NewUserRole().Where("uid=" + strconv.Itoa(user.Uid)).FindAll()
	if err != nil {
		logger.Errorf("获取用户 %s 角色 信息失败：%s", username, err)
		return nil
	}
	for _, userRole := range userRoleList {
		if len(user.Roleids) == 0 {
			user.Rolenames = []string{model.AllRole[userRole.Roleid].Name}
		} else {
			user.Rolenames = append(user.Rolenames, model.AllRole[userRole.Roleid].Name)
		}
	}
	return user
}

// 获得活跃用户
func FindActiveUsers(start, num int) []*model.UserActive {
	activeUsers, err := model.NewUserActive().Order("weight DESC").Limit(strconv.Itoa(start) + "," + strconv.Itoa(num)).FindAll()
	if err != nil {
		logger.Errorln("user service FindActiveUsers error:", err)
		return nil
	}
	return activeUsers
}

func FindNewUsers(start, num int) []*model.User {
	users, err := model.NewUser().Order("ctime DESC").Limit(strconv.Itoa(start) + "," + strconv.Itoa(num)).FindAll([]string{"uid", "username", "email", "avatar", "ctime"}...)
	if err != nil {
		logger.Errorln("user service FindNewUsers error:", err)
		return nil
	}
	return users
}

func FindUsersByPage(conds map[string]string, curPage, limit int) ([]*model.User, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	user := model.NewUser()

	limitStr := strconv.Itoa(curPage*limit) + "," + strconv.Itoa(limit)
	userList, err := user.Where(strings.Join(conditions, " AND ")).Limit(limitStr).
		FindAll()
	if err != nil {
		logger.Errorln("user service FindUsersByPage Error:", err)
		return nil, 0
	}

	total, err := user.Count()
	if err != nil {
		logger.Errorln("user service FindUsersByPage COUNT Error:", err)
		return nil, 0
	}

	return userList, total
}

var (
	ErrUsername = errors.New("用户名不存在")
	ErrPasswd   = errors.New("密码错误")
)

// 登录；成功返回用户登录信息(user_login)
func Login(username, passwd string) (*model.UserLogin, error) {
	userLogin := model.NewUserLogin()
	err := userLogin.Where("username=" + username + " OR email=" + username).Find()
	if err != nil {
		logger.Errorf("用户 %s 登录错误：%s", username, err)
		return nil, errors.New("内部错误，请稍后再试！")
	}
	// 校验用户
	if userLogin.Uid == 0 {
		logger.Infof("用户名 %s 不存在", username)
		return nil, ErrUsername
	}
	passcode := userLogin.GetPasscode()
	md5Passwd := util.Md5(passwd + passcode)
	logger.Debugf("passwd: %s, passcode: %s, md5passwd: %s, dbpasswd: %s", passwd, passcode, md5Passwd, userLogin.Passwd)
	if md5Passwd != userLogin.Passwd {
		logger.Infof("用户名 %s 填写的密码错误", username)
		return nil, ErrPasswd
	}

	// 登录，活跃度+1
	go IncUserWeight("uid="+strconv.Itoa(userLogin.Uid), 1)

	RecordLoginTime(username)

	return userLogin, nil
}

// 记录用户最后登录时间
func RecordLoginTime(username string) error {
	userLogin := model.NewUserLogin()
	err := userLogin.Set("login_time=" + time.Now().Format("2006-01-02 15:04:05")).Where("username=" + username).Update()
	if err != nil {
		logger.Errorf("记录用户 %s 登录时间错误：%s", username, err)
	}
	return err
}

// 更新用户密码（用户名或email）
func UpdatePasswd(username, passwd string) (string, error) {
	userLogin := model.NewUserLogin()
	passwd = userLogin.GenMd5Passwd(passwd)
	err := userLogin.Set("passwd=" + passwd + ",passcode=" + userLogin.GetPasscode()).Where("username=" + username + " OR email=" + username).Update()
	if err != nil {
		logger.Errorf("用户 %s 更新密码错误：%s", username, err)
		return "对不起，内部服务错误！", err
	}
	return "", nil
}

// 获取用户信息
func getUserInfos(uids map[int]int) map[int]*model.User {
	if len(uids) == 0 {
		return nil
	}
	// 获取用户信息
	inUids := util.Join(util.MapIntKeys(uids), ",")
	users, err := model.NewUser().Where("uid in(" + inUids + ")").FindAll()
	if err != nil {
		logger.Errorln("user service getUserInfos Error:", err)
		return map[int]*model.User{}
	}
	userMap := make(map[int]*model.User, len(users))
	for _, user := range users {
		userMap[user.Uid] = user
	}
	return userMap
}

// 会员总数
func CountUsers() int {
	total, err := model.NewUserLogin().Count()
	if err != nil {
		logger.Errorln("user service CountUsers error:", err)
		return 0
	}
	return total
}

// 构造update语句中的set部分子句
func GenSetClause(form url.Values, fields []string) string {
	stringBuilder := util.NewBuffer()
	for _, field := range fields {
		if form.Get(field) != "" {
			stringBuilder.Append(",").Append(field).Append("=").Append(form.Get(field))
		}
	}
	if stringBuilder.Len() > 0 {
		return stringBuilder.String()[1:]
	}
	return ""
}

// 增加或减少用户活跃度
func IncUserWeight(where string, weight int) {
	if err := model.NewUserActive().Where(where).Increment("weight", weight); err != nil {
		logger.Errorln("UserActive update Error:", err)
	}
}

func DecrUserWeight(where string, divide int) {
	if divide <= 0 {
		return
	}

	strSql := "UPDATE user_active SET weight=weight/" + strconv.Itoa(divide) + " WHERE " + where
	if result, err := model.NewUserActive().Exec(strSql); err != nil {
		logger.Errorln("UserActive update Error:", err)
	} else {
		n, _ := result.RowsAffected()
		logger.Debugln(strSql, "affected num:", n)
	}
}

// 获取 loginTime 之前没有登录的用户
func FindNotLoginUsers(loginTime string) (userList []*model.UserLogin, err error) {
	userLogin := model.NewUserLogin()
	userList, err = userLogin.Where("login_time<" + loginTime).FindAll()
	return
}
