// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

// 喜欢系统

import (
	"fmt"
	"net/http"

	"filter"
	"github.com/studygolang/mux"
	"service"
	"util"
)

// 喜欢（或取消喜欢）
// uri: /like/{objid:[0-9]+}.json
func LikeHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	user, _ := filter.CurrentUser(req)

	if !util.CheckInt(req.PostForm, "objtype") || !util.CheckInt(req.PostForm, "flag") {
		fmt.Fprint(rw, `{"ok": 0, "error":"参数错误"}`)
		return
	}

	uid := user["uid"].(int)
	objid := util.MustInt(vars["objid"])
	objtype := util.MustInt(req.PostFormValue("objtype"))
	likeFlag := util.MustInt(req.PostFormValue("flag"))

	err := service.LikeObject(uid, objid, objtype, likeFlag)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"服务器内部错误"}`)
		return
	}

	fmt.Fprint(rw, `{"ok": 1, "msg":"success", "data":""}`)
}
