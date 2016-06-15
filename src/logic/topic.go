// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"html/template"
	"model"
	"net/url"
	"time"
	"util"

	. "db"

	"github.com/fatih/structs"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/set"
	"golang.org/x/net/context"
)

type TopicLogic struct{}

var DefaultTopic = TopicLogic{}

// Publish 发布主题。入topics和topics_ex库
func (self TopicLogic) Publish(ctx context.Context, me *model.Me, form url.Values) (err error) {
	objLog := GetLogger(ctx)

	tid := goutils.MustInt(form.Get("tid"))
	if tid != 0 {
		topic := &model.Topic{}
		_, err = MasterDB.Id(tid).Get(topic)
		if err != nil {
			objLog.Errorln("Publish Topic find error:", err)
			return
		}

		if topic.Uid != me.Uid && !me.IsAdmin {
			err = NotModifyAuthorityErr
			return
		}

		_, err = self.Modify(ctx, me, form)
		if err != nil {
			objLog.Errorln("Publish Topic modif error:", err)
			return
		}
	} else {
		usernames := form.Get("usernames")
		form.Del("usernames")

		topic := &model.Topic{}
		err = schemaDecoder.Decode(topic, form)
		if err != nil {
			objLog.Errorln("TopicLogic Publish decode error:", err)
			return
		}
		topic.Uid = me.Uid
		topic.Lastreplytime = model.NewOftenTime()

		session := MasterDB.NewSession()
		defer session.Close()

		_, err = session.Insert(topic)
		if err != nil {
			session.Rollback()
			objLog.Errorln("TopicLogic Publish insert error:", err)
			return
		}

		topicEx := &model.TopicEx{
			Tid: topic.Tid,
		}

		_, err = session.Insert(topicEx)
		if err != nil {
			session.Rollback()
			objLog.Errorln("TopicLogic Publish Insert TopicEx error:", err)
			return
		}

		// 给 被@用户 发系统消息
		ext := map[string]interface{}{
			"objid":   topic.Tid,
			"objtype": model.TypeTopic,
			"uid":     me.Uid,
			"msgtype": model.MsgtypePublishAtMe,
		}
		go DefaultMessage.SendSysMsgAtUsernames(ctx, usernames, ext, 0)

		// 发布主题，活跃度+10
		go DefaultUser.IncrUserWeight("uid", me.Uid, 10)
	}

	return
}

// Modify 修改主题
// user 修改人的（有可能是作者或管理员）
func (TopicLogic) Modify(ctx context.Context, user *model.Me, form url.Values) (errMsg string, err error) {
	objLog := GetLogger(ctx)

	change := map[string]interface{}{
		"editor_uid": user.Uid,
	}

	fields := []string{"title", "content", "nid"}
	for _, field := range fields {
		change[field] = form.Get(field)
	}

	tid := form.Get("tid")
	_, err = MasterDB.Table(new(model.Topic)).Id(tid).Update(change)
	if err != nil {
		objLog.Errorf("更新主题 【%s】 信息失败：%s\n", tid, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	// 修改主题，活跃度+2
	go DefaultUser.IncrUserWeight("uid", user.Uid, 2)

	return
}

// FindAll 支持多页翻看
func (TopicLogic) FindAll(ctx context.Context, paginator *Paginator, orderBy string, querystring string, args ...interface{}) []map[string]interface{} {
	objLog := GetLogger(ctx)

	topicInfos := make([]*model.TopicInfo, 0)

	session := MasterDB.Join("INNER", "topics_ex", "topics.tid=topics_ex.tid")
	if querystring != "" {
		session.Where(querystring, args...)
	}
	err := session.OrderBy(orderBy).Limit(paginator.PerPage(), paginator.Offset()).Find(&topicInfos)
	if err != nil {
		objLog.Errorln("TopicLogic FindAll error:", err)
		return nil
	}

	uidSet := set.New(set.NonThreadSafe)
	nidSet := set.New(set.NonThreadSafe)
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
	for _, topic := range topics {
		topic.Node = GetNodeName(topic.Nid)
	}
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

// findByTids 获取多个主题详细信息 包内用
func (TopicLogic) findByTids(tids []int) map[int]*model.Topic {
	if len(tids) == 0 {
		return nil
	}

	topics := make(map[int]*model.Topic)
	err := MasterDB.In("tid", tids).Find(&topics)
	if err != nil {
		logger.Errorln("TopicLogic findByTids error:", err)
		return nil
	}
	return topics
}

// FindByTid 获得主题详细信息（包括详细回复）
func (self TopicLogic) FindByTid(ctx context.Context, tid int) (topicMap map[string]interface{}, replies []map[string]interface{}, err error) {
	objLog := GetLogger(ctx)

	topicInfo := &model.TopicInfo{}
	_, err = MasterDB.Join("INNER", "topics_ex", "topics.tid=topics_ex.tid").Where("topics.tid=?", tid).Get(topicInfo)
	if err != nil {
		objLog.Errorln("TopicLogic FindByTid get error:", err)
		return
	}

	topic := &topicInfo.Topic

	if topic.Tid == 0 {
		err = errors.New("The topic of tid is not exists")
		objLog.Errorln("TopicLogic FindByTid get error:", err)
		return
	}

	topicMap = make(map[string]interface{})
	structs.FillMap(topic, topicMap)
	structs.FillMap(topicInfo.TopicEx, topicMap)

	// 解析内容中的 @
	topicMap["content"] = self.decodeTopicContent(ctx, topic)

	// 节点名字
	topicMap["node"] = GetNodeName(topic.Nid)

	// 回复信息（评论）
	replies, owerUser, lastReplyUser := DefaultComment.FindObjComments(ctx, topic.Tid, model.TypeTopic, topic.Uid, topic.Lastreplyuid)
	topicMap["user"] = owerUser
	// 有人回复
	if topic.Lastreplyuid != 0 {
		topicMap["lastreplyusername"] = lastReplyUser.Username
	}

	if topic.EditorUid != 0 {
		editorUser := DefaultUser.FindOne(ctx, "uid", topic.EditorUid)
		topicMap["editor_username"] = editorUser.Username
	}

	return
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

// getOwner 通过tid获得话题的所有者
func (TopicLogic) getOwner(tid int) int {
	topic := &model.Topic{}
	_, err := MasterDB.Id(tid).Get(topic)
	if err != nil {
		logger.Errorln("topic logic getOwner Error:", err)
		return 0
	}
	return topic.Uid
}

func (TopicLogic) decodeTopicContent(ctx context.Context, topic *model.Topic) string {
	// 安全过滤
	content := template.HTMLEscapeString(topic.Content)

	// 允许内嵌 Wide iframe
	content = util.EmbedWide(content)

	// @别人
	return parseAtUser(ctx, content)
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

// 主题喜欢
type TopicLike struct{}

// 更新该主题的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self TopicLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Where("tid=?", objid).Incr("like", num).Update(new(model.TopicEx))
	if err != nil {
		logger.Errorln("更新主题喜欢数失败：", err)
	}
}

func (self TopicLike) String() string {
	return "topic"
}
