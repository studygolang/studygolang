// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"global"
	"logger"
	"model"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"util"
)

var (
	roleLocker sync.RWMutex
	Roles      []*model.Role
)

func FindRolesByPage(conds map[string]string, curPage, limit int) ([]*model.Role, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	role := model.NewRole()

	limitStr := strconv.Itoa(curPage*limit) + "," + strconv.Itoa(limit)
	roles, err := role.Where(strings.Join(conditions, " AND ")).Limit(limitStr).
		FindAll()
	if err != nil {
		return nil, 0
	}

	total, err := role.Count()
	if err != nil {
		return nil, 0
	}

	return roles, total
}

func FindRole(roleid string) *model.Role {
	if roleid == "" {
		return nil
	}

	role := model.NewRole()
	err := role.Where("roleid=" + roleid).Find()
	if err != nil {
		logger.Errorln("role FindRole error:", err)
		return nil
	}

	return role
}

func SaveRole(form url.Values, opUser string) (errMsg string, err error) {

	role := model.NewRole()
	role.Name = form.Get("name")
	role.OpUser = opUser

	roleid := form.Get("roleid")
	isNew := roleid == ""
	if isNew {
		role.Ctime = util.TimeNow()

		_, err = role.Insert()
	} else {
		role.Roleid, err = strconv.Atoi(roleid)
		if err != nil {
			errMsg = "roleid invalid"
			logger.Errorln(errMsg, ":", err)
			return
		}
		err = role.Persist(role)
	}

	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln(errMsg, ":", err)
		return
	}

	roleAuth := model.NewRoleAuthority()
	if !isNew {
		// 如果是更新角色，将之前的角色权限都删除
		roleAuth.Where("roleid=" + strconv.Itoa(role.Roleid)).Delete()
	}

	roleAuth.Roleid = role.Roleid
	roleAuth.OpUser = opUser

	// 增加角色拥有的权限
	for _, aid := range form["authorities[]"] {
		aid, err := strconv.Atoi(aid)
		if err != nil {
			continue
		}
		roleAuth.Aid = aid

		roleAuth.Insert()
	}

	global.RoleChan <- struct{}{}
	global.RoleAuthChan <- struct{}{}

	return
}

func DelRole(roleid string) error {
	err := model.NewRole().Where("roleid=" + roleid).Delete()
	model.NewRoleAuthority().Where("roleid=" + roleid).Delete()

	global.RoleChan <- struct{}{}
	global.RoleAuthChan <- struct{}{}

	return err
}

// 将所有 角色 加载到内存中；后台修改角色时，重新加载一次
func LoadRoles() error {
	roles, err := model.NewRole().FindAll()
	if err != nil {
		logger.Errorln("LoadRoles role read fail:", err)
		return err
	}

	roleLocker.Lock()
	defer roleLocker.Unlock()

	Roles = roles

	return nil
}
