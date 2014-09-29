// Copyright 2014 The StudyGolang Authors. All rights reserved.
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

// 所有规则（分页）
// /admin/crawl/rule/list
func RuleListHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	rules, total := service.FindRuleByPage(nil, curPage, limit)

	if rules == nil {
		logger.Errorln("[RuleListHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   rules,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/rule/list.html,/template/admin/rule/query.html")
	filter.SetData(req, data)
}

// /admin/crawl/rule/query.html
func RuleQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"domain"})

	rules, total := service.FindRuleByPage(conds, curPage, limit)

	if rules == nil {
		logger.Errorln("[RuleQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   rules,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/rule/query.html")
	filter.SetData(req, data)
}

// 新建规则
// /admin/crawl/rule/new
func NewRuleHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)
		username := user["username"].(string)

		errMsg, err := service.SaveRule(req.PostForm, username)
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "添加成功"
		}
	} else {
		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/rule/new.html")
	}

	filter.SetData(req, data)
}

func ModifyRuleHandler(rw http.ResponseWriter, req *http.Request) {
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

func DelRuleHandler(rw http.ResponseWriter, req *http.Request) {
}
