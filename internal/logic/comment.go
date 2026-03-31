// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"regexp"
	"strings"
	"time"

	. "github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/internal/model"

	"github.com/fatih/structs"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
	"github.com/polaris1119/set"
	"github.com/polaris1119/slices"
	"golang.org/x/sync/errgroup"
)

type CommentLogic struct{}

const (
	// Redis key 前缀用于评论楼层计数器
	commentFloorKeyPrefix = "comment:floor:"

	// Lua 脚本：原子性获取并递增楼层号
	getNextFloorScript = `
local key = KEYS[1]
local ttl = ARGV[1]

-- 尝试递增
local floor = redis.call('INCR', key)

-- 如果是第一次创建，设置过期时间
if floor == 1 then
    redis.call('EXPIRE', key, ttl)
end

return floor
`
)

// GetNextCommentFloor 获取下一个评论楼层号（原子性保证）
func GetNextCommentFloor(objid, objtype int) (int, error) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	key := fmt.Sprintf("%s%d:%d", commentFloorKeyPrefix, objtype, objid)
	ttl := 7 * 24 * 3600 // 7 天

	// 使用带前缀的 key
	prefixedKey := nosql.KeyPrefix + key

	// 执行 Lua 脚本（原子性）
	result, err := redisClient.Do("EVAL", getNextFloorScript, 1, prefixedKey, ttl)
	if err == nil {
		if floor, ok := result.(int64); ok {
			return int(floor), nil
		}
	}

	logger.Errorln("GetNextCommentFloor: Redis Lua script execution failed:", err)

	// 回退策略：使用数据库 + 分布式锁
	return getNextFloorWithLock(objid, objtype, key, ttl)
}

// getNextFloorWithLock 使用分布式锁从数据库获取楼层号
func getNextFloorWithLock(objid, objtype int, key string, ttl int) (int, error) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	prefixedKey := nosql.KeyPrefix + key
	lockKey := prefixedKey + ":lock"
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())

	// 尝试获取锁（最多等待 3 秒）
	locked := false
	for i := 0; i < 30; i++ {
		// SET NX EX（原子性获取锁）
		reply, err := redisClient.Do("SET", lockKey, lockValue, "NX", "EX", 10)
		if err == nil && reply != nil {
			locked = true
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	if !locked {
		return 0, errors.New("failed to acquire distributed lock: timeout")
	}

	defer redisClient.Do("DEL", lockKey) // 释放锁

	// 从数据库查询最大楼层号
	tmpCmt := &model.Comment{}
	_, err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).
		OrderBy("floor DESC").
		Get(tmpCmt)

	if err != nil {
		logger.Errorln("GetNextCommentFloor: database query failed:", err)
		return 1, nil // 默认从 1 开始
	}

	nextFloor := tmpCmt.Floor + 1

	// 初始化 Redis 计数器
	if err := redisClient.SET(key, nextFloor, ttl); err != nil {
		logger.Errorln("GetNextCommentFloor: failed to initialize Redis counter:", err)
	}

	return nextFloor, nil
}

var DefaultComment = CommentLogic{}

// FindObjComments 获得某个对象的所有评论
// owner: 被评论对象属主
// TODO:分页暂不做
func (self CommentLogic) FindObjComments(ctx context.Context, objid, objtype int, owner, lastCommentUid int /*, page, pageNum int*/) (comments []map[string]interface{}, ownerUser, lastReplyUser *model.User) {
	objLog := GetLogger(ctx)

	commentList := make([]*model.Comment, 0)
	err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).Find(&commentList)
	if err != nil {
		objLog.Errorln("CommentLogic FindObjComments Error:", err)
		return
	}

	uids := slices.StructsIntSlice(commentList, "Uid")

	// 避免某些情况下最后回复人没在回复列表中
	uids = append(uids, owner, lastCommentUid)

	// 获得用户信息
	userMap := DefaultUser.FindUserInfos(ctx, uids)
	ownerUser = userMap[owner]
	if lastCommentUid != 0 {
		lastReplyUser = userMap[lastCommentUid]
	}
	comments = make([]map[string]interface{}, 0, len(commentList))
	for _, comment := range commentList {
		tmpMap := structs.Map(comment)
		tmpMap["content"] = template.HTML(self.decodeCmtContent(ctx, comment))
		tmpMap["user"] = userMap[comment.Uid]
		comments = append(comments, tmpMap)
	}
	return
}

const CommentPerNum = 50

// FindObjectComments 获得某个对象的所有评论（新版）
func (self CommentLogic) FindObjectComments(ctx context.Context, objid, objtype, p int) (commentList []*model.Comment, replyComments []*model.Comment, pageNum int, err error) {
	objLog := GetLogger(ctx)

	total, err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).Count(new(model.Comment))
	if err != nil {
		objLog.Errorln("comment logic FindObjectComments count Error:", err)
		return
	}

	pageNum = int(math.Ceil(float64(total) / CommentPerNum))
	if p == 0 {
		p = pageNum
	}

	commentList = make([]*model.Comment, 0)
	err = MasterDB.Where("objid=? AND objtype=?", objid, objtype).Asc("cid").
		Limit(CommentPerNum, (p-1)*CommentPerNum).
		Find(&commentList)
	if err != nil {
		objLog.Errorln("comment logic FindObjectComments Error:", err)
	}

	floors := make([]interface{}, 0, len(commentList))
	for _, comment := range commentList {
		self.decodeCmtContentForShow(ctx, comment, true)

		if comment.ReplyFloor > 0 {
			floors = append(floors, comment.ReplyFloor)
		}
	}

	if len(floors) > 0 {
		replyComments = make([]*model.Comment, 0)
		err = MasterDB.Where("objid=? AND objtype=?", objid, objtype).In("floor", floors...).Find(&replyComments)
	}

	return
}

// FindComment 获得评论和额外两个评论
func (self CommentLogic) FindComment(ctx context.Context, cid, objid, objtype int) (*model.Comment, []*model.Comment) {
	objLog := GetLogger(ctx)

	comment := &model.Comment{}
	_, err := MasterDB.Where("cid=?", cid).Get(comment)
	if err != nil {
		objLog.Errorln("CommentLogic FindComment error:", err)
		return comment, nil
	}
	self.decodeCmtContentForShow(ctx, comment, false)

	comments := make([]*model.Comment, 0)
	err = MasterDB.Where("objid=? AND objtype=? AND cid!=?", objid, objtype, cid).
		Limit(2).Find(&comments)
	if err != nil {
		objLog.Errorln("CommentLogic FindComment Find more error:", err)
		return comment, nil
	}
	for _, cmt := range comments {
		self.decodeCmtContentForShow(ctx, cmt, false)
	}

	return comment, comments
}

// Total 评论总数(objtypes[0] 取某一类型的评论总数)
func (CommentLogic) Total(objtypes ...int) int64 {
	var (
		total int64
		err   error
	)
	if len(objtypes) > 0 {
		total, err = MasterDB.Where("objtype=?", objtypes[0]).Count(new(model.Comment))
	} else {
		total, err = MasterDB.Count(new(model.Comment))
	}
	if err != nil {
		logger.Errorln("CommentLogic Total error:", err)
	}
	return total
}

// FindRecent 获得最近的评论
// 如果 uid!=0，表示获取某人的评论；
// 如果 objtype!=-1，表示获取某类型的评论；
func (self CommentLogic) FindRecent(ctx context.Context, uid, objtype, limit int) []*model.Comment {
	dbSession := MasterDB.OrderBy("cid DESC").Limit(limit)

	if uid != 0 {
		dbSession.And("uid=?", uid)
	}
	if objtype != -1 {
		dbSession.And("objtype=?", objtype)
	}

	comments := make([]*model.Comment, 0)
	err := dbSession.Find(&comments)
	if err != nil {
		logger.Errorln("CommentLogic FindRecent error:", err)
		return nil
	}

	cmtMap := make(map[int][]*model.Comment, len(model.PathUrlMap))
	for _, comment := range comments {
		self.decodeCmtContent(ctx, comment)

		if _, ok := cmtMap[comment.Objtype]; !ok {
			cmtMap[comment.Objtype] = make([]*model.Comment, 0, 10)
		}

		cmtMap[comment.Objtype] = append(cmtMap[comment.Objtype], comment)
	}

	cmtObjs := []CommentObjecter{
		model.TypeTopic:     TopicComment{},
		model.TypeArticle:   ArticleComment{},
		model.TypeResource:  ResourceComment{},
		model.TypeWiki:      nil,
		model.TypeProject:   ProjectComment{},
		model.TypeBook:      BookComment{},
		model.TypeInterview: InterviewComment{},
	}
	for cmtType, cmts := range cmtMap {
		self.fillObjinfos(cmts, cmtObjs[cmtType])
	}

	return self.filterDelObjectCmt(comments)
}

// Publish 发表评论（或回复）。
// objid 注册的评论对象
// uid 评论人
func (self CommentLogic) Publish(ctx context.Context, uid, objid int, form url.Values) (*model.Comment, error) {
	objLog := GetLogger(ctx)

	objtype := goutils.MustInt(form.Get("objtype"))
	comment := &model.Comment{
		Objid:   objid,
		Objtype: objtype,
		Uid:     uid,
		Content: form.Get("content"),
	}

	// 使用 Redis 原子计数器获取下一个楼层号
	var err error
	comment.Floor, err = GetNextCommentFloor(objid, objtype)
	if err != nil {
		objLog.Errorln("post comment get next floor error:", err)
		return nil, err
	}

	// 检查是否重复提交（同一用户、同一内容）
	tmpCmt := &model.Comment{}
	_, err = MasterDB.Where("objid=? AND objtype=? AND uid=? AND content=?",
		objid, objtype, uid, comment.Content).Get(tmpCmt)
	if err == nil && tmpCmt.Cid > 0 {
		objLog.Infof("had post comment: %+v", *comment)
		return tmpCmt, nil
	}

	// 入评论库
	_, err = MasterDB.Insert(comment)
	if err != nil {
		objLog.Errorln("post comment service error:", err)
		return nil, err
	}
	self.decodeCmtContentForShow(ctx, comment, true)

	// 回调，不关心处理结果（有些对象可能不需要回调）
	// 使用 errgroup 管理 goroutine，非阻塞执行
	g, gCtx := errgroup.WithContext(context.Background())

	if commenter, ok := commenters[objtype]; ok {
		now := time.Now()

		objLog.Debugf("评论[objid:%d] [objtype:%d] [uid:%d] 成功，通知被评论者更新", objid, objtype, uid)
		// Capture values for closure
		cid := comment.Cid
		capturedCommenter := commenter
		g.Go(func() error {
			select {
			case <-gCtx.Done():
				return gCtx.Err()
			default:
				capturedCommenter.UpdateComment(cid, objid, uid, now)
				return nil
			}
		})

		DefaultFeed.updateComment(objid, objtype, uid, now)
	}

	// NotifyObservers
	cid := comment.Cid
	g.Go(func() error {
		select {
		case <-gCtx.Done():
			return gCtx.Err()
		default:
			commentObservable.NotifyObservers(uid, objtype, cid)
			return nil
		}
	})

	// sendSystemMsg
	capturedUid := uid
	capturedObjid := objid
	capturedObjtype := objtype
	capturedCid := comment.Cid
	capturedForm := form
	g.Go(func() error {
		select {
		case <-gCtx.Done():
			return gCtx.Err()
		default:
			self.sendSystemMsg(ctx, capturedUid, capturedObjid, capturedObjtype, capturedCid, capturedForm)
			return nil
		}
	})

	// Non-blocking wait for goroutines
	go func() {
		if err := g.Wait(); err != nil {
			objLog.Infof("Publish: async tasks error: %v", err)
		}
	}()

	return comment, nil
}

func (CommentLogic) sendSystemMsg(ctx context.Context, uid, objid, objtype, cid int, form url.Values) {
	// 给被评论对象所有者发系统消息 TODO: ext 考虑结构化
	ext := map[string]interface{}{
		"objid":   objid,
		"objtype": objtype,
		"cid":     cid,
		"uid":     uid,
	}

	to := 0
	switch objtype {
	case model.TypeTopic:
		to = DefaultTopic.getOwner(objid)
	case model.TypeArticle:
		to = DefaultArticle.getOwner(objid)
	case model.TypeResource:
		to = DefaultResource.getOwner(objid)
	case model.TypeWiki:
		to = DefaultWiki.getOwner(objid)
	case model.TypeProject:
		to = DefaultProject.getOwner(ctx, objid)
	}

	DefaultMessage.SendSystemMsgTo(ctx, to, objtype, ext)

	// @某人 发系统消息
	DefaultMessage.SendSysMsgAtUids(ctx, form.Get("uid"), ext, to)
	DefaultMessage.SendSysMsgAtUsernames(ctx, form.Get("usernames"), ext, to)
}

// Modify 修改评论信息
func (CommentLogic) Modify(ctx context.Context, cid int, content string) (errMsg string, err error) {
	objLog := GetLogger(ctx)

	_, err = MasterDB.Table(new(model.Comment)).ID(cid).Update(map[string]interface{}{"content": content})
	if err != nil {
		objLog.Errorf("更新评论内容 【%d】 失败：%s", cid, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	return
}

// fillObjinfos 填充评论对应的主体信息
func (CommentLogic) fillObjinfos(comments []*model.Comment, cmtObj CommentObjecter) {
	if len(comments) == 0 {
		return
	}
	count := len(comments)
	commentMap := make(map[int][]*model.Comment, count)
	idSet := set.New(set.NonThreadSafe)
	for _, comment := range comments {
		if _, ok := commentMap[comment.Objid]; !ok {
			commentMap[comment.Objid] = make([]*model.Comment, 0, count)
		}
		commentMap[comment.Objid] = append(commentMap[comment.Objid], comment)
		idSet.Add(comment.Objid)
	}
	cmtObj.SetObjinfo(set.IntSlice(idSet), commentMap)
}

// 提供给其他service调用（包内）
func (CommentLogic) findByIds(cids []int) map[int]*model.Comment {
	if len(cids) == 0 {
		return nil
	}

	comments := make(map[int]*model.Comment)
	err := MasterDB.In("cid", cids).Find(&comments)
	if err != nil {
		return nil
	}

	return comments
}

func (CommentLogic) FindById(cid int) (*model.Comment, error) {
	comment := &model.Comment{}
	_, err := MasterDB.Where("cid=?", cid).Get(comment)
	if err != nil {
		logger.Errorln("CommentLogic findById error:", err)
	}

	return comment, err
}

func (CommentLogic) decodeCmtContent(ctx context.Context, comment *model.Comment) string {
	// 安全过滤
	content := template.HTMLEscapeString(comment.Content)
	// @别人
	content = parseAtUser(ctx, content)

	// 回复某一楼层
	reg := regexp.MustCompile(`#(\d+)楼`)
	url := fmt.Sprintf("%s%d#comment", model.PathUrlMap[comment.Objtype], comment.Objid)
	content = reg.ReplaceAllString(content, `<a href="`+url+`$1" title="$1">#$1<span>楼</span></a>`)

	comment.Content = content

	return content
}

// decodeCmtContentForShow 采用引用的方式显示对其他楼层的回复
func (CommentLogic) decodeCmtContentForShow(ctx context.Context, comment *model.Comment, isEscape bool) {
	// 安全过滤
	content := template.HTMLEscapeString(comment.Content)

	// 回复某一楼层
	reg := regexp.MustCompile(`#(\d+)楼 @([a-zA-Z0-9_-]+)`)
	matches := reg.FindStringSubmatch(content)
	if len(matches) > 2 {
		comment.ReplyFloor = goutils.MustInt(matches[1])
		content = strings.TrimSpace(content[len(matches[0]):])
	}

	// @别人
	content = parseAtUser(ctx, content)

	comment.Content = content
}

// 填充 Comment 对象的 Objinfo 成员接口
// 评论属主应该实现该接口（以便填充 Objinfo 成员）
type CommentObjecter interface {
	// ids 是属主的主键 slice （comment 中的 objid）
	// commentMap 中的 key 是属主 id
	SetObjinfo(ids []int, commentMap map[int][]*model.Comment)
}

// FindAll 支持多页翻看
func (self CommentLogic) FindAll(ctx context.Context, paginator *Paginator, orderBy string, querystring string, args ...interface{}) []*model.Comment {
	objLog := GetLogger(ctx)

	comments := make([]*model.Comment, 0)
	session := MasterDB.OrderBy(orderBy)
	if querystring != "" {
		session.Where(querystring, args...)
	}
	err := session.Limit(paginator.PerPage(), paginator.Offset()).Find(&comments)
	if err != nil {
		objLog.Errorln("CommentLogical FindAll error:", err)
		return nil
	}

	cmtMap := make(map[int][]*model.Comment, len(model.PathUrlMap))
	for _, comment := range comments {
		self.decodeCmtContent(ctx, comment)
		if _, ok := cmtMap[comment.Objtype]; !ok {
			cmtMap[comment.Objtype] = make([]*model.Comment, 0, 10)
		}

		cmtMap[comment.Objtype] = append(cmtMap[comment.Objtype], comment)
	}

	cmtObjs := []CommentObjecter{
		model.TypeTopic:     TopicComment{},
		model.TypeArticle:   ArticleComment{},
		model.TypeResource:  ResourceComment{},
		model.TypeWiki:      nil,
		model.TypeProject:   ProjectComment{},
		model.TypeBook:      BookComment{},
		model.TypeInterview: InterviewComment{},
	}
	for cmtType, cmts := range cmtMap {
		self.fillObjinfos(cmts, cmtObjs[cmtType])
	}

	return self.filterDelObjectCmt(comments)
}

// Count 获取用户全部评论数
func (CommentLogic) Count(ctx context.Context, querystring string, args ...interface{}) int64 {
	objLog := GetLogger(ctx)

	var (
		total int64
		err   error
	)
	if querystring == "" {
		total, err = MasterDB.Count(new(model.Comment))
	} else {
		total, err = MasterDB.Where(querystring, args...).Count(new(model.Comment))
	}

	if err != nil {
		objLog.Errorln("CommentLogic Count error:", err)
	}

	return total
}

func (CommentLogic) filterDelObjectCmt(comments []*model.Comment) []*model.Comment {
	resultCmts := make([]*model.Comment, 0, len(comments))
	for _, comment := range comments {
		if comment.Objinfo != nil && len(comment.Objinfo) > 0 {
			resultCmts = append(resultCmts, comment)
		}
	}
	return resultCmts
}

// 回复赞（喜欢）
type CommentLike struct{}

// 更新该回复的赞
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self CommentLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Where("id=?", objid).Incr("likenum", num).Update(new(model.Comment))
	if err != nil {
		logger.Errorln("更新回复喜欢数失败：", err)
	}
}

func (self CommentLike) String() string {
	return "comment"
}
