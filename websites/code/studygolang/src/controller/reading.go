// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"net/http"
	"strconv"

	"filter"
	"github.com/studygolang/mux"
	"model"
	"service"
	"util"
)

// 晨读列表页
// uri: /readings
func ReadingsHandler(rw http.ResponseWriter, req *http.Request) {
	limit := 20

	lastId := req.FormValue("lastid")
	if lastId == "" {
		lastId = "0"
	}

	rtype, err := strconv.Atoi(req.FormValue("rtype"))
	if err != nil {
		rtype = model.RtypeGo
	}

	readings := service.FindReadings(lastId, "25", rtype)
	if readings == nil {
		// TODO:服务暂时不可用？
	}

	num := len(readings)
	if num == 0 {
		if lastId == "0" {
			util.Redirect(rw, req, "/")
		} else {
			util.Redirect(rw, req, "/readings")
		}

		return
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId != "0" {
		prevId, _ = strconv.Atoi(lastId)

		// 避免因为项目下线，导致判断错误（所以 > 5）
		if prevId-readings[0].Id > 5 {
			hasPrev = false
		} else {
			prevId += limit
			hasPrev = true
		}
	}

	if num > limit {
		hasNext = true
		readings = readings[:limit]
		nextId = readings[limit-1].Id
	} else {
		nextId = readings[num-1].Id
	}

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/readings/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeReadings": "active", "readings": readings, "page": pageInfo, "rtype": rtype})
}

// 点击 【我要晨读】，记录点击数，跳转
// uri: /readings/{id:[0-9]+}
func IReadingHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	url := service.IReading(vars["id"])

	util.Redirect(rw, req, url)
	return
}
