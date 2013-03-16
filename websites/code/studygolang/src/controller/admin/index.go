// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"config"
	"filter"
	"net/http"
)

var ROOT = config.ROOT

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	user, _ := filter.CurrentUser(req)
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/index.html")
	filter.SetData(req, map[string]interface{}{"user": user})
}
