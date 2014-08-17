// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package filter

import (
	"github.com/studygolang/mux"
	"net/http"
	//"service"
)

// 管理后台权限检查过滤器
type AdminFilter struct {
	*mux.EmptyFilter
}

func (this *AdminFilter) PreFilter(rw http.ResponseWriter, req *http.Request) bool {
	if user, ok := CurrentUser(req); ok {
		// 是管理员才能查看后台
		if isAdmin, ok := user["isadmin"].(bool); !ok || !isAdmin {
			return false
		}
		if req.RequestURI == "/admin" {
			return true
		}

		//return service.HasAuthority(user["uid"].(int), req.RequestURI)
	}
	return true
}

// 没有权限时，返回 403
func (this *AdminFilter) PreErrorHandle(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusForbidden)
}
