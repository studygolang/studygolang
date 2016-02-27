// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"model"

	"golang.org/x/net/context"

	. "db"
)

type UserLogic struct{}

var DefaultUserLogic = UserLogic{}

func (self UserLogic) FindUserInfos(ctx context.Context, uids []int) map[int]*model.User {
	objLog := GetLogger(ctx)

	var users []*model.User
	if DB.Where("uid in (?)", uids).Find(&users).RecordNotFound() {
		objLog.Infoln("user logic FindAll not record found:")
		return nil
	}

	usersMap := make(map[int]*model.User, len(users))
	for _, user := range users {
		if user == nil || user.Uid == 0 {
			continue
		}
		usersMap[user.Uid] = user
	}
	return usersMap
}
