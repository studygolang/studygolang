// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"model"
	"util"

	. "db"

	"golang.org/x/net/context"
)

type TopicLogic struct{}

func (*TopicLogic) FindAll(ctx context.Context) []map[string]interface{} {
	// objLog := GetLogger(ctx)

	var (
		count    = 10
		topics   = make([]*model.Topic, count)
		topicExs = make([]*model.TopicEx, count)
	)
	if DB.Limit(count).Find(&topics).Related(&topicExs, "Tid").RecordNotFound() {
		return nil
	}

	uids := make([]int, 0, count)
	nids := make([]int, count)

	for i, topic := range topics {
		uids = append(uids, topic.Uid)
		if topic.Lastreplyuid != 0 {
			uids = append(uids, topic.Lastreplyuid)
		}
		nids[i] = topic.Nid
	}
	usersMap := DefaultUserLogic.FindUserInfos(ctx, uids)

	data := make([]map[string]interface{}, 10)

	for i, topic := range topics {
		dest := make(map[string]interface{})

		// 有人回复
		if topic.Lastreplyuid != 0 {
			if user, ok := usersMap[topic.Lastreplyuid]; ok {
				dest["lastreplyusername"] = user.Username
			}
		}

		util.Struct2Map(dest, topic)
		util.Struct2Map(dest, topicExs[i])

		dest["user"] = usersMap[topic.Uid]
		// tmpMap["node"] = nodes[tmpMap["nid"].(int)]

		data[i] = dest
	}

	return data
}
