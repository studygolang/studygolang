// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"html/template"
	"model"
	"time"

	. "db"

	"github.com/fatih/set"
	"github.com/fatih/structs"
	"github.com/polaris1119/logger"
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

	uidSet := set.New()
	nidSet := set.New()
	for _, topicInfo := range topicInfos {
		uidSet.Add(topicInfo.Uid)
		if topicInfo.Lastreplyuid != 0 {
			uidSet.Add(topicInfo.Lastreplyuid)
		}
		nidSet.Add(topicInfo.Nid)
	}

	usersMap := DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))
	// 获取节点信息
	nodes := GetNodesName(set.IntSlice(nidSet))

	data := make([]map[string]interface{}, len(topicInfos))

	for i, topicInfo := range topicInfos {
		dest := make(map[string]interface{})

		// 有人回复
		if topicInfo.Lastreplyuid != 0 {
			if user, ok := usersMap[topicInfo.Lastreplyuid]; ok {
				dest["lastreplyusername"] = user.Username
			}
		}

		structs.FillMap(topicInfo.Topic, dest)
		structs.FillMap(topicInfo.TopicEx, dest)

		dest["user"] = usersMap[topicInfo.Uid]
		dest["node"] = nodes[topicInfo.Nid]

		data[i] = dest
	}

	return data
}

func (TopicLogic) FindLastList(beginTime string, limit int) ([]*model.Topic, error) {
	topics := make([]*model.Topic, 0)
	err := MasterDB.Where("ctime>? AND flag IN(?,?)", beginTime, model.FlagNoAudit, model.FlagNormal).
		OrderBy("tid DESC").Limit(limit).Find(&topics)

	return topics, err
}

// FindRecent 获得最近的主题(uids[0]，则获取某个用户最近的主题)
func (TopicLogic) FindRecent(limit int, uids ...int) []*model.Topic {
	dbSession := MasterDB.OrderBy("ctime DESC").Limit(limit)
	if len(uids) > 0 {
		dbSession.Where("uid=?", uids[0])
	}

	topics := make([]*model.Topic, 0)
	if err := dbSession.Find(&topics); err != nil {
		logger.Errorln("TopicLogic FindRecent error:", err)
	}
	// for _, topic := range topics {
	// 	topic.Node = GetNodeName(topic.Nid)
	// }
	return topics
}

// FindByNid 获得某个节点下的主题列表（侧边栏推荐）
func (TopicLogic) FindByNid(ctx context.Context, nid, curTid string) []*model.Topic {
	objLog := GetLogger(ctx)

	topics := make([]*model.Topic, 0)
	err := MasterDB.Where("nid=? AND tid!=?", nid, curTid).Limit(10).Find(&topics)
	if err != nil {
		objLog.Errorln("TopicLogic FindByNid Error:", err)
	}

	return topics
}

// FindByTids 获取多个主题详细信息
func (TopicLogic) FindByTids(tids []int) []*model.Topic {
	if len(tids) == 0 {
		return nil
	}

	topics := make([]*model.Topic, 0)
	err := MasterDB.In("tid", tids).Find(&topics)
	if err != nil {
		logger.Errorln("TopicLogic FindByTids error:", err)
		return nil
	}
	return topics
}

// FindHotNodes 获得热门节点
func (TopicLogic) FindHotNodes(ctx context.Context) []map[string]interface{} {
	objLog := GetLogger(ctx)

	strSql := "SELECT nid, COUNT(1) AS topicnum FROM topics GROUP BY nid ORDER BY topicnum DESC LIMIT 10"
	rows, err := MasterDB.DB().DB.Query(strSql)
	if err != nil {
		objLog.Errorln("TopicLogic FindHotNodes error:", err)
		return nil
	}
	nodes := make([]map[string]interface{}, 0, 10)
	for rows.Next() {
		var nid, topicnum int
		err = rows.Scan(&nid, &topicnum)
		if err != nil {
			objLog.Errorln("rows.Scan error:", err)
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

// Total 话题总数
func (TopicLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.Topic))
	if err != nil {
		logger.Errorln("TopicLogic Total error:", err)
	}
	return total
}

// JSEscape 安全过滤
func (TopicLogic) JSEscape(topics []*model.Topic) []*model.Topic {
	for i, topic := range topics {
		topics[i].Title = template.JSEscapeString(topic.Title)
		topics[i].Content = template.JSEscapeString(topic.Content)
	}
	return topics
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

// 话题回复（评论）
type TopicComment struct{}

// UpdateComment 更新该主题的回复信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self TopicComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新最后回复信息
	_, err := MasterDB.Table(new(model.Topic)).Id(objid).Update(map[string]interface{}{
		"lastreplyuid":  uid,
		"lastreplytime": cmttime,
	})
	if err != nil {
		logger.Errorln("更新主题最后回复人信息失败：", err)
	}

	// 更新回复数（TODO：暂时每次都更新表）
	_, err = MasterDB.Id(objid).Incr("reply", 1).Update(new(model.TopicEx))
	if err != nil {
		logger.Errorln("更新主题回复数失败：", err)
	}
}

func (self TopicComment) String() string {
	return "topic"
}

// 实现 CommentObjecter 接口
func (self TopicComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {

	topics := DefaultTopic.FindByTids(ids)
	if len(topics) == 0 {
		return
	}

	for _, topic := range topics {
		objinfo := make(map[string]interface{})
		objinfo["title"] = topic.Title
		objinfo["uri"] = model.PathUrlMap[model.TypeTopic]
		objinfo["type_name"] = model.TypeNameMap[model.TypeTopic]

		for _, comment := range commentMap[topic.Tid] {
			comment.Objinfo = objinfo
		}
	}
}
