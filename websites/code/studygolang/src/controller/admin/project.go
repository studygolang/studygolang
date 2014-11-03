// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"net/http"
	"strconv"

	"filter"
	"logger"
	"service"
)

// 所有开源项目（分页）
// /admin/community/project/list
func ProjectListHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	projects, total := service.FindProjectByPage(nil, curPage, limit)

	if projects == nil {
		logger.Errorln("[ProjectListHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   projects,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/project/list.html,/template/admin/project/query.html")
	filter.SetData(req, data)
}

// /admin/community/project/query.html
func ProjectQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"id", "domain", "title"})

	projects, total := service.FindProjectByPage(conds, curPage, limit)

	if projects == nil {
		logger.Errorln("[ProjectQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   projects,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/project/query.html")
	filter.SetData(req, data)
}

// 更新状态
// uri: /admin/community/project/update_status
func UpdateProjectStatusHandler(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PostFormValue("id"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	status, err := strconv.Atoi(req.FormValue("status"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var data = make(map[string]interface{})

	user, _ := filter.CurrentUser(req)

	err = service.UpdateProjectStatus(id, status, user["username"].(string))
	if err != nil {
		logger.Errorln("UpdateProjectStatusHandler error:", err)
		data["ok"] = 0
		data["error"] = "更新状态失败"
	} else {
		data["ok"] = 1
		data["msg"] = "更新状态成功"
	}

	filter.SetData(req, data)
}
