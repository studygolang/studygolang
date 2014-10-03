// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package filter

import (
	"config"
	"fmt"
	"github.com/studygolang/mux"
	"html/template"
	"logger"
	"net/http"
	"runtime"
	"service"
)

// 管理后台权限检查过滤器
type AdminFilter struct {
	*mux.EmptyFilter
}

func (this *AdminFilter) PreFilter(rw http.ResponseWriter, req *http.Request) bool {
	if user, ok := CurrentUser(req); ok {
		/*
			// 是管理员才能查看后台
			if isAdmin, ok := user["isadmin"].(bool); !ok || !isAdmin {
				return false
			}
		*/
		if req.URL.Path == "/admin" {
			return true
		}

		return service.HasAuthority(user["uid"].(int), req.URL.Path)
	}
	return true
}

// 没有权限时，返回 403
func (this *AdminFilter) PreErrorHandle(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusForbidden)

	tpl, err := template.ParseFiles(config.ROOT+"/template/admin/simple_base.html", config.ROOT+"/template/admin/403.html")
	if err != nil {
		logger.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", req.RequestURI, err)
		fmt.Fprint(rw, "403 Forbidden<br/>Go "+runtime.Version())
		return
	}

	// 当前用户信息
	me, _ := CurrentUser(req)
	data := map[string]interface{}{"me": me}

	if err = tpl.Execute(rw, data); err != nil {
		logger.Errorf("执行模板出错（Execute）：[%q] %s\n", req.RequestURI, err)
		fmt.Fprint(rw, "403 Forbidden<br/>Go "+runtime.Version())
		return
	}
}
