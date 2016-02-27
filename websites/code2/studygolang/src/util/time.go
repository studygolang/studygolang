// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package util

import (
	"time"
)

const TIME_LAYOUT_OFTEN = "2006-01-02 15:04:05"
const DATE_LAYOUT_OFTEN = "060102"

// 解析常用的日期时间格式：2014-01-11 16:18:00，东八区
func TimeParseOften(value string) (time.Time, error) {
	local, _ := time.LoadLocation("Local")
	return time.ParseInLocation(TIME_LAYOUT_OFTEN, value, local)
}

func TimeNow() string {
	return time.Now().Format(TIME_LAYOUT_OFTEN)
}

func DateNow() string {
	return time.Now().Format(DATE_LAYOUT_OFTEN)
}
