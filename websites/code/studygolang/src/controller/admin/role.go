// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"filter"
	"html/template"
	"logger"
	"net/http"
	"service"
	"strconv"
)

// 所有角色（分页）
func RoleListHandler(rw http.ResponseWriter, req *http.Request) {

	curPage, limit := parsePage(req)

	total := len(service.Roles)
	newLimit := limit
	if total < limit {
		newLimit = total
	}

	data := map[string]interface{}{
		"datalist":   service.Roles[curPage:newLimit],
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage + 1,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/role/list.html,/template/admin/role/query.html")
	filter.SetData(req, data)
}

func RoleQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"name"})

	roles, total := service.FindRolesByPage(conds, curPage, limit)

	if roles == nil {
		logger.Errorln("[RoleQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	tpl, err := template.ParseFiles(ROOT+"/template/admin/common_query.html", ROOT+"/template/admin/role/query.html")
	if err != nil {
		logger.Errorln("[RoleQueryHandler] parse file error:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   roles,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage + 1,
		"limit":      limit,
	}

	err = tpl.Execute(rw, data)
	if err != nil {
		logger.Errorln("[RoleQueryHandler] execute file error:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func NewRoleHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)
		username := user["username"].(string)

		errMsg, err := service.SaveRole(req.PostForm, username)
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "添加成功"
		}
	} else {
		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/role/new.html")
	}

	filter.SetData(req, data)
}

func ModifyRoleHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)
		username := user["username"].(string)

		errMsg, err := service.SaveRole(req.PostForm, username)
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "修改成功"
		}
	} else {
		role := service.FindRole(req.FormValue("roleid"))

		if role == nil || role.Roleid == 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/role/modify.html")

		data["role"] = role
	}

	filter.SetData(req, data)
}

func DelRoleHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	roleid := req.FormValue("roleid")

	if _, err := strconv.Atoi(roleid); err != nil {
		data["ok"] = 0
		data["error"] = "aid不是整型"

		filter.SetData(req, data)
		return
	}

	if err := service.DelRole(roleid); err != nil {
		data["ok"] = 0
		data["error"] = "删除失败！"
	} else {
		data["ok"] = 1
		data["msg"] = "删除成功！"
	}

	filter.SetData(req, data)
}
