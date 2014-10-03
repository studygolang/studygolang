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
	"util"
)

// 获取用户菜单
func GetUserMenu(uid int, uri string) ([]*model.Authority, map[int][]*model.Authority, int) {
	aidMap, err := userAuthority(strconv.Itoa(uid))
	if err != nil {
		return nil, nil, 0
	}

	authLocker.RLock()
	defer authLocker.RUnlock()

	userMenu1 := make([]*model.Authority, 0, 4)
	userMenu2 := make(map[int][]*model.Authority)
	curMenu1 := 0

	for _, authority := range Authorities {
		if _, ok := aidMap[authority.Aid]; ok {
			if authority.Menu1 == 0 {
				userMenu1 = append(userMenu1, authority)
				userMenu2[authority.Aid] = make([]*model.Authority, 0, 4)
			} else if authority.Menu2 == 0 {
				userMenu2[authority.Menu1] = append(userMenu2[authority.Menu1], authority)
			}
			if authority.Route == uri {
				curMenu1 = authority.Menu1
			}
		}
	}

	return userMenu1, userMenu2, curMenu1
}

// 获取整个菜单
func GetMenus() ([]*model.Authority, map[string][][]string) {
	var (
		menu1 = make([]*model.Authority, 0, 10)
		menu2 = make(map[string][][]string)
	)

	for _, authority := range Authorities {
		if authority.Menu1 == 0 {
			menu1 = append(menu1, authority)
			aid := strconv.Itoa(authority.Aid)
			menu2[aid] = make([][]string, 0, 4)
		} else if authority.Menu2 == 0 {
			m := strconv.Itoa(authority.Menu1)
			oneMenu2 := []string{strconv.Itoa(authority.Aid), authority.Name}
			menu2[m] = append(menu2[m], oneMenu2)
		}
	}

	return menu1, menu2
}

// 除了一级、二级菜单之外的权限（路由）
func GeneralAuthorities() map[int][]*model.Authority {
	auths := make(map[int][]*model.Authority)

	for _, authority := range Authorities {
		if authority.Menu1 == 0 {
			auths[authority.Aid] = make([]*model.Authority, 0, 8)
		} else if authority.Menu2 != 0 {
			auths[authority.Menu1] = append(auths[authority.Menu1], authority)
		}
	}

	return auths
}

// 判断用户是否有某个权限
func HasAuthority(uid int, route string) bool {
	aidMap, err := userAuthority(strconv.Itoa(uid))
	if err != nil {
		logger.Errorln("HasAuthority:Read user authority error:", err)
		return false
	}

	authLocker.RLock()
	defer authLocker.RUnlock()

	for _, authority := range Authorities {
		if _, ok := aidMap[authority.Aid]; ok {
			if route == authority.Route {
				return true
			}
		}
	}

	return false
}

func FindAuthoritiesByPage(conds map[string]string, curPage, limit int) ([]*model.Authority, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	authority := model.NewAuthority()

	limitStr := strconv.Itoa((curPage-1)*limit) + "," + strconv.Itoa(limit)
	auhtorities, err := authority.Where(strings.Join(conditions, " AND ")).Limit(limitStr).
		FindAll()
	if err != nil {
		return nil, 0
	}

	total, err := authority.Count()
	if err != nil {
		return nil, 0
	}

	return auhtorities, total
}

func FindAuthority(aid string) *model.Authority {
	if aid == "" {
		return nil
	}

	authority := model.NewAuthority()
	err := authority.Where("aid=" + aid).Find()
	if err != nil {
		logger.Errorln("authority FindAuthority error:", err)
		return nil
	}

	return authority
}

func SaveAuthority(form url.Values, opUser string) (errMsg string, err error) {
	authority := model.NewAuthority()
	err = util.ConvertAssign(authority, form)
	if err != nil {
		logger.Errorln("authority ConvertAssign error", err)
		errMsg = err.Error()
		return
	}

	authority.OpUser = opUser

	if authority.Aid != 0 {
		err = authority.Persist(authority)
	} else {
		authority.Ctime = util.TimeNow()

		_, err = authority.Insert()
	}

	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln(errMsg, ":", err)
		return
	}

	global.AuthorityChan <- struct{}{}

	return
}

func DelAuthority(aid string) error {
	err := model.NewAuthority().Where("aid=" + aid).Delete()

	global.AuthorityChan <- struct{}{}

	return err
}

func userAuthority(uid string) (map[int]bool, error) {
	userRoles, err := model.NewUserRole().Where("uid=" + uid).FindAll("roleid")
	if err != nil {
		logger.Errorln("userAuthority userole read fail:", err)
		return nil, err
	}

	roleAuthLocker.RLock()

	aidMap := make(map[int]bool)
	for _, userRole := range userRoles {
		for _, aid := range RoleAuthorities[userRole.Roleid] {
			aidMap[aid] = true
		}
	}

	roleAuthLocker.RUnlock()

	return aidMap, nil
}
