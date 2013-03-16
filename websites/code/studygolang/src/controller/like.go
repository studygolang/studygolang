// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

// 喜欢系统

import (
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"net/http"
	"service"
	"util"
)

// 喜欢（或取消喜欢）
// uri: /like/{objid:[0-9]+}.json
func LikeHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	user, _ := filter.CurrentUser(req)
	// 入库
	err := service.PostComment(util.MustInt(vars["objid"]), util.MustInt(req.FormValue("objtype")), user["uid"].(int), req.FormValue("content"), req.FormValue("objname"))
	if err != nil {
		fmt.Fprint(rw, `{"errno": 1, "error":"服务器内部错误"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "error":""}`)
}
