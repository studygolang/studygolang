// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"fmt"
	"math/rand"
	"model"
	"net/url"
	"strings"
	"time"
	"util"

	"github.com/polaris1119/slices"

	"github.com/go-validator/validator"
	"github.com/go-xorm/xorm"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"golang.org/x/net/context"

	. "db"
)

type UserLogic struct{}

var DefaultUser = UserLogic{}

// CreateUser 创建用户
func (self UserLogic) CreateUser(ctx context.Context, form url.Values) (errMsg string, err error) {
	objLog := GetLogger(ctx)

	if self.UserExists(ctx, "email", form.Get("email")) {
		errMsg = "该邮箱已注册过"
		err = errors.New(errMsg)
		return
	}
	if self.UserExists(ctx, "username", form.Get("username")) {
		errMsg = "用户名已存在"
		err = errors.New(errMsg)
		return
	}

	user := &model.User{}
	err = schemaDecoder.Decode(user, form)
	if err != nil {
		objLog.Errorln("user schema Decode error:", err)
		errMsg = err.Error()
		return
	}

	if err = validator.Validate(user); err != nil {
		objLog.Errorf("validate user error:%#v", err)

		// TODO: 暂时简单处理
		if errMap, ok := err.(validator.ErrorMap); ok {
			if _, ok = errMap["Username"]; ok {
				errMsg = "用户名不合法！"
			}
		} else {
			errMsg = err.Error()
		}
		return
	}

	if config.ConfigFile.MustBool("account", "verify_email", true) {
		if !user.IsRoot {
			// 避免前端伪造，传递 status=1
			user.Status = model.UserStatusNoAudit
		}
	} else {
		user.Status = model.UserStatusAudit
	}

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	err = self.doCreateUser(ctx, session, user, form.Get("passwd"))
	if err != nil {
		errMsg = "内部服务错误！"
		session.Rollback()
		objLog.Errorln("create user error:", err)
		return
	}

	if form.Get("id") != "" {
		id := goutils.MustInt(form.Get("id"))
		_, err = DefaultWechat.Bind(ctx, id, user.Uid, form.Get("userInfo"))
		if err != nil {
			session.Rollback()
			objLog.Errorln("bind wechat user error:", err)
			errMsg = err.Error()
			return
		}
	}

	session.Commit()

	return
}

// Update 更新用户信息
func (self UserLogic) Update(ctx context.Context, me *model.Me, form url.Values) (errMsg string, err error) {
	objLog := GetLogger(ctx)

	if form.Get("open") != "1" {
		form.Set("open", "0")
	}

	user := &model.User{}
	err = schemaDecoder.Decode(user, form)
	if err != nil {
		objLog.Errorln("userlogic update, schema decode error:", err)
		errMsg = "服务内部错误"
		return
	}

	cols := "name,open,city,company,github,weibo,website,monlog,introduce"
	// 变更了邮箱
	if user.Email != me.Email {
		cols += ",email,status"
		user.Status = model.UserStatusNoAudit
	}

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err = session.Id(me.Uid).Cols(cols).Update(user)
	if err != nil {
		session.Rollback()

		objLog.Errorf("更新用户 【%d】 信息失败：%s", me.Uid, err)
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			// TODO：被恶意注册？
			errMsg = "该邮箱地址被其他账号注册了"
		} else {
			errMsg = "对不起，服务器内部错误，请稍后再试！"
		}
		return
	}

	_, err = session.Table(new(model.UserLogin)).
		Where("uid=?", me.Uid).Update(map[string]interface{}{"email": me.Email})
	if err != nil {
		session.Rollback()
		objLog.Errorf("更新用户 【%d】 信息失败：%s", me.Uid, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	session.Commit()

	// 修改用户资料，活跃度+1
	go self.IncrUserWeight("uid", me.Uid, 1)

	return
}

// UpdateUserStatus 更新用户状态
func (UserLogic) UpdateUserStatus(ctx context.Context, uid, status int) error {
	objLog := GetLogger(ctx)

	_, err := MasterDB.Table(new(model.User)).Id(uid).Update(map[string]interface{}{"status": status})
	if err != nil {
		objLog.Errorf("更新用户 【%d】 状态失败：%s", uid, err)
	}

	return err
}

// ChangeAvatar 更换头像
func (UserLogic) ChangeAvatar(ctx context.Context, uid int, avatar string) (err error) {
	changeData := map[string]interface{}{"avatar": avatar}
	_, err = MasterDB.Table(new(model.User)).Id(uid).Update(changeData)
	if err == nil {
		_, err = MasterDB.Table(new(model.UserActive)).Id(uid).Update(changeData)
	}

	return
}

// UserExists 判断用户是否存在
func (UserLogic) UserExists(ctx context.Context, field, val string) bool {
	objLog := GetLogger(ctx)

	userLogin := &model.UserLogin{}
	_, err := MasterDB.Where(field+"=?", val).Get(userLogin)
	if err != nil || userLogin.Uid == 0 {
		if err != nil {
			objLog.Errorln("user logic UserExists error:", err)
		}
		return false
	}
	return true
}

// EmailOrUsernameExists 判断指定的邮箱（email）或用户名是否存在
func (UserLogic) EmailOrUsernameExists(ctx context.Context, email, username string) bool {
	objLog := GetLogger(ctx)

	userLogin := &model.UserLogin{}
	_, err := MasterDB.Where("email=?", email).Or("username=?", username).Get(userLogin)
	if err != nil || userLogin.Uid == 0 {
		if err != nil {
			objLog.Errorln("user logic EmailOrUsernameExists error:", err)
		}
		return false
	}
	return true
}

// FindUserInfos 获得用户信息，uniq 可能是 uid slice 或 username slice
func (self UserLogic) FindUserInfos(ctx context.Context, uniq interface{}) map[int]*model.User {
	objLog := GetLogger(ctx)

	field := "uid"
	if uids, ok := uniq.([]int); ok {
		if len(uids) == 0 {
			return nil
		}
	} else if usernames, ok := uniq.([]string); ok {
		if len(usernames) == 0 {
			return nil
		}
		field = "username"
	}

	usersMap := make(map[int]*model.User)
	if err := MasterDB.In(field, uniq).Find(&usersMap); err != nil {
		objLog.Errorln("user logic FindUserInfos not record found:", err)
		return nil
	}
	return usersMap
}

func (self UserLogic) FindOne(ctx context.Context, field string, val interface{}) *model.User {
	objLog := GetLogger(ctx)

	user := &model.User{}
	_, err := MasterDB.Where(field+"=?", val).Get(user)
	if err != nil {
		objLog.Errorln("user logic FindOne error:", err)
	}

	if user.Uid != 0 {
		if user.IsRoot {
			user.Roleids = []int{0}
			user.Rolenames = []string{"站长"}
			return user
		}

		// 获取用户角色信息
		userRoleList := make([]*model.UserRole, 0)
		err = MasterDB.Where("uid=?", user.Uid).OrderBy("roleid ASC").Find(&userRoleList)
		if err != nil {
			objLog.Errorf("获取用户 %s 角色 信息失败：%s", val, err)
			return nil
		}

		if roleNum := len(userRoleList); roleNum > 0 {
			user.Roleids = make([]int, roleNum)
			user.Rolenames = make([]string, roleNum)

			for i, userRole := range userRoleList {
				user.Roleids[i] = userRole.Roleid
				user.Rolenames[i] = Roles[userRole.Roleid-1].Name
			}
		}
	}
	return user
}

// 获取当前登录用户信息（常用信息）
func (self UserLogic) FindCurrentUser(ctx context.Context, username interface{}) *model.Me {
	objLog := GetLogger(ctx)

	user := &model.User{}
	var err error

	if uid, ok := username.(int); ok {
		_, err = MasterDB.Where("uid=? AND status<=?", uid, model.UserStatusAudit).Get(user)
	} else {
		_, err = MasterDB.Where("username=? AND status<=?", username, model.UserStatusAudit).Get(user)
	}

	if err != nil {
		objLog.Errorf("获取用户 %q 信息失败：%s", username, err)
		return &model.Me{}
	}
	if user.Uid == 0 {
		logger.Infof("用户 %q 不存在或状态不正常！", username)
		return &model.Me{}
	}

	me := &model.Me{
		Uid:       user.Uid,
		Username:  user.Username,
		Name:      user.Name,
		Monlog:    user.Monlog,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Status:    user.Status,
		IsRoot:    user.IsRoot,
		MsgNum:    DefaultMessage.FindNotReadMsgNum(ctx, user.Uid),
		DauAuth:   user.DauAuth,
		CreatedAt: time.Time(user.Ctime),

		Balance: user.Balance,
		Gold:    user.Gold,
		Silver:  user.Silver,
		Copper:  user.Copper,

		RoleIds: make([]int, 0, 2),
	}

	// TODO: 先每次都记录登录时间
	ip := ctx.Value("ip")
	go self.RecordLogin(user.Username, ip)

	if user.IsRoot {
		me.IsAdmin = true
		return me
	}

	// 获取角色信息
	userRoleList := make([]*model.UserRole, 0)
	err = MasterDB.Where("uid=?", user.Uid).Find(&userRoleList)
	if err != nil {
		logger.Errorf("获取用户 %q 角色 信息失败：%s", username, err)
		return me
	}
	for _, userRole := range userRoleList {
		me.RoleIds = append(me.RoleIds, userRole.Roleid)

		if userRole.Roleid <= model.AdminMinRoleId {
			// 是管理员
			me.IsAdmin = true
		}
	}

	return me
}

// findUsers 获得用户信息，包内使用。
// s 是包含用户 UID 的二维数组
func (self UserLogic) findUsers(ctx context.Context, s interface{}) []*model.User {
	objLog := GetLogger(ctx)

	uids := slices.StructsIntSlice(s, "Uid")

	users := make([]*model.User, 0)
	if err := MasterDB.In("uid", uids).Find(&users); err != nil {
		objLog.Errorln("user logic findUsers not record found:", err)
		return nil
	}
	return users
}

func (self UserLogic) findUser(ctx context.Context, uid int) *model.User {
	objLog := GetLogger(ctx)

	user := &model.User{}
	_, err := MasterDB.Id(uid).Get(user)
	if err != nil {
		objLog.Errorln("user logic findUser not record found:", err)
	}

	return user
}

// 会员总数
func (UserLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.User))
	if err != nil {
		logger.Errorln("UserLogic Total error:", err)
	}
	return total
}

var (
	ErrUsername = errors.New("用户名不存在")
	ErrPasswd   = errors.New("密码错误")
)

// Login 登录；成功返回用户登录信息(user_login)
func (self UserLogic) Login(ctx context.Context, username, passwd string) (*model.UserLogin, error) {
	objLog := GetLogger(ctx)

	userLogin := &model.UserLogin{}
	_, err := MasterDB.Where("username=? OR email=?", username, username).Get(userLogin)
	if err != nil {
		objLog.Errorf("user %q login failure: %s", username, err)
		return nil, errors.New("内部错误，请稍后再试！")
	}
	// 校验用户
	if userLogin.Uid == 0 {
		objLog.Infof("user %q is not exists!", username)
		return nil, ErrUsername
	}

	// 检验用户状态是否正常（未激活的可以登录，但不能发布信息）
	user := &model.User{}
	MasterDB.Id(userLogin.Uid).Get(user)
	if user.Status > model.UserStatusAudit {
		objLog.Infof("用户 %q 的状态非审核通过, 用户的状态值：%d", username, user.Status)
		var errMap = map[int]error{
			model.UserStatusRefuse: errors.New("您的账号审核拒绝"),
			model.UserStatusFreeze: errors.New("您的账号因为非法发布信息已被冻结，请联系管理员！"),
			model.UserStatusOutage: errors.New("您的账号因为非法发布信息已被停号，请联系管理员！"),
		}
		return nil, errMap[user.Status]
	}

	md5Passwd := goutils.Md5(passwd + userLogin.Passcode)
	objLog.Debugf("passwd: %s, passcode: %s, md5passwd: %s, dbpasswd: %s", passwd, userLogin.Passcode, md5Passwd, userLogin.Passwd)
	if md5Passwd != userLogin.Passwd {
		objLog.Infof("用户名 %q 填写的密码错误", username)
		return nil, ErrPasswd
	}

	go func() {
		self.IncrUserWeight("uid", userLogin.Uid, 1)
		ip := ctx.Value("ip")
		self.RecordLogin(username, ip)
	}()

	return userLogin, nil
}

// UpdatePasswd 更新用户密码
func (self UserLogic) UpdatePasswd(ctx context.Context, username, curPasswd, newPasswd string) (string, error) {
	userLogin := &model.UserLogin{}
	_, err := MasterDB.Where("username=?", username).Get(userLogin)
	if err != nil {
		return "用户不存在", err
	}

	if userLogin.Passwd != "" {
		_, err = self.Login(ctx, username, curPasswd)
		if err != nil {
			return "原密码填写错误", err
		}
	}

	userLogin = &model.UserLogin{
		Passwd: newPasswd,
	}
	err = userLogin.GenMd5Passwd()
	if err != nil {
		return err.Error(), err
	}

	changeData := map[string]interface{}{
		"passwd":   userLogin.Passwd,
		"passcode": userLogin.Passcode,
	}
	_, err = MasterDB.Table(userLogin).Where("username=?", username).Update(changeData)
	if err != nil {
		logger.Errorf("用户 %s 更新密码错误：%s", username, err)
		return "对不起，内部服务错误！", err
	}
	return "", nil
}

func (UserLogic) HasPasswd(ctx context.Context, uid int) bool {
	userLogin := &model.UserLogin{}
	_, err := MasterDB.Where("uid=?", uid).Get(userLogin)
	if err == nil && userLogin.Passwd != "" {
		return true
	}

	return false
}

func (self UserLogic) ResetPasswd(ctx context.Context, email, passwd string) (string, error) {
	objLog := GetLogger(ctx)

	userLogin := &model.UserLogin{
		Passwd: passwd,
	}
	err := userLogin.GenMd5Passwd()
	if err != nil {
		return err.Error(), err
	}

	changeData := map[string]interface{}{
		"passwd":   userLogin.Passwd,
		"passcode": userLogin.Passcode,
	}
	_, err = MasterDB.Table(userLogin).Where("email=?", email).Update(changeData)
	if err != nil {
		objLog.Errorf("用户 %s 更新密码错误：%s", email, err)
		return "对不起，内部服务错误！", err
	}
	return "", nil
}

// Activate 用户激活
func (self UserLogic) Activate(ctx context.Context, email, uuid string, timestamp int64, sign string) (*model.User, error) {
	objLog := GetLogger(ctx)

	realSign := DefaultEmail.genActivateSign(email, uuid, timestamp)
	if sign != realSign {
		return nil, errors.New("签名非法！")
	}

	user := self.FindOne(ctx, "email", email)
	if user.Uid == 0 {
		return nil, errors.New("邮箱非法")
	}

	user.Status = model.UserStatusAudit

	_, err := MasterDB.Id(user.Uid).Update(user)
	if err != nil {
		objLog.Errorf("activate [%s] failure:%s", email, err)
		return nil, err
	}

	return user, nil
}

// 增加或减少用户活跃度
func (UserLogic) IncrUserWeight(field string, value interface{}, weight int) {
	_, err := MasterDB.Where(field+"=?", value).Incr("weight", weight).Update(new(model.UserActive))
	if err != nil {
		logger.Errorln("UserActive update Error:", err)
	}
}

func (UserLogic) DecrUserWeight(field string, value interface{}, divide int) {
	if divide <= 0 {
		return
	}

	strSql := fmt.Sprintf("UPDATE user_active SET weight=weight/%d WHERE %s=?", divide, field)
	if result, err := MasterDB.Exec(strSql, value); err != nil {
		logger.Errorln("UserActive update Error:", err)
	} else {
		n, _ := result.RowsAffected()
		logger.Debugln(strSql, "affected num:", n)
	}
}

// RecordLogin 记录用户最后登录时间和 IP
func (UserLogic) RecordLogin(username string, ipinter interface{}) error {
	change := map[string]interface{}{
		"login_time": time.Now(),
	}
	if ip, ok := ipinter.(string); ok && ip != "" {
		change["login_ip"] = ip
	}
	_, err := MasterDB.Table(new(model.UserLogin)).Where("username=?", username).
		Update(change)
	if err != nil {
		logger.Errorf("记录用户 %q 登录错误：%s", username, err)
	}
	return err
}

// FindActiveUsers 获得活跃用户
func (UserLogic) FindActiveUsers(ctx context.Context, limit int, offset ...int) []*model.UserActive {
	objLog := GetLogger(ctx)

	activeUsers := make([]*model.UserActive, 0)
	err := MasterDB.OrderBy("weight DESC").Limit(limit, offset...).Find(&activeUsers)
	if err != nil {
		objLog.Errorln("UserLogic FindActiveUsers error:", err)
		return nil
	}
	return activeUsers
}

func (UserLogic) FindDAUUsers(ctx context.Context, uids []int) map[int]*model.User {
	objLog := GetLogger(ctx)

	users := make(map[int]*model.User)
	err := MasterDB.In("uid", uids).Find(&users)
	if err != nil {
		objLog.Errorln("UserLogic FindDAUUsers error:", err)
		return nil
	}
	return users
}

// FindNewUsers 最新加入会员
func (UserLogic) FindNewUsers(ctx context.Context, limit int, offset ...int) []*model.User {
	objLog := GetLogger(ctx)

	users := make([]*model.User, 0)
	err := MasterDB.OrderBy("ctime DESC").Limit(limit, offset...).Find(&users)
	if err != nil {
		objLog.Errorln("UserLogic FindNewUsers error:", err)
		return nil
	}
	return users
}

// 获取用户列表（分页）：后台用
func (UserLogic) FindUserByPage(ctx context.Context, conds map[string]string, curPage, limit int) ([]*model.User, int) {
	objLog := GetLogger(ctx)

	session := MasterDB.NewSession()

	for k, v := range conds {
		session.And(k+"=?", v)
	}

	totalSession := session.Clone()

	offset := (curPage - 1) * limit
	userList := make([]*model.User, 0)
	err := session.OrderBy("uid DESC").Limit(limit, offset).Find(&userList)
	if err != nil {
		objLog.Errorln("UserLogic find error:", err)
		return nil, 0
	}

	total, err := totalSession.Count(new(model.User))
	if err != nil {
		objLog.Errorln("UserLogic find count error:", err)
		return nil, 0
	}

	return userList, int(total)
}

func (self UserLogic) SetDauAuth(ctx context.Context, uid string, form url.Values) {
	user := self.FindOne(ctx, "uid", uid)
	user.DauAuth = 0

	for k := range form {
		switch k {
		case "topic":
			user.DauAuth |= model.DauAuthTopic
		case "article":
			user.DauAuth |= model.DauAuthArticle
		case "resource":
			user.DauAuth |= model.DauAuthResource
		case "project":
			user.DauAuth |= model.DauAuthProject
		case "wiki":
			user.DauAuth |= model.DauAuthWiki
		case "book":
			user.DauAuth |= model.DauAuthBook
		case "comment":
			user.DauAuth |= model.DauAuthComment
		case "top":
			user.DauAuth |= model.DauAuthTop
		}
	}

	MasterDB.Id(user.Uid).Update(user)
}

// GetUserMentions 获取 @ 的 suggest 列表
func (UserLogic) GetUserMentions(term string, limit int, isHttps bool) []map[string]string {
	userActives := make([]*model.UserActive, 0)
	err := MasterDB.Where("username like ?", "%"+term+"%").Desc("mtime").Limit(limit).Find(&userActives)
	if err != nil {
		logger.Errorln("UserLogic GetUserMentions Error:", err)
		return nil
	}

	users := make([]map[string]string, len(userActives))
	for i, userActive := range userActives {
		user := make(map[string]string, 2)
		user["username"] = userActive.Username
		user["avatar"] = util.Gravatar(userActive.Avatar, userActive.Email, 20, isHttps)
		users[i] = user
	}

	return users
}

// 获取 loginTime 之前没有登录的用户
func (UserLogic) FindNotLoginUsers(loginTime time.Time) (userList []*model.UserLogin, err error) {
	userList = make([]*model.UserLogin, 0)
	err = MasterDB.Where("login_time<?", loginTime).Find(&userList)
	return
}

// 邮件订阅或取消订阅
func (UserLogic) EmailSubscribe(ctx context.Context, uid, unsubscribe int) {
	_, err := MasterDB.Table(&model.User{}).Id(uid).Update(map[string]interface{}{"unsubscribe": unsubscribe})
	if err != nil {
		logger.Errorln("user:", uid, "Email Subscribe Error:", err)
	}
}

func (UserLogic) FindBindUsers(ctx context.Context, uid int) []*model.BindUser {
	bindUsers := make([]*model.BindUser, 0)
	err := MasterDB.Where("uid=?", uid).Find(&bindUsers)
	if err != nil {
		logger.Errorln("user:", uid, "FindBindUsers Error:", err)
	}
	return bindUsers
}

func (UserLogic) doCreateUser(ctx context.Context, session *xorm.Session, user *model.User, passwd ...string) error {

	if user.Avatar == "" && len(DefaultAvatars) > 0 {
		// 随机给一个默认头像
		user.Avatar = DefaultAvatars[rand.Intn(len(DefaultAvatars))]
	}
	user.Open = 0

	user.DauAuth = model.DefaultAuth

	_, err := session.Insert(user)
	if err != nil {
		return err
	}

	// 存用户登录信息
	userLogin := &model.UserLogin{
		Email:    user.Email,
		Username: user.Username,
		Uid:      user.Uid,
	}
	if len(passwd) > 0 {
		userLogin.Passwd = passwd[0]
		err = userLogin.GenMd5Passwd()
		if err != nil {
			return err
		}
	}

	if _, err = session.Insert(userLogin); err != nil {
		return err
	}

	if !user.IsRoot {
		// 存用户角色信息
		userRole := &model.UserRole{}
		// 默认为初级会员
		userRole.Roleid = Roles[len(Roles)-1].Roleid
		userRole.Uid = user.Uid
		if _, err = session.Insert(userRole); err != nil {
			return err
		}
	}

	// 存用户活跃信息，初始活跃+2
	userActive := &model.UserActive{
		Uid:      user.Uid,
		Username: user.Username,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Weight:   2,
	}
	if _, err = session.Insert(userActive); err != nil {
		return err
	}

	return nil
}
