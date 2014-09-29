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
)

// 所有帖子（分页）
// /admin/community/topic/list
func TopicListHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	topics, total := service.FindTopicsByPage(nil, curPage, limit)

	if topics == nil {
		logger.Errorln("[TopicListHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   topics,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/topic/list.html,/template/admin/topic/query.html")
	filter.SetData(req, data)
}

// /admin/community/topic/query.html
func TopicQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"title"})

	topics, total := service.FindTopicsByPage(conds, curPage, limit)

	if topics == nil {
		logger.Errorln("[TopicQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   topics,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/topic/query.html")
	filter.SetData(req, data)
}

func ModifyTopicHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)

		errMsg, err := service.ModifyTopic(user, req.PostForm)
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "修改成功"
		}
	} else {
		topic, replies, err := service.FindTopicByTid(req.FormValue("tid"))

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/topic/modify.html")
		data["topic"] = topic
		data["replies"] = replies
		data["nodes"] = service.GenNodes()
	}

	filter.SetData(req, data)
}

func DelTopicHandler(rw http.ResponseWriter, req *http.Request) {
}
