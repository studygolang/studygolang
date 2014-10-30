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

// 所有晨读（分页）
// /admin/reading/list
func ReadingListHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	readings, total := service.FindReadingByPage(nil, curPage, limit)

	if readings == nil {
		logger.Errorln("[ReadingListHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   readings,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/reading/list.html,/template/admin/reading/query.html")
	filter.SetData(req, data)
}

// /admin/reading/query.html
func ReadingQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"id", "rtype"})

	readings, total := service.FindReadingByPage(conds, curPage, limit)

	if readings == nil {
		logger.Errorln("[ReadingQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   readings,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/reading/query.html")
	filter.SetData(req, data)
}

// /admin/reading/publish
func PublishReadingHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)

		errMsg, err := service.SaveReading(req.PostForm, user["username"].(string))
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "操作成功"
		}
	} else {
		id, err := strconv.Atoi(req.FormValue("id"))
		if err == nil && id != 0 {
			data["reading"], err = service.FindReadingById(id)

			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/reading/modify.html")
	}

	filter.SetData(req, data)
}
