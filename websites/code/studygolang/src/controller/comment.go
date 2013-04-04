// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"net/http"
	"service"
	"util"
)

// 评论（或回复）
// uri: /comment/{objid:[0-9]+}.json
func CommentHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	user, _ := filter.CurrentUser(req)
	// 入库
	err := service.PostComment(user["uid"].(int), util.MustInt(vars["objid"]), req.Form)
	if err != nil {
		fmt.Fprint(rw, `{"errno": 1, "error":"服务器内部错误"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "error":""}`)
}
