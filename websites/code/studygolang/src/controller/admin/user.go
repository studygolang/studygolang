// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"filter"
	"logger"
	"net/http"
	"service"
	"strconv"
)

// 所有用户（分页）
func UserListHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	users, total := service.FindUsersByPage(nil, curPage, limit)

	if users == nil {
		logger.Errorln("[UsersHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   users,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/user/list.html,/template/admin/user/query.html")
	filter.SetData(req, data)
}

// /admin/user/user/query.html
func UserQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"uid", "name", "email"})

	users, total := service.FindUsersByPage(conds, curPage, limit)

	if users == nil {
		logger.Errorln("[UserQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   users,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/user/query.html")
	filter.SetData(req, data)
}

func UserDetailHandler(rw http.ResponseWriter, req *http.Request) {
	uid := req.FormValue("uid")

	if _, err := strconv.Atoi(uid); err != nil {
		logger.Errorln("[UserDetailHandler] invalid uid")
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	user := service.FindUserByUID(uid)
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/user/detail.html")
	filter.SetData(req, map[string]interface{}{"user": user})
}
