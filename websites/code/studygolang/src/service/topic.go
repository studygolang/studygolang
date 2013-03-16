// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"model"
	"strconv"
	"util"
)

// 发布帖子。入topics和topics_ex库
func PublishTopic(topic *model.Topic) (errMsg string, err error) {
	tid, err := topic.Insert()
	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln(errMsg, "：", err)
		return
	}

	// 存扩展信息
	topicEx := model.NewTopicEx()
	topicEx.Tid = tid
	_, err = topicEx.Insert()
	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln(errMsg, "：", err)
		return
	}
	return
}

// 获得帖子详细信息（包括详细回复）
// 为了避免转换，tid传string类型
func FindTopicByTid(tid string) (topicMap map[string]interface{}, replies []map[string]interface{}, err error) {
	condition := "tid=" + tid
	// 帖子信息
	topic := model.NewTopic()
	err = topic.Where(condition).Find()
	if err != nil {
		logger.Errorln("topic service FindTopicByTid Error:", err)
		return
	}
	// 帖子不存在
	if topic.Tid == 0 {
		return
	}
	topicMap = make(map[string]interface{})
	util.Struct2Map(topicMap, topic)
	topicEx := model.NewTopicEx()
	err = topicEx.Where(condition).Find()
	if err != nil {
		logger.Errorln("topic service FindTopicByTid Error:", err)
		return
	}
	if topicEx.Tid == 0 {
		return
	}
	util.Struct2Map(topicMap, topicEx)
	// 节点名字
	topicMap["node"] = model.GetNodeName(topic.Nid)

	// 回复信息（评论）
	replyList, err := model.NewComment().Where("objid=" + tid + " and objtype=" + strconv.Itoa(model.TYPE_TOPIC)).FindAll()
	if err != nil {
		logger.Errorln("topic service FindTopicByTid Error:", err)
		return
	}

	replyNum := len(replyList)
	uids := make(map[int]int, replyNum+1)
	uids[topic.Uid] = topic.Uid
	for _, reply := range replyList {
		uids[reply.Uid] = reply.Uid
	}

	// 获得用户信息
	userMap := getUserInfos(uids)
	topicMap["user"] = userMap[topic.Uid]

	// 有人回复
	if topic.Lastreplyuid != 0 {
		topicMap["lastreplyusername"] = userMap[topicMap["lastreplyuid"].(int)].Username
	}

	replies = make([]map[string]interface{}, 0, replyNum)
	for _, reply := range replyList {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, reply)
		tmpMap["user"] = userMap[reply.Uid]
		replies = append(replies, tmpMap)
	}
	return
}

// 获得帖子列表页需要的数据
// 如果order为空，则默认排序方式（之所以用不定参数，是为了可以不传）
func FindTopics(page, pageNum int, where string, orderSlice ...string) (topics []map[string]interface{}, total int) {
	if pageNum == 0 {
		pageNum = PAGE_NUM
	}
	var offset = 0
	if page > 1 {
		offset = (page - 1) * pageNum
	}
	// 即使传了多个，也只取第一个
	order := ""
	if len(orderSlice) > 0 {
		order = orderSlice[0]
	}
	return FindTopicsByWhere(where, order, strconv.Itoa(offset)+","+strconv.Itoa(pageNum))
}

// 获得某个节点下的帖子列表（侧边栏推荐）
func FindTopicsByNid(nid, curTid string) (topics []*model.Topic) {
	var err error
	topics, err = model.NewTopic().Where("nid=" + nid + " and tid!=" + curTid).Limit("0,10").FindAll()
	if err != nil {
		logger.Errorln("topic service FindTopicsByNid Error:", err)
		return
	}
	return
}

func FindTopicsByWhere(where, order, limit string) (topics []map[string]interface{}, total int) {
	topicObj := model.NewTopic()
	if where != "" {
		topicObj.Where(where)
	}
	if order != "" {
		topicObj.Order(order)
	}
	if limit != "" {
		topicObj.Limit(limit)
	}
	topicList, err := topicObj.FindAll()
	if err != nil {
		logger.Errorln("topic service FindTopics Error:", err)
		return
	}
	// 获得总帖子数
	total, err = topicObj.Count()
	if err != nil {
		logger.Errorln("topic service FindTopics Error:", err)
		return
	}
	count := len(topicList)
	tids := make([]int, count)
	uids := make(map[int]int)
	nids := make([]int, count)
	for i, topic := range topicList {
		tids[i] = topic.Tid
		uids[topic.Uid] = topic.Uid
		if topic.Lastreplyuid != 0 {
			uids[topic.Lastreplyuid] = topic.Lastreplyuid
		}
		nids[i] = topic.Nid
	}

	// 获取扩展信息（计数）
	topicExList, err := model.NewTopicEx().Where("tid in(" + util.Join(tids, ",") + ")").FindAll()
	if err != nil {
		logger.Errorln("topic service FindTopics Error:", err)
		return
	}
	topicExMap := make(map[int]*model.TopicEx, len(topicExList))
	for _, topicEx := range topicExList {
		topicExMap[topicEx.Tid] = topicEx
	}

	userMap := getUserInfos(uids)

	// 获取节点信息
	nodes := model.GetNodesName(nids)

	topics = make([]map[string]interface{}, count)
	for i, topic := range topicList {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, topic)
		util.Struct2Map(tmpMap, topicExMap[topic.Tid])
		tmpMap["user"] = userMap[topic.Uid]
		// 有人回复
		if tmpMap["lastreplyuid"].(int) != 0 {
			tmpMap["lastreplyusername"] = userMap[tmpMap["lastreplyuid"].(int)].Username
		}
		tmpMap["node"] = nodes[tmpMap["nid"].(int)]
		topics[i] = tmpMap
	}
	return
}

// 获得某个用户最近的帖子
func FindRecentTopics(uid int) []*model.Topic {
	topics, err := model.NewTopic().Where("uid=" + strconv.Itoa(uid)).Order("ctime DESC").Limit("0, 5").FindAll()
	if err != nil {
		logger.Errorln("topic service FindRecentTopics error:", err)
		return nil
	}
	for _, topic := range topics {
		topic.Node = model.GetNodeName(topic.Nid)
	}
	return topics
}

// 获得回复最多的10条帖子(TODO:避免一直显示相同的)
func FindHotTopics() []map[string]interface{} {
	topicExList, err := model.NewTopicEx().Order("reply DESC").Limit("0,10").FindAll()
	if err != nil {
		logger.Errorln("topic service FindHotReplies error:", err)
		return nil
	}
	tidMap := make(map[int]int, len(topicExList))
	topicExMap := make(map[int]*model.TopicEx, len(topicExList))
	for _, topicEx := range topicExList {
		tidMap[topicEx.Tid] = topicEx.Tid
		topicExMap[topicEx.Tid] = topicEx
	}
	tids := util.MapIntKeys(tidMap)
	topics := FindTopicsByTids(tids)
	if topics == nil {
		return nil
	}

	uidMap := make(map[int]int, len(topics))
	for _, topic := range topics {
		uidMap[topic.Uid] = topic.Uid
	}
	userMap := getUserInfos(uidMap)
	result := make([]map[string]interface{}, len(topics))
	for i, topic := range topics {
		oneTopic := make(map[string]interface{})
		util.Struct2Map(oneTopic, topic)
		util.Struct2Map(oneTopic, topicExMap[topic.Tid])
		oneTopic["user"] = userMap[topic.Uid]
		result[i] = oneTopic
	}
	return result
}

// 获得回复对应的主贴信息
func FindRecentReplies(comments []*model.Comment) []map[string]interface{} {
	if len(comments) == 0 {
		return nil
	}
	count := len(comments)
	commentMap := make(map[int]*model.Comment, count)
	tidMap := make(map[int]int, count)
	for _, comment := range comments {
		commentMap[comment.Objid] = comment
		tidMap[comment.Objid] = comment.Objid
	}
	tids := util.MapIntKeys(tidMap)
	topics := FindTopicsByTids(tids)
	if len(topics) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(topics))
	for i, topic := range topics {
		oneReply := make(map[string]interface{})
		oneReply["tid"] = topic.Tid
		oneReply["title"] = topic.Title
		oneReply["cmt_content"] = commentMap[topic.Tid].Content
		oneReply["replytime"] = commentMap[topic.Tid].Ctime
		result[i] = oneReply
	}
	return result
}

// 获取多个帖子详细信息
func FindTopicsByTids(tids []int) []*model.Topic {
	inTids := util.Join(tids, ",")
	topics, err := model.NewTopic().Where("tid in(" + inTids + ")").FindAll()
	if err != nil {
		logger.Errorln("topic service FindRecentReplies error:", err)
		return nil
	}
	return topics
}

// 增加话题浏览数（TODO:刷屏暂时不处理）
func IncrTopicView(tid string, uid int) {
	model.NewTopicEx().Where("tid="+tid).Increment("view", 1)
}

// 话题总数
func TopicsTotal() (total int) {
	total, err := model.NewTopic().Count()
	if err != nil {
		logger.Errorln("topic service TopicsTotal error:", err)
		return
	}
	return
}

// 获取用户信息
func getUserInfos(uids map[int]int) map[int]*model.User {
	// 获取用户信息
	inUids := util.Join(util.MapIntKeys(uids), ",")
	users, err := model.NewUser().Where("uid in(" + inUids + ")").FindAll()
	if err != nil {
		logger.Errorln("topic service getUserInfos Error:", err)
		return map[int]*model.User{}
	}
	userMap := make(map[int]*model.User, len(users))
	for _, user := range users {
		userMap[user.Uid] = user
	}
	return userMap
}

// 话题回复（评论）
type TopicComment struct{}

// 更新该帖子的回复信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self TopicComment) UpdateComment(cid, objid, uid int, cmttime string) {
	tid := strconv.Itoa(objid)
	// 更新最后回复信息
	stringBuilder := util.NewBuffer().Append("lastreplyuid=").AppendInt(uid).Append(",lastreplytime=").Append(cmttime)
	err := model.NewTopic().Set(stringBuilder.String()).Where("tid=" + tid).Update()
	if err != nil {
		logger.Errorln("更新帖子最后回复人信息失败：", err)
	}
	// 更新回复数（TODO：暂时每次都更新表）
	err = model.NewTopicEx().Where("tid="+tid).Increment("reply", 1)
	if err != nil {
		logger.Errorln("更新帖子回复数失败：", err)
	}
}
