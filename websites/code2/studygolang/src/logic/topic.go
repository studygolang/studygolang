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

var DefaultTopic = TopicLogic{}

func (TopicLogic) FindAll(ctx context.Context, paginator *Paginator, orderBy string, querystring string, args ...interface{}) []map[string]interface{} {
	objLog := GetLogger(ctx)

	var (
		count      = paginator.PerPage()
		topicInfos = make([]*model.TopicInfo, 0)
	)

	session := MasterDB.Join("INNER", "topics_ex", "topics.tid=topics_ex.tid")
	if querystring != "" {
		session.Where(querystring, args...)
	}
	err := session.OrderBy(orderBy).Limit(count, paginator.Offset()).Find(&topicInfos)
	if err != nil {
		objLog.Errorln("TopicLogic FindAll error:", err)
		return nil
	}

	uids := make([]int, 0, count)
	nids := make([]int, 0, count)

	for _, topic := range topicInfos {
		uids = append(uids, topic.Uid)
		if topic.Lastreplyuid != 0 {
			uids = append(uids, topic.Lastreplyuid)
		}
		nids = append(nids, topic.Nid)
	}
	usersMap := DefaultUser.FindUserInfos(ctx, uids)

	data := make([]map[string]interface{}, len(topicInfos))

	for i, topic := range topicInfos {
		dest := make(map[string]interface{})

		// 有人回复
		if topic.Lastreplyuid != 0 {
			if user, ok := usersMap[topic.Lastreplyuid]; ok {
				dest["lastreplyusername"] = user.Username
			}
		}

		util.Struct2Map(dest, topic.Topic)
		util.Struct2Map(dest, topic.TopicEx)

		dest["user"] = usersMap[topic.Uid]
		// tmpMap["node"] = nodes[tmpMap["nid"].(int)]

		data[i] = dest
	}

	return data
}

func (TopicLogic) Count(ctx context.Context, querystring string, args ...interface{}) int64 {
	objLog := GetLogger(ctx)

	var (
		total int64
		err   error
	)
	if querystring == "" {
		total, err = MasterDB.Count(new(model.Topic))
	} else {
		total, err = MasterDB.Where(querystring, args...).Count(new(model.Topic))
	}

	if err != nil {
		objLog.Errorln("TopicLogic Count error:", err)
	}

	return total
}
