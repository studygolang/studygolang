// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"net/http"

	"filter"
	"service"
)

// uri: /admin/tool/sitemap
func GenSitemapHandler(rw http.ResponseWriter, req *http.Request) {
	service.GenSitemap()
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/tool/sitemap.html")
}
