// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"math/rand"
	"model"
	"net/url"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"

	. "db"
)

var DefaultAvatars = []string{
	"gopher_aqua.jpg", "gopher_boy.jpg", "gopher_brown.jpg", "gopher_gentlemen.jpg",
	"gopher_strawberry.jpg", "gopher_strawberry_bg.jpg", "gopher_teal.jpg",
	"gopher01.png", "gopher02.png", "gopher03.png", "gopher04.png",
	"gopher05.png", "gopher06.png", "gopher07.png", "gopher08.png",
	"gopher09.png", "gopher10.png", "gopher11.png", "gopher12.png",
	"gopher13.png", "gopher14.png", "gopher15.png", "gopher16.png",
	"gopher17.png", "gopher18.png", "gopher19.png", "gopher20.png",
	"gopher21.png", "gopher22.png", "gopher23.png", "gopher24.png",
	"gopher25.png", "gopher26.png", "gopher27.png", "gopher28.png",
}

type UserLogic struct{}

var DefaultUser = UserLogic{}

// CreateUser 创建用户
func (self UserLogic) CreateUser(ctx context.Context, form url.Values) (errMsg string, err error) {
	objLog := GetLogger(ctx)

	if self.UserExists(ctx, "email", form.Get("email")) {
		err = errors.New("该邮箱已注册过")
		return
	}
	if self.UserExists(ctx, "username", form.Get("username")) {
		err = errors.New("用户名已存在")
		return
	}

	user := &model.User{}
	err = schemaDecoder.Decode(user, form)
	if err != nil {
		objLog.Errorln("user schema Decode error:", err)
		errMsg = err.Error()
		return
	}

	// 随机给一个默认头像
	user.Avatar = DefaultAvatars[rand.Intn(len(DefaultAvatars))]
	_, err = MasterDB.Insert(user)
	if err != nil {
		errMsg = "内部服务器错误"
		objLog.Errorln(errMsg, ":", err)
		return
	}

	// 存用户登录信息
	userLogin := &model.UserLogin{}
	err = schemaDecoder.Decode(userLogin, form)
	if err != nil {
		errMsg = err.Error()
		objLog.Errorln("CreateUser error:", err)
		return
	}
	userLogin.Uid = user.Uid
	if _, err = MasterDB.Insert(userLogin); err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln(errMsg, ":", err)
		return
	}

	// 存用户角色信息
	userRole := &model.UserRole{}
	// 默认为初级会员
	userRole.Roleid = Roles[len(Roles)-1].Roleid
	userRole.Uid = user.Uid
	if _, err = MasterDB.Insert(userRole); err != nil {
		objLog.Errorln("userRole insert Error:", err)
	}

	// 存用户活跃信息，初始活跃+2
	userActive := &model.UserActive{}
	userActive.Uid = user.Uid
	userActive.Username = user.Username
	userActive.Avatar = user.Avatar
	userActive.Email = user.Email
	userActive.Weight = 2
	if _, err = MasterDB.Insert(userActive); err != nil {
		objLog.Errorln("UserActive insert Error:", err)
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
	if err != nil || userLogin.Uid != 0 {
		if err != nil {
			objLog.Errorln("user logic EmailOrUsernameExists error:", err)
		}
		return false
	}
	return true
}

func (self UserLogic) FindUserInfos(ctx context.Context, uids []int) map[int]*model.User {
	objLog := GetLogger(ctx)

	var usersMap = make(map[int]*model.User)
	if err := MasterDB.In("uid", uids).Find(&usersMap); err != nil {
		objLog.Infoln("user logic FindAll not record found:")
		return nil
	}

	// usersMap := make(map[int]*model.User, len(users))
	// for _, user := range users {
	// 	if user == nil || user.Uid == 0 {
	// 		continue
	// 	}
	// 	usersMap[user.Uid] = user
	// }
	return usersMap
}

func (self UserLogic) Find(ctx context.Context) {

}
