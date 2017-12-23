// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"fmt"
	"html/template"
	"model"
	"net/url"
	"sync"
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
func (self TopicLogic) Publish(ctx context.Context, me *model.Me, form url.Values) (tid int, err error) {
	objLog := GetLogger(ctx)

	tid = goutils.MustInt(form.Get("tid"))
	if tid != 0 {
		topic := &model.Topic{}
		_, err = MasterDB.Id(tid).Get(topic)
		if err != nil {
			objLog.Errorln("Publish Topic find error:", err)
			return
		}

		if !CanEdit(me, topic) {
			err = NotModifyAuthorityErr
			return
		}

		_, err = self.Modify(ctx, me, form)
		if err != nil {
			objLog.Errorln("Publish Topic modify error:", err)
			return
		}

		nid := goutils.MustInt(form.Get("nid"))

		go func() {
			// 不是作者自己修改，且是调整节点，扣除铜币
			if topic.Uid != me.Uid && topic.Nid != nid {
				node := DefaultNode.FindOne(nid)
				award := -500
				if node.ShowIndex {
					award = -30
				}
				desc := fmt.Sprintf(`主题节点被管理员调整为 <a href="/go/%s">%s</a>`, node.Ename, node.Name)
				user := DefaultUser.FindOne(ctx, "uid", topic.Uid)
				DefaultUserRich.IncrUserRich(user, model.MissionTypeModify, award, desc)
			}

			if nid != topic.Nid {
				DefaultFeed.modifyTopicNode(tid, nid)
			}
		}()
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
		session.Begin()

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
		session.Commit()

		go func() {
			// 同一个首页不显示的节点，一天发布主题数超过3个，扣 1 千铜币
			topicNum, err := MasterDB.Where("uid=? AND ctime>?", me.Uid, time.Now().Format("2006-01-02 00:00:00")).Count(new(model.Topic))
			if err != nil {
				logger.Errorln("find today topic num error:", err)
				return
			}

			if topicNum > 3 {
				node := DefaultNode.FindOne(topic.Nid)
				if node.ShowIndex {
					return
				}

				award := -1000

				desc := fmt.Sprintf(`一天发布推广过多或 Spam 扣除铜币 %d 个`, -award)
				user := DefaultUser.FindOne(ctx, "uid", me.Uid)
				DefaultUserRich.IncrUserRich(user, model.MissionTypeSpam, award, desc)

				DefaultRank.GenDAURank(me.Uid, -1000)
			}
		}()

		// 发布动态
		DefaultFeed.publish(topic, topicEx)

		// 给 被@用户 发系统消息
		ext := map[string]interface{}{
			"objid":   topic.Tid,
			"objtype": model.TypeTopic,
			"uid":     me.Uid,
			"msgtype": model.MsgtypePublishAtMe,
		}
		go DefaultMessage.SendSysMsgAtUsernames(ctx, usernames, ext, 0)

		go publishObservable.NotifyObservers(me.Uid, model.TypeTopic, topic.Tid)

		tid = topic.Tid
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

	fields := []string{"title", "content", "nid", "permission"}
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

	go modifyObservable.NotifyObservers(user.Uid, model.TypeTopic, goutils.MustInt(tid))

	return
}

// Append 主题附言
func (self TopicLogic) Append(ctx context.Context, uid, tid int, content string) error {
	objLog := GetLogger(ctx)

	// 当前已经附言了几条，最多 3 条
	num, err := MasterDB.Where("tid=?", tid).Count(new(model.TopicAppend))
	if err != nil {
		objLog.Errorln("TopicLogic Append error:", err)
		return err
	}

	if num >= model.AppendMaxNum {
		return errors.New("不允许再发附言！")
	}

	topicAppend := &model.TopicAppend{
		Tid:     tid,
		Content: content,
	}
	_, err = MasterDB.Insert(topicAppend)

	if err != nil {
		objLog.Errorln("TopicLogic Append insert error:", err)
		return err
	}

	go appendObservable.NotifyObservers(uid, model.TypeTopic, tid)

	return nil
}

// SetTop 置顶
func (self TopicLogic) SetTop(ctx context.Context, me *model.Me, tid int) error {
	objLog := GetLogger(ctx)

	if !me.IsAdmin {
		topic := self.findByTid(tid)
		if topic.Tid == 0 || topic.Uid != me.Uid {
			return NotFoundErr
		}
	}

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.Table(new(model.Topic)).Id(tid).Update(map[string]interface{}{
		"top":      1,
		"top_time": time.Now().Unix(),
	})
	if err != nil {
		objLog.Errorln("TopicLogic SetTop error:", err)
		session.Rollback()
		return err
	}

	err = DefaultFeed.setTop(session, tid, model.TypeTopic, 1)
	if err != nil {
		objLog.Errorln("TopicLogic SetTop feed error:", err)
		session.Rollback()
		return err
	}

	session.Commit()

	go topObservable.NotifyObservers(me.Uid, model.TypeTopic, tid)

	return nil
}

// UnsetTop 取消置顶
func (self TopicLogic) UnsetTop(ctx context.Context, tid int) error {
	objLog := GetLogger(ctx)

	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.Table(new(model.Topic)).Id(tid).Update(map[string]interface{}{
		"top": 0,
	})
	if err != nil {
		objLog.Errorln("TopicLogic UnsetTop error:", err)
		session.Rollback()
		return err
	}

	err = DefaultFeed.setTop(session, tid, model.TypeTopic, 0)
	if err != nil {
		objLog.Errorln("TopicLogic UnsetTop feed error:", err)
		session.Rollback()
		return err
	}

	session.Commit()

	return nil
}

// AutoUnsetTop 自动取消置顶
func (self TopicLogic) AutoUnsetTop() error {
	topics := make([]*model.Topic, 0)
	err := MasterDB.Where("top=1").Find(&topics)
	if err != nil {
		logger.Errorln("TopicLogic AutoUnsetTop error:", err)
		return err
	}

	for _, topic := range topics {
		if topic.TopTime == 0 || topic.TopTime+86400 > time.Now().Unix() {
			continue
		}

		self.UnsetTop(nil, topic.Tid)
	}

	return nil
}

// FindAll 支持多页翻看
func (self TopicLogic) FindAll(ctx context.Context, paginator *Paginator, orderBy string, querystring string, args ...interface{}) []map[string]interface{} {
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

	return self.fillDataForTopicInfo(topicInfos)
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

func (self TopicLogic) FindFullinfoByTids(tids []int) []map[string]interface{} {
	topicInfoMap := make(map[int]*model.TopicInfo, 0)

	err := MasterDB.Join("INNER", "topics_ex", "topics.tid=topics_ex.tid").In("topics.tid", tids).Find(&topicInfoMap)
	if err != nil {
		logger.Errorln("TopicLogic FindFullinfoByTids error:", err)
		return nil
	}

	topicInfos := make([]*model.TopicInfo, 0, len(topicInfoMap))
	for _, tid := range tids {
		if topicInfo, ok := topicInfoMap[tid]; ok {
			topicInfos = append(topicInfos, topicInfo)
		}
	}

	return self.fillDataForTopicInfo(topicInfos)
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

	if topic.Flag > model.FlagNormal {
		err = errors.New("The topic of tid is not exists or delete")
		return
	}

	topicMap = make(map[string]interface{})
	structs.FillMap(topic, topicMap)
	structs.FillMap(topicInfo.TopicEx, topicMap)

	// 解析内容中的 @
	topicMap["content"] = self.decodeTopicContent(ctx, topic)

	// 节点
	topicMap["node"] = GetNode(topic.Nid)

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

// 获取列表（分页）：后台用
func (TopicLogic) FindByPage(ctx context.Context, conds map[string]string, curPage, limit int) ([]*model.Topic, int) {
	objLog := GetLogger(ctx)

	session := MasterDB.NewSession()

	for k, v := range conds {
		session.And(k+"=?", v)
	}

	totalSession := session.Clone()

	offset := (curPage - 1) * limit
	topicList := make([]*model.Topic, 0)
	err := session.OrderBy("tid DESC").Limit(limit, offset).Find(&topicList)
	if err != nil {
		objLog.Errorln("find error:", err)
		return nil, 0
	}

	total, err := totalSession.Count(new(model.Topic))
	if err != nil {
		objLog.Errorln("find count error:", err)
		return nil, 0
	}

	return topicList, int(total)
}

func (TopicLogic) FindAppend(ctx context.Context, tid int) []*model.TopicAppend {
	objLog := GetLogger(ctx)

	topicAppends := make([]*model.TopicAppend, 0)
	err := MasterDB.Where("tid=?", tid).Find(&topicAppends)
	if err != nil {
		objLog.Errorln("TopicLogic FindAppend error:", err)
	}

	return topicAppends
}

func (TopicLogic) findByTid(tid int) *model.Topic {
	topic := &model.Topic{}
	_, err := MasterDB.Where("tid=?", tid).Get(topic)
	if err != nil {
		logger.Errorln("TopicLogic findByTid error:", err)
	}
	return topic
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

func (TopicLogic) fillDataForTopicInfo(topicInfos []*model.TopicInfo) []map[string]interface{} {
	uidSet := set.New(set.NonThreadSafe)
	nidSet := set.New(set.NonThreadSafe)
	for _, topicInfo := range topicInfos {
		uidSet.Add(topicInfo.Uid)
		if topicInfo.Lastreplyuid != 0 {
			uidSet.Add(topicInfo.Lastreplyuid)
		}
		nidSet.Add(topicInfo.Nid)
	}

	usersMap := DefaultUser.FindUserInfos(nil, set.IntSlice(uidSet))
	// 获取节点信息
	nodes := GetNodesByNids(set.IntSlice(nidSet))

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

var (
	hotNodesCache  []map[string]interface{}
	hotNodesBegin  time.Time
	hotNodesLocker sync.Mutex
)

// FindHotNodes 获得热门节点
func (TopicLogic) FindHotNodes(ctx context.Context) []map[string]interface{} {
	hotNodesLocker.Lock()
	defer hotNodesLocker.Unlock()
	if !hotNodesBegin.IsZero() && hotNodesBegin.Add(1*time.Hour).Before(time.Now()) {
		return hotNodesCache
	}

	objLog := GetLogger(ctx)

	hotNum := 10

	lastWeek := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	strSql := fmt.Sprintf("SELECT nid, COUNT(1) AS topicnum FROM topics WHERE ctime>='%s' GROUP BY nid ORDER BY topicnum DESC LIMIT 15", lastWeek)
	rows, err := MasterDB.DB().DB.Query(strSql)
	if err != nil {
		objLog.Errorln("TopicLogic FindHotNodes error:", err)
		return nil
	}

	nids := make([]int, 0, 15)
	for rows.Next() {
		var nid, topicnum int
		err = rows.Scan(&nid, &topicnum)
		if err != nil {
			objLog.Errorln("rows.Scan error:", err)
			continue
		}

		nids = append(nids, nid)
	}

	nodes := make([]map[string]interface{}, 0, hotNum)

	topicNodes := GetNodesByNids(nids)
	for _, nid := range nids {
		topicNode := topicNodes[nid]
		if !topicNode.ShowIndex {
			continue
		}

		node := map[string]interface{}{
			"name":  topicNode.Name,
			"ename": topicNode.Ename,
			"nid":   topicNode.Nid,
		}
		nodes = append(nodes, node)
		if len(nodes) == hotNum {
			break
		}
	}

	hotNodesCache = nodes
	hotNodesBegin = time.Now()

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
	// 允许内嵌 Wide iframe
	content := util.EmbedWide(topic.Content)

	// @别人
	return parseAtUser(ctx, content)
}

// 话题回复（评论）
type TopicComment struct{}

// UpdateComment 更新该主题的回复信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self TopicComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	session := MasterDB.NewSession()
	defer session.Close()

	session.Begin()

	// 更新最后回复信息
	_, err := session.Table(new(model.Topic)).Id(objid).Update(map[string]interface{}{
		"lastreplyuid":  uid,
		"lastreplytime": cmttime,
	})
	if err != nil {
		logger.Errorln("更新主题最后回复人信息失败：", err)
		session.Rollback()
		return
	}

	// 更新回复数（TODO：暂时每次都更新表）
	_, err = MasterDB.Id(objid).Incr("reply", 1).Update(new(model.TopicUpEx))
	if err != nil {
		logger.Errorln("更新主题回复数失败：", err)
		session.Rollback()
		return
	}

	session.Commit()
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
	_, err := MasterDB.Where("tid=?", objid).Incr("like", num).Update(new(model.TopicUpEx))
	if err != nil {
		logger.Errorln("更新主题喜欢数失败：", err)
	}
}

func (self TopicLike) String() string {
	return "topic"
}
