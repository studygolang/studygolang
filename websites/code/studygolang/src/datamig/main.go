// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	"flag"
	"model"
	"service"
)

func main() {
	var uid, nid int
	var title, content string
	flag.IntVar(&uid, "u", 1, "用户uid")
	flag.IntVar(&nid, "n", 1, "节点nid")
	flag.StringVar(&title, "t", "", "帖子标题")
	flag.StringVar(&content, "c", "", "帖子内容")
	flag.Parse()

	// 入库
	topic := model.NewTopic()
	topic.Uid = uid
	topic.Nid = nid
	topic.Title = title
	topic.Content = content
	_, err := service.PublishTopic(topic)
	if err != nil {
		panic(err)
	}
}
