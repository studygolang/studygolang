// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service_test

import (
	. "service"
	"testing"
)

func TestFindNewBlogs(t *testing.T) {
	blogList := FindNewBlogs()
	if len(blogList) == 0 {
		t.Fatal("xxxx")
	}
	t.Log(len(blogList))
	for k, blog := range blogList {
		t.Log(k, blog)
		t.Log("===")
	}
}
