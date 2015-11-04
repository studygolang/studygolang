// Copyright 2015 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package filter

import (
	"model"
	"net/http"
	"service"

	"config"
	"logger"
	"util"

	"github.com/studygolang/mux"
)

// SensitiveFilter 敏感词过滤器
type SensitiveFilter struct {
	*mux.EmptyFilter
}

// PreFilter 执行 handler 之前的过滤方法
func (this *SensitiveFilter) PreFilter(rw http.ResponseWriter, req *http.Request) bool {
	logger.Debugln("SensitiveFilter PreFilter...")

	content := req.FormValue("content")
	title := req.FormValue("title")
	if title == "" && content == "" {
		return true
	}

	sensitive := config.Config["sensitive"]
	if util.HasSensitive(title, sensitive) || util.HasSensitive(content, sensitive) {
		// 把账号冻结
		curUser, _ := CurrentUser(req)
		service.UpdateUserStatus(curUser["uid"].(int), model.StatusFreeze)

		logger.Infoln("user=", curUser["uid"], "publish ad, title=", title, ";content=", content, ". freeze")

		return false
	}

	return true
}
