// Copyright 2015 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"net/http"

	"filter"
)

// Wide 的内嵌 iframe 的 playground
func PlaygroundHandler(rw http.ResponseWriter, req *http.Request) {
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/wide/playground.html")
}
