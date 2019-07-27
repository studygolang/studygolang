// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package util

import (
	"net/url"
	"strconv"
)

// 检测提交请求的参数是否是 int 类型
func CheckInt(form url.Values, field string) bool {
	_, err := strconv.Atoi(form.Get(field))
	if err != nil {
		return false
	}

	return true
}
