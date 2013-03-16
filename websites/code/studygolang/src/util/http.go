// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package util

import (
	"github.com/gorilla/context"
	"net/http"
)

func Redirect(rw http.ResponseWriter, req *http.Request, uri string) {
	// 避免跳转，context中没有清除
	context.Clear(req)

	http.Redirect(rw, req, uri, http.StatusFound)
}
