// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package util

import (
	"fmt"
	"strconv"
	"strings"
)

// 必须是int类型，否则panic
func MustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

// 将in slice通过sep连接起来
func Join(ins []int, sep string) string {
	strSlice := make([]string, len(ins))
	for i, one := range ins {
		strSlice[i] = strconv.Itoa(one)
	}
	return strings.Join(strSlice, sep)
}

// 获取头像
func Gravatar(avatar string, emailI interface{}, size uint16) string {
	if avatar != "" {
		return fmt.Sprintf("http://studygolang.qiniudn.com/avatar/%s?imageView2/2/w/%d", avatar, size)
	}

	email, ok := emailI.(string)
	if !ok {
		return fmt.Sprintf("http://studygolang.qiniudn.com/avatar/gopher28.png?imageView2/2/w/%d", size)
	}
	return fmt.Sprintf("http://gravatar.duoshuo.com/avatar/%s?s=%d", Md5(email), size)
}
