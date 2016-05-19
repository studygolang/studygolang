// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

// 图片
type Image struct {
	Pid    int `xorm:"pk autoincr"`
	Md5    string
	Path   string
	Size   int
	Width  int
	Height int
}
