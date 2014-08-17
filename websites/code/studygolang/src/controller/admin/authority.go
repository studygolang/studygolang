// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"encoding/json"
	"filter"
	"html/template"
	"logger"
	"net/http"
	"service"
	"strconv"
)

// 所有权限（分页）
func AuthListHandler(rw http.ResponseWriter, req *http.Request) {

	curPage, limit := parsePage(req)

	total := len(service.Authorities)
	newLimit := limit
	if total < limit {
		newLimit = total
	}

	data := map[string]interface{}{
		"datalist":   service.Authorities[curPage:newLimit],
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage + 1,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/authority/list.html,/template/admin/authority/query.html")
	filter.SetData(req, data)
}

func AuthQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"route", "name"})

	authorities, total := service.FindAuthoritiesByPage(conds, curPage, limit)

	if authorities == nil {
		logger.Errorln("[AuthQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	tpl, err := template.ParseFiles(ROOT+"/template/admin/common_query.html", ROOT+"/template/admin/authority/query.html")
	if err != nil {
		logger.Errorln("[AuthQueryHandler] parse file error:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   authorities,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage + 1,
		"limit":      limit,
	}

	err = tpl.Execute(rw, data)
	if err != nil {
		logger.Errorln("[AuthQueryHandler] execute file error:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func NewAuthorityHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)
		username := user["username"].(string)

		errMsg, err := service.SaveAuthority(req.PostForm, username)
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "添加成功"
		}
	} else {
		menu1, menu2 := service.GetMenus()
		allmenu2, _ := json.Marshal(menu2)

		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/authority/new.html")
		data["allmenu1"] = menu1
		data["allmenu2"] = string(allmenu2)
	}

	filter.SetData(req, data)
}

func ModifyAuthorityHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)
		username := user["username"].(string)

		errMsg, err := service.SaveAuthority(req.PostForm, username)
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "修改成功"
		}
	} else {
		menu1, menu2 := service.GetMenus()
		allmenu2, _ := json.Marshal(menu2)

		authority := service.FindAuthority(req.FormValue("aid"))

		if authority == nil || authority.Aid == 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/authority/modify.html")
		data["allmenu1"] = menu1
		data["allmenu2"] = string(allmenu2)
		data["authority"] = authority
	}

	filter.SetData(req, data)
}

func DelAuthorityHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	aid := req.FormValue("aid")

	if _, err := strconv.Atoi(aid); err != nil {
		data["ok"] = 0
		data["error"] = "aid不是整型"

		filter.SetData(req, data)
		return
	}

	if err := service.DelAuthority(aid); err != nil {
		data["ok"] = 0
		data["error"] = "删除失败！"
	} else {
		data["ok"] = 1
		data["msg"] = "删除成功！"
	}

	filter.SetData(req, data)
}
