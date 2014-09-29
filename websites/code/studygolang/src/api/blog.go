// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package api

import (
	"fmt"
	"net/http"
)

// 登录
// uri : /api/blog/category/all
func BlogCategoryHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, "true")
}
