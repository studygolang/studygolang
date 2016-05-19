// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"fmt"
	"time"
)

// Commenter 评论接口
type Commenter interface {
	fmt.Stringer
	// 评论回调接口，用于更新对象自身需要更新的数据
	UpdateComment(int, int, int, time.Time)
}

var commenters = make(map[int]Commenter)

// 注册评论对象，使得某种类型（主题、博客等）被评论了可以回调
func RegisterCommentObject(objtype int, commenter Commenter) {
	if commenter == nil {
		panic("logic: Register commenter is nil")
	}
	if _, dup := commenters[objtype]; dup {
		panic("logic: Register called twice for commenter " + commenter.String())
	}
	commenters[objtype] = commenter
}
