// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"html/template"
	"logger"
	"model"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"
)

// 发布帖子。入topics和topics_ex库
func PublishTopic(topic *model.Topic) (errMsg string, err error) {
	topic.Ctime = time.Now().Format("2006-01-02 15:04:05")
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

	// 发布帖子，活跃度+10
	go IncUserWeight("uid="+strconv.Itoa(topic.Uid), 10)

	return
}

// 修改帖子
// user 修改人的（有可能是作者或管理员）
func ModifyTopic(user map[string]interface{}, form url.Values) (errMsg string, err error) {
	uid := user["uid"].(int)
	form.Set("editor_uid", strconv.Itoa(uid))

	fields := []string{"title", "content", "nid", "editor_uid"}
	query, args := updateSetClause(form, fields)

	tid := form.Get("tid")

	err = model.NewTopic().Set(query, args...).Where("tid=" + tid).Update()
	if err != nil {
		logger.Errorf("更新帖子 【%s】 信息失败：%s\n", tid, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	username := user["username"].(string)
	// 修改帖子，活跃度+2
	go IncUserWeight("username="+username, 2)

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
		err = errors.New("The topic of tid is not exists")
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
	topicMap["node"] = GetNodeName(topic.Nid)

	// 回复信息（评论）
	replies, owerUser, lastReplyUser := FindObjComments(tid, strconv.Itoa(model.TYPE_TOPIC), topic.Uid, topic.Lastreplyuid)
	topicMap["user"] = owerUser
	// 有人回复
	if topic.Lastreplyuid != 0 {
		topicMap["lastreplyusername"] = lastReplyUser.Username
	}

	if topic.EditorUid != 0 {
		topicMap["editor_username"] = FindUsernameByUid(topic.EditorUid)
	}

	return
}

// 通过tid获得话题的所有者
func getTopicOwner(tid int) int {
	// 帖子信息
	topic := model.NewTopic()
	err := topic.Where("tid=" + strconv.Itoa(tid)).Find()
	if err != nil {
		logger.Errorln("topic service getTopicOwner Error:", err)
		return 0
	}
	return topic.Uid
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
	order := "mtime DESC"
	if len(orderSlice) > 0 && orderSlice[0] != "" {
		order = orderSlice[0]
	}
	return FindTopicsByWhere(where, order, strconv.Itoa(offset)+","+strconv.Itoa(pageNum))
}

// 获取话题列表（分页），目前供后台使用
func FindTopicsByPage(conds map[string]string, curPage, limit int) ([]*model.Topic, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	topic := model.NewTopic()

	limitStr := strconv.Itoa((curPage-1)*limit) + "," + strconv.Itoa(limit)
	topicList, err := topic.Where(strings.Join(conditions, " AND ")).Order("tid DESC").Limit(limitStr).
		FindAll()
	if err != nil {
		logger.Errorln("topic service FindTopicsByPage Error:", err)
		return nil, 0
	}

	total, err := topic.Count()
	if err != nil {
		logger.Errorln("topic service FindTopicsByPage COUNT Error:", err)
		return nil, 0
	}

	return topicList, total
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

// 获得社区最新公告（废弃）
func FindNoticeTopic() (topic *model.Topic) {
	topics, err := model.NewTopic().Where("nid=15").Limit("0,1").Order("mtime DESC").FindAll()
	if err != nil {
		logger.Errorln("topic service FindNoticeTopic Error:", err)
		return
	}
	if len(topics) > 0 {
		topic = topics[0]
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
		logger.Errorln("topic service topicObj.FindAll Error:", err)
		return
	}
	// 获得总帖子数
	total, err = topicObj.Count()
	if err != nil {
		logger.Errorln("topic service topicObj.Count Error:", err)
		return
	}
	count := len(topicList)
	tids := make([]int, count)
	uids := make([]int, 0, count)
	nids := make([]int, count)
	for i, topic := range topicList {
		tids[i] = topic.Tid
		uids = append(uids, topic.Uid)
		if topic.Lastreplyuid != 0 {
			uids = append(uids, topic.Lastreplyuid)
		}
		nids[i] = topic.Nid
	}

	// 获取扩展信息（计数）
	topicExList, err := model.NewTopicEx().Where("tid in(" + util.Join(tids, ",") + ")").FindAll()
	if err != nil {
		logger.Errorln("topic service NewTopicEx FindAll Error:", err)
		return
	}
	topicExMap := make(map[int]*model.TopicEx, len(topicExList))
	for _, topicEx := range topicExList {
		topicExMap[topicEx.Tid] = topicEx
	}

	userMap := GetUserInfos(uids)

	// 获取节点信息
	nodes := GetNodesName(nids)

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

// 获得最近的帖子(如果uid!=0，则获取某个用户最近的帖子)
func FindRecentTopics(uid int, limit string) []*model.Topic {
	cond := ""
	if uid != 0 {
		cond = "uid=" + strconv.Itoa(uid)
	}

	topics, err := model.NewTopic().Where(cond).Order("ctime DESC").Limit(limit).FindAll()
	if err != nil {
		logger.Errorln("topic service FindRecentTopics error:", err)
		return nil
	}
	for _, topic := range topics {
		topic.Node = GetNodeName(topic.Nid)
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

	uids := util.Models2Intslice(topics, "Uid")
	userMap := GetUserInfos(uids)

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

// 获取多个帖子详细信息
func FindTopicsByTids(tids []int) []*model.Topic {
	if len(tids) == 0 {
		return nil
	}
	inTids := util.Join(tids, ",")
	topics, err := model.NewTopic().Where("tid in(" + inTids + ")").FindAll()
	if err != nil {
		logger.Errorln("topic service FindTopicsByTids error:", err)
		return nil
	}
	return topics
}

// 提供给其他service调用（包内）
func getTopics(tids map[int]int) map[int]*model.Topic {
	topics := FindTopicsByTids(util.MapIntKeys(tids))
	topicMap := make(map[int]*model.Topic, len(topics))
	for _, topic := range topics {
		topicMap[topic.Tid] = topic
	}
	return topicMap
}

// 获得热门节点
func FindHotNodes() []map[string]interface{} {
	strSql := "SELECT nid, COUNT(1) AS topicnum FROM topics GROUP BY nid ORDER BY topicnum DESC LIMIT 10"
	rows, err := model.NewTopic().DoSql(strSql)
	if err != nil {
		logger.Errorln("topic service FindHotNodes error:", err)
		return nil
	}
	nodes := make([]map[string]interface{}, 0, 10)
	for rows.Next() {
		var nid, topicnum int
		err = rows.Scan(&nid, &topicnum)
		if err != nil {
			logger.Errorln("rows.Scan error:", err)
			continue
		}
		name := GetNodeName(nid)
		node := map[string]interface{}{
			"name": name,
			"nid":  nid,
		}
		nodes = append(nodes, node)
	}
	return nodes
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
	}
	return
}

// 安全过滤
func JSEscape(topics []*model.Topic) []*model.Topic {
	for i, topic := range topics {
		topics[i].Title = template.JSEscapeString(topic.Title)
		topics[i].Content = template.JSEscapeString(topic.Content)
	}
	return topics
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

// 实现 CommentObjecter 接口
func (self TopicComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	topics := FindTopicsByTids(ids)
	if len(topics) == 0 {
		return
	}

	for _, topic := range topics {
		objinfo := make(map[string]interface{})
		objinfo["title"] = topic.Title

		for _, comment := range commentMap[topic.Tid] {
			comment.Objinfo = objinfo
		}
	}
}
